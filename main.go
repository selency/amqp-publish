package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

var (
	uri           string
	exchange      string
	routingKey    string
	body          string
	inputFilePath string
)

func validateFlags() error {
	if uri == "" {
		return errors.New("uri cannot be blank")
	}

	if exchange == "" && routingKey == "" {
		return errors.New("exchange and routing-key cannot both be blank")
	}

	if body == "" && inputFilePath == "" {
		return errors.New("body and input-file cannot both be blank")
	}

	return nil
}

func init() {
	flag.StringVar(&uri, "uri", "", "AMQP URI amqp://<user>:<password>@<host>:<port>/[vhost]")
	flag.StringVar(&exchange, "exchange", "", "Exchange name")
	flag.StringVar(&routingKey, "routing-key", "", `Routing key. Use queue
        name with blank exchange to publish directly to queue.`)
	flag.StringVar(&body, "body", "", "Message body")
	flag.StringVar(&inputFilePath, "input-file", "", "Input file path")

	flag.Parse()

	err := validateFlags()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getMessages() ([]string, error) {
	messages := []string{}

	if inputFilePath != "" {
		b, err := ioutil.ReadFile(inputFilePath)

		if err != nil {
			return messages, errors.New("failed to read input file")
		}

		lines := strings.Split(string(b), "\n")

		for _, l := range lines {
			if l == "" {
				continue
			}

			messages = append(messages, l)
		}

	} else {
		messages = append(messages, body)
	}

	return messages, nil
}

func main() {
	connection, err := amqp.Dial(uri)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer connection.Close()

	channel, _ := connection.Channel()

	messages, err := getMessages()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf("%d messages to publish", len(messages))

	for _, m := range messages {
		channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(m),
			DeliveryMode:    amqp.Transient,
			Priority:        0,
		})
	}
}
