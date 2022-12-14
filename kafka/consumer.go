package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"kafka_test/postgres"
	"log"
	"os"
	"sync"
)

func NewConsumer(topic string) {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	brokers := []string{`localhost:9092`}

	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_0_0_0
	cfg.Consumer.Return.Errors = true
	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumer(brokers, cfg)
	if err != nil {
		log.Fatal(err)
	}

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		log.Fatal(err)
	}
	messages := make(chan *sarama.ConsumerMessage, 256)
	initialOffset := sarama.OffsetNewest

	tr := postgres.NewTaskRepository()

	wg := sync.WaitGroup{}

	for _, partition := range partitionList {
		partition := partition
		pc, err := consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			log.Fatal(err)
		}
		go func(pc sarama.PartitionConsumer) {
			wg.Add(1)
			defer wg.Done()
			fmt.Println("got new partition")
			for message := range pc.Messages() {
				messages <- message
				id, err := tr.PostTaskToDB(string(message.Value), message.Timestamp)
				if err != nil {
					fmt.Println(err)
				}

				fmt.Printf("Posted new task to DB:\n %s\nId = %d\n", string(message.Value), id)
			}
		}(pc)
	}
	wg.Wait()
}
