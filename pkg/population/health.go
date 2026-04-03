package population

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

// HealthManager manages population health data and reporting
type HealthManager struct {
	store  PopulationStore
	config HealthConfig
}

// PopulationStore defines the storage interface for population health data
type PopulationStore interface {
	GetPatients(ctx context.Context, criteria PatientCriteria) ([]Patient, error)
	GetConditions(ctx context.Context, criteria ConditionCriteria) ([]ConditionSummary, error)
	GetEncounters(ctx context.Context, criteria EncounterCriteria) ([]EncounterSummary, error)
	GetQualityMetrics(ctx context.Context, criteria QualityCriteria) ([]QualityMetric, error)
	StoreReport(ctx context.Context, report HealthReport) error
	GetReports(ctx context.Context, tenantID string) ([]HealthReport, error)
}

// HealthConfig contains population health configuration
type HealthConfig struct {
	Enabled           bool
	ReportRetentionDays int
	DefaultMeasureSet string
}

// Patient represents a patient in population health
type Patient struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Age         int                    `json:"age"`
	Gender      string                 `json:"gender"`
	Conditions  []string               `json:"conditions,omitempty"`
	RiskScore   float64                `json:"riskScore"`
	LastVisit   *time.Time             `json:"lastVisit,omitempty"`
	TenantID    string                 `json:"tenantId"`
	Location    string                 `json:"location,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PatientCriteria defines criteria for patient selection
type PatientCriteria struct {
	TenantID    string
	AgeMin      int
	AgeMax      int
	Gender      string
	Conditions  []string
	RiskMin     float64
	RiskMax     float64
	Location    string
	Limit       int
}

// ConditionSummary represents a condition summary
type ConditionSummary struct {
	Condition     string  `json:"condition"`
	ICD10Code     string  `json:"icd10Code"`
	PatientCount  int     `json:"patientCount"`
	Prevalence    float64 `json:"prevalence"`
	AvgAge        float64 `json:"avgAge"`
	GenderSplit   map[string]int `json:"genderSplit"`
}

// ConditionCriteria defines criteria for condition selection
type ConditionCriteria struct {
	TenantID   string
	Conditions []string
	TimeRange  TimeRange
}

// EncounterSummary represents an encounter summary
type EncounterSummary struct {
	Type         string  `json:"type"`
	Count        int     `json:"count"`
	AvgDuration  float64 `json:"avgDuration"` // minutes
	TotalCost    float64 `json:"totalCost"`
}

// EncounterCriteria defines criteria for encounter selection
type EncounterCriteria struct {
	TenantID  string
	Types     []string
	TimeRange TimeRange
}

// QualityMetric represents a quality measure metric
type QualityMetric struct {
	MeasureID   string  `json:"measureId"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Target      float64 `json:"target"`
	Actual      float64 `json:"actual"`
	Denominator int     `json:"denominator"`
	Numerator   int     `json:"numerator"`
	Compliance  float64 `json:"compliance"`
	Trend       string  `json:"trend"` // improving, stable, declining
}

