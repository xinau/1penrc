# 1PENRC

1penrc - Export 1Password secrets as environment variables.

**WARNING**: This software is still in it's infancy.
Bugs and breaking changes might occur.
Use at your own risk.

## INSTALLATION

1penrc uses the [1Password CLI](https://support.1password.com/command-line-getting-started/) (>=2.0.0) to interact with 1Password,
which therefore needs to be installed
and has to be available inside the users _$PATH_.

Currently, no pre-build binaries are provided.
To compile and install 1penrc manually,
Git and a recent version of the Go toolchain (>=1.18) is required.

```bash
git clone https://github.com/xinau/1penrc.git
cd 1penrc
go build -o 1penrc ./cmd/1penrc
sudo mv 1penrc /usr/local/bin
```

## USAGE

To export 1Password secrets as environment variables a configuration file defining an environment is needed.
The following example uses a configuration file which defines an environment of the name `example.com`.
This environment exposes a variable `API_TOKEN` with the value of the `token` secret inside the `API` section of the `Example` item inside `Private` vault.

```yaml
environments:
  secret_configs:
    - name: API_TOKEN
      secret: op://Private/Example/API/token
```

The environment can be sourced by referring to the environment by it's name. I.e.

```bash
source <(1penrc --config ./config.yaml example.com)
```

For more information take a look at the manpages link:[1penrc(1)](./docs/man/1penrc.1.md) and [1penrc(5)](./docs/man/1penrc.5.md).

## DISCLAIMER

As noted above this software is still in its early stages and breaking changes to behaviour might occur.
This software is build to ease the process of exporting 1Password item fields as environment
variables to be consumed by other tools on the CLI and as of now nothing else.

## LICENSE

This project is under [MIT license](./LICENSE).
