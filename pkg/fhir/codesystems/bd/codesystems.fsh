// ============================================================
// ZS-Enhanced-BD-CodeSystems.fsh
// ZarishSphere Enhanced Bangladesh Code Systems
// Extends BD-Core-FHIR-IG code systems with additional capabilities
// Includes refugee-specific codes, mobile health, surveillance codes
// ============================================================

// ----- Enhanced Identifier Types Code System -----
CodeSystem: ZSEnhancedBangladeshIdentifierType
Id: zs-enhanced-bd-identifier-type
Title: "ZarishSphere Enhanced Bangladesh Identifier Types"
Description: "Enhanced codes identifying the type of identifiers used in Bangladesh, including refugee and mobile health identifiers."
* ^url = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-identifier-type"
* ^caseSensitive = true
* ^experimental = false
* ^content = #complete

// BD Core identifiers maintained
* #NID "National ID"
    "Bangladesh National Identity Card number"
* #BRN "Birth Registration Number"
    "Bangladesh Birth Registration Certificate number"
* #UHID "Unique Health ID"
    "Bangladesh Unique Health Identifier for patient tracking"

// ZarishSphere Enhanced identifiers
* #UNHCR "UNHCR Registration Card"
    "UNHCR refugee registration card number"
* #SMARTCARD "UNHCR Smart Card"
    "UNHCR smart card identifier for refugees"
* #PASSPORT "Passport"
    "International passport number"
* #PHONE "Phone Number"
    "Mobile phone number as identifier"
* #ROHINGYA_ID "Rohingya Community ID"
    "Rohingya community-specific identifier"
* #CAMP_ID "Camp ID"
    "Refugee camp identification number"
* #HOUSEHOLD_ID "Household ID"
    "Household identification number in refugee camps"
* #MOBILE_APP_ID "Mobile App ID"
    "ZarishSphere mobile application user ID"
* #BIOMETRIC_ID "Biometric ID"
    "Biometric identifier for patient verification"
* #TEMP_ID "Temporary ID"
    "Temporary identifier for emergency situations"

ValueSet: ZSEnhancedBangladeshIdentifierTypeVS
Id: zs-enhanced-bd-identifier-type-valueset
Title: "ZarishSphere Enhanced Bangladesh Identifier Type"
Description: "Enhanced Bangladesh Standard Identifier types including refugee and mobile health identifiers"
* ^url = "https://fhir.zarishsphere.com/core/ValueSet/zs-enhanced-bd-identifier-type-valueset"
* include codes from system https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-identifier-type

// ----- Enhanced Vaccine Code System -----
CodeSystem: ZSEnhancedBDVaccineCS
Id: zs-enhanced-bd-vaccine-code
Title: "ZarishSphere Enhanced Bangladesh Vaccine Code System"
Description: "Enhanced vaccine codes used in Bangladesh EPI and immunization program, including refugee-specific vaccines and outbreak response vaccines."
* ^url = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-vaccine-code"
* ^status = #active
* ^caseSensitive = true
* ^experimental = false
* ^content = #complete

// BD Core vaccines maintained
* #BCG "BCG Vaccine"
    "Bacillus Calmette-Guérin vaccine, used against tuberculosis."
* #OPV "Oral Polio Vaccine (OPV)"
    "Live attenuated oral polio vaccine."
* #IPV "Inactivated Polio Vaccine (IPV)"
    "Inactivated polio vaccine."
* #PENTA "Pentavalent Vaccine"
    "DTP-HepB-Hib combined vaccine."
* #MR "Measles-Rubella (MR) Vaccine"
    "Combined measles and rubella vaccine."
* #TT "Tetanus Toxoid (TT) Vaccine"
    "Vaccine used for tetanus prevention."
* #PCV10 "Pneumococcal Conjugate Vaccine"
    "10-valent pneumococcal conjugate vaccine."
* #ROTA "Rotavirus Vaccine"
    "Live attenuated rotavirus vaccine for diarrheal disease prevention."
