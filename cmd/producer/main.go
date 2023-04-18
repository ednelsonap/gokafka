package main

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func main() {
	deliveryChan := make(chan kafka.Event)
	producer := NewKafkaProducer()
	Publish("mensagem", "teste", producer, nil, deliveryChan)

	go DeliveryReport(deliveryChan) // executa em outra thread para ser async
	producer.Flush(1000)

	// trecho de código responsável por pegar as informações da publicação da mensagem de forma síncrona
	/*e := <-deliveryChan
	  msg := e.(*kafka.Message)
	  if msg.TopicPartition.Error != nil {
	  	fmt.Println("Erro ao enviar")
	  } else {
	  	fmt.Println("Mensagem enviada:", msg.TopicPartition)
	  }*/
}

func NewKafkaProducer() *kafka.Producer {
	configMap := &kafka.ConfigMap{
		"bootstrap.servers":   "gokafka_kafka_1:9092",
		"delivery.timeout.ms": "0",
		"acks":                "all",
		"enable.idempotence":  "true",
	}
	p, err := kafka.NewProducer(configMap)
	if err != nil {
		log.Println(err.Error())
	}
	return p
}

func Publish(msg string, topic string, producer *kafka.Producer, key []byte, deliveryChan chan kafka.Event) error {
	message := &kafka.Message{
		Value:          []byte(msg),
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
	}
	err := producer.Produce(message, deliveryChan)
	if err != nil {
		return err
	}
	return nil
}

func DeliveryReport(deliveryChan chan kafka.Event) {
	for e := range deliveryChan {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Println("Erro ao enviar")
			} else {
				fmt.Println("Mensagem enviada", ev.TopicPartition)
				// anotar no banco de dados que a mensagem foi processada.
				// ex: confirmar que a trasferência bancária ocorreu.
			}
		}
	}
}
