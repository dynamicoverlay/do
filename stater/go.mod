module imabad.dev/do/stater

go 1.14

require (
	github.com/go-redis/redis/v8 v8.0.0-beta.5
	github.com/google/uuid v1.1.1
	github.com/gorilla/websocket v1.4.2
	github.com/joho/godotenv v1.3.0
	github.com/streadway/amqp v1.0.0
	imabad.dev/do/lib v0.0.0
)

replace imabad.dev/do/lib v0.0.0 => /home/stuart/DynamicOverlay/lib
