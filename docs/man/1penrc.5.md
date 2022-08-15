# 1PENRC(5)

## NAME

`1penrc/config.yaml` - 1penrc configuration file.

## DESCRIPTION

1penrc uses a configuration file for declaring which secrets belong to an environment and how these are exported.
By default this file is located at `1penrc/config.yaml` inside the user's configuration directory.

## CONFIGURATION

The configuration file is written in YAML format, defined by the schema below.
Brackets indicate that a parameter is optional.
For non-list parameters the value is set to the specified default.

```yaml
# 1Password client configuration.
client:
  # Name of the 1Password CLI executable.
  [ executable: <string> | default "op" ]

# List of environment configurations.
environments:
  [ - <environment_config> ... ]
```

### <environment_config>

An environment block specifies which variables belong to an environment
and how these are being obtained.

```yaml
# Environment name used for referencing.
name: <string>

# List of AWS STS provider configurations.
aws_sts_configs:
  [ - <aws_sts_configs> ... ]

# List of secret provider configurations.
secret_configs: 
  [ - <secret_config> ... ]

# List of value provider configurations.
value_configs: 
  [ - <value_config> ... ]
```

### <aws_sts_config>

The AWS STS provider block specifies an AWS role to assume
and 1Password secrets to use for authentication.

```yaml
# Name of the environment variable.
name: <string>

# Reference to 1Password secret used as AWS access key id.
access_key_id: <string>

# Reference to 1Password secret used as AWS secret access key.
secret_access_key: <string>

# ARN of the role to assume.
role_arn: <string>

# Time-To-Live duration of the AWS session. 
ttl: <duration | "1h" >

# Reference to 1Password secret used as multi-factor authentication token. Only
# applied if not empty.
mfa_token: <secret | "" >

# Serial number of the virtual multi-factor authentication device. If left
# empty the serial number of the first device assigned to the current user
# will be tried.
mfa_serial_number: <string | "" >
```

### <secret_config>

The secret provider block specifies the name
and 1Password secret to use as value of a variable.

```yaml
# Name of the environment variable.
name: <string>

# Reference to 1Password secret used as value.
secret: <secret>
```

### <value_config>

The value provider block specifies the name and value of a variable.

```yaml
# Name of the environment variable.
name: <string>

# Value of the environment variable.
value: <string>
```

## SEE ALSO

[**1penrc(1)**](./1penrc.1.md)
