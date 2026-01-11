# gauth

Fast, no-nonsense 2FA for your terminal, written in Go.

## tl;dr
- `gauth`: Show codes (with color-coded countdowns)
- `gauth -w`: Watch mode (auto-refresh)
- `gauth -a`: Add new account
- `gauth -d`: Delete account
- `gauth -l`: List all accounts
- `gauth -p`: Manage master password (AES-256 encryption)
- `gauth -i/-e`: andOTP backup support
- Saves to `$HOME/.gauth/gauth.json` (atomic writes)

## Installation
```bash
go install github.com/leeineian/gauth@latest
```

Or build from source:
```bash
git clone https://github.com/leeineian/gauth.git
cd gauth && go build -o gauth .
```

### Docker
```bash
# Build the image
docker build -t gauth .

# Run gauth (mounting your secrets)
docker run -it -v ~/.gauth:/root/.gauth gauth
```

## Usage

**Viewing codes**
```bash
./gauth
# built-in live mode (updates every second):
./gauth -w
```

**Managing Accounts**
```bash
# list all accounts
./gauth -l
```

**Adding/Deleting accounts**
```bash
# Add
./gauth -a

# Delete
./gauth -d
```

**Importing/Exporting**
```bash
./gauth -i
./gauth -e
```

**Security**
```bash
# set or change master password
# (leave empty to remove password)
./gauth -p
```
Once a password is set, your `gauth.json` is encrypted using AES-256-GCM.

## Requirements
- Go 1.25.5 or higher

## License
AGPL-3.0
