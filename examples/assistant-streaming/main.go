package main

import (
	"context"
	"fmt"

	"github.com/joschahenningsen/openai-go"
	"github.com/joschahenningsen/openai-go/option"
)

func main() {
	client := openai.NewClient(
		option.WithHeader("OpenAI-Beta", "assistants=v2"),
	)

	ctx := context.Background()

	// Create an assistant
	println("Create an assistant")
	assistant, err := client.Beta.Assistants.New(ctx, openai.BetaAssistantNewParams{
		Name:         openai.String("Math Tutor"),
		Instructions: openai.String("You are a personal math tutor. Write and run code to answer math questions."),
		Tools: openai.F([]openai.AssistantToolUnionParam{
			openai.CodeInterpreterToolParam{Type: openai.F(openai.CodeInterpreterToolTypeCodeInterpreter)},
		}),
		Model: openai.String("gpt-4-1106-preview"),
	})
	if err != nil {
		panic(err)
	}

	// Create a thread
	println("Create an thread")
	thread, err := client.Beta.Threads.New(ctx, openai.BetaThreadNewParams{})
	if err != nil {
		panic(err)
	}

	// Create a message in the thread
	println("Create a message")
	_, err = client.Beta.Threads.Messages.New(ctx, thread.ID, openai.BetaThreadMessageNewParams{
		Role: openai.F(openai.BetaThreadMessageNewParamsRoleAssistant),
		Content: openai.F([]openai.MessageContentPartParamUnion{
			openai.TextContentBlockParam{
				Type: openai.F(openai.TextContentBlockParamTypeText),
				Text: openai.String("I need to solve the equation `3x + 11 = 14`. Can you help me?"),
			},
		}),
	})
	if err != nil {
		panic(err)
	}

	// Create a run
	println("Create a run")
	stream := client.Beta.Threads.Runs.NewStreaming(ctx, thread.ID, openai.BetaThreadRunNewParams{
		AssistantID:  openai.String(assistant.ID),
		Instructions: openai.String("Please address the user as Jane Doe. The user has a premium account."),
	})
	if err != nil {
		panic(err)
	}

	for stream.Next() {
		evt := stream.Current()
		println(fmt.Sprintf("%T", evt.Data))
	}
}
