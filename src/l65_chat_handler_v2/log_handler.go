package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/tmc/langchaingo/callbacks"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

// LogHandler is a callback handler that prints to the standard output.
type LogInBoxHandler struct{}

var _ callbacks.Handler = LogInBoxHandler{}

func (l LogInBoxHandler) HandleLLMGenerateContentStart(_ context.Context, ms []llms.MessageContent) {
	fmt.Println()
	fmt.Println("Entering LLM with messages:")
	for _, m := range ms {
		// TODO: Implement logging of other content types
		var buf strings.Builder
		for _, t := range m.Parts {
			if t, ok := t.(llms.TextContent); ok {
				buf.WriteString(t.Text)
			}
		}
		//fmt.Println("Role:", m.Role)
		//fmt.Println("Text:", buf.String())
		out := fmt.Sprintf("Role: %s\nText: %s\n", m.Role, buf.String())
		printMessageWithBorder(out)

	}
	fmt.Println()
}

func (l LogInBoxHandler) HandleLLMGenerateContentEnd(_ context.Context, res *llms.ContentResponse) {
	fmt.Println()
	fmt.Println("Exiting LLM with response:")
	for _, c := range res.Choices {
		if c.Content != "" {
			//fmt.Println("Content:", c.Content)
			out := fmt.Sprintf("Content: %s", c.Content)
			printMessageWithBorder(out)
		}
		if c.StopReason != "" {
			//fmt.Println("StopReason:", c.StopReason)
			out := fmt.Sprintf("StopReason: %s", c.StopReason)
			printMessageWithBorder(out)
		}
		if len(c.GenerationInfo) > 0 {

			//fmt.Println("GenerationInfo:")
			out := "GenerationInfo:\n"
			for k, v := range c.GenerationInfo {
				//fmt.Printf("%20s: %v\n", k, v)
				out += fmt.Sprintf("%20s: %v\n", k, v)
			}
			printMessageWithBorder(out)
		}
		if c.FuncCall != nil {
			//fmt.Println("FuncCall: ", c.FuncCall.Name, c.FuncCall.Arguments)
			out := fmt.Sprintf("FuncCall: %s %s", c.FuncCall.Name, c.FuncCall.Arguments)
			printMessageWithBorder(out)
		}
	}
}

func (l LogInBoxHandler) HandleStreamingFunc(_ context.Context, chunk []byte) {
	//fmt.Println(string(chunk))
	printMessageWithBorder(string(chunk))
}

func (l LogInBoxHandler) HandleText(_ context.Context, text string) {
	//fmt.Println(text)
	printMessageWithBorder(text)
}

func (l LogInBoxHandler) HandleLLMStart(_ context.Context, prompts []string) {
	//fmt.Println("Entering LLM with prompts:", prompts)
	// add element to prompts at index 0
	prompts = append([]string{"Entering LLM with prompts:"}, prompts...)
	printMessagesWithBorder(prompts)
}

func (l LogInBoxHandler) HandleLLMError(_ context.Context, err error) {
	//fmt.Println("Exiting LLM with error:", err)
	printMessageWithBorder(fmt.Sprintf("Exiting LLM with error: %v", err))
}

func (l LogInBoxHandler) HandleChainStart(_ context.Context, inputs map[string]any) {
	//fmt.Println("Entering chain with inputs:", formatChainValues(inputs))
	printMessageWithBorder(fmt.Sprintf("Entering chain with inputs: %s", formatChainValues(inputs)))
}

func (l LogInBoxHandler) HandleChainEnd(_ context.Context, outputs map[string]any) {
	//fmt.Println("Exiting chain with outputs:", formatChainValues(outputs))
	printMessageWithBorder(fmt.Sprintf("Exiting chain with outputs: %s", formatChainValues(outputs)))
}

func (l LogInBoxHandler) HandleChainError(_ context.Context, err error) {
	//fmt.Println("Exiting chain with error:", err)
	printMessageWithBorder(fmt.Sprintf("Exiting chain with error: %v", err))
}

func (l LogInBoxHandler) HandleToolStart(_ context.Context, input string) {
	//fmt.Println("Entering tool with input:", removeNewLines(input))
	printMessageWithBorder(fmt.Sprintf("Entering tool with input: %s", removeNewLines(input)))
}

func (l LogInBoxHandler) HandleToolEnd(_ context.Context, output string) {
	//fmt.Println("Exiting tool with output:", removeNewLines(output))
	printMessageWithBorder(fmt.Sprintf("Exiting tool with output: %s", removeNewLines(output)))
}

func (l LogInBoxHandler) HandleToolError(_ context.Context, err error) {
	//fmt.Println("Exiting tool with error:", err)
	printMessageWithBorder(fmt.Sprintf("Exiting tool with error: %v", err))
}

func (l LogInBoxHandler) HandleAgentAction(_ context.Context, action schema.AgentAction) {
	//fmt.Println("Agent selected action:", formatAgentAction(action))
	printMessageWithBorder(fmt.Sprintf("Agent selected action: %s", formatAgentAction(action)))
}

func (l LogInBoxHandler) HandleAgentFinish(_ context.Context, finish schema.AgentFinish) {
	//fmt.Printf("Agent finish: %v \n", finish)
	printMessageWithBorder(fmt.Sprintf("Agent finish: %v", finish))
}

func (l LogInBoxHandler) HandleRetrieverStart(_ context.Context, query string) {
	//fmt.Println("Entering retriever with query:", removeNewLines(query))
	printMessageWithBorder(fmt.Sprintf("Entering retriever with query: %s", removeNewLines(query)))
}

func (l LogInBoxHandler) HandleRetrieverEnd(_ context.Context, query string, documents []schema.Document) {
	//fmt.Println("Exiting retriever with documents for query:", documents, query)
	printMessageWithBorder(fmt.Sprintf("Exiting retriever with documents for query: %s, documents: %v", removeNewLines(query), documents))
}

func formatChainValues(values map[string]any) string {
	output := ""
	for key, value := range values {
		output += fmt.Sprintf("\"%s\" : \"%s\", ", removeNewLines(key), removeNewLines(value))
	}

	return output
}

func formatAgentAction(action schema.AgentAction) string {
	return fmt.Sprintf("\"%s\" with input \"%s\"", removeNewLines(action.Tool), removeNewLines(action.ToolInput))
}

func removeNewLines(s any) string {
	return strings.ReplaceAll(fmt.Sprint(s), "\n", " ")
}

func printMessageWithBorder(msg string) {
	var style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	//style = style.SetString("TEST")

	fmt.Println(style.Render(msg))
}

func printMessagesWithBorder(messages []string) {
	var style = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("63"))

	// Join all messages with line breaks
	combinedMessage := strings.Join(messages, "\n")

	// Render the combined message with the border
	fmt.Println(style.Render(combinedMessage))
}
