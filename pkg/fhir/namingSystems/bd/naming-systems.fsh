// ============================================================
// ZS-Enhanced-BD-NamingSystems.fsh
// ZarishSphere Enhanced Bangladesh Naming Systems
// Extends BD-Core-FHIR-IG naming systems with additional capabilities
// Includes ICD-11 MMS, enhanced identifier systems, and local terminology
// ============================================================

// ----- Enhanced ICD-11 MMS Naming System -----
Instance: ZSICD11MMSBangladesh
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Enhanced ICD-11 MMS Naming System — Bangladesh"
Description: """
Enhanced ICD-11 Mortality and Morbidity Statistics (MMS) coding system
as a known and supported terminology within the ZarishSphere Bangladesh
national health information infrastructure.

Canonical system URI: http://id.who.int/icd/release/11/mms
Canonical authority: World Health Organization (WHO)

Preferred code form: short stem codes (e.g. 1A00, NC72.Z).
Linearization URIs are not used as code identifiers in this IG.

National terminology resolver (OCL):
  https://tr.ocl.dghs.gov.bd

Supported OCL operations (use `system=` parameter, not `url=`):
  - $validate-code: https://tr.ocl.dghs.gov.bd/api/fhir/CodeSystem/$validate-code
      ?system=http://id.who.int/icd/release/11/mms&code={code}
  - $lookup: https://tr.ocl.dghs.gov.bd/api/fhir/CodeSystem/$lookup
      ?system=http://id.who.int/icd/release/11/mms&code={code}

$expand is not supported — known OCL limitation.

Version 2025-01 is active in the national OCL instance with 36,941
imported concepts. The OCL resolver is an internal national service;
vendors do not interact with it directly. All vendor submissions are
validated at the HIE boundary via the Bangladesh ICD-11 Cluster Validator
at https://icd11.dghs.gov.bd/cluster/validate.

ZarishSphere enhancements:
  - Additional cluster expression validation
  - Refugee-specific disease coding
  - Enhanced surveillance integration
  - Mobile health coding support
"""

* name = "ZSICD11MMSBangladesh"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Enhanced ICD-11 MMS naming system for ZarishSphere Bangladesh"
* responsible = "WHO"
* type = #codesystem
* usage = "Diagnosis and condition coding in ZarishSphere Bangladesh"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://id.who.int/icd/release/11/mms"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Canonical ICD-11 MMS URI"

* uniqueId[+].type = #oid
* uniqueId[=].value = "2.16.840.1.113883.6.3"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Legacy OID for ICD-11"

// ----- Enhanced Bangladesh Identifier Naming System -----
Instance: ZSBangladeshIdentifiers
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Enhanced Bangladesh Identifiers"
Description: """
Enhanced identifier naming system for ZarishSphere Bangladesh.
Includes all BD-Core identifiers plus ZarishSphere-specific identifiers
for refugee populations, mobile health, and special programs.
"""

* name = "ZSBangladeshIdentifiers"
* status = #active
* kind = #identifier
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Enhanced identifier system for ZarishSphere Bangladesh healthcare"
* responsible = "DGHS/MoHFW, Bangladesh"
* type = #identifier
* usage = "Patient and resource identification in ZarishSphere Bangladesh"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://dghs.gov.bd/identifier/nid"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "National ID"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://dghs.gov.bd/identifier/brn"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Birth Registration Number"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://dghs.gov.bd/identifier/uhid"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Unique Health ID"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://unhcr.org/identifier/registration"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "UNHCR Registration"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://unhcr.org/identifier/smartcard"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "UNHCR Smart Card"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://zarishsphere.com/identifier/rohingya"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Rohingya Community ID"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://zarishsphere.com/identifier/camp"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Camp ID"

* uniqueId[+].type = #uri
* uniqueId[=].value = "http://zarishsphere.com/identifier/mobile-app"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Mobile App ID"

// ----- Enhanced Vaccine Naming System -----
Instance: ZSBangladeshVaccines
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Enhanced Bangladesh Vaccines"
Description: """
Enhanced vaccine naming system for ZarishSphere Bangladesh.
Includes BD-Core EPI vaccines plus ZarishSphere-specific vaccines
for refugee populations, outbreak response, and special campaigns.
"""

* name = "ZSBangladeshVaccines"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Enhanced vaccine coding system for ZarishSphere Bangladesh"
* responsible = "DGHS/MoHFW, Bangladesh"
* type = #codesystem
* usage = "Vaccine coding in ZarishSphere Bangladesh immunization program"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-vaccine-code"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Bangladesh EPI Vaccine Codes"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-vaccine-code"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Enhanced ZarishSphere Vaccine Codes"

