package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

const (
	// Profile URLs
	ProfileRohingyaPatient = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-patient"

	// Extension URLs for Identifiers
	ExtensionFCN        = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-fcn"
	ExtensionProgressID = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-progress-id"
	ExtensionMRN        = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-mrn"

	// Extension URLs for Location/Shelter
	ExtensionShelterNumber = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-shelter"
	ExtensionCamp          = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-camp"
	ExtensionBlock         = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-block"
	ExtensionSubBlock      = "https://health.zarishsphere.com/fhir/StructureDefinition/rohingya-sub-block"
)

// RohingyaPatient represents a r5.Patient resource localized for the Rohingya Response
type RohingyaPatient struct {
	r5.Patient
}

// NewRohingyaPatient creates a new localized r5.Patient for the Rohingya response
func NewRohingyaPatient() *RohingyaPatient {
	p := &RohingyaPatient{}
	profile := ProfileRohingyaPatient
	if p.Meta == nil {
		p.Meta = &fhir.Meta{}
	}
	p.Meta.Profile = []string{profile}
	return p
}

// AddRohingyaIdentifiers adds FCN, Progress ID, and MRN to the patient
func (p *RohingyaPatient) AddRohingyaIdentifiers(fcn, progressID, mrn string) {
	// Add FCN
	urlFCN := ExtensionFCN
	p.Extension = append(p.Extension, fhir.Extension{
		URL:         urlFCN,
		ValueString: &fcn,
	})

	// Add Progress ID
	urlPID := ExtensionProgressID
	p.Extension = append(p.Extension, fhir.Extension{
		URL:         urlPID,
		ValueString: &progressID,
	})

	// Add MRN
	urlMRN := ExtensionMRN
	p.Extension = append(p.Extension, fhir.Extension{
		URL:         urlMRN,
		ValueString: &mrn,
	})
}

// SetShelterLocation sets the detailed camp and shelter information
func (p *RohingyaPatient) SetShelterLocation(camp, block, subBlock, shelter string) {
	// We add these as extensions to the r5.Address or directly to the r5.Patient
	// For simplicity and direct access, we add them to the r5.Patient extensions

	urlCamp := ExtensionCamp
	p.Extension = append(p.Extension, fhir.Extension{URL: urlCamp, ValueString: &camp})

	urlBlock := ExtensionBlock
	p.Extension = append(p.Extension, fhir.Extension{URL: urlBlock, ValueString: &block})

	urlSubBlock := ExtensionSubBlock
	p.Extension = append(p.Extension, fhir.Extension{URL: urlSubBlock, ValueString: &subBlock})

	urlShelter := ExtensionShelterNumber
	p.Extension = append(p.Extension, fhir.Extension{URL: urlShelter, ValueString: &shelter})
}
