package util

import (
	"context"
	"github.com/Shopify/sarama"
	"go.uber.org/zap"
	logf "log"
	"os"
)

type KafkaConsumerProperties struct {
	Brokers []string `json:"brokers" yaml:"brokers"`
	Topics  []string `json:"topics" yaml:"topics"`
	Group   string   `json:"group" yaml:"group"`
	Verbose bool     `json:"verbose" yaml:"verbose"`
	Oldest  bool     `json:"oldest" yaml:"oldest"`
	Version string   `json:"version" yaml:"version"`
}

func NewKafkaConsumer(ctx context.Context, log *zap.Logger, conf KafkaConsumerProperties) *KConsumer {
	if len(conf.Brokers) == 0 || len(conf.Topics) == 0 {
		log.Panic("missing kafka brokers or topics, consumer will not be disabled.")
		return nil
	}
	if conf.Verbose {
		sarama.Logger = logf.New(os.Stdout, "[Sarama] ", logf.LstdFlags)
	}
	ver, err := sarama.ParseKafkaVersion(conf.Version)
	if err != nil {
		log.Panic("Error parsing Kafka version", zap.Error(err))
	}
	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()
	config.Version = ver

	if conf.Oldest {
		config.Consumer.Offsets.Initial = sarama.OffsetOldest
	}

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := &KConsumer{
		stop:  make(chan bool),
		ready: make(chan bool),
		log:   log,
	}

	//ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(conf.Brokers, conf.Group, config)
	if err != nil {
		log.Panic("Error creating consumer group client", zap.Error(err))
	}
	//
	topicArray := conf.Topics
	//
	consumer.mc = make(chan *sarama.ConsumerMessage)
	//
	go func() {
		defer func() {
			if err := client.Close(); err != nil {
				log.Panic("Error closing kafka consumer", zap.Error(err))
			}
			// 没有更多的消息需要处理
			close(consumer.mc)
		}()
		for {
			if err := client.Consume(ctx, topicArray, consumer); err != nil {
				log.Error("Error from consumer", zap.Error(err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()
	return consumer
}

// Consumer represents a Sarama consumer group consumer
type KConsumer struct {
	ready chan bool
	mc    chan *sarama.ConsumerMessage
	log   *zap.Logger
	stop  chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
// Do not invoke this method directly
func (consumer *KConsumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
// Do not invoke this method directly
func (consumer *KConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Do not invoke this method directly
func (consumer *KConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		//
		consumer.mc <- message
		//
		consumer.afterConsume(session, message)
	}
	consumer.stop <- true
	return nil
}

func (consumer *KConsumer) Read() <-chan *sarama.ConsumerMessage {
	return consumer.mc
}

func (consumer *KConsumer) WaitTerminated() {
	<-consumer.stop
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *KConsumer) afterConsume(session sarama.ConsumerGroupSession, message *sarama.ConsumerMessage) {
	session.MarkMessage(message, "")
	consumer.log.Debug("consumed kafka message", zap.String("tag", "KafkaMessage"), zap.String("data", string(message.Value)))
}
