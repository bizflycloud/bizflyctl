# Firewall Management Commands

Manage firewall rules and security groups for your servers and network interfaces.

## Overview

Firewalls provide network-level security by controlling inbound and outbound traffic. You can create firewall rules and apply them to servers or network interfaces.

## Commands

### List Firewalls

List all firewalls in your account:

```bash
bizfly firewall list
```

**Output:** Table showing:

-   ID
-   Name
-   Description
-   Rules Count
-   Servers Count
-   Created At

### Get Firewall Details

Get detailed information about a firewall (includes rules):

```bash
bizfly firewall get <firewall-id>
```

**Example:**

```bash
bizfly firewall get 02b28284-5a18-4a0e-9ecc-d5d1acaf7e7b
```

### Create Firewall

Create a new firewall:

```bash
bizfly firewall create --name <firewall-name>
```

**Required Flags:**

-   `--name`: Firewall name

**Example:**

```bash
bizfly firewall create --name web-firewall
```

### Delete Firewall

Delete one or more firewalls:

```bash
bizfly firewall delete <firewall-id> [firewall-id2] [firewall-id3] ...
```

**Example:**

```bash
bizfly firewall delete fw-123 fw-456
```

## Firewall Rules

### List Rules

List all rules in a firewall:

```bash
bizfly firewall rule list <firewall-id>
```

**Output:** Table showing:

-   ID
-   Description
-   Direction (ingress/egress)
-   Type
-   Ether Type
-   Protocol
-   CIDR
-   Port Range
-   Remote IP Prefix

### Create Rule

Create a new firewall rule:

```bash
bizfly firewall rule create <firewall-id> \
  --direction <ingress|egress> \
  --protocol <tcp|udp> \
  [options]
```

**Required Flags:**

-   `--direction`: `ingress` (inbound) or `egress` (outbound)
-   `--protocol`: `tcp` or `udp`

**Optional Flags:**

-   `--port-range <range>`: Port or port range (e.g., `80`, `80-90`) - default: `1-65535`
-   `--cidr <cidr>`: CIDR block (e.g., `10.0.0.0/24`) - default: `0.0.0.0/0`

**Examples:**

Allow HTTP traffic:

```bash
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 80 \
  --cidr 0.0.0.0/0
```

Allow SSH from specific IP:

```bash
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 22 \
  --cidr 203.0.113.0/24
```

Allow port range:

```bash
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 8000-8999 \
  --cidr 10.0.0.0/16
```

### Delete Rule

Delete a firewall rule:

```bash
bizfly firewall rule delete <firewall-id> <rule-id>
```

**Example:**

```bash
bizfly firewall rule delete fw-123 rule-456
```

## Firewall-Server Management

### List Servers

List servers that have a firewall applied:

```bash
bizfly firewall server list <firewall-id>
```

**Output:** Table showing:

-   ID (Server ID)
-   Name (Server Name)
-   Firewall ID

### Remove Servers

Remove servers from a firewall:

```bash
bizfly firewall server remove <firewall-id> <server-id1> <server-id2> ...
```

**Example:**

```bash
bizfly firewall server remove fw-123 server-456 server-789
```

## Examples

### Complete Firewall Workflow

```bash
# Create firewall
bizfly firewall create --name web-firewall

# Add rules
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 80 \
  --cidr 0.0.0.0/0

bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 443 \
  --cidr 0.0.0.0/0

# List rules
bizfly firewall rule list fw-123

# Apply to server (via network interface)
bizfly network-interface add-firewalls ni-123 \
  --firewall-id fw-123
```

### Web Server Firewall

```bash
# Create web firewall
bizfly firewall create --name web-server-firewall

# Allow HTTP
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 80 \
  --cidr 0.0.0.0/0

# Allow HTTPS
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 443 \
  --cidr 0.0.0.0/0

# Allow SSH from office
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 22 \
  --cidr 203.0.113.0/24
```

### Database Firewall (Restrictive)

```bash
# Create database firewall
bizfly firewall create --name database-firewall

# Allow PostgreSQL from app servers only
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 5432 \
  --cidr 10.0.1.0/24

# Allow SSH from management network
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 22 \
  --cidr 10.0.0.0/24
```

### Application Firewall

```bash
# Create application firewall
bizfly firewall create --name app-firewall

# Allow application port
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 8080 \
  --cidr 10.0.0.0/16

# Allow health check port
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 9090 \
  --cidr 10.0.0.0/16
```

## Common Use Cases

### Public Web Server

```bash
# Firewall allowing public HTTP/HTTPS
bizfly firewall create --name public-web

# HTTP
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 80

# HTTPS
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 443
```

### Internal Service

```bash
# Firewall for internal services only
bizfly firewall create --name internal-service

# Allow from internal network only
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 8080 \
  --cidr 10.0.0.0/16
```

### Development Environment

```bash
# More permissive firewall for development
bizfly firewall create --name dev-firewall

# Allow all from internal network
bizfly firewall rule create fw-123 \
  --direction ingress \
  --protocol tcp \
  --port-range 1-65535 \
  --cidr 10.0.0.0/16
```

## Port Range Examples

-   **Single port:** `80`
-   **Port range:** `8000-8999`
-   **All ports:** `1-65535` (default)

## CIDR Examples

-   **All IPs:** `0.0.0.0/0`
-   **Single IP:** `203.0.113.1/32`
-   **Subnet:** `10.0.0.0/24`
-   **Network:** `10.0.0.0/16`

## Troubleshooting

### Cannot Create Firewall Rule

-   Verify firewall exists
-   Check direction is `ingress` or `egress`
-   Ensure protocol is `tcp` or `udp`
-   Verify CIDR format is correct

### Rule Not Working

-   Check rule direction (ingress vs egress)
-   Verify port range includes the port you're using
-   Check CIDR includes your source/destination IP
-   Ensure firewall is applied to network interface

### Cannot Delete Firewall

-   Remove firewall from all network interfaces first
-   Check if firewall is in use
-   Wait for any pending operations

## Best Practices

1. **Principle of least privilege** - Only allow necessary ports and IPs
2. **Use descriptive names** - Include purpose in firewall name
3. **Separate firewalls by purpose** - Web, database, application, etc.
4. **Document rules** - Use descriptions or comments to explain rules
5. **Test rules** - Verify firewall rules work as expected
6. **Regular review** - Periodically review and clean up unused rules
7. **Use CIDR blocks** - Instead of allowing `0.0.0.0/0` when possible

## Security Considerations

-   **Default deny** - Firewalls deny all traffic by default
-   **Explicit allow** - Only explicitly allowed traffic is permitted
-   **Restrict SSH** - Limit SSH access to trusted IPs
-   **Separate environments** - Use different firewalls for prod/staging/dev
-   **Regular audits** - Review firewall rules regularly

## Related Commands

-   [Network Interface Management](network-interface.md) - Apply firewalls to network interfaces
-   [Server Management](server.md) - Create servers with firewalls
-   [VPC Management](vpc.md) - Network isolation
