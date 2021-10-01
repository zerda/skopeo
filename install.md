# Installing from packages

## Distribution Packages
`skopeo` may already be packaged in your distribution.

### Fedora

```sh
sudo dnf -y install skopeo
```

### RHEL/CentOS ≥ 8 and CentOS Stream

```sh
sudo dnf -y install skopeo
```

### RHEL/CentOS ≤ 7.x

```sh
sudo yum -y install skopeo
```

### openSUSE

```sh
sudo zypper install skopeo
```

### Alpine

```sh
sudo apk add skopeo
```

### macOS

```sh
brew install skopeo
```

### Nix / NixOS
```sh
$ nix-env -i skopeo
```

### Debian

The skopeo package is available on [Bullseye](https://packages.debian.org/bullseye/skopeo),
and Debian Testing and Unstable.

```bash
# Debian Bullseye, Testing or Unstable/Sid
sudo apt-get update
sudo apt-get -y install skopeo
```

### Raspberry Pi OS arm64 (beta)

Raspberry Pi OS uses the standard Debian's repositories,
so it is fully compatible with Debian's arm64 repository.
You can simply follow the [steps for Debian](#debian) to install Skopeo.


### Ubuntu

The skopeo package is available in the official repositories for Ubuntu 20.10
and newer.

```bash
# Ubuntu 20.10 and newer
sudo apt-get -y update
sudo apt-get -y install skopeo
```

The [Kubic project](https://build.opensuse.org/package/show/devel:kubic:libcontainers:stable/skopeo)
provides packages for Ubuntu 20.04 (it should also work with direct derivatives like Pop!\_OS).

```bash
. /etc/os-release
echo "deb https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${VERSION_ID}/ /" | sudo tee /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list
curl -L https://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${VERSION_ID}/Release.key | sudo apt-key add -
sudo apt-get update
sudo apt-get -y upgrade
sudo apt-get -y install skopeo
```

### Windows
Skopeo has not yet been packaged for Windows. There is an [open feature
request](https://github.com/containers/skopeo/issues/715) and contributions are
always welcome.


Otherwise, read on for building and installing it from source:

To build the `skopeo` binary you need at least Go 1.12.

There are two ways to build skopeo: in a container, or locally without a
container. Choose the one which better matches your needs and environment.

## Building from Source

### Building without a container

Building without a container requires a bit more manual work and setup in your
environment, but it is more flexible:

- It should work in more environments (e.g. for native macOS builds)
- It does not require root privileges (after dependencies are installed)
- It is faster, therefore more convenient for developing `skopeo`.

Install the necessary dependencies:

```bash
# Fedora:
sudo dnf install gpgme-devel libassuan-devel btrfs-progs-devel device-mapper-devel
```

```bash
# Ubuntu (`libbtrfs-dev` requires Ubuntu 18.10 and above):
sudo apt install libgpgme-dev libassuan-dev libbtrfs-dev libdevmapper-dev pkg-config
```

```bash
# macOS:
brew install gpgme
```

```bash
# openSUSE:
sudo zypper install libgpgme-devel device-mapper-devel libbtrfs-devel glib2-devel
```

Make sure to clone this repository in your `GOPATH` - otherwise compilation fails.

```bash
git clone https://github.com/containers/skopeo $GOPATH/src/github.com/containers/skopeo
cd $GOPATH/src/github.com/containers/skopeo && make bin/skopeo
```

By default the `make` command (make all) will build bin/skopeo and the documentation locally.

Building of documentation requires `go-md2man`. On systems that do not have this tool, the
document generation can be skipped by passing `DISABLE_DOCS=1`:
```
DISABLE_DOCS=1 make
```

### Building documentation

To build the manual you will need go-md2man.

```bash
# Debian:
sudo apt-get install go-md2man
```

```
# Fedora:
sudo dnf install go-md2man
```

```
# MacOS:
brew install go-md2man
```

Then

```bash
make docs
```

### Building in a container

Building in a container is simpler, but more restrictive:

- It requires the `podman` command and the ability to run Linux containers.
- The created executable is a Linux executable, and depends on dynamic libraries
  which may only be available only in a container of a similar Linux
  distribution.

```bash
$ make binary
```

### Installation

Finally, after the binary and documentation is built:

```bash
sudo make install
```
