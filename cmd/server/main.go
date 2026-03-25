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
		fmt.Printf("error creating connection to RabbitMQ: %v", err)
		os.Exit(1)
	}
	defer connection.Close()

	fmt.Println("Connection Successful!!")

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("\nProgram shutting down")
	os.Exit(1)
}
