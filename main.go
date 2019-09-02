package main

import (
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/satori/go.uuid"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	// brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9092").Strings()
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:32770").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("justatopic").String()
)

func main() {
	kingpin.Parse()
	fmt.Println("foo")

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Flush.MaxMessages = 500
	producer, err := sarama.NewAsyncProducer(*brokerList, config)
	if err != nil {
		log.Panicf("failed to setup the kafka producer: %s", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Panicf("failed to close the kafka producer: %s", err)
		}
	}()

	go func() {
		for err := range producer.Errors() {
			log.Panicf("failed to send msg (key %s): %s", err.Msg.Key, err.Err)
		}
	}()

	// UUID := uuid.NewV4()
	UUID := uuid.Must(uuid.NewV4())
	msg := "hello world"
	err = sendUpdate(producer.Input(), UUID.String(), msg)
	if err != nil {
		log.Panicf("failed to send update massage: %s", err)
	}

}

func sendUpdate(ch chan<- *sarama.ProducerMessage, UUID string, msg string) error {
	// bytes, err := proto.Marshal(msg)
	// if err != nil {
	// 	return fmt.Errorf("failed to serialize product delete massage: %s", err)
	// }
	ch <- &sarama.ProducerMessage{
		Topic: *topic,
		Key:   sarama.StringEncoder(UUID),
		// Value: sarama.ByteEncoder(bytes),
		Value: sarama.StringEncoder(msg),
	}
	return nil
}
