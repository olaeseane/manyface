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

	"github.com/gosuri/uitable"
	"manyface.net/internal/messenger"

	pb "manyface.net/grpc"
)

func ListFaces() {
	req, err := http.NewRequest(urls["GetFaces"][0], *ws+urls["GetFaces"][1], nil)
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
		table.AddRow(strconv.Itoa(i+1), yellow(f.ID), cyan(f.Name), f.Description)
	}
	fmt.Println("---")
	fmt.Println(table)
	fmt.Println("---")
}

func NewFace(name, descr string) {
	reqBody := []byte(fmt.Sprintf(`{"name": "%s","description": "%s"}`, name, descr))

	req, err := http.NewRequest(urls["CreateFace"][0], *ws+urls["CreateFace"][1], bytes.NewBuffer(reqBody))
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
		fmt.Println(green("The new face was created"))
	}
	fmt.Println("---")
}

func CreateConn(faceID, peerFaceID string) {
	reqBody := []byte(fmt.Sprintf(`{"face_user_id": "%s","face_peer_id": "%s"}`, faceID, peerFaceID))

	req, err := http.NewRequest(urls["CreateConn"][0], *ws+urls["CreateConn"][1], bytes.NewBuffer(reqBody))
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
		fmt.Println(red("The connection wasn't created, status code "), red(strconv.Itoa(resp.StatusCode)))
	} else {
		fmt.Println(green("The new connection was created"))
	}
	fmt.Println("---")
}

func ListConns() {
	req, err := http.NewRequest(urls["GetConns"][0], *ws+urls["GetConns"][1], nil)
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

	conns := []messenger.Conn{}
	err = json.Unmarshal(respBody, &conns)
	if err != nil {
		panic(err)
	}

	table := uitable.New()
	table.MaxColWidth = 80
	// table.Wrap = true
	table.AddRow("ID", "My Face", "Face Peer")
	for _, c := range conns {
		faceUser, _ := getFace(c.FaceUserID)
		facePeer, _ := getFace(c.FacePeerID)
		table.AddRow(cyan(strconv.Itoa(int(c.ID))), cyan(faceUser.Name)+" ("+yellow(c.FaceUserID)+")", cyan(facePeer.Name)+" ("+yellow(c.FacePeerID)+")")
	}
	fmt.Println("---")
	fmt.Println(table)
	fmt.Println("---")
}

func SendMsg(connID int64, message string) {
	request := &pb.SendRequest{
		Message:      message,
		ConnectionId: connID,
	}
	response, err := grpcCli.Send(context.Background(), request)
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Result)
}

func ListenMsg() {
	req, err := http.NewRequest(urls["GetConns"][0], *ws+urls["GetConns"][1], nil)
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

	conns := []messenger.Conn{}
	err = json.Unmarshal(respBody, &conns)
	if err != nil {
		panic(err)
	}

	// fmt.Println(conns)
	for _, c := range conns {
		go func(c messenger.Conn) {
			request := pb.ListenRequest{
				ConnectionId: c.ID,
			}
			// fmt.Printf("%+v\n", request)
			stream, err := grpcCli.Listen(context.Background(), &request)
			if err != nil {
				// fmt.Println(err)
				panic(err)
			}
			for {
				// fmt.Println("starting listening...")
				msg, err := stream.Recv()
				// fmt.Println(msg)
				if err == io.EOF {
					break
				}
				if err != nil {
					panic(err)
				}
				table := uitable.New()
				table.MaxColWidth = 80
				table.AddRow("From (Face)", "From (ID)", "Message", "Time")
				f, _ := getFace(msg.Sender)
				table.AddRow(cyan(f.Name), yellow(msg.Sender), green(msg.Content), time.Unix(msg.Timestamp, 0))
				fmt.Println(table)
			}
		}(c)
	}
}

func getFace(faceID string) (*messenger.Face, error) {
	req, err := http.NewRequest(urls["GetFace"][0], *ws+urls["GetFace"][1]+faceID, nil)
	req.Header.Add("session-id", loginResp.SessID)
	if err != nil {
		return nil, err
	}
	resp, err := httpCli.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var f messenger.Face
	json.Unmarshal(respBody, &f)
	return &f, nil
}