* #HPV "Human Papillomavirus (HPV) Vaccine"
    "Vaccine used for prevention of cervical cancer and HPV-related diseases."
* #COVID19 "COVID-19 Vaccine"
    "Vaccines against SARS-CoV-2 (various manufacturers)."

// ZarishSphere Enhanced vaccines
* #CHOLERA "Cholera Vaccine"
    "Oral cholera vaccine for outbreak response in refugee camps."
* #TYPHOID "Typhoid Vaccine"
    "Typhoid conjugate vaccine for high-risk populations."
* #HEPATITIS_A "Hepatitis A Vaccine"
    "Hepatitis A vaccine for outbreak prevention."
* #MENINGITIS "Meningococcal Vaccine"
    "Meningococcal vaccine for outbreak response."
* #RABIES "Rabies Vaccine"
    "Rabies vaccine for post-exposure prophylaxis."
* #YELLOW_FEVER "Yellow Fever Vaccine"
    "Yellow fever vaccine for travelers and outbreak response."
* #JAPANESE_ENCEPHALITIS "Japanese Encephalitis Vaccine"
    "Japanese encephalitis vaccine for endemic areas."
* #DIPHTHERIA "Diphtheria Vaccine"
    "Diphtheria vaccine for outbreak response."
* #PERTUSSIS "Pertussis Vaccine"
    "Pertussis vaccine for outbreak control."
* #INFLUENZA "Influenza Vaccine"
    "Seasonal influenza vaccine for high-risk groups."

ValueSet: ZSEnhancedBDVaccineVS
Id: zs-enhanced-bd-vaccine-valueset
Title: "ZarishSphere Enhanced Bangladesh Vaccine Value Set"
Description: "Enhanced Bangladesh vaccine codes including outbreak response and refugee-specific vaccines"
* ^url = "https://fhir.zarishsphere.com/core/ValueSet/zs-enhanced-bd-vaccine-valueset"
* include codes from system https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-vaccine-code

// ----- Enhanced Religions Code System -----
CodeSystem: ZSEnhancedBDReligionsCS
Id: zs-enhanced-bd-religions
Title: "ZarishSphere Enhanced Bangladesh Religions"
Description: "Enhanced religious affiliations in Bangladesh, including minority and refugee community religions."
* ^url = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-religions"
* ^caseSensitive = true
* ^experimental = false
* ^content = #complete

// BD Core religions maintained
* #islam "Islam"
    "Islamic faith"
* #hinduism "Hinduism"
    "Hindu faith"
* #buddhism "Buddhism"
    "Buddhist faith"
* #christianity "Christianity"
    "Christian faith"
* #other_religion "Other Religion"
    "Other religious affiliations"

// ZarishSphere Enhanced religions
* #rohingya_islamic "Rohingya Islamic"
    "Rohingya Islamic practices and traditions"
* #rohingya_cultural "Rohingya Cultural"
    "Rohingya cultural and traditional practices"
* #indigenous "Indigenous"
    "Indigenous religious practices"
* #animist "Animist"
    "Animist religious practices"
* #no_religion "No Religion"
    "No religious affiliation"
* #unknown "Unknown"
    "Religious affiliation unknown"

ValueSet: ZSEnhancedBDReligionsVS
Id: zs-enhanced-bd-religions-valueset
Title: "ZarishSphere Enhanced Bangladesh Religions"
Description: "Enhanced religious affiliations including minority and refugee community religions"
* ^url = "https://fhir.zarishsphere.com/core/ValueSet/zs-enhanced-bd-religions-valueset"
* include codes from system https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-religions

// ----- Enhanced Communicable Diseases Code System -----
CodeSystem: ZSEnhancedBDCommunicableDiseasesCS
Id: zs-enhanced-bd-communicable-diseases
Title: "ZarishSphere Enhanced Bangladesh Communicable Diseases"
Description: "Enhanced communicable diseases under surveillance in Bangladesh, including refugee-specific diseases."
* ^url = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-communicable-diseases"
* ^caseSensitive = true
* ^experimental = false
* ^content = #complete

