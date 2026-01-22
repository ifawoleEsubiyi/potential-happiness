# Examples

This directory contains example programs demonstrating how to use various features of Validatord.

## GitHub Models API Example

The `models_example.go` demonstrates how to use the GitHub Models API integration to interact with LLMs.

### Prerequisites

1. A GitHub personal access token with the `models` scope
   - Create one at: https://github.com/settings/tokens
   - Required scope: `models`

2. Set the token as an environment variable:
   ```bash
   export GITHUB_TOKEN="your_token_here"
   ```

### Running the Example

```bash
# From the repository root
go run examples/models_example.go
```

### What it demonstrates

1. **Simple Chat**: Send a single user message and get a response using the default model
2. **Chat with Specific Model**: Choose which model to use for your request
3. **Advanced Usage**: Use custom parameters like temperature, max tokens, and system messages

### Example Output

```
GitHub Models API Example
=========================

Initialized with default model: openai/gpt-4o-mini

Example 1: Simple chat
----------------------
Prompt: What is the capital of France?
Response: The capital of France is Paris.

Example 2: Chat with specific model
------------------------------------
Prompt: Explain recursion in one sentence
Model: openai/gpt-4o-mini
Response: Recursion is a programming technique where a function calls itself to solve smaller instances of the same problem.

Example 3: Advanced request with custom parameters
--------------------------------------------------
Request:
  Model: openai/gpt-4o-mini
  Temperature: 0.7
  Max Tokens: 100
  Messages: 2

Response: The three primary colors are red, blue, and yellow.
Model Used: openai/gpt-4o-mini

=========================
Example completed!
```

### Available Models

You can use any model from the GitHub Models marketplace:
- `openai/gpt-4o-mini` (default, fast and efficient)
- `openai/gpt-4` (most capable OpenAI model)
- `openai/gpt-4o` (optimized for speed and cost)
- And many more from OpenAI, Meta, DeepSeek, and other providers

See the full list at: https://github.com/marketplace/models

### Error Handling

The example includes error handling for common scenarios:
- Missing GitHub token
- Invalid model names
- API request failures
- Empty responses

### Learn More

- [GitHub Models Documentation](https://docs.github.com/github-models)
- [API Reference](https://docs.github.com/en/rest/models/inference)
- [Models Marketplace](https://github.com/marketplace/models)
