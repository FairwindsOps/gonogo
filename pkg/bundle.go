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



