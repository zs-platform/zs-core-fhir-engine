package predictive

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

// AnalyticsEngine provides predictive analytics capabilities
type AnalyticsEngine struct {
	dataStore DataStore
	models    map[string]PredictiveModel
	mu        sync.RWMutex
}

// DataStore defines the data storage interface
type DataStore interface {
	GetPatientHistory(ctx context.Context, patientID string) ([]DataPoint, error)
	GetPopulationData(ctx context.Context, criteria PopulationCriteria) ([]DataPoint, error)
	StoreTrend(ctx context.Context, trend Trend) error
	GetTrends(ctx context.Context, patientID string) ([]Trend, error)
}

// PredictiveModel defines the interface for predictive models
type PredictiveModel interface {
	Predict(ctx context.Context, data []DataPoint) (*Prediction, error)
	GetType() string
}

// DataPoint represents a single data point
type DataPoint struct {
	Timestamp   time.Time              `json:"timestamp"`
	PatientID   string                 `json:"patientId"`
	Type        string                 `json:"type"`
	Value       float64                `json:"value"`
	Unit        string                 `json:"unit,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// PopulationCriteria defines criteria for population selection
type PopulationCriteria struct {
	AgeRange    [2]int
	Gender      string
	Conditions  []string
	Location    string
	TenantID    string
	TimeRange   TimeRange
}

// TimeRange represents a time range
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Prediction represents a prediction result
type Prediction struct {
	Type        string                 `json:"type"`
	Value       float64                `json:"value"`
	Confidence  float64                `json:"confidence"`
	LowerBound  float64                `json:"lowerBound"`
	UpperBound  float64                `json:"upperBound"`
	Timestamp   time.Time              `json:"timestamp"`
	Factors     []RiskFactor           `json:"factors,omitempty"`
	Explanation string                 `json:"explanation,omitempty"`
}

// RiskFactor represents a contributing risk factor
type RiskFactor struct {
	Name        string  `json:"name"`
	Impact      float64 `json:"impact"` // -1 to 1
	Description string  `json:"description"`
}

// Trend represents a trend in health data
type Trend struct {
	ID          string    `json:"id"`
	PatientID   string    `json:"patientId"`
	Type        string    `json:"type"`
	Direction   string    `json:"direction"` // improving, worsening, stable
	StartValue  float64   `json:"startValue"`
	CurrentValue float64  `json:"currentValue"`
	ChangeRate  float64   `json:"changeRate"`
	StartDate   time.Time `json:"startDate"`
	EndDate     time.Time `json:"endDate"`
}

// NewAnalyticsEngine creates a new analytics engine
func NewAnalyticsEngine(dataStore DataStore) *AnalyticsEngine {
	return &AnalyticsEngine{
		dataStore: dataStore,
		models:    make(map[string]PredictiveModel),
	}
}

// RegisterModel registers a predictive model
func (ae *AnalyticsEngine) RegisterModel(model PredictiveModel) {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	
	ae.models[model.GetType()] = model
}

// PredictRisk predicts health risks for a patient
func (ae *AnalyticsEngine) PredictRisk(ctx context.Context, patientID string, riskType string) (*Prediction, error) {
	// Get patient history
	history, err := ae.dataStore.GetPatientHistory(ctx, patientID)
	if err != nil {
		return nil, fmt.Errorf("failed to get patient history: %w", err)
	}
	
	// Get model
	ae.mu.RLock()
	model, exists := ae.models[riskType]
	ae.mu.RUnlock()
	
	if !exists {
		return nil, fmt.Errorf("no model found for risk type: %s", riskType)
	}
	
	// Make prediction
	return model.Predict(ctx, history)
}

// AnalyzeTrends analyzes trends in patient data
func (ae *AnalyticsEngine) AnalyzeTrends(ctx context.Context, patientID string, metricType string) (*TrendAnalysis, error) {
	// Get patient history
	history, err := ae.dataStore.GetPatientHistory(ctx, patientID)
	if err != nil {
		return nil, err
	}
	
	// Filter by metric type
	var filtered []DataPoint
	for _, dp := range history {
		if dp.Type == metricType {
			filtered = append(filtered, dp)
		}
	}
	
	if len(filtered) < 2 {
		return nil, fmt.Errorf("insufficient data for trend analysis")
	}
	
	// Sort by timestamp
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})
	
	// Calculate trend
	startValue := filtered[0].Value
	currentValue := filtered[len(filtered)-1].Value
	change := currentValue - startValue
	
	// Determine direction
	direction := "stable"
	if change > 0.1 {
		direction = "worsening"
		if isPositiveMetric(metricType) {
			direction = "improving"
		}
	} else if change < -0.1 {
		direction = "improving"
		if isPositiveMetric(metricType) {
			direction = "worsening"
		}
	}
	
	// Calculate change rate
	timeSpan := filtered[len(filtered)-1].Timestamp.Sub(filtered[0].Timestamp).Hours() / 24
	changeRate := 0.0
	if timeSpan > 0 {
		changeRate = change / timeSpan
	}
	
	trend := Trend{
		ID:           generateTrendID(),
		PatientID:    patientID,
		Type:         metricType,
		Direction:    direction,
		StartValue:   startValue,
		CurrentValue: currentValue,
		ChangeRate:   changeRate,
		StartDate:    filtered[0].Timestamp,
		EndDate:      filtered[len(filtered)-1].Timestamp,
	}
	
	// Store trend
	if err := ae.dataStore.StoreTrend(ctx, trend); err != nil {
		return nil, err
	}
	
	return &TrendAnalysis{
		Trend:        trend,
		DataPoints:   len(filtered),
		TimeSpan:     timeSpan,
		Volatility:   calculateVolatility(filtered),
		Forecast:     ae.forecastNextValue(filtered),
	}, nil
}

// PopulationAnalysis analyzes population health trends
func (ae *AnalyticsEngine) PopulationAnalysis(ctx context.Context, criteria PopulationCriteria) (*PopulationReport, error) {
	// Get population data
	data, err := ae.dataStore.GetPopulationData(ctx, criteria)
	if err != nil {
		return nil, err
	}
	
	// Calculate statistics
	stats := calculatePopulationStats(data)
	
	// Identify trends
	trends := identifyPopulationTrends(data)
	
	// Risk stratification
	riskStrat := stratifyRisk(data)
	
	return &PopulationReport{
		Criteria:     criteria,
		GeneratedAt:  time.Now(),
		TotalPatients: stats.Total,
		Statistics:   stats,
		Trends:       trends,
		RiskStratification: riskStrat,
	}, nil
}

// isPositiveMetric returns true if higher values are better
func isPositiveMetric(metricType string) bool {
	positiveMetrics := []string{"exercise", "sleep_quality", "medication_adherence"}
	for _, m := range positiveMetrics {
		if m == metricType {
			return true
		}
	}
	return false
}

// calculateVolatility calculates data volatility
func calculateVolatility(data []DataPoint) float64 {
	if len(data) < 2 {
		return 0
	}
	
	var sum, mean, variance float64
	
	for _, dp := range data {
		sum += dp.Value
	}
	mean = sum / float64(len(data))
	
	for _, dp := range data {
		variance += (dp.Value - mean) * (dp.Value - mean)
	}
	
	return variance / float64(len(data))
}

// forecastNextValue forecasts the next value
func (ae *AnalyticsEngine) forecastNextValue(data []DataPoint) float64 {
	if len(data) < 2 {
		return 0
	}
	
	// Simple linear regression
	n := float64(len(data))
	var sumX, sumY, sumXY, sumX2 float64
	
	for i, dp := range data {
		x := float64(i)
		y := dp.Value
		sumX += x
		sumY += y
		sumXY += x * y
		sumX2 += x * x
	}
	
	slope := (n*sumXY - sumX*sumY) / (n*sumX2 - sumX*sumX)
	intercept := (sumY - slope*sumX) / n
	
	// Predict next value
	nextX := n
	return slope*nextX + intercept
}

// calculatePopulationStats calculates population statistics
func calculatePopulationStats(data []DataPoint) PopulationStats {
	if len(data) == 0 {
		return PopulationStats{}
	}
	
	var sum, min, max float64
	min = data[0].Value
	max = data[0].Value
	
	for _, dp := range data {
		sum += dp.Value
		if dp.Value < min {
			min = dp.Value
		}
		if dp.Value > max {
			max = dp.Value
		}
	}
	
	mean := sum / float64(len(data))
	
	return PopulationStats{
		Total: len(data),
		Mean:  mean,
		Min:   min,
		Max:   max,
	}
}

// identifyPopulationTrends identifies trends in population data
func identifyPopulationTrends(data []DataPoint) []TrendSummary {
	// Group by metric type
	byType := make(map[string][]DataPoint)
	for _, dp := range data {
		byType[dp.Type] = append(byType[dp.Type], dp)
	}
	
	var trends []TrendSummary
	
	for metricType, points := range byType {
		if len(points) < 2 {
			continue
		}
		
		// Sort by timestamp
		sort.Slice(points, func(i, j int) bool {
			return points[i].Timestamp.Before(points[j].Timestamp)
		})
		
		change := points[len(points)-1].Value - points[0].Value
		
		trends = append(trends, TrendSummary{
			MetricType: metricType,
			Change:     change,
			Direction:  determineDirection(change, metricType),
		})
	}
	
	return trends
}

// stratifyRisk stratifies population by risk level
func stratifyRisk(data []DataPoint) RiskStratification {
	low, medium, high := 0, 0, 0
	
	for _, dp := range data {
		switch dp.Metadata["risk_level"].(string) {
		case "low":
			low++
		case "medium":
			medium++
		case "high":
			high++
		}
	}
	
	return RiskStratification{
		LowRisk:    low,
		MediumRisk: medium,
		HighRisk:   high,
	}
}

// determineDirection determines trend direction
func determineDirection(change float64, metricType string) string {
	isPositive := isPositiveMetric(metricType)
	
	if change > 0 {
		if isPositive {
			return "improving"
		}
		return "worsening"
	} else if change < 0 {
		if isPositive {
			return "worsening"
		}
		return "improving"
	}
	
	return "stable"
}

// generateTrendID generates a unique trend ID
func generateTrendID() string {
	return fmt.Sprintf("trend-%d", time.Now().UnixNano())
}

// TrendAnalysis contains trend analysis results
type TrendAnalysis struct {
	Trend      Trend   `json:"trend"`
	DataPoints int     `json:"dataPoints"`
	TimeSpan   float64 `json:"timeSpan"` // days
	Volatility float64 `json:"volatility"`
	Forecast   float64 `json:"forecast"`
}

// PopulationReport contains population health analysis
type PopulationReport struct {
	Criteria          PopulationCriteria `json:"criteria"`
	GeneratedAt     time.Time          `json:"generatedAt"`
	TotalPatients   int                `json:"totalPatients"`
	Statistics      PopulationStats    `json:"statistics"`
	Trends          []TrendSummary     `json:"trends"`
	RiskStratification RiskStratification `json:"riskStratification"`
}

// PopulationStats contains population statistics
type PopulationStats struct {
	Total int     `json:"total"`
	Mean  float64 `json:"mean"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

// TrendSummary contains trend summary
type TrendSummary struct {
	MetricType string  `json:"metricType"`
	Change     float64 `json:"change"`
	Direction  string  `json:"direction"`
}

// RiskStratification contains risk stratification
type RiskStratification struct {
	LowRisk    int `json:"lowRisk"`
	MediumRisk int `json:"mediumRisk"`
	HighRisk   int `json:"highRisk"`
}

// PredictiveHandler handles HTTP requests for predictive analytics
type PredictiveHandler struct {
	engine *AnalyticsEngine
}

// NewPredictiveHandler creates a new predictive analytics HTTP handler
func NewPredictiveHandler(engine *AnalyticsEngine) *PredictiveHandler {
	return &PredictiveHandler{
		engine: engine,
	}
}

// RegisterRoutes registers predictive analytics endpoints
func (h *PredictiveHandler) RegisterRoutes(router chi.Router) {
	router.Post("/predictive/risk/{patientID}", h.handlePredictRisk)
	router.Get("/predictive/trends/{patientID}", h.handleAnalyzeTrends)
	router.Post("/predictive/population", h.handlePopulationAnalysis)
}

// handlePredictRisk handles POST /predictive/risk/{patientID}
func (h *PredictiveHandler) handlePredictRisk(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	patientID := chi.URLParam(r, "patientID")
	
	var req struct {
		RiskType string `json:"riskType"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	prediction, err := h.engine.PredictRisk(ctx, patientID, req.RiskType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(prediction)
}

// handleAnalyzeTrends handles GET /predictive/trends/{patientID}
func (h *PredictiveHandler) handleAnalyzeTrends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	patientID := chi.URLParam(r, "patientID")
	
	metricType := r.URL.Query().Get("metric")
	if metricType == "" {
		metricType = "blood_pressure"
	}
	
	analysis, err := h.engine.AnalyzeTrends(ctx, patientID, metricType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}

// handlePopulationAnalysis handles POST /predictive/population
func (h *PredictiveHandler) handlePopulationAnalysis(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	
	var criteria PopulationCriteria
	if err := json.NewDecoder(r.Body).Decode(&criteria); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	report, err := h.engine.PopulationAnalysis(ctx, criteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

// InMemoryDataStore implements DataStore with in-memory storage
type InMemoryDataStore struct {
	history map[string][]DataPoint
	trends  map[string][]Trend
	mu      sync.RWMutex
}

// NewInMemoryDataStore creates a new in-memory data store
func NewInMemoryDataStore() *InMemoryDataStore {
	return &InMemoryDataStore{
		history: make(map[string][]DataPoint),
		trends:  make(map[string][]Trend),
	}
}

// GetPatientHistory implements DataStore
func (s *InMemoryDataStore) GetPatientHistory(ctx context.Context, patientID string) ([]DataPoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.history[patientID], nil
}

// GetPopulationData implements DataStore
func (s *InMemoryDataStore) GetPopulationData(ctx context.Context, criteria PopulationCriteria) ([]DataPoint, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var results []DataPoint
	
	for _, history := range s.history {
		for _, dp := range history {
			// Apply criteria filters
			if criteria.TenantID != "" && dp.Metadata["tenant_id"] != criteria.TenantID {
				continue
			}
			
			results = append(results, dp)
		}
	}
	
	return results, nil
}

// StoreTrend implements DataStore
func (s *InMemoryDataStore) StoreTrend(ctx context.Context, trend Trend) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	s.trends[trend.PatientID] = append(s.trends[trend.PatientID], trend)
	return nil
}

// GetTrends implements DataStore
func (s *InMemoryDataStore) GetTrends(ctx context.Context, patientID string) ([]Trend, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	return s.trends[patientID], nil
}

// InMemoryPredictiveModel implements PredictiveModel with simple logic
type InMemoryPredictiveModel struct {
	modelType string
}

// NewInMemoryPredictiveModel creates a new in-memory predictive model
func NewInMemoryPredictiveModel(modelType string) *InMemoryPredictiveModel {
	return &InMemoryPredictiveModel{modelType: modelType}
}

// Predict implements PredictiveModel
func (m *InMemoryPredictiveModel) Predict(ctx context.Context, data []DataPoint) (*Prediction, error) {
	// Simple placeholder prediction logic
	return &Prediction{
		Type:       m.modelType,
		Value:      0.5,
		Confidence: 0.75,
		LowerBound: 0.3,
		UpperBound: 0.7,
		Timestamp:  time.Now(),
		Factors: []RiskFactor{
			{Name: "age", Impact: 0.2, Description: "Age factor"},
			{Name: "history", Impact: 0.3, Description: "Medical history"},
		},
	}, nil
}

// GetType implements PredictiveModel
func (m *InMemoryPredictiveModel) GetType() string {
	return m.modelType
}
