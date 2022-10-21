# [golangci-lint](https://golangci-lint.run/)

## MacOS [Install](https://golangci-lint.run/usage/install/#macos)

```bash
brew install golangci-lint
```

## [Run Locally](https://golangci-lint.run/usage/quick-start/)

```bash
golangci-lint run
```

## [VS-Code Integration](https://golangci-lint.run/usage/integrations/)

```json
{
    "go.lintTool": "golangci-lint",
    // This is recommended to not freeze the editor,
    // but it isn't catching stuff!
    // "go.lintFlags": [
    //     "--fast"
    // ],
}
```

Note that with the `lintTool` set to `golangci-lint`, the `Go` VS Code extension will `go install` golangci-lint, despite the fact that this is [explicitly recommended against](https://golangci-lint.run/usage/install/#install-from-source). ¯\_(ツ)_/¯