package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// TranslationService manages internationalization
type TranslationService struct {
	translations map[string]map[string]string
	mu           sync.RWMutex
	loadedLangs  map[string]bool
	basePath     string
}

// NewTranslationService creates a new translation service
func NewTranslationService(basePath string) *TranslationService {
	return &TranslationService{
		translations: make(map[string]map[string]string),
		loadedLangs:  make(map[string]bool),
		basePath:     basePath,
	}
}

// LoadTranslations loads translation files for all supported languages
func (ts *TranslationService) LoadTranslations() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Required languages for MVP
	requiredLangs := []string{"en", "bn"}
	
	for _, lang := range requiredLangs {
		if err := ts.loadLanguage(lang); err != nil {
			return fmt.Errorf("failed to load %s translations: %w", lang, err)
		}
		ts.loadedLangs[lang] = true
	}

	return nil
}

// loadLanguage loads translations for a specific language
func (ts *TranslationService) loadLanguage(lang string) error {
	filename := filepath.Join(ts.basePath, fmt.Sprintf("%s.json", lang))
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read translation file %s: %w", filename, err)
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return fmt.Errorf("failed to parse translation file %s: %w", filename, err)
	}

	ts.translations[lang] = translations
	return nil
}

// Translate translates a key to the specified language
func (ts *TranslationService) Translate(key, language string) string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	// Check if language is loaded
	if !ts.loadedLangs[language] {
		// Fallback to English
		language = "en"
	}

	// Try to get translation
	if langTranslations, exists := ts.translations[language]; exists {
		if translation, found := langTranslations[key]; found {
			return translation
		}
	}

	// Fallback to English if not found
	if language != "en" {
		if enTranslations, exists := ts.translations["en"]; exists {
			if translation, found := enTranslations[key]; found {
				return translation
			}
		}
	}

	// Return key if no translation found
	return key
}

// TranslateWithContext translates a key with context
func (ts *TranslationService) TranslateWithContext(key, language, context string) string {
	// Try context-specific key first
	contextKey := fmt.Sprintf("%s.%s", context, key)
	if translation := ts.Translate(contextKey, language); translation != contextKey {
		return translation
	}

	// Fallback to regular key
	return ts.Translate(key, language)
}

// ValidateKey validates that a translation key exists in all required languages
func (ts *TranslationService) ValidateKey(key string) []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	var missingLangs []string

	// Check required languages
	requiredLangs := []string{"en", "bn"}
	
	for _, lang := range requiredLangs {
		if !ts.loadedLangs[lang] {
			missingLangs = append(missingLangs, lang)
			continue
		}

		if langTranslations, exists := ts.translations[lang]; exists {
			if _, found := langTranslations[key]; !found {
				missingLangs = append(missingLangs, lang)
			}
		}
	}

	return missingLangs
}

// GetAllKeys returns all translation keys for a language
func (ts *TranslationService) GetAllKeys(language string) []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	var keys []string

	if langTranslations, exists := ts.translations[language]; exists {
		for key := range langTranslations {
			keys = append(keys, key)
		}
	}

	return keys
}

// ValidateAllKeys validates that all keys exist in all required languages
func (ts *TranslationService) ValidateAllKeys() map[string][]string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	validationErrors := make(map[string][]string)

	// Get all keys from English (base language)
	if enTranslations, exists := ts.translations["en"]; exists {
		for key := range enTranslations {
			missingLangs := ts.ValidateKey(key)
			if len(missingLangs) > 0 {
				validationErrors[key] = missingLangs
			}
		}
	}

	return validationErrors
}

// AddTranslation adds or updates a translation
func (ts *TranslationService) AddTranslation(language, key, value string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if ts.translations[language] == nil {
		ts.translations[language] = make(map[string]string)
	}

	ts.translations[language][key] = value
}

// RemoveTranslation removes a translation
func (ts *TranslationService) RemoveTranslation(language, key string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if langTranslations, exists := ts.translations[language]; exists {
		delete(langTranslations, key)
	}
}

// SaveTranslations saves translations to files
func (ts *TranslationService) SaveTranslations() error {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	for lang, translations := range ts.translations {
		filename := filepath.Join(ts.basePath, fmt.Sprintf("%s.json", lang))
		
		data, err := json.MarshalIndent(translations, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal %s translations: %w", lang, err)
		}

		if err := os.WriteFile(filename, data, 0644); err != nil {
			return fmt.Errorf("failed to write %s translation file: %w", lang, err)
		}
	}

	return nil
}

// GetSupportedLanguages returns list of supported languages
func (ts *TranslationService) GetSupportedLanguages() []string {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	var langs []string
	for lang := range ts.loadedLangs {
		langs = append(langs, lang)
	}

	return langs
}

// IsLanguageSupported checks if a language is supported
func (ts *TranslationService) IsLanguageSupported(language string) bool {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	return ts.loadedLangs[language]
}

// TranslationKey represents a translation key with metadata
type TranslationKey struct {
	Key         string   `json:"key"`
	Namespaces  []string `json:"namespaces"`
	Description string   `json:"description"`
	Required    bool     `json:"required"`
	Values      []string `json:"values,omitempty"` // For select/multiselect options
}

// KeyValidator validates translation keys against conventions
type KeyValidator struct {
	allowedNamespaces map[string]bool
	allowedPrefixes   map[string]bool
}

// NewKeyValidator creates a new key validator
func NewKeyValidator() *KeyValidator {
	return &KeyValidator{
		allowedNamespaces: map[string]bool{
			"forms":     true,
			"units":     true,
			"nav":       true,
			"errors":    true,
			"alerts":    true,
			"common":    true,
			"options":   true,
		},
		allowedPrefixes: map[string]bool{
			"forms.": true,
			"units.": true,
			"nav.":   true,
			"errors.": true,
			"alerts.": true,
		},
	}
}

