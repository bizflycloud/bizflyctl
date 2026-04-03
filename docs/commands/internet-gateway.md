# Internet Gateway Management Commands

Manage Internet Gateways for connecting VPC networks to the internet.

## Overview

Internet Gateways provide connectivity between VPC networks and the internet. Each VPC can have one Internet Gateway attached for outbound and inbound internet access.

## Commands

### List Internet Gateways

List all Internet Gateways:

```bash
bizfly internet-gateway list [options]
```

**Optional Flags:**

-   `--name <name>`: Filter by Internet Gateway name

**Example:**

```bash
bizfly internet-gateway list
bizfly internet-gateway list --name my-gateway
```

**Output:** Table showing:

-   ID
-   NAME
-   STATUS
-   VPC NAMES
-   AVAILABILITY ZONES
-   CREATED AT

### Get Internet Gateway Details

Get detailed information about an Internet Gateway:

```bash
bizfly internet-gateway get <internet-gateway-id>
```

**Example:**

```bash
bizfly internet-gateway get igw-123
```

### Create Internet Gateway

Create a new Internet Gateway:

```bash
bizfly internet-gateway create [options]
```

**Required Flags:**

-   `--name <name>`: Internet Gateway name

**Optional Flags:**

-   `--network-id <network-id>`: VPC network ID to attach (can specify multiple)
-   `--description <text>`: Internet Gateway description
-   `--availability-zone <zone>`: Availability zone for the Internet Gateway

**Example:**

```bash
# Create Internet Gateway without VPC
bizfly internet-gateway create --name my-igw

# Create Internet Gateway with VPC
bizfly internet-gateway create \
  --name my-igw \
  --network-id vpc-123

# Create Internet Gateway with availability zone
bizfly internet-gateway create \
  --name my-igw \
  --network-id vpc-123 \
  --availability-zone HN1
```

### Update Internet Gateway

Update an existing Internet Gateway:

```bash
bizfly internet-gateway update <internet-gateway-id> [options]
```

**Optional Flags:**

-   `--name <name>`: New Internet Gateway name
-   `--network-id <network-id>`: VPC network ID (can specify multiple)
-   `--description <text>`: New description

**Example:**

```bash
# Update name and description
bizfly internet-gateway update igw-123 \
  --name new-igw-name \
  --description "Updated description"

# Update VPC attachment
bizfly internet-gateway update igw-123 \
  --network-id vpc-456
```

### Detach VPC

Detach all VPCs from an Internet Gateway:

```bash
bizfly internet-gateway detach-vpc <internet-gateway-id>
```

**Example:**

```bash
bizfly internet-gateway detach-vpc igw-123
```

### Delete Internet Gateway

Delete an Internet Gateway:

```bash
bizfly internet-gateway delete <internet-gateway-id>
```

**Example:**

```bash
bizfly internet-gateway delete igw-123
```

## Examples

### Basic Internet Gateway Workflow

```bash
# List all Internet Gateways
bizfly internet-gateway list

# Create an Internet Gateway
bizfly internet-gateway create \
  --name production-igw \
  --description "Production Internet Gateway"

# Get Internet Gateway details
bizfly internet-gateway get igw-123

# Attach VPC
bizfly internet-gateway update igw-123 --network-id vpc-123

# Detach VPC (when reconfiguring network)
bizfly internet-gateway detach-vpc igw-123

# Delete Internet Gateway
bizfly internet-gateway delete igw-123
```

### Multi-VPC Setup

```bash
# Create separate Internet Gateways for each environment
bizfly internet-gateway create \
  --name prod-igw \
  --network-id prod-vpc-123

bizfly internet-gateway create \
  --name staging-igw \
  --network-id staging-vpc-123

bizfly internet-gateway create \
  --name dev-igw \
  --network-id dev-vpc-123
```

### Reconfiguring VPC Attachment

```bash
# Detach current VPC
bizfly internet-gateway detach-vpc igw-123

# Attach new VPC
bizfly internet-gateway update igw-123 --network-id new-vpc-456
```

## Common Use Cases

### Internet Access for VPC

```bash
# Create VPC
bizfly vpc create --name my-vpc --cidr 10.0.0.0/16

# Create Internet Gateway
bizfly internet-gateway create \
  --name my-igw \
  --network-id vpc-123

# The VPC can now access the internet through the gateway
```

### Isolated Environments

```bash
# Production environment
bizfly internet-gateway create \
  --name prod-gateway \
  --network-id prod-vpc

# Development environment
bizfly internet-gateway create \
  --name dev-gateway \
  --network-id dev-vpc
```

### High Availability Setup

```bash
# Create Internet Gateways in different AZs
bizfly internet-gateway create \
  --name ha-gateway-1 \
  --network-id vpc-hn1 \
  --availability-zone HN1

bizfly internet-gateway create \
  --name ha-gateway-2 \
  --network-id vpc-hn2 \
  --availability-zone hn2
```

## Troubleshooting

### Cannot Create Internet Gateway

-   Verify VPC exists: `bizfly vpc get <vpc-id>`
-   Check VPC network ID is correct
-   Ensure you have permissions to create Internet Gateways
-   Verify availability zone exists

### Cannot Update Internet Gateway

-   Check Internet Gateway exists: `bizfly internet-gateway get <id>`
-   Verify new VPC network ID is valid
-   Ensure no operations are in progress

### Cannot Delete Internet Gateway

-   Detach all VPCs first: `bizfly internet-gateway detach-vpc <id>`
-   Check for any active connections using the gateway
-   Wait for any pending operations to complete

### VPC Not Accessible from Internet

-   Verify Internet Gateway is attached to VPC
-   Check route tables are configured correctly
-   Ensure security groups allow traffic
-   Verify firewall rules are not blocking access

## Best Practices

1. **One IGW per VPC** - Each VPC typically needs one Internet Gateway
2. **Use descriptive names** - Include environment or purpose in name
3. **Plan availability zones** - Consider AZ redundancy for production
4. **Clean up unused gateways** - Delete Internet Gateways you no longer need
5. **Document network topology** - Keep track of which VPCs are connected
6. **Monitor usage** - Track Internet Gateway activity and metrics

## Related Commands

-   [VPC Management](vpc.md) - Create and manage VPCs
-   [Network Interface Management](network-interface.md) - Manage network interfaces
-   [Server Management](server.md) - Manage servers in VPCs
