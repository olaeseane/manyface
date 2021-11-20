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
	"strings"
	"time"

	"github.com/ttacon/chalk"
	"google.golang.org/grpc"
)

const grpcAddr string = "127.0.0.1:5300"

// TODO: remove global vars
var (
	s = flag.String("s", "http://localhost:8080", "manyface server")
	u = flag.String("u", "user3", "manyface user")
	p = flag.String("p", "welcome", "manyface password")

	httpCli  = &http.Client{Timeout: time.Second}
	grpcConn *grpc.ClientConn

	loginResp = LoginResp{}

	helpMessage string = `#Commands
/nf <name> <description>			- Create a new face with given name and description (name without spaces)
/faces (f)					- List of my faces
/conn <face id> <peer id>			- Create a connection with given peer
/text <face id> <peer id> <message>			- Send a message to given peer from the given face
/whoami						- Show user id and device id
/rooms						- List of rooms which the client is joined to
/members <room id>				- List of joined members to given room
/quit (q)					- Quit app`

	green, cyan, red, yellow func(string) string

	token string

	urls = map[string][]string{
		"Register":   {"POST", "/api/reg"},
		"Login":      {"POST", "/api/login"},
		"CreateFace": {"POST", "/api/face"},
		"GetFace":    {"GET", "/api/face/:FACE_ID"},
		"DelFace":    {"DELETE", "/api/face/:FACE_ID"},
		"GetFaces":   {"GET", "/api/faces"},
		"CreateConn": {"POST", "/api/conn"},
		"DeleteConn": {"DELETE", "/api/conn"},
	}
)

func main() {
	flag.Parse()
	if *u == "" || *p == "" || *s == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	green = chalk.Green.NewStyle().WithBackground(chalk.Black).Style
	cyan = chalk.Cyan.NewStyle().Style
	red = chalk.Red.NewStyle().Style
	yellow = chalk.Yellow.NewStyle().Style

	reqBody := []byte(fmt.Sprintf(`{"username": "%s","password": "%s"}`, *u, *p))
	req, err := http.NewRequest(urls["Login"][0], *s+urls["Login"][1], bytes.NewBuffer(reqBody))
	// req.Header.Add("AccessToken", srv.AccessToken)
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
	json.Unmarshal(respBody, &loginResp)

	fmt.Println("Logging into", cyan(*s), "as", cyan(*u))

	fmt.Println(helpMessage)

	grpcConn := getGrpcConn(grpcAddr)
	defer grpcConn.Close()

	commandCh := make(chan string)
	go readCommands(commandCh)

	for cmd := range commandCh {
		params := strings.Split(cmd, " ")
		switch params[0] {
		case "/faces", "/f":
			ListFaces()
		case "/nf":
			params := strings.SplitN(cmd, " ", 3)
			NewFace(params[1], params[2])
			/*		case "/conn":
						CreateConn(params[1], params[2])
					case "/text":
						params := strings.SplitN(cmd, " ", 4)
						SendMessage(params[1], params[2], params[3])
			*/
		// case "/whoami":
		// 	resp, err := client.Whoami()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	fmt.Printf("User %s, Device %s\n", cyan(resp.UserID.String()), cyan(resp.DeviceID.String()))
		// case "/rooms":
		// 	resp, err := client.JoinedRooms()
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	for _, r := range resp.JoinedRooms {
		// 		var out interface{}
		// 		client.GetRoomAccountData(r, "", out)
		// 		// fmt.Println(out.(string))
		// 		fmt.Println(cyan(string(r)))
		// 	}
		// case "/members":
		// 	// resp, err := client.Members(id.RoomID(params[1]))
		// 	resp, err := client.JoinedMembers(id.RoomID(params[1]))
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// 	for k, m := range resp.Joined {
		// 		fmt.Printf("%s %s\n", cyan(string(k)), cyan(*m.DisplayName))
		// 	}
		case "/help", "/h":
			fmt.Println(yellow(helpMessage))
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

/*
func DoGrpc() {

	// cli := pb.NewMessengerClient(conn)

	req1 := &pb.SendRequest{
		Message:        "this is new 3 test message from User 1 to User 3",
		SenderFaceId:   "ada2da68aace03fa2891efbb9314f2c1",
		ReceiverFaceId: "cd456136f58c0abb5e80bd0b360d9228",
	}
	res, err := client.Send(context.Background(), req1)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	req2 := pb.ListenRequest{
		SenderFaceId:   "ada2da68aace03fa2891efbb9314f2c1",
		ReceiverFaceId: "cd456136f58c0abb5e80bd0b360d9228",
	}
	stream, _ := client.Listen(context.Background(), &req2)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)
	}

}
*/

/*
func main2() {
	var err error

	green = chalk.Green.NewStyle().WithBackground(chalk.Black).Style
	cyan = chalk.Cyan.NewStyle().Style
	red = chalk.Red.NewStyle().Style
	yellow = chalk.Yellow.NewStyle().Style

	flag.Parse()
	// if *username == "" || *password == "" || *homeserver == "" {
	// 	_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	// 	flag.PrintDefaults()
	// 	os.Exit(1)
	// }
	if *homeserver == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Logging into", *homeserver)
	// fmt.Println("Logging into", *homeserver, "as", *username)
	client, err = mautrix.NewClient(*homeserver, "", "")
	if err != nil {
		panic(err)
	}
	// _, err = client.Login(&mautrix.ReqLogin{
	// 	Type:             "m.login.password",
	// 	Identifier:       mautrix.UserIdentifier{Type: mautrix.IdentifierTypeUser, User: *username},
	// 	Password:         *password,
	// 	StoreCredentials: true,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Login successful")

	syncer := client.Syncer.(*mautrix.DefaultSyncer)
	syncer.OnEventType(event.EventMessage, func(source mautrix.EventSource, evt *event.Event) {
		fmt.Printf("<%s> %s (%s/%s) - %s - %s\n", cyan(string(evt.Sender)), green(evt.Content.AsMessage().Body), evt.Type.String(), evt.ID, evt.RoomID, evt.ToDeviceID)
	})

	commandCh := make(chan string)

	go readCommands(commandCh)

	go client.Sync()
	defer client.StopSync()

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
		case "/text":
			params := strings.SplitN(cmd, " ", 4)
			SendMessage(params[1], params[2], params[3])
		case "/whoami":
			resp, err := client.Whoami()
			if err != nil {
				panic(err)
			}
			fmt.Printf("User %s, Device %s\n", cyan(resp.UserID.String()), cyan(resp.DeviceID.String()))
		case "/rooms":
			resp, err := client.JoinedRooms()
			if err != nil {
				panic(err)
			}
			for _, r := range resp.JoinedRooms {
				var out interface{}
				client.GetRoomAccountData(r, "", out)
				// fmt.Println(out.(string))
				fmt.Println(cyan(string(r)))
			}
		case "/members":
			// resp, err := client.Members(id.RoomID(params[1]))
			resp, err := client.JoinedMembers(id.RoomID(params[1]))
			if err != nil {
				panic(err)
			}
			for k, m := range resp.Joined {
				fmt.Printf("%s %s\n", cyan(string(k)), cyan(*m.DisplayName))
			}
		case "/help", "/h":
			fmt.Println(yellow(helpMessage))
		default:
			fmt.Println(red("Unknown command, for help type - /help (h)"))
		}
	}
}

*/
