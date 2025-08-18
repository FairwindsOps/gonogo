#!/bin/bash

# Example script demonstrating webhook JSON response functionality
# This script shows how to use the automatic JSON response handling for warning messages

echo "🔧 Webhook JSON Response Example"
echo "================================="

# Example 1: Basic usage with JSON response
echo ""
echo "Example 1: Basic usage with JSON response"
echo "-----------------------------------------"
echo "Command: gonogo generate --bundle pkg/bundle/bundles/nginx-ingress.yaml --webhook https://your-webhook.com/api"
echo "This will:"
echo "  - Send the nginx-ingress bundle data to the webhook"
echo "  - Automatically detect and handle JSON response with addons and warnings"
echo "  - Process and display the received warnings grouped by addon"

# Example 2: Save response to file
echo ""
echo "Example 2: Save warning messages to file"
echo "----------------------------------------"
echo "Command: gonogo generate --bundle pkg/bundle/bundles/nginx-ingress.yaml --webhook https://your-webhook.com/api --output-bundle warnings.yaml"
echo "This will:"
echo "  - Send the bundle data to the webhook"
echo "  - Automatically detect and handle JSON response with addons and warnings"
echo "  - Save the warnings grouped by addon to warnings.yaml"

# Example 3: Dry run testing
echo ""
echo "Example 3: Dry run testing"
echo "--------------------------"
echo "Command: gonogo generate --bundle pkg/bundle/bundles/nginx-ingress.yaml --webhook https://your-webhook.com/api --dry-run"
echo "This will:"
echo "  - Test webhook functionality without Kubernetes connection"
echo "  - Send bundle data to webhook"
echo "  - Automatically process any JSON response with addons and warnings"

# Example 4: Single release with JSON response
echo ""
echo "Example 4: Single release with JSON response"
echo "--------------------------------------------"
echo "Command: gonogo generate my-release --desired-version 1.2.3 --repo https://my-repo.com --webhook https://your-webhook.com/api"
echo "This will:"
echo "  - Process a single Helm release"
echo "  - Send release data to webhook"
echo "  - Automatically detect and process JSON response with warning messages"

# Expected JSON response format
echo ""
echo "Expected JSON Response Format"
echo "-----------------------------"
echo "The webhook should return a JSON response like this:"
echo ""
echo '{'
echo '  "addons": ['
echo '    {'
echo '      "name": "gng-nginx",'
echo '      "warnings": ['
echo '        "This is a warning",'
echo '        "This is another warning",'
echo '        "This is a third warning"'
echo '      ]'
echo '    }'
echo '  ],'
echo '  "status": "success"'
echo '}'
echo ""
echo "Where:"
echo "  - status: Status string indicating success or failure"
echo "  - addons: Array of addon objects, each containing:"
echo "    - name: The name of the addon"
echo "    - warnings: Array of warning message strings for that addon"

# Testing with mock webhook
echo ""
echo "Testing with Mock Webhook"
echo "-------------------------"
echo "To test this functionality, you can use a mock webhook service like:"
echo "  - webhook.site"
echo "  - mockbin.org"
echo "  - httpbin.org"
echo ""
echo "Example:"
echo "  gonogo generate --bundle pkg/bundle/bundles/nginx-ingress.yaml --webhook https://webhook.site/your-unique-id --dry-run"

echo ""
echo "✅ Example script completed!"
