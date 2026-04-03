// ============================================================
// ZS-Enhanced-BD-Condition.fsh
// ZarishSphere Enhanced Bangladesh Condition Profile
// Extends BD-Core-FHIR-IG with ICD-11 MMS support and enhanced features
// Includes communicable disease surveillance, outbreak detection, cluster validation
// ============================================================

Profile: ZSEnhancedBDConditionProfile
Parent: Condition
Id: zs-enhanced-bd-condition
Title: "ZarishSphere Enhanced Bangladesh Condition Profile (ICD-11)"
Description: """
Enhanced Condition profile for ZarishSphere building on Bangladesh Core FHIR IG.
- All BD-Core ICD-11 MMS requirements maintained with cluster expression support
- Additional ZarishSphere extensions: communicable disease surveillance, outbreak detection
- Enhanced privacy and security for sensitive health conditions
- Integration with national disease surveillance systems
- Support for refugee-specific health conditions
"""

// ----- BD Core ICD-11 Requirements (Maintained) -----
* code 1..1 MS
* code from bd-condition-icd11-diagnosis-valueset (preferred)
* code ^comment = """
Condition.code SHALL contain at least one coding conforming to the
coding[stem] slice with system = http://id.who.int/icd/release/11/mms.

Stem code rules:
  - The stem code SHALL be a Diagnosis or Finding class ICD-11 MMS concept.
  - This restriction is enforced at runtime via OCL ValueSet $validate-code
    against the Bangladesh ICD-11 MMS Condition ValueSet.
  - Stem-only codes SHALL be validated via OCL $validate-code.
  - Substance, Organism, Device, Anatomy, and Misc class concepts SHALL NOT
    appear as standalone stem codes in Condition.code.

Cluster expression rules:
  - When a concept requires postcoordination, the full cluster expression
    SHALL be carried in the icd11-cluster-expression extension on coding[stem].
  - The icd11-cluster-expression extension SHALL only be present when the
    expression contains at least one satellite code joined by & or / operators.
  - Satellite codes in the cluster expression are exempt from the
    Diagnosis/Finding class restriction.
  - Cluster expressions SHALL be validated against the Bangladesh ICD-11
    Cluster Validator at https://icd11.dghs.gov.bd/cluster/validate
    prior to submission to the HIE.

Additional local codings are permitted alongside the mandatory ICD-11 stem
(slicing is open). Cluster expressions are typically sourced from the WHO
Electronic Coding Tool (ECT) at the point of care.
"""

* code.coding ^slicing.discriminator.type = #value
* code.coding ^slicing.discriminator.path = "system"
* code.coding ^slicing.rules = #open
* code.coding ^slicing.description = "Slice requiring exactly one ICD-11 MMS stem code. Additional local codings permitted."

* code.coding contains stem 1..1
* code.coding[stem] ^short = "Mandatory ICD-11 MMS stem code"
* code.coding[stem] ^definition = """
Exactly one ICD-11 MMS stem code is required. The stem code SHALL be a
Diagnosis or Finding class concept. When the condition requires
postcoordination, the full cluster expression is carried in the
icd11-cluster-expression extension on this coding element.
"""
* code.coding[stem].system 1..1
* code.coding[stem].system = "http://id.who.int/icd/release/11/mms" (exactly)
* code.coding[stem].code 1..1
* code.coding[stem].extension contains
    https://fhir.dghs.gov.bd/core/StructureDefinition/icd11-cluster-expression named clusterExpression 0..1

// ----- ZarishSphere Enhanced Extensions -----

// Communicable Disease Surveillance
* extension contains https://fhir.zarishsphere.com/StructureDefinition/condition-communicable-disease named communicableDisease 0..1
* extension[communicableDisease].extension contains
    diseaseCode 0..1 and
    reportingRequired 0..1 and
    outbreakStatus 0..1 and
    surveillanceCategory 0..1 and
    reportingDeadline 0..1 and
    notifiedTo 0..* and
    caseClassification 0..1

// Disease Severity and Risk Assessment
* extension contains https://fhir.zarishsphere.com/StructureDefinition/condition-severity-risk named severityRisk 0..1
* extension[severityRisk].extension contains
    clinicalSeverity 0..1 and
    publicHealthRisk 0..1 and
    transmissionRisk 0..1 and
    mortalityRisk 0..1 and
    complicationRisk 0..1

// Refugee-Specific Health Conditions
* extension contains https://fhir.zarishsphere.com/StructureDefinition/condition-refugee-context named refugeeContext 0..1
* extension[refugeeContext].extension contains
    campRelated 0..1 and
    displacementRelated 0..1 and
    traumaRelated 0..1 and
    nutritionalDeficiency 0..1 and
    waterborneDisease 0..1

