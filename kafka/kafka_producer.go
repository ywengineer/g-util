package util

import (
	"github.com/Shopify/sarama"
	"github.com/ywengineer/g-util/util"
	"go.uber.org/zap"
	"time"
)

type KafkaProducerProperties struct {
	Brokers           []string `json:"brokers" yaml:"brokers"`
	CertFile          string   `json:"cert_file" yaml:"cert-file"`
	KeyFile           string   `json:"key_file" yaml:"key-file"`
	CaFile            string   `json:"ca_file" yaml:"ca-file"`
	VerifySSL         bool     `json:"verify_ssl" yaml:"verify-ssl"`
	Compression       string   `json:"compression" yaml:"compression"`
	IgnoreSendSuccess bool     `json:"ignore_send_success" yaml:"ignore-send-success"`
}

//"none",
//"gzip",
//"snappy",
//"lz4",
//"zstd",
func checkCompression(compression string) sarama.CompressionCodec {
	switch compression {
	case "none":
		return sarama.CompressionNone
	case "gzip":
		return sarama.CompressionGZIP
	case "snappy":
		return sarama.CompressionSnappy
	case "lz4":
		return sarama.CompressionLZ4
	case "zstd":
		return sarama.CompressionZSTD
	default:
		util.Panic("unsupported kafka's compression: %s. %s", compression, `"none","gzip","snappy","lz4","zstd"`)
	}
	return sarama.CompressionNone
}

func NewKafkaSyncProducerFromConf(conf KafkaProducerProperties) sarama.SyncProducer {
	return NewKafkaSyncProducer(conf.Brokers, &conf.CertFile, &conf.KeyFile, &conf.CaFile, conf.VerifySSL, checkCompression(conf.Compression))
}

func NewKafkaAsyncProducerFromConf(conf KafkaProducerProperties, errLog *zap.Logger) sarama.AsyncProducer {
	return NewKafkaAsyncProducer(conf.Brokers, &conf.CertFile, &conf.KeyFile, &conf.CaFile, conf.VerifySSL, conf.IgnoreSendSuccess, errLog, checkCompression(conf.Compression))
}

func NewKafkaSyncProducer(brokerList []string, certFile, keyFile, caFile *string, verifySsl bool, compression sarama.CompressionCodec) sarama.SyncProducer {
	if len(brokerList) == 0 {
		util.Panic("no brokers for kafka sync sender")
	}
	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	//
	config.Producer.Compression = compression
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 2                    // Retry up to 10 times to produce the message
	//
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Producer.Idempotent = true
	//
	tlsConfig := util.CreateTlsConfiguration(certFile, keyFile, caFile, verifySsl)
	if tlsConfig != nil {
		config.Net.TLS.Config = tlsConfig
		config.Net.TLS.Enable = true
	}
	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.
	producer, err := sarama.NewSyncProducer(brokerList, config)
	//
	if err != nil {
		util.Panic("Failed to start Kafka producer: %v", err)
	}
	//
	return producer
}

func NewKafkaAsyncProducer(brokerList []string, certFile, keyFile, caFile *string,
	verifySsl, ignoreSendSuccess bool,
	errorLogger *zap.Logger,
	compression sarama.CompressionCodec) sarama.AsyncProducer {
	if len(brokerList) == 0 {
		util.Panic("no brokers for kafka async sender")
	}
	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	//
	tlsConfig := util.CreateTlsConfiguration(certFile, keyFile, caFile, verifySsl)
	if tlsConfig != nil {
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	//
	config.Producer.Idempotent = true
	config.Producer.RequiredAcks = sarama.WaitForLocal // Only wait for the leader to ack
	config.Producer.Compression = compression          // Compress messages
	config.Producer.Flush.Frequency = 1 * time.Second  // Flush batches every 1s

	producer, err := sarama.NewAsyncProducer(brokerList, config)

	if err != nil {
		util.Panic("Failed to start Kafka producer: %v", err)
	}
	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			errorLogger.Error("Failed to write entry to kafka broker", zap.Error(err), zap.String("tag", "KafkaAsyncProducer"))
		}
	}()
	//
	if ignoreSendSuccess {
		go func() {
			for range producer.Successes() {
			}
		}()
	}
	return producer
}
