package main

type LoginResp struct {
	Body struct {
		User struct {
			UserID string `json:"user_id"`
			SessID string `json:"sess_id"`
		} `json:"user"`
	} `json:"body"`
}

type GetFacesResp struct {
	Body struct {
		Faces []struct {
			FaceID   string `json:"face_id"`
			Nick     string `json:"nick"`
			Purpose  string `json:"purpose"`
			Bio      string `json:"bio"`
			Comments string `json:"comments"`
			Server   string `json:"server"`
			UserID   string `json:"user_id"`
		} `json:"faces"`
	} `json:"body"`
}

type GetFaceResp struct {
	Body struct {
		Face struct {
			FaceID   string `json:"face_id"`
			Nick     string `json:"nick"`
			Purpose  string `json:"purpose"`
			Bio      string `json:"bio"`
			Comments string `json:"comments"`
			Server   string `json:"server"`
			UserID   string `json:"user_id"`
		} `json:"face"`
	} `json:"body"`
}

type GetConnsResp struct {
	Body struct {
		Connections []struct {
			ConnID     int64  `json:"conn_id"`
			FaceUserID string `json:"face_user_id"`
			FacePeerID string `json:"face_peer_id"`
		} `json:"сonnections"`
	} `json:"body"`
}
