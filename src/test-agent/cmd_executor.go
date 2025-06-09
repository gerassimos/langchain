package main

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
)

// CmdExecutor is a tool that can do math.
type CmdExecutor struct {
	CallbacksHandler callbacks.Handler
}

var _ tools.Tool = CmdExecutor{}

func (c CmdExecutor) Description() string {
	return `Linux command to get the current working directory.`
}

// Name returns the name of the tool.
func (c CmdExecutor) Name() string {
	return "CmdExecutor"
}

func (c CmdExecutor) Call(ctx context.Context, input string) (string, error) {
	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// cmd := exec.Command("sh", "-c", input)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return "", err
	// }
	// result := string(output)
	fmt.Printf("CmdExecutor Call input: %s\n", input)
	result := "/tmp/test1"

	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}
