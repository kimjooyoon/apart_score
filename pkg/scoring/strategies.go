package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
	"fmt"
	"math"
	"sort"
)

type StrategyType string

const (
	StrategyWeightedSum   StrategyType = "weighted_sum"
	StrategyGeometricMean StrategyType = "geometric_mean"
	StrategyMinMax        StrategyType = "min_max"
	StrategyHarmonicMean  StrategyType = "harmonic_mean"
)

type ApartmentData struct {
	ID       string                                      `json:"id"`
	Name     string                                      `json:"name"`
	Scores   map[metadata.MetadataType]shared.ScoreValue `json:"scores"`
	Location string                                      `json:"location"`
}
type RankingResult struct {
	Apartment  ApartmentData                             `json:"apartment"`
	Score      float64                                   `json:"score"`
	Rank       int                                       `json:"rank"`
	Percentile float64                                   `json:"percentile"`
	Method     ScoringMethod                             `json:"method"`
	Weights    [metadata.MetadataTypeCount]shared.Weight `json:"weights"`
}
type RankingsSummary struct {
	TotalApartments int             `json:"total_apartments"`
	Strategy        StrategyType    `json:"strategy"`
	TopRanked       []RankingResult `json:"top_ranked"`
	ScoreRange      struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
		Avg float64 `json:"avg"`
	} `json:"score_range"`
}

// CalculateWithStrategy calculates scores using maps (legacy compatibility).
func CalculateWithStrategy(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight,
	strategy StrategyType) (ScoreResult, error) {
	// Convert maps to arrays for better performance
	scoreArray := shared.ScoreArray{}
	weightArray := shared.WeightArray{}

	for mt, score := range scores {
		scoreArray[int(mt)] = score
	}
	for mt, weight := range weights {
		weightArray[int(mt)] = weight
	}

	return CalculateWithStrategyArray(scoreArray, weightArray, strategy)
}

// CalculateWithStrategyArray calculates scores using arrays (recommended for performance).
func CalculateWithStrategyArray(scores [14]shared.ScoreValue,
	weights [14]shared.Weight,
	strategy StrategyType) (ScoreResult, error) {
	if err := validateStrategyInputsArray(scores, weights); err != nil {
		return ScoreResult{}, err
	}

	result := ScoreResult{Method: MethodWeightedSum}

	switch strategy {
	case StrategyWeightedSum:
		totalWeightedSum := 0.0
		totalWeight := shared.Weight(0)
		for i := range scores {
			rawScore := scores[i]
			weight := weights[i]
			weightedScore := shared.MulDivWeight(rawScore, weight)
			result.RawScores[i] = rawScore
			result.Weights[i] = weight
			result.WeightedScores[i] = weightedScore.ToFloat()
			totalWeightedSum += weightedScore.ToFloat()
			totalWeight += weight
		}
		if totalWeight > 0 {
			result.TotalScore = totalWeightedSum / (float64(totalWeight) / float64(shared.WeightScale))
		}
		result.Method = MethodWeightedSum

	case StrategyGeometricMean:
		minScore := shared.ScoreValueFromFloat(0.1)
		logSum := 0.0
		totalWeight := shared.Weight(0)
		for i := range scores {
			rawScore := scores[i]
			weight := weights[i]
			if rawScore < minScore {
				rawScore = minScore
			}
			logVal := math.Log(rawScore.ToFloat()) * weight.ToFloat()
			result.RawScores[i] = scores[i]
			result.Weights[i] = weight
			result.WeightedScores[i] = rawScore.ToFloat()
			logSum += logVal
			totalWeight += weight
		}
		if totalWeight > 0 {
			result.TotalScore = math.Exp(logSum / totalWeight.ToFloat())
		}
		result.Method = MethodGeometricMean

	case StrategyMinMax:
		minScore := shared.ScoreValueFromFloat(100.0)
		for i := range scores {
			rawScore := scores[i]
			weight := weights[i]
			result.RawScores[i] = rawScore
			result.Weights[i] = weight
			weightedScore := shared.MulDivWeight(rawScore, weight)
			result.WeightedScores[i] = weightedScore.ToFloat()
			if weightedScore < minScore {
				minScore = weightedScore
			}
		}
		result.TotalScore = minScore.ToFloat()
		result.Method = MethodMinMax

	case StrategyHarmonicMean:
		minScore := shared.ScoreValueFromFloat(0.1)
		weightedHarmonicSum := 0.0
		totalWeight := shared.Weight(0)
		for i := range scores {
			rawScore := scores[i]
			weight := weights[i]
			if rawScore < minScore {
				rawScore = minScore
			}
			weightedHarmonicSum += weight.ToFloat() / rawScore.ToFloat()
			result.RawScores[i] = scores[i]
			result.Weights[i] = weight
			result.WeightedScores[i] = rawScore.ToFloat()
			totalWeight += weight
		}
		if weightedHarmonicSum > 0 && totalWeight > 0 {
			result.TotalScore = totalWeight.ToFloat() / weightedHarmonicSum
		}
		result.Method = MethodHarmonicMean

	default:
		return ScoreResult{}, fmt.Errorf("지원하지 않는 전략: %s", strategy)
	}

	return result, nil
}
func validateStrategyInputsArray(scores [14]shared.ScoreValue,
	weights [14]shared.Weight) error {
	for i, score := range scores {
		if score < 0 || score > 100*shared.ScoreScale {
			mt := metadata.MetadataType(i)
			return fmt.Errorf("잘못된 점수 범위 (%s: %.1f)", mt.String(), score.ToFloat())
		}
	}
	totalWeight := shared.Weight(0)
	for i, weight := range weights {
		if weight < 0 || weight > shared.WeightScale {
			mt := metadata.MetadataType(i)
			return fmt.Errorf("잘못된 가중치 범위 (%s: %.3f)", mt.String(), weight.ToFloat())
		}
		totalWeight += weight
	}
	if totalWeight < shared.WeightScale-1 || totalWeight > shared.WeightScale+1 {
		return fmt.Errorf("가중치 합계가 1000이 아닙니다 (현재: %d)", totalWeight)
	}
	return nil
}

