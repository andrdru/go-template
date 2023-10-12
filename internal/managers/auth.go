package managers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/andrdru/go-template/internal/ctxsess"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/andrdru/go-template/internal/entities"
	"github.com/andrdru/go-template/internal/middlewares"
	"github.com/andrdru/go-template/internal/repos"
)

type Auth struct {
	userRepo *repos.User
}

const (
	headerUserSession = "X-User-Session"

	cookieTokenStoreDuration = 3 * 30 * 24 * time.Hour
)

func NewAuth(userRepo *repos.User) *Auth {
	return &Auth{
		userRepo: userRepo,
	}
}

func (a *Auth) Check(r *http.Request) (ctx context.Context, err error) {
	cookie, err := r.Cookie(headerUserSession)
	if err != nil {
		return nil, fmt.Errorf("get cookie: %w", err)
	}

	session, err := sessionDataFromCookie(cookie)
	if err != nil {
		return nil, fmt.Errorf("sessionDataFromCookie: %w", err)
	}

	session, err = a.getSessionByToken(r.Context(), session.Token)
	if err != nil {
		return nil, fmt.Errorf("getSessionByToken: %w", err)
	}

	return ctxsess.Set(r.Context(), session), nil
}

func (a *Auth) Login(
	ctx context.Context,
	w http.ResponseWriter,
	session entities.Session,
) error {
	getUser, err := a.userRepo.User(ctx, session.Email)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	if !checkPasswordHash(session.Pass, getUser.Passhash) {
		return entities.ErrNotAllowed
	}

	session.UserID = getUser.ID
	session.Token = uuid.NewString()

	err = a.userRepo.CreateSession(ctx, session)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	err = setSessionCookie(w, &session)
	if err != nil {
		return fmt.Errorf("setSessionCookie: %w", err)
	}

	return nil
}

func (a *Auth) getSessionByToken(ctx context.Context, token string) (session *entities.Session, err error) {
	userSession, err := a.userRepo.Session(ctx, token)
	if err != nil {
		if errors.Is(err, entities.ErrNotFound) {
			return nil, middlewares.ErrNotAllowed
		}
		return nil, fmt.Errorf("get session: %w", err)
	}

	return &userSession, nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func sessionDataFromCookie(cookie *http.Cookie) (session *entities.Session, err error) {
	data, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	session = &entities.Session{}
	err = json.Unmarshal(data, session)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return session, nil
}

func setSessionCookie(w http.ResponseWriter, session *entities.Session) error {
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     headerUserSession,
		Value:    base64.URLEncoding.EncodeToString(data),
		Expires:  time.Now().Add(cookieTokenStoreDuration),
		Path:     "/",
		HttpOnly: true,
	})

	return nil
}
