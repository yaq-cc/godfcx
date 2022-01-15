package godfcx

import (
	"encoding/json"
	"net/http"
)

type IntentInfo struct {
	LastMatchedIntent string  `json:"lastMatchedIntent,omitempty"`
	DisplayName       string  `json:"displayName,omitempty"`
	Confidence        float64 `json:"confidence,omitempty"`
}

type PageInfo struct {
	CurrentPage string `json:"currentPage,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type SessionInfo struct {
	Session    string            `json:"session,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type FulfillmentInfo struct {
	Tag string `json:"tag,omitempty"`
}

type Text struct {
	Text                      []string `json:"text,omitempty"`
	RedactedText              []string `json:"redactedText,omitempty"`
	AllowPlaybackInterruption bool     `json:"allowPlaybackInterruption,omitempty"`
}

func (t *Text) isMessage() {}

type Source struct {
	Text string `json:"text,omitempty"`
	SSML string `json:"ssml,omitempty"`
}

type OutputAudioText struct {
	AllowPlaybackInterruption bool `json:"allowPlaybackInterruption,omitempty"`
	Source                    Source
}

func (o *OutputAudioText) isMessage() {}

type Messages struct {
	Text         Text   `json:"text,omitempty"`
	ResponseType string `json:"responseType,omitempty"`
	Source       string `json:"source,omitempty"`
}

type WebhookRequest struct {
	DetectIntentResponseID string          `json:"detectIntentResponseId,omitempty"`
	IntentInfo             IntentInfo      `json:"intentInfo,omitempty"`
	PageInfo               PageInfo        `json:"pageInfo,omitempty"`
	SessionInfo            SessionInfo     `json:"sessionInfo,omitempty"`
	FulfillmentInfo        FulfillmentInfo `json:"fulfillmentInfo,omitempty"`
	Messages               []Messages      `json:"messages,omitempty"`
	Text                   string          `json:"text,omitempty"`
	LanguageCode           string          `json:"languageCode,omitempty"`
}

func (wr *WebhookRequest) FromRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(wr)
	if err != nil {
		return err
	}
	return nil
}

type Message interface {
	isMessage()
}

type FulfillmentResponse struct {
	// https://pkg.go.dev/google.golang.org/genproto/googleapis/cloud/dialogflow/cx/v3beta1#ResponseMessage
	Messages      []Message
	MergeBehavior int
}

type WebhookResponse struct {
	FulfillmentResponse FulfillmentResponse `json:"fulfillmentResponse,omitempty"`
	PageInfo            PageInfo            `json:"pageInfo,omitempty"`
	SessionInfo         SessionInfo         `json:"sessionInfo,omitempty"`
	Payload             map[string]string   `json:"payload,omitempty"`
}

func (wr *WebhookResponse) TextResponse(w http.ResponseWriter, msgs ...string) {
	t := Text{
		Text:                      msgs,
		AllowPlaybackInterruption: true,
	}

	r := WebhookResponse{
		FulfillmentResponse: FulfillmentResponse{
			Messages:      []Message{&t},
			MergeBehavior: 0,
		},
	}
	json.NewEncoder(w).Encode(r)
}

func (wr *WebhookResponse) SSMLResponse(w http.ResponseWriter, msg string) {
	t := OutputAudioText{
		AllowPlaybackInterruption: true,
		Source: Source{
			SSML: msg,
		},
	}

	r := WebhookResponse{
		FulfillmentResponse: FulfillmentResponse{
			Messages:      []Message{&t},
			MergeBehavior: 0,
		},
	}
	json.NewEncoder(w).Encode(r)
}