# VPC Management Commands

Manage Virtual Private Cloud (VPC) networks for network isolation and security.

## Overview

VPCs provide isolated network environments for your resources. You can create multiple VPCs to separate different environments or applications.

## Commands

### List VPCs

List all VPCs in your account:

```bash
bizfly vpc list
```

**Output:** Table showing:

-   ID
-   Name
-   MTU
-   CIDR
-   Description
-   Tags
-   Created At
-   Is Default
-   Zones (Availability Zones)

### Get VPC Details

Get detailed information about a specific VPC:

```bash
bizfly vpc get <vpc-id>
```

**Example:**

```bash
bizfly vpc get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
```

### Create VPC

Create a new VPC:

```bash
bizfly vpc create --name <vpc-name> [options]
```

**Required Flags:**

-   `--name`: VPC name

**Optional Flags:**

-   `--description <text>`: VPC description
-   `--cidr <cidr>`: CIDR block for the VPC (e.g., `10.0.0.0/16`)
-   `--is-default <true|false>`: Set as default VPC - default: `false`

**Example:**

```bash
bizfly vpc create \
  --name production-vpc \
  --description "Production environment VPC" \
  --cidr 10.0.0.0/16
```

### Update VPC

Update VPC properties:

```bash
bizfly vpc update <vpc-id> [options]
```

**Optional Flags:**

-   `--name <name>`: New VPC name
-   `--description <text>`: New description
-   `--cidr <cidr>`: New CIDR block
-   `--is-default <true|false>`: Set as default VPC

**Example:**

```bash
bizfly vpc update vpc-123 \
  --name updated-vpc-name \
  --description "Updated description"
```

### Delete VPC

Delete a VPC:

```bash
bizfly vpc delete <vpc-id>
```

**Warning:** Ensure no resources are using the VPC before deletion.

**Example:**

```bash
bizfly vpc delete vpc-123
```

## Examples

### Basic VPC Workflow

```bash
# List all VPCs
bizfly vpc list

# Create a new VPC
bizfly vpc create \
  --name development-vpc \
  --description "Development environment" \
  --cidr 10.1.0.0/16

# Get VPC details
bizfly vpc get vpc-123

# Update VPC
bizfly vpc update vpc-123 \
  --description "Updated development VPC"

# Delete VPC (when no longer needed)
bizfly vpc delete vpc-123
```

### Multi-Environment Setup

```bash
# Production VPC
bizfly vpc create \
  --name prod-vpc \
  --description "Production environment" \
  --cidr 10.0.0.0/16

# Staging VPC
bizfly vpc create \
  --name staging-vpc \
  --description "Staging environment" \
  --cidr 10.1.0.0/16

# Development VPC
bizfly vpc create \
  --name dev-vpc \
  --description "Development environment" \
  --cidr 10.2.0.0/16
```

### Creating Servers in VPC

```bash
# Create server with network interface in VPC
bizfly network-interface create <vpc-network-id> \
  --name prod-network-interface \
  --server-id server-123

# Or attach VPC to existing server
bizfly server add-vpc server-123 --vpc-ids vpc-123
```

## Common Use Cases

### Environment Isolation

```bash
# Create separate VPCs for each environment
bizfly vpc create --name prod --cidr 10.0.0.0/16
bizfly vpc create --name staging --cidr 10.1.0.0/16
bizfly vpc create --name dev --cidr 10.2.0.0/16
```

### Application Isolation

```bash
# Separate VPCs for different applications
bizfly vpc create --name web-app-vpc --cidr 10.10.0.0/16
bizfly vpc create --name api-vpc --cidr 10.20.0.0/16
bizfly vpc create --name database-vpc --cidr 10.30.0.0/16
```

## CIDR Planning

When creating VPCs, plan your CIDR blocks to avoid conflicts:

-   **Production:** `10.0.0.0/16` (65,536 IPs)
-   **Staging:** `10.1.0.0/16` (65,536 IPs)
-   **Development:** `10.2.0.0/16` (65,536 IPs)

Or use smaller subnets:

-   **Small VPC:** `10.0.0.0/24` (256 IPs)
-   **Medium VPC:** `10.0.0.0/20` (4,096 IPs)

## Troubleshooting

### Cannot Delete VPC

-   Ensure no servers are attached to the VPC
-   Check for network interfaces using the VPC
-   Verify no load balancers are using the VPC
-   Wait for all resources to be removed

### CIDR Conflict

-   Choose non-overlapping CIDR blocks
-   Use different IP ranges for different VPCs
-   Plan your IP space before creating VPCs

### VPC Not Visible

-   Check you're in the correct region
-   Verify VPC was created successfully
-   List all VPCs to confirm

## Best Practices

1. **Plan CIDR blocks** - Avoid overlapping IP ranges
2. **Use descriptive names** - Include environment or purpose
3. **Document VPC purposes** - Use descriptions effectively
4. **Isolate environments** - Separate production, staging, dev
5. **Clean up unused VPCs** - Delete VPCs you no longer need
6. **Use default VPC wisely** - Only one default VPC per region

## Related Commands

-   [Network Interface Management](network-interface.md) - Create network interfaces in VPCs
-   [Server Management](server.md) - Attach servers to VPCs
-   [Firewall Management](firewall.md) - Configure firewall rules for VPCs
