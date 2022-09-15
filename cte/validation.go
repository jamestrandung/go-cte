package cte

import (
	"reflect"
	"strings"
)

func isComplete(e Engine, planValue reflect.Value) error {
	sd := newStructDisassembler()
	sd.self().extractAvailableMethods(planValue.Type())

	var cs componentStack
	rootPlanName := extractFullNameFromType(planValue.Type())

	var verifyFn func(planName string, curPlanValue reflect.Value) error
	verifyFn = func(planName string, curPlanValue reflect.Value) error {
		ap := e.findAnalyzedPlan(planName, curPlanValue)

		cs = cs.push(planName)
		defer func() {
			cs = cs.pop()
		}()

		for _, h := range ap.preHooks {
			err := verifyComponentCompleteness(h.metadata, sd, cs, rootPlanName, reflect.TypeOf(h.hook).String(), planValue.Type().String())
			if err != nil {
				return err
			}
		}

		for _, component := range ap.components {
			if c, ok := e.computers[component.id]; ok {
				err := verifyComponentCompleteness(c.metadata, sd, cs, rootPlanName, component.id, planValue.Type().String())
				if err != nil {
					return err
				}
			}

			if _, ok := e.plans[component.id]; ok {
				nestedPlanValue := func() reflect.Value {
					if curPlanValue.Kind() == reflect.Pointer {
						return curPlanValue.Elem().Field(component.fieldIdx)
					}

					return curPlanValue.Field(component.fieldIdx)
				}()

				if err := verifyFn(component.id, nestedPlanValue); err != nil {
					return err
				}
			}
		}

		for _, h := range ap.postHooks {
			err := verifyComponentCompleteness(h.metadata, sd, cs, rootPlanName, reflect.TypeOf(h.hook).String(), planValue.Type().String())
			if err != nil {
				return err
			}
		}

		return nil
	}

	return verifyFn(rootPlanName, planValue)
}

var verifyComponentCompleteness = func(
	pm parsedMetadata,
	sd structDisassembler,
	cs componentStack,
	rootPlanName string,
	componentID string,
	planType string,
) error {
	expectedInout, ok := pm.getInoutInterface()
	if !ok {
		return ErrInoutMetaMissing.Err(componentID)
	}

	err := isInterfaceSatisfied(sd, expectedInout, rootPlanName)
	if err != nil {
		cs = cs.push(componentID)
		return ErrPlanNotMeetingInoutRequirements.Err(planType, expectedInout, err.Error(), cs)
	}

	return nil
}

var isInterfaceSatisfied = func(sd structDisassembler, expectedInterface reflect.Type, rootPlanName string) error {
	for i := 0; i < expectedInterface.NumMethod(); i++ {
		rm := expectedInterface.Method(i)

		requiredMethod := extractMethodDetails(rm, false)

		methodSet, ok := sd.self().findAvailableMethods(requiredMethod.name)
		if !ok {
			return ErrPlanMissingMethod.Err(requiredMethod)
		}

		if methodSet.Count() > 1 {
			methodLocations := sd.self().findMethodLocations(methodSet, rootPlanName)
			return ErrPlanHavingAmbiguousMethods.Err(requiredMethod, toString(methodSet), strings.Join(methodLocations, "; "))
		}

		foundMethod := methodSet.Items()[0]

		if !foundMethod.hasSameSignature(requiredMethod) {
			return ErrPlanHavingMethodButSignatureMismatched.Err(requiredMethod, foundMethod)
		}

		if sd.self().isAvailableMoreThanOnce(foundMethod) {
			methodLocations := sd.self().findMethodLocations(methodSet, rootPlanName)
			return ErrPlanHavingSameMethodRegisteredMoreThanOnce.Err(foundMethod, strings.Join(methodLocations, "; "))
		}
	}

	return nil
}

type componentStack []string

func (s componentStack) push(componentName string) componentStack {
	return append(s, componentName)
}

func (s componentStack) pop() componentStack {
	return s[0 : len(s)-1]
}

func (s componentStack) clone() componentStack {
	result := make([]string, 0, len(s))

	for _, c := range s {
		result = append(result, c)
	}

	return result
}

func (s componentStack) String() string {
	return strings.Join(s, " >> ")
}
