# 📁 gDriveSyncer

**gDriveSyncer** is a command-line tool for synchronizing files and folders between your local machine and Google Drive. It supports efficient uploads, versioning, and change detection — making it ideal for backups, automation, and collaboration workflows.

---

## ✨ Features

- 🔼 **Push Changes**: Uploads only modified files to Google Drive.
- 🕒 **Versioning**: Maintains historical versions of your uploads.
- 🚀 **Multi-worker Uploads**: Upload files in parallel for faster sync.
- ⚙️ **Initialization**: Easily create a new sync configuration.
- 🐞 **Verbose Logging**: Enable detailed logs for debugging and monitoring.
- 🔐 **OAuth2 Authentication**: Secure integration with Google Drive.

---

## 📦 Installation

### Using Homebrew (macOS/Linux)

```bash
brew tap a-zinc/gdrivesyncer
brew install gdrivesyncer
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/a-zinc/gDriveSyncer.git
cd gDriveSyncer

# Build the binary
go build -o gDriveSyncer .
```

**Note**: Requires Go 1.18 or higher installed.

## 🔑 Authentication Setup

1. Go to [Google Cloud Console](https://console.cloud.google.com/).
2. Create a project and enable the Google Drive API.
3. Under "Credentials", create an OAuth 2.0 Client ID (type: Desktop App).
4. Download the `credentials.json` file.
5. Place it in your home directory under `~/.gdrive/credentials.json`:

```bash
mkdir -p ~/.gdrive
mv credentials.json ~/.gdrive/
```

## 🚀 Usage

### 🆕 Initialize Sync Configuration

```bash
gdrivesyncer --create
```

Generates a new `gdrive_init.json` file in the current directory.

### 📤 Push Local Changes to Google Drive

```bash
gdrivesyncer --push
```

Syncs updated files and folders based on your configuration.

### ⚙️ Command-line Options

| Flag | Description |
|------|-------------|
| `--create` | Create a new sync configuration |
| `--push` | Upload local changes to Google Drive |
| `--workers` | Set number of parallel upload workers |
| `--verbose` | Enable detailed logging during execution |

### ✅ Example

```bash
gdrivesyncer --push --workers=4 --verbose
```

## 🔧 Requirements

- Go 1.18 or higher
- Google Cloud Platform account
- Google Drive API enabled
- Valid OAuth2 credentials

## 📝 Configuration

gDriveSyncer uses a `gdrive_init.json` configuration file that is automatically generated when you run the `--create` command.
