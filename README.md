# 🔐 secret-hub

A lightweight CLI tool for encrypting, decrypting, and securely storing secrets using AES encryption --- built in Go.

## 🚀 Features

- 🔒 AES-based encryption & decryption

- 📁 Local secret storage

- 🧪 Simple CLI interface

- 🛠️ Built with Go modules

## 📦 Installation

Clone the repo and build:

bash

```
git clone https://github.com/rogerio-castellano/secret-hub.git
cd secret-hub
go build -o secret-hub

```

## 🧰 Usage

### 🔐 Encrypt a secret

bash

```
./secret-hub encrypt --key my-secret-key --value "super-sensitive-data"

```

### 🔓 Decrypt a secret

bash

```
./secret-hub decrypt --key my-secret-key --value "<encrypted-string>"

```

### 📁 Store a secret

bash

```
./secret-hub store --key my-secret-key --name db_password --value "p@ssw0rd"

```

### 📤 Retrieve a stored secret

bash

```
./secret-hub get --key my-secret-key --name db_password

```

## 🛠 Tech Stack

| Component | Description                                        |
| --------- | -------------------------------------------------- |
| Language  | Go                                                 |
| Crypto    | AES                                                |
| CLI       | Cobra (if used)                                    |
| Storage   | Local file or memory (depending on implementation) |

## 📜 License

MIT --- free to use, modify, and share.
