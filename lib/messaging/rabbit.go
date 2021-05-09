package messaging

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"imabad.dev/do/lib/utils"
)

var connection *amqp.Connection
var channel *amqp.Channel
var queues = map[string]RegisteredQueue{}
var done chan bool

//Setup sets up a new RabbitMQ connection, queues and listeners using the central configuration file
func Setup() {
	config := utils.GetConfig()
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%v", config.Messaging.Username, config.Messaging.Password, config.Messaging.Host, config.Messaging.Port))
	if err != nil {
		log.Fatal("Could not connect to Queue", err)
		return
	}
	connection = conn
	channel, err = connection.Channel()
	if err != nil {
		log.Fatal("Could not get channel", err)
		return
	}
	done = make(chan bool, 1)
}

//Close closes the current open connections to RabbitMQ
func Close() {
	connection.Close()
	channel.Close()
	close(done)
}

//RegisteredQueue is a wrapper for the queue object
type RegisteredQueue struct {
	Name     string
	Callback QueueListener
	Queue    amqp.Queue
}

//QueueListener is a callback function used when a queue receives a message
type QueueListener func([]byte, amqp.Delivery)

//RegisterQueue registers a new queue in the RabbitMQ messaging system
func RegisterQueue(queue RegisteredQueue) RegisteredQueue {
	if len(queue.Name) > 0 {
		if _, ok := queues[queue.Name]; !ok {
			q, err := channel.QueueDeclare(queue.Name, false, false, false, false, nil)
			if err != nil {
				log.Fatalf("Could not create %s queue %e", queue.Name, err)
				return queue
			}
			queue.Queue = q
			queues[queue.Name] = queue
			startListeningForMessages(queue)
		}
	} else {
		log.Println("Registering queue with no name")
		q, err := channel.QueueDeclare("", false, false, false, false, nil)
		if err != nil {
			log.Fatalf("Could not create %s queue %e", queue.Name, err)
			return queue
		}
		queue.Name = q.Name
		log.Println("Queue declared with no name has name,", q.Name)
		queue.Queue = q
		queues[queue.Name] = queue
		startListeningForMessages(queue)
	}
	return queue
}

func startListeningForMessages(q RegisteredQueue) {
	msgs, err := channel.Consume(q.Queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Could not consume queue", err)
		return
	}
	go func() {
		for d := range msgs {
			if q.Callback != nil {
				q.Callback(d.Body, d)
			}
		}
		if <-done {
			return
		}
	}()
}

//Publish publishes a message on the desired queue
func Publish(queue RegisteredQueue, contentType string, body []byte) {
	channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	})
}

func PublishWithName(queueName string, contentType string, body []byte) {
	channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	})
}

func PublishWithReply(queue RegisteredQueue, replyTo string,
	correlationID string, contentType string, body []byte) {
	log.Println("Sending with replyto, ", replyTo)
	channel.Publish("", queue.Name, false, false, amqp.Publishing{
		ContentType:   contentType,
		Body:          body,
		ReplyTo:       replyTo,
		CorrelationId: correlationID,
	})
}

func Reply(queueName string, correlationID string, contentType string, body []byte) {
	channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType:   contentType,
		Body:          body,
		CorrelationId: correlationID,
	})
}
