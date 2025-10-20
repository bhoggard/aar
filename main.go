package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	sourceFolder     = "_aar"
	archiveFolder    = "_aar_processed"
	screenshotDir    = "./screenshots"
	screenshotWidth  = 1280
	screenshotHeight = 800
)

var (
	limit  = flag.Int("limit", 0, "Maximum emails to process (default: 0 = all)")
	dryRun = flag.Bool("dry-run", false, "Preview operations without making changes")
)

func main() {
	flag.Parse()

	// Get API key from environment
	apiKey := os.Getenv("FASTMAIL_AAR_KEY")
	if apiKey == "" {
		log.Fatal("FASTMAIL_AAR_KEY environment variable is required")
	}

	fmt.Println("Starting email screenshot generator...")

	// Create JMAP client
	client, err := NewJMAPClient(apiKey)
	if err != nil {
		log.Fatalf("Failed to create JMAP client: %v", err)
	}
	fmt.Println("✓ Connected to JMAP server")

	// Find source mailbox
	sourceMailbox, err := client.FindMailboxByName(sourceFolder)
	if err != nil {
		log.Fatalf("Failed to find source folder '%s': %v", sourceFolder, err)
	}

	// Find archive mailbox
	archiveMailbox, err := client.FindMailboxByName(archiveFolder)
	if err != nil {
		log.Fatalf("Failed to find archive folder '%s': %v", archiveFolder, err)
	}

	// Get emails from source folder
	emailIDs, err := client.GetEmailsInMailbox(sourceMailbox.ID, *limit)
	if err != nil {
		log.Fatalf("Failed to retrieve emails: %v", err)
	}

	emailCount := len(emailIDs)
	if emailCount == 0 {
		fmt.Printf("No emails found in folder '%s'\n", sourceFolder)
		return
	}

	fmt.Printf("Found %d email(s) in folder '%s'\n", emailCount, sourceFolder)

	if *dryRun {
		fmt.Println("\nDRY RUN MODE - No changes will be made")
		fmt.Printf("Would process %d emails:\n", emailCount)
		for i, id := range emailIDs {
			fmt.Printf("  %d. Email ID: %s\n", i+1, id)
		}
		return
	}

	// Create screenshot generator
	generator, err := NewScreenshotGenerator(screenshotDir, screenshotWidth, screenshotHeight)
	if err != nil {
		log.Fatalf("Failed to create screenshot generator: %v", err)
	}

	// Process emails
	var processedCount, failedCount int
	for i, emailID := range emailIDs {
		fmt.Printf("\nProcessing email %d/%d (ID: %s)...\n", i+1, emailCount, emailID)

		// Get email details
		emails, err := client.GetEmails([]string{emailID})
		if err != nil {
			log.Printf("  ✗ Failed to fetch email: %v", err)
			failedCount++
			continue
		}

		if len(emails) == 0 {
			log.Printf("  ✗ Email not found")
			failedCount++
			continue
		}

		email := emails[0]
		fmt.Printf("  Subject: %s\n", email.Subject)

		// Extract HTML content
		htmlContent := extractHTMLContent(email)
		if htmlContent == "" {
			log.Printf("  ✗ No HTML content found")
			failedCount++
			continue
		}

		// Generate screenshot
		screenshotPath, err := generator.GenerateScreenshot(emailID, htmlContent)
		if err != nil {
			log.Printf("  ✗ Failed to generate screenshot: %v", err)
			failedCount++
			continue
		}
		fmt.Printf("  ✓ Screenshot generated: %s\n", screenshotPath)

		// Move email to archive folder
		if err := client.MoveEmail(emailID, sourceMailbox.ID, archiveMailbox.ID); err != nil {
			log.Printf("  ✗ Failed to move email to archive: %v", err)
			failedCount++
			continue
		}
		fmt.Printf("  ✓ Moved to archive folder\n")

		processedCount++
	}

	// Print summary
	fmt.Printf("\n" + "=== Summary ===\n")
	fmt.Printf("Total emails: %d\n", emailCount)
	fmt.Printf("Successfully processed: %d\n", processedCount)
	fmt.Printf("Failed: %d\n", failedCount)
}

// extractHTMLContent extracts HTML content from an email
func extractHTMLContent(email Email) string {
	if len(email.HTMLBody) == 0 {
		return ""
	}

	// Get the first HTML body part
	partID := email.HTMLBody[0].PartID

	// Get the body value
	if bodyValue, ok := email.BodyValues[partID]; ok {
		return bodyValue.Value
	}

	return ""
}