func validateStrategyInputs(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) error {
	for mt, score := range scores {
		if score < 0 || score > 100*shared.ScoreScale {
			return fmt.Errorf("잘못된 점수 범위 (%s: %.1f)", mt.String(), score.ToFloat())
		}
	}
	totalWeight := shared.Weight(0)
	for _, mt := range shared.FastAllMetadataTypes() {
		weight := weights[mt]
		if weight < 0 || weight > shared.WeightScale {
			return fmt.Errorf("잘못된 가중치 범위 (%s: %.3f)", mt.String(), weight.ToFloat())
		}
		totalWeight += weight
	}
	if totalWeight < shared.WeightScale-1 || totalWeight > shared.WeightScale+1 {
		return fmt.Errorf("가중치 합계가 1000이 아닙니다 (현재: %d)", totalWeight)
	}
	return nil
}
func GetAvailableStrategies() []StrategyType {
	return []StrategyType{
		StrategyWeightedSum,
		StrategyGeometricMean,
		StrategyMinMax,
		StrategyHarmonicMean,
	}
}

// StrategyGuidelines defines when and how to use each calculation strategy.
var StrategyGuidelines = map[StrategyType]StrategyGuide{
	StrategyWeightedSum: {
		UseCase:     "일반적인 선형 평가",
		BestFor:     "균형 잡힌 의사결정, 대부분의 평가 상황",
		WhenToUse:   "특별한 제약이 없을 때, 직관적인 평가가 필요할 때",
		Limitations: "극단값에 덜 민감함, 모든 요소의 중요도가 유사할 때 부적합",
		Example:     "아파트 구매, 임대 수익 평가, 일반적인 투자 분석",
		Strengths:   []string{"직관적", "계산 간단", "결과 예측 쉬움"},
		Weaknesses:  []string{"균형 평가 부족", "약점 요소 무시 가능"},
	},
	StrategyGeometricMean: {
		UseCase:     "균형 요구 평가",
		BestFor:     "모든 요소가 골고루 만족되어야 하는 상황",
		WhenToUse:   "최소 기준이 엄격한 평가, 가족 주택, 장기 거주",
		Limitations: "계산 복잡성, 낮은 점수의 영향이 과도하게 큼",
		Example:     "가족 주택 선택, 교육 환경 평가, 안전 우선 평가",
		Strengths:   []string{"균형 보장", "약점 요소 강조", "공정성 높음"},
		Weaknesses:  []string{"계산 복잡", "결과 예측 어려움", "극단적 페널티"},
	},
	StrategyMinMax: {
		UseCase:     "최소 요구사항 평가",
		BestFor:     "필수 조건이 엄격한 상황",
		WhenToUse:   "안전 기준, 법적 요구사항, 하드 커트라인 존재",
		Limitations: "다른 요소의 장점 무시, 너무 엄격할 수 있음",
		Example:     "안전 기준 평가, 최소 주거 조건, 자격 요건 확인",
		Strengths:   []string{"안전성 보장", "명확한 기준", "패스/페일 명확"},
		Weaknesses:  []string{"유연성 부족", "균형 무시", "과도한 엄격함"},
	},
	StrategyHarmonicMean: {
		UseCase:     "역수 관계 평가",
		BestFor:     "비용 대비 성능, 효율성 중심 평가",
		WhenToUse:   "가격 대비 품질, 에너지 효율, 투자 수익률",
		Limitations: "이해하기 어려움, 실용적 사용 사례 제한적",
		Example:     "가격 대비 성능, 연비 평가, 투자 효율성",
		Strengths:   []string{"효율성 강조", "비용 고려", "균형적 역수 평가"},
		Weaknesses:  []string{"복잡성", "직관성 부족", "제한적 사용성"},
	},
}

