package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
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
				//fmt.Printf("[%d] %s\n", partition, message.Value)
			}
		}(pc)
	}
	fmt.Println("Here we are")
	wg.Wait()
	fmt.Println("Here we are NOW")
}
