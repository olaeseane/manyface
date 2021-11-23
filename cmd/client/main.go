package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ttacon/chalk"
	"google.golang.org/grpc"

	pb "manyface.net/grpc"
)

// const grpcAddr string = "127.0.0.1:5300"

// TODO: remove global vars
var (
	ws = flag.String("ws", "http://localhost:8080", "manyface server")
	gs = flag.String("gs", "127.0.0.1:5300", "manyface server")
	u  = flag.String("u", "user3", "manyface user")
	p  = flag.String("p", "welcome", "manyface password")

	httpCli = &http.Client{Timeout: time.Second * 5}

	grpcConn *grpc.ClientConn
	grpcCli  pb.MessengerClient

	loginResp = LoginResp{}

	helpMessage string = `---
#Commands
/nf <name> <description>			- Create a new face with given name and description (name without spaces)
/faces (f)					- List of my faces
/conn <my face id> <peer face id>		- Create a connection with given peer
/conns (cc) 					- List of my connections
/text <conn id> <message>			- Send a message to given peer from the given face
/quit (q)					- Quit app
---
`

	green, cyan, red, yellow, magenta, blue func(string) string

	urls = map[string][]string{
		"Register":   {"POST", "/api/reg"},
		"Login":      {"POST", "/api/login"},
		"CreateFace": {"POST", "/api/face"},
		"GetFace":    {"GET", "/api/face/"},    // :FACE_ID
		"DelFace":    {"DELETE", "/api/face/"}, // :FACE_ID
		"GetFaces":   {"GET", "/api/faces"},
		"CreateConn": {"POST", "/api/conn"},
		"DeleteConn": {"DELETE", "/api/conn"},
		"GetConns":   {"GET", "/api/conns"},
	}
)

func main() {
	flag.Parse()
	if *u == "" || *p == "" || *ws == "" || *gs == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	green = chalk.Green.NewStyle().WithBackground(chalk.Black).Style
	cyan = chalk.Cyan.NewStyle().Style
	red = chalk.Red.NewStyle().Style
	yellow = chalk.Yellow.NewStyle().Style
	blue = chalk.Blue.NewStyle().Style
	magenta = chalk.Magenta.NewStyle().Style

	reqBody := []byte(fmt.Sprintf(`{"username": "%s","password": "%s"}`, *u, *p))
	req, err := http.NewRequest(urls["Login"][0], *ws+urls["Login"][1], bytes.NewBuffer(reqBody))
	if err != nil {
		panic(err)
	}

	resp, err := httpCli.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(red("Can't login"))
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(respBody, &loginResp)

	fmt.Println("Logging into", cyan(*ws), "and", cyan(*gs), "as", cyan(*u))
	fmt.Println(magenta(helpMessage))

	grpcConn = getGrpcConn(*gs)
	defer grpcConn.Close()
	grpcCli = pb.NewMessengerClient(grpcConn)

	commandCh := make(chan string)
	go readCommands(commandCh)

	ListenMsg()

	for cmd := range commandCh {
		params := strings.Split(cmd, " ")
		switch params[0] {
		case "/faces", "/f":
			ListFaces()
		case "/nf":
			params := strings.SplitN(cmd, " ", 3)
			NewFace(params[1], params[2])
		case "/conn":
			CreateConn(params[1], params[2])
		case "/conns", "/cc":
			ListConns()
		case "/text":
			params := strings.SplitN(cmd, " ", 3)
			i, _ := strconv.Atoi(params[1])
			SendMsg(int64(i), params[2])
		case "/help", "/h":
			fmt.Println(magenta(helpMessage))
		default:
			fmt.Println(red("Unknown command, for help type - /help (h)"))
		}
	}
}

func readCommands(commandCh chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		scanner.Scan()
		command := scanner.Text()
		if command == "/quit" || command == "/q" {
			fmt.Println("Finishing...")
			close(commandCh)
			return
		}
		commandCh <- scanner.Text()
	}
}

func getGrpcConn(grpcAddr string) *grpc.ClientConn {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	grcpConn, err := grpc.Dial(
		grpcAddr,
		opts...,
	)
	if err != nil {
		panic(err)
	}
	return grcpConn
}