// QualityCriteria defines criteria for quality metrics
type QualityCriteria struct {
	TenantID    string
	Measures    []string
	TimeRange   TimeRange
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// NewHealthManager creates a new population health manager
func NewHealthManager(store PopulationStore, config HealthConfig) *HealthManager {
	return &HealthManager{
		store:  store,
		config: config,
	}
}

// GetPatientRegistry returns the patient registry for a population
func (hm *HealthManager) GetPatientRegistry(ctx context.Context, criteria PatientCriteria) (*PatientRegistry, error) {
	patients, err := hm.store.GetPatients(ctx, criteria)
	if err != nil {
		return nil, err
	}
	
	// Calculate registry statistics
	stats := calculateRegistryStats(patients)
	
	return &PatientRegistry{
		Criteria:   criteria,
		GeneratedAt: time.Now(),
		TotalPatients: len(patients),
		Patients:   patients,
		Statistics: stats,
	}, nil
}

// GetConditionReport returns a condition prevalence report
func (hm *HealthManager) GetConditionReport(ctx context.Context, criteria ConditionCriteria) (*ConditionReport, error) {
	conditions, err := hm.store.GetConditions(ctx, criteria)
	if err != nil {
		return nil, err
	}
	
	// Sort by patient count
	sort.Slice(conditions, func(i, j int) bool {
		return conditions[i].PatientCount > conditions[j].PatientCount
	})
	
	return &ConditionReport{
		Criteria:     criteria,
		GeneratedAt:  time.Now(),
		Conditions:   conditions,
		TotalPatients: sumPatientCounts(conditions),
	}, nil
}

// GetUtilizationReport returns a care utilization report
func (hm *HealthManager) GetUtilizationReport(ctx context.Context, criteria EncounterCriteria) (*UtilizationReport, error) {
	encounters, err := hm.store.GetEncounters(ctx, criteria)
	if err != nil {
		return nil, err
	}
	
	// Calculate utilization metrics
	totalEncounters := 0
	totalCost := 0.0
	
	for _, e := range encounters {
		totalEncounters += e.Count
		totalCost += e.TotalCost
	}
	
	return &UtilizationReport{
		Criteria:        criteria,
		GeneratedAt:     time.Now(),
		Encounters:      encounters,
		TotalEncounters: totalEncounters,
		TotalCost:       totalCost,
		TimeRange:       criteria.TimeRange,
	}, nil
}

// GetQualityReport returns a quality measures report
func (hm *HealthManager) GetQualityReport(ctx context.Context, criteria QualityCriteria) (*QualityReport, error) {
	metrics, err := hm.store.GetQualityMetrics(ctx, criteria)
	if err != nil {
		return nil, err
	}
	
	// Calculate overall score
	overallScore := calculateOverallQualityScore(metrics)
	
	return &QualityReport{
		Criteria:      criteria,
		GeneratedAt:   time.Now(),
		Metrics:       metrics,
		OverallScore:  overallScore,
		MeasureCount:  len(metrics),
	}, nil
}

// GenerateHealthReport generates a comprehensive health report
func (hm *HealthManager) GenerateHealthReport(ctx context.Context, tenantID string, timeRange TimeRange) (*HealthReport, error) {
	// Get patient registry
	patientCriteria := PatientCriteria{
		TenantID: tenantID,
	}
	
	registry, err := hm.GetPatientRegistry(ctx, patientCriteria)
	if err != nil {
		return nil, err
	}
	
	// Get condition report
	conditionCriteria := ConditionCriteria{
		TenantID:  tenantID,
		TimeRange: timeRange,
	}
	
	conditionReport, err := hm.GetConditionReport(ctx, conditionCriteria)
	if err != nil {
		return nil, err
	}
	
	// Get utilization report
	encounterCriteria := EncounterCriteria{
		TenantID:  tenantID,
		TimeRange: timeRange,
	}
	
	utilizationReport, err := hm.GetUtilizationReport(ctx, encounterCriteria)
	if err != nil {
		return nil, err
	}
	
	// Get quality report
	qualityCriteria := QualityCriteria{
		TenantID:  tenantID,
		TimeRange: timeRange,
	}
	
	qualityReport, err := hm.GetQualityReport(ctx, qualityCriteria)
	if err != nil {
		return nil, err
	}
	
	report := &HealthReport{
		ID:                generateReportID(),
		TenantID:          tenantID,
		GeneratedAt:       time.Now(),
		TimeRange:         timeRange,
		PatientRegistry:   registry,
		ConditionReport:   conditionReport,
		UtilizationReport: utilizationReport,
		QualityReport:     qualityReport,
	}
	
	// Store report
	if err := hm.store.StoreReport(ctx, *report); err != nil {
		return nil, err
	}
	
	return report, nil
}

// IdentifyCareGaps identifies care gaps in the population
func (hm *HealthManager) IdentifyCareGaps(ctx context.Context, tenantID string) ([]CareGap, error) {
	// Get patient registry
	registry, err := hm.GetPatientRegistry(ctx, PatientCriteria{TenantID: tenantID})
	if err != nil {
		return nil, err
	}
	
	gaps := make([]CareGap, 0)
	
	// Identify various care gaps
	for _, patient := range registry.Patients {
		// Check for overdue visits
		if patient.LastVisit == nil || time.Since(*patient.LastVisit) > 365*24*time.Hour {
			gaps = append(gaps, CareGap{
				PatientID:     patient.ID,
				PatientName:   patient.Name,
				Type:          "overdue_visit",
				Description:   "Annual wellness visit overdue",
				Priority:      "medium",
				DaysOverdue:   int(time.Since(*patient.LastVisit).Hours() / 24),
			})
		}
		
		// Check for uncontrolled conditions
		for _, condition := range patient.Conditions {
			if isUncontrolled(condition) {
				gaps = append(gaps, CareGap{
					PatientID:     patient.ID,
					PatientName:   patient.Name,
					Type:          "uncontrolled_condition",
					Description:   fmt.Sprintf("Uncontrolled %s", condition),
					Priority:      "high",
					Condition:     condition,
				})
			}
		}
		
		// Check for high risk patients without care plan
		if patient.RiskScore > 0.7 {
			hasCarePlan := false // Check if patient has active care plan
			if !hasCarePlan {
				gaps = append(gaps, CareGap{
					PatientID:     patient.ID,
					PatientName:   patient.Name,
					Type:          "missing_care_plan",
					Description:   "High-risk patient without care plan",
					Priority:      "high",
					RiskScore:     patient.RiskScore,
				})
			}
		}
	}
	
	return gaps, nil
}

// CalculateRiskScores calculates risk scores for a population
func (hm *HealthManager) CalculateRiskScores(ctx context.Context, criteria PatientCriteria) ([]RiskAssessment, error) {
	patients, err := hm.store.GetPatients(ctx, criteria)
	if err != nil {
		return nil, err
	}
	
	assessments := make([]RiskAssessment, 0)
	
	for _, patient := range patients {
		score := calculateRiskScore(patient)
		
		assessments = append(assessments, RiskAssessment{
			PatientID:   patient.ID,
			PatientName: patient.Name,
			RiskScore:   score,
			RiskLevel:   getRiskLevel(score),
			Factors:     getRiskFactors(patient),
			CalculatedAt: time.Now(),
		})
	}
	
	return assessments, nil
}

// Helper functions

func calculateRegistryStats(patients []Patient) RegistryStatistics {
	if len(patients) == 0 {
		return RegistryStatistics{}
	}
	
	var totalAge, totalRisk float64
	genderSplit := make(map[string]int)
	conditionCounts := make(map[string]int)
	
	for _, p := range patients {
		totalAge += float64(p.Age)
		totalRisk += p.RiskScore
		genderSplit[p.Gender]++
		
		for _, c := range p.Conditions {
			conditionCounts[c]++
		}
	}
	
	return RegistryStatistics{
		AvgAge:          totalAge / float64(len(patients)),
		AvgRiskScore:    totalRisk / float64(len(patients)),
		GenderSplit:     genderSplit,
		TopConditions:   getTopConditions(conditionCounts, 5),
	}
}

func getTopConditions(counts map[string]int, limit int) []ConditionCount {
	type pair struct {
		condition string
		count     int
	}
	
	pairs := make([]pair, 0, len(counts))
	for c, count := range counts {
		pairs = append(pairs, pair{c, count})
	}
	
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	
	if limit > len(pairs) {
		limit = len(pairs)
	}
	
	result := make([]ConditionCount, limit)
	for i := 0; i < limit; i++ {
		result[i] = ConditionCount{
			Condition: pairs[i].condition,
			Count:     pairs[i].count,
		}
	}
	
	return result
}

func sumPatientCounts(conditions []ConditionSummary) int {
	seen := make(map[string]bool)
	total := 0
	
	for _, c := range conditions {
		if !seen[c.Condition] {
			seen[c.Condition] = true
			total += c.PatientCount
		}
	}
	
	return total
}

func calculateOverallQualityScore(metrics []QualityMetric) float64 {
	if len(metrics) == 0 {
		return 0
	}
	
	var total float64
	for _, m := range metrics {
		total += m.Compliance
	}
	
	return total / float64(len(metrics))
}

func isUncontrolled(condition string) bool {
	// Simplified check - in production would use clinical rules
	uncontrolledConditions := []string{"uncontrolled_diabetes", "uncontrolled_hypertension", "severe_asthma"}
	for _, c := range uncontrolledConditions {
		if c == condition {
			return true
		}
	}
	return false
}

func calculateRiskScore(patient Patient) float64 {
	// Simplified risk score calculation
	score := 0.0
	
	// Age factor
	if patient.Age > 65 {
		score += 0.2
	}
	
	// Condition factor
	score += float64(len(patient.Conditions)) * 0.1
	
	// Cap at 1.0
	if score > 1.0 {
		score = 1.0
	}
	
	return score
}

func getRiskLevel(score float64) string {
	if score < 0.3 {
		return "low"
	} else if score < 0.7 {
		return "medium"
	}
	return "high"
}

func getRiskFactors(patient Patient) []string {
	factors := make([]string, 0)
	
	if patient.Age > 65 {
		factors = append(factors, "advanced_age")
	}
	
	if len(patient.Conditions) > 2 {
		factors = append(factors, "multiple_conditions")
	}
	
	return factors
}

func generateReportID() string {
	return fmt.Sprintf("rpt-%d", time.Now().UnixNano())
}

// PatientRegistry contains patient registry data
type PatientRegistry struct {
	Criteria      PatientCriteria    `json:"criteria"`
	GeneratedAt   time.Time          `json:"generatedAt"`
	TotalPatients int                `json:"totalPatients"`
	Patients      []Patient          `json:"patients,omitempty"`
	Statistics    RegistryStatistics `json:"statistics"`
}

// RegistryStatistics contains registry statistics
type RegistryStatistics struct {
	AvgAge        float64            `json:"avgAge"`
	AvgRiskScore  float64            `json:"avgRiskScore"`
	GenderSplit   map[string]int     `json:"genderSplit"`
	TopConditions []ConditionCount   `json:"topConditions"`
}

// ConditionCount represents a condition count
type ConditionCount struct {
	Condition string `json:"condition"`
	Count     int    `json:"count"`
}

// ConditionReport contains condition prevalence data
type ConditionReport struct {
	Criteria      ConditionCriteria  `json:"criteria"`
	GeneratedAt   time.Time          `json:"generatedAt"`
	Conditions    []ConditionSummary `json:"conditions"`
	TotalPatients int                `json:"totalPatients"`
}

// UtilizationReport contains care utilization data
type UtilizationReport struct {
	Criteria        EncounterCriteria  `json:"criteria"`
	GeneratedAt     time.Time          `json:"generatedAt"`
	Encounters      []EncounterSummary `json:"encounters"`
	TotalEncounters int                `json:"totalEncounters"`
	TotalCost       float64            `json:"totalCost"`
	TimeRange       TimeRange          `json:"timeRange"`
}

// QualityReport contains quality measures data
type QualityReport struct {
	Criteria     QualityCriteria `json:"criteria"`
	GeneratedAt  time.Time       `json:"generatedAt"`
	Metrics      []QualityMetric `json:"metrics"`
	OverallScore float64         `json:"overallScore"`
	MeasureCount int             `json:"measureCount"`
}

// HealthReport contains comprehensive health data
type HealthReport struct {
	ID                string             `json:"id"`
	TenantID          string             `json:"tenantId"`
	GeneratedAt       time.Time          `json:"generatedAt"`
	TimeRange         TimeRange          `json:"timeRange"`
	PatientRegistry   *PatientRegistry   `json:"patientRegistry,omitempty"`
	ConditionReport   *ConditionReport   `json:"conditionReport,omitempty"`
	UtilizationReport *UtilizationReport `json:"utilizationReport,omitempty"`
	QualityReport     *QualityReport     `json:"qualityReport,omitempty"`
}

// CareGap represents an identified care gap
type CareGap struct {
	PatientID   string  `json:"patientId"`
	PatientName string  `json:"patientName"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	DaysOverdue int     `json:"daysOverdue,omitempty"`
	Condition   string  `json:"condition,omitempty"`
	RiskScore   float64 `json:"riskScore,omitempty"`
}

// RiskAssessment represents a patient risk assessment
type RiskAssessment struct {
	PatientID    string   `json:"patientId"`
	PatientName  string   `json:"patientName"`
	RiskScore    float64  `json:"riskScore"`
	RiskLevel    string   `json:"riskLevel"`
	Factors      []string `json:"factors"`
	CalculatedAt time.Time `json:"calculatedAt"`
}

// PopulationHandler handles HTTP requests for population health
type PopulationHandler struct {
	manager *HealthManager
}

// NewPopulationHandler creates a new population health HTTP handler
func NewPopulationHandler(manager *HealthManager) *PopulationHandler {
	return &PopulationHandler{
		manager: manager,
	}
}

// RegisterRoutes registers population health endpoints
func (h *PopulationHandler) RegisterRoutes(router chi.Router) {
	router.Get("/population/registry", h.handleGetRegistry)
	router.Get("/population/conditions", h.handleGetConditions)
	router.Get("/population/utilization", h.handleGetUtilization)
	router.Get("/population/quality", h.handleGetQuality)
	router.Post("/population/report", h.handleGenerateReport)
	router.Get("/population/caregaps", h.handleGetCareGaps)
	router.Get("/population/risk", h.handleCalculateRisk)
}

// handleGetRegistry handles GET /population/registry
func (h *PopulationHandler) handleGetRegistry(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	criteria := PatientCriteria{
		TenantID: r.Header.Get("X-Tenant-ID"),
	}
	
	registry, err := h.manager.GetPatientRegistry(ctx, criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(registry)
}

// handleGetConditions handles GET /population/conditions
func (h *PopulationHandler) handleGetConditions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	criteria := ConditionCriteria{
		TenantID: r.Header.Get("X-Tenant-ID"),
	}
	
	report, err := h.manager.GetConditionReport(ctx, criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// handleGetUtilization handles GET /population/utilization
func (h *PopulationHandler) handleGetUtilization(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	criteria := EncounterCriteria{
		TenantID: r.Header.Get("X-Tenant-ID"),
		TimeRange: TimeRange{
			Start: time.Now().AddDate(0, -1, 0),
			End:   time.Now(),
		},
	}
	
	report, err := h.manager.GetUtilizationReport(ctx, criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// handleGetQuality handles GET /population/quality
func (h *PopulationHandler) handleGetQuality(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	criteria := QualityCriteria{
		TenantID: r.Header.Get("X-Tenant-ID"),
		TimeRange: TimeRange{
			Start: time.Now().AddDate(0, -3, 0),
			End:   time.Now(),
		},
	}
	
	report, err := h.manager.GetQualityReport(ctx, criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// handleGenerateReport handles POST /population/report
func (h *PopulationHandler) handleGenerateReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")
	
	var req struct {
		From time.Time `json:"from"`
		To   time.Time `json:"to"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	report, err := h.manager.GenerateHealthReport(ctx, tenantID, TimeRange{
		Start: req.From,
		End:   req.To,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

// handleGetCareGaps handles GET /population/caregaps
func (h *PopulationHandler) handleGetCareGaps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID := r.Header.Get("X-Tenant-ID")
	
	gaps, err := h.manager.IdentifyCareGaps(ctx, tenantID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(gaps)
}

// handleCalculateRisk handles GET /population/risk
func (h *PopulationHandler) handleCalculateRisk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	criteria := PatientCriteria{
		TenantID: r.Header.Get("X-Tenant-ID"),
	}
	
	assessments, err := h.manager.CalculateRiskScores(ctx, criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(assessments)
}

// InMemoryPopulationStore implements PopulationStore with in-memory storage
type InMemoryPopulationStore struct {
	patients map[string]Patient
	reports  map[string]HealthReport
	mu       sync.RWMutex
}

// NewInMemoryPopulationStore creates a new in-memory population store
func NewInMemoryPopulationStore() *InMemoryPopulationStore {
	return &InMemoryPopulationStore{
		patients: make(map[string]Patient),
		reports:  make(map[string]HealthReport),
	}
}

// GetPatients implements PopulationStore
func (s *InMemoryPopulationStore) GetPatients(ctx context.Context, criteria PatientCriteria) ([]Patient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []Patient
	
	for _, patient := range s.patients {
		if criteria.TenantID != "" && patient.TenantID != criteria.TenantID {
			continue
		}
		if criteria.Gender != "" && patient.Gender != criteria.Gender {
			continue
		}
		if criteria.AgeMin > 0 && patient.Age < criteria.AgeMin {
			continue
		}
		if criteria.AgeMax > 0 && patient.Age > criteria.AgeMax {
			continue
		}
		
		results = append(results, patient)
	}
	
	return results, nil
}

// GetConditions implements PopulationStore
func (s *InMemoryPopulationStore) GetConditions(ctx context.Context, criteria ConditionCriteria) ([]ConditionSummary, error) {
	return []ConditionSummary{}, nil
}

// GetEncounters implements PopulationStore
func (s *InMemoryPopulationStore) GetEncounters(ctx context.Context, criteria EncounterCriteria) ([]EncounterSummary, error) {
	return []EncounterSummary{}, nil
}

// GetQualityMetrics implements PopulationStore
func (s *InMemoryPopulationStore) GetQualityMetrics(ctx context.Context, criteria QualityCriteria) ([]QualityMetric, error) {
	return []QualityMetric{}, nil
}

// StoreReport implements PopulationStore
func (s *InMemoryPopulationStore) StoreReport(ctx context.Context, report HealthReport) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.reports[report.ID] = report
	return nil
}

// GetReports implements PopulationStore
func (s *InMemoryPopulationStore) GetReports(ctx context.Context, tenantID string) ([]HealthReport, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []HealthReport
	for _, report := range s.reports {
		if report.TenantID == tenantID {
			results = append(results, report)
		}
	}
	
	return results, nil
}
