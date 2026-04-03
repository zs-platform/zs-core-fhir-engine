package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

// Bangladesh Organization Profile Constants
const (
	ProfileBDOrganization = "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-organization"

	// Organization Types for Bangladesh
	OrganizationTypeDGHS = "DGHS" // Directorate General of Health Services
	OrganizationTypeDMC  = "DMC"  // District/Medical College Hospital
	OrganizationTypeUHC  = "UHC"  // Upazila Health Complex
	OrganizationTypeUHCP = "UHCP" // Union Health and Family Welfare Center
	OrganizationTypeCC   = "CC"   // Community Clinic
	OrganizationTypeNGO  = "NGO"  // Non-Governmental Organization
	OrganizationTypePVT  = "PVT"  // Private Hospital/Clinic
	OrganizationTypeUN   = "UN"   // United Nations Agency
	OrganizationTypeGOVT = "GOVT" // Other Government Organization

	// Facility Ownership Types
	OwnershipTypePublic  = "PUBLIC"  // Public/Government
	OwnershipTypePrivate = "PRIVATE" // Private
	OwnershipTypeNGO     = "NGO"     // Non-Governmental
	OwnershipTypeJoint   = "JOINT"   // Public-Private Partnership
)

// BDOrganization represents a Bangladesh-specific Organization profile
type BDOrganization struct {
	r5.Organization
}

// NewBDOrganization creates a new Bangladesh Organization
func NewBDOrganization() *BDOrganization {
	org := &BDOrganization{
		Organization: r5.Organization{},
	}

	// Set profile
	if org.Meta == nil {
		org.Meta = &fhir.Meta{}
	}
	org.Meta.Profile = []string{ProfileBDOrganization}

	return org
}

// SetBangladeshOrganizationType sets Bangladesh-specific organization type
func (o *BDOrganization) SetBangladeshOrganizationType(orgType string) {
	// Create coding for Bangladesh organization type
	coding := r5.Coding{
		System:  stringPtr("https://fhir.dghs.gov.bd/core/CodeSystem/bd-organization-type"),
		Code:    stringPtr(orgType),
		Display: stringPtr(getOrganizationTypeDisplay(orgType)),
	}

	typeCodeable := r5.CodeableConcept{
		Coding: []r5.Coding{coding},
		Text:   stringPtr(getOrganizationTypeDisplay(orgType)),
	}

	o.Type = []r5.CodeableConcept{typeCodeable}
}

// SetOwnership sets ownership information
func (o *BDOrganization) SetOwnership(ownership string) {
	// Add ownership extension
	if o.Extension == nil {
		o.Extension = []fhir.Extension{}
	}

	ownershipExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-ownership",
	}
	ownershipExt.ValueString = &ownership
	o.Extension = append(o.Extension, ownershipExt)
}

// SetAdministrativeInfo sets administrative division information
func (o *BDOrganization) SetAdministrativeInfo(division, district, upazila string) {
	if o.Extension == nil {
		o.Extension = []fhir.Extension{}
	}

	// Division extension
	if division != "" {
		divExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-division",
		}
		divExt.ValueString = &division
		o.Extension = append(o.Extension, divExt)
	}

	// District extension
	if district != "" {
		distExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-district",
		}
		distExt.ValueString = &district
		o.Extension = append(o.Extension, distExt)
	}

	// Upazila extension
	if upazila != "" {
		upazilaExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-upazila",
		}
		upazilaExt.ValueString = &upazila
		o.Extension = append(o.Extension, upazilaExt)
	}
}

// SetFacilityLevel sets healthcare facility level
func (o *BDOrganization) SetFacilityLevel(level string) {
	// Add facility level extension
	if o.Extension == nil {
		o.Extension = []fhir.Extension{}
	}

	levelExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-facility-level",
	}
	levelExt.ValueString = &level
	o.Extension = append(o.Extension, levelExt)
}

// SetBedCapacity sets bed capacity information
func (o *BDOrganization) SetBedCapacity(totalBeds, availableBeds int) {
	// Add bed capacity extension
	if o.Extension == nil {
		o.Extension = []fhir.Extension{}
	}

	// Total beds extension
	totalBedsExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-total-beds",
	}
	totalBedsInt := totalBeds
	totalBedsExt.ValueInteger = &totalBedsInt
	o.Extension = append(o.Extension, totalBedsExt)

	// Available beds extension
	availableBedsExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-available-beds",
	}
	availableBedsInt := availableBeds
	availableBedsExt.ValueInteger = &availableBedsInt
	o.Extension = append(o.Extension, availableBedsExt)
}

// ValidateBDOrganization validates Bangladesh-specific requirements
func (o *BDOrganization) ValidateBDOrganization() []string {
	var errors []string

	// Check if organization type is set
	if o.Type == nil || len(o.Type) == 0 {
		errors = append(errors, "Bangladesh organization must have an organization type")
	}

	// Validate organization type
	if o.Type != nil && len(o.Type) > 0 && o.Type[0].Coding != nil && len(o.Type[0].Coding) > 0 {
		orgType := o.Type[0].Coding[0].Code
		if orgType != nil && !isValidBangladeshOrganizationType(*orgType) {
			errors = append(errors, "Invalid Bangladesh organization type: "+*orgType)
		}
	}

	// Check if name is set
	if o.Name == nil || *o.Name == "" {
		errors = append(errors, "Bangladesh organization must have a name")
	}

	return errors
}

// Helper functions

func getOrganizationTypeDisplay(orgType string) string {
	switch orgType {
	case OrganizationTypeDGHS:
		return "Directorate General of Health Services"
	case OrganizationTypeDMC:
		return "District/Medical College Hospital"
	case OrganizationTypeUHC:
		return "Upazila Health Complex"
	case OrganizationTypeUHCP:
		return "Union Health and Family Welfare Center"
	case OrganizationTypeCC:
		return "Community Clinic"
	case OrganizationTypeNGO:
		return "Non-Governmental Organization"
	case OrganizationTypePVT:
		return "Private Hospital/Clinic"
	case OrganizationTypeUN:
		return "United Nations Agency"
	case OrganizationTypeGOVT:
		return "Other Government Organization"
	default:
		return "Unknown"
	}
}

func isValidBangladeshOrganizationType(orgType string) bool {
	validTypes := []string{
		OrganizationTypeDGHS,
		OrganizationTypeDMC,
		OrganizationTypeUHC,
		OrganizationTypeUHCP,
		OrganizationTypeCC,
		OrganizationTypeNGO,
		OrganizationTypePVT,
		OrganizationTypeUN,
		OrganizationTypeGOVT,
	}

	for _, validType := range validTypes {
		if orgType == validType {
			return true
		}
	}
	return false
}