// ----- Enhanced Facility Naming System -----
Instance: ZSBangladeshFacilities
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Enhanced Bangladesh Facilities"
Description: """
Enhanced facility naming system for ZarishSphere Bangladesh.
Includes BD-Core facility types plus ZarishSphere-specific facilities
for refugee camps, mobile clinics, and emergency response.
"""

* name = "ZSBangladeshFacilities"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Enhanced facility coding system for ZarishSphere Bangladesh"
* responsible = "DGHS/MoHFW, Bangladesh"
* type = #codesystem
* usage = "Facility coding in ZarishSphere Bangladesh healthcare network"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-facility-types"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Bangladesh Facility Types"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-facility-types"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Enhanced ZarishSphere Facility Types"

// ----- Enhanced Communicable Disease Naming System -----
Instance: ZSBangladeshCommunicableDiseases
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Enhanced Bangladesh Communicable Diseases"
Description: """
Enhanced communicable disease naming system for ZarishSphere Bangladesh.
Includes BD-Core diseases plus ZarishSphere-specific diseases
for refugee populations and outbreak surveillance.
"""

* name = "ZSBangladeshCommunicableDiseases"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Enhanced communicable disease coding for ZarishSphere Bangladesh"
* responsible = "IEDCR, Bangladesh"
* type = #codesystem
* usage = "Disease coding in ZarishSphere Bangladesh surveillance system"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-communicable-diseases"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Bangladesh Communicable Diseases"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-communicable-diseases"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Enhanced ZarishSphere Communicable Diseases"

// ----- Enhanced Religious Affiliation Naming System -----
Instance: ZSBangladeshReligions
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Enhanced Bangladesh Religions"
Description: """
Enhanced religious affiliation naming system for ZarishSphere Bangladesh.
Includes BD-Core religions plus ZarishSphere-specific religious affiliations
for Rohingya and minority communities.
"""

* name = "ZSBangladeshReligions"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Enhanced religious affiliation coding for ZarishSphere Bangladesh"
* responsible = "DGHS/MoHFW, Bangladesh"
* type = #codesystem
* usage = "Religious affiliation coding in ZarishSphere Bangladesh"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.dghs.gov.bd/core/CodeSystem/bd-religions"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Bangladesh Religions"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-religions"
* uniqueId[=].preferred = false
* uniqueId[=].comment = "Enhanced ZarishSphere Religions"

// ----- Enhanced Rohingya Status Naming System -----
Instance: ZSRohingyaStatus
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Rohingya Status"
Description: """
Rohingya community status and classification naming system for ZarishSphere.
Includes protection status, vulnerability assessment, and special care requirements.
"""

* name = "ZSRohingyaStatus"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Rohingya community status coding for ZarishSphere"
* responsible = "UNHCR/DGHS, Bangladesh"
* type = #codesystem
* usage = "Rohingya status coding in ZarishSphere refugee health system"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-rohingya-status"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Rohingya Status Codes"

// ----- Enhanced Privacy Classification Naming System -----
Instance: ZSPrivacyClassification
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Privacy Classification"
Description: """
Privacy classification naming system for ZarishSphere.
Includes classification levels for healthcare data, especially for vulnerable populations.
"""

* name = "ZSPrivacyClassification"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Privacy classification coding for ZarishSphere"
* responsible = "DGHS/Data Protection Authority, Bangladesh"
* type = #codesystem
* usage = "Privacy classification in ZarishSphere healthcare system"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-privacy-classification"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Privacy Classification Codes"

// ----- Enhanced Mobile Health Naming System -----
Instance: ZSMobileHealth
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Mobile Health"
Description: """
Mobile health integration naming system for ZarishSphere.
Includes app consent, device registration, and data synchronization codes.
"""

* name = "ZSMobileHealth"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Mobile health coding for ZarishSphere"
* responsible = "Zarishsphere Digital Health Team"
* type = #codesystem
* usage = "Mobile health coding in ZarishSphere digital health system"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-mobile-health"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Mobile Health Codes"

// ----- Enhanced Cold Chain Naming System -----
Instance: ZSColdChain
InstanceOf: NamingSystem
Usage: #definition
Title: "ZarishSphere Cold Chain"
Description: """
Cold chain management naming system for ZarishSphere.
Includes temperature monitoring, excursion tracking, and quality control codes.
"""

* name = "ZSColdChain"
* status = #active
* kind = #codesystem
* date = "2026-04-03"
* publisher = "ZarishSphere Health Authority"
* description = "Cold chain management coding for ZarishSphere"
* responsible = "DGHS/Expanded Program on Immunization, Bangladesh"
* type = #codesystem
* usage = "Cold chain coding in ZarishSphere immunization system"

* uniqueId[+].type = #uri
* uniqueId[=].value = "https://fhir.zarishsphere.com/core/CodeSystem/zs-cold-chain"
* uniqueId[=].preferred = true
* uniqueId[=].comment = "Cold Chain Codes"
