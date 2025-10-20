# Email Screenshot Generator

A command-line Go application that reads emails from a specified JMAP folder, generates screenshots of each email's content, and automatically moves processed emails to an archive folder.

## Features

- JMAP protocol support for email access (Fastmail)
- Batch processing of multiple emails
- Screenshot generation with customizable dimensions
- Automatic email archiving after processing
- Dry-run mode for testing
- Detailed progress logging

## Prerequisites

- Go 1.19 or later
- Chrome/Chromium browser (for headless screenshot generation)
- Fastmail account with API access

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd aar
```

2. Install dependencies:
```bash
go mod download
```

3. Set the FASTMAIL_AAR_KEY environment variable:
```bash
export FASTMAIL_AAR_KEY=your_fastmail_api_key_here
```

## Getting a Fastmail API Key

1. Log in to your Fastmail account
2. Go to Settings > Privacy & Security
3. Under "API Tokens", click "Create API Token"
4. Give it a name (e.g., "Email Screenshot Generator")
5. Grant the following permissions:
   - Mail: Read-write (for reading and moving emails)
6. Copy the generated token and set it as the FASTMAIL_AAR_KEY environment variable

## Mailbox Setup

The application requires two mailboxes in your Fastmail account:

- **`_aar`** - Source folder containing emails to process
- **`_aar_processed`** - Archive folder for processed emails

Create these mailboxes in Fastmail before running the application.

## Usage

### Basic Usage

Process all emails in the source folder:
```bash
go run .
```

Or build and run:
```bash
go build -o email-screenshot-generator
./email-screenshot-generator
```

### Command-Line Flags

**Limit the number of emails to process:**
```bash
./email-screenshot-generator -limit 10
```

**Dry-run mode (preview without making changes):**
```bash
./email-screenshot-generator -dry-run
```

**Combine flags:**
```bash
./email-screenshot-generator -limit 5 -dry-run
```

## Configuration

The application uses the following configuration:

| Setting | Description | Value |
|---------|-------------|-------|
| `FASTMAIL_AAR_KEY` | Fastmail API key (environment variable, required) | - |
| Screenshot directory | Directory to save screenshots | `./screenshots` |
| Screenshot width | Screenshot width in pixels | `1280` |
| Screenshot height | Screenshot height in pixels | `800` |
| Source folder | Mailbox to read emails from | `_aar` |
| Archive folder | Mailbox to move processed emails to | `_aar_processed` |

Screenshot dimensions and folder names are defined as constants in the code.

## Output

Screenshots are saved to the configured output directory with the naming format:
```
email_<emailID>.png
```

Example output:
```
Starting email screenshot generator...
✓ Connected to JMAP server
Found 5 email(s) in folder '_aar'

Processing email 1/5 (ID: M123abc)...
  Subject: Welcome to our service
  ✓ Screenshot generated: screenshots/email_M123abc.png
  ✓ Moved to archive folder

Processing email 2/5 (ID: M456def)...
  Subject: Your monthly report
  ✓ Screenshot generated: screenshots/email_M456def.png
  ✓ Moved to archive folder

=== Summary ===
Total emails: 5
Successfully processed: 5
Failed: 0
```

## Error Handling

The application will:
- Continue processing remaining emails if one fails
- Log detailed error messages for each failure
- Provide a summary of successes and failures at the end
- Exit with clear error messages for authentication or connection issues

## Development

### Project Structure

```
.
├── main.go           # Main application and orchestration
├── jmap.go           # JMAP client implementation
├── screenshot.go     # Screenshot generation
├── go.mod            # Go module dependencies
└── README.md         # This file
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o email-screenshot-generator
```

## Troubleshooting

**"FASTMAIL_AAR_KEY environment variable is required"**
- Make sure you have set the FASTMAIL_AAR_KEY environment variable with your Fastmail API key

**"Failed to find source folder '_aar'"**
- Create the `_aar` mailbox in your Fastmail account

**"Failed to find archive folder '_aar_processed'"**
- Create the `_aar_processed` mailbox in your Fastmail account

**"Failed to generate screenshot"**
- Ensure Chrome/Chromium is installed on your system
- Check that the HTML content is valid

## License

MIT

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.
