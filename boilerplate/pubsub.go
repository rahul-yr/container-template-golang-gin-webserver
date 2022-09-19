package boilerplate

import (
	"context"
	"errors"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
)

type PubsubClient struct {
	ServiceAccount   string
	GCPProjectId     string
	pubsubClient     *pubsub.Client
	PubsubClientOnce sync.Once
	pubsubTopic      *pubsub.Topic
	PubsubTopicOnce  sync.Once
	TopicName        string
}

func (pb *PubsubClient) Validate() error {
	if pb.ServiceAccount == "" {
		return errors.New("Missing service account for pubsub")
	}
	if pb.GCPProjectId == "" {
		return errors.New("Missing GCP project id for pubsub")
	}
	if pb.TopicName == "" {
		return errors.New("Missing topic name for pubsub")
	}
	return nil
}

// PubsubClient returns a singleton instance of pubsub client
func (pb *PubsubClient) GetSetupPubsubClient() *pubsub.Client {
	pb.PubsubClientOnce.Do(func() {
		if pb.ServiceAccount == "" {
			log.Fatalf("Missing service account for pubsub")
		}
		if pb.GCPProjectId == "" {
			log.Fatalf("Missing GCP project id for pubsub")
		}
		optns := []option.ClientOption{
			option.WithCredentialsFile(pb.ServiceAccount),
		}
		client, err := pubsub.NewClient(context.Background(), pb.GCPProjectId, optns...)
		if err != nil {
			log.Fatalf("pubsub.NewClient: %v", err)
		}
		pb.pubsubClient = client
	})
	return pb.pubsubClient
}

// PubsubTopic returns a singleton instance of pubsub topic
func (pb *PubsubClient) GetSetupPubsubTopic() *pubsub.Topic {
	pb.PubsubTopicOnce.Do(func() {
		client := pb.GetSetupPubsubClient()
		if pb.TopicName == "" {
			log.Fatalf("Missing topic name for pubsub")
		}
		topic := client.Topic(pb.TopicName)
		pb.pubsubTopic = topic
	})
	return pb.pubsubTopic
}

// PublishMessage publishes a message to pubsub
func (pb *PubsubClient) PublishMessage(ctx context.Context, msg *pubsub.Message) error {
	topic := pb.GetSetupPubsubTopic()
	_, err := topic.Publish(ctx, msg).Get(ctx)
	return err
}
