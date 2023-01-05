package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"kafka_test/cache"
	"kafka_test/repository"
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

	tr := repository.NewTaskRepository()
	ch := cache.NewCache(0, 0)

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
			for message := range pc.Messages() {
				messages <- message
				id, err := tr.PostTaskToDB(string(message.Value), message.Timestamp)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Printf("Posted new task to DB:\n %s\nId = %d\n", string(message.Value), id)
				ch.Set(id, string(message.Value), message.Timestamp, 0)
				fmt.Println(ch.Tasks[id])

			}
		}(pc)
	}
	wg.Wait()
}
