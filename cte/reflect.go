package cte

import (
	"fmt"
	"github.com/jamestrandung/go-data-structure/set"
	"reflect"
	"strings"
)

type method struct {
	owningType string
	name       string
	arguments  string // comma-separated argument types
	outputs    string // comma-separated argument types
}

func (m method) hasSameSignature(other method) bool {
	return m.name == other.name &&
		m.arguments == other.arguments &&
		m.outputs == other.outputs
}

func extractMethodDetails(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
	var arguments []string
	for i := 0; i < rm.Type.NumIn(); i++ {
		if i == 0 && ignoreFirstReceiverArgument {
			continue
		}

		arguments = append(arguments, rm.Type.In(i).String())
	}

	var outputs []string
	for i := 0; i < rm.Type.NumOut(); i++ {
		outputs = append(outputs, rm.Type.Out(i).String())
	}

	return method{
		name:      rm.Name,
		arguments: strings.Join(arguments, ","),
		outputs:   strings.Join(outputs, ","),
	}
}

func (m method) String() string {
	methodSignature := func() string {
		if strings.Contains(m.outputs, ",") {
			return fmt.Sprintf("%s(%s) (%s)", m.name, m.arguments, m.outputs)
		}

		if m.outputs == "" {
			return fmt.Sprintf("%s(%s)", m.name, m.arguments)
		}

		return fmt.Sprintf("%s(%s) %s", m.name, m.arguments, m.outputs)
	}()

	if m.owningType != "" {
		return m.owningType + "." + methodSignature
	}

	return methodSignature
}

type structDisassembler struct {
	availableMethods             map[string]set.HashSet[method]
	methodsAvailableMoreThanOnce set.HashSet[method]
}

func newStructDisassembler() structDisassembler {
	return structDisassembler{
		availableMethods: make(map[string]set.HashSet[method]),
	}
}

func (sd structDisassembler) addAvailableMethod(m method) {
	methodSet, ok := sd.availableMethods[m.name]
	if !ok {
		methodSet = set.NewHashSet[method]()
		sd.availableMethods[m.name] = methodSet
	}

	if methodSet.Has(m) {
		sd.methodsAvailableMoreThanOnce.Add(m)
		return
	}

	methodSet.Add(m)
}

func (sd structDisassembler) isAvailableMoreThanOnce(m method) bool {
	return sd.methodsAvailableMoreThanOnce.Has(m)
}

func (sd structDisassembler) extractAvailableMethods(t reflect.Type) []method {
	var hoistedMethods []method

	actualType := t
	if actualType.Kind() == reflect.Pointer {
		actualType = t.Elem()
	}

	if actualType.Kind() == reflect.Struct {
		for i := 0; i < actualType.NumField(); i++ {
			rf := actualType.Field(i)

			// Extract methods from embedded fields
			if rf.Anonymous {
				childMethods := sd.extractAvailableMethods(rf.Type)
				hoistedMethods = append(hoistedMethods, childMethods...)
			}
		}
	}

	var ownMethods []method

	for i := 0; i < t.NumMethod(); i++ {
		rm := t.Method(i)

		m := extractMethodDetails(rm, true)
		m.owningType = t.PkgPath() + "/" + t.Name()

		isHoistedMethod := func() bool {
			for j := 0; j < len(hoistedMethods); j++ {
				hm := hoistedMethods[j]

				if hm.hasSameSignature(m) {
					return true
				}
			}

			return false
		}()

		// If a parent struct has a method carrying the same signature as one
		// that is available in an embedded field, it means this is a hoisted
		// method or the parent is overriding the same method with its own
		// implementation. In either case, we can assume this is a hoisted
		// method as it doesn't make a difference to how we should validate
		// code templates.
		if isHoistedMethod {
			continue
		}

		ownMethods = append(ownMethods, m)
		sd.addAvailableMethod(m)
	}

	var allMethods []method
	allMethods = append(allMethods, hoistedMethods...)
	allMethods = append(allMethods, ownMethods...)

	return allMethods
}
