package messenger

import (
	"encoding/json"
	"image/png"
	"io/ioutil"
	"net/http"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/julienschmidt/httprouter"
	"manyface.net/internal/utils"
)

// POST /api/face
func (h *MessengerHandler) CreateFace(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't read body", h.Logger)
		return
	}

	face := &Face{UserID: sess.UserID}
	if err = json.Unmarshal(body, face); err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't unmarshal body json", h.Logger)
		return
	}
	face.ID, err = h.Srv.CreateFace(face.Nick, face.Purpose, face.Bio, face.Comments, face.Server, face.UserID)
	if err != nil || face.ID == "" {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't create face", h.Logger)
		return
	}

	utils.RespJSON(w, http.StatusCreated, map[string]interface{}{
		"face": &face,
	})

	h.Logger.Infof("The face %v - %v was created", face.ID, face.Nick)
}

// GET /api/face/{FACE_ID}
func (h *MessengerHandler) GetFace(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())
	face, err := h.Srv.GetFaceByID(ps.ByName("FACE_ID"), sess.UserID)
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't get face", h.Logger)
		return
	}

	utils.RespJSON(w, http.StatusOK, map[string]interface{}{
		"face": &face,
	})

	h.Logger.Infof("Got the face %v ", face.ID)
}

// GET /api/faces
func (h *MessengerHandler) GetFaces(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	userID := sess.UserID
	faces, err := h.Srv.GetFacesByUser(userID)
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't get faces for user "+userID, h.Logger)
		return
	}

	utils.RespJSON(w, http.StatusOK, map[string]interface{}{
		"faces": &faces,
	})

	h.Logger.Infof("Got the faces for user %v ", userID)
}

// DELETE /api/face/{FACE_ID}
func (h *MessengerHandler) DelFace(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())
	faceID := ps.ByName("FACE_ID")
	if err := h.Srv.DelFaceByID(faceID, sess.UserID); err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't delete face "+faceID, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.Logger.Infof("The face %v was deleted", faceID)
}

func (h *MessengerHandler) UpdFace(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// r.ParseMultipartForm(5 * 1024 * 1025) // TODO: save or remove?
	sess, _ := h.SM.GetFromCtx(r.Context())
	faceID := ps.ByName("FACE_ID")

	body := r.FormValue("face")
	face := &Face{ID: faceID, UserID: sess.UserID}
	if err := json.Unmarshal([]byte(body), face); err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't unmarshal face json", h.Logger)
		return
	}

	avatar, _, err := r.FormFile("avatar")
	if err != nil {
		utils.RespJSONError(w, http.StatusBadRequest, err, "Can't read avatar", h.Logger)
		return
	}
	defer avatar.Close()

	err = h.Srv.UpdFace(face)
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't update face", h.Logger)
		return
	}
	err = h.BS.Put(avatar, face.ID+".png")
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't save avatar", h.Logger)
		return
	}

	utils.RespJSON(w, http.StatusOK, map[string]interface{}{
		"face": &face,
	})
	h.Logger.Infof("The face %v - %v was update", face.ID, face.Nick)
}

func (h *MessengerHandler) GetFaceQR(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	faceID := ps.ByName("FACE_ID")

	qrCode, err := qr.Encode(faceID, qr.L, qr.Auto)
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't get qr for face "+faceID, h.Logger)
		return
	}
	qrCode, err = barcode.Scale(qrCode, 512, 512)
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't get qr for face "+faceID, h.Logger)
		return
	}
	w.WriteHeader(http.StatusOK)
	png.Encode(w, qrCode)
}

func (h *MessengerHandler) GetFaceAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	faceID := ps.ByName("FACE_ID")

	avatar, err := h.BS.Get(faceID + ".png")
	if err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't get avatar for face "+faceID, h.Logger)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Write(avatar)
}

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
