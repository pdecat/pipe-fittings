package app_specific_connection

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/turbot/pipe-fittings/connection"
	"github.com/turbot/pipe-fittings/perr"
	"github.com/zclconf/go-cty/cty"
	"reflect"
)

var BaseConnectionCtyType = cty.Capsule("BaseConnectionCtyType", reflect.TypeOf(&connection.ConnectionImpl{}))

func NewPipelingConnection(connectionType, shortName string, declRange hcl.Range) (connection.PipelingConnection, error) {
	ctor, exists := ConnectionTypeRegistry[connectionType]
	if !exists {
		return nil, perr.BadRequestWithMessage("Invalid connection type " + connectionType)
	}

	return ctor(shortName, declRange), nil
}

func ConnectionCtyType(connectionType string) cty.Type {
	// NOTE: if not type is provided, just return the type of ConnectionImpl
	if connectionType == "" {
		return BaseConnectionCtyType
	}

	ctor, exists := ConnectionTypeRegistry[connectionType]
	if !exists {
		return cty.NilType
	}
	// instantiate connection
	inst := ctor("", hcl.Range{})
	goType := reflect.TypeOf(inst)
	// dereference pointer
	if goType.Kind() == reflect.Ptr {
		goType = goType.Elem()
	}

	return cty.Capsule(connectionType, goType)
}

func DefaultPipelingConnections() (map[string]connection.PipelingConnection, error) {
	conns := make(map[string]connection.PipelingConnection)

	for k := range ConnectionTypeRegistry {
		// check for a default connection of this type
		defaultConnection, exists := DefaultConnections[k]
		// if no default was specified, create an empty connection
		if !exists {
			var err error
			defaultConnection, err = NewPipelingConnection(k, "default", hcl.Range{})
			if err != nil {
				return nil, err
			}
		}

		conns[k+".default"] = defaultConnection
	}

	return conns, nil
}
