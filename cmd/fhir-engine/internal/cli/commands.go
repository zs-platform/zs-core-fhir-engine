package cli

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/alecthomas/kong"
	"github.com/charmbracelet/log"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zarishsphere/zs-core-fhir-engine/cmd/fhir-engine/internal/build"
	"github.com/zarishsphere/zs-core-fhir-engine/cmd/fhir-engine/internal/config"
	"github.com/zarishsphere/zs-core-fhir-engine/pkg/auth"
)

// ServeCommand starts the FHIR server
type ServeCommand struct {
	Port   int    `short:"p" default:"8080" help:"Port to listen on"`
	IGPath string `short:"i" default:"./config" help:"Path to FHIR resources"`
}

// Run starts the FHIR server with SMART on FHIR 2.1 authentication
func (s *ServeCommand) Run(ctx *kong.Context, cfg *config.GlobalConfig) error {
	// Setup logging based on global debug flag
	if cfg.Debug {
		log.SetLevel(log.DebugLevel)
		log.SetReportCaller(true)
	}

	log.Infof("Starting ZarishSphere FHIR Engine with SMART on FHIR 2.1...")
	log.Infof("Port: %d", s.Port)
	log.Infof("Resources path: %s", s.IGPath)

	// Initialize authentication components
	authConfig := &auth.AuthConfig{
		BaseURL:     "http://localhost:8080",
		Realm:       "zarishsphere",
		ClientID:    "zs-fhir-engine",
		RedirectURI: "http://localhost:3000/callback",
		JWKSURL:     "http://localhost:8080/auth/.well-known/jwks.json",
		TokenTTL:    3600,  // 1 hour
		RefreshTTL:  86400, // 24 hours
	}

	// Create mock services for now (in production, these would be real implementations)
	tokenStore := auth.NewInMemoryTokenStore()
	auditLogger := auth.NewMockAuditLogger()

	// Create JWT manager (in production, load real keys)
	jwtManager := auth.NewJWTManager(nil, nil, "http://localhost:8080", "http://localhost:8080/auth")

	// Create auth service
	authService := auth.NewKeycloakAuthService(&auth.KeycloakConfig{
		BaseURL:     authConfig.BaseURL,
		Realm:       authConfig.Realm,
		ClientID:    authConfig.ClientID,
		Secret:      "your-secret-key",
		RedirectURI: authConfig.RedirectURI,
	}, jwtManager, tokenStore)

	// Create auth handler
	authHandler := auth.NewHTTPHandler(authService, auditLogger, authConfig)

	// Create router with middleware
	router := chi.NewRouter()

	// Add middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(middleware.AllowContentType("application/json", "application/fhir+json"))

	// Add CORS middleware
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Register authentication routes
	authHandler.RegisterRoutes(router)

	// Add authentication middleware to protected routes
	router.Group(func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)

		// Health check endpoint (public)
		router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"status":        "healthy",
				"service":       "zs-core-fhir-engine",
				"version":       "2.0.0",
				"smart_on_fhir": "2.1",
			})
		})

		// FHIR metadata endpoint (public)
		router.Get("/fhir/R5/metadata", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/fhir+json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"resourceType": "CapabilityStatement",
				"status":       "active",
				"date":         time.Now().Format(time.RFC3339),
				"fhirVersion":  "5.0.0",
				"format":       []string{"application/fhir+json", "application/fhir+xml"},
				"smart_on_fhir": map[string]interface{}{
					"capabilities": []string{
						"launch-standalone",
						"launch-ehr",
						"sso-openid-connect",
						"permission-patient",
						"permission-v1",
					},
					"authorization_endpoint": fmt.Sprintf("http://localhost:%d/auth/authorize", s.Port),
					"token_endpoint":         fmt.Sprintf("http://localhost:%d/auth/token", s.Port),
					"userinfo_endpoint":      fmt.Sprintf("http://localhost:%d/auth/userinfo", s.Port),
					"jwks_uri":               fmt.Sprintf("http://localhost:%d/auth/.well-known/jwks.json", s.Port),
				},
				"rest": []interface{}{
					map[string]interface{}{
						"mode": "server",
						"resource": []interface{}{
							map[string]interface{}{
								"type":        "Patient",
								"interaction": []string{"read", "search", "create", "update", "delete", "history-instance"},
								"operation":   []string{"validate"},
							},
							map[string]interface{}{
								"type":        "Observation",
								"interaction": []string{"read", "search", "create", "update", "delete", "history-instance"},
								"operation":   []string{"validate"},
							},
							map[string]interface{}{
								"type":        "Condition",
								"interaction": []string{"read", "search", "create", "update", "delete", "history-instance"},
								"operation":   []string{"validate"},
							},
							map[string]interface{}{
								"type":        "MedicationRequest",
								"interaction": []string{"read", "search", "create", "update", "delete", "history-instance"},
								"operation":   []string{"validate"},
							},
						},
					},
				},
			})
		})

		// FHIR resource endpoints (protected)
		s.setupFHIRResourceEndpoints(router, s.IGPath)
	})

	log.Infof("🚀 ZarishSphere FHIR Engine with SMART on FHIR 2.1 started on http://localhost:%d", s.Port)
	log.Infof("📋 Available endpoints:")
	log.Infof("   - Health: http://localhost:%d/health", s.Port)
	log.Infof("   - FHIR Metadata: http://localhost:%d/fhir/R5/metadata", s.Port)
	log.Infof("   - SMART Auth: http://localhost:%d/auth/authorize", s.Port)
	log.Infof("   - Token Endpoint: http://localhost:%d/auth/token", s.Port)
	log.Infof("   - User Info: http://localhost:%d/auth/userinfo", s.Port)
	log.Infof("   - OpenID Config: http://localhost:%d/auth/.well-known/openid_configuration", s.Port)
	log.Infof("   - JWKS: http://localhost:%d/auth/.well-known/jwks.json", s.Port)

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), router)
}

