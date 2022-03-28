package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
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
	"google.golang.org/grpc/credentials"

	"manyface.net/internal/messenger"
)

// const grpcAddr string = "127.0.0.1:5300"
const caCertPath = "../../configs/ssl/kyma.pem"

// const caCertPath = "../../configs/ssl/server.crt"

// TODO: remove global vars
var (
	// fRest = flag.String("rest", "http://localhost:8080", "manyface rest server")
	fRest = flag.String("rest", "https://manyface.a131084.kyma.ondemand.com:80", "manyface rest server")
	// fGrpc = flag.String("grpc", "127.0.0.1:5300", "manyface grpc server")
	fGrpc = flag.String("grpc", "grpc.a131084.kyma.ondemand.com:5300", "manyface grpc server")
	// u  = flag.String("u", "xaafZ6kkKxX8SvCPFhMHZ", "manyface user") // another user mC6FneCs2VbK26Fkt9IBp - synapse
	// u   = flag.String("u", "NfMlpnrGSWknJgh6Wt5uN", "manyface user") // another user RlPqljVc5JGuBKrMSBGUs - conduit;
	fUser = flag.String("user", "kVsd6HAQUtGZujPVg-jt8", "manyface user") // kyma user
	fPwd  = flag.String("pwd", "welcome", "manyface password")
	// fTls  = flag.Bool("tls", false, "tls connection")
	fTls = flag.Bool("tls", false, "tls connection")

	httpCli = &http.Client{Timeout: time.Second * 5}

	grpcConn *grpc.ClientConn
	grpcCli  messenger.MessengerClient

	loginResp = LoginResp{}

	globalCtx, globalCancel = context.WithCancel(context.Background())

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
		"GetConns":   {"GET", "/api/v2beta1/conn"},
	}
)

func main() {
	flag.Parse()

	if isFlagPassed("-h") || isFlagPassed("--help") || isFlagPassed("-help") {
		flag.PrintDefaults()
	}
	if *fUser == "" || *fPwd == "" || *fRest == "" || *fGrpc == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("user - %v\nrest - %v\ngrpc - %v\ntls - %v\n\n", *fUser, *fRest, *fGrpc, *fTls)

	green = chalk.Green.NewStyle().WithBackground(chalk.Black).Style
	cyan = chalk.Cyan.NewStyle().Style
	red = chalk.Red.NewStyle().Style
	yellow = chalk.Yellow.NewStyle().Style
	blue = chalk.Blue.NewStyle().Style
	magenta = chalk.Magenta.NewStyle().Style

	if *fTls {
		caCert, err := ioutil.ReadFile(caCertPath)
		if err != nil {
			panic(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		httpCli = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					// RootCAs:            caCertPool,
					InsecureSkipVerify: true,
				},
			},
			Timeout: time.Second * 15,
		}
	}

	req, err := http.NewRequest(urls["Login"][0], *fRest+urls["Login"][1], nil)
	fmt.Printf("%+v\n", req)
	req.SetBasicAuth(*fUser, *fPwd)
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

	fmt.Println("Logging into", cyan(*fRest), "and", cyan(*fGrpc), "as", cyan(*fUser))
	fmt.Println(magenta(helpMessage))

	grpcConn = getGrpcConn(*fGrpc)
	defer grpcConn.Close()
	grpcCli = messenger.NewMessengerClient(grpcConn)
	fmt.Println(grpcConn)

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
			close(commandCh)
			globalCancel()
			fmt.Println("Finishing...")
			return
		}
		commandCh <- scanner.Text()
	}
}

func getGrpcConn(grpcAddr string) *grpc.ClientConn {
	// opts := []grpc.DialOption{
	// 	grpc.WithInsecure(),
	// }
	opts := grpc.WithInsecure()
	if *fTls {
		// caCertFile := "../../configs/ssl/ca.crt"
		// caCertFile := "../../configs/ssl/kyma.cer"
		caCertFile := caCertPath
		creds, sslErr := credentials.NewClientTLSFromFile(caCertFile, "")
		if sslErr != nil {
			panic(fmt.Sprintf("Error while loading CA trust certificate: %v", sslErr))
		}
		opts = grpc.WithTransportCredentials(creds)
	}

	grcpConn, err := grpc.Dial(
		grpcAddr,
		opts,
	)
	if err != nil {
		panic(err)
	}
	return grcpConn
}
