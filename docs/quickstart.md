---
meta:
  - name: description
    content: "Quickstart instructions to install GoNoGo. GoNoGo depends on a bundle spec to run its validation. See documenation on creating the bundle spec for more details "
---
# QuickStart

First, follow the install instructions to install gonogo.

## Bundle Creation

GoNoGo depends on a bundle spec to run its validation. See documentation on creating the bundle spec for more details.

## Running Against Bundles

GoNoGo can be run against one more more bundles files by either using the `-b` flag to pass in separate files
```
gonogo check -b /path/to/file/1 -b /path/to/file/2
```

or the `-d` flag to specify a directory containing bundles
```
gonogo check -d /path/to/dir
```

You can also run GoNoGo with no flags and it will use the curated bundle files found in the `pkg/bundle/bundles` directory of this repo.

In all cases the resulting output should be a json document with a list of found cluster addons as specified in your bundle file. For each cluster addon in the list, you should see the output of the fields you defined in your spec. For example:

```
{
 "Addons": [
  {
   "Name": "cert-manager",
   "Versions": {
    "Current": "v1.5.0",
    "Upgrade": "1.7.0"
   },
   "UpgradeConfidence": 0,
   "ActionItems": [
    {
     "ResourceNamespace": "cert-manager",
     "ResourceKind": "Deployment",
     "ResourceName": "cert-manager-cainjector",
     "Title": "Found cert with removed apiversion",
     "Description": "A deprecated or removed annotation was found",
     "Remediation": "Please update your ingress annotations to use the current versions. See https://cert-manager.io/docs/release-notes/release-notes-0.11/ for details",
     "EventType": "",
     "Severity": "0.1",
     "Category": "Reliability"
    }
   ],
   "Notes": "",
   "Warnings": [
    "no schema available, unable to validate release"
   ]
  }
 ]
}
```

This indicates that the version of `cert-manager` running in the cluster falls between versions `v1.5.0` and `1.70` and has triggered the OPA check defined in the bundle, which is looking for deprecated or removed annotations.

