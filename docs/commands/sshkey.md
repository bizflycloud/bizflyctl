# SSH Key Management Commands

Manage SSH keys for secure server access.

## Overview

SSH keys provide secure authentication to servers without passwords. You can create, list, and delete SSH keys in Bizfly Cloud.

## Commands

### List SSH Keys

List all SSH keys in your account:

```bash
bizfly ssh-key list
```

**Output:** Table showing:

-   Name
-   Fingerprint

### Create SSH Key

Create a new SSH key:

```bash
bizfly ssh-key create --name <key-name> --public-key <public-key>
```

**Required Flags:**

-   `--name`: SSH key name
-   `--public-key`: Public key content (or path to public key file)

**Example (with public key content):**

```bash
bizfly ssh-key create \
  --name my-key \
  --public-key "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAB..."
```

**Example (with public key file):**

```bash
bizfly ssh-key create \
  --name my-key \
  --public-key "$(cat ~/.ssh/id_rsa.pub)"
```

### Delete SSH Key

Delete an SSH key:

```bash
bizfly ssh-key delete <key-name>
```

**Example:**

```bash
bizfly ssh-key delete my-key
```

## Examples

### Complete SSH Key Workflow

```bash
# Generate SSH key pair (if you don't have one)
ssh-keygen -t rsa -b 4096 -C "your-email@example.com"

# List existing keys
bizfly ssh-key list

# Create SSH key in Bizfly Cloud
bizfly ssh-key create \
  --name production-key \
  --public-key "$(cat ~/.ssh/id_rsa.pub)"

# Use key when creating server
bizfly server create \
  --name my-server \
  --flavor nix.3c_6g \
  --image-id <image-id> \
  --rootdisk-size 40 \
  --ssh-key production-key

# Delete key (when no longer needed)
bizfly ssh-key delete production-key
```

### Using Existing SSH Key

```bash
# If you already have an SSH key
bizfly ssh-key create \
  --name my-existing-key \
  --public-key "$(cat ~/.ssh/id_rsa.pub)"
```

### Multiple Keys for Different Environments

```bash
# Production key
bizfly ssh-key create \
  --name prod-key \
  --public-key "$(cat ~/.ssh/prod_key.pub)"

# Development key
bizfly ssh-key create \
  --name dev-key \
  --public-key "$(cat ~/.ssh/dev_key.pub)"
```

## Common Use Cases

### Setting Up New Server Access

```bash
# 1. Generate new key pair
ssh-keygen -t rsa -b 4096 -f ~/.ssh/bizfly_key

# 2. Add to Bizfly Cloud
bizfly ssh-key create \
  --name bizfly-key \
  --public-key "$(cat ~/.ssh/bizfly_key.pub)"

# 3. Use when creating server
bizfly server create \
  --name new-server \
  --ssh-key bizfly-key \
  ...
```

## Troubleshooting

### Cannot Connect to Server

-   Verify SSH key is added to server: `bizfly server get <server-id>`
-   Check key name is correct
-   Ensure public key format is correct
-   Verify server has WAN IP assigned

### Key Already Exists

-   Use a different key name
-   Or delete existing key first: `bizfly ssh-key delete <name>`

## Best Practices

1. **Use strong keys** - Generate keys with at least 2048 bits (4096 recommended)
2. **Use descriptive names** - Include environment or purpose in name
3. **Keep private keys secure** - Never share private keys
4. **Rotate keys regularly** - Update keys periodically for security
5. **Use different keys** - Separate keys for different environments

## Related Commands

-   [Server Management](server.md) - Use SSH keys when creating servers
