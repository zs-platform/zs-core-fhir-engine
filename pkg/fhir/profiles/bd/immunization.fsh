// ============================================================
// ZS-Enhanced-BD-Immunization.fsh
// ZarishSphere Enhanced Bangladesh Immunization Profile
// Extends BD-Core-FHIR-IG with comprehensive vaccine management
// Includes refugee immunization tracking, cold chain management, adverse event monitoring
// ============================================================

Profile: ZSEnhancedBDImmunizationProfile
Parent: Immunization
Id: zs-enhanced-bd-immunization
Title: "ZarishSphere Enhanced Bangladesh Immunization Profile"
Description: """
Enhanced Immunization profile for ZarishSphere building on Bangladesh Core FHIR IG.
- All BD-Core immunization requirements maintained with EPI vaccine support
- Additional ZarishSphere extensions: refugee immunization tracking, cold chain management
- Enhanced adverse event monitoring and reporting
- Integration with national immunization program
- Support for special immunization campaigns and outbreak response
"""

// ----- BD Core Immunization Requirements (Maintained) -----
* identifier 1..*
* identifier ^short = "Unique identifier"
* identifier ^definition = "Unique identifier for the vaccination event"

* reasonReference 0..*
* reasonReference only Reference(Condition or Observation or DiagnosticReport)

* vaccineCode 1..1
* vaccineCode from BDVaccineVS (required)

* manufacturer 0..1
* manufacturer ^short = "Manufacturer"
* manufacturer ^definition = "Vaccine manufacturer"
* manufacturer only Reference(ZSEnhancedBDOrganizationProfile)

* lotNumber 0..1
* lotNumber ^short = "Vaccine Lot Number"
* lotNumber ^definition = "Vaccine lot or batch number"

* expirationDate 0..1
* expirationDate ^short = "Expiration Date"
* expirationDate ^definition = "Expiration date of vaccine lot"

* patient 1..1
* patient ^definition = "The patient receiving the vaccine"
* patient only Reference(ZSEnhancedBDPatientProfile)

* encounter 1..1
* encounter ^definition = "Encounter during which vaccine was administered"
* encounter only Reference(ZSEnhancedBDEncounterProfile)

* occurrence[x] 1..1

* location 0..1
* location ^definition = "Location where vaccine was administered"
* location only Reference(ZSEnhancedBDLocationProfile)

* site 0..1
* site ^definition = "Body site of administration"
* site from BDImmunizationSiteVS

* route 0..1
* route ^definition = "Route of administration"
* route from BDImmunizationRouteVS

* doseQuantity 0..1
* doseQuantity ^definition = "Amount of vaccine administered"
* doseQuantity.system = "http://unitsofmeasure.org"

* performer 0..*
* performer ^definition = "Individual who performed the immunization"
* performer.actor only Reference(ZSEnhancedBDPractitionerProfile)

* reaction 0..*
* reaction ^definition = "Adverse reaction following immunization"
* reaction.detail only Reference(ZSEnhancedBDObservationProfile)

// ----- ZarishSphere Enhanced Extensions -----

// Cold Chain Management
* extension contains https://fhir.zarishsphere.com/StructureDefinition/immunization-cold-chain named coldChain 0..1
* extension[coldChain].extension contains
    storageTemperature 0..1 and
    temperatureExcursions 0..* and
    refrigeratorId 0..1 and
    lastTemperatureCheck 0..1 and
    coldChainStatus 0..1

// Campaign Information
* extension contains https://fhir.zarishsphere.com/StructureDefinition/immunization-campaign named campaign 0..1
* extension[campaign].extension contains
    campaignId 0..1 and
    campaignName 0..1 and
    campaignType 0..1 and
    targetPopulation 0..1 and
    campaignDates 0..1

// Refugee Immunization Tracking
* extension contains https://fhir.zarishsphere.com/StructureDefinition/immunization-refugee-context named refugeeContext 0..1
* extension[refugeeContext].extension contains
    campBased 0..1 and
    specialSchedule 0..1 and
    documentationChallenges 0..1 and
    followUpPlan 0..1

// Adverse Event Monitoring
* extension contains https://fhir.zarishsphere.com/StructureDefinition/immunization-adverse-event named adverseEvent 0..1
* extension[adverseEvent].extension contains
    eventSeverity 0..1 and
    eventOnset 0..1 and
    eventDuration 0..1 and
    medicalAttentionRequired 0..1 and
    reportedTo 0..*

// Vaccine Effectiveness
* extension contains https://fhir.zarishsphere.com/StructureDefinition/immunization-effectiveness named effectiveness 0..1
* extension[effectiveness].extension contains
    seroconversion 0..1 and
    antibodyTiter 0..1 and
    protectionLevel 0..1 and
    duration 0..1

// Integration with National Systems
* extension contains https://fhir.zarishsphere.com/StructureDefinition/immunization-national-integration named nationalIntegration 0..1
* extension[nationalIntegration].extension contains
    epiRegistryId 0..1 and
    dhis2Reported 0..1 and
    reportingDate 0..1 and
    dataQuality 0..1

