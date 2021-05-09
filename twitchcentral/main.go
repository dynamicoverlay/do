package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
	"imabad.dev/do/twitchcentral/twitch"
	"io"
	"log"
	rand2 "math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var twitchOAuthConfig *oauth2.Config
var redisClient *redis.Client
var aesCipher cipher.Block
var authCodeQueue amqp.Queue
var channel *amqp.Channel

type AddAuthCodeRequest struct {
	Code     string `json:"code"`
	TwitchID string `json:"twitchID"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateRandomString(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand2.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	fmt.Println("Starting app...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not load .env, defaulting")
	}

	tokenKey := os.Getenv("TOKEN_KEY")
	aesCipher, err = aes.NewCipher([]byte(tokenKey))
	if err != nil {
		fmt.Println("Error creating new AES cipher", err)
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%v", os.Getenv("HERMES_USERNAME"), os.Getenv("HERMES_PASSWORD"), os.Getenv("HERMES_ADDRESS"), os.Getenv("HERMES_PORT")))
	if err != nil {
		log.Fatal("Could not connect to Queue", err)
		return
	}
	defer conn.Close()
	channel, err = conn.Channel()
	if err != nil {
		log.Fatal("Could not get channel", err)
		return
	}
	defer channel.Close()
	authCodeQueue, err = channel.QueueDeclare("authCodes", false, false, false, false, nil)
	if err != nil {
		log.Fatal("Could not create email queue", err)
		return
	}
	twitchOAuthConfig = &oauth2.Config{
		ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
		ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("TWITCH_REDIRECT"),
		Scopes:       []string{"user:read:email"},
		Endpoint:     endpoints.Twitch,
	}

	fmt.Println("Starting HTTP server on port 8083")
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := newWebserver()
	go gracefulShutdown(server, quit, done)

	fmt.Println("Server is ready to handle requests at :8083")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Errorf("could not listen on %s: %v\n", ":8083", err)
	}
	<-done
	fmt.Println("Server stopped")
}

func AddAuthCode(request AddAuthCodeRequest) error {
	bodyBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	channel.Publish("", authCodeQueue.Name, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        bodyBytes,
	})
	return nil
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

	router.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		redirectTo := request.URL.Query().Get("redirect")
		fmt.Printf("Redirect is %s\n", redirectTo)
		if len(redirectTo) > 0 {
			fmt.Println("Set cookie to", redirectTo)
			http.SetCookie(writer, &http.Cookie{Name: "redirect", Value: redirectTo})
		}
		fmt.Println("Redirect URI: ", len(twitchOAuthConfig.RedirectURL), twitchOAuthConfig.RedirectURL)
		url := twitchOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
		//fmt.Fprint(writer, "HUH")
		http.Redirect(writer, request, url, 301)
	})

	router.HandleFunc("/auth", func(writer http.ResponseWriter, request *http.Request) {
		ctx := context.Background()
		redirectTo, err := request.Cookie("redirect")
		if err != nil {
			fmt.Fprintf(writer, "Invalid redirect, %v", err)
			return
		}
		token, err := twitchOAuthConfig.Exchange(ctx, request.URL.Query().Get("code"))
		if err != nil {
			fmt.Fprint(writer, "Invalid redirect")
			return
		}
		user, err := twitch.GetUser(token.AccessToken, os.Getenv("TWITCH_CLIENT_ID"))
		if err != nil {
			fmt.Fprint(writer, "Invalid code")
			return
		}
		gcm, err := cipher.NewGCM(aesCipher)
		if err != nil {
			fmt.Fprint(writer, "Invalid cipher")
			return
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			fmt.Fprint(writer, "Error generating nonce", err)
			return
		}
		err = redisClient.Set(ctx, fmt.Sprintf("token:%v", user.ID), gcm.Seal(nonce, nonce, []byte(token.RefreshToken), nil), 0).Err()
		if err != nil {
			fmt.Fprint(writer, "Error saving to database", err)
			return
		}
		authCode := GenerateRandomString(20)
		err = AddAuthCode(AddAuthCodeRequest{
			Code:     authCode,
			TwitchID: user.ID,
		})
		if err != nil {
			fmt.Fprint(writer, "Error saving to database", err)
			return
		}
		http.Redirect(writer, request, fmt.Sprintf("%s?code=%v", redirectTo.Value, authCode), http.StatusTemporaryRedirect)
	})

	return &http.Server{
		Addr:    ":8083",
		Handler: router,
	}
}
