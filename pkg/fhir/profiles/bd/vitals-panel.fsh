Profile: ZSVitalsPanelV2
Id: zs-vitals-panel-v2
Parent: Observation
Title: "ZarishSphere Vitals Panel Profile v2"
Description: """
Comprehensive vital signs panel for ZarishSphere with Bangladesh-specific clinical protocols,
enhanced privacy controls, and real-time alerting capabilities.

Key Features:
- Complete vital signs panel (BP, HR, RR, Temp, SpO2)
- Pain assessment with localization
- BMI calculation with growth charts
- Clinical decision support alerts
- Bangladesh-specific reference ranges
- Pediatric and adult adaptations
"""

* status 1..1 MS
* status = #final

* category 1..1 MS
* category = http://terminology.hl7.org/CodeSystem/observation-category#vital-signs

* code 1..1 MS
* code.coding 1..1 MS
* code.coding.system = "http://loinc.org"
* code.coding.code = #85353-1
* code.coding.display = "Vital signs panel"

* subject 1..1 MS
* subject only Reference(Patient or Group)

* encounter 0..1 MS
* encounter only Reference(Encounter)

* effective[x] 1..1 MS
* effective[x] only dateTime

* performer 0..* MS
* performer only Reference(Practitioner | PractitionerRole | Organization | CareTeam | Patient | RelatedPerson)

* component 5..* MS
* component ^slicing.discriminator.type = #value
* component ^slicing.discriminator.path = "code.coding.code"
* component ^slicing.rules = #open
* component contains
    systolicBP 0..1 and
    diastolicBP 0..1 and
    heartRate 0..1 and
    respiratoryRate 0..1 and
    temperature 0..1 and
    oxygenSaturation 0..1 and
    height 0..1 and
    weight 0..1 and
    bmi 0..1 and
    bloodGlucose 0..1 and
    painScale 0..1 and
    painLocation 0..1 and
    consciousnessLevel 0..1 and
    generalAppearance 0..1

// ----- Blood Pressure Components -----
* component[systolicBP].code.coding.system = "http://loinc.org"
* component[systolicBP].code.coding.code = #8480-6
* component[systolicBP].code.coding.display = "Systolic blood pressure"
* component[systolicBP].valueQuantity 1..1 MS
* component[systolicBP].valueQuantity.unit = "mm[Hg]"
* component[systolicBP].valueQuantity.system = "http://unitsofmeasure.org"
* component[systolicBP].interpretation 0..1 MS
* component[systolicBP].interpretation from VitalSignInterpretationVS

* component[diastolicBP].code.coding.system = "http://loinc.org"
* component[diastolicBP].code.coding.code = #8462-4
* component[diastolicBP].code.coding.display = "Diastolic blood pressure"
* component[diastolicBP].valueQuantity 1..1 MS
* component[diastolicBP].valueQuantity.unit = "mm[Hg]"
* component[diastolicBP].valueQuantity.system = "http://unitsofmeasure.org"
* component[diastolicBP].interpretation 0..1 MS
* component[diastolicBP].interpretation from VitalSignInterpretationVS

// ----- Heart Rate Component -----
* component[heartRate].code.coding.system = "http://loinc.org"
* component[heartRate].code.coding.code = #8867-4
* component[heartRate].code.coding.display = "Heart rate"
* component[heartRate].valueQuantity 1..1 MS
* component[heartRate].valueQuantity.unit = "/min"
* component[heartRate].valueQuantity.system = "http://unitsofmeasure.org"
* component[heartRate].interpretation 0..1 MS
* component[heartRate].interpretation from VitalSignInterpretationVS

// ----- Respiratory Rate Component -----
* component[respiratoryRate].code.coding.system = "http://loinc.org"
* component[respiratoryRate].code.coding.code = #9279-1
* component[respiratoryRate].code.coding.display = "Respiratory rate"
* component[respiratoryRate].valueQuantity 1..1 MS
* component[respiratoryRate].valueQuantity.unit = "/min"
* component[respiratoryRate].valueQuantity.system = "http://unitsofmeasure.org"
* component[respiratoryRate].interpretation 0..1 MS
* component[respiratoryRate].interpretation from VitalSignInterpretationVS

// ----- Temperature Component -----
* component[temperature].code.coding.system = "http://loinc.org"
* component[temperature].code.coding.code = #8310-5
* component[temperature].code.coding.display = "Body temperature"
* component[temperature].valueQuantity 1..1 MS
* component[temperature].valueQuantity.unit = "Cel"
* component[temperature].valueQuantity.system = "http://unitsofmeasure.org"
* component[temperature].interpretation 0..1 MS
* component[temperature].interpretation from VitalSignInterpretationVS

// ----- Oxygen Saturation Component -----
* component[oxygenSaturation].code.coding.system = "http://loinc.org"
* component[oxygenSaturation].code.coding.code = #59408-5
* component[oxygenSaturation].code.coding.display = "Oxygen saturation"
* component[oxygenSaturation].valueQuantity 1..1 MS
* component[oxygenSaturation].valueQuantity.unit = "%"
* component[oxygenSaturation].valueQuantity.system = "http://unitsofmeasure.org"
* component[oxygenSaturation].interpretation 0..1 MS
* component[oxygenSaturation].interpretation from VitalSignInterpretationVS

