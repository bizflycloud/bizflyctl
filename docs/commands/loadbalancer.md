# Load Balancer Management Commands

Manage load balancers, listeners, pools, and health monitors for distributing traffic across servers.

## Overview

Load balancers distribute incoming traffic across multiple servers to improve availability and performance. This guide covers managing load balancers and their components.

## Load Balancer Commands

### List Load Balancers

List all load balancers:

```bash
bizfly loadbalancer list
```

**Output:** Table showing:

-   ID
-   Name
-   Network Type
-   IP Address
-   Operating Status
-   Type

### Get Load Balancer Details

Get detailed information about a load balancer:

```bash
bizfly loadbalancer get <loadbalancer-id>
```

**Example:**

```bash
bizfly loadbalancer get fd554aac-9ab1-11ea-b09d-bbaf82f02f58
```

### Create Load Balancer

Create a new load balancer:

```bash
bizfly loadbalancer create \
  --name <name> \
  [options]
```

**Required Flags:**

-   `--name`: Load balancer name

**Optional Flags:**

-   `--type <type>`: Load balancer type (`small`, `medium`, `large`, `xtralarge`) - default: `medium`
-   `--network-type <type>`: Network type (`external` or `internal`) - default: `external`
-   `--network-id <id>`: VPC network ID (for internal load balancers)
-   `--description <text>`: Description
-   `--listener-name <name>`: Listener name - default: `Default Listener`
-   `--listener-protocol <protocol>`: Listener protocol (`HTTP`, `HTTPS`, `TCP`, `UDP`) - default: `HTTP`
-   `--listener-port <port>`: Listener port - default: `80`
-   `--pool-name <name>`: Pool name - default: `Default`
-   `--algorithm <algorithm>`: Load balancing algorithm (`ROUND_ROBIN`, `LEAST_CONNECTIONS`, `SOURCE_IP`) - default: `ROUND_ROBIN`
-   `--tls-ref <ref>`: TLS certificate reference (for HTTPS)
-   Health monitor options (see Health Monitor section)

**Example:**

```bash
bizfly loadbalancer create \
  --name web-lb \
  --type large \
  --network-type external \
  --listener-protocol HTTP \
  --listener-port 80 \
  --algorithm ROUND_ROBIN
```

### Delete Load Balancer

Delete one or more load balancers:

```bash
bizfly loadbalancer delete <loadbalancer-id> [loadbalancer-id2] ...
```

**Note:** This cascades to delete associated listeners, pools, and health monitors.

**Example:**

```bash
bizfly loadbalancer delete lb-123 lb-456
```

### Resize Load Balancer

Resize a load balancer to a different type:

```bash
bizfly loadbalancer resize <loadbalancer-id> <new-type>
```

**Types:** `small`, `medium`, `large`, `xtralarge`

**Example:**

```bash
bizfly loadbalancer resize lb-123 large
```

## Listener Commands

Listeners define the protocol and port on which the load balancer accepts traffic.

### List Listeners

List all listeners in a load balancer:

```bash
bizfly loadbalancer listener list <loadbalancer-id>
```

**Output:** Table showing:

-   ID
-   Name
-   Protocol
-   Protocol Port
-   Operating Status
-   Default Pool ID

### Get Listener

Get detailed information about a listener:

```bash
bizfly loadbalancer listener get <listener-id>
```

### Create Listener

Create a new listener:

```bash
bizfly loadbalancer listener create <loadbalancer-id> \
  --name <name> \
  --protocol <protocol> \
  --protocol-port <port> \
  [options]
```

**Required Flags:**

-   `--name`: Listener name
-   `--protocol`: Protocol (`HTTP`, `HTTPS`, `TCP`, `UDP`)
-   `--protocol-port`: Port number

**Optional Flags:**

-   `--default-pool-id <id>`: Default pool ID
-   `--description <text>`: Description

**Example:**

```bash
bizfly loadbalancer listener create lb-123 \
  --name https-listener \
  --protocol HTTPS \
  --protocol-port 443 \
  --default-pool-id pool-456
```

### Update Listener

Update a listener:

```bash
bizfly loadbalancer listener update <listener-id> \
  [options]
```

**Optional Flags:**

-   `--name <name>`: New name
-   `--description <text>`: New description
-   `--default-pool-id <id>`: New default pool ID
-   `--tls-ref <ref>`: TLS certificate reference

**Example:**

```bash
bizfly loadbalancer listener update listener-123 \
  --default-pool-id pool-789 \
  --tls-ref cert-456
```

### Delete Listener

Delete a listener:

```bash
bizfly loadbalancer listener delete <listener-id>
```

## Pool Commands

Pools define the group of servers that receive traffic from the load balancer.

### List Pools

List all pools in a load balancer:

```bash
bizfly loadbalancer pool list <loadbalancer-id>
```

**Output:** Table showing:

-   ID
-   Name
-   Algorithm
-   Protocol
-   Operating Status

### Get Pool

Get detailed information about a pool:

```bash
bizfly loadbalancer pool get <pool-id>
```

### Create Pool

Create a new pool:

```bash
bizfly loadbalancer pool create <loadbalancer-id> \
  --name <name> \
  --protocol <protocol> \
  [options]
```

**Required Flags:**

-   `--name`: Pool name
-   `--protocol`: Protocol (`HTTP`, `HTTPS`, `TCP`, `UDP`)

**Optional Flags:**

-   `--algorithm <algorithm>`: Load balancing algorithm (`ROUND_ROBIN`, `LEAST_CONNECTIONS`, `SOURCE_IP`) - default: `ROUND_ROBIN`
-   `--listener-id <id>`: Associated listener ID
-   `--session-persistence-type <type>`: Session persistence type (`HTTP_COOKIE`, `APP_COOKIE`)
-   `--cookie-name <name>`: Cookie name (for `APP_COOKIE`)

