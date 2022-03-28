package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc/status"

	"github.com/gosuri/uitable"
	"manyface.net/internal/messenger"
)

func ListFaces() {
	req, err := http.NewRequest(urls["GetFaces"][0], *fRest+urls["GetFaces"][1], nil)
	req.Header.Add("session-id", loginResp.Body.User.SessID)
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

	var faces GetFacesResp

	err = json.Unmarshal(respBody, &faces)
	if err != nil {
		panic(err)
	}

	table := uitable.New()
	table.MaxColWidth = 80
	// table.Wrap = true
	table.AddRow("#", "ID", "Nick", "Purpose")
	for i, f := range faces.Body.Faces {
		table.AddRow(strconv.Itoa(i+1), yellow(f.FaceID), cyan(f.Nick), f.Purpose)
	}
	fmt.Println("---")
	fmt.Println(table)
	fmt.Println("---")

}

func NewFace(nick, purpose, bio, comments, server string) {
	reqBody := []byte(fmt.Sprintf(`{"nick": "%s","purpose": "%s","bio": "%s","comments": "%s","server": "%s"}`, nick, purpose, bio, comments, server))

	req, err := http.NewRequest(urls["CreateFace"][0], *fRest+urls["CreateFace"][1], bytes.NewBuffer(reqBody))
	req.Header.Add("session-id", loginResp.Body.User.SessID)
	if err != nil {
		panic(err)
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	if resp.StatusCode != http.StatusCreated {
		fmt.Println(red("The face wasn't created, status code "), red(strconv.Itoa(resp.StatusCode)))
	} else {
		fmt.Println(green("The new face was created"))
	}
	fmt.Println("---")
}

func CreateConn(faceID, peerFaceID string) {
	reqBody := []byte(fmt.Sprintf(`{"face_user_id": "%s","face_peer_id": "%s"}`, faceID, peerFaceID))

	req, err := http.NewRequest(urls["CreateConn"][0], *fRest+urls["CreateConn"][1], bytes.NewBuffer(reqBody))
	req.Header.Add("session-id", loginResp.Body.User.SessID)
	if err != nil {
		panic(err)
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("---")
	if resp.StatusCode != http.StatusCreated {
		fmt.Println(red("The connection wasn't created, status code "), red(strconv.Itoa(resp.StatusCode)))
	} else {
		fmt.Println(green("The new connection was created"))
	}
	fmt.Println("---")
}

func ListConns() {
	req, err := http.NewRequest(urls["GetConns"][0], *fRest+urls["GetConns"][1], nil)
	req.Header.Add("session-id", loginResp.Body.User.SessID)
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

	var conns GetConnsResp
	err = json.Unmarshal(respBody, &conns)
	if err != nil {
		panic(err)
	}

	table := uitable.New()
	table.MaxColWidth = 80
	// table.Wrap = true
	table.AddRow("ID", "My Face", "Face Peer")
	for _, c := range conns.Body.Connections {
		faceUser := getFace(c.FaceUserID).Body.Face
		// fmt.Println("faceUser=", &faceUser)
		facePeer := getFace(c.FacePeerID).Body.Face
		// fmt.Println("facePeer=", facePeer)
		table.AddRow(cyan(strconv.Itoa(int(c.ConnID))), cyan(faceUser.Nick)+" ("+yellow(c.FaceUserID)+")", cyan(facePeer.Nick)+" ("+yellow(c.FacePeerID)+")")
	}
	fmt.Println("---")
	fmt.Println(table)
	fmt.Println("---")
}

func SendMsg(connID int64, message string) {
	request := &messenger.SendRequest{
		Message:      message,
		ConnectionId: connID,
	}
	response, err := grpcCli.Send(context.Background(), request)
	if err != nil {
		errStatus, ok := status.FromError(err)
		if ok {
			fmt.Println(red(fmt.Sprintf("%v, %v", errStatus.Message(), errStatus.Code())))
		}
	}
	fmt.Println(green(fmt.Sprintf("message sent to %v", response.ConnectionId)))
}

func ListenMsg() {
	req, err := http.NewRequest(urls["GetConns"][0], *fRest+urls["GetConns"][1], nil)
	req.Header.Add("session-id", loginResp.Body.User.SessID)
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

	var conns GetConnsResp
	err = json.Unmarshal(respBody, &conns)
	if err != nil {
		panic(err)
	}

	for _, c := range conns.Body.Connections {
		go func(connID int64) {
			request := messenger.ListenRequest{
				ConnectionId: connID,
			}
			// fmt.Printf("%+v\n", request)
			stream, err := grpcCli.Listen(globalCtx, &request)
			if err != nil {
				errStatus, ok := status.FromError(err)
				if ok {
					fmt.Println(red(fmt.Sprintf("%v, %v", errStatus.Message(), errStatus.Code())))
				}
			}
			go func() {
				<-globalCtx.Done()
				fmt.Printf("clean up %v\n", connID)
			}()
			for {

				// fmt.Println("starting listening...")
				msg, err := stream.Recv()
				// fmt.Println(msg)
				if err == io.EOF {
					return
				}
				if err != nil {
					errStatus, ok := status.FromError(err)
					if ok {
						fmt.Println(red(fmt.Sprintf("%v, %v", errStatus.Message(), errStatus.Code())))
					}
				}
				table := uitable.New()
				table.MaxColWidth = 80
				table.AddRow("From", "To", "Message", "Time")
				fS := getFace(msg.SenderFaceId).Body.Face
				fR := getFace(msg.ReceiverFaceId).Body.Face
				// fmt.Println(f.Nick)
				table.AddRow(cyan(fS.Nick)+" ("+yellow(msg.SenderFaceId)+")", cyan(fR.Nick)+" ("+yellow(msg.ReceiverFaceId)+")", green(msg.Message), time.Unix(msg.Timestamp, 0))
				fmt.Println(table)
			}
		}(c.ConnID)
	}
}

func getFace(faceID string) *GetFaceResp {
	req, err := http.NewRequest(urls["GetFace"][0], *fRest+urls["GetFace"][1]+faceID, nil)
	req.Header.Add("session-id", loginResp.Body.User.SessID)
	if err != nil {
		return nil
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		return nil
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var f GetFaceResp
	json.Unmarshal(respBody, &f)
	return &f
}
