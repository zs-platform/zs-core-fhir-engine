// ============================================================
// ZS-Enhanced-BD-Patient.fsh
// ZarishSphere Enhanced Bangladesh Patient Profile
// Extends BD-Core-FHIR-IG with additional ZarishSphere capabilities
// Includes Rohingya support, GPS tracking, comprehensive consent management
// ============================================================

Alias: $translation = http://hl7.org/fhir/StructureDefinition/translation

Profile: ZSEnhancedBDPatientProfile
Parent: Patient
Id: zs-enhanced-bd-patient
Title: "ZarishSphere Enhanced Bangladesh Patient Profile"
Description: """
Enhanced Patient profile for ZarishSphere building on Bangladesh Core FHIR IG.
- All BD-Core requirements maintained (NID, BRN, UHID, bilingual names)
- Additional ZarishSphere extensions: GPS location, Rohingya support, comprehensive consent
- Privacy and security enhancements for refugee populations
- Mobile health integration capabilities
- Multi-language support extended to include Myanmar, Urdu, Hindi, Thai
"""

// ----- BD Core Requirements (Maintained) -----
// Require exactly one HumanName with bilingual support
* name 1..1 MS
* name.use 1..1
* name.use = #official (exactly)

// Require a text element with translation extensions
* name.text 1..1 MS
* name.text.extension 2..* MS
* name.text.extension contains
    $translation named nameEn 1..1 MS and
    $translation named nameBn 1..1 MS

// Constraints on English name
* name.text.extension[nameEn].extension[lang].valueCode = #en (exactly)
* name.text.extension[nameEn].extension[content] 1..1 MS

// Constraints on Bangla name
* name.text.extension[nameBn].extension[lang].valueCode = #bn (exactly)
* name.text.extension[nameBn].extension[content] 1..1 MS

// ----- Bangladesh Identifiers (Enhanced) -----
* identifier 1..*
* identifier ^slicing.discriminator.type = #value
* identifier ^slicing.discriminator.path = "system"
* identifier ^slicing.rules = #open
* identifier contains
    NID 0..1 and
    BRN 0..1 and
    UHID 0..1 and
    UNHCR 0..1 and
    SMARTCARD 0..1 and
    PASSPORT 0..1 and
    PHONE 0..1 and
    ROHINGYA_ID 0..1

// BD Core identifiers maintained
* identifier[NID].system = "http://dghs.gov.bd/identifier/nid"
* identifier[NID].type.coding.code = #NID
* identifier[NID].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[NID].type from BangladeshIdentifierTypeVS (extensible)

* identifier[BRN].system = "http://dghs.gov.bd/identifier/brn"
* identifier[BRN].type.coding.code = #BRN
* identifier[BRN].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[BRN].type from BangladeshIdentifierTypeVS (extensible)

* identifier[UHID].system = "http://dghs.gov.bd/identifier/uhid"
* identifier[UHID].type.coding.code = #UHID
* identifier[UHID].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[UHID].type from BangladeshIdentifierTypeVS (extensible)

// ----- ZarishSphere Enhanced Identifiers -----
* identifier[UNHCR].system = "http://unhcr.org/identifier/registration"
* identifier[UNHCR].type.coding.code = #UNHCR
* identifier[UNHCR].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[UNHCR].type from BangladeshIdentifierTypeVS (extensible)

* identifier[SMARTCARD].system = "http://unhcr.org/identifier/smartcard"
* identifier[SMARTCARD].type.coding.code = #SMARTCARD
* identifier[SMARTCARD].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[SMARTCARD].type from BangladeshIdentifierTypeVS (extensible)

* identifier[PASSPORT].system = "http://passport.gov.bd/identifier/passport"
* identifier[PASSPORT].type.coding.code = #PASSPORT
* identifier[PASSPORT].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[PASSPORT].type from BangladeshIdentifierTypeVS (extensible)

* identifier[PHONE].system = "tel"
* identifier[PHONE].type.coding.code = #PHONE
* identifier[PHONE].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[PHONE].type from BangladeshIdentifierTypeVS (extensible)

