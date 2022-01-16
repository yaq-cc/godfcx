package respmgr

import (
	"context"
	"testing"

	"cloud.google.com/go/firestore"
)

func TestRegister(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "holy-diver-297719")
	if err != nil {
		t.Fail()
	}
	tmp := ResponseTemplateDefinition{
		HandlerName: "HelloHandler",
		MappedVariables: map[string]string{
			"first-name": "firstName",
			"last-name":  "lastName",
		},
		CalculatedVariables: []string{"varOne", "varTwo"},
		ResponseTemplate:    "{firstName} {lastName} {varOne} {varTwo}",
	}
	rtd := ResponseTemplateDefinitions{
		AgentName: "sample-go-test-agent",
		Templates: map[string]ResponseTemplateDefinition{
			"/goodbye": tmp,
		},
	}
	regErr := rtd.Register(ctx, client)
	if regErr != nil {
		t.Fail()
	}
}

func TestRebuild(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "holy-diver-297719")
	if err != nil {
		t.Fail()
	}
	rtd := ResponseTemplateDefinitions{
		AgentName: "sample-go-test-agent",
	}
	rebErr := rtd.Rebuild(ctx, client)
	if rebErr != nil {
		t.Fail()
	}
	for ep, tmp := range rtd.Templates {
		t.Log(ep, " ~ ", tmp)
	}
}

func TestAddTemplateTo(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "holy-diver-297719")
	if err != nil {
		t.Fail()
	}
	rtd := ResponseTemplateDefinitions{
		AgentName: "sample-go-test-agent",
	}
	rebErr := rtd.Rebuild(ctx, client)
	if rebErr != nil {
		t.Fail()
	}

	tmp := ResponseTemplateDefinition{
		Endpoint:    "/rtd-test",
		HandlerName: "HelloHandler",
		MappedVariables: map[string]string{
			"first-name": "firstName",
			"last-name":  "lastName",
		},
		CalculatedVariables: []string{"varOne", "varTwo"},
		ResponseTemplate:    "{firstName} {lastName} {varOne} {varTwo}",
	}

	tmp.AddTemplateTo(&rtd)
	for ep, tmp := range rtd.Templates {
		t.Log(ep, " ~ ", tmp)
	}
}

func TestRegisterAddedTemplates(t *testing.T) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "holy-diver-297719")
	if err != nil {
		t.Fail()
	}
	rtd := ResponseTemplateDefinitions{
		AgentName: "sample-go-test-agent",
	}
	rebErr := rtd.Rebuild(ctx, client)
	if rebErr != nil {
		t.Fail()
	}

	tmp := ResponseTemplateDefinition{
		Endpoint:    "/rtd-test",
		HandlerName: "HelloHandler",
		MappedVariables: map[string]string{
			"first-name": "firstName",
			"last-name":  "lastName",
		},
		CalculatedVariables: []string{"varOne", "varTwo"},
		ResponseTemplate:    "{firstName} {lastName} {varOne} {varTwo}",
	}

	tmp.AddTemplateTo(&rtd)
	for ep, tmp := range rtd.Templates {
		t.Log(ep, " ~ ", tmp)
	}

	// Adding a comment
	regErr := rtd.Register(ctx, client)
	if regErr != nil {
		t.Log(regErr)
		t.Fail()
	}
}
