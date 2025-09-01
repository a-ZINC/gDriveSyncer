# 📁 gDriveSyncer

gDriveSyncer is a command-line tool for synchronizing files and folders between your local machine and Google Drive. It supports efficient uploads, versioning, change detection, and visual exploration of your Drive contents — making it ideal for backups, automation, and collaboration workflows.

## ✨ Features

- 🔼 **Push Changes**: Uploads only modified files to Google Drive
- 📋 **List & Explore**: Visually browse your Google Drive files and folders in a tree format
- 🔍 **Search Files**: Find files and folders by name or ID across Drive and shared folders
- 📥 **Clone Remote Folders**: Download entire folders from Google Drive to your local machine
- 🕒 **Versioning**: Maintains historical versions of your uploads
- 🚀 **Multi-worker Uploads**: Upload files in parallel for faster sync
- ⚙️ **Initialization**: Easily create a new sync configuration
- 🐞 **Verbose Logging**: Enable detailed logs for debugging and monitoring
- 🔐 **OAuth2 Authentication**: Secure integration with Google Drive

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

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a project and enable the Google Drive API
3. Under "Credentials", create an OAuth 2.0 Client ID (type: Desktop App)
4. Download the `credentials.json` file
5. Place it in your home directory under `~/.gdrive/credentials.json`:

```bash
mkdir -p ~/.gdrive
mv credentials.json ~/.gdrive/
```

## 🚀 Usage

### 1. Initialize Sync Configuration

Set up your local workspace for Google Drive sync and authentication.

```bash
gdrivesyncer --create
```

Prompts for Google authentication and creates a config file.

### 2. Push Local Changes to Google Drive

Upload only changed files and folders to your Drive. Supports parallel uploads and verbose logging.

```bash
gdrivesyncer --push
```

**Options:**
- `--workers` / `-w`: Number of parallel upload workers
- `--verbose` / `-v`: Enable detailed logging

### 3. List & Explore Google Drive Contents

Visualize your Drive and shared folders in a tree format. Filter by file type and location (Drive/Shared).

```bash
gdrivesyncer list [options]
```

**Options:**
- `-t` / `--type`: 0=All, 1=Drive only, 2=Shared only
- `-x` / `--file-type`: 0=Both, 1=Folders only, 2=Files only

### 4. Search Files and Folders

Search for files/folders by name or ID in your Drive and shared folders. Supports searching by partial name, exact ID, or root folder ID.

```bash
gdrivesyncer search "report.pdf"
gdrivesyncer search "<file_or_folder_id>"
```

Searches both your Drive and shared items. Shows results with icons, color, size, and parent info.

### 5. Clone Remote Folder

Download a folder from Google Drive to your local machine.

```bash
gdrivesyncer clone <folderId> --destination <local_path>
```

**Options:**
- `--folderId` / `-f`: ID of the Google Drive folder to clone
- `--destination` / `-d`: Local path to save the cloned files

## 🖥️ List Output Example

```
📁 MyDrive (root)
├── 📁 Projects (folderId)
│   ├── 📄 report.pdf (fileId) [1.2 MB]
│   └── 📁 Docs (folderId)
│       └── 📄 notes.txt (fileId) [4 KB]
└── 📄 todo.txt (fileId) [2 KB]

📁 Shared with me
└── 📁 TeamFolder (folderId)
    └── 📄 shared_doc.docx (fileId) [3.5 MB]
```

## ⚙️ Command-line Options

### General Options

| Flag | Description |
|------|-------------|
| `--create` | Create a new sync configuration |
| `--push` | Upload local changes to Google Drive |
| `--workers` | Set number of parallel upload workers |
| `--verbose` | Enable detailed logging during execution |

### List Command Options

| Flag | Long Flag | Description |
|------|-----------|-------------|
| `-t` | `--type` | 0: All, 1: Drive, 2: Shared |
| `-x` | `--file-type` | 0: Both, 1: Folder, 2: File |
| `-v` | `--verbose` | Enable detailed logging |

## ✅ Examples

```bash
# Push with multiple workers and verbose logging
gdrivesyncer --push --workers=4 --verbose

# List all files and folders
gdrivesyncer list -t 0

# List only folders in your Drive with verbose output
gdrivesyncer list -t 1 -x 1 -v

# List only files shared with you
gdrivesyncer list -t 2 -x 2

# Search for a specific file
gdrivesyncer search "my-document.pdf"

# Clone a folder to local directory
gdrivesyncer clone 1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms --destination ./downloaded-folder
```

## 🔧 Requirements

- Go 1.18 or higher
- Google Cloud Platform account
- Google Drive API enabled
- Valid OAuth2 credentials

## 📝 Configuration

gDriveSyncer uses a `gdrive_init.json` configuration file that is automatically generated when you run the `--create` command.
