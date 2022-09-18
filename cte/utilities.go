package cte

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/jamestrandung/go-data-structure/set"
)

func swallowErrPlanExecutionEndingEarly(err error) error {
	// Execution was intentionally ended by clients
	if err == ErrPlanExecutionEndingEarly || err == ErrRootPlanExecutionEndingEarly {
		return nil
	}

	return err
}

func extractFullNameFromValue(v any) string {
	return extractFullNameFromType(reflect.TypeOf(v))
}

var extractFullNameFromType = func(t reflect.Type) string {
	t = extractNonPointerType(t)

	return t.PkgPath() + "/" + t.Name()
}

func extractShortName(fullName string) string {
	shortNameIdx := strings.LastIndex(fullName, "/")
	return fullName[shortNameIdx+1:]
}

func extractFieldTypes(field reflect.StructField) (isPointerType bool, valueType reflect.Type, pointerType reflect.Type) {
	rawFieldType := field.Type
	isPointerType = rawFieldType.Kind() == reflect.Pointer

	valueType = rawFieldType
	if isPointerType {
		valueType = rawFieldType.Elem()
	}

	pointerType = reflect.PointerTo(valueType)

	return
}

var extractUnderlyingType = func(v reflect.Value) reflect.Type {
	return extractNonPointerType(v.Type())
}

var extractNonPointerType = func(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Pointer {
		return t.Elem()
	}

	return t
}

func toString[T comparable](s set.HashSet[T]) string {
	all := make([]string, 0, s.Count())
	for _, element := range s.Items() {
		all = append(all, fmt.Sprintf("%v", element))
	}

	sort.Strings(all)

	return strings.Join(all, ", ")
}
