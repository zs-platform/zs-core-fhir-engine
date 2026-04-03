package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

// Bangladesh Practitioner Profile Constants
const (
	ProfileBDPractitioner = "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-practitioner"

	// Practitioner Qualifications for Bangladesh
	QualificationMBBS      = "MBBS"      // Bachelor of Medicine, Bachelor of Surgery
	QualificationBDS       = "BDS"       // Bachelor of Dental Surgery
	QualificationMD        = "MD"        // Doctor of Medicine
	QualificationMS        = "MS"        // Master of Surgery
	QualificationMPH       = "MPH"       // Master of Public Health
	QualificationPhD       = "PhD"       // Doctor of Philosophy
	QualificationNursing   = "NURSING"   // Nursing
	QualificationMidwifery = "MIDWIFERY" // Midwifery
	QualificationOther     = "OTHER"     // Other

	// Registration Bodies
	RegBodyBMDC = "BMDC" // Bangladesh Medical and Dental Council
	RegBodyBNS  = "BNS"  // Bangladesh Nursing Council
	RegBodyBSMM = "BSMM" // Bangladesh State Medical and Dental Council
)

// BDPractitioner represents a Bangladesh-specific Practitioner profile
type BDPractitioner struct {
	r5.Practitioner
}

// NewBDPractitioner creates a new Bangladesh Practitioner
func NewBDPractitioner() *BDPractitioner {
	practitioner := &BDPractitioner{
		Practitioner: r5.Practitioner{},
	}

	// Set profile
	if practitioner.Meta == nil {
		practitioner.Meta = &fhir.Meta{}
	}
	practitioner.Meta.Profile = []string{ProfileBDPractitioner}

	return practitioner
}

// SetBangladeshQualification sets Bangladesh-specific qualification
func (p *BDPractitioner) SetBangladeshQualification(qualification, institution, year string) {
	// Create qualification extension
	if p.Extension == nil {
		p.Extension = []fhir.Extension{}
	}

	qualificationExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-qualification",
	}

	// Build qualification details
	qualificationDetails := qualification + " from " + institution + " (" + year + ")"
	qualificationExt.ValueString = &qualificationDetails

	p.Extension = append(p.Extension, qualificationExt)
}

// SetRegistrationInfo sets registration information
func (p *BDPractitioner) SetRegistrationInfo(regNumber, regBody string) {
	// Create registration extension
	if p.Extension == nil {
		p.Extension = []fhir.Extension{}
	}

	registrationExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-registration",
	}

	// Build registration details
	registration := regNumber + " - " + regBody
	registrationExt.ValueString = &registration

	p.Extension = append(p.Extension, registrationExt)
}

// SetSpecialization sets practitioner specialization
func (p *BDPractitioner) SetSpecialization(specialization string) {
	// Create specialization extension
	if p.Extension == nil {
		p.Extension = []fhir.Extension{}
	}

	specializationExt := fhir.Extension{
		URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-specialization",
	}
	specializationExt.ValueString = &specialization

	p.Extension = append(p.Extension, specializationExt)
}

// SetAdministrativeInfo sets administrative division information
func (p *BDPractitioner) SetAdministrativeInfo(division, district, upazila string) {
	if p.Extension == nil {
		p.Extension = []fhir.Extension{}
	}

	// Division extension
	if division != "" {
		divExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-division",
		}
		divExt.ValueString = &division
		p.Extension = append(p.Extension, divExt)
	}

	// District extension
	if district != "" {
		distExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-district",
		}
		distExt.ValueString = &district
		p.Extension = append(p.Extension, distExt)
	}

	// Upazila extension
	if upazila != "" {
		upazilaExt := fhir.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-upazila",
		}
		upazilaExt.ValueString = &upazila
		p.Extension = append(p.Extension, upazilaExt)
	}
}

// ValidateBDPractitioner validates Bangladesh-specific requirements
func (p *BDPractitioner) ValidateBDPractitioner() []string {
	var errors []string

	// Check if name is set
	if p.Name == nil || len(p.Name) == 0 {
		errors = append(errors, "Bangladesh practitioner must have a name")
	}

	// Check if at least one identifier is set
	if p.Identifier == nil || len(p.Identifier) == 0 {
		errors = append(errors, "Bangladesh practitioner must have at least one identifier")
	}

	return errors
}

// Helper functions

func getQualificationDisplay(qualification string) string {
	switch qualification {
	case QualificationMBBS:
		return "Bachelor of Medicine, Bachelor of Surgery"
	case QualificationBDS:
		return "Bachelor of Dental Surgery"
	case QualificationMD:
		return "Doctor of Medicine"
	case QualificationMS:
		return "Master of Surgery"
	case QualificationMPH:
		return "Master of Public Health"
	case QualificationPhD:
		return "Doctor of Philosophy"
	case QualificationNursing:
		return "Nursing"
	case QualificationMidwifery:
		return "Midwifery"
	case QualificationOther:
		return "Other"
	default:
		return "Unknown"
	}
}

func getRegBodyDisplay(regBody string) string {
	switch regBody {
	case RegBodyBMDC:
		return "Bangladesh Medical and Dental Council"
	case RegBodyBNS:
		return "Bangladesh Nursing Council"
	case RegBodyBSMM:
		return "Bangladesh State Medical and Dental Council"
	default:
		return "Unknown"
	}
}

func isValidBangladeshQualification(qualification string) bool {
	validQuals := []string{
		QualificationMBBS,
		QualificationBDS,
		QualificationMD,
		QualificationMS,
		QualificationMPH,
		QualificationPhD,
		QualificationNursing,
		QualificationMidwifery,
		QualificationOther,
	}

	for _, validQual := range validQuals {
		if qualification == validQual {
			return true
		}
	}
	return false
}

func isValidBangladeshRegBody(regBody string) bool {
	validBodies := []string{
		RegBodyBMDC,
		RegBodyBNS,
		RegBodyBSMM,
	}

	for _, validBody := range validBodies {
		if regBody == validBody {
			return true
		}
	}
	return false
}
