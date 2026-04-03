package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

// Bangladesh-specific CodeSystems and ValueSets based on DGHS BD-Core-FHIR-IG

// Identifier Types
const (
	CodeSystemBDIdentifierType = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-identifier-type"

	IdentifierTypeNID  = "NID"         // National ID
	IdentifierTypeBRN  = "BRN"         // Birth Registration Number
	IdentifierTypeUHID = "UHID"        // Unique Health ID
	IdentifierTypeFCN  = "FCN"         // Family Counting Number (Rohingya)
	IdentifierTypePID  = "PROGRESS_ID" // Progress ID (Rohingya)
	IdentifierTypeMRN  = "MRN"         // Medical Record Number
)

// Religions
const (
	CodeSystemBDReligion = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-religions"

	ReligionIslam     = "ISLAM"
	ReligionHinduism  = "HINDUISM"
	ReligionBuddhism  = "BUDDHISM"
	ReligionChristian = "CHRISTIAN"
	ReligionOther     = "OTHER"
	ReligionNone      = "NONE"
)

// Administrative Levels
const (
	CodeSystemBDAdminLevel = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-admin-level"

	AdminLevelDivisionCode = "DIV"
	AdminLevelDistrictCode = "DIS"
	AdminLevelUpazilaCode  = "UPA"
	AdminLevelUnionCode    = "UNI"
	AdminLevelCityCode     = "CITY"
	AdminLevelWardCode     = "WARD"
)

// Healthcare Facility Types
const (
	CodeSystemBDFacilityType = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-facility-type"

	FacilityTypeMedicalCollege       = "MEDICAL_COLLEGE"
	FacilityTypeDistrictHospital     = "DISTRICT_HOSPITAL"
	FacilityTypeUpazilaHealthComplex = "UHC"
	FacilityTypeUnionHealthCenter    = "UHC"
	FacilityTypeCommunityClinic      = "CC"
	FacilityTypeSpecializedHospital  = "SPECIALIZED"
)

// BDValueSet represents a Bangladesh-specific ValueSet
type BDValueSet struct {
	r5.ValueSet
}

// NewBDIdentifierTypeValueSet creates the Bangladesh identifier type ValueSet
func NewBDIdentifierTypeValueSet() *BDValueSet {
	active := "active"
	valueSet := &BDValueSet{
		ValueSet: r5.ValueSet{
			Status:      active,
			Name:        stringPtr("BangladeshIdentifierTypeVS"),
			Title:       stringPtr("Bangladesh Identifier Type Value Set"),
			Description: stringPtr("Types of identifiers used in Bangladesh healthcare system"),
			Compose: &r5.ValueSetCompose{
				Include: []r5.ValueSetComposeInclude{
					{
						System: stringPtr(CodeSystemBDIdentifierType),
						Concept: []r5.ValueSetComposeIncludeConcept{
							{
								Code:    IdentifierTypeNID,
								Display: stringPtr("National ID"),
							},
							{
								Code:    IdentifierTypeBRN,
								Display: stringPtr("Birth Registration Number"),
							},
							{
								Code:    IdentifierTypeUHID,
								Display: stringPtr("Unique Health ID"),
							},
							{
								Code:    IdentifierTypeFCN,
								Display: stringPtr("Family Counting Number"),
							},
							{
								Code:    IdentifierTypePID,
								Display: stringPtr("Progress ID"),
							},
							{
								Code:    IdentifierTypeMRN,
								Display: stringPtr("Medical Record Number"),
							},
						},
					},
				},
			},
		},
	}

	return valueSet
}

