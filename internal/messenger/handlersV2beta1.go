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
func (h *MessengerHandler) CreateFaceV2beta1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	face.ID, err = h.Srv.CreateFaceV2beta1(face.Nick, face.Purpose, face.Bio, face.Comments, face.UserID)
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
func (h *MessengerHandler) GetFaceV2beta1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())
	face, err := h.Srv.GetFaceByIDV2beta1(ps.ByName("FACE_ID"), sess.UserID)
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
func (h *MessengerHandler) GetFacesV2beta1(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())

	userID := sess.UserID
	faces, err := h.Srv.GetFacesByUserV2beta1(userID)
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
func (h *MessengerHandler) DelFaceV2beta1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sess, _ := h.SM.GetFromCtx(r.Context())
	faceID := ps.ByName("FACE_ID")
	if err := h.Srv.DelFaceByIDV2beta1(faceID, sess.UserID); err != nil {
		utils.RespJSONError(w, http.StatusInternalServerError, err, "Can't delete face "+faceID, h.Logger)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	h.Logger.Infof("The face %v was deleted", faceID)
}

func (h *MessengerHandler) UpdFaceV2beta1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	err = h.Srv.UpdFaceV2beta1(face)
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

func (h *MessengerHandler) GetFaceQRV2beta1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (h *MessengerHandler) GetFaceAvatarV2beta1(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
