# Hall Monitor

<div align="center" class="no-border">
  <img src="/img/dawg-the-hall-monitor.png" alt="Hall Monitor Placeholder Image">
  <br>
  <h3>Do you have a pass?</h3>
</div>


Hall Monitor is a go/no go utility to help users determine whether a specific cluster addon is ready for upgrade.

## Purpose
A number factors can affect whether the upgrade of an addon (like cert-manager, nginx ingress, etc) will be successful or not. For example, some addon upgrades require a specific api to be available in the cluster, or a specific version of the Kubernetes cluster in general. Or perhaps an addon has deprecated a particular annotation and you want to make sure your upgraded addon doesn't include those deprecated annotations. Rather than having to manually assess each addon, Hall Monitor enables you to create a specification (called a bundle spec) that you can populate with checks for the upgraded version, and run those against your cluster to get an upgrade confidence score.

For example, `cert-manager` [changed a number of annotations](https://cert-manager.io/docs/installation/upgrading/upgrading-0.10-0.11/#additional-annotation-changes) in the upgrade from `0.10` to `0.11`.With Hall Monitor you can add an OPA check to your bundle spec looking for instances of that annotation in the affected cluster resources and be warned about it before you do the upgrade.

## Documentation

# Installation
TBD

# Usage
// Using hall-monitor for the command here as a placeholder
```
hall-monitor --help
The Kubernetes Add-On Upgrade Validation Bundle is a spec that can be used to define and then discover if an add-on upgrade is safe to perform.

Usage:
  hall-monitor [flags]
  hall-monitor [command]

Available Commands:
  check       Check for Helm releases that can be updated
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  version     Prints the current version of the tool.

Flags:
  -h, --help      help for hall-monitor
  -v, --v Level   number for the log level verbosity

Use "hall-monitor [command] --help" for more information about a command.
```

Pass in a bundle spec with the addon definitions that you want to check
```
hall-monitor check /path/to/bundle/bundle.yaml
```



<!-- Begin boilerplate -->
## Join the Fairwinds Open Source Community

The goal of the Fairwinds Community is to exchange ideas, influence the open source roadmap,
and network with fellow Kubernetes users.
[Chat with us on Slack](https://join.slack.com/t/fairwindscommunity/shared_invite/zt-e3c6vj4l-3lIH6dvKqzWII5fSSFDi1g)
or
[join the user group](https://www.fairwinds.com/open-source-software-user-group) to get involved!

<a href="https://www.fairwinds.com/t-shirt-offer?utm_source=pluto&utm_medium=pluto&utm_campaign=pluto-tshirt">
  <img src="https://www.fairwinds.com/hubfs/Doc_Banners/Fairwinds_OSS_User_Group_740x125_v6.png" alt="Love Fairwinds Open Source? Share your business email and job title and we'll send you>
</a>

## Other Projects from Fairwinds

Enjoying Pluto? Check out some of our other projects:
* [Polaris](https://github.com/FairwindsOps/Polaris) - Audit, enforce, and build policies for Kubernetes resources, including over 20 built-in checks for best practices
* [Goldilocks](https://github.com/FairwindsOps/Goldilocks) - Right-size your Kubernetes Deployments by compare your memory and CPU settings against actual usage
* [Nova](https://github.com/FairwindsOps/Nova) - Check to see if any of your Helm charts have updates available
* [rbac-manager](https://github.com/FairwindsOps/rbac-manager) - Simplify the management of RBAC in your Kubernetes clusters

Or [check out the full list](https://www.fairwinds.com/open-source-software?utm_source=pluto&utm_medium=pluto&utm_campaign=pluto)
## Fairwinds Insights
If you're interested in running Pluto in multiple clusters,
tracking the results over time, integrating with Slack, Datadog, and Jira,
or unlocking other functionality, check out
[Fairwinds Insights](https://www.fairwinds.com/pluto-user-insights-demo?utm_source=pluto&utm_medium=pluto&utm_campaign=pluto),
a platform for auditing and enforcing policy in Kubernetes clusters.

<a href="https://www.fairwinds.com/pluto-user-insights-demo?utm_source=pluto&utm_medium=ad&utm_campaign=plutoad">
  <img src="https://www.fairwinds.com/hubfs/Doc_Banners/Fairwinds_Pluto_Ad.png" alt="Fairwinds Insights" />
</a>