// Treatment and Management Information
* extension contains https://fhir.zarishsphere.com/StructureDefinition/condition-treatment-plan named treatmentPlan 0..1
* extension[treatmentPlan].extension contains
    treatmentProtocol 0..1 and
    medicationRequired 0..* and
    followUpRequired 0..1 and
    referralNeeded 0..1 and
    isolationRequired 0..1

// Outbreak and Cluster Detection
* extension contains https://fhir.zarishsphere.com/StructureDefinition/condition-outbreak-info named outbreakInfo 0..1
* extension[outbreakInfo].extension contains
    outbreakId 0..1 and
    clusterId 0..1 and
    indexCase 0..1 and
    transmissionChain 0..1 and
    investigationStatus 0..1

// Privacy and Confidentiality
* extension contains https://fhir.zarishsphere.com/StructureDefinition/condition-privacy named privacy 0..1
* extension[privacy].extension contains
    confidentialityLevel 0..1 and
    restrictedAccess 0..1 and
    consentRequired 0..1 and
    dataSharingLimitations 0..*

// ----- Enhanced Clinical Status Support -----
* clinicalStatus 1..1 MS
* clinicalStatus from http://hl7.org/CodeSystem/condition-clinical-status (required)

* verificationStatus 1..1 MS
* verificationStatus from http://hl7.org/CodeSystem/condition-ver-status (required)

* category 1..* MS
* category from http://terminology.hl7.org/CodeSystem/condition-category (required)

* severity 0..1 MS
* severity from http://hl7.org/CodeSystem/condition-severity (preferred)

// ----- Enhanced Subject and Encounter Support -----
* subject 1..1 MS
* subject only Reference(ZSEnhancedBDPatientProfile)

* encounter 0..1 MS
* encounter only Reference(ZSEnhancedBDEncounterProfile)

* onset[x] 1..1 MS
* onset[x] only dateTime or Period or Range or Age

* abatement[x] 0..1 MS
* abatement[x] only dateTime or Period or Range or Age

// ----- Enhanced Evidence Support -----
* evidence 0..* MS
* evidence.detail only Reference(Observation or DiagnosticReport or QuestionnaireResponse)
* evidence.code from http://snomed.info/sct (preferred)

// ----- Enhanced Body Site Support -----
* bodySite 0..* MS
* bodySite from http://snomed.info/sct (preferred)

// ----- Enhanced Stage Support -----
* stage 0..* MS
* stage.summary 0..1 MS
* stage.summary from http://snomed.info/sct (preferred)
* stage.assessment 0..1 MS
* stage.assessment only Reference(Observation or DiagnosticReport or QuestionnaireResponse)

// ----- Enhanced Note Support -----
* note 0..* MS
* note.author[x] only Reference(Patient or Practitioner or RelatedPerson or Organization) or string
* note.time 0..1 MS
* note.text 1..1 MS

// ----- Bangladesh-Specific Validation Rules -----
* invariant zs-bd-icd11-stem-only "All ICD-11 stem codes SHALL be Diagnosis or Finding class concepts" {
  %resource.code.coding.where(system = 'http://id.who.int/icd/release/11/mms').all(
    code.memberOf('http://fhir.dghs.gov.bd/core/ValueSet/bd-icd11-diagnosis-finding-valueset')
  )
}

* invariant zs-bd-cluster-expression "Cluster expressions SHALL be validated before submission" {
  %resource.code.coding.where(system = 'http://id.who.int/icd/release/11/mms').extension.where(url = 'https://fhir.dghs.gov.bd/core/StructureDefinition/icd11-cluster-expression').exists()
  implies (
    %resource.code.coding.where(system = 'http://id.who.int/icd/release/11/mms').extension.where(url = 'https://fhir.dghs.gov.bd/core/StructureDefinition/icd11-cluster-expression').valueString.matches('.*&.*|.*\\/.*')
  )
}

* invariant zs-communicable-disease-reporting "Communicable diseases SHALL be reported within 24 hours" {
  %resource.extension.where(url = 'https://fhir.zarishsphere.com/StructureDefinition/condition-communicable-disease').extension.where(url = 'reportingRequired').valueBoolean = true
  implies (
    %resource.extension.where(url = 'https://fhir.zarishsphere.com/StructureDefinition/condition-communicable-disease').extension.where(url = 'reportingDeadline').exists()
  )
}
