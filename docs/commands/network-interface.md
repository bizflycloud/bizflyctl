# Network Interface Management Commands

Manage network interfaces for connecting servers to VPC networks.

## Overview

Network interfaces connect servers to VPC networks, allowing them to communicate within the VPC. You can attach multiple network interfaces to a server for different network configurations.

## Commands

### List Network Interfaces

List all network interfaces:

```bash
bizfly network-interface list [options]
```

**Optional Flags:**

-   `--vpc-network-id <id>`: Filter by VPC network ID
-   `--status <status>`: Filter by status
-   `--type <type>`: Filter by type

**Example:**

```bash
bizfly network-interface list
bizfly network-interface list --vpc-network-id vpc-123
```

**Output:** Table showing:

-   ID
-   Name
-   Status
-   Network ID
-   Device ID (Server ID)
-   IP Address
-   IP Version
-   Security Groups (Firewalls)
-   Created At
-   Updated At

### Get Network Interface Details

Get detailed information about a network interface:

```bash
bizfly network-interface get <network-interface-id>
```

**Example:**

```bash
bizfly network-interface get ni-123
```

### Create Network Interface

Create a new network interface:

```bash
bizfly network-interface create <vpc-network-id> [options]
```

**Required Arguments:**

-   `<vpc-network-id>`: The VPC network ID

**Optional Flags:**

-   `--name <name>`: Network interface name
-   `--fixed-ip <ip>`: Specific IP address to assign
-   `--server-id <id>`: Attach to server immediately

**Example:**

```bash
bizfly network-interface create vpc-123 \
  --name web-server-nic \
  --server-id server-456
```

### Delete Network Interface

Delete a network interface:

```bash
bizfly network-interface delete <network-interface-id>
```

**Note:** Network interface must be detached from server first.

**Example:**

```bash
bizfly network-interface delete ni-123
```

### Attach Server

Attach a server to a network interface:

```bash
bizfly network-interface attach-server <network-interface-id> --server-id <server-id>
```

**Example:**

```bash
bizfly network-interface attach-server ni-123 --server-id server-456
```

### Detach Server

Detach a server from a network interface:

```bash
bizfly network-interface detach-server <network-interface-id>
```

**Example:**

```bash
bizfly network-interface detach-server ni-123
```

### Add Firewalls

Add firewall rules to a network interface:

```bash
bizfly network-interface add-firewalls <network-interface-id> \
  --firewall-id <firewall-id1> \
  --firewall-id <firewall-id2>
```

**Example:**

```bash
bizfly network-interface add-firewalls ni-123 \
  --firewall-id fw-123 \
  --firewall-id fw-456
```

### Remove Firewalls

Remove firewall rules from a network interface:

```bash
bizfly network-interface remove-firewalls <network-interface-id> \
  --firewall-id <firewall-id1> \
  --firewall-id <firewall-id2>
```

**Example:**

```bash
bizfly network-interface remove-firewalls ni-123 \
  --firewall-id fw-123
```

## Examples

### Complete Network Interface Lifecycle

```bash
# List network interfaces
bizfly network-interface list

# Create network interface in VPC
bizfly network-interface create vpc-123 \
  --name app-nic \
  --server-id server-456

# Get network interface details
bizfly network-interface get ni-123

# Add firewall rules
bizfly network-interface add-firewalls ni-123 \
  --firewall-id web-firewall

# Detach from server
bizfly network-interface detach-server ni-123

# Delete network interface
bizfly network-interface delete ni-123
```

### Multi-NIC Server Setup

```bash
# Create server
bizfly server create --name multi-nic-server ...

# Create first network interface (production VPC)
bizfly network-interface create prod-vpc-123 \
  --name prod-nic \
  --server-id server-456

# Create second network interface (management VPC)
bizfly network-interface create mgmt-vpc-123 \
  --name mgmt-nic \
  --server-id server-456
```

### Network Interface with Specific IP

```bash
# Create network interface with fixed IP
bizfly network-interface create vpc-123 \
  --name db-nic \
  --fixed-ip 10.0.0.100 \
  --server-id db-server-456
```

## Common Use Cases

### Web Server in VPC

```bash
# Create network interface for web server
bizfly network-interface create web-vpc-123 \
  --name web-server-nic \
  --server-id web-server-456

# Add web firewall
bizfly network-interface add-firewalls ni-123 \
  --firewall-id web-firewall
```

### Database Server Isolation

```bash
# Create network interface in database VPC
bizfly network-interface create db-vpc-123 \
  --name db-nic \
  --server-id db-server-456

# Add database firewall (restrictive rules)
bizfly network-interface add-firewalls ni-123 \
  --firewall-id db-firewall
```

### Multi-Environment Access

```bash
# Server with access to both prod and staging
bizfly network-interface create prod-vpc-123 \
  --name prod-nic \
  --server-id shared-server-456

bizfly network-interface create staging-vpc-123 \
  --name staging-nic \
  --server-id shared-server-456
```

## Troubleshooting

### Cannot Create Network Interface

-   Verify VPC exists: `bizfly vpc get <vpc-id>`
-   Check VPC network ID is correct
-   Ensure you have available IP addresses in the VPC
-   Verify server is in the same region

### Cannot Attach to Server

-   Check server status (must be running or stopped)
-   Verify network interface is not already attached
-   Ensure server and VPC are in the same region
-   Check server limits for network interfaces

### Cannot Delete Network Interface

-   Detach from server first: `bizfly network-interface detach-server <id>`
-   Wait for detachment to complete
-   Check for any pending operations

### IP Address Conflicts

-   Use `--fixed-ip` to specify a specific IP
-   Or let the system assign automatically
-   Check existing IPs in the VPC

## Best Practices

1. **Use descriptive names** - Include server or purpose in name
2. **Plan IP addresses** - Use fixed IPs for important services
3. **Apply firewalls** - Add appropriate firewall rules for security
4. **One NIC per VPC** - Typically one network interface per VPC per server
5. **Clean up unused interfaces** - Delete interfaces you no longer need
6. **Document network topology** - Keep track of which servers are in which VPCs

## Related Commands

-   [VPC Management](vpc.md) - Create and manage VPCs
-   [Server Management](server.md) - Manage servers
-   [Firewall Management](firewall.md) - Configure firewall rules
