# 🔐 secret-hub

A lightweight CLI tool for encrypting, decrypting, and securely storing secrets using AES encryption --- built in Go.

## 🚀 Features

- 🔒 AES-based encryption & decryption

- 📁 Local secret storage

- 🧪 Simple CLI interface

- 🛠️ Built with Go modules

## 📦 Installation

Clone the repo and build:

```bash
git clone https://github.com/rogerio-castellano/secret-hub.git
cd secret-hub
go build
```

## 🔧 Commands Overview

| Command      | Description                                                      |
| ------------ | ---------------------------------------------------------------- |
| generate-key | Create a new encryption key                                      |
| encrypt      | Encrypt a file                                                   |
| decrypt      | Decrypt a file                                                   |
| store        | Encrypt and store a secret by name                               |
| get          | Retrieve and decrypt a secret by name                            |
| delete       | Delete a stored secret by name                                   |
| list         | List all stored secret names                                     |
| search       | Search for stored secret names by substring                      |
| import       | Import secrets from .env, JSON, or YAML and store them encrypted |
| export       | Export decrypted secrets in .env, JSON, or YAML format           |

## 🧰 Usage

### 🔑 Generate a key

```bash
secret-hub generate-key [--out <keyfile>]
```

#### Examples

Generate a new encryption key and save it as `custom-key.bin`

```bash
secret-hub generate-key --out custom-key.bin
```

Generate a new encryption key and save it using the default file name (`key.bin`)

```bash
secret-hub generate-key
```

### 🔐 Encrypt a secret

```bash
secret-hub encrypt --input <plaintext> --output <ciphertext> [--key <keyfile>] [--base64]
```

#### Examples

Encrypt `secret.txt` using the key provided in `key.bin`; output saved as `secret.enc`

```bash
secret-hub encrypt --in secret.txt --out secret.enc --key key.bin
```

Encrypt `secret.txt` using default key and settings from the configuration file

```bash
secret-hub encrypt --in secret.txt --out secret.enc
```

Encrypt `secret.txt` using `key.bin` and encode the output in Base64 format

```bash
secret-hub encrypt --in secret.txt --out secret.enc --key key.bin --base64
```

### 🔓 Decrypt a secret

```bash
secret-hub decrypt --input <ciphertext> --output <plaintext> [--key <keyfile>] [--base64]
```

#### Examples

Decrypt `secret.enc` using the key from `key.bin`; output saved as `secret-dec.txt`

```bash
secret-hub decrypt --in secret.enc --out secret-dec.txt --key key.bin
```

Decrypt `secret.enc` using default key and settings from the configuration file

```bash
secret-hub decrypt --in secret.enc --out secret-dec.txt
```

Decrypt Base64-encoded `secret.enc` using `key.bin`; decoded output saved as `secret-dec.txt`

```bash
secret-hub decrypt --in secret.enc --out secret-dec.txt --key key.bin --base64
```

### 📁 Encrypt and store a secret by name

```bash
secret-hub store --name <secret_name> --value <secret_value> [--key <keyfile>] [--storage <filepath>] [--force]
```

#### Examples

Store the secret 'db_password' with the value "p@ssw0rd" in `secret-store.json`, using `key.bin` to encrypt

```bash
secret-hub store --key key.bin --storage=secret-store.json --name db_password --value "p@ssw0rd"
```

Store the secret 'db_password' with the value "p@ssw0rd" using default encryption key and storage settings

```bash
secret-hub store --name db_password --value "p@ssw0rd"
```

Overwrite existing 'db_password' secret with the new value "p@ssw0rd" in `secret-store.json`, using `key.bin` and --force to overwrite if a secret with the same name already exists

```bash
secret-hub store --key key.bin --storage=secret-store.json --name db_password --value "p@ssw0rd" --force
```

### 📤 Retrieve and decrypt a secret by name

```bash
secret-hub get --name <secret_name> [--key <keyfile>] [--storage <filepath>]
```

#### Examples

Retrieve the secret named 'db_password' from local storage file `secret-store.json` using the decryption key from `key.bin`

```bash
secret-hub get --name db_password --key key.bin --storage secret-store.json
```

Retrieve the secret named 'db_password' using default configuration and key settings

```bash
secret-hub get --name db_password
```

### 🗑️ Delete a stored secret by name

```bash
delete --name <secret_name> [--storage <filepath>]
```

#### Examples

Delete the secret named 'mysecret' from the `secret-store.json` storage file

```bash
secret-hub delete --name mysecret --storage secret-store.json
```

Delete the secret named 'db_password' using the default storage location

```bash
secret-hub delete --name db_password
```

### 📋 List all stored secret names

```bash
secret-hub list [--storage <filepath>] [--json] [--pretty]
```

#### Examples

List all stored secrets from the file `secret-store.json`

