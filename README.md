# secret-hub

CLI to encrypt/decrypt &amp; store secrets using AES

## Features

- Encrypt and decrypt secrets using AES-256-GCM
- Store secrets securely in local files
- Simple CLI interface for managing secrets
- Cross-platform support

## Installation

```sh
go install github.com/yourusername/secret-hub@latest
```

## Usage

### Encrypt a secret

```sh
secret-hub encrypt --key <your-key> --in secret.txt --out secret.enc
```

### Decrypt a secret

```sh
secret-hub decrypt --key <your-key> --in secret.enc --out secret.txt
```

### Store a secret

```sh
secret-hub store --key <your-key> --name <secret-name> --value <secret-value>
```

### Retrieve a secret

```sh
secret-hub get --key <your-key> --name <secret-name>
```

## Security

- Uses AES-256-GCM for encryption
- Secrets are never stored in plaintext
- Key management is the user's responsibility

## Contributing

Contributions are welcome! Please open issues or pull requests.

## License

MIT License