// BD Core diseases maintained
* #CHOLERA "Cholera"
    "Cholera infection"
* #DIARRHEA "Diarrheal Disease"
    "Acute diarrheal disease"
* #DENGUE "Dengue Fever"
    "Dengue viral infection"
* #MALARIA "Malaria"
    "Malaria parasitic infection"
* #TUBERCULOSIS "Tuberculosis"
    "Mycobacterium tuberculosis infection"
* #COVID19 "COVID-19"
    "SARS-CoV-2 infection"
* #MEASLES "Measles"
    "Measles viral infection"
* #POLIO "Poliomyelitis"
    "Poliovirus infection"

// ZarishSphere Enhanced diseases
* #ACUTE_RESPIRATORY_INFECTION "Acute Respiratory Infection"
    "Acute respiratory infection in refugee camps"
* #SKIN_INFECTION "Skin Infection"
    "Skin infections in crowded settings"
* #MALNUTRITION_COMPLICATIONS "Malnutrition Complications"
    "Complications arising from severe malnutrition"
* #WATERBORNE_DISEASE "Waterborne Disease"
    "Diseases from contaminated water sources"
* #VECTOR_BORNE_DISEASE "Vector-Borne Disease"
    "Diseases transmitted by vectors in camp settings"
* #TYPHOID "Typhoid Fever"
    "Typhoid fever in high-density populations"
* #HEPATITIS_A "Hepatitis A"
    "Hepatitis A in poor sanitation settings"
* #SCABIES "Scabies"
    "Scabies in crowded living conditions"
* #MENTAL_HEALTH_CRISIS "Mental Health Crisis"
    "Mental health emergencies in refugee populations"

ValueSet: ZSEnhancedBDCommunicableDiseasesVS
Id: zs-enhanced-bd-communicable-diseases-valueset
Title: "ZarishSphere Enhanced Bangladesh Communicable Diseases"
Description: "Enhanced communicable diseases under surveillance including refugee-specific diseases"
* ^url = "https://fhir.zarishsphere.com/core/ValueSet/zs-enhanced-bd-communicable-diseases-valueset"
* include codes from system https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-communicable-diseases

// ----- Enhanced Facility Types Code System -----
CodeSystem: ZSEnhancedBDFacilityTypesCS
Id: zs-enhanced-bd-facility-types
Title: "ZarishSphere Enhanced Bangladesh Facility Types"
Description: "Enhanced healthcare facility types in Bangladesh, including refugee-specific facilities."
* ^url = "https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-facility-types"
* ^caseSensitive = true
* ^experimental = false
* ^content = #complete

// BD Core facility types maintained
* #primary-health-center "Primary Health Center"
    "Primary level healthcare facility"
* #community-clinic "Community Clinic"
    "Community-based healthcare facility"
* #upazila-health-complex "Upazila Health Complex"
    "Upazila level healthcare facility"
* #district-hospital "District Hospital"
    "District level hospital"
* #medical-college-hospital "Medical College Hospital"
    "Teaching hospital with medical college"
* #specialized-hospital "Specialized Hospital"
    "Specialized healthcare facility"
* #laboratory "Laboratory"
    "Diagnostic laboratory facility"
* #emergency-department "Emergency Department"
    "Emergency care facility"

// ZarishSphere Enhanced facility types
* #refugee-health-center "Refugee Health Center"
    "Healthcare facility in refugee camps"
* #mobile-clinic "Mobile Clinic"
    "Mobile healthcare unit for outreach"
* #field-hospital "Field Hospital"
    "Temporary hospital for emergency response"
* #isolation-center "Isolation Center"
    "Facility for infectious disease isolation"
* #nutrition-center "Nutrition Center"
    "Facility for nutrition rehabilitation"
