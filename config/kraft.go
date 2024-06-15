package config

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
)

type Kraft struct {
	producer sarama.SyncProducer
	consumer sarama.Consumer
	brokers  []string
}

func newKafkaClient() (*Kraft, error) {
	connStr := fmt.Sprintf("%s:%s", os.Getenv("KRAFT_HOST"), os.Getenv("KRAFT_PORT"))
	client := &Kraft{
		brokers: []string{connStr},
	}
	if err := client.connectProducer(); err != nil {
		return nil, err
	}
	if err := client.connectConsumer(); err != nil {
		return nil, err
	}
	return client, nil
}

// connectProducer connects to Kafka as a producer
func (k *Kraft) connectProducer() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = os.Getenv("KRAFT_USER")
	config.Net.SASL.Password = os.Getenv("KRAFT_PASS")
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
	}

	producer, err := sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		return err
	}
	k.producer = producer
	return nil
}

// connectConsumer connects to Kafka as a consumer
func (k *Kraft) connectConsumer() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = os.Getenv("KRAFT_USER")
	config.Net.SASL.Password = os.Getenv("KRAFT_PASS")
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = &tls.Config{
		InsecureSkipVerify: true,
	}

	consumer, err := sarama.NewConsumer(k.brokers, config)
	if err != nil {
		return err
	}
	k.consumer = consumer
	return nil
}

// PushMessageToQueue sends a message to a Kafka topic
func (k *Kraft) PushMessageToQueue(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	return nil
}

// ConsumeMessages consumes messages from a Kafka topic
func (k *Kraft) ConsumeMessages(topic string, partition int32) {
	consumer, err := k.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Error consuming partition: %v\n", err)
		return
	}
	defer consumer.Close()

	for {
		select {
		case msg := <-consumer.Messages():
			log.Printf("Consumed message: %s\n", string(msg.Value))
		case err := <-consumer.Errors():
			log.Printf("Error consuming: %v\n", err)
		}
	}
}

// Close closes the Kafka producer and consumer
func (k *Kraft) Close() {
	if k.producer != nil {
		k.producer.Close()
	}
	if k.consumer != nil {
		k.consumer.Close()
	}
}
