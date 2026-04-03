package icd11

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

const (
	SystemICD11 = "http://id.who.int/icd/release/11/mms"
)

// NewCoding creates a new FHIR r5.Coding for an ICD-11 code
func NewCoding(code, display string) r5.Coding {
	system := SystemICD11
	return r5.Coding{
		System:  &system,
		Code:    &code,
		Display: &display,
	}
}

// NewCodeableConcept creates a new FHIR r5.CodeableConcept for an ICD-11 code
func NewCodeableConcept(code, display string) r5.CodeableConcept {
	coding := NewCoding(code, display)
	return r5.CodeableConcept{
		Coding: []r5.Coding{coding},
		Text:   &display,
	}
}