// StrategyGuide provides detailed guidance for using a calculation strategy.
type StrategyGuide struct {
	UseCase     string   // 주요 사용 사례
	BestFor     string   // 가장 적합한 상황
	WhenToUse   string   // 구체적인 사용 시점
	Limitations string   // 제한사항 및 주의점
	Example     string   // 실제 적용 예시
	Strengths   []string // 장점들
	Weaknesses  []string // 단점들
}

// RecommendStrategy recommends the best calculation strategy based on user profile.
func RecommendStrategy(userProfile map[string]interface{}) StrategyType {
	// 기본값: Weighted Sum
	if userProfile == nil {
		return StrategyWeightedSum
	}

	// 가족 구성원 수 확인
	familySize, hasFamily := userProfile["family_size"].(int)
	if hasFamily && familySize > 3 {
		return StrategyGeometricMean // 모든 요소 균형 필요
	}

	// 예산 제한 확인
	budget, hasBudget := userProfile["budget_constraint"].(bool)
	if hasBudget && budget {
		return StrategyMinMax // 최소 요구사항 우선
	}

	// 투자 목적 확인
	investment, hasInvestment := userProfile["investment_focus"].(bool)
	if hasInvestment && investment {
		return StrategyHarmonicMean // 효율성 중심
	}

	// 기본 추천
	return StrategyWeightedSum
}

func GetStrategyDescription(strategy StrategyType) string {
	switch strategy {
	case StrategyWeightedSum:
		return "각 메타데이터의 점수에 가중치를 곱한 후 합계를 계산하는 기본 전략입니다."
	case StrategyGeometricMean:
		return "모든 요소가 균형을 이루어야 하는 경우에 적합한 전략입니다. 하나의 낮은 점수가 전체 점수를 크게 낮춥니다."
	case StrategyMinMax:
		return "모든 요소가 일정 수준 이상이어야 하는 경우에 적합합니다. 가장 낮은 점수가 전체 점수를 결정합니다."
	case StrategyHarmonicMean:
		return "낮은 점수에 매우 민감하게 반응하는 전략입니다. 모든 요소가 고르게 중요할 때 사용합니다."
	default:
		return "알 수 없는 전략입니다."
	}
}

// CalculateWithPipeline performs scoring using a custom calculation pipeline.
func CalculateWithPipeline(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight,
	pipeline CalculationPipeline) (ScoreResult, error) {

	result := ScoreResult{
		Method:   MethodWeightedSum, // 기본값
		Scenario: ScenarioBalanced,  // 기본값
	}

	// 우선순위에 따라 스텝 정렬 (낮은 우선순위가 먼저 실행)
	sortedSteps := make([]CalculationStep, len(pipeline.Steps))
	copy(sortedSteps, pipeline.Steps)
	sort.Slice(sortedSteps, func(i, j int) bool {
		return sortedSteps[i].Priority < sortedSteps[j].Priority
	})

	// 각 스텝 실행
	totalScore := 0.0
	for _, step := range sortedSteps {
		// 현재까지의 결과로 조건 확인
		tempResult := result
		tempResult.TotalScore = totalScore

		if step.Condition == nil || step.Condition(tempResult) {
			stepScore := step.Calculator(scores, weights)
			totalScore += stepScore
		}
	}

	result.TotalScore = totalScore
	return result, nil
}

