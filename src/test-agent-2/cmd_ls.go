package main

import (
	"context"
	"fmt"

	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/tools"
)

// Linux command to print the current working directory "pwd".
type CmdLs struct {
	CallbacksHandler callbacks.Handler
}

var _ tools.Tool = CmdLs{}

func (c CmdLs) Description() string {
	return `Linux command to list directory contents.`
}

// Name returns the name of the tool.
// Note If I will return "CmdPwd" then is not working
// Note If I will return any Human Readable string related to directory then is not working
// I get the following error:
// unable to parse agent output: I should use the CmdPwd command to find out the current working directory.

func (c CmdLs) Name() string {
	return "Cmd002"
}

func (c CmdLs) Call(ctx context.Context, input string) (string, error) {
	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolStart(ctx, input)
	}

	// cmd := exec.Command("sh", "-c", input)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return "", err
	// }
	// result := string(output)
	fmt.Printf("CmdLs Call input: %s\n", input)
	result := "/tmp/test2"

	if c.CallbacksHandler != nil {
		c.CallbacksHandler.HandleToolEnd(ctx, result)
	}

	return result, nil
}
