package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"manyface.net/internal/blobstorage"
	"manyface.net/internal/config"
	"manyface.net/internal/messenger"
	"manyface.net/internal/middleware"
	"manyface.net/internal/session"
	"manyface.net/internal/user"
)

// @title Manyface proxy server
// @version 0.2.1
// @BasePath /api/v2beta1

const appName = "manyface"

func main() {
	// Load configuration
	cfg := &config.Config{}
	if err := config.Read(appName, cfg); err != nil {
		panic(err) // TODO: logger?
	}

	// Logger
	logger := initLogger(cfg)

	// BLOBStorage
	storage, err := blobstorage.NewFSStorage(cfg.Data.BLOB)
	if err != nil {
		logger.Fatal("cant create blobstorage", err)
	}

	// Database
	db, err := sql.Open("sqlite3", cfg.Data.DB)
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
	mtrxServer := fmt.Sprintf("%s://%s:%s", cfg.Matrix.Protocol, cfg.Matrix.Host, cfg.Matrix.Port)
	srv := messenger.NewProxy(db, logger, mtrxServer)
	fmt.Printf("Using matrix server at %v\n", mtrxServer)

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
		BS:     storage,
	}

	// Start grpc server
	listener, err := net.Listen("tcp", ":"+cfg.Grpc.Port)
	// listener, err := net.Listen("tcp", "localhost:5300") // NOTE: remove if deploy into k8s
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err) // TODO: remove?
		logger.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	/*
		tls := true
		if tls {
			certFile := "./configs/ssl/server.crt"
			keyFile := "./configs/ssl/server.pem"
			creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
			if sslErr != nil {
				grpclog.Fatalf("Failed loading certificates: %v", sslErr) // TODO: remove?
				logger.Fatalf("Failed loading certificates: %v", sslErr)
			}
			opts = append(opts, grpc.Creds(creds))
		}
	*/
	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer) // register reflection service on gRPC server
	messenger.RegisterMessengerServer(grpcServer, srv)
	logger.Infof("Starting grpc server at :%v", cfg.Grpc.Port)
	fmt.Printf("Starting grpc server at :%v\n", cfg.Grpc.Port)
	go grpcServer.Serve(listener)

	// Routes and middleware
	router := httprouter.New()

	router.POST("/api/v2beta1/user", userHandler.Register)
	router.GET("/api/v2beta1/user", userHandler.Login)

	router.POST("/api/v2beta1/face", messengerHandler.CreateFace)
	router.GET("/api/v2beta1/face/:FACE_ID", messengerHandler.GetFace)
	router.GET("/api/v2beta1/faces", messengerHandler.GetFaces)
	router.DELETE("/api/v2beta1/face/:FACE_ID", messengerHandler.DelFace)
	router.PUT("/api/v2beta1/face/:FACE_ID", messengerHandler.UpdFace)
	router.GET("/api/v2beta1/qr/:FACE_ID", messengerHandler.GetFaceQR)
	router.GET("/api/v2beta1/avatar/:FACE_ID", messengerHandler.GetFaceAvatar)

	router.POST("/api/v2beta1/conn", messengerHandler.CreateConn)
	router.DELETE("/api/v2beta1/conn", messengerHandler.DeleteConn)
	router.GET("/api/v2beta1/conn", messengerHandler.GetConns)
	mux := middleware.Auth(logger, sm, router)

	// Start rest api server
	logger.Infof("Starting rest api server at :%v", cfg.Rest.Port)
	fmt.Printf("Starting rest api server at :%v\n", cfg.Rest.Port)
	http.ListenAndServe(":"+cfg.Rest.Port, mux)
	// http.ListenAndServe("localhost:8080", mux) // NOTE: remove if deploy into k8s

	if err != nil {
		logger.Fatalf("Can't start rest api server at :%v port, %v", cfg.Rest.Port, err)
		return
	}

}

func initLogger(cfg *config.Config) *zap.SugaredLogger {
	// option with standard settings
	// logger, _ := zap.NewDevelopment() // or zap.NewProduction()
	// defer logger.Sync()
	// sugarLogger := logger.Sugar()

	writerSyncer := func() zapcore.WriteSyncer {
		file, _ := os.Create(cfg.Log.File)
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
