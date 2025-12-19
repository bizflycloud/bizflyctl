# Bizfly Cloud CLI (bizflyctl) Documentation

Welcome to the Bizfly Cloud CLI documentation! This guide will help you install, configure, and use `bizflyctl` to manage your Bizfly Cloud resources from the command line.

## Table of Contents

1. [Installation](installation.md)
2. [Authentication](authentication.md)
3. [Configuration](configuration.md)
4. [Command Reference](#command-reference)
    - [Login](commands/login.md)
    - [Server Management](commands/server.md)
    - [Volume Management](commands/volume.md)
    - [Snapshot Management](commands/snapshot.md)
    - [VPC Management](commands/vpc.md)
    - [Network Interface Management](commands/network-interface.md)
    - [Firewall Management](commands/firewall.md)
    - [Load Balancer Management](commands/loadbalancer.md)
    - [Kubernetes Management](commands/kubernetes.md)
    - [Image Management](commands/image.md)
    - [SSH Key Management](commands/sshkey.md)
    - [Flavor Management](commands/flavor.md)
    - [DNS Management](commands/dns.md)
    - [WAN IP Management](commands/wan-ip.md)
    - [Container Registry](commands/container-registry.md)
    - [Custom Image](commands/custom-image.md)
    - [IAM](commands/iam.md)
    - [Scheduled Volume Backup](commands/scheduled-volume-backup.md)
    - [Alert Management](commands/alert.md)

## Quick Start

1. **Install bizflyctl** - See [Installation Guide](installation.md)
2. **Authenticate** - Run `bizfly login` or configure credentials (see [Authentication](authentication.md))
3. **Start Managing Resources** - Use commands like `bizfly server list` to see your servers

## Getting Help

-   Use `bizfly --help` to see all available commands
-   Use `bizfly <command> --help` to see help for a specific command
-   Check the [Command Reference](#command-reference) section for detailed documentation

## Examples

### List all servers

```bash
bizfly server list
```

### Create a new server

```bash
bizfly server create --name my-server --flavor nix.3c_6g --image-id <image-id> --rootdisk-size 40
```

### List all volumes

```bash
bizfly volume list
```

## Support

For issues, questions, or contributions, please visit the [GitHub repository](https://github.com/bizflycloud/bizflyctl).