// NewBDReligionValueSet creates the Bangladesh religion ValueSet
func NewBDReligionValueSet() *BDValueSet {
	active := "active"
	valueSet := &BDValueSet{
		ValueSet: r5.ValueSet{
			Status:      active,
			Name:        stringPtr("BangladeshReligionVS"),
			Title:       stringPtr("Bangladesh Religion Value Set"),
			Description: stringPtr("Religions recognized in Bangladesh"),
			Compose: &r5.ValueSetCompose{
				Include: []r5.ValueSetComposeInclude{
					{
						System: stringPtr(CodeSystemBDReligion),
						Concept: []r5.ValueSetComposeIncludeConcept{
							{
								Code:    ReligionIslam,
								Display: stringPtr("Islam"),
							},
							{
								Code:    ReligionHinduism,
								Display: stringPtr("Hinduism"),
							},
							{
								Code:    ReligionBuddhism,
								Display: stringPtr("Buddhism"),
							},
							{
								Code:    ReligionChristian,
								Display: stringPtr("Christianity"),
							},
							{
								Code:    ReligionOther,
								Display: stringPtr("Other"),
							},
							{
								Code:    ReligionNone,
								Display: stringPtr("None"),
							},
						},
					},
				},
			},
		},
	}

	return valueSet
}

// NewBDFacilityTypeValueSet creates the Bangladesh facility type ValueSet
func NewBDFacilityTypeValueSet() *BDValueSet {
	active := "active"
	valueSet := &BDValueSet{
		ValueSet: r5.ValueSet{
			Status:      active,
			Name:        stringPtr("BangladeshFacilityTypeVS"),
			Title:       stringPtr("Bangladesh Healthcare Facility Type Value Set"),
			Description: stringPtr("Types of healthcare facilities in Bangladesh"),
			Compose: &r5.ValueSetCompose{
				Include: []r5.ValueSetComposeInclude{
					{
						System: stringPtr(CodeSystemBDFacilityType),
						Concept: []r5.ValueSetComposeIncludeConcept{
							{
								Code:    FacilityTypeMedicalCollege,
								Display: stringPtr("Medical College Hospital"),
							},
							{
								Code:    FacilityTypeDistrictHospital,
								Display: stringPtr("District Hospital"),
							},
							{
								Code:    FacilityTypeUpazilaHealthComplex,
								Display: stringPtr("Upazila Health Complex"),
							},
							{
								Code:    FacilityTypeUnionHealthCenter,
								Display: stringPtr("Union Health Center"),
							},
							{
								Code:    FacilityTypeCommunityClinic,
								Display: stringPtr("Community Clinic"),
							},
							{
								Code:    FacilityTypeSpecializedHospital,
								Display: stringPtr("Specialized Hospital"),
							},
						},
					},
				},
			},
		},
	}

	return valueSet
}

// GetBangladeshValueSets returns all Bangladesh-specific ValueSets
func GetBangladeshValueSets() []*BDValueSet {
	return []*BDValueSet{
		NewBDIdentifierTypeValueSet(),
		NewBDReligionValueSet(),
		NewBDFacilityTypeValueSet(),
	}
}

// ValidateIdentifierType checks if an identifier type is valid for Bangladesh
func ValidateIdentifierType(identifierType string) bool {
	validTypes := []string{
		IdentifierTypeNID,
		IdentifierTypeBRN,
		IdentifierTypeUHID,
		IdentifierTypeFCN,
		IdentifierTypePID,
		IdentifierTypeMRN,
	}

	for _, validType := range validTypes {
		if identifierType == validType {
			return true
		}
	}
	return false
}

// ValidateReligion checks if a religion code is valid for Bangladesh
func ValidateReligion(religion string) bool {
	validReligions := []string{
		ReligionIslam,
		ReligionHinduism,
		ReligionBuddhism,
		ReligionChristian,
		ReligionOther,
		ReligionNone,
	}

	for _, validReligion := range validReligions {
		if religion == validReligion {
			return true
		}
	}
	return false
}

// ValidateFacilityType checks if a facility type is valid for Bangladesh
func ValidateFacilityType(facilityType string) bool {
	validTypes := []string{
		FacilityTypeMedicalCollege,
		FacilityTypeDistrictHospital,
		FacilityTypeUpazilaHealthComplex,
		FacilityTypeUnionHealthCenter,
		FacilityTypeCommunityClinic,
		FacilityTypeSpecializedHospital,
	}

	for _, validType := range validTypes {
		if facilityType == validType {
			return true
		}
	}
	return false
}
