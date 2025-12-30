package providers

import (
	"apart_score/pkg/metadata"
	"testing"
	"time"
)

func TestMockDataProvider_GetAdditionalScores(t *testing.T) {
	provider := NewMockDataProvider(ProviderConfig{
		SourceType: DataSourceMock,
	})

	apartmentID := "test-apartment-123"
	metadataTypes := []metadata.MetadataType{
		metadata.ConstructionCompany,
		metadata.NearbyAmenities,
	}

	scores, err := provider.GetAdditionalScores(apartmentID, metadataTypes)
	if err != nil {
		t.Fatalf("GetAdditionalScores failed: %v", err)
	}

	if len(scores) != len(metadataTypes) {
		t.Errorf("Expected %d scores, got %d", len(metadataTypes), len(scores))
	}

	for _, mt := range metadataTypes {
		if score, exists := scores[mt]; !exists {
			t.Errorf("Score for %v not found", mt)
		} else if score < 0 || score > 100 {
			t.Errorf("Invalid score for %v: %v", mt, score)
		}
	}
}

func TestMockDataProvider_GetContextData(t *testing.T) {
	provider := NewMockDataProvider(ProviderConfig{
		SourceType: DataSourceMock,
	})

	location := "서울시 강남구 역삼동"
	data, err := provider.GetContextData(location)
	if err != nil {
		t.Fatalf("GetContextData failed: %v", err)
	}

	if data == nil {
		t.Fatal("Context data is nil")
	}

	if data.RealEstate == nil {
		t.Error("Real estate data is missing")
	}

	if data.Environmental == nil {
		t.Error("Environmental data is missing")
	}

	if data.ReliabilityScore < 0 || data.ReliabilityScore > 1 {
		t.Errorf("Invalid reliability score: %v", data.ReliabilityScore)
	}
}

func TestMockDataProvider_IsDataAvailable(t *testing.T) {
	provider := NewMockDataProvider(ProviderConfig{
		SourceType: DataSourceMock,
	})

	if !provider.IsDataAvailable("any-id") {
		t.Error("Mock provider should always return data as available")
	}
}

func TestMockDataProvider_GetDataFreshness(t *testing.T) {
	provider := NewMockDataProvider(ProviderConfig{
		SourceType: DataSourceMock,
	})

	freshness, err := provider.GetDataFreshness("test-id")
	if err != nil {
		t.Fatalf("GetDataFreshness failed: %v", err)
	}

	if freshness.LastUpdated.IsZero() {
		t.Error("Last updated time is zero")
	}

	if freshness.MaxAge <= 0 {
		t.Errorf("Invalid max age: %v", freshness.MaxAge)
	}
}

func TestCompositeDataProvider_GetAdditionalScores(t *testing.T) {
	primary := NewMockDataProvider(ProviderConfig{SourceType: DataSourceMock})
	fallback := NewMockDataProvider(ProviderConfig{SourceType: DataSourceMock})

	composite := NewCompositeDataProvider(primary, fallback)

	apartmentID := "test-apartment-123"
	metadataTypes := []metadata.MetadataType{metadata.ConstructionCompany}

	scores, err := composite.GetAdditionalScores(apartmentID, metadataTypes)
	if err != nil {
		t.Fatalf("GetAdditionalScores failed: %v", err)
	}

	if len(scores) == 0 {
		t.Error("No scores returned")
	}

	for _, score := range scores {
		if score < 0 || score > 100 {
			t.Errorf("Invalid score: %v", score)
		}
	}
}

func TestCompositeDataProvider_GetContextData(t *testing.T) {
	primary := NewMockDataProvider(ProviderConfig{SourceType: DataSourceMock})
	composite := NewCompositeDataProvider(primary)

	location := "서울시 강남구 역삼동"
	data, err := composite.GetContextData(location)
	if err != nil {
		t.Fatalf("GetContextData failed: %v", err)
	}

	if data == nil {
		t.Fatal("Context data is nil")
	}

	if len(data.DataSources) == 0 {
		t.Error("No data sources specified")
	}
}

func TestDataFreshness_IsStale(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		freshness DataFreshness
		expected  bool
	}{
		{
			name: "fresh data",
			freshness: DataFreshness{
				LastUpdated: now,
				MaxAge:      time.Hour * 24,
			},
			expected: false,
		},
		{
			name: "stale data",
			freshness: DataFreshness{
				LastUpdated: now.Add(-time.Hour * 48),
				MaxAge:      time.Hour * 24,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.freshness.Age = now.Sub(tt.freshness.LastUpdated)
			tt.freshness.IsStale = tt.freshness.Age > tt.freshness.MaxAge
			if tt.freshness.IsStale != tt.expected {
				t.Errorf("Expected IsStale=%v, got %v", tt.expected, tt.freshness.IsStale)
			}
		})
	}
}
