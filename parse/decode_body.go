package parse

import (
	"reflect"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/turbot/go-kit/helpers"
	"github.com/turbot/pipe-fittings/hclhelpers"
	"github.com/turbot/pipe-fittings/modconfig"
)

func DecodeHclBody(body hcl.Body, evalCtx *hcl.EvalContext, resourceProvider modconfig.ResourceProvider, resource modconfig.HclResource) (diags hcl.Diagnostics) {
	defer func() {
		if r := recover(); r != nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "unexpected error in DecodeHclBody",
				Detail:   helpers.ToError(r).Error()})
		}
	}()

	nestedStructs, moreDiags := GetNestedStructValsRecursive(resource)
	diags = append(diags, moreDiags...)
	// get the schema for this resource
	schema := getResourceSchema(resource, nestedStructs)
	// handle invalid block types
	moreDiags = validateHcl(resource.GetBlockType(), body.(*hclsyntax.Body), schema)
	diags = append(diags, moreDiags...)

	moreDiags = decodeHclBodyIntoStruct(body, evalCtx, resourceProvider, resource)
	diags = append(diags, moreDiags...)

	for _, nestedStruct := range nestedStructs {
		moreDiags := decodeHclBodyIntoStruct(body, evalCtx, resourceProvider, nestedStruct)
		diags = append(diags, moreDiags...)
	}

	return diags
}

func decodeHclBodyIntoStruct(body hcl.Body, evalCtx *hcl.EvalContext, resourceProvider modconfig.ResourceProvider, resource any) hcl.Diagnostics {
	var diags hcl.Diagnostics
	// call decodeHclBodyIntoStruct to do actual decode
	moreDiags := gohcl.DecodeBody(body, evalCtx, resource)
	diags = append(diags, moreDiags...)

	// resolve any resource references using the resource map, rather than relying on the EvalCtx
	// (which does not work with nested struct vals)
	moreDiags = resolveReferences(body, resourceProvider, resource)
	diags = append(diags, moreDiags...)
	return diags
}

func GetSchemaForStruct(t reflect.Type) *hcl.BodySchema {
	var schema = &hcl.BodySchema{}
	// get all hcl tags
	for i := 0; i < t.NumField(); i++ {
		tag := t.FieldByIndex([]int{i}).Tag.Get("hcl")
		if tag == "" {
			continue
		}
		if idx := strings.LastIndex(tag, ",block"); idx != -1 {
			blockName := tag[:idx]
			schema.Blocks = append(schema.Blocks, hcl.BlockHeaderSchema{Type: blockName})
		} else {
			attributeName := strings.Split(tag, ",")[0]
			if attributeName != "" {
				schema.Attributes = append(schema.Attributes, hcl.AttributeSchema{Name: attributeName})
			}
		}
	}
	return schema
}

// rather than relying on the evaluation context to resolve resource references
// (which has the issue that when deserializing from cty we do not receive all base struct values)
// instead resolve the reference by parsing the resource name and finding the resource in the ResourceMap
// and use this resource to set the target property
func resolveReferences(body hcl.Body, modResourcesProvider modconfig.ResourceProvider, val any) (diags hcl.Diagnostics) {
	defer func() {
		if r := recover(); r != nil {
			if r := recover(); r != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "unexpected error in resolveReferences",
					Detail:   helpers.ToError(r).Error()})
			}
		}
	}()
	attributes := body.(*hclsyntax.Body).Attributes
	rv := reflect.ValueOf(val)
	for rv.Type().Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	ty := rv.Type()
	if ty.Kind() != reflect.Struct {
		return
	}

	ct := ty.NumField()
	for i := 0; i < ct; i++ {
		field := ty.Field(i)
		fieldVal := rv.Field(i)
		// get hcl attribute tag (if any) tag
		hclAttribute := getHclAttributeTag(field)
		if hclAttribute == "" {
			continue
		}
		if fieldVal.Type().Kind() == reflect.Pointer && !fieldVal.IsNil() {
			fieldVal = fieldVal.Elem()
		}
		if fieldVal.Kind() == reflect.Struct {
			v := fieldVal.Addr().Interface()
			if _, ok := v.(modconfig.HclResource); ok {
				if hclVal, ok := attributes[hclAttribute]; ok {
					if scopeTraversal, ok := hclVal.Expr.(*hclsyntax.ScopeTraversalExpr); ok {
						path := hclhelpers.TraversalAsString(scopeTraversal.Traversal)
						if parsedName, err := modconfig.ParseResourceName(path); err == nil {
							if r, ok := modResourcesProvider.GetResource(parsedName); ok {
								f := rv.FieldByName(field.Name)
								if f.IsValid() && f.CanSet() {
									targetVal := reflect.ValueOf(r)
									f.Set(targetVal)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func getHclAttributeTag(field reflect.StructField) string {
	tag := field.Tag.Get("hcl")
	if tag == "" {
		return ""
	}

	comma := strings.Index(tag, ",")
	var name, kind string
	if comma != -1 {
		name = tag[:comma]
		kind = tag[comma+1:]
	} else {
		name = tag
		kind = "attr"
	}

	switch kind {
	case "attr":
		return name
	default:
		return ""
	}
}

func GetNestedStructValsRecursive(val any) ([]any, hcl.Diagnostics) {
	nested, diags := getNestedStructVals(val)
	res := nested

	for _, n := range nested {
		nestedVals, moreDiags := GetNestedStructValsRecursive(n)
		diags = append(diags, moreDiags...)
		res = append(res, nestedVals...)
	}
	return res, diags

}

// GetNestedStructVals return a slice of any nested structs within val
func getNestedStructVals(val any) (_ []any, diags hcl.Diagnostics) {
	defer func() {
		if r := recover(); r != nil {
			if r := recover(); r != nil {
				diags = append(diags, &hcl.Diagnostic{
					Severity: hcl.DiagError,
					Summary:  "unexpected error in resolveReferences",
					Detail:   helpers.ToError(r).Error()})
			}
		}
	}()

	rv := reflect.ValueOf(val)
	for rv.Type().Kind() == reflect.Pointer {
		rv = rv.Elem()
	}
	ty := rv.Type()
	if ty.Kind() != reflect.Struct {
		return nil, nil
	}
	ct := ty.NumField()
	var res []any
	for i := 0; i < ct; i++ {
		field := ty.Field(i)
		fieldVal := rv.Field(i)
		if field.Anonymous && fieldVal.Kind() == reflect.Struct {
			res = append(res, fieldVal.Addr().Interface())
		}
	}
	return res, nil
}
