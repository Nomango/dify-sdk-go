package test

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"testing"

	"github.com/langgenius/dify-sdk-go"
)

var (
	host         = "https://dify.labex.dev"
	apiSecretKey = "app-JjNGGcKRFxgLu6Kfg4tIALDl"
)

func TestAPI(t *testing.T) {
	var c = &dify.ClientConfig{
		Host:         host,
		ApiSecretKey: apiSecretKey,
	}
	var client = dify.NewClientWithConfig(c)

	ctx := context.Background()

	stream, err := client.Api().ChatMessagesStream(ctx, &dify.ChatMessageRequest{
		Query: "你是谁?",
		User:  "user-123",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	defer stream.Close()

	for {
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
			return
		default:
			resp, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				t.Fatal(err.Error())
			}
			content, _ := json.Marshal(resp)
			log.Println("Event", string(content))
		}
	}
}

func TestMessages(t *testing.T) {
	var cId = "ec373942-2d17-4f11-89bb-f9bbf863ebcc"
	var err error
	ctx := context.Background()

	// messages
	var messageReq = &dify.MessagesRequest{
		ConversationID: cId,
		User:           "jiuquan AI",
	}

	var client = dify.NewClient(host, apiSecretKey)

	var msg *dify.MessagesResponse
	if msg, err = client.Api().Messages(ctx, messageReq); err != nil {
		t.Fatal(err.Error())
		return
	}
	j, _ := json.Marshal(msg)
	t.Log(string(j))
}

func TestMessagesFeedbacks(t *testing.T) {
	var client = dify.NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var id = "72d3dc0f-a6d5-4b5e-8510-bec0611a6048"

	var res *dify.MessagesFeedbacksResponse
	if res, err = client.Api().MessagesFeedbacks(ctx, &dify.MessagesFeedbacksRequest{
		MessageID: id,
		Rating:    dify.FeedbackLike,
		User:      "jiuquan AI",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}

func TestConversations(t *testing.T) {
	var client = dify.NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var res *dify.ConversationsResponse
	if res, err = client.Api().Conversations(ctx, &dify.ConversationsRequest{
		User: "jiuquan AI",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}

func TestConversationsRename(t *testing.T) {
	var client = dify.NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var res *dify.ConversationsRenamingResponse
	if res, err = client.Api().ConversationsRenaming(ctx, &dify.ConversationsRenamingRequest{
		ConversationID: "ec373942-2d17-4f11-89bb-f9bbf863ebcc",
		Name:           "rename!!!",
		User:           "jiuquan AI",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}

func TestParameters(t *testing.T) {
	var client = dify.NewClient(host, apiSecretKey)
	var err error
	ctx := context.Background()

	var res *dify.ParametersResponse
	if res, err = client.Api().Parameters(ctx, &dify.ParametersRequest{
		User: "jiuquan AI",
	}); err != nil {
		t.Fatal(err.Error())
	}

	j, _ := json.Marshal(res)

	log.Println(string(j))
}
