package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/delivery"
	mongo2 "github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/repository/mongo"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type App struct {
	httpServer *http.Server

	authUseCase auth.UseCase
}

func NewApp() *App {
	db := initDB()

	userRepo := mongo2.NewUserRepository(db, viper.GetString("mongo.authCollection"))
	chatRepo := mongo2.NewChatRepository(db, viper.GetString("mongo.chatCollection"))
	messageRepo := mongo2.NewMessageRepository(db, viper.GetString("mongo.chatCollection"),
		viper.GetString("mongo.messageCollection"))
	authUseCase := usecase.NewAuthorizer(
		userRepo,
		chatRepo,
		messageRepo,
		viper.GetString("auth.hash_salt"),
		[]byte(viper.GetString("auth.signing_key")),
		viper.GetDuration("auth.token_ttl")*time.Second,
	)

	return &App{
		authUseCase: authUseCase,
	}
}

func (a *App) Run(port string) error {
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	authApi := router.Group("/auth")
	messageApi := router.Group("/message")
	chatApi := router.Group("/chat")

	delivery.RegisterHTTPAuthEndpoints(authApi, a.authUseCase)
	delivery.RegisterHTTPMessageEndpoints(messageApi, a.authUseCase)
	delivery.RegisterHTTPChatEndpoints(chatApi, a.authUseCase)

	a.httpServer = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatalf(err.Error())
		log.Fatalf("Error occured while establishing connection to mongoDB")

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.name"))
}
