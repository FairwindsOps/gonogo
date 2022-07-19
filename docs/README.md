<div align="center" class="no-border">
  <img src="/img/gonogo-logo.png" alt="GoNoGo Logo">
  <br>
  <h3>Determine the upgrade confidence of your cluster addons.</h3>
  <a href="https://github.com/FairwindsOps/gonogo/releases">
    <img src="https://img.shields.io/github/v/release/FairwindsOps/gonogo">
  </a>
  <a href="https://goreportcard.com/report/github.com/FairwindsOps/gonogo">
    <img src="https://goreportcard.com/badge/github.com/FairwindsOps/gonogo">
  </a>
  <a href="https://circleci.com/gh/FairwindsOps/gonogo.svg">
    <img src="https://circleci.com/gh/FairwindsOps/gonogo.svg?style=svg">
  </a>
  <a href="https://insights.fairwinds.com/gh/FairwindsOps/gonogo">
    <img src="https://insights.fairwinds.com/v0/gh/FairwindsOps/gonogo/badge.svg">
  </a>
</div>


GoNoGo is a utility to help users determine upgrade confidence around Kubernetes cluster addons.

## Purpose
A number factors can affect whether the upgrade of an addon (like cert-manager, nginx ingress, etc) will be successful or not. For example, some addon upgrades require a specific api to be available in the cluster, or a specific version of the Kubernetes cluster in general. Or perhaps an addon has deprecated a particular annotation and you want to make sure your upgraded addon doesn't include those deprecated annotations. Rather than having to manually assess each addon, GoNoGo enables you to create a specification (called a bundle spec) that you can populate with checks for the upgraded version, and run those against your cluster to get an upgrade confidence score.

For example, `cert-manager` [changed a number of annotations](https://cert-manager.io/docs/installation/upgrading/upgrading-0.10-0.11/#additional-annotation-changes) in the upgrade from `0.10` to `0.11`.With Hall Monitor you can add an OPA check to your bundle spec looking for instances of that annotation in the affected cluster resources and be warned about it before you do the upgrade.

## Documentation

# Installation
TBD

# Usage
```
gonogo --help
The Kubernetes Add-On Upgrade Validation Bundle is a spec that can be used to define and then discover if an add-on upgrade is safe to perform.

Usage:
  gonogo [flags]
  gonogo [command]

Available Commands:
  check       Check for Helm releases that can be updated
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Prints the current version of the tool.

Flags:
  -h, --help      help for gonogo
  -v, --v Level   number for the log level verbosity

Use "gonogo [command] --help" for more information about a command.
```

Pass in a bundle spec with the addon definitions that you want to check
```
gonogo check /path/to/bundle/bundle.yaml
```



<!-- Begin boilerplate -->
## Join the Fairwinds Open Source Community

The goal of the Fairwinds Community is to exchange ideas, influence the open source roadmap,
and network with fellow Kubernetes users.
[Chat with us on Slack](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g)
or
[join the user group](https://www.fairwinds.com/open-source-software-user-group) to get involved!

<a href="https://www.fairwinds.com/t-shirt-offer?utm_source=gonogo&utm_medium=gonogo&utm_campaign=gonogo-tshirt">
 <img src="https://www.fairwinds.com/hubfs/Doc_Banners/Fairwinds_OSS_User_Group_740x125_v6.png" alt="Love Fairwinds Open Source? Share your business email and job title and we'll send you a free Fairwinds t-shirt!" />
</a>

## Other Projects from Fairwinds

Enjoying GoNoGo? Check out some of our other projects:
* [Polaris](https://github.com/FairwindsOps/Polaris) - Audit, enforce, and build policies for Kubernetes resources, including over 20 built-in checks for best practices
* [Goldilocks](https://github.com/FairwindsOps/Goldilocks) - Right-size your Kubernetes Deployments by compare your memory and CPU settings against actual usage
* [Nova](https://github.com/FairwindsOps/Nova) - Check to see if any of your Helm charts have updates available
* [rbac-manager](https://github.com/FairwindsOps/rbac-manager) - Simplify the management of RBAC in your Kubernetes clusters

Or [check out the full list](https://www.fairwinds.com/open-source-software?utm_source=gonogo&utm_medium=gonogo&utm_campaign=gonogo)



