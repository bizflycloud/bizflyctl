# Configuration Guide

This guide explains how to configure `bizflyctl` with default settings to simplify your workflow.

## Configuration File Location

The configuration file is located at:

-   **macOS/Linux:** `~/.bizfly.yaml`
-   **Windows:** `%USERPROFILE%\.bizfly.yaml`

## Configuration File Format

The configuration file uses YAML format. Here's an example:

```yaml
# Authentication
email: your-email@example.com
password: your-password
# OR use application credentials
app_credential_id: your-app-credential-id
app_credential_secret: your-app-credential-secret

# Region (default: HaNoi)
region: HaNoi

# Project ID (optional)
project_id: your-project-id

# Authentication token (automatically set by 'bizfly login')
auth_token: your-auth-token
```

## Configuration Options

### Authentication Options

-   **email**: Your Bizfly Cloud email address
-   **password**: Your Bizfly Cloud password
-   **app_credential_id**: Application credential ID (alternative to email/password)
-   **app_credential_secret**: Application credential secret
-   **auth_token**: Authentication token (set automatically by `bizfly login`)

### General Options

-   **region**: Default region to use for operations
    -   Options: `HaNoi`, `HoChiMinh`, etc.
    -   Default: `HaNoi`
-   **project_id**: Default project ID for operations

## Creating the Configuration File

### Method 1: Using `bizfly login`

The easiest way to create the configuration file is by using the login command:

```bash
bizfly login
```

This will create the configuration file and save your authentication token.

### Method 2: Manual Creation

Create the file manually:

```bash
# macOS/Linux
nano ~/.bizfly.yaml

# Windows (PowerShell)
notepad $env:USERPROFILE\.bizfly.yaml
```

Then add your configuration:

```yaml
email: your-email@example.com
password: your-password
region: HaNoi
project_id: your-project-id
```

## Setting File Permissions

On macOS and Linux, set restrictive permissions on your configuration file:

```bash
chmod 600 ~/.bizfly.yaml
```

This ensures only you can read and write the file.

## Using Custom Configuration File

You can specify a custom configuration file location using the `--config` flag:

```bash
bizfly --config /path/to/custom/config.yaml server list
```

## Overriding Configuration

You can override configuration file settings using:

1. **Command-line flags:**

    ```bash
    bizfly --region HoChiMinh server list
    ```

2. **Environment variables:**
    ```bash
    export BIZFLY_CLOUD_REGION=HoChiMinh
    bizfly server list
    ```

## Priority Order

Settings are resolved in this order (highest to lowest priority):

1. Command-line flags
2. Environment variables
3. Configuration file
4. Default values

## Example Configuration

Here's a complete example configuration file:

```yaml
# Primary authentication method
email: user@example.com
password: secure-password-here

# Alternative: Application credentials
# app_credential_id: app-cred-id
# app_credential_secret: app-cred-secret

# Default region
region: HaNoi

# Default project
project_id: 12345678-1234-1234-1234-123456789012

# Auth token (automatically managed)
auth_token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

## Verifying Configuration

Test your configuration:

```bash
# This will show which config file is being used
bizfly server list
```

The output will include: `Using config file: /path/to/.bizfly.yaml`

## Troubleshooting

### Configuration file not found

If you see an error about the configuration file not being found:

1. Create the file manually (see "Creating the Configuration File" above)
2. Or use `bizfly login` to create it automatically

### Configuration not being used

1. Check the file location and name (`~/.bizfly.yaml`)
2. Verify YAML syntax is correct
3. Check file permissions (should be readable)
4. Use `--config` flag to specify the exact path

### Invalid YAML syntax

Common YAML syntax errors:

-   Missing colons after keys
-   Incorrect indentation
-   Special characters not quoted

Use a YAML validator to check your syntax.

## Security Considerations

1. **Never commit the configuration file** to version control
2. **Set restrictive permissions** (600 on Unix systems)
3. **Use application credentials** for automated scripts
4. **Rotate credentials regularly**
5. **Use environment variables** in shared environments

## Next Steps

-   [Command Reference](README.md#command-reference) - Learn available commands
-   [Server Management](commands/server.md) - Manage cloud servers
-   [Volume Management](commands/volume.md) - Manage storage volumes
