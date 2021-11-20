package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "manyface.net/grpc"
	"manyface.net/internal/messenger"
	"manyface.net/internal/middleware"
	"manyface.net/internal/session"
	"manyface.net/internal/user"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

// TODO: move to configuration?
const (
	webPort        = "8080"
	grpcPort       = ":5300"
	homeMtrxServer = "http://localhost:8008"
)

const (
	logFile = "./server.log"
	dbFile  = "../../db/data.db"
)

func main() {
	// Logger
	logger := initLogger()

	// Database
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		logger.Fatal(err)
	}

	// SessionManager
	sm := session.NewSessionManager(db)

	// MessengerServer
	srv := messenger.NewServer(db, logger)
	// go srv.StartSyncers()

	// Handlers
	userHandler := &user.UserHandler{
		Logger: logger,
		Repo:   user.NewRepo(db),
		SM:     sm,
	}
	messengerHandler := &messenger.MessengerHandler{
		Logger: logger,
		Srv:    srv,
		SM:     sm,
	}

	// Start grpc server
	listener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err) // TODO: remove?
		logger.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterMessengerServer(grpcServer, srv)
	logger.Infof("Starting grpc server at :%v", grpcPort)
	fmt.Printf("Starting grpc server at :%v\n", grpcPort)
	go grpcServer.Serve(listener)

	// Routes and middleware
	router := httprouter.New()
	router.POST("/api/reg", userHandler.Register)
	router.POST("/api/login", userHandler.Login)
	router.POST("/api/face", messengerHandler.CreateFace)
	router.GET("/api/face/:FACE_ID", messengerHandler.GetFace)
	router.DELETE("/api/face/:FACE_ID", messengerHandler.DelFace)
	router.GET("/api/faces", messengerHandler.GetFaces)
	router.POST("/api/conn", messengerHandler.CreateConn)
	router.DELETE("/api/conn", messengerHandler.DeleteConn)
	// router.POST("/api/conns", messengerHandler.GetConns)
	mux := middleware.Auth(logger, sm, router)

	// Start web server
	logger.Infof("Starting web server at :%v", webPort)
	fmt.Printf("Starting web server at :%v\n", webPort)
	http.ListenAndServe("localhost:"+webPort, mux)
	if err != nil {
		logger.Fatalf("Can't start web server at :%v port, %v", webPort, err)
		return
	}

	// srv.wg.Wait()
}

func initLogger() *zap.SugaredLogger {
	// option with standard settings
	// logger, _ := zap.NewDevelopment() // or zap.NewProduction()
	// defer logger.Sync()
	// sugarLogger := logger.Sugar()

	writerSyncer := func() zapcore.WriteSyncer {
		file, _ := os.Create(logFile)
		return zapcore.AddSync(file)
	}()
	encoder := func() zapcore.Encoder {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		return zapcore.NewConsoleEncoder(encoderConfig) // or zapcore.NewJSONEncoder
	}()

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)
	logger := zap.New(core)
	return logger.Sugar()
}
