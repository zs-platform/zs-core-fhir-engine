// ============================================================
// ZS-Enhanced-BD-Extensions.fsh
// ZarishSphere Enhanced Bangladesh Extensions
// Extends BD-Core-FHIR-IG extensions with additional capabilities
// Includes refugee-specific extensions, mobile health, surveillance extensions
// ============================================================

// ----- Enhanced GPS Location Extension -----
Extension: ZSGPSLocationExtension
Id: zs-gps-location
Title: "ZarishSphere GPS Location"
Description: """
Enhanced GPS location extension for patient and facility location tracking.
Includes altitude, accuracy, device type, and recording timestamp.
"""
Parent: Extension
* ^url = "https://fhir.zarishsphere.com/StructureDefinition/zs-gps-location"
* ^status = #active
* ^context[+].type = #element
* ^context[=].expression = "Patient"
* ^context[+].type = #element
* ^context[=].expression = "Location"

* extension contains
    latitude 1..1 and
    longitude 1..1 and
    accuracy 0..1 and
    altitude 0..1 and
    recordedAt 1..1 and
    deviceType 0..1 and
    coordinateSystem 0..1

* extension[latitude].url = "latitude"
* extension[latitude].value[x] only decimal
* extension[latitude].valueDecimal ^short = "Latitude coordinate"
* extension[latitude].valueDecimal ^definition = "Geographic latitude in decimal degrees"

* extension[longitude].url = "longitude"
* extension[longitude].value[x] only decimal
* extension[longitude].valueDecimal ^short = "Longitude coordinate"
* extension[longitude].valueDecimal ^definition = "Geographic longitude in decimal degrees"

* extension[accuracy].url = "accuracy"
* extension[accuracy].value[x] only decimal
* extension[accuracy].valueDecimal ^short = "GPS accuracy"
* extension[accuracy].valueDecimal ^definition = "GPS accuracy in meters"

* extension[altitude].url = "altitude"
* extension[altitude].value[x] only decimal
* extension[altitude].valueDecimal ^short = "Altitude"
* extension[altitude].valueDecimal ^definition = "Altitude in meters above sea level"

* extension[recordedAt].url = "recordedAt"
* extension[recordedAt].value[x] only dateTime
* extension[recordedAt].valueDateTime ^short = "Recording timestamp"
* extension[recordedAt].valueDateTime ^definition = "When the GPS coordinates were recorded"

* extension[deviceType].url = "deviceType"
* extension[deviceType].value[x] only CodeableConcept
* extension[deviceType].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/gps-device-types-valueset

* extension[coordinateSystem].url = "coordinateSystem"
* extension[coordinateSystem].value[x] only string
* extension[coordinateSystem].valueString ^short = "Coordinate system"
* extension[coordinateSystem].valueString ^definition = "Coordinate system used (WGS84, etc.)"

// ----- Enhanced Rohingya Status Extension -----
Extension: ZSRohingyaStatusExtension
Id: zs-rohingya-status
Title: "ZarishSphere Rohingya Status"
Description: """
Extension for Rohingya community status and specific healthcare needs.
Includes protection status, vulnerability assessment, and special care requirements.
"""
Parent: Extension
* ^url = "https://fhir.zarishsphere.com/StructureDefinition/zs-rohingya-status"
* ^status = #active
* ^context[+].type = #element
* ^context[=].expression = "Patient"

* extension contains
    status 1..1 and
    protectionStatus 0..1 and
    vulnerabilityLevel 0..1 and
    specialNeeds 0..* and
    campInformation 0..1 and
    familySeparation 0..1 and
    traumaHistory 0..1 and
    languageSupport 0..*

* extension[status].url = "status"
* extension[status].value[x] only CodeableConcept
* extension[status].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/rohingya-status-valueset

* extension[protectionStatus].url = "protectionStatus"
* extension[protectionStatus].value[x] only CodeableConcept
* extension[protectionStatus].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/protection-status-valueset

* extension[vulnerabilityLevel].url = "vulnerabilityLevel"
* extension[vulnerabilityLevel].value[x] only CodeableConcept
* extension[vulnerabilityLevel].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/vulnerability-level-valueset

* extension[specialNeeds].url = "specialNeeds"
* extension[specialNeeds].value[x] only CodeableConcept
* extension[specialNeeds].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/special-needs-valueset

* extension[campInformation].url = "campInformation"
* extension[campInformation].extension contains
    campName 0..1 and
    blockName 0..1 and
    householdNumber 0..1 and
    arrivalDate 0..1 and
    familySize 0..1

