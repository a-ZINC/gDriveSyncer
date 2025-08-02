# 📁 gDriveSyncer
**gDriveSyncer** is a command-line tool for synchronizing files and folders between your local machine and Google Drive. It supports efficient uploads, versioning, change detection, and visual exploration of your Drive contents — making it ideal for backups, automation, and collaboration workflows.

---

## ✨ Features
- 🔼 **Push Changes**: Uploads only modified files to Google Drive.
- 📋 **List & Explore**: Visually browse your Google Drive files and folders in a tree format.
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

### 📋 List & Explore Google Drive Contents
The **List** feature lets you visually explore your Google Drive and shared folders from the command line. It displays files and folders in a tree format, showing names, types, sizes, and IDs, with color and icons for clarity.

#### List All Files (Drive + Shared)
```bash
gdrivesyncer list -t 0
```

#### List Only Your Drive
```bash
gdrivesyncer list -t 1
```

#### List Only Shared With Me
```bash
gdrivesyncer list -t 2
```

#### Filter by File Type
- **Folders only**: `-x 1`
- **Files only**: `-x 2`
- **Both**: `-x 0`

Example: List only folders in your Drive
```bash
gdrivesyncer list -t 1 -x 1
```

#### 🖥️ List Output Example
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

### ⚙️ Command-line Options

#### General Options
| Flag | Description |
|------|-------------|
| `--create` | Create a new sync configuration |
| `--push` | Upload local changes to Google Drive |
| `--workers` | Set number of parallel upload workers |
| `--verbose` | Enable detailed logging during execution |

#### List Command Options
| Flag | Long Flag | Description |
|------|-----------|-------------|
| `-t` | `--type` | 0: All, 1: Drive, 2: Shared |
| `-x` | `--file-type` | 0: Both, 1: Folder, 2: File |
| `-v` | `--verbose` | Enable detailed logging |

### ✅ Examples
```bash
# Push with multiple workers and verbose logging
gdrivesyncer --push --workers=4 --verbose

# List all files and folders
gdrivesyncer list -t 0

# List only folders in your Drive with verbose output
gdrivesyncer list -t 1 -x 1 -v

# List only files shared with you
gdrivesyncer list -t 2 -x 2
```

## 🔧 Requirements
- Go 1.18 or higher
- Google Cloud Platform account
- Google Drive API enabled
- Valid OAuth2 credentials

## 📝 Configuration
gDriveSyncer uses a `gdrive_init.json` configuration file that is automatically generated when you run the `--create` command.
