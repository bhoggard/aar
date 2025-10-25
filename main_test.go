package main

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

// MockEmailClient is a mock implementation of EmailClient
type MockEmailClient struct {
	mailboxes      map[string]*Mailbox
	emails         map[string][]string
	emailDetails   map[string]Email
	moveEmailError error
	getEmailsError error
}

func NewMockEmailClient() *MockEmailClient {
	return &MockEmailClient{
		mailboxes:    make(map[string]*Mailbox),
		emails:       make(map[string][]string),
		emailDetails: make(map[string]Email),
	}
}

func (m *MockEmailClient) FindMailboxByName(name string) (*Mailbox, error) {
	if mailbox, ok := m.mailboxes[name]; ok {
		return mailbox, nil
	}
	return nil, errors.New("mailbox not found")
}

func (m *MockEmailClient) GetEmailsInMailbox(mailboxID string, limit int) ([]string, error) {
	if m.getEmailsError != nil {
		return nil, m.getEmailsError
	}
	if emails, ok := m.emails[mailboxID]; ok {
		if limit > 0 && len(emails) > limit {
			return emails[:limit], nil
		}
		return emails, nil
	}
	return []string{}, nil
}

func (m *MockEmailClient) GetEmails(emailIDs []string) ([]Email, error) {
	var result []Email
	for _, id := range emailIDs {
		if email, ok := m.emailDetails[id]; ok {
			result = append(result, email)
		}
	}
	return result, nil
}

func (m *MockEmailClient) MoveEmail(emailID, sourceMailboxID, targetMailboxID string) error {
	return m.moveEmailError
}

// MockScreenshotService is a mock implementation of ScreenshotService
type MockScreenshotService struct {
	generatedScreenshots map[string]string
	generateError        error
}

func NewMockScreenshotService() *MockScreenshotService {
	return &MockScreenshotService{
		generatedScreenshots: make(map[string]string),
	}
}

func (m *MockScreenshotService) GenerateScreenshot(timestamp, emailID, htmlContent string) (string, error) {
	if m.generateError != nil {
		return "", m.generateError
	}
	path := "screenshots/" + timestamp + "-" + emailID + ".png"
	m.generatedScreenshots[emailID] = path
	return path, nil
}

// Test successful processing of emails
func TestProcessEmails_Success(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()

	// Setup test data
	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}
	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}
	client.emails["src-123"] = []string{"email1", "email2"}
	client.emailDetails["email1"] = Email{
		ID:         "email1",
		Subject:    "Test Email 1",
		ReceivedAt: "2025-10-24T14:30:00Z",
		HTMLBody:   []HTMLBodyPart{{PartID: "part1", Type: "text/html"}},
		BodyValues: map[string]BodyValue{
			"part1": {Value: "<html><body>Test content</body></html>"},
		},
	}
	client.emailDetails["email2"] = Email{
		ID:         "email2",
		Subject:    "Test Email 2",
		ReceivedAt: "2025-10-24T14:35:00Z",
		HTMLBody:   []HTMLBodyPart{{PartID: "part1", Type: "text/html"}},
		BodyValues: map[string]BodyValue{
			"part1": {Value: "<html><body>Test content 2</body></html>"},
		},
	}

	var output bytes.Buffer
	result, err := processEmails(client, generator, 0, false, &output)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.TotalCount != 2 {
		t.Errorf("Expected TotalCount=2, got %d", result.TotalCount)
	}

	if result.ProcessedCount != 2 {
		t.Errorf("Expected ProcessedCount=2, got %d", result.ProcessedCount)
	}

	if result.FailedCount != 0 {
		t.Errorf("Expected FailedCount=0, got %d", result.FailedCount)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Test Email 1") {
		t.Error("Output should contain 'Test Email 1'")
	}
	if !strings.Contains(outputStr, "Test Email 2") {
		t.Error("Output should contain 'Test Email 2'")
	}
}

// Test dry run mode
func TestProcessEmails_DryRun(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()

	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}
	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}
	client.emails["src-123"] = []string{"email1"}

	var output bytes.Buffer
	result, err := processEmails(client, generator, 0, true, &output)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ProcessedCount != 0 {
		t.Errorf("Expected ProcessedCount=0 in dry run, got %d", result.ProcessedCount)
	}

	if result.TotalCount != 1 {
		t.Errorf("Expected TotalCount=1, got %d", result.TotalCount)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "DRY RUN MODE") {
		t.Error("Output should contain 'DRY RUN MODE'")
	}
}

// Test no emails found
func TestProcessEmails_NoEmails(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()

	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}
	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}
	client.emails["src-123"] = []string{}

	var output bytes.Buffer
	result, err := processEmails(client, generator, 0, false, &output)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.TotalCount != 0 {
		t.Errorf("Expected TotalCount=0, got %d", result.TotalCount)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "No emails found") {
		t.Error("Output should contain 'No emails found'")
	}
}

// Test error when source folder not found
func TestProcessEmails_SourceFolderNotFound(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()

	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}

	var output bytes.Buffer
	_, err := processEmails(client, generator, 0, false, &output)

	if err == nil {
		t.Fatal("Expected error when source folder not found")
	}

	if !strings.Contains(err.Error(), "failed to find source folder") {
		t.Errorf("Expected source folder error, got: %v", err)
	}
}

