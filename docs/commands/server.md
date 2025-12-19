# Server Management Commands

Manage your Bizfly Cloud servers (virtual machines) using these commands.

## Overview

The `server` command group provides operations to create, list, manage, and delete cloud servers.

## Commands

### List Servers

List all servers in your account:

```bash
bizfly server list
```

**Output:** Table showing:

-   ID
-   Name
-   Zone (Availability Zone)
-   Key Name (SSH key)
-   Status
-   Flavor
-   Category
-   LAN IP
-   WAN IP
-   Created At

### Get Server Details

Get detailed information about a specific server:

```bash
bizfly server get <server-id>
```

**Example:**

```bash
bizfly server get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
```

**Output:** Detailed information including attached volumes.

### Create Server

Create a new server:

```bash
bizfly server create \
  --name <server-name> \
  --flavor <flavor-name> \
  --rootdisk-size <size-in-gb> \
  [options]
```

**Required Flags:**

-   `--name`: Server name
-   `--flavor`: Flavor name (e.g., `nix.3c_6g`)
-   `--rootdisk-size`: Root disk size in GB (minimum 20GB)

**Optional Flags:**

-   `--image-id <id>`: Create from an OS image
-   `--volume-id <id>`: Create from an existing volume
-   `--snapshot-id <id>`: Create from a snapshot
-   `--category <type>`: Server category (`basic`, `premium`, `enterprise`) - default: `premium`
-   `--availability-zone <zone>`: Availability zone (e.g., `HN1`) - default: `HN1`
-   `--rootdisk-type <type>`: Root disk type (`HDD` or `SSD`) - default: `HDD`
-   `--rootdisk-volume-type <type>`: Root disk volume type (e.g., `PREMIUM-HDD1`)
-   `--ssh-key <key-name>`: SSH key name
-   `--network-plan <plan>`: Network plan (`free_bandwidth` or `free_datatransfer`)
-   `--net-interface <id>`: Network interface IDs (can be specified multiple times)
-   `--firewall <id>`: Firewall IDs (can be specified multiple times)
-   `--billing-plan <plan>`: Billing plan (`saving_plan` or `on_demand`) - default: `saving_plan`
-   `--is-created-wan-ip <true|false>`: Create WAN IP - default: `true`

**Examples:**

Create server from image:

```bash
bizfly server create \
  --name my-web-server \
  --flavor nix.3c_6g \
  --image-id abc123-def456-ghi789 \
  --rootdisk-size 40 \
  --category premium \
  --ssh-key my-ssh-key
```

Create server from snapshot:

```bash
bizfly server create \
  --name restored-server \
  --flavor nix.3c_6g \
  --snapshot-id snap-123 \
  --rootdisk-size 50
```

**Output:** Returns a task ID for server creation.

### Delete Server

Delete one or more servers:

```bash
bizfly server delete <server-id> [server-id2] [server-id3] ...
```

**Options:**

-   `--delete-rootdisk <true|false>`: Delete root disk with server - default: `true`

**Example:**

```bash
bizfly server delete fd554aac-9ab1-11ea-b09d-bbaf82f02f58
```

**Example (multiple servers):**

```bash
bizfly server delete server-id-1 server-id-2 server-id-3
```

### Server Actions

#### Start Server

Start a stopped server:

```bash
bizfly server start <server-id>
```

#### Stop Server

Stop a running server:

```bash
bizfly server stop <server-id>
```

#### Reboot Server (Soft)

Perform a soft reboot:

```bash
bizfly server reboot <server-id>
```

#### Hard Reboot

Perform a hard reboot:

```bash
bizfly server hard reboot <server-id>
```

#### Resize Server

Resize a server to a different flavor:

```bash
bizfly server resize <server-id> --flavor <new-flavor-name>
```

**Example:**

```bash
bizfly server resize server-123 --flavor nix.6c_12g
```

#### Rename Server

Rename a server:

```bash
bizfly server rename <server-id> --name <new-name>
```

**Example:**

```bash
bizfly server rename server-123 --name my-new-server-name
```

