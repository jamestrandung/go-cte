package cte

import (
	"reflect"
)

type metaType string

const (
	metaTypeKey   metaType = "key"
	metaTypeInout metaType = "inout"
)

type MetadataProvider interface {
	Metadata() any
}

func extractMetadata(mp MetadataProvider) map[metaType]reflect.Type {
	result := make(map[metaType]reflect.Type)

	metadata := mp.Metadata()
	if metadata == nil {
		panic(ErrNilMetadata.Err(reflect.TypeOf(mp)))
	}

	rt := reflect.TypeOf(metadata)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		result[metaType(field.Name)] = field.Type
	}

	return result
}
