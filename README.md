# gauth

Fast, no-nonsense 2FA for your terminal, written in Go.

## tl;dr
- `gauth`: Show codes (with color-coded countdowns)
- `gauth entry add`: Interactive setup
- `gauth passwd`: Manage master password (AES-256 encryption)
- `gauth import/export`: andOTP backup support
- Saves to `$HOME/.gauth/gauth.json` (atomic writes)

## Installation
```bash
go install github.com/leeineian/gauth/cmd/gauth@latest
```

Or build from source:
```bash
git clone https://github.com/leeineian/gauth.git
cd gauth && go build -o gauth ./cmd/gauth
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
# build-in live mode (updates every second):
./gauth -w
```

**Adding accounts**
```bash
./gauth entry add
```

**Importing/Exporting**
```bash
./gauth import -f accounts.json -p <password>
./gauth export -f backup.json -p <password>
```

**Security**
```bash
# set or change master password
./gauth passwd
```
Once a password is set, your `gauth.json` is encrypted using AES-256-GCM.

## Requirements
- Go 1.25.5+

## License
AGPL-3.0
