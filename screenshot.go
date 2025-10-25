package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

// ScreenshotGenerator handles screenshot generation
type ScreenshotGenerator struct {
	outputDir string
	width     int
	height    int
}

// NewScreenshotGenerator creates a new screenshot generator
func NewScreenshotGenerator(outputDir string, width, height int) (*ScreenshotGenerator, error) {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	return &ScreenshotGenerator{
		outputDir: outputDir,
		width:     width,
		height:    height,
	}, nil
}

// GenerateScreenshot creates a screenshot from HTML content
func (s *ScreenshotGenerator) GenerateScreenshot(timestamp, htmlContent string) (string, error) {
	// Parse the timestamp
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return "", fmt.Errorf("failed to parse timestamp: %w", err)
	}

	// Format timestamp as yyyy-mm-dd-hh-mm-ss
	formattedTime := t.Format("2006-01-02-15-04-05")

	// Create output filename
	outputPath := filepath.Join(s.outputDir, fmt.Sprintf("%s.png", formattedTime))

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create chromedp context
	allocCtx, allocCancel := chromedp.NewContext(ctx)
	defer allocCancel()

	// Prepare HTML with base structure
	fullHTML := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body {
            margin: 20px;
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            font-size: 14px;
            line-height: 1.5;
        }
        img {
            max-width: 100%%;
            height: auto;
        }
    </style>
</head>
<body>
%s
</body>
</html>`, htmlContent)

	// Create a data URL from the HTML
	dataURL := "data:text/html;charset=utf-8," + url.PathEscape(fullHTML)

	// Run chromedp tasks
	var buf []byte
	if err := chromedp.Run(allocCtx,
		chromedp.EmulateViewport(int64(s.width), int64(s.height)),
		chromedp.Navigate(dataURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(500*time.Millisecond), // Give time for rendering
		chromedp.FullScreenshot(&buf, 90),
	); err != nil {
		return "", fmt.Errorf("failed to generate screenshot: %w", err)
	}

	// Write screenshot to file
	if err := os.WriteFile(outputPath, buf, 0644); err != nil {
		return "", fmt.Errorf("failed to write screenshot: %w", err)
	}

	return outputPath, nil
}
