---
meta:
  - name: description
    content: "Fairwinds GoNoGo | Documentation"
---
# Bundle Creation

GoNoGo relies on a file called a bundle spec. The bundle spec is a yaml document that defines the addons to check against, and the various conditions to check for prior to upgrading said addons. The top level key is called `addons` and contains a list of maps of conditions to check. An example bundle spec with one entry for the `cert-manager` addon could look like this:

```
addons:
- name: cert-manager
  versions:
    start: 1.5.0
    end: 1.7.0
  notes: A text field with general notes
  source:
    chart: cert-manager
    repository: https://charts.jetstack.io
  warnings:
  - "warning 1"
  - "warning 2"
  compatible_k8s_versions:
    max: 1.21
    min: 1.18
  necessary_api_versions:
  - apps/v1
  - v1
  values_schema: ""
  resources:
  - "v1/secrets"
  opa_checks:
  - >
    package Fairwinds
    removedAPIVersions[actionItem] {
        input.kind == "Deployment"
        input.metadata.annotations[k] = _
        startswith(k, "")
        actionItem := {
          "title": "Found cert with removed apiversion",
          "description": "A deprecated or removed annotation was found",
          "severity": 0.1,
          "remediation": "Please update your ingress annotations to use the current versions. See https://cert-manager.io/docs/release-notes/release-notes-0.11/ for details",
          "category": "Reliability"
        }
      }
```

Here is a breakdown of the different fields or keys

-- versions: the begin and end versions of the addon to be checked for in the cluster.
-- notes: free-form string value that can be used for internal information
-- source: gonogo will check to see if there is a values.schema.json file in the source chart repo. More info on using this schema validation can be found in this (https://austindewey.com/2020/06/13/helm-tricks-input-validation-with-values-schema-json/)[article].
-- warning: a free-form string value that can be used for internal information
-- compatible_k8s_versions: the begin and end cluster versions supported by the addon
-- necessary_api_versions: apis that must be present in the cluster for the addon to succeed
-- values_schema: string value that can be used to define inline schema validation
-- resources: a list of cluster objects to be checked during OPA validation
-- opa_checks: a string value that can be used to define inline OPA policies using (https://medium.com/@mathurvarun98/how-to-write-great-rego-policies-dc6117679c9f)[Rego].