package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/andrdru/go-template/internal/ctxsess"
)

type (
	UserGetResp struct {
		ID int64 `json:"id"`
	}
)

func (a *API) UserGet(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	m := NewMessage()

	sd := ctxsess.Get(r.Context())

	m.Data = UserGetResp{ID: sd.UserID}

	_ = m.Return(w)
}