* identifier[ROHINGYA_ID].system = "http://zarishsphere.com/identifier/rohingya"
* identifier[ROHINGYA_ID].type.coding.code = #ROHINGYA_ID
* identifier[ROHINGYA_ID].type.coding.system = "https://fhir.zarishsphere.com/ValueSet/rohingya-identifier-types-valueset"

// ----- BD Core Extensions (Maintained) -----
* birthDate ^comment = "If exact date of birth is partially or completely unknown, Implementers SHALL populate this element with the date of birth information listed on the patient's government-issued identification."

* maritalStatus from http://hl7.org/fhir/ValueSet/marital-status
* deceased[x] only dateTime

// Religion using standard HL7 extension (maintained)
* extension contains http://hl7.org/fhir/StructureDefinition/patient-religion named religion 0..1
* extension[religion].valueCodeableConcept from https://fhir.dghs.gov.bd/core/ValueSet/bd-religions-valueset

* address 1..* MS
* address only BDAddress

// ----- ZarishSphere Enhanced Extensions -----

// GPS Location Support
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-gps-location named gpsLocation 0..1
* extension[gpsLocation].extension contains
    latitude 0..1 and
    longitude 0..1 and
    accuracy 0..1 and
    altitude 0..1 and
    recordedAt 0..1 and
    deviceType 0..1

// Rohingya Community Support
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-rohingya-status named rohingyaStatus 0..1
* extension[rohingyaStatus].valueCodeableConcept from https://fhir.dghs.gov.bd/core/ValueSet/rohingya-status-valueset

// Camp Information
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-camp-info named campInfo 0..1
* extension[campInfo].extension contains
    campName 0..1 and
    blockName 0..1 and
    campType 0..1 and
    householdNumber 0..1 and
    arrivalDate 0..1 and
    familySize 0..1

// Comprehensive Consent Management
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-treatment-consent named treatmentConsent 0..1
* extension[treatmentConsent].extension contains
    consentGiven 0..1 and
    consentDate 0..1 and
    recordedBy 0..1 and
    version 0..1 and
    witnessName 0..1 and
    consentMethod 0..1

* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-data-sharing-consent named dataSharingConsent 0..1
* extension[dataSharingConsent].extension contains
    consentGiven 0..1 and
    consentDate 0..1 and
    purposes 0..* and
    restrictions 0..* and
    expiryDate 0..1 and
    withdrawalDate 0..1

* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-research-consent named researchConsent 0..1
* extension[researchConsent].extension contains
    consentGiven 0..1 and
    consentDate 0..1 and
    studyTypes 0..* and
    withdrawalDate 0..1 and
    guardianConsent 0..1

// Privacy Classification
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-privacy-classification named privacyClassification 0..1
* extension[privacyClassification].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/privacy-classification-valueset

// Mobile Health Integration
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-mobile-health named mobileHealth 0..1
* extension[mobileHealth].extension contains
    appConsent 0..1 and
    deviceRegistration 0..1 and
    dataSyncEnabled 0..1 and
    lastSyncDate 0..1

// Emergency Contact Enhancement
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-emergency-contact named emergencyContact 0..*
* extension[emergencyContact].extension contains
    contactName 0..1 and
    relationship 0..1 and
    phone 0..1 and
    priority 0..1

// ----- Enhanced Multi-language Support -----
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-preferred-languages named preferredLanguages 0..*
* extension[preferredLanguages].extension contains
    language 0..1 and
    proficiency 0..1 and
    preference 0..1

// ----- Healthcare Access Information -----
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-healthcare-access named healthcareAccess 0..1
* extension[healthcareAccess].extension contains
    primaryFacility 0..1 and
    insuranceStatus 0..1 and
    transportationAccess 0..1 and
    digitalLiteracy 0..1

// ----- Vulnerable Population Support -----
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-vulnerability named vulnerability 0..*
* extension[vulnerability].extension contains
    vulnerabilityType 0..1 and
    severity 0..1 and
    supportNeeds 0..* and
    referralRequired 0..1
