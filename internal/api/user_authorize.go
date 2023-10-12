package api

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/andrdru/go-template/internal/entities"
	"github.com/julienschmidt/httprouter"
)

//go:generate easyjson

type (
	//easyjson:json
	UserAuthorizeReq struct {
		Email string `json:"email"`
		Pass  string `json:"pass"`
	}
)

const (
	// HeaderIP .
	HeaderIP = "X-Real-IP"
	// HeaderUserAgent .
	HeaderUserAgent = "User-Agent"
)

func (a *API) UserAuthorize(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	message := NewMessage()

	req := &UserAuthorizeReq{}
	err := ReadRequest(r.Body, req)
	if err != nil {
		message.SetError(Error(err.Error()), Code(http.StatusBadRequest))
		_ = message.Return(w)
		return
	}

	if !req.Validate(message) {
		message.SetError(Code(http.StatusBadRequest))
		_ = message.Return(w)
		return
	}

	session := entities.Session{
		Extra: entities.SessionExtra{
			IP:        r.Header.Get(HeaderIP),
			UserAgent: r.Header.Get(HeaderUserAgent),
		},
		Email: req.Email,
		Pass:  req.Pass,
	}

	err = a.authManager.Login(r.Context(), w, session)
	if err != nil {
		if !errors.Is(err, entities.ErrNotFound) {
			message.SetError(Code(http.StatusForbidden))
			_ = message.Return(w)
			return
		}

		a.logger.Error("login", slog.Any("error", err))

		message.SetError(OptInternalError)
		_ = message.Return(w)
		return
	}

	_ = message.Return(w)
}

func (v *UserAuthorizeReq) Validate(message *Message) (ok bool) {
	ok = true
	if v.Email == "" {
		ok = false
		message.SetError(MapError("email", "should not be empty"))
	}

	if v.Pass == "" {
		ok = false
		message.SetError(MapError("pass", "should not be empty"))
	}

	return ok
}
