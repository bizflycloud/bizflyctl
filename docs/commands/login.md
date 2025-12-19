# Login Command

The `login` command authenticates you with Bizfly Cloud using a browser-based authentication flow.

## Usage

```bash
bizfly login [--project-id PROJECT_ID]
```

## Description

The `login` command provides an interactive way to authenticate with Bizfly Cloud. It:

1. Starts a local HTTP server on port 15995
2. Opens your default web browser to the Bizfly Cloud login page
3. Waits for you to complete authentication in the browser
4. Automatically saves your authentication token to the configuration file

### Project-Scoped Tokens

If you provide the `--project-id` flag, the login command will:

- Exchange the root token for a project-scoped token
- Save the project-scoped token to your configuration file

If you don't provide `--project-id`, the root token is saved as-is.

## How It Works

1. When you run `bizfly login`, a local server starts on `localhost:15995`
2. Your browser opens to: `https://id.bizflycloud.vn/login?service=http://localhost:15995/callback`
3. You log in with your Bizfly Cloud credentials
4. After successful login, the browser redirects back to the local server
5. The authentication token is extracted and saved to `~/.bizfly.yaml`

## Examples

### Basic Login

```bash
bizfly login
```

Output:

```
Opening browser to login: https://id.bizflycloud.vn/login?service=http://localhost:15995/callback
Login successful! Token saved to config file.
```

### Login with Project ID

To exchange the root token for a project-scoped token:

```bash
bizfly login --project-id 092a68edd01b4afca8f4d4783395727c
```

Output:

```
Opening browser to login: https://id.bizflycloud.vn/login?service=http://localhost:15995/callback
Login successful! Project-scoped token saved to config file.
```

This will exchange the root token for a project-scoped token that is limited to the specified project.

### Login with Manual Browser Opening

If the browser doesn't open automatically, you'll see:

```
Opening browser to login: https://id.bizflycloud.vn/login?service=http://localhost:15995/callback
Failed to open browser: <error>
Please open the URL manually.
```

Simply copy the URL and open it in your browser.

## Authentication Token

After successful login, the authentication token is saved to your configuration file:

-   **Location:** `~/.bizfly.yaml` (macOS/Linux) or `%USERPROFILE%\.bizfly.yaml` (Windows)
-   **Field:** `auth_token`

The token is used for subsequent API calls, so you don't need to log in again until it expires.

### Token Types

-   **Root Token**: Obtained when logging in without `--project-id`. This token has access to all projects.
-   **Project-Scoped Token**: Obtained when logging in with `--project-id`. This token is limited to the specified project and is automatically exchanged from the root token.

## Troubleshooting

### Port 15995 Already in Use

If port 15995 is already in use, you may see:

```
failed to start local server: listen tcp localhost:15995: bind: address already in use
```

**Solution:**

-   Close any other applications using port 15995
-   Or wait a few moments and try again

### Browser Doesn't Open Automatically

If your browser doesn't open:

1. Copy the login URL from the output
2. Open it manually in your browser
3. Complete the login process

### Login Failed

If login fails, check:

1. **Internet connection** - Ensure you can reach `id.bizflycloud.vn`
2. **Credentials** - Verify your Bizfly Cloud account credentials
3. **Firewall** - Ensure localhost:15995 is not blocked
4. **Browser** - Try a different browser if issues persist

### Token Not Saved

If the token isn't saved:

1. Check file permissions on `~/.bizfly.yaml`
2. Ensure you have write access to your home directory
3. Try running with appropriate permissions

## Security Notes

-   The local server only listens on `localhost` (127.0.0.1), so it's not accessible from other machines
-   The authentication token is stored in plain text in the configuration file
-   Set appropriate file permissions: `chmod 600 ~/.bizfly.yaml`
-   Never share your configuration file or authentication token

## Alternative Authentication Methods

If browser login doesn't work for you, consider:

1. **Configuration file** - See [Authentication Guide](../authentication.md#method-2-configuration-file)
2. **Environment variables** - See [Authentication Guide](../authentication.md#method-3-environment-variables)
3. **Application credentials** - See [Authentication Guide](../authentication.md#method-4-application-credentials)

## Related Commands

-   [Authentication Guide](../authentication.md) - Other authentication methods
-   [Configuration Guide](../configuration.md) - Managing configuration files