// CreateFamilyPipeline creates a family-oriented calculation pipeline.
func CreateFamilyPipeline() CalculationPipeline {
	return CalculationPipeline{
		Name:        "가족 중심 평가",
		Description: "학군, 크기/가격 균형, 교통 접근성을 고려한 가족 중심 평가",
		Steps: []CalculationStep{
			{
				Name:        "학군 우선 평가",
				Description: "학군 점수를 40% 가중치로 평가",
				Priority:    1,
				Condition:   nil, // 항상 실행
				Calculator: func(scores map[metadata.MetadataType]shared.ScoreValue, weights map[metadata.MetadataType]shared.Weight) float64 {
					schoolScore := scores[metadata.SchoolDistrict].ToFloat()
					return schoolScore * 0.4
				},
			},
			{
				Name:        "크기/가격 균형",
				Description: "아파트 크기와 가격의 균형을 40%로 평가",
				Priority:    2,
				Condition:   nil,
				Calculator: func(scores map[metadata.MetadataType]shared.ScoreValue, weights map[metadata.MetadataType]shared.Weight) float64 {
					sizeScore := scores[metadata.ApartmentSize].ToFloat() * 0.6
					priceScore := scores[metadata.MaintenanceFee].ToFloat() * 0.4
					return (sizeScore + priceScore) * 0.4
				},
			},
			{
				Name:        "교통 보너스",
				Description: "교통 접근성이 좋으면 보너스 점수",
				Priority:    3,
				Condition: func(result ScoreResult) bool {
					return result.TotalScore > 60 // 기본 점수가 60점 이상일 때만 보너스
				},
				Calculator: func(scores map[metadata.MetadataType]shared.ScoreValue, weights map[metadata.MetadataType]shared.Weight) float64 {
					transportScore := scores[metadata.TransportationAccess].ToFloat()
					if transportScore >= 85 {
						return 5 * 0.2
					} else if transportScore >= 75 {
						return 2 * 0.2
					}
					return 0
				},
			},
		},
	}
}
func CalculateRankings(apartments []ApartmentData, weights map[metadata.MetadataType]shared.Weight, strategy StrategyType) (*RankingsSummary, error) {
	if len(apartments) == 0 {
		return nil, fmt.Errorf("순위를 매길 아파트가 없습니다")
	}
	if err := validateStrategyInputs(apartments[0].Scores, weights); err != nil {
		return nil, fmt.Errorf("입력 검증 실패: %w", err)
	}
	var rankings []RankingResult
	var totalScore float64
	minScore := 100.0
	maxScore := 0.0
	for _, apt := range apartments {
		result, err := CalculateWithStrategy(apt.Scores, weights, strategy)
		if err != nil {
			return nil, fmt.Errorf("아파트 %s 점수 계산 실패: %w", apt.ID, err)
		}
		ranking := RankingResult{
			Apartment: apt,
			Score:     result.TotalScore,
			Method:    result.Method,
			Weights:   result.Weights,
		}
		rankings = append(rankings, ranking)
		totalScore += result.TotalScore
		if result.TotalScore < minScore {
			minScore = result.TotalScore
		}
		if result.TotalScore > maxScore {
			maxScore = result.TotalScore
		}
	}
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].Score > rankings[j].Score
	})
	for i := range rankings {
		rankings[i].Rank = i + 1
		if maxScore > minScore {
			rankings[i].Percentile = float64(rankings[i].Score-minScore) / float64(maxScore-minScore) * 100.0
		} else {
			rankings[i].Percentile = 100.0
		}
	}
	summary := &RankingsSummary{
		TotalApartments: len(apartments),
		Strategy:        strategy,
		TopRanked:       rankings,
	}
	summary.ScoreRange.Min = minScore
	summary.ScoreRange.Max = maxScore
	summary.ScoreRange.Avg = totalScore / float64(len(apartments))
	return summary, nil
}
func FormatRankings(summary *RankingsSummary, limit int) string {
	if summary == nil {
		return "순위 데이터가 없습니다."
	}
	output := fmt.Sprintf("🏆 아파트 순위표 (%s 전략)\n", GetStrategyDescription(summary.Strategy))
	output += "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	output += fmt.Sprintf("총 아파트 수: %d개\n", summary.TotalApartments)
	output += fmt.Sprintf("점수 범위: %.1f - %.1f (평균: %.1f)\n", summary.ScoreRange.Min, summary.ScoreRange.Max, summary.ScoreRange.Avg)
	output += "\n📊 순위 결과:\n"
	displayCount := len(summary.TopRanked)
	if limit > 0 && limit < displayCount {
		displayCount = limit
	}
	for i := 0; i < displayCount; i++ {
		ranking := summary.TopRanked[i]
		rankEmoji := getRankEmoji(i + 1)
		output += fmt.Sprintf("%s %d위: %s (%.1f점, %.1f%%)\n",
			rankEmoji,
			ranking.Rank,
			ranking.Apartment.Name,
			ranking.Score,
			ranking.Percentile)
	}
	if displayCount < len(summary.TopRanked) {
		output += fmt.Sprintf("\n... 외 %d개 아파트", len(summary.TopRanked)-displayCount)
	}
	return output
}
func getRankEmoji(rank int) string {
	switch rank {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		return "🏠"
	}
}
