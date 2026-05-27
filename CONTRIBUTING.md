# Contributing

Keep this repository public-safe by default.

- Do not add real credentials, tokens, cookies, browser profiles, private local
  paths, private tenant state, or machine-specific automation.
- Use example values in docs and tests.
- Keep session state under `.open-ralph/`, which is gitignored.
- Add tests for new command behavior.
- Run `go test ./...` and `gitleaks detect --source . --no-git --redact`
  before opening a pull request.

