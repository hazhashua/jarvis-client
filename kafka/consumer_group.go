package kafka

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
)

type exampleConsumerGroupHandler struct{}

func (exampleConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (exampleConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (exampleConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// for msg := range claim.Messages() {
	// 	fmt.Printf("Message topic:%q partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
	// 	sess.MarkMessage(msg, "")
	// }
	fmt.Println("consume topic: ", claim.Topic())
	return nil
}

func (exampleConsumerGroupHandler) GetTopicName(claim sarama.ConsumerGroupClaim) {
	fmt.Println("consume group: ", claim.Topic())
}

func NewTestConfig() *sarama.Config {
	configure := sarama.NewConfig()
	// 设置从最头开始读topic数据
	configure.Consumer.Offsets.Initial = sarama.OffsetOldest
	return configure
}

func ConsumeTest() {
	config := NewTestConfig()
	config.Version = sarama.V2_7_0_0 // specify appropriate version
	config.Consumer.Return.Errors = true

	group, err := sarama.NewConsumerGroup([]string{"bigdata-dev01:9092"}, "flinktohbase", config)
	if err != nil {
		panic(err)
	}
	defer func() { _ = group.Close() }()

	// Track errors
	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	// Iterate over consumer sessions.
	ctx := context.Background()
	for {
		topics := []string{"my_topic"}
		handler := exampleConsumerGroupHandler{}

		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			panic(err)
		}

	}
}
