package dify

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Metadata struct {
	Usage struct {
		PromptTokens        int     `json:"prompt_tokens"`
		PromptUnitPrice     string  `json:"prompt_unit_price"`
		PromptPriceUnit     string  `json:"prompt_price_unit"`
		PromptPrice         string  `json:"prompt_price"`
		CompletionTokens    int     `json:"completion_tokens"`
		CompletionUnitPrice string  `json:"completion_unit_price"`
		CompletionPriceUnit string  `json:"completion_price_unit"`
		CompletionPrice     string  `json:"completion_price"`
		TotalTokens         int     `json:"total_tokens"`
		TotalPrice          string  `json:"total_price"`
		Currency            string  `json:"currency"`
		Latency             float64 `json:"latency"`
	} `json:"usage"`

	RetrieverResource struct {
		Position     int     `json:"position"`
		DatasetID    string  `json:"dataset_id"`
		DatasetName  string  `json:"dataset_name"`
		DocumentID   string  `json:"document_id"`
		DocumentName string  `json:"document_name"`
		SegmentID    string  `json:"segment_id"`
		Score        float64 `json:"score"`
		Content      string  `json:"content"`
	} `json:"retriever_resources"`
}

type API struct {
	c      *Client
	secret string
}

func (api *API) WithSecret(secret string) *API {
	api.secret = secret
	return api
}

func (api *API) getSecret() string {
	if api.secret != "" {
		return api.secret
	}
	return api.c.getAPISecret()
}

func (api *API) createBaseRequest(ctx context.Context, method, apiUrl string, body interface{}) (*http.Request, error) {
	var b io.Reader
	if body != nil {
		reqBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		b = bytes.NewBuffer(reqBytes)
	} else {
		b = http.NoBody
	}
	req, err := http.NewRequestWithContext(ctx, method, api.c.getHost()+apiUrl, b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+api.getSecret())
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	return req, nil
}
