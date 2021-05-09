package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/friendsofgo/graphiql"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	"github.com/streadway/amqp"
	"imabad.dev/do/api/handlers"
	"imabad.dev/do/api/messages"
	"imabad.dev/do/api/resolvers"
	"imabad.dev/do/api/schemas"
	consts "imabad.dev/do/api/utils"
	"imabad.dev/do/lib/db"
	"imabad.dev/do/lib/messaging"
	"imabad.dev/do/lib/models"
	"imabad.dev/do/lib/utils"
)

//Db is the DB

func main() {
	utils.LoadConfig()
	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	consts.Db = db
	defer db.Close()

	messaging.Setup()
	defer messaging.Close()
	messaging.RegisterQueue(messages.EmailQueue)
	messaging.RegisterQueue(messaging.RegisteredQueue{
		Name: "authCodes",
		Callback: func(bytes []byte, d amqp.Delivery) {
			var request messages.AddAuthCodeRequest
			json.Unmarshal(bytes, &request)
			log.Printf("Received message %s", bytes)
			if err != nil {
				log.Printf("Failed to send email %v", err)
			}
		},
	})
	messaging.RegisterQueue(messages.OverlayQueue)

	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.EmailVerification{})
	db.AutoMigrate(&models.Overlay{})
	db.AutoMigrate(&models.Module{})
	db.AutoMigrate(&models.OverlayModule{})
	fmt.Println("Starting HTTP server on port 8080")
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := newWebserver()
	go gracefulShutdown(server, quit, done)

	fmt.Println("Server is ready to handle requests at :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Errorf("could not listen on %s: %v", ":8080", err)
	}

	<-done
	fmt.Println("Server stopped")
}

func gracefulShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	fmt.Println("Shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Could not gracefully shutdown server, %v ", err)
	}
	close(done)
}

func newWebserver() *http.Server {
	router := http.NewServeMux()
	schema := graphql.MustParseSchema(*schemas.NewSchema(), &resolvers.RootResolver{Db: consts.Db})

	graphiqlHandler, err := graphiql.NewGraphiqlHandler("/query")
	if err != nil {
		panic(err)
	}
	router.Handle("/graphql", graphiqlHandler)
	router.Handle("/query", handlers.AuthCheckMiddleware(&relay.Handler{Schema: schema}))
	router.HandleFunc("/email", func(writer http.ResponseWriter, request *http.Request) {
		err := messages.SendEmail(messages.EmailRequest{
			To:      "stuart@pomeroys.site",
			Subject: "New Email!",
			Content: "You suck!",
		})
		if err != nil {
			fmt.Fprint(writer, err)
		} else {
			fmt.Fprint(writer, "Email sent!")
		}
	})
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		// Enable Debugging for testing, consider disabling in production
		Debug: true,
	})
	corsRouter := c.Handler(router)
	return &http.Server{
		Addr:    ":8080",
		Handler: corsRouter,
	}
}