* #mental-health-center "Mental Health Center"
    "Facility for mental health services"
* #rehabilitation-center "Rehabilitation Center"
    "Facility for rehabilitation services"
* #quarantine-facility "Quarantine Facility"
    "Facility for quarantine purposes"
* #vaccination-center "Vaccination Center"
    "Facility for mass vaccination campaigns"
* #outreach-post "Outreach Post"
    "Remote healthcare service point"

ValueSet: ZSEnhancedBDFacilityTypesVS
Id: zs-enhanced-bd-facility-types-valueset
Title: "ZarishSphere Enhanced Bangladesh Facility Types"
Description: "Enhanced healthcare facility types including refugee-specific facilities"
* ^url = "https://fhir.zarishsphere.com/core/ValueSet/zs-enhanced-bd-facility-types-valueset"
* include codes from system https://fhir.zarishsphere.com/core/CodeSystem/zs-enhanced-bd-facility-types

// ----- Enhanced Rohingya Status Code System -----
CodeSystem: ZSRohingyaStatusCS
Id: zs-rohingya-status
Title: "ZarishSphere Rohingya Status"
Description: "Rohingya community status and classification codes for healthcare services."
* ^url = "https://fhir.zarishsphere.com/core/CodeSystem/zs-rohingya-status"
* ^caseSensitive = true
* ^experimental = false
* ^content = #complete

* #refugee "Refugee"
    "Registered refugee under UNHCR protection"
* #asylum_seeker "Asylum Seeker"
    "Individual seeking asylum status"
* #undocumented "Undocumented"
    "Undocumented Rohingya individual"
* #returnee "Returnee"
    "Rohingya who returned from displacement"
* #host_community "Host Community"
    "Member of host community living near refugee camps"
* #camp_resident "Camp Resident"
    "Individual residing in refugee camps"
* #urban_refugee "Urban Refugee"
    "Refugee living in urban areas"
* #vulnerable "Vulnerable"
    "Vulnerable Rohingya individual requiring special protection"
* #elderly "Elderly"
    "Elderly Rohingya person (65+ years)"
* #unaccompanied_minor "Unaccompanied Minor"
    "Minor without accompanying adult"
* #person_with_disability "Person with Disability"
    "Rohingya person with disability"

ValueSet: ZSRohingyaStatusVS
Id: zs-rohingya-status-valueset
Title: "ZarishSphere Rohingya Status"
Description: "Rohingya community status and classification codes"
* ^url = "https://fhir.zarishsphere.com/core/ValueSet/zs-rohingya-status-valueset"
* include codes from system https://fhir.zarishsphere.com/core/CodeSystem/zs-rohingya-status

// ----- Enhanced Privacy Classification Code System -----
CodeSystem: ZSPrivacyClassificationCS
Id: zs-privacy-classification
Title: "ZarishSphere Privacy Classification"
Description: "Privacy classification levels for healthcare data, especially for vulnerable populations."
* ^url = "https://fhir.zarishsphere.com/core/CodeSystem/zs-privacy-classification"
* ^caseSensitive = true
* ^experimental = false
* ^content = #complete

* #public "Public"
    "Information that can be freely disclosed"
* #restricted "Restricted"
    "Information with restricted access"
* #confidential "Confidential"
    "Confidential patient information"
* #sensitive "Sensitive"
    "Sensitive health information requiring special protection"
* #highly_sensitive "Highly Sensitive"
    "Highly sensitive information (e.g., refugee status, trauma)"
* #critical "Critical"
    "Critical information requiring highest level protection"

ValueSet: ZSPrivacyClassificationVS
Id: zs-privacy-classification-valueset
Title: "ZarishSphere Privacy Classification"
Description: "Privacy classification levels for healthcare data"
* ^url = "https://fhir.zarishsphere.com/core/ValueSet/zs-privacy-classification-valueset"
* include codes from system https://fhir.zarishsphere.com/core/CodeSystem/zs-privacy-classification
