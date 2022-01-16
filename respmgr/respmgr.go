package respmgr

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

type ResponseTemplateDefinition struct {
	Endpoint            string            `firestore:"-"`
	HandlerName         string            `firestore:"handler-name,omitempty"`
	MappedVariables     map[string]string `firestore:"mapped-variables,omitempty"`
	CalculatedVariables []string          `firestore:"calculated-variables,omitempty"`
	ResponseTemplate    string            `firestore:"response-template,omitempty"`
}

func (rt *ResponseTemplateDefinition) AddTemplateTo(rtd *ResponseTemplateDefinitions) {
	rtd.Templates[rt.Endpoint] = *rt
}

type ResponseTemplateDefinitions struct {
	AgentName string `firestore:"-"`
	// Since Templates is a map, it doesn't get tags.
	Templates map[string]ResponseTemplateDefinition
}

func (rtd *ResponseTemplateDefinitions) Rebuild(ctx context.Context, client *firestore.Client) error {
	snap, err := client.Collection("dialogflow-agents").Doc(rtd.AgentName).Get(ctx)
	if err != nil {
		return err
	}
	snap.DataTo((&rtd.Templates))
	return nil
}

// Synchronizes the ResponseTemplateDefinitions from Firestore.  
func (rtd *ResponseTemplateDefinitions) Initialize(ctx context.Context, client *firestore.Client) {
	err := rtd.Rebuild(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
}

// Registers all response templates for a given agent
func (rtd *ResponseTemplateDefinitions) Register(ctx context.Context, client *firestore.Client) error {
	_, err := client.Collection("dialogflow-agents").Doc(rtd.AgentName).Set(ctx, rtd.Templates)
	if err != nil {
		return err
	}
	return nil
}
