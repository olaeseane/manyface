package messenger

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"manyface.net/internal/utils"
)

/*
// POST /api/face
func (h *MessengerHandler) CreateFace(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't read body", h.Logger)
		return
	}

	face := &Face{UserID: sess.UserID}
	if err = json.Unmarshal(body, face); err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't unmarshal body json", h.Logger)
		return
	}
	face.ID, err = h.Srv.CreateFace(face.Name, face.Description, face.UserID)
	if err != nil || face.ID == "" {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't create face", h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(face)
	w.Write(b)

	h.Logger.Infof("The face %v - %v was created", face.ID, face.Name)
}

// GET /api/face/{FACE_ID}
func (h *MessengerHandler) GetFace(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	face, err := h.Srv.GetFaceByID(ps.ByName("FACE_ID"))
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't get face", h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(face)
	w.Write(b)

	h.Logger.Infof("Got the face %v ", face.ID)
}

// DELETE /api/face/{FACE_ID}
func (h *MessengerHandler) DelFace(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())
	f := ps.ByName("FACE_ID")
	if err := h.Srv.DelFaceByID(f, sess.UserID); err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't delete face "+f, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	h.Logger.Infof("The face %v was deleted", f)
}

// GET /api/faces
func (h *MessengerHandler) GetFaces(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	u := sess.UserID
	faces, err := h.Srv.GetFacesByUser(u)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't get faces for user "+strconv.Itoa(int(u)), h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(faces)
	w.Write(b)

	h.Logger.Infof("Got the faces for user %v ", u)
}
*/

// POST /api/connection
func (h *MessengerHandler) CreateConn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't read body", h.Logger)
		return
	}

	conn := &Conn{}
	if err = json.Unmarshal(body, conn); err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't unmarshal body json", h.Logger)
		return
	}

	connections, err := h.Srv.CreateConn(sess.UserID, conn.FaceUserID, conn.FacePeerID)
	if err != nil || len(connections) != 2 {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't create connection", h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(connections)
	w.Write(b)

	h.Logger.Infof("The connections %v and %v was created", connections[0].ID, connections[1].ID)
}

// DELETE /api/connection
func (h *MessengerHandler) DeleteConn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't read body", h.Logger)
		return
	}

	conn := &Conn{}
	if err = json.Unmarshal(body, conn); err != nil {
		utils.HandleError(w, err, http.StatusBadRequest, "Can't unmarshal body json", h.Logger)
		return
	}

	err = h.Srv.DeleteConn(sess.UserID, conn.FaceUserID, conn.FacePeerID)
	if err != nil {
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't delete connection", h.Logger)
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
		utils.HandleError(w, err, http.StatusInternalServerError, "Can't get connections for user "+u, h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(conns)
	w.Write(b)

	h.Logger.Infof("Got the connections for user %v ", u)
}
