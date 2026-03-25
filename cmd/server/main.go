package main

import (
	"fmt"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connectString = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectString)
	if err != nil {
		fmt.Printf("error creaing connection to RabbitMQ: %v", err)
	}
	defer connection.Close()

	fmt.Println("Connection successful!!")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Program shutting down")
	os.Exit(1)
}
