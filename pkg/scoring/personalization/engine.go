package personalization

import (
	"apart_score/pkg/scoring"
	"fmt"
	"math"
	"sort"
	"time"
)

type DefaultPersonalizationEngine struct {
	userProfiles   map[string]*UserProfile
	interactions   []UserInteraction
	maxHistorySize int
}

func NewDefaultPersonalizationEngine() *DefaultPersonalizationEngine {
	return &DefaultPersonalizationEngine{
		userProfiles:   make(map[string]*UserProfile),
		interactions:   make([]UserInteraction, 0),
		maxHistorySize: 1000,
	}
}
func (e *DefaultPersonalizationEngine) LearnFromInteraction(userID string, apartmentID string, interaction UserInteraction) error {
	interaction.UserID = userID
	interaction.ApartmentID = apartmentID
	interaction.Timestamp = time.Now()
	e.interactions = append(e.interactions, interaction)
	if len(e.interactions) > e.maxHistorySize {
		e.interactions = e.interactions[len(e.interactions)-e.maxHistorySize:]
	}
	e.updateUserProfile(userID, interaction)
	return nil
}
func (e *DefaultPersonalizationEngine) GetPersonalizedWeights(userID string, baseWeights map[string]float64) (map[string]float64, error) {
	profile, exists := e.userProfiles[userID]
	if !exists {
		return baseWeights, nil
	}
	personalized := make(map[string]float64)
	for feature, baseWeight := range baseWeights {
		personalizedWeight := baseWeight
		if preference, exists := profile.Preferences.FeatureWeights[feature]; exists {
			personalizedWeight = baseWeight*0.7 + preference*0.3
		}
		if adjustment := e.calculateBehaviorAdjustment(userID, feature); adjustment != 0 {
			personalizedWeight += adjustment * 0.1
			personalizedWeight = math.Max(0.0, math.Min(1.0, personalizedWeight))
		}
		personalized[feature] = personalizedWeight
	}
	return personalized, nil
}
func (e *DefaultPersonalizationEngine) RecommendApartments(userID string, candidates []scoring.ScoreResult, limit int) ([]Recommendation, error) {
	if limit <= 0 {
		limit = 10
	}
	profile, exists := e.userProfiles[userID]
	if !exists {
		return e.fallbackRecommendation(candidates, limit)
	}
	var recommendations []Recommendation
	for _, candidate := range candidates {
		score := e.calculateRecommendationScore(userID, candidate, profile)
		recommendation := Recommendation{
			ApartmentID:  "unknown",
			Score:        score,
			Confidence:   e.calculateConfidence(profile),
			Reason:       e.generateReason(candidate, profile),
			SimilarUsers: e.countSimilarUsers(userID, candidate),
			Features:     e.extractMatchingFeatures(candidate, profile),
		}
		recommendations = append(recommendations, recommendation)
	}
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}
	return recommendations, nil
}
func (e *DefaultPersonalizationEngine) GetUserProfile(userID string) (*UserProfile, error) {
	profile, exists := e.userProfiles[userID]
	if !exists {
		return nil, fmt.Errorf("사용자 프로필을 찾을 수 없습니다: %s", userID)
	}
	return profile, nil
}
func (e *DefaultPersonalizationEngine) updateUserProfile(userID string, interaction UserInteraction) {
	profile, exists := e.userProfiles[userID]
	if !exists {
		profile = &UserProfile{
			UserID:      userID,
			Preferences: UserPreferences{},
			Behavior:    UserBehavior{},
			LastUpdated: time.Now(),
		}
		e.userProfiles[userID] = profile
	}
	e.updatePreferences(profile, interaction)
	e.updateBehavior(profile, interaction)
	profile.LastUpdated = time.Now()
	profile.InteractionCount++
	profile.ProfileScore = e.calculateProfileScore(profile)
}
func (e *DefaultPersonalizationEngine) updatePreferences(profile *UserProfile, interaction UserInteraction) {
	if profile.Preferences.FeatureWeights == nil {
		profile.Preferences.FeatureWeights = map[string]float64{
			"Distance to Station": 0.15,
			"Elevator Presence":   0.07,
			"Construction Year":   0.10,
			"Parking":             0.06,
		}
	}
	switch interaction.Action {
	case InteractionLike:
		for feature := range profile.Preferences.FeatureWeights {
			profile.Preferences.FeatureWeights[feature] += 0.05
		}
	case InteractionDislike:
		for feature := range profile.Preferences.FeatureWeights {
			profile.Preferences.FeatureWeights[feature] -= 0.05
		}
	}
	for feature, weight := range profile.Preferences.FeatureWeights {
		profile.Preferences.FeatureWeights[feature] = math.Max(0.0, math.Min(1.0, weight))
	}
}
func (e *DefaultPersonalizationEngine) updateBehavior(profile *UserProfile, interaction UserInteraction) {
	profile.Behavior.InteractionHistory = append(profile.Behavior.InteractionHistory, interaction)
	maxHistory := 100
	if len(profile.Behavior.InteractionHistory) > maxHistory {
		profile.Behavior.InteractionHistory = profile.Behavior.InteractionHistory[len(profile.Behavior.InteractionHistory)-maxHistory:]
	}
	totalDuration := time.Duration(0)
	for _, hist := range profile.Behavior.InteractionHistory {
		totalDuration += hist.Duration
	}
	profile.Behavior.AvgSessionDuration = totalDuration / time.Duration(len(profile.Behavior.InteractionHistory))
}
func (e *DefaultPersonalizationEngine) calculateRecommendationScore(userID string, candidate scoring.ScoreResult, profile *UserProfile) float64 {
	baseScore := float64(candidate.TotalScore) / 100.0
	preferenceMultiplier := 1.0
	if profile.Preferences.PreferredScoreRange.Min > 0 && profile.Preferences.PreferredScoreRange.Max > 0 {
		if candidate.TotalScore >= profile.Preferences.PreferredScoreRange.Min &&
			candidate.TotalScore <= profile.Preferences.PreferredScoreRange.Max {
			preferenceMultiplier = 1.2
		}
	}
	behaviorMultiplier := e.calculateBehaviorMultiplier(userID, candidate)
	return baseScore * preferenceMultiplier * behaviorMultiplier
}
func (e *DefaultPersonalizationEngine) calculateBehaviorMultiplier(userID string, candidate scoring.ScoreResult) float64 {
	multiplier := 1.0
	recentInteractions := e.getRecentInteractions(userID, 10)
	positiveCount := 0
	negativeCount := 0
	for _, interaction := range recentInteractions {
		switch interaction.Action {
		case InteractionLike, InteractionSave, InteractionContact:
			positiveCount++
		case InteractionDislike:
			negativeCount++
		}
	}
	total := positiveCount + negativeCount
	if total > 0 {
		positiveRatio := float64(positiveCount) / float64(total)
		multiplier = 0.8 + positiveRatio*0.4
	}
	return multiplier
}
func (e *DefaultPersonalizationEngine) calculateConfidence(profile *UserProfile) float64 {
	baseConfidence := 0.5
	baseConfidence += profile.ProfileScore * 0.3
	interactionFactor := math.Min(float64(profile.InteractionCount)/100.0, 1.0)
	baseConfidence += interactionFactor * 0.2
	return math.Min(baseConfidence, 1.0)
}
func (e *DefaultPersonalizationEngine) generateReason(candidate scoring.ScoreResult, profile *UserProfile) string {
	score := candidate.TotalScore
	if score >= scoring.ScoreValue(profile.Preferences.PreferredScoreRange.Max) {
		return "귀하의 선호 점수 범위를 만족하는 우수한 옵션입니다"
	} else if score >= scoring.ScoreValue(profile.Preferences.PreferredScoreRange.Min) {
		return "귀하의 선호 점수 범위에 적합한 옵션입니다"
	} else {
		return "개선의 여지가 있지만 고려할 만한 옵션입니다"
	}
}
func (e *DefaultPersonalizationEngine) countSimilarUsers(userID string, candidate scoring.ScoreResult) int {
	return int(math.Min(float64(len(e.userProfiles)-1), 10))
}
func (e *DefaultPersonalizationEngine) extractMatchingFeatures(candidate scoring.ScoreResult, profile *UserProfile) []string {
	features := []string{}
	if candidate.TotalScore >= scoring.ScoreValue(profile.Preferences.PreferredScoreRange.Min) {
		features = append(features, "선호 점수 범위")
	}
	features = append(features, "기본 점수")
	return features
}
func (e *DefaultPersonalizationEngine) fallbackRecommendation(candidates []scoring.ScoreResult, limit int) ([]Recommendation, error) {
	var recommendations []Recommendation
	for i, candidate := range candidates {
		recommendation := Recommendation{
			ApartmentID:  fmt.Sprintf("candidate-%d", i),
			Score:        float64(candidate.TotalScore) / 100.0,
			Confidence:   0.5,
			Reason:       "점수 기반 추천",
			SimilarUsers: 0,
			Features:     []string{"점수 기반"},
		}
		recommendations = append(recommendations, recommendation)
	}
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}
	return recommendations, nil
}
func (e *DefaultPersonalizationEngine) calculateProfileScore(profile *UserProfile) float64 {
	score := 0.0
	if len(profile.Preferences.FeatureWeights) > 0 {
		score += 0.25
	}
	if profile.InteractionCount > 10 {
		score += 0.25
	} else if profile.InteractionCount > 0 {
		score += 0.125
	}
	if len(profile.Behavior.InteractionHistory) > 5 {
		score += 0.25
	}
	if time.Since(profile.LastUpdated) < time.Hour*24*7 {
		score += 0.25
	}
	return score
}
func (e *DefaultPersonalizationEngine) getRecentInteractions(userID string, limit int) []UserInteraction {
	var recent []UserInteraction
	for i := len(e.interactions) - 1; i >= 0 && len(recent) < limit; i-- {
		if e.interactions[i].UserID == userID {
			recent = append(recent, e.interactions[i])
		}
	}
	return recent
}
func (e *DefaultPersonalizationEngine) calculateBehaviorAdjustment(userID, feature string) float64 {
	recentInteractions := e.getRecentInteractions(userID, 5)
	if len(recentInteractions) == 0 {
		return 0.0
	}
	positiveCount := 0
	for _, interaction := range recentInteractions {
		switch interaction.Action {
		case InteractionLike, InteractionSave, InteractionContact:
			positiveCount++
		}
	}
	positiveRatio := float64(positiveCount) / float64(len(recentInteractions))
	if positiveRatio >= 0.8 {
		return 0.05
	} else if positiveRatio <= 0.6 {
		return -0.05
	}
	return 0.0
}
