package parse

import (
	"fmt"
	"slices"

	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/v2/hclhelpers"
	"github.com/turbot/pipe-fittings/v2/modconfig"
)

func ResolveChildrenFromNames(childNames []string, block *hcl.Block, supportedChildren []string, parseCtx *ModParseContext) ([]modconfig.ModTreeItem, hcl.Diagnostics) {
	// TODO #validate validate all children are same type (i.e. we do not support detections and controls in same tree) https://github.com/turbot/pipe-fittings/v2/issues/612
	var diags hcl.Diagnostics
	diags = checkForDuplicateChildren(childNames, block)
	if diags.HasErrors() {
		return nil, diags
	}

	// find the children in the eval context and populate control children
	children := make([]modconfig.ModTreeItem, len(childNames))

	for i, childName := range childNames {
		parsedName, err := modconfig.ParseResourceName(childName)
		if err != nil || !slices.Contains(supportedChildren, parsedName.ItemType) {
			diags = append(diags, childErrorDiagnostic(childName, block))
			continue
		}

		// now get the resource from the parent mod
		// find the mod which owns this resource - it may be either the current mod, or one of it's direct dependencies
		var mod = parseCtx.GetMod(parsedName.Mod)
		if mod == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Could not resolve mod for child %s", childName),
				Subject:  &block.TypeRange,
			})
			break
		}

		resource, found := mod.GetResource(parsedName)
		// ensure this item is a mod tree item
		child, ok := resource.(modconfig.ModTreeItem)
		if !found || !ok {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Could not resolve child %s", childName),
				Subject:  &block.TypeRange,
			})
			continue
		}

		children[i] = child
	}
	if diags.HasErrors() {
		return nil, diags
	}

	return children, nil
}

func checkForDuplicateChildren(names []string, block *hcl.Block) hcl.Diagnostics {
	var diags hcl.Diagnostics
	// validate each child name appears only once
	nameMap := make(map[string]int)
	for _, n := range names {
		nameCount := nameMap[n]
		// raise an error if this name appears more than once (but only raise 1 error per name)
		if nameCount == 1 {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("'%s.%s' has duplicate child name '%s'", block.Type, block.Labels[0], n),
				Subject:  hclhelpers.BlockRangePointer(block)})
		}
		nameMap[n] = nameCount + 1
	}

	return diags
}

func childErrorDiagnostic(childName string, block *hcl.Block) *hcl.Diagnostic {
	return &hcl.Diagnostic{
		Severity: hcl.DiagError,
		Summary:  fmt.Sprintf("Invalid child %s", childName),
		Subject:  &block.TypeRange,
	}
}

func GetChildNameStringsFromModTreeItem(children []modconfig.ModTreeItem) []string {
	res := make([]string, len(children))
	for i, n := range children {
		res[i] = n.Name()
	}
	return res
}
