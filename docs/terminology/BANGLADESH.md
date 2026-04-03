# Bangladesh Divisions

The terminology service includes Bangladesh administrative divisions for use in patient addresses and facility locations.

## System URL

```
https://health.zarishsphere.com/fhir/ValueSet/bd-divisions
```

## Usage

```bash
# Get all Bangladesh divisions
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=https://health.zarishsphere.com/fhir/ValueSet/bd-divisions"

# Filter by text
curl "http://localhost:8080/fhir/ValueSet/\$expand?url=https://health.zarishsphere.com/fhir/ValueSet/bd-divisions&filter=Dhaka"
```

## Divisions

Bangladesh has 8 divisions:

| Code | Name (English) | Name (বাংলা) | Capital |
|------|-----------------|--------------|---------|
| DH | Dhaka | ঢাকা | Dhaka |
| CTG | Chattogram | চট্টগ্রাম | Chattogram |
| SYL | Sylhet | সিলেট | Sylhet |
| KHU | Khulna | খুলনা | Khulna |
| RAJ | Rajshahi | রাজশাহী | Rajshahi |
| RAN | Rangpur | রংপুর | Rangpur |
| BAR | Barisal | বরিশাল | Barisal |
| MYM | Mymensingh | ময়মনসিংহ | Mymensingh |

## Districts

### Dhaka Division (DH)

| Code | District |
|------|----------|
| DHA | Dhaka |
| GAZ | Gazipur |
| NAR | Narayanganj |
| MUN | Munshiganj |
| MAN | Manikganj |
| FAR | Faridpur |
| GOP | Gopalganj |
| SHA | Shariatpur |
| MAD | Madaripur |
| TANG | Tangail |
| KISH | Kishoreganj |
| NET | Netrakona |
| JAM | Jamalpur |

### Chattogram Division (CTG)

| Code | District |
|------|----------|
| CTG | Chattogram |
| COX | Cox's Bazar |
| BAN | Bandarban |
| RANG | Rangamati |
| KHAW | Khagrachari |
| FEN | Feni |
| NOA | Noakhali |
| LAK | Lakshmipur |
| COM | Comilla |
| BRA | Brahmanbaria |
| CHA | Chandpur |

## Usage in FHIR Resources

### Patient Address

```json
{
  "resourceType": "Patient",
  "address": [
    {
      "use": "home",
      "district": "Dhaka",
      "state": "Dhaka",
      "postalCode": "1205",
      "country": "Bangladesh"
    }
  ]
}
```

### Location Resource

```json
{
  "resourceType": "Location",
  "name": "Dhaka Medical College Hospital",
  "type": {
    "coding": [
      {
        "system": "http://terminology.hl7.org/CodeSystem/v3-ServiceType",
        "code": "HOSP",
        "display": "Hospital"
      }
    ]
  },
  "address": {
    "district": "Dhaka",
    "state": "Dhaka",
    "country": "Bangladesh"
  }
}
```

## Extending the List

To add more districts or update codes, modify `cmd/zs-core-fhir/terminology.go`:

```go
func StartTerminologyServer(port int) {
    server := NewTerminologyServer()
    
    // Add Bangladesh divisions
    server.AddConcept(
        "https://health.zarishsphere.com/fhir/ValueSet/bd-divisions",
        "DH",
        "Dhaka"
    )
    server.AddConcept(
        "https://health.zarishsphere.com/fhir/ValueSet/bd-divisions",
        "CTG",
        "Chattogram"
    )
    // ... add more
}
```
