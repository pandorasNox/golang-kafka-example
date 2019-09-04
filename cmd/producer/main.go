package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	// ,localhost:32772,localhost:32773
	//go run . --brokerList localhost:32771 --brokerList localhost:32772 --brokerList localhost:32772
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("kafka1:9092").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("justatopic").String()
)

func main() {
	kingpin.Parse()
	fmt.Println("init...")

	fmt.Println("")
	// fmt.Printf("%v", brokerList)
	// fmt.Println("")
	// fmt.Printf("len: %v", len(*brokerList))
	// fmt.Println("")

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500
	producer, err := sarama.NewAsyncProducer(*brokerList, config)
	if err != nil {
		log.Panicf("failed to setup the kafka producer: %s", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Printf("failed to close the kafka producer: %s", err)
		}
	}()

	go func() {
		for err := range producer.Errors() {
			log.Printf("failed to send msg (key %s): %s", err.Msg.Key, err.Err)
		}
	}()

	msg := "hello world"
	fmt.Println("sending...")
	// for {
	// 	producer.Input() <- &sarama.ProducerMessage{Topic: *topic, Value: sarama.StringEncoder(msg)}
	// }
	// producer.Input() <- &sarama.ProducerMessage{Topic: *topic, Value: sarama.StringEncoder(msg)}
	producer.Input() <- &sarama.ProducerMessage{Topic: *topic, Value: sarama.StringEncoder(msg)}

	fmt.Println("done")
}
