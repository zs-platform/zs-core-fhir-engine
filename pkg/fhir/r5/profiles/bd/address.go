package bd

import (
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/fhir/r5"
)

// BDAdministrativeDivisions represents Bangladesh's administrative hierarchy
const (
	// Profile URLs
	ProfileBDAddress = "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-address"

	// Administrative level codes
	AdminLevelDivision = "1" // Division
	AdminLevelDistrict = "2" // District/Zila
	AdminLevelUpazila  = "3" // Upazila/Thana
	AdminLevelUnion    = "4" // Union Parishad
	AdminLevelCity     = "5" // City Corporation
	AdminLevelWard     = "6" // Ward
)

// BDAddress represents a Bangladesh-specific address profile
type BDAddress struct {
	r5.Address
}

// NewBDAddress creates a new Bangladesh address
func NewBDAddress() *BDAddress {
	address := &BDAddress{
		Address: r5.Address{},
	}

	// Set Bangladesh as default country
	country := "BD"
	address.Country = &country

	return address
}

// SetAdministrativeLevels sets the administrative division hierarchy
func (a *BDAddress) SetAdministrativeLevels(division, district, upazila string) {
	// Add administrative extensions
	if a.Extension == nil {
		a.Extension = []r5.Extension{}
	}

	// Division extension
	if division != "" {
		divExt := r5.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-division",
		}
		divExt.ValueString = &division
		a.Extension = append(a.Extension, divExt)
	}

	// District extension
	if district != "" {
		distExt := r5.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-district",
		}
		distExt.ValueString = &district
		a.Extension = append(a.Extension, distExt)
	}

	// Upazila extension
	if upazila != "" {
		upazilaExt := r5.Extension{
			URL: "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-upazila",
		}
		upazilaExt.ValueString = &upazila
		a.Extension = append(a.Extension, upazilaExt)
	}
}

// SetCity sets the city name
func (a *BDAddress) SetCity(city string) {
	a.City = stringPtr(city)
}

// SetPostalCode sets the postal code
func (a *BDAddress) SetPostalCode(postalCode string) {
	a.PostalCode = stringPtr(postalCode)
}

// AddLine adds address line
func (a *BDAddress) AddLine(line string) {
	if a.Line == nil {
		a.Line = []string{}
	}
	a.Line = append(a.Line, line)
}

// GetAdministrativeLevels returns the administrative hierarchy
func (a *BDAddress) GetAdministrativeLevels() (division, district, upazila string) {
	if a.Extension == nil {
		return "", "", ""
	}

	for _, ext := range a.Extension {
		switch ext.URL {
		case "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-division":
			if ext.ValueString != nil {
				division = *ext.ValueString
			}
		case "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-district":
			if ext.ValueString != nil {
				district = *ext.ValueString
			}
		case "https://fhir.dghs.gov.bd/core/StructureDefinition/bd-upazila":
			if ext.ValueString != nil {
				upazila = *ext.ValueString
			}
		}
	}

	return division, district, upazila
}

// Bangladesh Administrative Divisions Data
var (
	// Divisions (8 divisions)
	BDDivisions = []string{
		"Barishal", "Chattogram", "Dhaka", "Khulna",
		"Mymensingh", "Rajshahi", "Rangpur", "Sylhet",
	}

	// Major Districts (sample - full list would be 64 districts)
	BDDistricts = map[string][]string{
		"Dhaka":      {"Dhaka", "Faridpur", "Gazipur", "Gopalganj", "Kishoreganj", "Manikganj", "Munshiganj", "Narayanganj", "Narsingdi", "Rajbari", "Shariatpur", "Tangail"},
		"Chattogram": {"Bandarban", "Brahmanbaria", "Chandpur", "Chattogram", "Comilla", "Cox's Bazar", "Feni", "Khagrachhari", "Lakshmipur", "Noakhali", "Rangamati"},
		"Rajshahi":   {"Bogura", "Chapainawabganj", "Joypurhat", "Naogaon", "Natore", "Pabna", "Rajshahi", "Sirajganj"},
		"Khulna":     {"Bagerhat", "Chuadanga", "Jessore", "Jhenaidah", "Khulna", "Kushtia", "Magura", "Meherpur", "Narail", "Satkhira"},
		"Barishal":   {"Barguna", "Barishal", "Bhola", "Jhalokati", "Patuakhali", "Pirojpur"},
		"Sylhet":     {"Habiganj", "Moulvibazar", "Sunamganj", "Sylhet"},
		"Rangpur":    {"Dinajpur", "Gaibandha", "Kurigram", "Lalmonirhat", "Nilphamari", "Panchagarh", "Rangpur", "Thakurgaon"},
		"Mymensingh": {"Jamalpur", "Mymensingh", "Netrokona", "Sherpur"},
	}

	// Rohingya Camp Locations (Cox's Bazar District)
	RohingyaCamps = []string{
		"Kutupalong", "Nayapara", "Leda", "Shamlapur", "Unchiprang",
		"Balukhali", "Thangkhali", "Jamtoli", "Hakimpara", "Camp 25",
	}
)

// IsValidDivision checks if the division is valid
func IsValidDivision(division string) bool {
	for _, div := range BDDivisions {
		if div == division {
			return true
		}
	}
	return false
}

// IsValidDistrict checks if the district is valid for the given division
func IsValidDistrict(division, district string) bool {
	districts, exists := BDDistricts[division]
	if !exists {
		return false
	}

	for _, dist := range districts {
		if dist == district {
			return true
		}
	}
	return false
}

// IsRohingyaCamp checks if the location is a Rohingya camp
func IsRohingyaCamp(location string) bool {
	for _, camp := range RohingyaCamps {
		if camp == location {
			return true
		}
	}
	return false
}

// GetDistrictsByDivision returns all districts in a division
func GetDistrictsByDivision(division string) []string {
	return BDDistricts[division]
}

// GetAllDivisions returns all divisions
func GetAllDivisions() []string {
	return BDDivisions
}

// GetAllRohingyaCamps returns all Rohingya camps
func GetAllRohingyaCamps() []string {
	return RohingyaCamps
}
