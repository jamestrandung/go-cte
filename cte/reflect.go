package cte

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jamestrandung/go-data-structure/set"
)

type methodLocation struct {
	rootPlanName string
	componentStack
}

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

var extractMethodDetails = func(rm reflect.Method, ignoreFirstReceiverArgument bool) method {
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

//go:generate mockery --name iStructDisassembler --case=underscore --inpackage
type iStructDisassembler interface {
	isAvailableMoreThanOnce(m method) bool
	findAvailableMethods(name string) (set.HashSet[method], bool)
	findMethodLocations(methodSet set.HashSet[method], rootPlanName string) []string
	addAvailableMethod(rootPlanName string, cs componentStack, m method)
	extractAvailableMethods(t reflect.Type) []method
	performMethodExtraction(t reflect.Type, rootPlanName string, cs componentStack) []method
	extractChildMethods(t reflect.Type, rootPlanName string, cs componentStack) []method
	extractOwnMethods(t reflect.Type, rootPlanName string, cs componentStack, hoistedMethods []method) []method
}

type structDisassembler struct {
	itself                       iStructDisassembler
	availableMethods             map[string]set.HashSet[method]
	methodsAvailableMoreThanOnce set.HashSet[method]
	methodLocations              map[method][]methodLocation
}

func newStructDisassembler() *structDisassembler {
	result := &structDisassembler{
		availableMethods:             make(map[string]set.HashSet[method]),
		methodsAvailableMoreThanOnce: set.NewHashSet[method](),
		methodLocations:              make(map[method][]methodLocation),
	}

	result.itself = result

	return result
}

func (sd *structDisassembler) isAvailableMoreThanOnce(m method) bool {
	return sd.methodsAvailableMoreThanOnce.Has(m)
}

func (sd *structDisassembler) findAvailableMethods(name string) (set.HashSet[method], bool) {
	result, ok := sd.availableMethods[name]
	return result, ok && result.Count() > 0
}

func (sd *structDisassembler) findMethodLocations(methodSet set.HashSet[method], rootPlanName string) []string {
	var methodLocations []string
	for _, m := range methodSet.Items() {
		for _, ml := range sd.methodLocations[m] {
			if ml.rootPlanName == rootPlanName {
				methodLocations = append(methodLocations, ml.componentStack.String())
			}
		}
	}

	return methodLocations
}

func (sd *structDisassembler) addAvailableMethod(rootPlanName string, cs componentStack, m method) {
	methodSet, ok := sd.availableMethods[m.name]
	if !ok {
		methodSet = set.NewHashSet[method]()
		sd.availableMethods[m.name] = methodSet
	}

	// Even if a method is registered twice, they will be declared at different locations.
	// Hence, the provided component stack will never be the same and thus, we must record
	// the location before checking for duplicate.
	sd.methodLocations[m] = append(
		sd.methodLocations[m], methodLocation{
			rootPlanName:   rootPlanName,
			componentStack: cs.clone(),
		},
	)

	if methodSet.Has(m) {
		sd.methodsAvailableMoreThanOnce.Add(m)
		return
	}

	methodSet.Add(m)
}

func (sd *structDisassembler) extractAvailableMethods(t reflect.Type) []method {
	var cs componentStack
	return sd.performMethodExtraction(t, extractFullNameFromType(t), cs)
}

func (sd *structDisassembler) performMethodExtraction(t reflect.Type, rootPlanName string, cs componentStack) []method {
	cs = cs.push(extractFullNameFromType(t))
	defer func() {
		cs = cs.pop()
	}()

	hoistedMethods := sd.itself.extractChildMethods(t, rootPlanName, cs)
	ownMethods := sd.itself.extractOwnMethods(t, rootPlanName, cs, hoistedMethods)

	var allMethods []method
	allMethods = append(allMethods, hoistedMethods...)
	allMethods = append(allMethods, ownMethods...)

	return allMethods
}

func (sd *structDisassembler) extractChildMethods(t reflect.Type, rootPlanName string, cs componentStack) []method {
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
				childMethods := sd.itself.performMethodExtraction(rf.Type, rootPlanName, cs)
				hoistedMethods = append(hoistedMethods, childMethods...)
			}
		}
	}

	return hoistedMethods
}

func (sd *structDisassembler) extractOwnMethods(t reflect.Type, rootPlanName string, cs componentStack, hoistedMethods []method) []method {
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
		sd.itself.addAvailableMethod(rootPlanName, cs, m)
	}

	return ownMethods
}
