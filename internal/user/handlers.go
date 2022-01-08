package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"manyface.net/internal/utils"
)

/*
// POST /api/reg
func (h *UserHandler) RegisterV1beta1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	if u.ID, err = h.Repo.RegisterV1beta1(u.Username, u.Password); err != nil || u.ID == -1 {
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

	h.Logger.Infof("The user %v has registered", u.ID)
}

// POST /api/login
func (h *UserHandler) LoginV1beta1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	if u.ID, err = h.Repo.LoginV1beta1(u.Username, u.Password); err != nil || u.ID == -1 {
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
*/

// POST /api/user
func (h *UserHandler) RegisterV2beta1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {

		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't read body", h.Logger)
		return
	}

	u := &User{}
	if err = json.Unmarshal(body, u); err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't unmarshal body json", h.Logger)
		return
	}
	if u.Password == "" {
		utils.RespJSONError(w, http.StatusBadRequest, nil, "Password empty", h.Logger)
		return
	}
	var mnemonic []string
	if mnemonic, u.ID, err = h.Repo.RegisterV2beta1(u.Password); err != nil || u.ID == -1 {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't register user", h.Logger)
		return
	}

	var sessID string
	if sessID, err = h.SM.Create(u.ID); err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't create session", h.Logger)
	}

	utils.RespJSON(w, http.StatusCreated, map[string]interface{}{
		"user": struct {
			UserID   int64    `json:"user_id"`
			SessID   string   `json:"sess_id"`
			Mnemonic []string `json:"mnemonic"`
		}{UserID: u.ID, SessID: sessID, Mnemonic: mnemonic},
	})

	h.Logger.Infof("The user %v has registered", u.Username)
}

// GET /api/user
func (h *UserHandler) LoginV2beta1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't read body", h.Logger)
		return
	}

	u := &User{}
	if err = json.Unmarshal(body, u); err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't unmarshal body json", h.Logger)
		return
	}
	if u.ID, err = h.Repo.LoginV2beta1(u.ID, u.Password, u.Mnemonic); err != nil || u.ID == -1 {
		utils.RespJSONError(w, http.StatusUnauthorized, err, "Can't login", h.Logger)
		return
	}

	var sessID string
	if sessID, err = h.SM.Create(u.ID); err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't create session", h.Logger)
	}

	utils.RespJSON(w, http.StatusOK, map[string]interface{}{
		"user": struct {
			UserID int64  `json:"user_id"`
			SessID string `json:"sess_id"`
		}{UserID: u.ID, SessID: sessID},
	})

	h.Logger.Infof("The user %v has logged in", u.Username)
}

// GET /api/user
func (h *UserHandler) LoginV3beta1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uid, p, ok := r.BasicAuth()
	if !ok {
		utils.RespJSONError(w, http.StatusUnauthorized, nil, "Error parsing basic auth", h.Logger)
		return
	}
	userID, err := strconv.Atoi(uid)
	if err != nil {
		utils.RespJSONError(w, http.StatusUnauthorized, nil, "Wrong user id", h.Logger)
		return
	}

	if err := h.Repo.LoginV3beta1(int64(userID), p); err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't login", h.Logger)
		return
	}

	var sessID string
	if sessID, err = h.SM.Create(int64(userID)); err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't create session", h.Logger)
	}

	utils.RespJSON(w, http.StatusOK, map[string]interface{}{
		"user": struct {
			UserID int64  `json:"user_id"`
			SessID string `json:"sess_id"`
		}{UserID: int64(userID), SessID: sessID},
	})

	h.Logger.Infof("The user %v has logged in", userID)
}