* extension[familySeparation].url = "familySeparation"
* extension[familySeparation].value[x] only boolean
* extension[familySeparation].valueBoolean ^short = "Family separation"
* extension[familySeparation].valueBoolean ^definition = "Whether the individual is separated from family members"

// ----- Enhanced Camp Information Extension -----
Extension: ZSCampInfoExtension
Id: zs-camp-info
Title: "ZarishSphere Camp Information"
Description: """
Comprehensive camp information extension for refugee healthcare.
Includes camp demographics, services, and infrastructure information.
"""
Parent: Extension
* ^url = "https://fhir.zarishsphere.com/StructureDefinition/zs-camp-info"
* ^status = #active
* ^context[+].type = #element
* ^context[=].expression = "Patient"
* ^context[+].type = #element
* ^context[=].expression = "Location"

* extension contains
    campName 1..1 and
    campType 1..1 and
    blockName 0..1 and
    householdNumber 0..1 and
    arrivalDate 0..1 and
    familySize 0..1 and
    shelterType 0..1 and
    waterAccess 0..1 and
    sanitationAccess 0..1 and
    foodSecurity 0..1 and
    healthcareAccess 0..1 and
    educationAccess 0..1 and
    protectionServices 0..1 and
    campPopulation 0..1 and
    servicesAvailable 0..*

* extension[campName].url = "campName"
* extension[campName].value[x] only string
* extension[campName].valueString ^short = "Camp name"
* extension[campName].valueString ^definition = "Official name of the refugee camp"

* extension[campType].url = "campType"
* extension[campType].value[x] only CodeableConcept
* extension[campType].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/camp-types-valueset

* extension[shelterType].url = "shelterType"
* extension[shelterType].value[x] only CodeableConcept
* extension[shelterType].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/shelter-types-valueset

* extension[waterAccess].url = "waterAccess"
* extension[waterAccess].value[x] only CodeableConcept
* extension[waterAccess].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/water-access-valueset

* extension[sanitationAccess].url = "sanitationAccess"
* extension[sanitationAccess].value[x] only CodeableConcept
* extension[sanitationAccess].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/sanitation-access-valueset

// ----- Enhanced Consent Management Extension -----
Extension: ZSConsentManagementExtension
Id: zs-consent-management
Title: "ZarishSphere Consent Management"
Description: """
Comprehensive consent management extension for healthcare data sharing.
Includes treatment consent, data sharing consent, and research consent.
"""
Parent: Extension
* ^url = "https://fhir.zarishsphere.com/StructureDefinition/zs-consent-management"
* ^status = #active
* ^context[+].type = #element
* ^context[=].expression = "Patient"
* ^context[+].type = #element
* ^context[=].expression = "Consent"

* extension contains
    treatmentConsent 0..1 and
    dataSharingConsent 0..1 and
    researchConsent 0..1 and
    emergencyConsent 0..1 and
    guardianConsent 0..1 and
    consentHistory 0..*

* extension[treatmentConsent].url = "treatmentConsent"
* extension[treatmentConsent].extension contains
    consentGiven 1..1 and
    consentDate 1..1 and
    recordedBy 1..1 and
    version 0..1 and
    witnessName 0..1 and
    consentMethod 0..1 and
    language 0..1 and
    understandingLevel 0..1

* extension[dataSharingConsent].url = "dataSharingConsent"
* extension[dataSharingConsent].extension contains
    consentGiven 1..1 and
    consentDate 1..1 and
    purposes 0..* and
    restrictions 0..* and
    expiryDate 0..1 and
    withdrawalDate 0..1 and
    sharingPartners 0..*

* extension[researchConsent].url = "researchConsent"
* extension[researchConsent].extension contains
    consentGiven 1..1 and
    consentDate 1..1 and
    studyTypes 0..* and
    withdrawalDate 0..1 and
    guardianConsent 0..1 and
    compensation 0..1 and
    risksExplained 0..1

// ----- Enhanced Communicable Disease Extension -----
Extension: ZSCommunicableDiseaseExtension
Id: zs-communicable-disease
Title: "ZarishSphere Communicable Disease"
Description: """
Enhanced communicable disease surveillance extension.
Includes outbreak detection, reporting requirements, and public health response.
"""
Parent: Extension
* ^url = "https://fhir.zarishsphere.com/StructureDefinition/zs-communicable-disease"
* ^status = #active
* ^context[+].type = #element
* ^context[=].expression = "Condition"
* ^context[+].type = #element
* ^context[=].expression = "Observation"

