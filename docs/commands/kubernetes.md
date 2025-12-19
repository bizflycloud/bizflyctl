# Kubernetes Management Commands

Manage Kubernetes clusters, worker pools, and nodes on Bizfly Cloud.

## Overview

Bizfly Cloud provides managed Kubernetes clusters. You can create clusters, manage worker pools, and configure nodes using these commands.

## Cluster Commands

### List Clusters

List all Kubernetes clusters:

```bash
bizfly kubernetes list
```

**Output:** Table showing:

-   ID (Cluster UID)
-   Name
-   VPC Network ID
-   Worker Pools Count
-   Cluster Status
-   Tags
-   Created At
-   Cluster Version

### Get Cluster Details

Get detailed information about a cluster:

```bash
bizfly kubernetes get <cluster-id>
```

**Example:**

```bash
bizfly kubernetes get xfbxsws38dcs8o94
```

**Output:** Includes worker pool IDs and detailed cluster information.

### Create Cluster

Create a new Kubernetes cluster with worker pools.

**Using flags:**

```bash
bizfly kubernetes create \
  --name <cluster-name> \
  --version <version-id> \
  --vpc-network-id <vpc-id> \
  --worker-pool "<pool-config>" \
  [options]
```

**Using config file:**

```bash
bizfly kubernetes create --config-file <config-file.yml>
```

**Required Flags:**

-   `--name`: Cluster name
-   `--version`: Kubernetes version ID
-   `--vpc-network-id`: VPC network ID
-   `--worker-pool`: Worker pool configuration (can specify multiple)

**Optional Flags:**

-   `--tag <tag>`: Tags (can specify multiple)
-   `--config-file`: Path to YAML configuration file

**Worker Pool Configuration Format:**

```
name=<name>;flavor=<flavor>;profile_type=<type>;volume_type=<type>;volume_size=<size>;availability_zone=<zone>;desired_size=<size>;min_size=<size>;max_size=<size>;labels=<key=value>;taints=<key=value:effect>
```

**Example with flags:**

```bash
bizfly kubernetes create \
  --name production-cluster \
  --version 5f7d3a91d857155ad4993a32 \
  --vpc-network-id 145bed1f-a7f7-4f88-ab3d-ce2fc95a4e71 \
  --tag production \
  --tag env=prod \
  --worker-pool "name=worker-pool-1;flavor=nix.3c_6g;profile_type=premium;volume_type=PREMIUM-HDD1;volume_size=40;availability_zone=HN1;desired_size=2;min_size=1;max_size=10;labels=env=prod;taints=app=demo:NoSchedule"
```

**Example config file** (`create_cluster.yml`):

```yaml
name: production-cluster
version: 5f7d3a91d857155ad4993a32
vpc_network_id: 145bed1f-a7f7-4f88-ab3d-ce2fc95a4e71
tags:
    - production
    - env=prod
worker_pools:
    - name: worker-pool-1
      flavor: nix.3c_6g
      profile_type: premium
      volume_type: PREMIUM-HDD1
      volume_size: 40
      availability_zone: HN1
      desired_size: 2
      min_size: 1
      max_size: 10
      enable_autoscaling: true
      labels:
          env: prod
      taints:
          - key: app
            value: demo
            effect: NoSchedule
```

### Delete Cluster

Delete a Kubernetes cluster and all its worker pools:

```bash
bizfly kubernetes delete <cluster-id>
```

**Warning:** This will delete the cluster and all worker pools.

**Example:**

```bash
bizfly kubernetes delete xfbxsws38dcs8o94
```

## Worker Pool Commands

### List Worker Pools

List worker pools in a cluster:

```bash
bizfly kubernetes workerpool list <cluster-id>
```

### Get Worker Pool

Get detailed information about a worker pool:

```bash
bizfly kubernetes workerpool get <cluster-id> <worker-pool-id>
```

**Example:**

```bash
bizfly kubernetes workerpool get cluster-123 pool-456
```

**Output:** Table showing:

-   ID
-   Name
-   Version
-   Flavor
-   Volume Size
-   Volume Type
-   Nodes
-   Enabled AutoScaling
-   Min Size
-   Max Size
-   Created At

