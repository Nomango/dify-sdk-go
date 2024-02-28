package dify

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatMessageStreamResponse struct {
	Raw       []byte `json:"-"`
	Event     string `json:"event"`
	TaskID    string `json:"task_id"`
	MessageID string `json:"message_id"`

	Data struct {
		Message        *ChatMessageStreamEventMessage        `json:"message,omitempty"`
		AgentMessage   *ChatMessageStreamEventAgentMessage   `json:"agent_message,omitempty"`
		AgentThought   *ChatMessageStreamEventAgentThought   `json:"agent_thought,omitempty"`
		MessageFile    *ChatMessageStreamEventMessageFile    `json:"message_file,omitempty"`
		MessageReplace *ChatMessageStreamEventMessageReplace `json:"message_replace,omitempty"`
		MessageEnd     *ChatMessageStreamEventMessageEnd     `json:"message_end,omitempty"`
		Error          *ChatMessageStreamEventError          `json:"error,omitempty"`
	} `json:"_data"`
}

type ChatMessageStreamEventMessage struct {
	ConversationID string `json:"conversation_id"`
	Answer         string `json:"answer"`
	CreatedAt      int64  `json:"created_at"`
}

type ChatMessageStreamEventAgentMessage struct {
	ConversationID string `json:"conversation_id"`
	Answer         string `json:"answer"`
	CreatedAt      int64  `json:"created_at"`
}

type ChatMessageStreamEventAgentThought struct {
	ID             string   `json:"id"`
	ConversationID string   `json:"conversation_id"`
	Position       int64    `json:"position"`
	Thought        string   `json:"thought"`
	Observation    string   `json:"observation"`
	Tool           string   `json:"tool"`
	ToolInput      string   `json:"tool_input"`
	MessageFiles   []string `json:"message_files"`
	CreatedAt      int64    `json:"created_at"`
}

type ChatMessageStreamEventMessageFile struct {
	ID             string `json:"id"`
	ConversationID string `json:"conversation_id"`
	Type           string `json:"type"`
	BelongsTo      string `json:"belongs_to"`
	URL            string `json:"url"`
}

type ChatMessageStreamEventMessageReplace struct {
	ConversationID string `json:"conversation_id"`
	Answer         string `json:"answer"`
	CreatedAt      int64  `json:"created_at"`
}

type ChatMessageStreamEventMessageEnd struct {
	Metadata Metadata `json:"metadata"`
}

type ChatMessageStreamEventError struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type StreamReader[T any] interface {
	Recv() (T, error)
	Close()
}

func (api *API) ChatMessagesStream(ctx context.Context, req *ChatMessageRequest) (StreamReader[ChatMessageStreamResponse], error) {
	req.ResponseMode = "streaming"
	httpReq, err := api.createBaseRequest(ctx, http.MethodPost, "/v1/chat-messages", req)
	if err != nil {
		return nil, err
	}

	httpResp, err := api.c.sendRequest(httpReq)
	if err != nil {
		return nil, err
	}

	return &chatMessageStream{
		httpResp: httpResp,
		reader:   bufio.NewReader(httpResp.Body),
	}, nil
}

type chatMessageStream struct {
	httpResp   *http.Response
	reader     *bufio.Reader
	isFinished bool
}

func (s *chatMessageStream) Recv() (resp ChatMessageStreamResponse, err error) {
	if s.isFinished {
		err = io.EOF
		return
	}
	for {
		rawLine, readErr := s.reader.ReadBytes('\n')
		if readErr != nil {
			err = fmt.Errorf("error reading chat message: %w", readErr)
			return
		}
		noSpaceLine := bytes.TrimSpace(rawLine)
		if !bytes.HasPrefix(noSpaceLine, []byte("data: ")) {
			continue
		}
		noHeaderLine := bytes.TrimPrefix(noSpaceLine, []byte("data: "))
		if !bytes.HasPrefix(noHeaderLine, []byte(`{"event":`)) {
			err = fmt.Errorf("error chat message event: %s", string(noHeaderLine))
			return
		}
		resp.Raw = noHeaderLine
		if unmarshalErr := json.Unmarshal(noHeaderLine, &resp); unmarshalErr != nil {
			err = fmt.Errorf("error unmarshal chat message: %w", unmarshalErr)
			return
		}
		var data any
		switch resp.Event {
		case "message":
			data = &resp.Data.Message
		case "agent_message":
			data = &resp.Data.AgentMessage
		case "agent_thought":
			data = &resp.Data.AgentThought
		case "message_file":
			data = &resp.Data.MessageFile
		case "message_replace":
			data = &resp.Data.MessageReplace
		case "message_end":
			data = &resp.Data.MessageEnd
			s.isFinished = true
		case "error":
			data = &resp.Data.Error
		case "ping":
			continue
		default:
			return
		}
		if unmarshalErr := json.Unmarshal(noHeaderLine, data); unmarshalErr != nil {
			err = fmt.Errorf("error unmarshal chat message data: %w", unmarshalErr)
			return
		}
		return
	}
}

func (s *chatMessageStream) Close() {
	s.httpResp.Body.Close()
}
