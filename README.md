# amqp-publish

A simple tool to publish messages to RabbitMQ from the command line.

## Setup

Download the [latest release](https://github.com/selency/amqp-publish/releases) binary and save it to `/usr/local/bin` or any executable path.

## Usage

Publish to exchange

```SHELL
amqp-publish --uri="amqp://admin:password@localhost:5672/" --exchange="foo" --routing-key="awesome-routing-key" --body="hello, world!"
```

Publish the `bar` queue directly, using RabbitMQ default exchange

```SHELL
amqp-publish --uri="amqp://admin:password@localhost:5672/" --exchange="" --routing-key="bar" --body="hello, world!"
```

Publish using an input file

```SHELL
# Use one message per line
echo '{"foo": "bar-1"}' >> file.txt
echo '{"foo": "bar-2"}' >> file.txt
echo '{"foo": "bar-3"}' >> file.txt

amqp-publish --uri="amqp://admin:password@localhost:5672/" --exchange="foo" --routing-key="bar" --input-file="file.txt"
```

## Credit

Streadway's [AMQP Go library](https://github.com/streadway/amqp).