// ----- Enhanced Vaccine Code Support -----
* vaccineCode.coding ^slicing.discriminator.type = #value
* vaccineCode.coding ^slicing.discriminator.path = "system"
* vaccineCode.coding ^slicing.rules = #open
* vaccineCode.coding contains
    bdEpi 1..1 and
    international 0..1 and
    manufacturer 0..1

* vaccineCode.coding[bdEpi].system 1..1
* vaccineCode.coding[bdEpi].system = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-vaccine-code"
* vaccineCode.coding[bdEpi].code 1..1

* vaccineCode.coding[international].system = "http://hl7.org/fhir/sid/cvx"
* vaccineCode.coding[international].code 0..1

* vaccineCode.coding[manufacturer].system = "http://hl7.org/fhir/sid/mtx"
* vaccineCode.coding[manufacturer].code 0..1

// ----- Enhanced Dose Information -----
* doseSequence 0..1 MS
* doseSequence ^short = "Dose sequence within series"
* doseSequence ^definition = "The dose number in the vaccination series"

* series 0..1 MS
* series ^short = "Vaccine series name"
* series ^definition = "The name of the vaccine series this dose belongs to"

* seriesDoses 0..1 MS
* seriesDoses ^short = "Number of doses in series"
* seriesDoses ^definition = "The total number of doses in the vaccine series"

// ----- Enhanced Administration Information -----
* administrationSite 0..1 MS
* administrationSite ^short = "Body site of administration"
* administrationSite from BDImmunizationSiteVS

* route 0..1 MS
* route ^short = "Route of administration"
* route from BDImmunizationRouteVS

// ----- Enhanced Performer Information -----
* performer ^slicing.discriminator.type = #value
* performer ^slicing.discriminator.path = "function.coding.code"
* performer ^slicing.rules = #open
* performer contains
    administrator 0..1 and
    verifier 0..1 and
    witness 0..1

* performer[administrator].function.coding.code = #administering
* performer[administrator].actor only Reference(ZSEnhancedBDPractitionerProfile)

* performer[verifier].function.coding.code = #verifier
* performer[verifier].actor only Reference(ZSEnhancedBDPractitionerProfile)

* performer[witness].function.coding.code = #witness
* performer[witness].actor only Reference(ZSEnhancedBDPractitionerProfile or RelatedPerson)

// ----- Enhanced Reason Codes -----
* reasonCode 0..* MS
* reasonCode from http://snomed.info/sct (preferred)

* reasonReference 0..* MS
* reasonReference only Reference(Condition or Observation or DiagnosticReport)

// ----- Enhanced Subpotent Reason -----
* subpotentReason 0..* MS
* subpotentReason from http://hl7.org/fhir/ValueSet/immunization-subpotent-reason (preferred)

// ----- Enhanced Education Support -----
* education 0..* MS
* education.documentType 0..1 MS
* education.presentationDateTime 0..1 MS
* education.presentationTitle 0..1 MS
* education.publicationDate 0..1 MS
* education.presentationFormat 0..1 MS

// ----- Program Eligibility -----
* programEligibility 0..* MS
* programEligibility.program 0..1 MS
* programEligibility.program from https://fhir.dghs.gov.bd/core/ValueSet/bd-immunization-programs (required)
* programEligibility.programStatus 0..1 MS
* programEligibility.programStatus from http://hl7.org/fhir/ValueSet/immunization-program-eligibility (required)

// ----- Funding Source -----
* fundingSource 0..1 MS
* fundingSource from https://fhir.dghs.gov.bd/core/ValueSet/bd-immunization-funding-sources (required)

// ----- Bangladesh-Specific Validation Rules -----
* invariant zs-bd-vaccine-schedule "Vaccine doses SHALL follow Bangladesh EPI schedule" {
  %resource.vaccineCode.coding.where(system = 'https://fhir.dghs.gov.bd/core/CodeSystem/bd-vaccine-code').exists()
  implies (
    %resource.doseSequence.exists() and %resource.series.exists()
  )
}

* invariant zs-cold-chain-temperature "Cold chain temperature SHALL be monitored for all vaccines" {
  %resource.extension.where(url = 'https://fhir.zarishsphere.com/StructureDefinition/immunization-cold-chain').exists()
  implies (
    %resource.extension.where(url = 'https://fhir.zarishsphere.com/StructureDefinition/immunization-cold-chain').extension.where(url = 'storageTemperature').exists()
  )
}

* invariant zs-refugee-documentation "Refugee immunizations SHALL include special documentation" {
  %resource.patient.resolve().extension.where(url = 'https://fhir.zarishsphere.com/StructureDefinition/patient-rohingya-status').exists()
  implies (
    %resource.extension.where(url = 'https://fhir.zarishsphere.com/StructureDefinition/immunization-refugee-context').exists()
  )
}