### Network Management

#### Add VPC to Server

Attach VPC networks to a server:

```bash
bizfly server add-vpc <server-id> --vpc-ids <vpc-id1> --vpc-ids <vpc-id2>
```

**Example:**

```bash
bizfly server add-vpc server-123 --vpc-ids vpc-1 --vpc-ids vpc-2
```

#### Remove VPC from Server

Detach VPC networks from a server:

```bash
bizfly server remove-vpc <server-id> --vpc-ids <vpc-id1> --vpc-ids <vpc-id2>
```

**Example:**

```bash
bizfly server remove-vpc server-123 --vpc-ids vpc-1
```

#### Change Network Plan

Change the network plan of a server:

```bash
bizfly server change-network-plan <server-id> --network-plan <plan>
```

**Options:**

-   `--network-plan`: `free_bandwidth` or `free_datatransfer`

**Example:**

```bash
bizfly server change-network-plan server-123 --network-plan free_datatransfer
```

### Billing Management

#### Switch Billing Plan

Switch between billing plans:

```bash
bizfly server switch-billing-plan <server-id> --billing-plan <plan>
```

**Options:**

-   `--billing-plan`: `saving_plan` or `on_demand`

**Example:**

```bash
bizfly server switch-billing-plan server-123 --billing-plan on_demand
```

### List Server Types

List available server types:

```bash
bizfly server list-types
```

**Output:** Table showing:

-   ID
-   Name
-   Enabled status
-   Compute class

## Examples

### Complete Server Lifecycle

```bash
# List available flavors
bizfly flavor list

# Create a server
bizfly server create \
  --name production-web \
  --flavor nix.3c_6g \
  --image-id ubuntu-20.04 \
  --rootdisk-size 50 \
  --category premium \
  --ssh-key my-key

# Get server details
bizfly server get <server-id>

# Resize server
bizfly server resize <server-id> --flavor nix.6c_12g

# Stop server
bizfly server stop <server-id>

# Start server
bizfly server start <server-id>

# Delete server
bizfly server delete <server-id>
```

## Common Use Cases

### Creating a Development Server

```bash
bizfly server create \
  --name dev-server \
  --flavor nix.2c_4g \
  --image-id ubuntu-20.04 \
  --rootdisk-size 30 \
  --category basic \
  --ssh-key dev-key
```

### Creating a Production Server with Firewall

```bash
bizfly server create \
  --name prod-server \
  --flavor nix.6c_12g \
  --image-id ubuntu-20.04 \
  --rootdisk-size 100 \
  --category enterprise \
  --firewall web-firewall \
  --firewall db-firewall \
  --ssh-key prod-key
```

### Cloning a Server from Snapshot

```bash
# Create snapshot first (see snapshot commands)
bizfly snapshot create <volume-id> --name backup-snapshot

# Create new server from snapshot
bizfly server create \
  --name cloned-server \
  --flavor nix.3c_6g \
  --snapshot-id <snapshot-id> \
  --rootdisk-size 50
```

## Troubleshooting

### Server Creation Fails

-   Verify flavor name is correct: `bizfly flavor list`
-   Check image ID exists: `bizfly image list`
-   Ensure root disk size meets minimum (20GB)
-   Verify availability zone is correct

### Server Won't Start

-   Check server status: `bizfly server get <server-id>`
-   Verify billing plan is active
-   Check for resource limits

### Cannot Connect via SSH

-   Verify SSH key is correct: `bizfly sshkey list`
-   Check firewall rules allow SSH (port 22)
-   Verify WAN IP is assigned: `bizfly server get <server-id>`

## Related Commands

-   [Volume Management](volume.md) - Attach/detach volumes
-   [Snapshot Management](snapshot.md) - Create server snapshots
-   [Firewall Management](firewall.md) - Configure firewall rules
-   [VPC Management](vpc.md) - Manage VPC networks
-   [SSH Key Management](sshkey.md) - Manage SSH keys
-   [Flavor Management](flavor.md) - List available flavors