// Test error when archive folder not found
func TestProcessEmails_ArchiveFolderNotFound(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()

	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}

	var output bytes.Buffer
	_, err := processEmails(client, generator, 0, false, &output)

	if err == nil {
		t.Fatal("Expected error when archive folder not found")
	}

	if !strings.Contains(err.Error(), "failed to find archive folder") {
		t.Errorf("Expected archive folder error, got: %v", err)
	}
}

// Test screenshot generation error
func TestProcessEmails_ScreenshotError(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()
	generator.generateError = errors.New("screenshot generation failed")

	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}
	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}
	client.emails["src-123"] = []string{"email1"}
	client.emailDetails["email1"] = Email{
		ID:         "email1",
		Subject:    "Test Email",
		ReceivedAt: "2025-10-24T14:30:00Z",
		HTMLBody:   []HTMLBodyPart{{PartID: "part1", Type: "text/html"}},
		BodyValues: map[string]BodyValue{
			"part1": {Value: "<html><body>Test</body></html>"},
		},
	}

	var output bytes.Buffer
	result, err := processEmails(client, generator, 0, false, &output)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.FailedCount != 1 {
		t.Errorf("Expected FailedCount=1, got %d", result.FailedCount)
	}

	if result.ProcessedCount != 0 {
		t.Errorf("Expected ProcessedCount=0, got %d", result.ProcessedCount)
	}
}

// Test move email error
func TestProcessEmails_MoveEmailError(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()
	client.moveEmailError = errors.New("move failed")

	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}
	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}
	client.emails["src-123"] = []string{"email1"}
	client.emailDetails["email1"] = Email{
		ID:         "email1",
		Subject:    "Test Email",
		ReceivedAt: "2025-10-24T14:30:00Z",
		HTMLBody:   []HTMLBodyPart{{PartID: "part1", Type: "text/html"}},
		BodyValues: map[string]BodyValue{
			"part1": {Value: "<html><body>Test</body></html>"},
		},
	}

	var output bytes.Buffer
	result, err := processEmails(client, generator, 0, false, &output)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.FailedCount != 1 {
		t.Errorf("Expected FailedCount=1, got %d", result.FailedCount)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "Failed to move email to archive") {
		t.Error("Output should contain move error message")
	}
}

// Test email with no HTML content
func TestProcessEmails_NoHTMLContent(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()

	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}
	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}
	client.emails["src-123"] = []string{"email1"}
	client.emailDetails["email1"] = Email{
		ID:         "email1",
		Subject:    "Text Only Email",
		ReceivedAt: "2025-10-24T14:30:00Z",
		HTMLBody:   []HTMLBodyPart{},
	}

	var output bytes.Buffer
	result, err := processEmails(client, generator, 0, false, &output)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.FailedCount != 1 {
		t.Errorf("Expected FailedCount=1, got %d", result.FailedCount)
	}

	outputStr := output.String()
	if !strings.Contains(outputStr, "No HTML content found") {
		t.Error("Output should contain 'No HTML content found'")
	}
}

// Test limit parameter
func TestProcessEmails_WithLimit(t *testing.T) {
	client := NewMockEmailClient()
	generator := NewMockScreenshotService()

	client.mailboxes[sourceFolder] = &Mailbox{ID: "src-123", Name: sourceFolder}
	client.mailboxes[archiveFolder] = &Mailbox{ID: "arch-456", Name: archiveFolder}
	client.emails["src-123"] = []string{"email1", "email2", "email3"}

	for i := 1; i <= 3; i++ {
		id := "email" + string(rune('0'+i))
		client.emailDetails[id] = Email{
			ID:         id,
			Subject:    "Test Email " + string(rune('0'+i)),
			ReceivedAt: "2025-10-24T14:30:00Z",
			HTMLBody:   []HTMLBodyPart{{PartID: "part1", Type: "text/html"}},
			BodyValues: map[string]BodyValue{
				"part1": {Value: "<html><body>Test</body></html>"},
			},
		}
	}

	var output bytes.Buffer
	result, err := processEmails(client, generator, 2, false, &output)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.TotalCount != 2 {
		t.Errorf("Expected TotalCount=2 with limit, got %d", result.TotalCount)
	}
}

// Test extractHTMLContent function
func TestExtractHTMLContent(t *testing.T) {
	tests := []struct {
		name     string
		email    Email
		expected string
	}{
		{
			name: "Valid HTML content",
			email: Email{
				HTMLBody: []HTMLBodyPart{{PartID: "part1"}},
				BodyValues: map[string]BodyValue{
					"part1": {Value: "<html><body>Test</body></html>"},
				},
			},
			expected: "<html><body>Test</body></html>",
		},
		{
			name: "No HTML body",
			email: Email{
				HTMLBody: []HTMLBodyPart{},
			},
			expected: "",
		},
		{
			name: "Missing body value",
			email: Email{
				HTMLBody:   []HTMLBodyPart{{PartID: "part1"}},
				BodyValues: map[string]BodyValue{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractHTMLContent(tt.email)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}
