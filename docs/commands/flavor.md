# Flavor Management Commands

List available server flavors (instance types) with their specifications.

## Overview

Flavors define the compute resources (CPU, RAM) available for servers. Each flavor has a specific combination of vCPUs and RAM.

## Commands

### List Flavors

List all available flavors:

```bash
bizfly flavor list [options]
```

**Optional Flags:**

-   `--vcpus <count>`: Filter by number of vCPUs
-   `--ram <gb>`: Filter by RAM in GB
-   `--category <category>`: Filter by category

**Example:**

```bash
bizfly flavor list
bizfly flavor list --vcpus 4
bizfly flavor list --category premium
```

**Output:** Table showing:

-   ID
-   Name
-   CPU (vCPUs)
-   RAM (GB)
-   Category

## Examples

### Find Appropriate Flavor

```bash
# List all flavors
bizfly flavor list

# Find flavors with 4 vCPUs
bizfly flavor list --vcpus 4

# Find flavors with 8GB RAM
bizfly flavor list --ram 8

# Find premium category flavors
bizfly flavor list --category premium
```

### Using Flavor in Server Creation

```bash
# List available flavors
bizfly flavor list

# Create server with specific flavor
bizfly server create \
  --name my-server \
  --flavor nix.3c_6g \
  --image-id <image-id> \
  --rootdisk-size 40
```

## Common Flavor Categories

-   **Basic**: Entry-level instances
-   **Premium**: Standard instances
-   **Enterprise**: High-performance instances

## Flavor Naming Convention

Flavors typically follow the pattern: `nix.Xc_Yg`

-   `X`: Number of vCPUs
-   `Y`: RAM in GB

**Examples:**

-   `nix.2c_4g`: 2 vCPUs, 4GB RAM
-   `nix.3c_6g`: 3 vCPUs, 6GB RAM
-   `nix.6c_12g`: 6 vCPUs, 12GB RAM

## Common Use Cases

### Finding Right Flavor for Workload

```bash
# For small application
bizfly flavor list --vcpus 2 --ram 4

# For medium application
bizfly flavor list --vcpus 4 --ram 8

# For large application
bizfly flavor list --vcpus 8 --ram 16
```

### Development vs Production

```bash
# Development (smaller)
bizfly server create --flavor nix.2c_4g ...

# Production (larger)
bizfly server create --flavor nix.6c_12g ...
```

## Related Commands

-   [Server Management](server.md) - Use flavors when creating servers
-   [Kubernetes Management](kubernetes.md) - Use flavors for worker pools
