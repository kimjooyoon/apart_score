package scoring

import "apart_score/pkg/metadata"

type ScoringManager struct {
	scorer  Scorer
	profile ScoringProfile
}

func NewScoringManager(profile ScoringProfile) (*ScoringManager, error) {
	scorer := NewDefaultScorer(profile.Method)
	if profile.Weights == nil {
		profile.Weights = GetScenarioWeights(profile.Scenario)
	}
	return &ScoringManager{
		scorer:  scorer,
		profile: profile,
	}, nil
}
func (sm *ScoringManager) CalculateScore(scores map[metadata.MetadataType]ScoreValue) (*ScoreResult, error) {
	result, err := sm.scorer.Calculate(scores, sm.profile.Weights)
	if err != nil {
		return nil, err
	}
	result.Scenario = sm.profile.Scenario
	return result, nil
}
func (sm *ScoringManager) GetProfile() ScoringProfile {
	return sm.profile
}
func (sm *ScoringManager) UpdateWeights(weights map[metadata.MetadataType]Weight) error {
	if err := sm.scorer.ValidateWeights(weights); err != nil {
		return err
	}
	sm.profile.Weights = weights
	return nil
}
func QuickScore(scores map[metadata.MetadataType]ScoreValue, scenario ScoringScenario) (*ScoreResult, error) {
	profile := ScoringProfile{
		Name:     string(scenario),
		Method:   MethodWeightedSum,
		Scenario: scenario,
	}
	manager, err := NewScoringManager(profile)
	if err != nil {
		return nil, err
	}
	return manager.CalculateScore(scores)
}
func CreateCustomProfile(name, description string, weights map[metadata.MetadataType]Weight) (ScoringProfile, error) {
	if err := NewDefaultScorer(MethodWeightedSum).ValidateWeights(weights); err != nil {
		return ScoringProfile{}, err
	}
	return ScoringProfile{
		Name:        name,
		Description: description,
		Method:      MethodWeightedSum,
		Weights:     weights,
	}, nil
}
func GetPresetProfiles() []ScoringProfile {
	return []ScoringProfile{
		{
			Name:        "균형 잡힌 선택",
			Description: "모든 요소를 균형 있게 고려",
			Method:      MethodWeightedSum,
			Scenario:    ScenarioBalanced,
		},
		{
			Name:        "교통 중심",
			Description: "대중교통 접근성이 중요",
			Method:      MethodWeightedSum,
			Scenario:    ScenarioTransportation,
		},
		{
			Name:        "교육 우선",
			Description: "학군과 교육 환경 중시",
			Method:      MethodWeightedSum,
			Scenario:    ScenarioEducation,
		},
		{
			Name:        "가성비 중시",
			Description: "가격 대비 가치 우선",
			Method:      MethodWeightedSum,
			Scenario:    ScenarioCostEffective,
		},
		{
			Name:        "가족 친화적",
			Description: "가족 단위 거주에 적합",
			Method:      MethodWeightedSum,
			Scenario:    ScenarioFamilyFriendly,
		},
		{
			Name:        "투자 가치",
			Description: "장기 투자 관점에서 고려",
			Method:      MethodWeightedSum,
			Scenario:    ScenarioInvestment,
		},
	}
}