### Add Worker Pool

Add a new worker pool to a cluster.

**Using flags:**

```bash
bizfly kubernetes workerpool add <cluster-id> \
  --worker-pool "<pool-config>"
```

**Using config file:**

```bash
bizfly kubernetes workerpool add <cluster-id> \
  --config-file <config-file.yml>
```

**Example:**

```bash
bizfly kubernetes workerpool add cluster-123 \
  --worker-pool "name=worker-pool-2;flavor=nix.6c_12g;profile_type=premium;volume_type=PREMIUM-HDD1;volume_size=80;availability_zone=HN1;desired_size=3;min_size=2;max_size=20;labels=env=prod;taints=app=demo:NoSchedule"
```

**Example config file** (`add_pools.yml`):

```yaml
worker_pools:
    - name: worker-pool-2
      flavor: nix.6c_12g
      profile_type: premium
      volume_type: PREMIUM-HDD1
      volume_size: 80
      availability_zone: HN1
      desired_size: 3
      min_size: 2
      max_size: 20
      enable_autoscaling: true
      labels:
          env: prod
      taints:
          - key: app
            value: demo
            effect: NoSchedule
```

### Update Worker Pool

Update worker pool configuration:

```bash
bizfly kubernetes workerpool update <cluster-id> <worker-pool-id> \
  --desired-size <size> \
  --min-size <size> \
  --max-size <size> \
  --autoscaling <true|false>
```

**Required Flags:**

-   `--desired-size`: Desired number of nodes
-   `--min-size`: Minimum number of nodes
-   `--max-size`: Maximum number of nodes

**Optional Flags:**

-   `--autoscaling`: Enable or disable autoscaling

**Example:**

```bash
bizfly kubernetes workerpool update cluster-123 pool-456 \
  --desired-size 5 \
  --min-size 3 \
  --max-size 15 \
  --autoscaling true
```

### Delete Worker Pool

Delete a worker pool from a cluster:

```bash
bizfly kubernetes workerpool delete <cluster-id> <worker-pool-id>
```

**Example:**

```bash
bizfly kubernetes workerpool delete cluster-123 pool-456
```

## Node Commands

### Recycle Node

Recycle (replace) a node in a worker pool:

```bash
bizfly kubernetes workerpool node recycle <cluster-id> <worker-pool-id> <node-id>
```

**Example:**

```bash
bizfly kubernetes workerpool node recycle cluster-123 pool-456 node-789
```

### Delete Node

Delete a node from a worker pool:

```bash
bizfly kubernetes workerpool node delete <cluster-id> <worker-pool-id> <node-id>
```

**Example:**

```bash
bizfly kubernetes workerpool node delete cluster-123 pool-456 node-789
```

## Kubeconfig Commands

### Get Kubeconfig

Download the kubeconfig file for a cluster:

```bash
bizfly kubernetes kubeconfig get <cluster-id> [options]
```

**Optional Flags:**

-   `--output <path>`: Output path (file or directory) - default: current directory
-   `--expire-time <seconds>`: Kubeconfig expiration time - default: `3000`

**Examples:**

Save to current directory:

```bash
bizfly kubernetes kubeconfig get cluster-123
# Creates: cluster-123.kubeconfig
```

Save to specific file:

```bash
bizfly kubernetes kubeconfig get cluster-123 \
  --output /path/to/kubeconfig.yaml
```

Save to directory:

```bash
bizfly kubernetes kubeconfig get cluster-123 \
  --output /path/to/configs/
# Creates: /path/to/configs/cluster-123.kubeconfig
```

With custom expiration:

```bash
bizfly kubernetes kubeconfig get cluster-123 \
  --expire-time 7200
```

## Examples

### Complete Cluster Lifecycle

