# Image Management Commands

List and manage OS images available for creating servers.

## Overview

Images are operating system templates used to create servers. Bizfly Cloud provides various OS images including Ubuntu, CentOS, Debian, and others.

## Commands

### List Images

List all available OS images:

```bash
bizfly image list
```

**Output:** Table showing:

-   ID (Image ID)
-   Distribution (OS distribution name)
-   Version (OS version)

**Example:**

```bash
bizfly image list
```

**Output Example:**

```
+--------------------------------------+--------------+------------------+
|                  ID                  | Distribution |     Version      |
+--------------------------------------+--------------+------------------+
| abc123-def456-ghi789                 | Ubuntu       | 20.04 LTS       |
| xyz789-abc123-def456                 | CentOS       | 8               |
+--------------------------------------+--------------+------------------+
```

## Usage Examples

### Find Image for Server Creation

```bash
# List available images
bizfly image list

# Use image ID when creating server
bizfly server create \
  --name my-server \
  --flavor nix.3c_6g \
  --image-id abc123-def456-ghi789 \
  --rootdisk-size 40
```

## Common Use Cases

### Finding Ubuntu Image

```bash
bizfly image list | grep Ubuntu
```

### Finding Latest OS Version

```bash
# List all images and filter
bizfly image list
```

## Related Commands

-   [Server Management](server.md) - Create servers from images
-   [Custom Image](custom-image.md) - Create custom images
