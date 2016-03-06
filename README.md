# cloudctl

`cloudctl` is a command line interface to various public and private clouds.

It is currently in very early development, all suggestions are welcome! If you've got any ideas please open an issue on our [Github page](https://github.com/nicr9/cloudctl/issues/new).

## Installation

If you have already set up Go and your $GOPATH then you just need to install like so:

```bash
go get github.com/nicr9/cloudctl
go install github.com/nicr9/cloudctl
```

## Usage

The current trend embraced by the enterprise is the cloud. This comes with a lot of benifits but it's early days yet so there are a lot of players solving the same problems. You will likely find yourself working on many teams/projects that each have their own style of cloud deployment (or a team that deploys to multiple providers for redundancy). There's no reason you shouldn't have a tool that works seamlessly across them all.

`cloudctl` was designed to simplify interfacing with multiple cloud accounts. Configure a wide array of options for each cloud and specify your target at runtime with the `--cloud` flag. I recommend that you create an alias for each target:

```bash
alias work="cloudctl --cloud work"
alias play="cloudctl --cloud play"
```

## Commands

If you're just starting out, you may have an AWS account. Let's find out how to connect to that.

First you'll need to initialise a new cloud configuration:

```bash
cloudctl -c my_aws init
```

This will create a config file `~/.cloudinit/my_aws.yaml` and will add a basic config from a template. You can find out more about the config options for each cloud provider supported in the [config](#Config) section below.

## Config

At the moment the only top level configuration option is `provider` which defaults to `aws`.

| Option | Default value | Supported values |
| --- | --- | --- |
| provider | aws | aws |

### AWS

To interact with an AWS stack, you'll need to provide access and secret keys. Your options are either to provide these as config options or environment variables.

| Option | Default value | Environment variable |
| --- | --- | --- |
| access_key | "" | AWS_ACCESS_KEY_ID |
| secret_key | "" | AWS_SECRET_ACCESS_KEY |
