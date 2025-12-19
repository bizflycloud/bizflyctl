# Snapshot Management Commands

Create and manage volume snapshots for backup and disaster recovery.

## Overview

Snapshots are point-in-time copies of volumes. They can be used to:

-   Backup data
-   Create new volumes
-   Create new servers
-   Restore volumes to previous states

## Commands

### List Snapshots

List all snapshots in your account:

```bash
bizfly snapshot list [options]
```

**Optional Flags:**

-   `--volume-id <id>`: Filter snapshots by volume ID

**Example:**

```bash
bizfly snapshot list
bizfly snapshot list --volume-id vol-123
```

**Output:** Table showing:

-   ID
-   Name
-   Status
-   Size (GB)
-   Type
-   Created At
-   Volume ID
-   Billing Plan
-   Zone

### Get Snapshot Details

Get detailed information about a specific snapshot:

```bash
bizfly snapshot get <snapshot-id>
```

**Example:**

```bash
bizfly snapshot get 5af19947-566d-48a1-bc45-93666086951f
```

### Create Snapshot

Create a snapshot of a volume:

```bash
bizfly snapshot create <volume-id> --name <snapshot-name>
```

**Required Arguments:**

-   `<volume-id>`: The volume ID to snapshot

**Required Flags:**

-   `--name`: Snapshot name

**Example:**

```bash
bizfly snapshot create vol-123 --name daily-backup-2024-01-15
```

**Note:** The snapshot is created with `force: true`, which means it will be created even if the volume is attached to a running server.

### Delete Snapshot

Delete one or more snapshots:

```bash
bizfly snapshot delete <snapshot-id> [snapshot-id2] [snapshot-id3] ...
```

**Example:**

```bash
bizfly snapshot delete snap-123
```

**Example (multiple snapshots):**

```bash
bizfly snapshot delete snap-123 snap-456 snap-789
```

## Examples

### Basic Snapshot Workflow

```bash
# List all snapshots
bizfly snapshot list

# Create snapshot
bizfly snapshot create vol-123 --name backup-before-upgrade

# Get snapshot details
bizfly snapshot get snap-456

# Delete old snapshot
bizfly snapshot delete snap-789
```

### Backup Strategy

```bash
# Create daily backup
bizfly snapshot create vol-123 --name daily-$(date +%Y-%m-%d)

# List snapshots for a specific volume
bizfly snapshot list --volume-id vol-123

# Keep only last 7 days of snapshots
# (Delete older ones manually or via script)
```

### Using Snapshots to Create New Resources

```bash
# Create snapshot
bizfly snapshot create vol-123 --name production-backup

# Create new volume from snapshot
bizfly volume create \
  --name restored-data \
  --size 100 \
  --snapshot-id snap-456

# Create new server from snapshot
bizfly server create \
  --name restored-server \
  --flavor nix.3c_6g \
  --snapshot-id snap-456 \
  --rootdisk-size 50
```

### Restore Volume from Snapshot

```bash
# Restore existing volume
bizfly volume restore vol-123 --snapshot-id snap-456

# Or create new volume from snapshot
bizfly volume create \
  --name restored-volume \
  --size 100 \
  --snapshot-id snap-456
```

## Common Use Cases

### Pre-Upgrade Backup

```bash
# Before upgrading application
bizfly snapshot create vol-123 --name pre-upgrade-$(date +%Y%m%d)

# Perform upgrade...

# If upgrade fails, restore
bizfly volume restore vol-123 --snapshot-id snap-456
```

### Database Backup

```bash
# Create database backup snapshot
bizfly snapshot create db-volume-123 --name db-backup-$(date +%Y-%m-%d)

# Create development database from production backup
bizfly volume create \
  --name dev-db \
  --size 200 \
  --snapshot-id snap-456
```

### Disaster Recovery

```bash
# Regular backup schedule
bizfly snapshot create vol-123 --name weekly-backup-$(date +%Y-W%V)

# In case of disaster, restore from latest snapshot
bizfly volume restore vol-123 --snapshot-id snap-latest
```

## Best Practices

1. **Name snapshots descriptively** - Include date, purpose, or version
2. **Create snapshots before major changes** - Upgrades, migrations, etc.
3. **Regular backup schedule** - Automate with scripts or scheduled backups
4. **Clean up old snapshots** - Delete snapshots you no longer need
5. **Test restore procedures** - Periodically verify you can restore from snapshots
6. **Document snapshot purpose** - Use descriptive names and notes

## Troubleshooting

### Snapshot Creation Fails

-   Verify volume exists: `bizfly volume get <volume-id>`
-   Check volume status (should be available)
-   Ensure you have sufficient quota
-   Wait for any pending volume operations to complete

### Cannot Delete Snapshot

-   Check if snapshot is being used by any volumes or servers
-   Verify snapshot status
-   Wait for any pending operations

### Snapshot Status Stuck

If snapshot status is stuck:

1. Check volume status
2. Wait a few minutes for operation to complete
3. If still stuck, contact support

## Cost Considerations

-   Snapshots consume storage space
-   You're charged for snapshot storage
-   Delete old snapshots to reduce costs
-   Consider snapshot retention policies

## Related Commands

-   [Volume Management](volume.md) - Create volumes from snapshots
-   [Server Management](server.md) - Create servers from snapshots
-   [Scheduled Volume Backup](scheduled-volume-backup.md) - Automate snapshot creation