```bash
secret-hub list --storage secret-store.json
```

List all stored secrets and format the output as JSON for programmatic use or integration

```bash
secret-hub list --json
```

List all stored secrets using default storage location

```bash
secret-hub list
```

### 🔍 Search for stored secret names by substring

```bash
search --query <substring> [--storage <filepath>]
```

#### Examples

Search for secrets matching the keyword 'db' within the `secret-store.json` storage file

```bash
secret-hub search --query db --storage secret-store.json
```

Search for secrets matching the keyword 'db' using default configuration and storage

```bash
  secret-hub search --query db
```

### 📥 Import secrets and store them encrypted

```bash
secret-hub import --file <path> --format <env|json|yaml> --key <keyfile> [--storage <filepath>] [--force] [--skip-existing] [--dry-run] [--only <keys>] [--rename <keys>] [--exclude <keys>] [--summary=json] [--quiet]
```

#### Examples

Import from a .env file into custom storage

```bash
secret-hub import --file secrets.env --format env --key test-key.bin --storage imported-secrets.json
```

Import from JSON and overwrite existing secrets

```bash
secret-hub import --file secrets.json --format json --key test-key.bin --force
```

Import from YAML and skip secrets that already exist

```bash
secret-hub import --file secrets.yaml --format yaml --skip-existing
```

Dry-run preview without encrypting or saving

```bash
secret-hub import --file secrets.json --format json --key test-key.bin --dry-run
```

Import only specific secrets (e.g., API_KEY and DB_PASSWORD)

```bash
secret-hub import config.env --key superkey123 --only API_KEY,DB_PASSWORD
```

Import prepending dev\_ to each key (e.g. DB_PASS becomes dev_DB_PASS)

```bash
secret-hub import --file ./config.env --key masterKey123 --prefix dev\_
```

Import selectively with prefixing; only DB_PASS and API_KEY will be imported, with keys remapped to prod_DB_PASS and prod_API_KEY.

```bash
secret-hub import --file ./config.env --key masterKey123 --only DB*PASS,API_KEY --prefix prod*
```

Dry run with preview of prefixed output. Preview what will be saved as test\_-prefixed keys, allowing validation before committing.

```bash
secret-hub import --file ./secrets.yaml --key masterKey123 --prefix test\_ --dry-run
```

Import the secrets from secrets.env, renaming DB_PASS to DATABASE_PASSWORD and API_KEY to SERVICE_API_KEY.

```bash
secret-hub import secrets.env --rename DB_PASS=DATABASE_PASSWORD --rename API_KEY=SERVICE_API_KEY
```

Import only DB_USER and DB_PASS, and their keys are renamed to USERNAME and PASSWORD

```bash
secret-hub import config.yaml --only DB_USER DB_PASS --rename DB_USER=USERNAME --rename DB_PASS=PASSWORD
```

Import the TOKEN secret as prod_ACCESS_TOKEN in storage, utilizing the combined effect of --prefix and --rename

```bash
secret-hub import app.env --prefix prod\_ --rename TOKEN=ACCESS_TOKEN
```

Dry run prefixing each key with dev\_ (e.g., dev_API_KEY), renaming DATABASE_URL to db_url, and excluding DEBUG_MODE and LOCAL_SECRET—even if they're in the file.

```bash
secret-hub import --file=secrets.env --format env --prefix dev\_ --rename DATABASE_URL=db_url --exclude DEBUG_MODE,LOCAL_SECRET --dry-run --file=secrets.env
```

Import JSON secrets with a summary in machine-readable format

```bash
secret-hub import --file secrets.json --key team-key --summary=json
```

### 📥 Export decrypted secrets

````bash
export --format <env|json|yaml> [--key <keyfile>] [--storage <filepath>] [--summary <json|text>] [--quiet] [--output <filename>]
```

#### Examples

Export all stored secrets in decrypted form using 'env' format and the key from ```key.bin```
```bash
secret-hub export --format env --key key.bin

````

Export all stored secrets in decrypted form using 'json' format and `key.bin`; redirect output to `exported-secrets.json`

````bash
secret-hub export --format json --key ```key.bin``` > ```exported-secrets.json```
````

Export all stored secrets from `secret-store.json` in 'yaml' format using the key from `key.bin`

````bash
secret-hub export --format yaml --key ```key.bin``` --storage ```secret-store.json```
````

Export all stored secrets in 'env' format, save to `.env.generated`, and include a brief summary of exported items

````bash
secret-hub export --format env --summary --output ```.env.generated```
````

## 🛠 Tech Stack

| Component | Description  |
| --------- | ------------ |
| Language  | Go           |
| Crypto    | AES          |
| CLI       | Cobra, Viper |
| Storage   | Local file   |

## 📜 License

MIT --- free to use, modify, and share. See [LICENSE](LICENSE) for full details.
