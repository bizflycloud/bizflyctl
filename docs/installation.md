# Installation Guide

This guide covers different methods to install `bizflyctl` on your system.

## Prerequisites

-   A Bizfly Cloud account
-   Command line access (Terminal on macOS/Linux, PowerShell on Windows)

## Installation Methods

### Method 1: Using Homebrew (macOS)

1. Add the Bizfly Cloud tap:

```bash
brew tap bizflycloud/bizflyctl
```

2. Install bizflyctl:

```bash
brew install bizflyctl
```

### Method 2: Build from Source

1. **Clone the repository:**

```bash
git clone https://github.com/bizflycloud/bizflyctl
cd bizflyctl
```

2. **Install Go** (if not already installed):

    - Visit [golang.org](https://golang.org/dl/) to download and install Go
    - Ensure Go is in your PATH

3. **Build the binary:**

```bash
go build -o bizfly main.go
```

4. **Install the binary** (optional):
    - **macOS/Linux:** Copy to `/usr/local/bin`:
        ```bash
        sudo cp bizfly /usr/local/bin/
        ```
    - **Windows:** Add the directory containing `bizfly.exe` to your PATH

### Method 3: Download Pre-built Binaries

1. Navigate to the [GitHub Releases page](https://github.com/bizflycloud/bizflyctl/releases)

2. Download the appropriate archive for your platform:

    - **macOS:** `bizflyctl_*_darwin_amd64.tar.gz`
    - **Linux:** `bizflyctl_*_linux_amd64.tar.gz`
    - **Windows:** `bizflyctl_*_windows_amd64.tar.gz`

3. Extract the archive:

    ```bash
    tar -xzf bizflyctl_*_<platform>_amd64.tar.gz
    ```

4. **Install the binary:**
    - **macOS/Linux:** Copy to `/usr/local/bin`:
        ```bash
        sudo cp bizfly /usr/local/bin/
        ```
    - **Windows:** Add the directory containing `bizfly.exe` to your PATH

## Verify Installation

After installation, verify that `bizflyctl` is working:

```bash
bizfly --version
```

You should see the version number. If you see a "command not found" error, ensure the binary is in your PATH.

## Next Steps

After installation, proceed to:

1. [Authentication Guide](authentication.md) - Set up your credentials
2. [Configuration Guide](configuration.md) - Configure default settings
