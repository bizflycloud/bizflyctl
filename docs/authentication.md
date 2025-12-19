# Authentication Guide

`bizflyctl` supports multiple authentication methods. This guide explains how to authenticate and configure your credentials.

## Authentication Methods

### Method 1: Browser Login (Recommended)

The easiest way to authenticate is using the browser login command:

```bash
bizfly login
```

Or with a project ID to get a project-scoped token:

```bash
bizfly login --project-id YOUR_PROJECT_ID
```

This command will:

1. Open your default web browser
2. Redirect you to the Bizfly Cloud login page
3. After successful login, automatically save your authentication token
4. If `--project-id` is provided, exchange the root token for a project-scoped token

**Note:** The token is saved to `~/.bizfly.yaml` (or `%USERPROFILE%\.bizfly.yaml` on Windows).

**Project-Scoped Tokens:** When you provide `--project-id`, the login command exchanges the root token for a project-scoped token that is limited to the specified project. This is useful when you want to restrict access to a specific project.

### Method 2: Configuration File

Create a configuration file at `~/.bizfly.yaml` (or `%USERPROFILE%\.bizfly.yaml` on Windows):

```yaml
email: your-email@example.com
password: your-password
region: HaNoi
project_id: your-project-id
```

**Security Note:** Make sure to set appropriate file permissions:

```bash
chmod 600 ~/.bizfly.yaml
```

### Method 3: Environment Variables

Set the following environment variables:

```bash
export BIZFLY_CLOUD_EMAIL="your-email@example.com"
export BIZFLY_CLOUD_PASSWORD="your-password"
export BIZFLY_CLOUD_REGION="HaNoi"  # Optional, default is HaNoi
export BIZFLY_CLOUD_PROJECT_ID="your-project-id"  # Optional
```

**Windows PowerShell:**

```powershell
$env:BIZFLY_CLOUD_EMAIL="your-email@example.com"
$env:BIZFLY_CLOUD_PASSWORD="your-password"
$env:BIZFLY_CLOUD_REGION="HaNoi"
$env:BIZFLY_CLOUD_PROJECT_ID="your-project-id"
```

### Method 4: Application Credentials

You can also use application credentials for authentication:

**Using environment variables:**

```bash
export BIZFLY_CLOUD_APP_CREDENTIAL_ID="your-app-credential-id"
export BIZFLY_CLOUD_APP_CREDENTIAL_SECRET="your-app-credential-secret"
```

**Using configuration file:**

```yaml
app_credential_id: your-app-credential-id
app_credential_secret: your-app-credential-secret
region: HaNoi
project_id: your-project-id
```

**Using command-line flags:**

```bash
bizfly --app-credential-id <id> --app-credential-secret <secret> server list
```

## Command-Line Flags

You can also provide credentials directly via command-line flags:

```bash
bizfly --email your-email@example.com --password your-password --region HaNoi server list
```

**Note:** Command-line flags take precedence over configuration files and environment variables.

## Priority Order

Authentication credentials are resolved in the following order (highest to lowest priority):

1. Command-line flags (`--email`, `--password`, etc.)
2. Environment variables (`BIZFLY_CLOUD_*`)
3. Configuration file (`~/.bizfly.yaml`)
4. Stored authentication token (from `bizfly login`)

## Regions

Available regions:

-   `HaNoi` (default)
-   `HoChiMinh`
-   Other regions as supported by Bizfly Cloud

You can specify the region using:

-   Configuration file: `region: HaNoi`
-   Environment variable: `BIZFLY_CLOUD_REGION=HaNoi`
-   Command-line flag: `--region HaNoi`

## Verifying Authentication

Test your authentication by listing your servers:

```bash
bizfly server list
```

If authentication is successful, you'll see a list of your servers. If not, you'll see an authentication error.

## Troubleshooting

### "Authentication failed" error

1. Verify your credentials are correct
2. Check that your account is active
3. Ensure the region is correct
4. Try logging in again with `bizfly login`

### "Token expired" error

If your token has expired, run:

```bash
bizfly login
```

This will refresh your authentication token.

### Configuration file not found

If the configuration file doesn't exist, create it manually or use `bizfly login` to generate it automatically.

## Security Best Practices

1. **Never commit credentials to version control** - Add `~/.bizfly.yaml` to your `.gitignore`
2. **Use application credentials** for automated scripts instead of user credentials
3. **Set proper file permissions** on configuration files (600 on Unix systems)
4. **Use environment variables** in CI/CD pipelines instead of storing credentials in files
5. **Rotate credentials regularly** for security

## Next Steps

After authentication, proceed to:

-   [Configuration Guide](configuration.md) - Configure default settings
-   [Command Reference](README.md#command-reference) - Learn how to use commands
