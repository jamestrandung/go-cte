package cte

import (
	"reflect"
)

type metaType string

const (
	metaTypeKey      metaType = "key"
	metaTypeComputer metaType = "computer"
	metaTypeInout    metaType = "inout"
)

type MetadataProvider interface {
	Metadata() any
}

func extractMetadata(mp MetadataProvider) map[metaType]reflect.Type {
	result := make(map[metaType]reflect.Type)

	metadata := mp.Metadata()
	if metadata == nil {
		panic(ErrMissingMetadata.Err(reflect.TypeOf(mp)))
	}

	rt := reflect.TypeOf(metadata)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		fieldType := func() reflect.Type {
			if field.Type.Kind() == reflect.Pointer {
				return field.Type.Elem()
			}

			return field.Type
		}()

		result[metaType(field.Name)] = fieldType
	}

	return result
}

type metadataParser struct{}

func parseMetadata(metadata any) map[metaType]reflect.Type {
	if metadata == nil {
		panic(ErrNilMetadata)
	}

	result := make(map[metaType]reflect.Type)

	rt := reflect.TypeOf(metadata)
	if rt.Kind() == reflect.Pointer {
		rt = rt.Elem()
	}

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		fieldType := func() reflect.Type {
			if field.Type.Kind() == reflect.Pointer {
				return field.Type.Elem()
			}

			return field.Type
		}()

		result[metaType(field.Name)] = fieldType
	}

	return result
}
