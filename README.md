# bizflyctl

Command line to interact with Bizfly Cloud

### Install `bizflyctl`

There are three ways to install the bizflyctl.

### Using homebrew on macOS:
```shell script
brew tap bizflycloud/bizflyctl
brew install bizflyctl
```

#### Build from source code

- Clone the source code

```shell script
git clone https://github.com/bizflycloud/bizflyctl
```

- Install golang

- Build

```shell script
go build -o bizfly main.go
```

#### Download the latest release in Github release pages

- Navigate to [release page](https://github.com/bizflycloud/bizflyctl/releases). Download the tar.gz file with your platform (Linux, Windows and MacOS).

- Extract the tar.gz file


#### Install and configure
## Using config file
- Copy bizfly binary to `/usr/local/bin` with Linux and MacOS

- Configure user name and password. Create a file `.bizfly.yaml` in your `$HOME` directory

```
email: <your email>
password: <your password>
```

## Using environment variables
- Configure these environment for using bizflyctl:
  - `BIZFLY_CLOUD_EMAIL` (Required)
  - `BIZFLY_CLOUD_PASSWORD` (Required)
  - `BIZFLY_CLOUD_REGION` (Optional. Default value is HN)
  - `BIZFLY_CLOUD_PROJECT_NAME` (Optional)

### Example

```shell script
➜  bizflycli git:(master) ✗ bizfly --help
Bizfly Cloud Command Line

Usage:
  bizfly [command]

Available Commands:
  help        Help about any command
  server      Bizfly Cloud Server Interaction
  snapshot    Bizfly Cloud Snapshot Interaction
  volume      Bizfly Cloud Volume Interaction

Flags:
      --config string     config file (default is $HOME/.bizfly.yaml)
      --email string      Your Bizfly Cloud Email
  -h, --help              help for bizfly
      --password string   Your Bizfly Cloud Password
  -t, --toggle            Help message for toggle

Use "bizfly [command] --help" for more information about a command.

```

- Example Get snapshot

```shell script
➜  bizflycli git:(master) ✗ bizfly snapshot get 5af19947-566d-48a1-bc45-93666086951f
+--------------------------------------+----------------------------+-----------+------+--------------------------------------+----------------------------+--------------------------------------+
|                  ID                  |            NAME            |  STATUS   | SIZE |                 TYPE                 |         CREATED AT         |              VOLUME ID               |
+--------------------------------------+----------------------------+-----------+------+--------------------------------------+----------------------------+--------------------------------------+
| 5af19947-566d-48a1-bc45-93666086951f | snapshot-15-38-41-4-5-2019 | available |   20 | ec6fb900-1ae0-4e9e-90e0-53a6063f95e1 | 2019-05-04T06:38:48.000000 | 172b31b6-1d2c-4421-9e89-c74e28d0d77d |
+--------------------------------------+----------------------------+-----------+------+--------------------------------------+----------------------------+--------------------------------------+

```
