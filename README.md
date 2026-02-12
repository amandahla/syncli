# SynCLI

A command-line interface (CLI) tool for interacting with Synapse, the Matrix homeserver.

For interacting with Matrix, this CLI uses the [Mautrix](https://maunium.net/go/mautrix/) Go library.

## Requirements
- Go version **go1.26.0 linux/amd64**

## Features
- Retrieve and manage Matrix spaces
- Flexible configuration and debugging options

## Installation

1. Clone the repository:
   ```sh
   git clone <repo-url>
   cd syncli
   ```
2. Build the binary or run other tasks using Makefile:
   ```sh
   make build
   make test
   make lint
   make fmt
   ```
   The lint target uses [golangci-lint](https://github.com/golangci/golangci-lint). Make sure it is installed.

## Usage

Run the CLI with:
```sh
./syncli [command] [flags]
```

### Example Commands
- Get spaces:
  ```sh
  ./syncli get spaces --debug
  ```

## Project Structure
- `main.go`: Entry point for the CLI
- `cmd/`: Command definitions (root, get, spaces, etc.)
- `internal/`: Internal logic (config, printer, requester, synapse API)

## Configuration

Required configuration:
- `base_url`: Synapse Matrix homeserver URL
- `access_token`: Access token for authentication

For more information about obtaining and using an access token, refer to the [Element Admin API documentation](https://docs.element.io/latest/element-support/advanced-administration/getting-started-using-the-admin-api/#promoting-a-matrix-account-to-admin).

Configuration options can be set via configuration file (default `$HOME/.syncli.yaml`) or environment variables. See `--help` for details on available options.

## Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.
See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## License
This project is licensed under the Apache License. See the [LICENSE](LICENSE) file for details.
