package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"message-emmitter/message"
	"message-emmitter/settings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

var wg = sync.WaitGroup{}

func publish(message string) {
	fmt.Println(message)
	wg.Done()
}

func publishMessage(ctx context.Context, body string, ch *amqp.Channel, mark string, settings settings.Settings) {
	err := ch.PublishWithContext(ctx,
		settings.Rabbit.Exhange,
		settings.Rabbit.RoutingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent message for: %v\n", mark)
	wg.Done()
}

func main() {

	var fileName string

	settings := &settings.Settings{}
	settings.GetSettings()

	conn, err := amqp.Dial(settings.GetRabbitmqUrl())
	log.Println(settings.GetRabbitmqUrl())
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	flag.StringVar(&fileName, "f", "temp.txt", "Filename")
	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		m := message.Message{}
		m.Body.Password = scanner.Text()
		m.Body.Password = scanner.Text()
		body, err := json.Marshal(m)
		failOnError(err, "Unable to marshal item")
		wg.Add(1)
		go publishMessage(ctx, string(body), ch, m.Body.Password, *settings)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