// ----- Height Component -----
* component[height].code.coding.system = "http://loinc.org"
* component[height].code.coding.code = #8302-2
* component[height].code.coding.display = "Body height"
* component[height].valueQuantity 1..1 MS
* component[height].valueQuantity.unit = "cm"
* component[height].valueQuantity.system = "http://unitsofmeasure.org"

// ----- Weight Component -----
* component[weight].code.coding.system = "http://loinc.org"
* component[weight].code.coding.code = #29463-7
* component[weight].code.coding.display = "Body weight"
* component[weight].valueQuantity 1..1 MS
* component[weight].valueQuantity.unit = "kg"
* component[weight].valueQuantity.system = "http://unitsofmeasure.org"

// ----- BMI Component -----
* component[bmi].code.coding.system = "http://loinc.org"
* component[bmi].code.coding.code = #39156-5
* component[bmi].code.coding.display = "Body mass index (BMI) [Ratio]"
* component[bmi].valueQuantity 1..1 MS
* component[bmi].valueQuantity.unit = "kg/m2"
* component[bmi].valueQuantity.system = "http://unitsofmeasure.org"
* component[bmi].interpretation 0..1 MS
* component[bmi].interpretation from BMIInterpretationVS

// ----- Blood Glucose Component -----
* component[bloodGlucose].code.coding.system = "http://loinc.org"
* component[bloodGlucose].code.coding.code = #2345-7
* component[bloodGlucose].code.coding.display = "Glucose [Moles/volume] in Blood"
* component[bloodGlucose].valueQuantity 1..1 MS
* component[bloodGlucose].valueQuantity.unit = "mmol/L"
* component[bloodGlucose].valueQuantity.system = "http://unitsofmeasure.org"
* component[bloodGlucose].interpretation 0..1 MS
* component[bloodGlucose].interpretation from GlucoseInterpretationVS

// ----- Pain Assessment Components -----
* component[painScale].code.coding.system = "http://loinc.org"
* component[painScale].code.coding.code = #72514-3
* component[painScale].code.coding.display = "Pain severity scale - 0-10 verbal numeric rating scale"
* component[painScale].valueQuantity 1..1 MS
* component[painScale].valueQuantity.unit = "{score}"
* component[painScale].interpretation 0..1 MS
* component[painScale].interpretation from PainInterpretationVS

* component[painLocation].code.coding.system = "http://loinc.org"
* component[painLocation].code.coding.code = #85358-1
* component[painLocation].code.coding.display = "Location of pain"
* component[painLocation].valueString 0..1 MS

// ----- Consciousness Level Component -----
* component[consciousnessLevel].code.coding.system = "http://loinc.org"
* component[consciousnessLevel].code.coding.code = #67788-6
* component[consciousnessLevel].code.coding.display = "Level of consciousness"
* component[consciousnessLevel].valueCodeableConcept 1..1 MS
* component[consciousnessLevel].valueCodeableConcept from ConsciousnessLevelVS

// ----- General Appearance Component -----
* component[generalAppearance].code.coding.system = "http://loinc.org"
* component[generalAppearance].code.coding.code = #54128-1
* component[generalAppearance].code.coding.display = "General appearance"
* component[generalAppearance].valueString 0..1 MS

// ----- Bangladesh Extensions -----
// Clinical Context
* extension contains https://fhir.zarishsphere.com/StructureDefinition/vitals-clinical-context named clinicalContext 0..1
* extension[clinicalContext].extension contains
    measurementLocation 0..1 and
    measurementDevice 0..1 and
    clinicianNotes 0..1 and
    triageCategory 0..1

// Pediatric Growth Data
* extension contains https://fhir.zarishsphere.com/StructureDefinition/vitals-pediatric-growth named pediatricGrowth 0..1
* extension[pediatricGrowth].extension contains
    ageInMonths 0..1 and
    heightForAge 0..1 and
    weightForAge 0..1 and
    bmiForAge 0..1 and
    growthPercentiles 0..*

// Bangladesh Reference Ranges
* extension contains https://fhir.zarishsphere.com/StructureDefinition/vitals-bd-reference-ranges named bdReferenceRanges 0..1
* extension[bdReferenceRanges].extension contains
    adultNormalRanges 0..1 and
    pediatricNormalRanges 0..1 and
    elderlyAdjustments 0..1 and
    regionalVariations 0..1

// Alert Thresholds
* extension contains https://fhir.zarishsphere.com/StructureDefinition/vitals-alert-thresholds named alertThresholds 0..1
* extension[alertThresholds].extension contains
    criticalThresholds 0..1 and
    warningThresholds 0..1 and
    ageSpecificThresholds 0..* and
    conditionSpecificThresholds 0..*

// Quality Control
* extension contains https://fhir.zarishsphere.com/StructureDefinition/vitals-quality-control named qualityControl 0..1
* extension[qualityControl].extension contains
    measurementMethod 0..1 and
    deviceCalibration 0..1 and
    dataQuality 0..1 and
    qualityScore 0..1
