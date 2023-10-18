package main

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type GenerateContainer struct {
	Token string
}

// FIXME: should be a secret as soon as supported
func initOpenAIClient(token string) *openai.Client {
	return openai.NewClient(token)
}

// Excract code from openai response
func cleanResponse(resp openai.ChatCompletionResponse) string {
	re := regexp.MustCompile("```.*\n(.+)```")
	content := resp.Choices[0].Message.Content
	match := re.FindStringSubmatch(content)
	if len(match) > 0 {
		content = match[0]
	}
	return strings.Trim(content, " \n")
}

func (ctr *Container) Generate(ctx context.Context, token, prompt string) (*Container, error) {
	c := initOpenAIClient(token)

	resp, err := c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Temperature: 0.0,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You will only generate a single Dockerfile based on the user request,
				only include the Dockerfile content in your response.
				The Dockerfile always installs bash as a package dependency.
				Do not include instructions on how to use it, just print the content of the Dockerfile.`,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	content := cleanResponse(resp)
	fmt.Printf("Dockerfile:\n---\n%s\n---\n", content)

	// TODO: support passed context dir as argument
	contextDir := dag.Directory()

	return contextDir.
		WithNewFile("Dockerfile", content).
		DockerBuild(DirectoryDockerBuildOpts{}), nil
}

func (ctr *Container) GenerateBashScript(ctx context.Context, token, prompt string) (*Container, error) {
	c := initOpenAIClient(token)

	resp, err := c.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:       openai.GPT3Dot5Turbo,
		Temperature: 0.0,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: `You will only generate a single Bash script based on the user request,
				only include the script content in your response.
				Do not include instructions on how to use it, just print the content of the script.`,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	content := cleanResponse(resp)
	fmt.Printf("Bash script:\n---\n%s\n---\n", content)

	return ctr.WithNewFile("script.sh", ContainerWithNewFileOpts{Contents: content}).
		WithExec([]string{"/bin/bash", "script.sh"}).
		Sync(ctx)
}
