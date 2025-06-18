package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/tools"
	"os"
)

func createHTMLReport(htmlContent string) (string, error) {
	fileName := "report.html"
	file, err := os.Create(fileName)
	if err != nil {
		return "", fmt.Errorf("error creating HTML file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(htmlContent)
	if err != nil {
		return "", fmt.Errorf("error writing to HTML file: %w", err)
	}
	return "HTML file successfully created", nil
}

type WriteHtmlReport struct {
}

var _ tools.Tool = WriteHtmlReport{}

// Description returns a string describing the calculator tool.
func (c WriteHtmlReport) Description() string {
	return "Given an html string will create an html file to disk. Use this tool whenever someone asks for a report."
}

// Name returns the name of the tool.
func (c WriteHtmlReport) Name() string {
	return "write_html_report"
}

func (c WriteHtmlReport) Call(ctx context.Context, input string) (string, error) {
	return createHTMLReport(input)
}
