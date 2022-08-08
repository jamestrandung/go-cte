package cte

import (
	"reflect"
)

func isComplete(e Engine, rp reflect.Value) error {
	planName := extractFullNameFromType(extractUnderlyingType(rp))

	sd := newStructDisassembler()

	sd.extractAvailableMethods(rp.Type())

	var verifyFn func(planName string) error
	verifyFn = func(planName string) error {
		ap := e.findAnalyzedPlan(planName, rp)

		for _, h := range ap.preHooks {
			expectedInout, ok := h.metadata.getInoutInterface()
			if !ok {
				return ErrInoutMetaMissing.Err(reflect.TypeOf(h.hook))
			}

			err := isInterfaceSatisfied(sd, expectedInout)
			if err != nil {
				return ErrPlanNotMeetingInoutRequirements.Err(rp.Type(), expectedInout, err.Error())
			}
		}

		for _, component := range ap.components {
			if c, ok := e.computers[component.id]; ok {
				expectedInout, ok := c.metadata.getInoutInterface()
				if !ok {
					return ErrInoutMetaMissing.Err(component.id)
				}

				err := isInterfaceSatisfied(sd, expectedInout)
				if err != nil {
					return ErrPlanNotMeetingInoutRequirements.Err(rp.Type(), expectedInout, err.Error())
				}
			}

			if _, ok := e.plans[component.id]; ok {
				if err := verifyFn(component.id); err != nil {
					return err
				}
			}
		}

		for _, h := range ap.postHooks {
			expectedInout, ok := h.metadata.getInoutInterface()
			if !ok {
				return ErrInoutMetaMissing.Err(reflect.TypeOf(h.hook))
			}

			err := isInterfaceSatisfied(sd, expectedInout)
			if err != nil {
				return ErrPlanNotMeetingInoutRequirements.Err(rp.Type(), expectedInout, err.Error())
			}
		}

		return nil
	}

	return verifyFn(planName)
}

func isInterfaceSatisfied(sd structDisassembler, expectedInterface reflect.Type) error {
	for i := 0; i < expectedInterface.NumMethod(); i++ {
		rm := expectedInterface.Method(i)

		requiredMethod := extractMethodDetails(rm, false)

		methodSet, ok := sd.availableMethods[requiredMethod.name]
		if !ok {
			return ErrPlanMissingMethod.Err(requiredMethod)
		}

		if methodSet.Count() > 1 {
			return ErrPlanHavingAmbiguousMethods.Err(requiredMethod, methodSet)
		}

		foundMethod := methodSet.Items()[0]

		if !foundMethod.hasSameSignature(requiredMethod) {
			return ErrPlanHavingMethodButSignatureMismatched.Err(requiredMethod, foundMethod)
		}

		if sd.isAvailableMoreThanOnce(foundMethod) {
			return ErrPlanHavingSameMethodRegisteredMoreThanOnce.Err(foundMethod)
		}
	}

	return nil
}
