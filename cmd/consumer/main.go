package main

import (
	"log"

	"github.com/Shopify/sarama"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka1:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("justatopic").String()
)

func main() {
	log.Println("Hello, world!")

	kingpin.Parse()

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewClient(*brokerList, config)
	if err != nil {
		log.Panicf("failed to setup kafka client: %s", err)
	}
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		log.Panicf("failed to setup kafka consumer: %s", err)
	}
	partition, err := consumer.ConsumePartition(*topic, 0, 0)
	if err != nil {
		log.Panicf("failed to setup kafka partition: %s", err)
	}

	log.Printf("starting to process %v", topic)

	for {
		select {
		case err := <-partition.Errors():
			log.Printf("failure from kafka consumer: %s", err)

		case msg := <-partition.Messages():
			value := string(msg.Value)
			log.Printf("recieved message with offset %v: %v", msg.Offset, value)
		}
	}
}
