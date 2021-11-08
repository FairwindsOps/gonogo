/*
Copyright © 2021 FairwindsOps Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bundle

import "fmt"

type Source struct {
	Chart string
	Repository string
}

type Version struct {
	Start string
	End string
}

type Bundle struct {
	Name string
	Version Version
	Notes string
	Source Source
	Warnings []string
	Compatible_k8s_versions []string
	Necessary_api_versions []string
	Values_schame string
	Opa_checks []string
}



