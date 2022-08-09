# 1PENRC(1)

## NAME

1penrc - Export 1Password secrets as environment variables.

## SYNOPSIS

`1penrc [--config CONFIG_FILE] ENVIRONMENT`

## DESCRIPTION

1penrc is a command line utility for exporting secrets from 1Password as environment variables.
This is done by declaring which secrets belong to an environment and how these are exported in a configuration file.
The output from 1penrc can than be used by shell builtins like `eval` or `source` similar to `op signin`.

## OPTIONS

`--help`

:   Show the help message and exit

`--config CONFIG_FILE`

:   Use alternate configuration file.

## FILES

`CONFIG_FILE`

:   Configuration file defining each environment. By default it's located at `1penrc/config.yml` inside the user's configuration directory. See [**1penrc(5)**](./1penrc.5.md) for more details.

## EXAMPLES

The following shows how to export variables of an environment named `example.com` declared in an alternate configuration file.

```bash
source <(1penrc --config ./config.yaml example.com)
```

## BUGS

Bugs are tracked and can be reported at https://github.com/xinau/1penrc/issues.

## SEE ALSO

[**1penrc(5)**](./1penrc.5.md)
