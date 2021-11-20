package main

type LoginResp struct {
	UserID int64  `json:"user_id,omitempty"`
	SessID string `json:"sess_id,omitempty"`
}