// ValidateKey validates a translation key against conventions
func (kv *KeyValidator) ValidateKey(key string) error {
	// Key must contain a dot (namespace.key format)
	if !strings.Contains(key, ".") {
		return fmt.Errorf("key must be in format 'namespace.key': %s", key)
	}

	parts := strings.Split(key, ".")
	if len(parts) < 2 {
		return fmt.Errorf("key must be in format 'namespace.key': %s", key)
	}

	namespace := parts[0]
	
	// Check if namespace is allowed
	if !kv.allowedNamespaces[namespace] {
		return fmt.Errorf("invalid namespace '%s' in key '%s'. Allowed namespaces: %v", 
			namespace, key, kv.getAllowedNamespaces())
	}

	// Check key format
	if strings.Contains(key, " ") {
		return fmt.Errorf("key cannot contain spaces: %s", key)
	}

	if strings.Contains(key, "..") {
		return fmt.Errorf("key cannot contain consecutive dots: %s", key)
	}

	// Check for valid characters
	for _, char := range key {
		if !isValidKeyChar(char) {
			return fmt.Errorf("invalid character '%c' in key '%s'. Only lowercase letters, numbers, dots, and underscores allowed", 
				char, key)
		}
	}

	return nil
}

// isValidKeyChar checks if a character is valid in a translation key
func isValidKeyChar(char rune) bool {
	return (char >= 'a' && char <= 'z') || 
		   (char >= '0' && char <= '9') || 
		   char == '.' || 
		   char == '_' || 
		   char == '-'
}

// getAllowedNamespaces returns list of allowed namespaces
func (kv *KeyValidator) getAllowedNamespaces() []string {
	namespaces := make([]string, 0, len(kv.allowedNamespaces))
	for ns := range kv.allowedNamespaces {
		namespaces = append(namespaces, ns)
	}
	return namespaces
}

// TranslationValidator validates translation completeness
type TranslationValidator struct {
	service *TranslationService
	validator *KeyValidator
}

// NewTranslationValidator creates a new translation validator
func NewTranslationValidator(service *TranslationService) *TranslationValidator {
	return &TranslationValidator{
		service:  service,
		validator: NewKeyValidator(),
	}
}

// ValidateTranslations validates all translations
func (tv *TranslationValidator) ValidateTranslations() *ValidationResult {
	result := &ValidationResult{
		IsValid:    true,
		Errors:     []ValidationError{},
		Warnings:   []ValidationWarning{},
		Summary:    make(map[string]int),
	}

	// Check required languages are loaded
	requiredLangs := []string{"en", "bn"}
	for _, lang := range requiredLangs {
		if !tv.service.IsLanguageSupported(lang) {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Code:    "missing_language",
				Message: fmt.Sprintf("Required language '%s' is not loaded", lang),
				Language: lang,
			})
		}
	}

	// Validate all keys
	validationErrors := tv.service.ValidateAllKeys()
	for key, missingLangs := range validationErrors {
		// Validate key format
		if err := tv.validator.ValidateKey(key); err != nil {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Code:    "invalid_key_format",
				Message: err.Error(),
				Key:     key,
			})
			continue
		}

		// Check missing translations
		for _, lang := range missingLangs {
			result.IsValid = false
			result.Errors = append(result.Errors, ValidationError{
				Code:    "missing_translation",
				Message: fmt.Sprintf("Key '%s' is missing translation in language '%s'", key, lang),
				Key:     key,
				Language: lang,
			})
		}
	}

	// Update summary
	result.Summary["total_errors"] = len(result.Errors)
	result.Summary["total_warnings"] = len(result.Warnings)
	result.Summary["total_keys"] = len(tv.service.GetAllKeys("en"))
	result.Summary["supported_languages"] = len(tv.service.GetSupportedLanguages())

	return result
}

// ValidationResult represents validation results
type ValidationResult struct {
	IsValid    bool                `json:"is_valid"`
	Errors     []ValidationError    `json:"errors"`
	Warnings   []ValidationWarning  `json:"warnings"`
	Summary    map[string]int       `json:"summary"`
	Timestamp  string              `json:"timestamp"`
}

// ValidationError represents a validation error
type ValidationError struct {
	Code     string `json:"code"`
	Message  string `json:"message"`
	Key      string `json:"key,omitempty"`
	Language string `json:"language,omitempty"`
}

// ValidationWarning represents a validation warning
type ValidationWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Key     string `json:"key,omitempty"`
}

// Helper function to create translation service
func CreateTranslationService(basePath string) (*TranslationService, error) {
	service := NewTranslationService(basePath)
	
	if err := service.LoadTranslations(); err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	// Validate that required languages are loaded
	validator := NewTranslationValidator(service)
	result := validator.ValidateTranslations()
	
	if !result.IsValid {
		return nil, fmt.Errorf("translation validation failed: %d errors", len(result.Errors))
	}

	return service, nil
}

// Helper function to validate translation key format
func ValidateTranslationKey(key string) error {
	validator := NewKeyValidator()
	return validator.ValidateKey(key)
}

// Helper function to extract namespace from key
func ExtractNamespace(key string) string {
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

// Helper function to extract key part from namespace.key
func ExtractKeyPart(key string) string {
	parts := strings.Split(key, ".")
	if len(parts) > 1 {
		return strings.Join(parts[1:], ".")
	}
	return key
}

// Helper function to check if key is in namespace
func IsKeyInNamespace(key, namespace string) bool {
	return ExtractNamespace(key) == namespace
}

// Helper function to generate translation key
func GenerateTranslationKey(namespace, key string) string {
	return fmt.Sprintf("%s.%s", namespace, key)
}
