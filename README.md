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

The following examples assume that you have an AWS account.

### init

First thing you'll need to do when you wish to connect to a new cloud is to create a configuration:

```bash
cloudctl -c my_aws init
```

This will create a config file `~/.cloudinit/my_aws.yaml` based on a template with some sensible defaults. More than likely you'll need to change some of these, you can edit the config in vim using the [config-edit](#config-edit) command.

### config-edit

Opens up the cloud configuration file in vim. This is usually `~/.cloudctl/default.yaml` unless you specify a particular cloud with `-c` like so:

```bash
cloudctl -c my_cloud config-edit # opens ~/.cloudctl/my_cloud.yaml
```

You can find out more about the config options for each cloud provider supported in the [config](#Config) section below.

### ls

`cloudctl ls` will list all the machines running in a particular cloud. It also lists valuable details for each machine such as IP addresses.

The term "machine" here varies from cloud to cloud but typically refers to a VM of some sort (e.g. EC2 instances for AWS, droplet for Digital Ocean).

### show

`cloudctl show <machine_id>` prints out any information available to cloudctl about that machine. This is typically a print-out of the provider API object representing that machine and so will vary from cloud provider to cloud provider.

### ssh

Connect to a specific machine over ssh.

Bacause APIs don't provide a way of inspecting the default user account on a machine, you'll need to know this ahead of time. The command should look like this:

```bash
cloudctl ssh <user>@<machine_id>
```

### rm

The `rm` command takes a list of machine ids and terminates them.

**N.B.:** There will be no confirmation dialog and the command returns immediately with no opertunity to cancel: consider yourself warned.

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
