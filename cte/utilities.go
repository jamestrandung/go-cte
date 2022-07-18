package cte

import (
	"reflect"
	"strings"
)

func swallowErrPlanExecutionEndingEarly(err error) error {
	// Execution was intentionally ended by clients
	if err == ErrPlanExecutionEndingEarly || err == ErrRootPlanExecutionEndingEarly {
		return nil
	}

	return err
}

func extractFullNameFromValue(v any) string {
	rv := reflect.ValueOf(v)

	rt := func() reflect.Type {
		if rv.Kind() == reflect.Pointer {
			return rv.Elem().Type()
		}

		return rv.Type()
	}()

	return extractFullNameFromType(rt)
}

func extractFullNameFromType(t reflect.Type) string {
	return t.PkgPath() + "/" + t.Name()
}

func extractShortName(fullName string) string {
	shortNameIdx := strings.LastIndex(fullName, "/")
	return fullName[shortNameIdx+1:]
}
