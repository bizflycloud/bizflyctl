# Volume Management Commands

Manage block storage volumes in Bizfly Cloud.

## Overview

Volumes are persistent block storage devices that can be attached to servers. They persist independently of server lifecycles.

## Commands

### List Volumes

List all volumes in your account:

```bash
bizfly volume list
```

**Output:** Table showing:

-   ID
-   Name
-   Description
-   Status
-   Size (GB)
-   Created At
-   Volume Type
-   Snapshot ID
-   Billing Plan
-   Zone (Availability Zone)
-   Attached Server

### Get Volume Details

Get detailed information about a specific volume:

```bash
bizfly volume get <volume-id>
```

**Example:**

```bash
bizfly volume get 9e580b1a-0526-460b-9a6f-d8f80130bda8
```

### Create Volume

Create a new volume:

```bash
bizfly volume create \
  --name <volume-name> \
  --size <size-in-gb> \
  [options]
```

**Required Flags:**

-   `--name`: Volume name
-   `--size`: Volume size in GB

**Optional Flags:**

-   `--type <type>`: Volume type (`HDD` or `SSD`) - default: `HDD`
-   `--category <category>`: Volume category (`premium`, `enterprise`, `basic`) - default: `premium`
-   `--availability-zone <zone>`: Availability zone (e.g., `HN1`) - default: `HN1`
-   `--snapshot-id <id>`: Create volume from snapshot
-   `--server-id <id>`: Create and attach to server immediately
-   `--description <text>`: Volume description
-   `--billing-plan <plan>`: Billing plan (`saving_plan` or `on_demand`) - default: `saving_plan`

**Examples:**

Create a basic volume:

```bash
bizfly volume create \
  --name my-data-volume \
  --size 100 \
  --type SSD \
  --category premium
```

Create volume from snapshot:

```bash
bizfly volume create \
  --name restored-volume \
  --size 50 \
  --snapshot-id snap-123
```

Create and attach to server:

```bash
bizfly volume create \
  --name app-data \
  --size 200 \
  --server-id server-123
```

### Delete Volume

Delete one or more volumes:

```bash
bizfly volume delete <volume-id> [volume-id2] [volume-id3] ...
```

**Note:** Volumes must be detached from servers before deletion.

**Example:**

```bash
bizfly volume delete vol-123 vol-456
```

### Attach Volume

Attach a volume to a server:

```bash
bizfly volume attach <volume-id> <server-id>
```

**Example:**

```bash
bizfly volume attach vol-123 server-456
```

### Detach Volume

Detach a volume from a server:

```bash
bizfly volume detach <volume-id> <server-id>
```

**Example:**

```bash
bizfly volume detach vol-123 server-456
```

### Extend Volume

Increase the size of a volume:

```bash
bizfly volume extend <volume-id> --size <new-size-in-gb>
```

**Note:** You can only increase volume size, not decrease it.

**Example:**

```bash
bizfly volume extend vol-123 --size 200
```

### Restore Volume

Restore a volume from a snapshot:

```bash
bizfly volume restore <volume-id> --snapshot-id <snapshot-id>
```

**Warning:** This operation will overwrite all data on the volume.

**Example:**

```bash
bizfly volume restore vol-123 --snapshot-id snap-456
```

### Update Volume

Update volume metadata:

```bash
bizfly volume patch <volume-id> --description <new-description>
```

**Example:**

```bash
bizfly volume patch vol-123 --description "Updated description"
```

### List Volume Types

List available volume types:

```bash
bizfly volume list-types [options]
```

**Optional Flags:**

-   `--category <category>`: Filter by category
-   `--availability-zone <zone>`: Filter by availability zone

**Example:**

```bash
bizfly volume list-types --category premium --availability-zone HN1
```

**Output:** Table showing:

-   Type
-   Category
-   Availability Zones

## Examples

### Complete Volume Lifecycle

```bash
# List available volume types
bizfly volume list-types

# Create a volume
bizfly volume create \
  --name database-storage \
  --size 500 \
  --type SSD \
  --category enterprise

# Attach to server
bizfly volume attach vol-123 server-456

# Extend volume
bizfly volume extend vol-123 --size 1000

# Create snapshot (see snapshot commands)
bizfly snapshot create vol-123 --name backup

# Detach volume
bizfly volume detach vol-123 server-456

# Delete volume
bizfly volume delete vol-123
```

### Backup and Restore Workflow

```bash
# Create snapshot of volume
bizfly snapshot create vol-123 --name daily-backup

# Restore volume from snapshot
bizfly volume restore vol-123 --snapshot-id snap-456

# Or create new volume from snapshot
bizfly volume create \
  --name restored-data \
  --size 100 \
  --snapshot-id snap-456
```

## Common Use Cases

### Database Storage

```bash
# Create large SSD volume for database
bizfly volume create \
  --name postgres-data \
  --size 1000 \
  --type SSD \
  --category enterprise \
  --billing-plan saving_plan

# Attach to database server
bizfly volume attach vol-123 db-server-456
```

### Application Data Storage

```bash
# Create volume for application
bizfly volume create \
  --name app-files \
  --size 200 \
  --type HDD \
  --category premium

# Attach to application server
bizfly volume attach vol-123 app-server-456
```

### Disaster Recovery

```bash
# Create snapshot before major changes
bizfly snapshot create vol-123 --name pre-upgrade-backup

# If something goes wrong, restore
bizfly volume restore vol-123 --snapshot-id snap-456
```

## Troubleshooting

### Cannot Delete Volume

-   Ensure volume is detached from all servers
-   Check volume status: `bizfly volume get <volume-id>`
-   Wait for any pending operations to complete

### Cannot Attach Volume

-   Verify volume and server are in the same availability zone
-   Check server status (must be running or stopped, not in error state)
-   Ensure volume is not already attached to another server

### Volume Extension Fails

-   Verify new size is larger than current size
-   Check available quota
-   Ensure volume is in a valid state

### Volume Not Visible in Server

After attaching a volume:

1. Check volume attachment: `bizfly volume get <volume-id>`
2. On Linux, use `lsblk` or `fdisk -l` to see new block device
3. May need to format and mount the volume

## Best Practices

1. **Use snapshots regularly** for important data
2. **Choose appropriate volume type** (SSD for performance, HDD for cost)
3. **Select correct availability zone** matching your servers
4. **Use appropriate billing plan** based on usage patterns
5. **Detach volumes** before deleting servers if you want to keep data
6. **Monitor volume usage** and extend before running out of space

## Related Commands

-   [Snapshot Management](snapshot.md) - Create and manage snapshots
-   [Server Management](server.md) - Attach volumes to servers
-   [Scheduled Volume Backup](scheduled-volume-backup.md) - Automate backups
