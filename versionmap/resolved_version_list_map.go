package versionmap

import (
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

// ResolvedVersionListMap represents a map of ResolvedVersionConstraint arrays, keyed by dependency name
type ResolvedVersionListMap map[string][]*InstalledModVersion

// Add appends the version constraint to the list for the given name
func (m ResolvedVersionListMap) Add(name string, versionConstraint *InstalledModVersion) {
	// if there is already an entry for the same name, replace it
	m[name] = []*InstalledModVersion{versionConstraint}
}

// Remove removes the given version constraint from the list for the given name
func (m ResolvedVersionListMap) Remove(name string, constraint *ResolvedVersionConstraint) {
	var res []*InstalledModVersion
	for _, c := range m[name] {
		if !c.Equals(constraint) {
			res = append(res, c)
		}
	}
	m[name] = res
}

// FlatMap converts the ResolvedVersionListMap map into a map keyed by the FULL dependency name (i.e. including version(
func (m ResolvedVersionListMap) FlatMap() map[string]*InstalledModVersion {
	var res = make(map[string]*InstalledModVersion)
	for name, versions := range m {
		for _, version := range versions {
			key := modconfig.BuildModDependencyPath(name, &version.DependencyVersion)
			res[key] = version
		}
	}
	return res
}

// FlatNames converts the ResolvedVersionListMap map into a string array of full names
func (m ResolvedVersionListMap) FlatNames() []string {
	var res []string
	for name, versions := range m {
		for _, version := range versions {
			res = append(res, modconfig.BuildModDependencyPath(name, &version.DependencyVersion))
		}
	}
	return res
}
