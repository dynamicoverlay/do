package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"imabad.dev/do/lib/messaging"
)

var redisClient *redis.Client
var changeStateQueue amqp.Queue
var channel *amqp.Channel

var connectedClients sync.Map
var authedClients []string
var authRequests sync.Map

func checkOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
}

type UpdateStateMessage struct {
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"value"`
}

func updateState(context context.Context, r messaging.ChangeStateRequest) {
	var state map[string]interface{}
	val, err := redisClient.Get(context, r.Overlay).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &state)
	} else if err == redis.Nil {
		state = map[string]interface{}{}
	} else {
		log.Println("Failed to get state", err)
	}
	state[r.Key] = r.Value
	newBytes, err := json.Marshal(state)
	if err != nil {
		log.Printf("Failed to marshal new state,", err)
		return
	}
	err = redisClient.Set(context, r.Overlay, string(newBytes), 0).Err()
	if err != nil {
		log.Printf("Failed to update state for overlay %s %e", r.Overlay, err)
		return
	}
	if value, ok := connectedClients.Load(r.Overlay); ok {
		clients := value.([]*websocket.Conn)
		log.Printf("Publishing to clients")
		BroadcastToClients(clients, UpdateStateMessage{
			Message: "updateState",
			Data:    state,
		})
	}
}

func GetCurrentState(overlay string) map[string]interface{} {
	var state map[string]interface{}
	val, err := redisClient.Get(context.Background(), overlay).Result()
	if err == nil {
		json.Unmarshal([]byte(val), &state)
	} else if err == redis.Nil {
		state = map[string]interface{}{}
	} else {
		log.Println("Failed to get state", err)
	}
	return state
}

func BroadcastToClients(clients []*websocket.Conn, message interface{}) {
	for _, client := range clients {
		client.WriteJSON(message)
	}
}

var OverlayStateChangeQueue = messaging.RegisteredQueue{
	Name: "changeOverlayState",
	Callback: func(bytes []byte, d amqp.Delivery) {
		var request messaging.ChangeStateRequest
		json.Unmarshal(bytes, &request)
		log.Printf("Received message %s", bytes)
		updateState(context.Background(), request)
	},
}

var OverlayAuthQueue = messaging.RegisteredQueue{
	Name: "overlay-auth",
}

var OverlayAuthReplies = messaging.RegisteredQueue{
	Name: "",
	Callback: func(bytes []byte, d amqp.Delivery) {
		log.Println("received reply")
		var response messaging.OverlayAuthResponse
		err := json.Unmarshal(bytes, &response)
		if err != nil {
			log.Println("Failed to decode message", err)
			return
		}
		if conn, ok := authRequests.Load(d.CorrelationId); ok {
			if response.Authenticated {
				handleWSMessage(response.Identifier, conn.(*websocket.Conn))
			} else {
				conn.(*websocket.Conn).WriteJSON(AuthResponseMessage{Success: false, Message: "authRequest"})
			}
		}
	},
}

func main() {
	fmt.Println("Starting app...")
	err := godotenv.Load()
	if err != nil {
		log.Print("Could not load .env, defaulting")
	}
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	messaging.Setup()
	defer messaging.Close()

	messaging.RegisterQueue(OverlayStateChangeQueue)
	OverlayAuthReplies = messaging.RegisterQueue(OverlayAuthReplies)
	fmt.Println("Name of replies queue is", OverlayAuthReplies.Name)
	fmt.Println("Starting Websocket server")
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := newWebServer()
	go gracefulShutdown(server, quit, done)

	fmt.Println("Server is ready to handle requests")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", ":8083", err)
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

func isAuthed(identifier string) bool {
	for _, client := range authedClients {
		if client == identifier {
			return true
		}
	}
	return false
}

var pongWait = 60 * time.Second

func handleWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	var authMessage *AuthMessage
	//	c.SetReadDeadline(time.Now().Add(pongWait))
	//	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	identifier := uuid.New().String()
	for {
		err := c.ReadJSON(&authMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			} else if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Printf("Client closed connection")
			}
			break
		}
		if authMessage != nil && !isAuthed(identifier) {
			log.Println("Is not authed, doing auth verification")
			authMessage.verify(identifier, c)
		}
		//err = c.WriteMessage(mt, message)
		//if err != nil {
		//		log.Println("write:", err)
		//			break
		//		}
	}
	fmt.Println("Client Disconnected")
	if authMessage != nil {
		log.Println("Removing client from connected clients")
		if value, ok := connectedClients.Load(authMessage.Overlay); ok {
			log.Println("Found list of clients to remove from")
			clients := value.([]*websocket.Conn)
			log.Println("Length of clients is", len(clients))
			for k, v := range clients {
				if v == c {
					log.Println("Found client in list, removing them")
					clients[k] = clients[len(clients)-1]
					clients[len(clients)-1] = nil
					clients = clients[:len(clients)-1]
					log.Println("New length of clients is, ", len(clients))
					connectedClients.Store(authMessage.Overlay, clients)
					break
				}
			}
		}
	}
}

type AuthMessage struct {
	Overlay  string `json:"overlay"`
	Passcode string `json:"passcode"`
}

type AuthResponseMessage struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (auth *AuthMessage) verify(identifier string, client *websocket.Conn) {
	if len(auth.Overlay) <= 0 || len(auth.Passcode) <= 0 {
		client.WriteJSON(AuthResponseMessage{Success: false, Message: "authRequest"})
		return
	}
	authRequest := messaging.OverlayAuthRequest{
		Identifier:    auth.Overlay,
		Pin:           auth.Passcode,
		CorrelationID: identifier,
	}
	bytes, err := json.Marshal(authRequest)
	if err != nil {
		log.Println("Error encoding auth request", err)
		return
	}
	authRequests.Store(identifier, client)
	messaging.PublishWithReply(OverlayAuthQueue, OverlayAuthReplies.Queue.Name, identifier, "application/json", bytes)
	log.Println("Sent auth verification request")
}

// https://myoverlay.io/o/{overlay-uuid}/{overlay-passcode}
func handleWSMessage(overlay string, conn *websocket.Conn) {
	//	if !authMessage.verify()
	//		conn.WriteJSON(AuthResponseMessage{Success: false, Message: "Invalid keys"})
	//		return false
	//	}
	foundValue, ok := connectedClients.Load(overlay)
	var overlayClients []*websocket.Conn
	if !ok {
		overlayClients = []*websocket.Conn{}
	} else {
		overlayClients = foundValue.([]*websocket.Conn)
	}
	overlayClients = append(overlayClients, conn)
	connectedClients.Store(overlay, overlayClients)
	conn.WriteJSON(AuthResponseMessage{Success: true, Message: "authRequest"})
	conn.WriteJSON(UpdateStateMessage{
		Message: "updateState",
		Data:    GetCurrentState(overlay),
	})
}

func newWebServer() *http.Server {
	router := http.NewServeMux()
	router.HandleFunc("/ws", handleWS)

	return &http.Server{
		Addr:    ":8083",
		Handler: router,
	}
}
