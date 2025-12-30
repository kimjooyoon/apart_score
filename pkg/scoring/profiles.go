package scoring

import "apart_score/pkg/metadata"

// ScoringManager는 스코어링 관리를 위한 편의 클래스
type ScoringManager struct {
	scorer  Scorer
	profile ScoringProfile
}

// NewScoringManager는 새로운 스코어링 매니저를 생성합니다.
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

// CalculateScore는 주어진 점수들로 최종 점수를 계산합니다.
func (sm *ScoringManager) CalculateScore(scores map[metadata.MetadataType]ScoreValue) (*ScoreResult, error) {
	result, err := sm.scorer.Calculate(scores, sm.profile.Weights)
	if err != nil {
		return nil, err
	}
	result.Scenario = sm.profile.Scenario
	return result, nil
}

// GetProfile은 현재 프로필을 반환합니다.
func (sm *ScoringManager) GetProfile() ScoringProfile {
	return sm.profile
}

// UpdateWeights는 가중치를 업데이트합니다.
func (sm *ScoringManager) UpdateWeights(weights map[metadata.MetadataType]Weight) error {
	if err := sm.scorer.ValidateWeights(weights); err != nil {
		return err
	}
	sm.profile.Weights = weights
	return nil
}

// QuickScore는 간단한 점수 계산을 위한 헬퍼 함수
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

// CreateCustomProfile은 커스텀 프로필을 생성합니다.
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

// GetPresetProfiles는 미리 정의된 프로필들을 반환합니다.
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
