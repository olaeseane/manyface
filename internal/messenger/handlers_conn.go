package messenger

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"manyface.net/internal/utils"
)

// POST /api/conn
func (h *MessengerHandler) CreateConn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't read body", h.Logger)
		return
	}

	conn := &Conn{}
	if err = json.Unmarshal(body, conn); err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't unmarshal body json", h.Logger)
		return
	}

	connections, err := h.Srv.CreateConn(sess.UserID, conn.FaceUserID, conn.FacePeerID)
	if err != nil || len(connections) != 2 {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't create connection", h.Logger)
		return
	}

	utils.RespJSON(w, http.StatusCreated, map[string]interface{}{
		"сonnections": connections,
	})
	h.Logger.Infof("The connections %v and %v was created", connections[0].ID, connections[1].ID)
}

// DELETE /api/conn
func (h *MessengerHandler) DeleteConn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't read body", h.Logger)
		return
	}

	conn := &Conn{}
	if err = json.Unmarshal(body, conn); err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't unmarshal body json", h.Logger)
		return
	}

	err = h.Srv.DeleteConn(sess.UserID, conn.FaceUserID, conn.FacePeerID)
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't delete connection", h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.Logger.Infof("The connections were deleted")
}

// GET /api/conns
func (h *MessengerHandler) GetConns(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	u := sess.UserID
	conns, err := h.Srv.GetConnsByUser(u)
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't get connections for user "+u, h.Logger)
		return
	}

	utils.RespJSON(w, http.StatusOK, map[string]interface{}{
		"сonnections": conns,
	})

	h.Logger.Infof("Got the connections for user %v ", u)
}