* extension contains
    diseaseCode 1..1 and
    reportingRequired 1..1 and
    outbreakStatus 0..1 and
    surveillanceCategory 0..1 and
    reportingDeadline 0..1 and
    notifiedTo 0..* and
    caseClassification 0..1 and
    investigationStatus 0..1 and
    controlMeasures 0..* and
    contactTracing 0..1 and
    quarantineRequirements 0..1 and
    laboratoryConfirmation 0..1

* extension[diseaseCode].url = "diseaseCode"
* extension[diseaseCode].value[x] only CodeableConcept
* extension[diseaseCode].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/communicable-diseases-valueset

* extension[reportingRequired].url = "reportingRequired"
* extension[reportingRequired].value[x] only boolean
* extension[reportingRequired].valueBoolean ^short = "Reporting required"
* extension[reportingRequired].valueBoolean ^definition = "Whether this disease requires immediate reporting"

* extension[outbreakStatus].url = "outbreakStatus"
* extension[outbreakStatus].value[x] only CodeableConcept
* extension[outbreakStatus].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/outbreak-status-valueset

* extension[surveillanceCategory].url = "surveillanceCategory"
* extension[surveillanceCategory].value[x] only CodeableConcept
* extension[surveillanceCategory].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/surveillance-categories-valueset

// ----- Enhanced Mobile Health Extension -----
Extension: ZSMobileHealthExtension
Id: zs-mobile-health
Title: "ZarishSphere Mobile Health"
Description: """
Mobile health integration extension for patient engagement and remote monitoring.
Includes app consent, device registration, and data synchronization.
"""
Parent: Extension
* ^url = "https://fhir.zarishsphere.com/StructureDefinition/zs-mobile-health"
* ^status = #active
* ^context[+].type = #element
* ^context[=].expression = "Patient"

* extension contains
    appConsent 0..1 and
    deviceRegistration 0..1 and
    dataSyncEnabled 0..1 and
    lastSyncDate 0..1 and
    preferredLanguage 0..1 and
    literacyLevel 0..1 and
    connectivityAccess 0..1 and
    deviceType 0..1 and
    notificationPreferences 0..* and
    remoteMonitoring 0..1 and
    digitalSkills 0..1

* extension[appConsent].url = "appConsent"
* extension[appConsent].extension contains
    consentGiven 1..1 and
    consentDate 1..1 and
    appVersion 0..1 and
    dataTypes 0..* and
    withdrawalDate 0..1

* extension[deviceRegistration].url = "deviceRegistration"
* extension[deviceRegistration].extension contains
    deviceId 0..1 and
    deviceType 0..1 and
    registrationDate 0..1 and
    lastActiveDate 0..1 and
    appVersion 0..1

* extension[dataSyncEnabled].url = "dataSyncEnabled"
* extension[dataSyncEnabled].value[x] only boolean
* extension[dataSyncEnabled].valueBoolean ^short = "Data sync enabled"
* extension[dataSyncEnabled].valueBoolean ^definition = "Whether data synchronization is enabled for this patient"

// ----- Enhanced Cold Chain Extension -----
Extension: ZSColdChainExtension
Id: zs-cold-chain
Title: "ZarishSphere Cold Chain"
Description: """
Cold chain management extension for vaccine and medication storage.
Includes temperature monitoring, excursion tracking, and quality control.
"""
Parent: Extension
* ^url = "https://fhir.zarishsphere.com/StructureDefinition/zs-cold-chain"
* ^status = #active
* ^context[+].type = #element
* ^context[=].expression = "Immunization"
* ^context[+].type = #element
* ^context[=].expression = "MedicationAdministration"

* extension contains
    storageTemperature 0..1 and
    temperatureExcursions 0..* and
    refrigeratorId 0..1 and
    lastTemperatureCheck 0..1 and
    coldChainStatus 0..1 and
    monitoringDevice 0..1 and
    qualityControl 0..1 and
    transportConditions 0..1 and
    storageLocation 0..1 and
    responsiblePerson 0..1

* extension[storageTemperature].url = "storageTemperature"
* extension[storageTemperature].extension contains
    currentTemp 1..1 and
    minTemp 0..1 and
    maxTemp 0..1 and
    unit 1..1 and
    recordedAt 1..1

* extension[temperatureExcursions].url = "temperatureExcursions"
* extension[temperatureExcursions].extension contains
    excursionStart 1..1 and
    excursionEnd 0..1 and
    minTemp 0..1 and
    maxTemp 0..1 and
    duration 0..1 and
    impact 0..1 and
    correctiveAction 0..1

* extension[coldChainStatus].url = "coldChainStatus"
* extension[coldChainStatus].value[x] only CodeableConcept
* extension[coldChainStatus].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/cold-chain-status-valueset