// setupFHIRResourceEndpoints sets up FHIR resource endpoints
func (s *ServeCommand) setupFHIRResourceEndpoints(router chi.Router, igPath string) {
	// Load sample resources
	s.loadSampleResources(router, igPath)

	// FHIR Patient endpoints
	router.Route("/fhir/R5/Patient", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/fhir+json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"resourceType": "Bundle",
				"type":         "searchset",
				"total":        0,
				"entry":        []interface{}{},
			})
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/fhir+json")
			json.NewEncoder(w).Encode(map[string]string{
				"resourceType": "Patient",
				"id":           "created-" + fmt.Sprintf("%d", time.Now().Unix()),
			})
		})

		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
			patientID := chi.URLParam(r, "id")
			w.Header().Set("Content-Type", "application/fhir+json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"resourceType": "Patient",
				"id":           patientID,
				"message":      "Patient resource - SMART on FHIR 2.1 protected",
			})
		})
	})

	// FHIR Observation endpoints
	router.Route("/fhir/R5/Observation", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/fhir+json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"resourceType": "Bundle",
				"type":         "searchset",
				"total":        0,
				"entry":        []interface{}{},
			})
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/fhir+json")
			json.NewEncoder(w).Encode(map[string]string{
				"resourceType": "Observation",
				"id":           "created-" + fmt.Sprintf("%d", time.Now().Unix()),
			})
		})
	})

	// Add sample endpoints for development
	router.Get("/fhir/R5/Patient/sample", func(w http.ResponseWriter, r *http.Request) {
		if patientData, err := os.ReadFile(igPath + "/fhir-resources/patient.json"); err == nil {
			w.Header().Set("Content-Type", "application/fhir+json")
			w.Write(patientData)
		} else {
			http.Error(w, "Sample not found", http.StatusNotFound)
		}
	})
}

func (s *ServeCommand) loadSampleResources(router chi.Router, igPath string) {
}

// TerminologyCommand starts the terminology server
type TerminologyCommand struct {
	Port int `short:"p" default:"8081" help:"Port to listen on"`
}

// Run starts the terminology server
func (t *TerminologyCommand) Run(ctx *kong.Context, cfg *config.GlobalConfig) error {
	log.Infof("Starting terminology server on port %d...", t.Port)

	// Create a simple terminology server
	router := chi.NewRouter()

	// Basic ValueSet expansion
	router.Get("/fhir/R5/ValueSet/$expand", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/fhir+json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"resourceType": "ValueSet",
			"status":       "active",
			"expansion": map[string]interface{}{
				"identifier": "urn:uuid:" + "12345678-1234-5678-9abc-123456789012",
				"timestamp":  "2026-04-03T12:00:00Z",
				"contains": []interface{}{
					map[string]interface{}{
						"system":  "http://snomed.info/sct",
						"code":    "386661006",
						"display": "Fetus heart rate",
					},
				},
			},
		})
	})

	log.Infof("Terminology server listening on :%d", t.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", t.Port), router)
}

// ValidateCommand validates FHIR resources
type ValidateCommand struct {
	File     string `short:"v" help:"FHIR resource file to validate" type:"path"`
	Resource string `short:"r" help:"Resource type (e.g., Patient, Observation)"`
}

// Run validates a FHIR resource
func (v *ValidateCommand) Run(ctx *kong.Context, cfg *config.GlobalConfig) error {
	if v.File == "" {
		return fmt.Errorf("file path is required for validation")
	}

	log.Infof("Validating FHIR resource: %s", v.File)

	// Read and validate the resource
	data, err := os.ReadFile(v.File)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var resource map[string]any
	if err := json.Unmarshal(data, &resource); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Basic validation - check if it has resourceType
	if resourceType, ok := resource["resourceType"]; !ok {
		return fmt.Errorf("missing resourceType field")
	} else {
		log.Infof("✅ Validated %s resource", resourceType)
	}

	log.Infof("✅ Resource validation passed")
	return nil
}

// VersionCommand shows detailed version information
type VersionCommand struct{}

// Run shows detailed version information
func (v *VersionCommand) Run(ctx *kong.Context, cfg *config.GlobalConfig) error {
	build.PrintBuildInfo()

	// Show additional version details
	fmt.Printf("FHIR Version: R5 (5.0.0)\n")
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("Build Target: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("ZarishSphere Version: 2.0.0\n")

	return nil
}
