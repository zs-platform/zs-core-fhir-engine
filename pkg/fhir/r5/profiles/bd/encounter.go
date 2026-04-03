package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

// Bangladesh Encounter Profile Constants
const (
	ProfileBDEncounter = "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-encounter"

	// Encounter Class Values for Bangladesh
	EncounterClassAMB  = "AMB"  // Ambulatory
	EncounterClassEMER = "EMER" // Emergency
	EncounterClassIMP  = "IMP"  // Inpatient
	EncounterClassVR   = "VR"   // Virtual

	// Encounter Type Codes for Bangladesh
	EncounterTypeOPD   = "OPD"   // Outpatient Department
	EncounterTypeIPD   = "IPD"   // Inpatient Department
	EncounterTypeER    = "ER"    // Emergency Room
	EncounterTypeANC   = "ANC"   // Antenatal Care
	EncounterTypePNC   = "PNC"   // Postnatal Care
	EncounterTypeEPI   = "EPI"   // Extended Program on Immunization
	EncounterTypeNCD   = "NCD"   // Non-Communicable Disease
	EncounterTypeCOVID = "COVID" // COVID-19 Care
)

// BDEncounter represents a Bangladesh-specific Encounter profile
type BDEncounter struct {
	r5.Encounter
}

// NewBDEncounter creates a new Bangladesh Encounter
func NewBDEncounter() *BDEncounter {
	encounter := &BDEncounter{
		Encounter: r5.Encounter{},
	}

	// Set profile
	if encounter.Meta == nil {
		encounter.Meta = &fhir.Meta{}
	}
	encounter.Meta.Profile = []string{ProfileBDEncounter}

	return encounter
}

// SetBangladeshEncounterType sets Bangladesh-specific encounter type
func (e *BDEncounter) SetBangladeshEncounterType(encounterType string) {
	// Create coding for Bangladesh encounter type
	coding := r5.Coding{
		System:  stringPtr("https://fhir.dghs.gov.bd/core/CodeSystem/bd-encounter-type"),
		Code:    stringPtr(encounterType),
		Display: stringPtr(getEncounterTypeDisplay(encounterType)),
	}

	typeCodeable := r5.CodeableConcept{
		Coding: []r5.Coding{coding},
		Text:   stringPtr(getEncounterTypeDisplay(encounterType)),
	}

	e.Type = []r5.CodeableConcept{typeCodeable}

	// Set class based on encounter type
	class := getEncounterClass(encounterType)
	classCodeable := r5.CodeableConcept{
		Coding: []r5.Coding{
			{
				System:  stringPtr("http://terminology.hl7.org/CodeSystem/v3-ActEncounterCode"),
				Code:    stringPtr(class),
				Display: stringPtr(getEncounterClassDisplay(class)),
			},
		},
	}
	e.Class = []r5.CodeableConcept{classCodeable}
}

// SetFacilityReference sets reference to Bangladesh healthcare facility
func (e *BDEncounter) SetFacilityReference(facilityID, facilityName, facilityType string) {
	reference := r5.Reference{
		Reference: stringPtr("Organization/" + facilityID),
		Type:      stringPtr("Organization"),
		Display:   stringPtr(facilityName),
	}

	e.ServiceProvider = &reference

	// Add facility type extension
	if e.Extension == nil {
		e.Extension = []fhir.Extension{}
	}

	facilityTypeExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-facility-type",
	}
	facilityTypeExt.ValueString = &facilityType
	e.Extension = append(e.Extension, facilityTypeExt)
}

// SetAdministrativeInfo sets administrative division information
func (e *BDEncounter) SetAdministrativeInfo(division, district, upazila string) {
	if e.Extension == nil {
		e.Extension = []fhir.Extension{}
	}

	// Division extension
	if division != "" {
		divExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-division",
		}
		divExt.ValueString = &division
		e.Extension = append(e.Extension, divExt)
	}

	// District extension
	if district != "" {
		distExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-district",
		}
		distExt.ValueString = &district
		e.Extension = append(e.Extension, distExt)
	}

	// Upazila extension
	if upazila != "" {
		upazilaExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-upazila",
		}
		upazilaExt.ValueString = &upazila
		e.Extension = append(e.Extension, upazilaExt)
	}
}

