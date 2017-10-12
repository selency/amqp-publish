package main

import (
    "fmt"
    "os"
    "flag"
    "errors"
    "github.com/streadway/amqp"
)

var (
    uri string
    exchange string
    routingKey string
    body string
)

func validateFlags() error {
    if uri == "" {
        return errors.New("uri cannot be blank")
    }

    if exchange == "" && routingKey == "" {
        return errors.New("exchange and routing-key cannot both be blank")
    }

    if body == "" {
        return errors.New("body cannot be blank")
    }

    return nil
}

func init() {
    flag.StringVar(&uri, "uri", "", "AMQP URI amqp://<user>:<password>@<host>:<port>/[vhost]")
    flag.StringVar(&exchange, "exchange", "", "Exchange name")
    flag.StringVar(&routingKey, "routing-key", "", `Routing key. Use queue
        name with blank exchange to publish directly to queue.`)
    flag.StringVar(&body, "body", "", "Message body")

    flag.Parse()

    err := validateFlags()

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}

func main() {
    connection, err := amqp.Dial(uri)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    defer connection.Close()

    channel, _ := connection.Channel()

    channel.Publish(exchange, routingKey, false, false, amqp.Publishing{
        Headers: amqp.Table{},
        ContentType: "text/plain",
        ContentEncoding: "",
        Body: []byte(body),
        DeliveryMode: amqp.Transient,
        Priority: 0,
    })
}
