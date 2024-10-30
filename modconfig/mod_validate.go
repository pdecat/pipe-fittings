package modconfig

import (
	"fmt"

	"github.com/turbot/pipe-fittings/error_helpers"
)

// ensure we have resolved all Children in the resource tree
func (m *Mod) validateResourceTree() error {
	var errors []error
	for _, child := range m.GetChildren() {
		if err := m.validateChildren(child); err != nil {
			errors = append(errors, err)
		}
	}
	return error_helpers.CombineErrorsWithPrefix(fmt.Sprintf("failed to resolve Children for %d resources", len(errors)), errors...)
}

func (m *Mod) validateChildren(item ModTreeItem) error {
	missing := 0
	for _, child := range item.GetChildren() {
		if child == nil {
			missing++

		}
	}
	if missing > 0 {
		return fmt.Errorf("%s has %d unresolved Children", item.Name(), missing)
	}
	return nil
}
