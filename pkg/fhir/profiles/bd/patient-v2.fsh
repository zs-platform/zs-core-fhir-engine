Alias: $translation = http://hl7.org/fhir/StructureDefinition/translation

Profile: ZSPatientV2
Id: zs-patient-v2
Parent: Patient
Title: "ZarishSphere Patient Profile v2"
Description: """
Enhanced Patient profile for ZarishSphere with Bangladesh DGHS compliance, 
multilingual support, and comprehensive privacy controls.

Key Features:
- Bilingual name requirements (Bangla + English)
- Bangladesh-specific identifiers (NID, BRN, UHID, UNHCR)
- Enhanced privacy and consent management
- GPS location support
- Rohingya community support
- GDPR/HIPAA compliance extensions
"""

// ----- Name Requirements -----
* name 1..* MS
* name.use 1..1 MS
* name.use = #official (exactly)

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

// ----- Identifier Requirements -----
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
    PHONE 0..1

// National ID
* identifier[NID].system = "http://dghs.gov.bd/identifier/nid"
* identifier[NID].type.coding.code = #NID
* identifier[NID].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[NID].type from BangladeshIdentifierTypeVS (extensible)

// Birth Registration Number
* identifier[BRN].system = "http://dghs.gov.bd/identifier/brn"
* identifier[BRN].type.coding.code = #BRN
* identifier[BRN].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[BRN].type from BangladeshIdentifierTypeVS (extensible)

// Unique Health ID
* identifier[UHID].system = "http://dghs.gov.bd/identifier/uhid"
* identifier[UHID].type.coding.code = #UHID
* identifier[UHID].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[UHID].type from BangladeshIdentifierTypeVS (extensible)

// UNHCR Registration
* identifier[UNHCR].system = "http://unhcr.org/identifier/registration"
* identifier[UNHCR].type.coding.code = #UNHCR
* identifier[UNHCR].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[UNHCR].type from BangladeshIdentifierTypeVS (extensible)

// UNHCR Smart Card
* identifier[SMARTCARD].system = "http://unhcr.org/identifier/smartcard"
* identifier[SMARTCARD].type.coding.code = #SMARTCARD
* identifier[SMARTCARD].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[SMARTCARD].type from BangladeshIdentifierTypeVS (extensible)

// Passport
* identifier[PASSPORT].system = "http://passport.gov.bd/identifier/passport"
* identifier[PASSPORT].type.coding.code = #PASSPORT
* identifier[PASSPORT].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[PASSPORT].type from BangladeshIdentifierTypeVS (extensible)

// Phone Number
* identifier[PHONE].system = "tel"
* identifier[PHONE].type.coding.code = #PHONE
* identifier[PHONE].type.coding.system = "https://fhir.dghs.gov.bd/core/ValueSet/bd-identifier-type-valueset"
* identifier[PHONE].type from BangladeshIdentifierTypeVS (extensible)

// ----- Demographics -----
* gender 1..1 MS
* gender from http://hl7.org/fhir/ValueSet/administrative-gender

* birthDate 1..1 MS
* birthDate ^comment = "If exact date of birth is partially or completely unknown, implementers SHALL populate this element with date of birth information listed on patient's government-issued identification."

* maritalStatus from http://hl7.org/fhir/ValueSet/marital-status

// ----- Bangladesh Extensions -----
// Religion
* extension contains http://hl7.org/fhir/StructureDefinition/patient-religion named religion 0..1
* extension[religion].valueCodeableConcept from https://fhir.dghs.gov.bd/core/ValueSet/bd-religions-valueset

// Nationality
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-nationality named nationality 0..1
* extension[nationality].valueCodeableConcept from https://fhir.dghs.gov.bd/core/ValueSet/bd-nationalities-valueset

// Rohingya Community Support
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-rohingya-status named rohingyaStatus 0..1
* extension[rohingyaStatus].valueCodeableConcept from https://fhir.dghs.gov.bd/core/ValueSet/rohingya-status-valueset

// Camp Information
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-camp-info named campInfo 0..1
* extension[campInfo].extension contains
    campName 0..1 and
    blockName 0..1 and
    campType 0..1

// GPS Location
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-gps-location named gpsLocation 0..1
* extension[gpsLocation].extension contains
    latitude 0..1 and
    longitude 0..1 and
    accuracy 0..1 and
    recordedAt 0..1

// ----- Privacy & Consent -----
// Treatment Consent
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-treatment-consent named treatmentConsent 0..1
* extension[treatmentConsent].extension contains
    consentGiven 0..1 and
    consentDate 0..1 and
    recordedBy 0..1 and
    version 0..1

// Data Sharing Consent
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-data-sharing-consent named dataSharingConsent 0..1
* extension[dataSharingConsent].extension contains
    consentGiven 0..1 and
    consentDate 0..1 and
    purposes 0..* and
    restrictions 0..*

// Research Consent
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-research-consent named researchConsent 0..1
* extension[researchConsent].extension contains
    consentGiven 0..1 and
    consentDate 0..1 and
    studyTypes 0..* and
    withdrawalDate 0..1

// Privacy Classification
* extension contains https://fhir.zarishsphere.com/StructureDefinition/patient-privacy-classification named privacyClassification 0..1
* extension[privacyClassification].valueCodeableConcept from https://fhir.zarishsphere.com/ValueSet/privacy-classification-valueset

// ----- Address -----
* address 1..* MS
* address only BDAddressV2

// ----- Communication -----
* telecom 0..* MS
* telecom ^slicing.discriminator.type = #value
* telecom ^slicing.discriminator.path = "system"
* telecom ^slicing.rules = #open
* telecom contains
    phone 0..1 and
    email 0..1

* telecom[phone].system = "phone"
* telecom[phone].use = #mobile
* telecom[phone].rank 1

* telecom[email].system = "email"
* telecom[email].use = #home
* telecom[email].rank 2

// ----- Contact -----
* contact 0..* MS
* contact.relationship from http://hl7.org/fhir/ValueSet/patient-contactrelationship
* contact.name 1..1 MS
