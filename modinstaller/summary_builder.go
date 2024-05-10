package modinstaller

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/v2/constants"
	"github.com/turbot/pipe-fittings/v2/utils"
	"github.com/turbot/pipe-fittings/v2/versionmap"
)

const (
	VerbInstalled   = "Installed"
	VerbUninstalled = "Uninstalled"
	VerbUpgraded    = "Upgraded"
	VerbDowngraded  = "Downgraded"
	VerbPruned      = "Pruned"
)

var dryRunVerbs = map[string]string{
	VerbInstalled:   "Would install",
	VerbUninstalled: "Would uninstall",
	VerbUpgraded:    "Would upgrade",
	VerbDowngraded:  "Would downgrade",
	VerbPruned:      "Would prune",
}

func getVerb(verb string) string {
	if viper.GetBool(constants.ArgDryRun) {
		verb = dryRunVerbs[verb]
	}
	return verb
}

func BuildInstallSummary(installData *InstallData) string {
	/* - walk source lock, building path without version number
	   - look for corresponding node in new lock, check for added/deleted/updated/downgraded

	 OR

	- walk old lock, down tree starting at root and traversing down to leaf nodes

	*/

	// for now treat an install as update - we only install deps which are in the mod.sp but missing in the mod folder
	modDependencyPath := installData.WorkspaceMod.GetInstallCacheKey()
	installCount, installedTreeString := getInstallationResultString(installData.Installed, modDependencyPath, installData.NewLock)
	uninstallCount, uninstalledTreeString := getInstallationResultString(installData.Uninstalled, modDependencyPath, installData.NewLock)
	upgradeCount, upgradeTreeString := getInstallationResultString(installData.Upgraded, modDependencyPath, installData.NewLock)
	downgradeCount, downgradeTreeString := getInstallationResultString(installData.Downgraded, modDependencyPath, installData.NewLock)

	var installString, upgradeString, downgradeString, uninstallString string
	if installCount > 0 {
		verb := getVerb(VerbInstalled)
		installString = fmt.Sprintf("\n%s %d %s:\n\n%s\n", verb, installCount, utils.Pluralize("mod", installCount), installedTreeString)
	}
	if uninstallCount > 0 {
		verb := getVerb(VerbUninstalled)
		uninstallString = fmt.Sprintf("\n%s %d %s:\n\n%s\n", verb, uninstallCount, utils.Pluralize("mod", uninstallCount), uninstalledTreeString)
	}
	if upgradeCount > 0 {
		verb := getVerb(VerbUpgraded)
		upgradeString = fmt.Sprintf("\n%s %d %s:\n\n%s\n", verb, upgradeCount, utils.Pluralize("mod", upgradeCount), upgradeTreeString)
	}
	if downgradeCount > 0 {
		verb := getVerb(VerbDowngraded)
		downgradeString = fmt.Sprintf("\n%s %d %s:\n\n%s\n", verb, downgradeCount, utils.Pluralize("mod", downgradeCount), downgradeTreeString)
	}

	if installCount+uninstallCount+upgradeCount+downgradeCount == 0 {
		if len(installData.Lock.InstallCache) == 0 {
			return "No mods are installed"
		}
		return "All mods are up to date"
	}
	return fmt.Sprintf("%s%s%s%s", installString, upgradeString, downgradeString, uninstallString)
}

func getInstallationResultString(items versionmap.InstalledDependencyVersionsMap, modDependencyPath string, lock *versionmap.WorkspaceLock) (int, string) {
	var res string
	count := len(items.FlatMap())
	if count > 0 {
		tree := items.GetDependencyTree(modDependencyPath, lock)
		res = tree.String()
	}
	return count, res
}

func BuildUninstallSummary(installData *InstallData) string {
	// for now treat an install as update - we only install deps which are in the mod.sp but missing in the mod folder
	uninstallCount := len(installData.Uninstalled.FlatMap())
	if uninstallCount == 0 {
		return "Nothing uninstalled"
	}
	uninstalledTree := installData.GetUninstalledTree()

	verb := getVerb(VerbUninstalled)
	return fmt.Sprintf("\n%s %d %s:\n\n%s", verb, uninstallCount, utils.Pluralize("mod", uninstallCount), uninstalledTree.String())
}
