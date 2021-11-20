package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gosuri/uitable"
	"manyface.net/internal/messenger"
)

func ListFaces() {
	req, err := http.NewRequest(urls["GetFaces"][0], *s+urls["GetFaces"][1], nil)
	req.Header.Add("session-id", loginResp.SessID)
	if err != nil {
		panic(err)
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		panic(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	faces := []messenger.Face{}
	err = json.Unmarshal(respBody, &faces)
	if err != nil {
		panic(err)
	}

	table := uitable.New()
	table.MaxColWidth = 80
	// table.Wrap = true
	table.AddRow("#", "ID", "Name", "Description")
	for i, f := range faces {
		table.AddRow(strconv.Itoa(i+1), cyan(f.ID), cyan(f.Name), cyan(f.Description))
		// table.AddRow("ID:", cyan(f.ID))
		// table.AddRow("Name:", cyan(f.Name))
		// table.AddRow("Description:", cyan(f.Description))
		// table.AddRow("") // blank
	}
	fmt.Println("---")
	fmt.Println(table)
	fmt.Println("---")
}

func NewFace(name, descr string) {
	reqBody := []byte(fmt.Sprintf(`{"name": "%s","description": "%s"}`, name, descr))

	req, err := http.NewRequest(urls["CreateFace"][0], *s+urls["CreateFace"][1], bytes.NewBuffer(reqBody))
	req.Header.Add("session-id", loginResp.SessID)
	if err != nil {
		panic(err)
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	if resp.StatusCode != http.StatusOK {
		fmt.Println(red("The face wasn't created, status code "), red(strconv.Itoa(resp.StatusCode)))
	} else {
		fmt.Println(green("New face was created"))
	}
	fmt.Println("---")
}

/*
func CreateConn(faceID, matrixPeerUserID string) {
	var faceName, matrixUserID, matrixPassword, matrixAccessToken, matrixRoomID string

	err := db.QueryRow("SELECT name FROM face WHERE face_id = ?", faceID).Scan(&faceName)
	if err != nil {
		fmt.Println(red("Wrong face id"))
		return
	}
	matrixUserID = RandStringRunes(10)
	matrixPassword = RandStringRunes(10)

	res1, err := client.RegisterDummy(&mautrix.ReqRegister{
		Username: matrixUserID,
		Password: matrixPassword,
	})
	if err != nil {
		panic(err)
	}
	client.SetCredentials(res1.UserID, res1.AccessToken)
	matrixAccessToken = res1.AccessToken
	matrixUserID = res1.UserID.String()

	res2, err := client.CreateRoom(
		&mautrix.ReqCreateRoom{
			Visibility: "private",
			// RoomAliasName: "RoomAliasName",
			Name:     faceName,
			Invite:   []id.UserID{id.UserID(matrixPeerUserID)},
			Preset:   "trusted_private_chat",
			IsDirect: true,
		})
	if err != nil {
		panic(err)
	}
	matrixRoomID = string(res2.RoomID)

	res3, err := db.Exec("INSERT INTO connection (username, password, room_id, peer_id, access_token, face_id) VALUES (?, ?, ?, ?, ?, ?)",
		matrixUserID, matrixPassword, matrixRoomID, matrixPeerUserID, matrixAccessToken, faceID)
	if err != nil {
		panic(err)
	}
	rowCnt, err := res3.RowsAffected()
	if rowCnt != 1 || err != nil {
		fmt.Println(red("Error created new connection"))
		return
	}

	client.Logout()
	client.ClearCredentials()
	fmt.Printf("Created connection with user %s for face %s\n", matrixPeerUserID, faceName)
}

func SendMessage(faceID, matrixPeerUserID, message string) {
	var matrixUserID, matrixPassword, matrixAccessToken, matrixRoomID string

	err := db.
		QueryRow("SELECT username, access_token, room_id, password FROM connection WHERE face_id = ? AND peer_id = ?", faceID, matrixPeerUserID).
		Scan(&matrixUserID, &matrixAccessToken, &matrixRoomID, &matrixPassword)
	if err != nil {
		panic(err)
	}

	_, err = client.Login(&mautrix.ReqLogin{
		Type:             "m.login.password",
		Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: matrixUserID},
		Password:         matrixPassword,
		StoreCredentials: true,
	})
	if err != nil {
		panic(err)
	}
	_, err = client.SendText(id.RoomID(matrixRoomID), message)
	if err != nil {
		panic(err)
	}

	client.Logout()
	client.ClearCredentials()
}
*/
