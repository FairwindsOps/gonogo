# Webhook JSON Response Support

This document describes the new functionality for handling JSON responses from webhooks that contain warning messages as a list of strings.

## Overview

In the future, webhook integrations will return JSON responses containing warning messages that need to be processed and displayed. This feature provides backward compatibility while adding support for the new JSON response format.

## Expected JSON Response Format

The webhook is expected to return a JSON response with the following structure:

```json
{
  "addons": [
    {
      "name": "gng-nginx",
      "warnings": [
        "This is a warning",
        "This is another warning",
        "This is a third warning"
      ]
    }
  ],
  "status": "success"
}
```

### Response Fields

- `status`: Status string indicating success or failure (e.g., "success", "error")
- `addons`: Array of addon objects, each containing:
  - `name`: The name of the addon
  - `warnings`: Array of warning message strings for that addon

## Usage

### Command Line Flags

One new flag has been added to support JSON response handling:

- `--output-bundle`: Specify output file path for warning messages

### Examples

#### Basic Usage with JSON Response

```bash
# Send bundle data to webhook and automatically handle JSON response with warning messages
gonogo generate --bundle my-bundle.yaml --webhook https://my-webhook.com/api
```

#### Save Warning Messages to File

```bash
# Send bundle data to webhook and save warning messages to file
gonogo generate --bundle my-bundle.yaml --webhook https://my-webhook.com/api --output-bundle warnings.yaml
```

#### Dry Run with JSON Response

```bash
# Test webhook functionality with automatic JSON response handling
gonogo generate --bundle my-bundle.yaml --webhook https://my-webhook.com/api --dry-run
```

#### Single Release with JSON Response

```bash
# Process single release with automatic JSON response handling
gonogo generate my-release --desired-version 1.2.3 --repo https://my-repo.com --webhook https://my-webhook.com/api
```

## Processing Flow

1. **Send Data**: The tool sends the bundle data to the webhook URL
2. **Receive Response**: The webhook returns a JSON response with addons and their warnings
3. **Validate Status**: Check that the response status is "success"
4. **Process Addons**: Extract and process addons and their warning messages
5. **Output**: The warnings are displayed grouped by addon to stdout or saved to a file

## Error Handling

The tool handles various error scenarios:

- **Empty Response**: Returns error if no response is received from webhook
- **Non-Success Status**: Returns error if webhook returns non-success status
- **No Addons**: Continues normally if no addons are received
- **File Write Errors**: Returns error if unable to write output file

## Backward Compatibility

The new functionality is fully backward compatible:

- Existing webhook usage continues to work without changes
- JSON response handling is automatic - no flags required
- All existing flags and functionality remain unchanged

## Testing

You can test the JSON response functionality using the `--dry-run` flag:

```bash
gonogo generate --bundle my-bundle.yaml --webhook https://my-webhook.com/api --dry-run
```

This will send the bundle data to the webhook and automatically process any JSON response received, without requiring a Kubernetes connection.

## Implementation Details

### Data Structures

```go
// WebhookJSONResponse represents the expected JSON response from webhook
type WebhookJSONResponse struct {
    Addons []WebhookAddon `json:"addons"`
    Status string         `json:"status"`
}

// WebhookAddon represents an addon with warnings in the webhook response
type WebhookAddon struct {
    Name     string   `json:"name"`
    Warnings []string `json:"warnings"`
}
```

### Key Functions

- `sendToWebhookWithResponse()`: Sends data to webhook and optionally expects JSON response
- `processWebhookJSONResponse()`: Processes JSON response and displays warnings grouped by addon
- `convertYAMLToJSON()`: Converts YAML content to JSON for webhook transmission

## Future Enhancements

Potential future enhancements could include:

- Support for different response formats
- Validation of bundle data before processing
- Support for partial updates to existing bundles
- Integration with version control systems
- Support for different output formats (JSON, YAML, etc.)
