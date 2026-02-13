# Contributing to SynCLI

These guidelines are designed for both human and AI agents (e.g., Gemini, Claude, Copilot) contributing to SynCLI. Please follow these practices to ensure code quality, safety, and maintainability.

## 1. Concurrency & Safety (The "Synapse First" Rule)

- **Mutex for Shared Data:** When updating shared data in concurrent operations, always use a mutex to ensure thread safety.
- **Bounded Concurrency:** Never launch unbounded goroutines. Limit concurrent API calls (default: 5) using semaphores (channels) or `errgroup`.
- **Fail-Fast:** Use `golang.org/x/sync/errgroup` with context cancellation. If one API call in a batch fails, immediately cancel all pending and in-flight requests for that command.
- **Backoff & Retries:** Do not retry 5xx errors blindly. Implement exponential backoff if retries are necessary.

## 2. Idiomatic Go & Architecture

- **Required Flags:** 
  - If a required flag is not set, the CLI should fail with exit code 1 and provide a clear error message.
- **Cobra/Viper Integration:**
  - Bind flags to Viper in `init()`.
- **Project Structure:**
  - `cmd/`: Only CLI orchestration (parsing flags, calling internal logic).
  - `internal/synapse/`: Matrix/Synapse client logic, decoupled from Cobra.
- **Effective Go:**
  - Avoid `init()` functions for logic; use them only for flag registration.
  - Prefer `io.Reader`/`io.Writer` over passing byte slices.
  - Prefer Generics for collection-handling (slices) functions to avoid slice-to-interface conversion overhead.

## 3. Defensive Programming & HTTP

- **Context Propagation:** Every function interacting with the network must accept a `context.Context` as its first argument.
- **Status Codes:** Always check `resp.StatusCode`. Do not assume 200 OK.
- **Explicit Error Wrapping:** Use `fmt.Errorf("context message: %w", err)` to preserve error types and provide context.
- **Zero Panics:** No `panic()` or `log.Fatal()` inside `internal/`. Return errors to the `cmd/` layer and let the CLI handle the exit.

## 4. Logging & Privacy

- **Debug Logging for Synapse:** Always add debug logs for interactions with Synapse, including request details and responses (excluding sensitive information).
- **Structured Logging:** Use logrus with fields.
- **Sanitization:** Ensure Authorization headers and access_token strings never reach Debug logs. Implement a "Redacting RoundTripper" or manual checks.
- **Verbosity:**
  - Debug: API payloads and internal state
  - Info: High-level progress
  - Error: Actionable failures

## 5. Quality Assurance

- **API Testing Requirement:** All new functions added under `internal/synapse` that interact with the Synapse API must include corresponding unit tests.
- **Testing Guidance:** Use table-driven tests and mock clients to simulate API responses and errors. See `internal/synapse/rooms_test.go` for an example.
- **Pre-PR Checklist:** Before submitting a pull request, ensure that `make audit` passes with no errors and that the binary builds successfully with `make build`.
- **Makefile:** Use the provided Makefile for build, test, lint, and formatting tasks.

## 6. Documentation & Maintenance

- **Comments:** Use Godoc style (complete sentences, starting with the symbol name).
- **README:** Update the README only when the "interface" changes (new commands, flags, or configuration keys).

---

For questions or clarifications, open an issue or start a discussion. Thank you for contributing to SynCLI!
