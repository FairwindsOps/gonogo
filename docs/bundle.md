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

- **versions**: the begin and end versions of the addon to be checked for in the cluster.
- **notes**: free-form string value that can be used for internal information
- **source**: gonogo will check to see if there is a values.schema.json file in the source chart repo. More info on using this schema validation can be found in this (https://austindewey.com/2020/06/13/helm-tricks-input-validation-with-values-schema-json/)[article].
- **warning**: a free-form string value that can be used for internal information
- **compatible_k8s_versions**: the begin and end cluster versions supported by the addon
- **necessary_api_versions**: apis that must be present in the cluster for the addon to succeed
- **values_schema**: string value that can be used to define inline (schema validation)[https://helm.sh/docs/topics/charts/#schema-files]
- **resources**: a list of cluster objects to be checked during OPA validation
- **opa_checks**: a string value that can be used to define inline OPA policies using (https://medium.com/@mathurvarun98/how-to-write-great-rego-policies-dc6117679c9f)[Rego].

Example of specifying a `values_schema` value:

```
values_schema: |
    {
      "$schema": "http://json-schema.org/schema#",
      "type": "object",
      "properties": {
        "image": {
          "type": "object",
          "required": [
            "repository",
            "pullPolicy"
            ],
            "properties": {
              "repository": {
                "type": "string",
                "pattern": "^[a-z0-9]+"
              },
              "pullPolicy": {
                "type": "string",
                "pattern": "Always"
              }
            }
          }
        }
      }
```

# How GoNoGO Uses the Bundle
GoNoGo first compares the list of addons in your bundle spec to the Helm releases in you cluster. It only runs checks against addons that have a successfully deployed release in your Kubernetes cluster.

It will then check to see if there are user-defined values in use for the release. If it finds that there are, GoNoGo will attempt to validate those values against a schema. It will first look to see if you have specified a value for the `values_schema` key, and validate against that entry. If you do not specify the `values_schema` key, GoNoGo will attempt to look at the upstream chart repo for a `values.json.schema` file and use that as the schema. If there is none present GoNoGo will move on to the next check, OPA checks.

If you have specified a value for the `opa_checks` key, GoNoGo will run your OPA check against the individual yaml files found in the Helm release for the addon. If you have also specified a `resources` value, GoNoGo will also run your OPA check against object yaml in your cluster of that resource type. This allows you to check for resources that are not included in the Helm chart. For example, with `cert-manager` there are deprecated annotations that are used in objects/yaml not included in the `cert-manager` chart itself, but rather in `ingress` objects. This allows you to specify reviewing all ingress objects in your cluster for the deprecated annotation.

Finally GoNoGo runs checks against the values you provide for the K8s version and API versions and your cluster info.

