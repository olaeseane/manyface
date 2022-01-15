package main

import (
	"bufio"
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

	"manyface.net/internal/messenger"
)

// const grpcAddr string = "127.0.0.1:5300"

// TODO: remove global vars
var (
	ws = flag.String("ws", "http://localhost:8080", "manyface server")
	gs = flag.String("gs", "127.0.0.1:5300", "manyface server")
	u  = flag.String("u", "l1RaTFQdzwlM2TpOX5xs_", "manyface user")
	p  = flag.String("p", "welcome", "manyface password")

	httpCli = &http.Client{Timeout: time.Second * 5}

	grpcConn *grpc.ClientConn
	grpcCli  messenger.MessengerClient

	loginResp = LoginResp{}

	helpMessage string = `---
#Commands
/nf <nick> <purpose> <bio> <comments> <server>		- Create a new face with given name and description (name without spaces)
/faces (f)						- List of my faces
/conn <my face id> <peer face id>			- Create a connection with given peer
/conns (cc) 						- List of my connections
/text <conn id> <message>				- Send a message to given peer from the given face
/quit (q)						- Quit app
---
`

	green, cyan, red, yellow, magenta, blue func(string) string

	urls = map[string][]string{
		"Register":   {"POST", "/api/v2beta1/user"},
		"Login":      {"GET", "/api/v2beta1/user"},
		"CreateFace": {"POST", "/api/v2beta1/face"},
		"GetFace":    {"GET", "/api/v2beta1/face/"},    // :FACE_ID
		"DelFace":    {"DELETE", "/api/v2beta1/face/"}, // :FACE_ID
		"GetFaces":   {"GET", "/api/v2beta1/faces"},
		"CreateConn": {"POST", "/api/v2beta1/conn"},
		"DeleteConn": {"DELETE", "/api/v2beta1/conn"},
		"GetConns":   {"GET", "/api/v2beta1/conns"},
	}
)

func main() {
	flag.Parse()

	if isFlagPassed("-h") || isFlagPassed("--help") || isFlagPassed("-help") {
		flag.PrintDefaults()
	}
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

	req, err := http.NewRequest(urls["Login"][0], *ws+urls["Login"][1], nil)
	req.SetBasicAuth(*u, *p)
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
	grpcCli = messenger.NewMessengerClient(grpcConn)

	commandCh := make(chan string)
	go readCommands(commandCh)

	ListenMsg()

	for cmd := range commandCh {
		params := strings.Split(cmd, " ")
		switch params[0] {
		case "/faces", "/f":
			ListFaces()
		case "/nf":
			params := strings.SplitN(cmd, " ", 6)
			NewFace(params[1], params[2], params[3], params[4], params[5])
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

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
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
