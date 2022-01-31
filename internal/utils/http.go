package utils

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// TODO: remove func
func HandleError(w http.ResponseWriter, err error, status int, msg string, logger *zap.SugaredLogger) {
	var errMsg string
	if err != nil {
		errMsg = msg + " - " + err.Error()
	} else {
		errMsg = msg
	}
	http.Error(w, errMsg, status)
	logger.Error(errMsg)
}

type ServerResponse struct {
	Body  interface{} `json:"body,omitempty"`
	Error string      `json:"error,omitempty"`
}

func RespJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	respJSON, _ := json.Marshal(&ServerResponse{
		Body: body,
	})
	w.Write(respJSON)
}

func RespJSONError(w http.ResponseWriter, status int, err error, text string, logger *zap.SugaredLogger) {
	if err != nil {
		text = text + " - " + err.Error()
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	respJSON, _ := json.Marshal(&ServerResponse{
		Error: text,
	})
	w.Write(respJSON)
	logger.Error(text)
}
