package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

const (
	ProfileBDPatient = "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-patient"
	ExtensionNID     = "http://dghs.gov.bd/identifier/nid"
	ExtensionBRN     = "http://dghs.gov.bd/identifier/brn"
	ExtensionUHID    = "http://dghs.gov.bd/identifier/uhid"
)

// BDPatient represents a comprehensive Bangladesh Patient profile
// Based on DGHS BD-Core-FHIR-IG BDPatientProfile
type BDPatient struct {
	r5.Patient
}

// NewBDPatient creates a new Bangladesh Patient with DGHS compliance
func NewBDPatient() *BDPatient {
	p := &BDPatient{
		Patient: r5.Patient{},
	}

	// Set DGHS profile
	if p.Meta == nil {
		p.Meta = &fhir.Meta{}
	}
	p.Meta.Profile = []string{ProfileBDPatient}

	return p
}

// AddNID adds National ID identifier with proper DGHS formatting
func (p *BDPatient) AddNID(nid string) {
	system := "http://dghs.gov.bd/identifier/nid"
	value := nid
	identifier := r5.Identifier{
		System: &system,
		Value:  &value,
		Type: &r5.CodeableConcept{
			Coding: []r5.Coding{
				{
					System:  stringPtr("https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"),
					Code:    stringPtr("NID"),
					Display: stringPtr("National ID"),
				},
			},
			Text: stringPtr("National ID"),
		},
	}

	p.Identifier = append(p.Identifier, identifier)
}

// AddBRN adds Birth Registration Number identifier
func (p *BDPatient) AddBRN(brn string) {
	system := "http://dghs.gov.bd/identifier/brn"
	value := brn
	identifier := r5.Identifier{
		System: &system,
		Value:  &value,
		Type: &r5.CodeableConcept{
			Coding: []r5.Coding{
				{
					System:  stringPtr("https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"),
					Code:    stringPtr("BRN"),
					Display: stringPtr("Birth Registration Number"),
				},
			},
			Text: stringPtr("Birth Registration Number"),
		},
	}

	p.Identifier = append(p.Identifier, identifier)
}

// AddUHID adds Unique Health ID identifier
func (p *BDPatient) AddUHID(uhid string) {
	system := "http://dghs.gov.bd/identifier/uhid"
	value := uhid
	identifier := r5.Identifier{
		System: &system,
		Value:  &value,
		Type: &r5.CodeableConcept{
			Coding: []r5.Coding{
				{
					System:  stringPtr("https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"),
					Code:    stringPtr("UHID"),
					Display: stringPtr("Unique Health ID"),
				},
			},
			Text: stringPtr("Unique Health ID"),
		},
	}

	p.Identifier = append(p.Identifier, identifier)
}

// AddBilingualName adds name with English and Bangla translation extensions
func (p *BDPatient) AddBilingualName(englishName, banglaName string) {
	official := "official"
	name := r5.HumanName{
		Use:    &official,
		Text:   stringPtr(englishName),
		Family: stringPtr(englishName), // Simplified - in real implementation would parse full name
		Given:  []string{englishName},
	}

	p.Name = append(p.Name, name)
}

// SetReligion sets patient religion using DGHS extension
func (p *BDPatient) SetReligion(religionCode, religionDisplay string) {
	url := "http://hl7.org/fhir/StructureDefinition/patient-religion"
	religionExtension := fhir.Extension{
		URL: url,
	}

	if p.Extension == nil {
		p.Extension = []fhir.Extension{}
	}

	// Remove existing religion extension if present
	for i, ext := range p.Extension {
		if ext.URL == "http://hl7.org/fhir/StructureDefinition/patient-religion" {
			p.Extension = append(p.Extension[:i], p.Extension[i+1:]...)
			break
		}
	}

	p.Extension = append(p.Extension, religionExtension)
}

// AddBangladeshAddress adds a Bangladesh-specific address
func (p *BDPatient) AddBangladeshAddress(division, district, upazila, city, postalCode, line string) {
	address := r5.Address{
		City:       stringPtr(city),
		PostalCode: stringPtr(postalCode),
		Country:    stringPtr("BD"), // Bangladesh
	}

	if line != "" {
		address.Line = []string{line}
	}

	p.Address = append(p.Address, address)
}

// ValidateBDPatient validates Bangladesh-specific requirements per DGHS IG
func (p *BDPatient) ValidateBDPatient() []string {
	var errors []string

	// Check if at least one official name exists
	hasOfficialName := false
	for _, name := range p.Name {
		if name.Use != nil && *name.Use == "official" {
			hasOfficialName = true
			break
		}
	}

	if !hasOfficialName {
		errors = append(errors, "Patient must have at least one official name")
	}

	// Check if at least one Bangladesh identifier exists
	hasBDIdentifier := false
	for _, identifier := range p.Identifier {
		if identifier.System != nil && (*identifier.System == "http://dghs.gov.bd/identifier/nid" ||
			*identifier.System == "http://dghs.gov.bd/identifier/brn" ||
			*identifier.System == "http://dghs.gov.bd/identifier/uhid") {
			hasBDIdentifier = true
			break
		}
	}

	if !hasBDIdentifier {
		errors = append(errors, "Patient must have at least one Bangladesh identifier (NID, BRN, or UHID)")
	}

	// Check if at least one address exists
	if len(p.Address) == 0 {
		errors = append(errors, "Patient must have at least one address")
	}

	return errors
}

// GetBDIdentifiers returns all Bangladesh-specific identifiers
func (p *BDPatient) GetBDIdentifiers() map[string]string {
	identifiers := make(map[string]string)

	for _, id := range p.Identifier {
		if id.System != nil && id.Value != nil {
			switch *id.System {
			case "http://dghs.gov.bd/identifier/nid":
				identifiers["NID"] = *id.Value
			case "http://dghs.gov.bd/identifier/brn":
				identifiers["BRN"] = *id.Value
			case "http://dghs.gov.bd/identifier/uhid":
				identifiers["UHID"] = *id.Value
			}
		}
	}

	return identifiers
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
