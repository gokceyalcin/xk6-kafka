package kafka

import (
	"crypto/tls"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/segmentio/kafka-go/sasl/aws_msk_iam"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/dop251/goja"
	kafkago "github.com/segmentio/kafka-go"
	"go.k6.io/k6/js/common"
)

type ConnectionConfig struct {
	Address string     `json:"address"`
	SASL    SASLConfig `json:"sasl"`
	TLS     TLSConfig  `json:"tls"`
}

// connectionClass is a constructor for the Connection object in JS
// that creates a new connection for creating, listing and deleting topics,
// e.g. new Connection(...).
// nolint: funlen
func (k *Kafka) connectionClass(call goja.ConstructorCall) *goja.Object {
	runtime := k.vu.Runtime()
	var connectionConfig *ConnectionConfig
	if len(call.Arguments) == 0 {
		common.Throw(runtime, ErrNotEnoughArguments)
	}

	if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
		if b, err := json.Marshal(params); err != nil {
			common.Throw(runtime, err)
		} else {
			if err = json.Unmarshal(b, &connectionConfig); err != nil {
				common.Throw(runtime, err)
			}
		}
	}

	connection := k.getKafkaControllerConnection(connectionConfig)

	connectionObject := runtime.NewObject()
	// This is the connection object itself
	if err := connectionObject.Set("This", connection); err != nil {
		common.Throw(runtime, err)
	}

	err := connectionObject.Set("createTopic", func(call goja.FunctionCall) goja.Value {
		var topicConfig *kafkago.TopicConfig
		if len(call.Arguments) == 0 {
			common.Throw(runtime, ErrNotEnoughArguments)
		}

		if params, ok := call.Argument(0).Export().(map[string]interface{}); ok {
			if b, err := json.Marshal(params); err != nil {
				common.Throw(runtime, err)
			} else {
				if err = json.Unmarshal(b, &topicConfig); err != nil {
					common.Throw(runtime, err)
				}
			}
		}

		k.createTopic(connection, topicConfig)
		return goja.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("deleteTopic", func(call goja.FunctionCall) goja.Value {
		if len(call.Arguments) > 0 {
			if topic, ok := call.Argument(0).Export().(string); !ok {
				common.Throw(runtime, ErrNotEnoughArguments)
			} else {
				k.deleteTopic(connection, topic)
			}
		}

		return goja.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("listTopics", func(call goja.FunctionCall) goja.Value {
		topics := k.listTopics(connection)
		return runtime.ToValue(topics)
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	err = connectionObject.Set("close", func(call goja.FunctionCall) goja.Value {
		if err := connection.Close(); err != nil {
			common.Throw(runtime, err)
		}

		return goja.Undefined()
	})
	if err != nil {
		common.Throw(runtime, err)
	}

	return connectionObject
}

// getKafkaControllerConnection returns a kafka controller connection with a given node address.
// It will also try to use the auth and TLS settings to create a secure connection. The connection
// should be closed after use.
func (k *Kafka) getKafkaControllerConnection(connectionConfig *ConnectionConfig) *kafkago.Conn {
	dialer, wrappedError := GetDialer(connectionConfig.SASL, connectionConfig.TLS)
	if wrappedError != nil {
		logger.WithField("error", wrappedError).Error(wrappedError)
		if dialer == nil {
			common.Throw(k.vu.Runtime(), wrappedError)
			return nil
		}
	}

	isConnectionWithIAM := os.Getenv("CONNECT_AWS_IAM")
	awsRegion := aws.String(os.Getenv("AWS_DEFAULT_REGION"))
	sess, errSession := session.NewSession(&aws.Config{
		Region: awsRegion,
	})
	if errSession != nil {
		log.Panicf("Error creating session: %s", errSession)
	}

	if isConnectionWithIAM != "" {
		dialer.SASLMechanism = &aws_msk_iam.Mechanism{
			Signer: sigv4.NewSigner(sess.Config.Credentials),
			Region: *awsRegion,
		}
		dialer.TLS = &tls.Config{}
	}

	ctx := k.vu.Context()
	if ctx == nil {
		err := NewXk6KafkaError(noContextError, "No context.", nil)
		logger.WithField("error", err).Info(err)
		common.Throw(k.vu.Runtime(), err)
		return nil
	}

	conn, err := dialer.DialContext(ctx, "tcp", connectionConfig.Address)
	if err != nil {
		wrappedError := NewXk6KafkaError(dialerError, "Failed to create dialer.", err)
		logger.WithField("error", wrappedError).Error(wrappedError)
		common.Throw(k.vu.Runtime(), wrappedError)
		return nil
	}

	controller, err := conn.Controller()
	if err != nil {
		wrappedError := NewXk6KafkaError(failedGetController, "Failed to get controller.", err)
		logger.WithField("error", wrappedError).Error(wrappedError)
		common.Throw(k.vu.Runtime(), wrappedError)
		return nil
	}

	controllerConn, err := dialer.DialContext(
		ctx, "tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		wrappedError := NewXk6KafkaError(failedGetController, "Failed to get controller.", err)
		logger.WithField("error", wrappedError).Error(wrappedError)
		common.Throw(k.vu.Runtime(), wrappedError)
		return nil
	}

	return controllerConn
}

// createTopic creates a topic with the given name, partitions, replication factor and compression.
// It will also try to use the auth and TLS settings to create a secure connection. If the topic
// already exists, it will do no-op.
func (k *Kafka) createTopic(conn *kafkago.Conn, topicConfig *kafkago.TopicConfig) {
	if topicConfig.NumPartitions <= 0 {
		topicConfig.NumPartitions = 1
	}

	if topicConfig.ReplicationFactor <= 0 {
		topicConfig.ReplicationFactor = 1
	}

	err := conn.CreateTopics(*topicConfig)
	if err != nil {
		wrappedError := NewXk6KafkaError(failedCreateTopic, "Failed to create topic.", err)
		logger.WithField("error", wrappedError).Error(wrappedError)
		common.Throw(k.vu.Runtime(), wrappedError)
	}
}

// deleteTopic deletes the given topic from the given address. It will also try to
// use the auth and TLS settings to create a secure connection. If the topic
// does not exist, it will raise an error.
func (k *Kafka) deleteTopic(conn *kafkago.Conn, topic string) {
	err := conn.DeleteTopics([]string{topic}...)
	if err != nil {
		wrappedError := NewXk6KafkaError(failedDeleteTopic, "Failed to delete topic.", err)
		logger.WithField("error", wrappedError).Error(wrappedError)
		common.Throw(k.vu.Runtime(), wrappedError)
	}
}

// listTopics lists the topics from the given address. It will also try to
// use the auth and TLS settings to create a secure connection. If the topic
// does not exist, it will raise an error.
func (k *Kafka) listTopics(conn *kafkago.Conn) []string {
	partitions, err := conn.ReadPartitions()
	if err != nil {
		wrappedError := NewXk6KafkaError(failedReadPartitions, "Failed to read partitions.", err)
		logger.WithField("error", wrappedError).Error(wrappedError)
		common.Throw(k.vu.Runtime(), wrappedError)
		return nil
	}

	// There should be a better way to return unique set of
	// topics instead of looping over them twice
	topicSet := map[string]struct{}{}

	for _, partition := range partitions {
		topicSet[partition.Topic] = struct{}{}
	}

	topics := make([]string, 0, len(topicSet))
	for topic := range topicSet {
		topics = append(topics, topic)
	}

	return topics
}
