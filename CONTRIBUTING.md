# Contributing

Keep this repository small, portable, and reviewable.

- Do not add real credentials, tokens, cookies, host-specific paths, or user
  data.
- Use example values in docs and tests.
- Add tests for new command behavior.
- Run `go test ./...` and `gitleaks detect --source . --no-git --redact`
  before opening a pull request.
