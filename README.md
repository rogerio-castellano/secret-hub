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
go build -o secret-hub
```

## 🧰 Usage

### 🔑 Generate a key

```bash
./secret-hub generate-key --out mykey.bin
```

### 🔐 Encrypt a secret

```bash
./secret-hub encrypt --in secret.txt --out secret.enc --key mykey.bin
```

### 🔓 Decrypt a secret

```bash
./secret-hub decrypt --in secret.enc --out secret-dec.txt --key mykey.bin
```

### 📥 Store a secret

```bash
./secret-hub store --key mykey.bin --name db_password --value "p@ssw0rd"
```

### 📤 Retrieve a stored secret

```bash
./secret-hub get --key my-secret-key --name db_password
```

### 🗑️ Delete a stored secret

```bash
./secret-hub delete --name db_password
```

### 📋 List stored secrets

```bash
./secret-hub list
```

## 🛠 Tech Stack

| Component | Description  |
| --------- | ------------ |
| Language  | Go           |
| Crypto    | AES          |
| CLI       | Cobra, Viper |
| Storage   | Local file   |

## 📜 License

MIT --- free to use, modify, and share. See [LICENSE](LICENSE) for full details.
