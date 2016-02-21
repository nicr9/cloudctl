# cloudctl

`cloudctl` is a command line interface to various public and private clouds.

It is currently in very early development so there is almost nothin to document here. Come back soon!

## Installation

If you have already set up Go and your $GOPATH then you just need to install like so:

```bash
go install github.com/nicr9/cloudctl
```

## Usage

`cloudctl` requires you to create a config file for each cloud you interact with.

If you're just starting out, you may have an AWS account. Let's find out how to connect to that.

First you'll need to initialise a new cloud configuration:

```bash
cloudctl -c my_aws init
```

This will create a config file `~/.cloudinit/my_aws.yaml` and will add a basic config from a template. You can find out more about the config options for each cloud provider supported in the [config](#Config) section below.

## Config

At the moment the only top level configuration option is `platform` which defaults to `aws`.

| Option | Default value | Supported values |
| --- | --- | --- |
| platform | aws | aws |

### AWS

To interact with an AWS stack, you'll need to provide access and secret keys. Your options are either to provide these as config options or environment variables.

| Option | Default value | Environment variable |
| --- | --- | --- |
| access_key | "" | AWS_ACCESS_KEY_ID |
| secret_key | "" | AWS_SECRET_ACCESS_KEY |
