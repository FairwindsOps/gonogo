---
meta:
  - name: description
    content: "The generate command allows you to generate bundle files and send data to webhooks for automated processing."
---
# Generate Command

The `generate` command in GoNoGo allows you to generate helm release version information for specific applications or bundles, and optionally send this data to webhooks for automated processing.

## Usage

```bash
gonogo generate [helm-release-name] [flags]
```

## Modes

### Individual Mode

Generate information for a single helm release:

```bash
gonogo generate my-app --repo https://charts.example.com --desired-version 1.2.3
```

**Required flags for individual mode:**
- `--desired-version` or `-V`: The target version you want to upgrade to
- `--repo` or `-r`: The helm repository URL

**Optional flags:**
- `--bundle-output` or `-B`: Generate a bundle file with the specified path
- `--webhook`: Send data to an n8n webhook URL
- `--webhook-api-key`: API key for webhook authentication (or use `GNG_API_KEY` env var)
- `--analyze`: Enable OpenAI-powered upgrade analysis
- `--output` or `-o`: Output format (text, json)

### Bundle Mode

Process multiple addons from a bundle file:

```bash
gonogo generate --bundle my-bundle.yaml
```

**Required flags for bundle mode:**
- `--bundle` or `-b`: Path to the bundle file

**Optional flags:**
- `--webhook`: Send data to an n8n webhook URL
- `--webhook-api-key`: API key for webhook authentication (or use `GNG_API_KEY` env var)
- `--analyze`: Enable OpenAI-powered upgrade analysis
- `--output` or `-o`: Output format (text, json)

## Bundle Generation

When using individual mode, you can generate a properly formatted bundle file using the `--bundle-output` flag:

```bash
gonogo generate gng-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --desired-version 4.13.0 \
  --bundle-output nginx-bundle.yaml
```

This will create a bundle file with the following structure:

```yaml
addons:
- name: gng-nginx
  versions:
    current: ""
    desired: 4.13.0
  notes: ""
  source:
    chart: ingress-nginx
    repository: https://kubernetes.github.io/ingress-nginx
  warnings: []
  compatible_k8s_versions:
    min: ""
    max: ""
  necessary_api_versions: []
  values_schema: ""
  opa_checks: []
  resources: []
```

The generated bundle file includes all necessary fields with proper types, ready for you to customize with your specific validation rules.

## Webhook Integration

You can send the generated data to an n8n webhook for automated processing:

```bash
gonogo generate my-app \
  --repo https://charts.example.com \
  --desired-version 1.2.3 \
  --webhook https://your-n8n-workflow.com/webhook \
  --webhook-api-key your-api-key
```

The webhook will receive a JSON payload containing:
- Cluster version information
- Current and desired versions
- Release metadata
- Optional OpenAI analysis results

## OpenAI Analysis

Enable AI-powered upgrade analysis with the `--analyze` flag:

```bash
gonogo generate my-app \
  --repo https://charts.example.com \
  --desired-version 1.2.3 \
  --analyze \
  --openai-api-key your-openai-key
```

This provides insights into:
- Breaking changes between versions
- CRD changes
- Upgrade considerations
- Potential compatibility issues

## Dry Run Mode

Test webhook functionality without connecting to Kubernetes:

```bash
gonogo generate my-app \
  --repo https://charts.example.com \
  --desired-version 1.2.3 \
  --webhook https://your-n8n-workflow.com/webhook \
  --dry-run
```

This creates mock data to test your webhook integration.

## Examples

### Generate a bundle file for nginx-ingress:

```bash
gonogo generate gng-nginx \
  --repo https://kubernetes.github.io/ingress-nginx \
  --desired-version 4.13.0 \
  --bundle-output nginx-ingress-bundle.yaml
```

### Send data to webhook with bundle generation:

```bash
gonogo generate cert-manager \
  --repo https://charts.jetstack.io \
  --desired-version 1.13.0 \
  --webhook https://your-workflow.com/webhook \
  --bundle-output cert-manager-bundle.yaml
```

### Process existing bundle with webhook:

```bash
gonogo generate --bundle my-bundle.yaml --webhook https://your-workflow.com/webhook
```

## Output Formats

The command supports multiple output formats:

- **Text** (default): Human-readable output
- **JSON**: Structured data for programmatic use

```bash
gonogo generate my-app --repo https://charts.example.com --desired-version 1.2.3 --output json
```