// SetEmergencyInfo sets emergency-related information
func (e *BDEncounter) SetEmergencyInfo(isEmergency bool, emergencyReason string) {
	if isEmergency {
		e.Priority = &r5.CodeableConcept{
			Coding: []r5.Coding{
				{
					System:  stringPtr("http://terminology.hl7.org/CodeSystem/v3-ActPriority"),
					Code:    stringPtr("STAT"),
					Display: stringPtr("Emergency"),
				},
			},
		}
	}

	if emergencyReason != "" {
		if e.Reason == nil {
			e.Reason = []r5.EncounterReason{}
		}

		reason := r5.EncounterReason{
			Use: []r5.CodeableConcept{
				{
					Coding: []r5.Coding{
						{
							System:  stringPtr("https://fhir.dghs.gov.bd/core/CodeSystem/bd-emergency-reason"),
							Code:    stringPtr("EMERGENCY_REASON"),
							Display: stringPtr(emergencyReason),
						},
					},
					Text: stringPtr(emergencyReason),
				},
			},
		}

		e.Reason = append(e.Reason, reason)
	}
}

// ValidateBDEncounter validates Bangladesh-specific requirements
func (e *BDEncounter) ValidateBDEncounter() []string {
	var errors []string

	// Check if encounter type is set
	if e.Type == nil || len(e.Type) == 0 {
		errors = append(errors, "Bangladesh encounter must have an encounter type")
	}

	// Validate encounter type
	if e.Type != nil && len(e.Type) > 0 && e.Type[0].Coding != nil && len(e.Type[0].Coding) > 0 {
		encType := e.Type[0].Coding[0].Code
		if encType != nil && !isValidBangladeshEncounterType(*encType) {
			errors = append(errors, "Invalid Bangladesh encounter type: "+*encType)
		}
	}

	// Check if class is set
	if e.Class == nil || len(e.Class) == 0 {
		errors = append(errors, "Bangladesh encounter must have a class")
	}

	// Validate administrative information
	hasDistrict := false

	if e.Extension != nil {
		for _, ext := range e.Extension {
			switch ext.URL {
			case "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-district":
				hasDistrict = true
			}
		}
	}

	// Require at least district level
	if !hasDistrict {
		errors = append(errors, "Bangladesh encounter must specify district")
	}

	return errors
}

// Helper functions

func getEncounterTypeDisplay(encounterType string) string {
	switch encounterType {
	case EncounterTypeOPD:
		return "Outpatient Department"
	case EncounterTypeIPD:
		return "Inpatient Department"
	case EncounterTypeER:
		return "Emergency Room"
	case EncounterTypeANC:
		return "Antenatal Care"
	case EncounterTypePNC:
		return "Postnatal Care"
	case EncounterTypeEPI:
		return "Extended Program on Immunization"
	case EncounterTypeNCD:
		return "Non-Communicable Disease"
	case EncounterTypeCOVID:
		return "COVID-19 Care"
	default:
		return "Unknown"
	}
}

func getEncounterClass(encounterType string) string {
	switch encounterType {
	case EncounterTypeOPD, EncounterTypeANC, EncounterTypePNC, EncounterTypeEPI, EncounterTypeNCD:
		return EncounterClassAMB
	case EncounterTypeIPD:
		return EncounterClassIMP
	case EncounterTypeER, EncounterTypeCOVID:
		return EncounterClassEMER
	default:
		return EncounterClassAMB
	}
}

func getEncounterClassDisplay(class string) string {
	switch class {
	case EncounterClassAMB:
		return "Ambulatory"
	case EncounterClassEMER:
		return "Emergency"
	case EncounterClassIMP:
		return "Inpatient"
	case EncounterClassVR:
		return "Virtual"
	default:
		return "Unknown"
	}
}

func isValidBangladeshEncounterType(encounterType string) bool {
	validTypes := []string{
		EncounterTypeOPD,
		EncounterTypeIPD,
		EncounterTypeER,
		EncounterTypeANC,
		EncounterTypePNC,
		EncounterTypeEPI,
		EncounterTypeNCD,
		EncounterTypeCOVID,
	}

	for _, validType := range validTypes {
		if encounterType == validType {
			return true
		}
	}
	return false
}
