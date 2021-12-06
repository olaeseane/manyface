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
	"manyface.net/internal/config"
	"manyface.net/internal/messenger"
	"manyface.net/internal/middleware"
	"manyface.net/internal/session"
	"manyface.net/internal/user"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

const appName = "manyface"

func main() {
	// Load configuration
	cfg := &config.Config{}
	if err := config.Read(appName, cfg); err != nil {
		panic(err) // TODO: logger?
	}

	// Logger
	logger := initLogger(cfg)

	// Database
	db, err := sql.Open("sqlite3", cfg.DB.File)
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
	srv := messenger.NewServer(db, logger, fmt.Sprintf("%s://%s:%s", cfg.Matrix.Protocol, cfg.Matrix.Host, cfg.Matrix.Port))
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
	listener, err := net.Listen("tcp", ":"+cfg.Grpc.Port)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err) // TODO: remove?
		logger.Fatalf("failed to listen: %v", err)
	}
	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	messenger.RegisterMessengerServer(grpcServer, srv)
	logger.Infof("Starting grpc server at :%v", cfg.Grpc.Port)
	fmt.Printf("Starting grpc server at :%v\n", cfg.Grpc.Port)
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
	router.GET("/api/conns", messengerHandler.GetConns)
	mux := middleware.Auth(logger, sm, router)

	// Start rest api server
	logger.Infof("Starting rest api server at :%v", cfg.Rest.Port)
	fmt.Printf("Starting rest api server at :%v\n", cfg.Rest.Port)
	http.ListenAndServe(cfg.Rest.Host+":"+cfg.Rest.Port, mux)
	if err != nil {
		logger.Fatalf("Can't start rest api server at :%v port, %v", cfg.Rest.Port, err)
		return
	}

	// srv.wg.Wait()

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
