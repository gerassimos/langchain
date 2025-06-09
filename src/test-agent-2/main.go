package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("Start")
	// We can construct an LLMChain from a PromptTemplate and an LLM.
	llm, err := openai.New()
	if err != nil {
		return err
	}

	agentTools := []tools.Tool{
		CmdPwd{},
		CmdLs{},
	}

	//This will actually create the agent and assign it to the executor
	//View Udemy S6L53 03:34
	// executor, err := agents.Initialize(
	// 	llm,
	// 	agentTools,
	// 	agents.ZeroShotReactDescription,
	// 	agents.WithMaxIterations(3),
	// )
	// if err != nil {
	// 	return err
	// }

	agent := agents.NewOneShotAgent(llm, agentTools)
	// agent := agents.NewOpenAIFunctionsAgent(llm, agentTools)

	executor := agents.NewExecutor(
		agent,
		agents.WithMaxIterations(3),
	)
	if err != nil {
		return err
	}

	// Prompt the agent
	prompt := "What is the current working directory?"
	// output, err := agents.Run(context.Background(), agent, prompt)
	// if err != nil {
	// 	return err
	// }
	answer, err := chains.Run(context.Background(), executor, prompt)
	if err != nil {
		return err
	}
	fmt.Println(answer)

	prompt = "Which are the files in the current working directory?"
	// output, err := agents.Run(context.Background(), agent, prompt)
	// if err != nil {
	// 	return err
	// }
	answer, err = chains.Run(context.Background(), executor, prompt)
	if err != nil {
		return err
	}

	fmt.Println(answer)
	return nil
}
