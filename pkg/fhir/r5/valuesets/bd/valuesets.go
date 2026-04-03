package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

const (
	SystemBDDivisions = "https://health.zarishsphere.com/fhir/ValueSet/bd-divisions"
	SystemBDDistricts = "https://health.zarishsphere.com/fhir/ValueSet/bd-districts"
)

// Bangladesh Divisions
var Divisions = map[string]string{
	"DH": "Dhaka",
	"CH": "Chattogram",
	"RJ": "Rajshahi",
	"KH": "Khulna",
	"BR": "Barishal",
	"SY": "Sylhet",
	"RG": "Rangpur",
	"MY": "Mymensingh",
}

// GetDivisionCoding returns a FHIR r5.Coding for a Bangladesh division
func GetDivisionCoding(code string) *r5.Coding {
	if display, ok := Divisions[code]; ok {
		system := SystemBDDivisions
		return &r5.Coding{
			System:  &system,
			Code:    &code,
			Display: &display,
		}
	}
	return nil
}

const (
	SystemRohingyaCamps = "https://health.zarishsphere.com/fhir/ValueSet/rohingya-camps"
)

// Rohingya Camps Example List
var RohingyaCamps = map[string]string{
	"C1E": "Camp 1E",
	"C1W": "Camp 1W",
	"C2E": "Camp 2E",
	"C2W": "Camp 2W",
	"C3":  "Camp 3",
	"C4":  "Camp 4",
	"KTP": "Kutupalong RC",
	"NYP": "Nayapara RC",
}

// GetCampCoding returns a FHIR r5.Coding for a Rohingya camp
func GetCampCoding(code string) *r5.Coding {
	if display, ok := RohingyaCamps[code]; ok {
		system := SystemRohingyaCamps
		return &r5.Coding{
			System:  &system,
			Code:    &code,
			Display: &display,
		}
	}
	return nil
}
