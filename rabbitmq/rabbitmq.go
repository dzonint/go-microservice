package rabbitmq

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dzonint/go-microservice/config"
	"github.com/dzonint/go-microservice/data"
	"github.com/streadway/amqp"
)

type RabbitMQService struct{}

func NewRabbitMQService() *RabbitMQService {
	return &RabbitMQService{}
}

func handleError(err error, msg string, isFatal bool) {
	if err != nil {
		if isFatal != false {
			log.Fatalf("%s: %s", msg, err)
		} else {
			log.Printf("%s: %s", msg, err)
		}
	}
}

func (rmq *RabbitMQService) GenerateConsumer() {
	conn, err := amqp.Dial(config.Config.RabbitMQUrl)
	handleError(err, "[Consumer] Can't connect to AMQP", true)
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "[Consumer] Can't create a amqpChannel", true)
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("users", true, false, false, false, nil)
	handleError(err, "[Consumer] Could not declare `users` queue", true)

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "[Consumer] Could not configure QoS", true)

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "[Consumer] Could not register consumer", true)

	log.Printf("Consumer ready, PID: %d", os.Getpid())
	for d := range messageChannel {
		log.Printf("[Consumer] Received a message: %s", d.Body)

		user := &data.User{}

		err := json.Unmarshal(d.Body, user)

		if err != nil {
			log.Printf("[Consumer] Error decoding JSON: %s", err)
		}

		err = data.AddUser(user)
		if err != nil {
			log.Printf("[Consumer] Error adding user: %s", err)
		} else {
			log.Printf("[Consumer] User %v added successfuly", user.Email)
		}

		if err := d.Ack(false); err != nil {
			log.Printf("[Consumer] Error acknowledging message : %s", err)
		} else {
			log.Printf("[Consumer] Acknowledged message")
		}
	}
}

func (rmq *RabbitMQService) fetchUser() ([]byte, error) {
	res, err := http.Get(config.Config.GenerateUserUrl)
	handleError(err, "Could not fetch user", false)
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func (rmq *RabbitMQService) GenerateProducer() {
	conn, err := amqp.Dial(config.Config.RabbitMQUrl)
	handleError(err, "[Producer] Can't connect to AMQP", true)
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "[Producer] Can't create a amqpChannel", true)
	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare("users", true, false, false, false, nil)
	handleError(err, "[Producer] Could not declare `users` queue", true)

	log.Printf("Producer ready. Producer will publish new user every 10 seconds.")
	for {
		user, err := rmq.fetchUser()
		handleError(err, "Error", true)
		if err != nil {
			log.Printf("[Producer] Error fetching user: %s", err)
		}
		msg := amqp.Publishing{
			Body: user,
		}
		err = amqpChannel.Publish("", queue.Name, false, false, msg)
		handleError(err, "[Producer] Error publishing message:", true)

		log.Printf("[Producer] AddUser: %v", string(user))
		time.Sleep(10 * time.Second)
	}
}
