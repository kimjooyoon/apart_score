package personalization

import (
	"apart_score/pkg/scoring"
	"testing"
	"time"
)

func TestDefaultPersonalizationEngine_LearnFromInteraction(t *testing.T) {
	engine := NewDefaultPersonalizationEngine()

	userID := "user123"
	apartmentID := "apt456"

	interaction := UserInteraction{
		Action:      InteractionLike,
		Duration:    time.Minute * 5,
		ScoreViewed: true,
		ContactMade: false,
		Timestamp:   time.Now(),
	}

	err := engine.LearnFromInteraction(userID, apartmentID, interaction)
	if err != nil {
		t.Fatalf("LearnFromInteraction failed: %v", err)
	}

	// 프로필이 생성되었는지 확인
	profile, err := engine.GetUserProfile(userID)
	if err != nil {
		t.Fatalf("GetUserProfile failed: %v", err)
	}

	if profile.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, profile.UserID)
	}

	if profile.InteractionCount != 1 {
		t.Errorf("Expected interaction count 1, got %d", profile.InteractionCount)
	}
}

func TestDefaultPersonalizationEngine_GetPersonalizedWeights(t *testing.T) {
	engine := NewDefaultPersonalizationEngine()

	userID := "user123"
	baseWeights := map[string]float64{
		"Distance to Station": 0.15,
		"Elevator Presence":   0.07,
		"Construction Year":   0.10,
	}

	// 프로필이 없는 경우 기본 가중치 반환
	weights, err := engine.GetPersonalizedWeights(userID, baseWeights)
	if err != nil {
		t.Fatalf("GetPersonalizedWeights failed: %v", err)
	}

	// 기본 가중치와 동일해야 함
	for feature, expected := range baseWeights {
		if actual, exists := weights[feature]; !exists || actual != expected {
			t.Errorf("Expected weight %v for %s, got %v", expected, feature, actual)
		}
	}
}

func TestDefaultPersonalizationEngine_RecommendApartments(t *testing.T) {
	engine := NewDefaultPersonalizationEngine()

	userID := "user123"
	candidates := []scoring.ScoreResult{
		{TotalScore: 85.0},
		{TotalScore: 90.0},
		{TotalScore: 75.0},
	}

	recommendations, err := engine.RecommendApartments(userID, candidates, 2)
	if err != nil {
		t.Fatalf("RecommendApartments failed: %v", err)
	}

	if len(recommendations) != 2 {
		t.Errorf("Expected 2 recommendations, got %d", len(recommendations))
	}

	// 점수 기준 내림차순 정렬 확인
	if recommendations[0].Score < recommendations[1].Score {
		t.Error("Recommendations should be sorted by score descending")
	}

	// 각 추천의 필수 필드 확인
	for i, rec := range recommendations {
		if rec.Score < 0 || rec.Score > 1 {
			t.Errorf("Recommendation %d has invalid score: %v", i, rec.Score)
		}
		if rec.Confidence < 0 || rec.Confidence > 1 {
			t.Errorf("Recommendation %d has invalid confidence: %v", i, rec.Confidence)
		}
	}
}

func TestProfileScore(t *testing.T) {
	engine := NewDefaultPersonalizationEngine()
	userID := "test-user"

	// 프로필 생성 (상호작용 추가)
	interaction := UserInteraction{
		Action:    InteractionLike,
		Duration:  time.Minute * 5,
		Timestamp: time.Now(),
	}

	// 여러 상호작용 추가로 프로필 완성도 높임
	for i := 0; i < 15; i++ {
		err := engine.LearnFromInteraction(userID, "apt123", interaction)
		if err != nil {
			t.Fatalf("LearnFromInteraction failed: %v", err)
		}
	}

	// 프로필 조회
	profile, err := engine.GetUserProfile(userID)
	if err != nil {
		t.Fatalf("GetUserProfile failed: %v", err)
	}

	// 프로필 점수가 0보다 크고 1보다 작거나 같아야 함
	if profile.ProfileScore <= 0 || profile.ProfileScore > 1 {
		t.Errorf("Invalid profile score: %v", profile.ProfileScore)
	}

	// 디버깅을 위해 실제 값 출력
	t.Logf("Profile details: InteractionCount=%d, FeatureWeights=%d, History=%d, LastUpdated=%v",
		profile.InteractionCount,
		len(profile.Preferences.FeatureWeights),
		len(profile.Behavior.InteractionHistory),
		time.Since(profile.LastUpdated))

	// 프로필 점수가 정상 범위인지 확인 (완벽한 프로필이라도 1.0는 어려움)
	if profile.ProfileScore >= 0.5 { // 50% 이상이면 성공
		t.Logf("Profile score acceptable: %v", profile.ProfileScore)
	} else {
		t.Errorf("Profile score too low: %v", profile.ProfileScore)
	}
}

func TestInteractionTypes(t *testing.T) {
	// 상호작용 타입들이 제대로 정의되었는지 확인
	expectedTypes := []InteractionType{
		InteractionView,
		InteractionLike,
		InteractionDislike,
		InteractionSave,
		InteractionShare,
		InteractionContact,
		InteractionCompare,
		InteractionSearch,
	}

	if len(expectedTypes) == 0 {
		t.Error("Interaction types not properly defined")
	}

	// 각 타입이 비어있지 않은지 확인
	for _, interactionType := range expectedTypes {
		if interactionType == "" {
			t.Error("Empty interaction type found")
		}
	}
}
