# Email Screenshot Generator - CLAUDE.md

## Project Overview

A command-line Go application that reads emails from a specified JMAP folder, generates screenshots of each email’s content, and automatically moves processed emails to an archive folder.

## Core Functionality

### Primary Goals

1. **Email Reading**: Connect to a JMAP server and authenticate using provided credentials
1. **Email Processing**: Retrieve emails from a configured source folder
1. **Screenshot Generation**: Convert each email’s HTML content to a PNG screenshot
1. **Email Management**: Move processed emails to a designated archive folder
1. **Progress Tracking**: Display progress and logging for all operations

### Key Features

- JMAP protocol support for email access
- Batch processing of multiple emails
- Configurable source and destination folders
- Screenshot generation with customizable dimensions
- Error handling and retry logic for failed emails
- Detailed logging of all operations

## Technical Architecture

### Dependencies

- `github.com/jmap/jmap-go` - JMAP protocol implementation
- `github.com/chromedp/chromedp` - Headless Chrome for screenshot generation
- `github.com/joho/godotenv` - Environment configuration loading
- Standard library packages: `os`, `fmt`, `flag`, `log`, `net/http`

### Key Components

**JMAP Client**

- Establishes authenticated connection to JMAP server at https://api.fastmail.com using the API key FASTMAIL_AAR_KEY
- Handles account discovery and mailbox selection
- Manages session state

**Email Processor**

- Creates screenshot directory if it does not already exist
- Retrieves email list from source folder named _aar
- Filters emails based on status (unread, unseen, etc.)
- Handles pagination for large email sets

**Screenshot Generator**

- Uses headless Chrome (Chromedp) to render HTML content
- Converts email HTML to PNG screenshots
- Saves screenshots with email ID as filename

**Folder Manager**

- Moves processed emails to archive folder named _aar_processed
- Tracks success/failure of each move operation

## Configuration

### Environment Variables

```
FASTMAIL_AAR_KEY
SCREENSHOT_DIR=./screenshots
SCREENSHOT_WIDTH=1280
SCREENSHOT_HEIGHT=800
```

### Command-Line Flags

```
-limit int         Maximum emails to process (default: 0 = all)
-dry-run bool      Preview operations without making changes
```

## Usage

### Basic Setup

1. Create a `.env` file with JMAP credentials and configuration
1. Build the application: `go build -o email-screenshot-generator`
1. Run with default configuration: `./email-screenshot-generator`
1. Run with custom flags: `./email-screenshot-generator -limit 10`

### Dry-Run Mode

Test your configuration without modifying emails:

```bash
./email-screenshot-generator -dry-run
```

## Data Flow

1. **Authentication** → Connect to JMAP server with credentials
1. **Mailbox Discovery** → Locate source and archive folders by ID
1. **Email Retrieval** → Query JMAP for emails in source folder
1. **Processing Loop**:
- Fetch email metadata and HTML body
- Generate screenshot via Chromedp
- Save screenshot to output directory
- Move email to archive folder
- Log completion status
1. **Completion Report** → Display summary of processed emails

## Error Handling

- Connection failures: Retry with exponential backoff
- Invalid credentials: Exit with clear error message
- Screenshot generation failures: Log error and continue with next email
- Network timeouts: Configurable timeout with retry logic

## Expected Output

```
Starting email screenshot generator...
Connected to JMAP server
Found source folder: Inbox (5 emails)
Processing email 1/5...
  ✓ Screenshot generated: screenshots/email_12345.png
  ✓ Moved to archive folder
Processing email 2/5...
  ✓ Screenshot generated: screenshots/email_12346.png
  ✓ Moved to archive folder
...
Completed: 5 emails processed, 0 failed
```

## Development Notes

### Code Formatting

**IMPORTANT: When working on this project, ALWAYS run `go fmt ./...` after making any changes to Go files.**

- Run `go fmt ./...` before running tests
- Run `go fmt ./...` before committing code
- All Go code must be properly formatted according to Go standards
- Use `go test -v` to verify tests pass after formatting

### Testing Considerations

- Mock JMAP server responses for unit testing
- Use test fixtures with sample HTML emails
- Validate screenshot file creation and naming
- Verify folder move operations

### Performance Optimization

- Implement concurrent email processing with goroutines
- Batch JMAP requests where possible
- Cache mailbox IDs to reduce lookups
- Implement rate limiting to respect server limits

### Security Best Practices

- Never log passwords or credentials
- Use environment variables for sensitive data
- Validate JMAP server SSL certificates
- Implement request signing for JMAP authentication

## Future Enhancements

- Support for multiple email accounts
- Configurable screenshot formats (PDF, JPEG)
- Email filtering by date range or sender
- Webhook notifications on completion
- Database tracking of processed emails
- Web UI for configuration and monitoring