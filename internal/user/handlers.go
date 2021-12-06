package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"manyface.net/internal/utils"
)

// POST /api/reg
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't read body", h.Logger)
		return
	}

	u := &User{}
	if err = json.Unmarshal(body, u); err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't unmarshal body json", h.Logger)
		return
	}
	if u.ID, err = h.Repo.Register(u.Username, u.Password); err != nil || u.ID == -1 {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't register user", h.Logger)
		return
	}

	var sessID string
	if sessID, err = h.SM.Create(u.ID); err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't create session", h.Logger)
	}

	w.WriteHeader(http.StatusOK)
	resp := struct {
		UserID int64  `json:"user_id"`
		SessID string `json:"sess_id"`
	}{UserID: u.ID, SessID: sessID}
	b, _ := json.Marshal(resp)
	w.Write(b)

	h.Logger.Infof("The user %v has registered", u.Username)
}

// POST /api/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't read body", h.Logger)
		return
	}

	u := &User{}
	if err = json.Unmarshal(body, u); err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't unmarshal body json", h.Logger)
		return
	}
	if u.ID, err = h.Repo.Login(u.Username, u.Password); err != nil || u.ID == -1 {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't login", h.Logger)
		return
	}

	var sessID string
	if sessID, err = h.SM.Create(u.ID); err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't create session", h.Logger)
	}

	w.WriteHeader(http.StatusOK)
	resp := struct {
		UserID int64  `json:"user_id"`
		SessID string `json:"sess_id"`
	}{UserID: u.ID, SessID: sessID}
	b, _ := json.Marshal(resp)
	w.Write(b)

	h.Logger.Infof("The user %v has logged in", u.Username)
}