**Example:**

```bash
bizfly loadbalancer pool create lb-123 \
  --name web-pool \
  --protocol HTTP \
  --algorithm LEAST_CONNECTIONS \
  --session-persistence-type HTTP_COOKIE
```

### Delete Pool

Delete a pool:

```bash
bizfly loadbalancer pool delete <pool-id>
```

## Health Monitor Commands

Health monitors check the health of pool members and remove unhealthy servers from rotation.

### Get Health Monitor

Get health monitor details:

```bash
bizfly loadbalancer health-monitor get <health-monitor-id>
```

**Output:** Table showing:

-   ID
-   Name
-   Type
-   Delay
-   Max Retries
-   Timeout
-   Operating Status
-   Domain Name
-   URL Path

### Create Health Monitor

Create a health monitor for a pool:

```bash
bizfly loadbalancer health-monitor create <pool-id> \
  --name <name> \
  --type <type> \
  [options]
```

**Required Flags:**

-   `--name`: Health monitor name
-   `--type`: Monitor type (`HTTP`, `HTTPS`, `TCP`, `UDP`)

**Optional Flags:**

-   `--delay <seconds>`: Delay between checks - default: `5`
-   `--timeout <seconds>`: Timeout for each check - default: `5`
-   `--max-retries <count>`: Maximum retries before marking unhealthy - default: `3`
-   `--url-path <path>`: URL path for HTTP/HTTPS checks - default: `/`
-   `--expected-codes <codes>`: Expected HTTP status codes (e.g., `200`, `200-299`)
-   `--method <method>`: HTTP method (`GET`, `POST`, etc.) - default: `GET`
-   `--max-retries-down <count>`: Retries before marking down - default: `3`

**Example:**

```bash
bizfly loadbalancer health-monitor create pool-123 \
  --name web-health-check \
  --type HTTP \
  --delay 10 \
  --timeout 5 \
  --max-retries 3 \
  --url-path /health \
  --expected-codes 200
```

### Update Health Monitor

Update a health monitor:

```bash
bizfly loadbalancer health-monitor update <health-monitor-id> \
  [options]
```

**Optional Flags:** Same as create command

**Example:**

```bash
bizfly loadbalancer health-monitor update hm-123 \
  --delay 15 \
  --max-retries 5
```

### Delete Health Monitor

Delete a health monitor:

```bash
bizfly loadbalancer health-monitor delete <health-monitor-id>
```

## Examples

### Complete Load Balancer Setup

```bash
# Create load balancer
bizfly loadbalancer create \
  --name web-lb \
  --type medium \
  --network-type external

# Create listener
bizfly loadbalancer listener create lb-123 \
  --name http-listener \
  --protocol HTTP \
  --protocol-port 80

# Create pool
bizfly loadbalancer pool create lb-123 \
  --name web-pool \
  --protocol HTTP \
  --algorithm ROUND_ROBIN

# Create health monitor
bizfly loadbalancer health-monitor create pool-456 \
  --name web-health \
  --type HTTP \
  --url-path /health

# Update listener to use pool
bizfly loadbalancer listener update listener-123 \
  --default-pool-id pool-456
```

### HTTPS Load Balancer

```bash
# Create load balancer
bizfly loadbalancer create \
  --name https-lb \
  --type large \
  --network-type external

# Create HTTPS listener with TLS
bizfly loadbalancer listener create lb-123 \
  --name https-listener \
  --protocol HTTPS \
  --protocol-port 443 \
  --tls-ref cert-123

# Create pool
bizfly loadbalancer pool create lb-123 \
  --name https-pool \
  --protocol HTTPS \
  --algorithm LEAST_CONNECTIONS

# Create health monitor
bizfly loadbalancer health-monitor create pool-456 \
  --name https-health \
  --type HTTPS \
  --url-path /health \
  --expected-codes 200
```

### Internal Load Balancer

```bash
# Create internal load balancer
bizfly loadbalancer create \
  --name internal-lb \
  --type medium \
  --network-type internal \
  --network-id vpc-123
```

## Common Use Cases

### Web Application Load Balancer

```bash
bizfly loadbalancer create \
  --name web-app-lb \
  --type large \
  --listener-protocol HTTP \
  --listener-port 80 \
  --algorithm ROUND_ROBIN
```

### High Availability Setup

```bash
# Create multiple load balancers in different zones
bizfly loadbalancer create --name lb-zone-1 --type large
bizfly loadbalancer create --name lb-zone-2 --type large
```

## Troubleshooting

### Load Balancer Creation Fails

-   Verify VPC network ID exists (for internal load balancers)
-   Check available quota
-   Ensure network type is correct

### Health Monitor Not Working

-   Verify URL path is accessible
-   Check expected status codes
-   Ensure servers respond to health checks
-   Verify firewall rules allow health check traffic

### Cannot Delete Load Balancer

-   Remove all listeners first
-   Remove all pools
-   Wait for operations to complete

## Best Practices

1. **Choose appropriate type** - Match load balancer size to traffic
2. **Use health monitors** - Monitor backend server health
3. **Configure appropriate algorithms** - ROUND_ROBIN for general use, LEAST_CONNECTIONS for long connections
4. **Use session persistence** - For stateful applications
5. **Monitor operating status** - Check load balancer health regularly
6. **Use HTTPS** - For secure traffic
7. **Scale appropriately** - Resize load balancer as traffic grows

## Related Commands

-   [Server Management](server.md) - Backend servers
-   [VPC Management](vpc.md) - Network configuration
-   [Firewall Management](firewall.md) - Security rules