```bash
# Create cluster
bizfly kubernetes create \
  --name my-cluster \
  --version <version-id> \
  --vpc-network-id <vpc-id> \
  --worker-pool "name=pool1;flavor=nix.3c_6g;profile_type=premium;volume_type=PREMIUM-HDD1;volume_size=40;availability_zone=HN1;desired_size=2;min_size=1;max_size=10"

# List clusters
bizfly kubernetes list

# Get cluster details
bizfly kubernetes get cluster-123

# Get kubeconfig
bizfly kubernetes kubeconfig get cluster-123

# Add worker pool
bizfly kubernetes workerpool add cluster-123 \
  --worker-pool "name=pool2;flavor=nix.6c_12g;profile_type=premium;volume_type=PREMIUM-HDD1;volume_size=80;availability_zone=HN1;desired_size=3;min_size=2;max_size=20"

# Update worker pool
bizfly kubernetes workerpool update cluster-123 pool-456 \
  --desired-size 5 \
  --min-size 3 \
  --max-size 15

# Delete cluster
bizfly kubernetes delete cluster-123
```

### Using Configuration Files

**create_cluster.yml:**

```yaml
name: production-cluster
version: 5f7d3a91d857155ad4993a32
vpc_network_id: 145bed1f-a7f7-4f88-ab3d-ce2fc95a4e71
tags:
    - production
worker_pools:
    - name: worker-pool-1
      flavor: nix.3c_6g
      profile_type: premium
      volume_type: PREMIUM-HDD1
      volume_size: 40
      availability_zone: HN1
      desired_size: 2
      min_size: 1
      max_size: 10
      enable_autoscaling: true
```

**Create cluster:**

```bash
bizfly kubernetes create --config-file create_cluster.yml
```

## Worker Pool Configuration Parameters

### Required Parameters

-   `name`: Worker pool name
-   `flavor`: Server flavor (e.g., `nix.3c_6g`)
-   `profile_type`: Profile type (`premium`, `basic`, `enterprise`)
-   `volume_type`: Volume type (e.g., `PREMIUM-HDD1`)
-   `volume_size`: Volume size in GB
-   `availability_zone`: Availability zone (e.g., `HN1`)
-   `desired_size`: Desired number of nodes
-   `min_size`: Minimum number of nodes
-   `max_size`: Maximum number of nodes

### Optional Parameters

-   `enable_autoscaling`: Enable autoscaling (true/false)
-   `labels`: Node labels (key=value pairs, comma-separated)
-   `taints`: Node taints (key=value:effect, comma-separated)
    -   Effect values: `NoSchedule`, `PreferNoSchedule`, `NoExecute`

## Common Use Cases

### Production Cluster

```bash
bizfly kubernetes create \
  --name prod-cluster \
  --version <version-id> \
  --vpc-network-id <vpc-id> \
  --tag production \
  --worker-pool "name=prod-workers;flavor=nix.6c_12g;profile_type=enterprise;volume_type=PREMIUM-HDD1;volume_size=100;availability_zone=HN1;desired_size=5;min_size=3;max_size=20;enable_autoscaling=true"
```

### Development Cluster

```bash
bizfly kubernetes create \
  --name dev-cluster \
  --version <version-id> \
  --vpc-network-id <vpc-id> \
  --tag development \
  --worker-pool "name=dev-workers;flavor=nix.2c_4g;profile_type=basic;volume_type=PREMIUM-HDD1;volume_size=40;availability_zone=HN1;desired_size=1;min_size=1;max_size=5"
```

## Troubleshooting

### Cluster Creation Fails

-   Verify VPC network ID exists: `bizfly vpc list`
-   Check Kubernetes version ID is valid
-   Ensure flavor is available in the region
-   Verify availability zone is correct

### Worker Pool Creation Fails

-   Check cluster status (must be active)
-   Verify flavor is available
-   Ensure volume type exists in availability zone
-   Check autoscaling limits

### Cannot Get Kubeconfig

-   Verify cluster is in active state
-   Check cluster ID is correct
-   Ensure you have proper permissions

## Best Practices

1. **Use configuration files** for complex cluster setups
2. **Enable autoscaling** for production workloads
3. **Use appropriate flavors** based on workload requirements
4. **Set proper min/max sizes** for autoscaling
5. **Use labels and taints** for node scheduling
6. **Regular updates** - Keep clusters updated
7. **Backup configurations** - Save cluster config files

## Related Commands

-   [VPC Management](vpc.md) - Create VPCs for clusters
-   [Server Management](server.md) - Understand flavors and zones
-   [Volume Management](volume.md) - Understand volume types
