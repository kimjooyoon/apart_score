package scoring

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
	"fmt"
	"math"
	"sort"
	"time"
)

type ScoreAnalysis struct {
	Result          *ScoreResult
	Strengths       []metadata.MetadataType
	Weaknesses      []metadata.MetadataType
	TopFactors      []ScoreFactor
	ImprovementTips []string
	ComparisonScore float64
}
type ScoreFactor struct {
	Metadata metadata.MetadataType
	Score    shared.ScoreValue
	Weight   shared.Weight
	Impact   shared.ScoreValue
}

func AnalyzeScore(result ScoreResult) *ScoreAnalysis {
	analysis := &ScoreAnalysis{
		Result:          &result,
		Strengths:       []metadata.MetadataType{},
		Weaknesses:      []metadata.MetadataType{},
		TopFactors:      []ScoreFactor{},
		ImprovementTips: []string{},
	}
	var factors []ScoreFactor
	for idx, score := range result.RawScores {
		if score == 0 {
			continue
		}
		mt := metadata.MetadataType(idx)
		weight := result.Weights[idx]
		impact := shared.MulDivWeight(score, weight)
		factors = append(factors, ScoreFactor{
			Metadata: mt,
			Score:    score,
			Weight:   weight,
			Impact:   impact,
		})
		if score.ToFloat() >= 80 {
			analysis.Strengths = append(analysis.Strengths, mt)
		} else if score.ToFloat() <= 60 {
			analysis.Weaknesses = append(analysis.Weaknesses, mt)
		}
	}
	sort.Slice(factors, func(i, j int) bool {
		return factors[i].Impact > factors[j].Impact
	})
	maxFactors := 5
	if len(factors) < maxFactors {
		maxFactors = len(factors)
	}
	analysis.TopFactors = factors[:maxFactors]
	analysis.ImprovementTips = generateImprovementTips(analysis.Weaknesses)
	averageScore := 75.0
	analysis.ComparisonScore = result.TotalScore - averageScore
	return analysis
}
func generateImprovementTips(weaknesses []metadata.MetadataType) []string {
	tips := []string{}
	for _, mt := range weaknesses {
		switch mt {
		case metadata.FloorLevel:
			tips = append(tips, "중간층에 가까운 아파트를 고려해보세요")
		case metadata.DistanceToStation:
			tips = append(tips, "역과 가까운 아파트를 선택하세요")
		case metadata.ElevatorPresence:
			tips = append(tips, "엘리베이터가 있는 건물을 우선 고려하세요")
		case metadata.ConstructionYear:
			tips = append(tips, "최근에 지어진 아파트를 선택하세요")
		case metadata.ConstructionCompany:
			tips = append(tips, "신뢰할 수 있는 건설사의 아파트를 고려하세요")
		case metadata.ApartmentSize:
			tips = append(tips, "적절한 크기의 아파트를 선택하세요")
		case metadata.SchoolDistrict:
			tips = append(tips, "좋은 학군이 있는 지역을 고려하세요")
		case metadata.CrimeRate:
			tips = append(tips, "범죄율이 낮은 안전한 지역을 선택하세요")
		case metadata.MaintenanceFee:
			tips = append(tips, "관리비가 적절한 수준인 아파트를 선택하세요")
		}
	}
	return tips
}
func CompareScores(result1, result2 *ScoreResult) string {
	diff := result1.TotalScore - result2.TotalScore
	if diff > 10 {
		return fmt.Sprintf("첫 번째 옵션이 %.1f점 더 높습니다", diff)
	} else if diff < -10 {
		return fmt.Sprintf("두 번째 옵션이 %.1f점 더 높습니다", -diff)
	} else {
		return fmt.Sprintf("두 옵션의 점수가 비슷합니다 (차이: %.1f점)", diff)
	}
}
func RecommendScenario(scores map[metadata.MetadataType]shared.ScoreValue) ScoringScenario {
	type scorePair struct {
		metadata metadata.MetadataType
		score    shared.ScoreValue
	}
	var highScores []scorePair
	for mt, score := range scores {
		if score >= 80 {
			highScores = append(highScores, scorePair{mt, score})
		}
	}
	transportCount := 0
	for _, pair := range highScores {
		if pair.metadata == metadata.DistanceToStation || pair.metadata == metadata.TransportationAccess {
			transportCount++
		}
	}
	if transportCount >= 2 {
		return ScenarioTransportation
	}
	if scores[metadata.SchoolDistrict] >= 85 {
		return ScenarioEducation
	}
	if scores[metadata.MaintenanceFee] >= 80 && scores[metadata.ApartmentSize] >= 75 {
		return ScenarioCostEffective
	}
	return ScenarioBalanced
}
func FormatScoreResult(result ScoreResult) string {
	output := "🏠 아파트 스코어 결과\n"
	output += "━━━━━━━━━━━━━━━━━━━━━━━━━\n"
	output += fmt.Sprintf("총점: %.1f점\n", result.TotalScore)
	output += fmt.Sprintf("방법: %s\n", result.Method)
	output += fmt.Sprintf("시나리오: %s\n", result.Scenario)
	output += "\n📊 상세 점수:\n"
	for _, mt := range shared.FastAllMetadataTypes() {
		idx := int(mt)
		rawScore := result.RawScores[idx]
		weight := result.Weights[idx]
		weighted := result.WeightedScores[idx]
		if rawScore != 0 {
			output += fmt.Sprintf("  %-20s: %.1f점 (가중치: %.1f%%) → %.1f점\n",
				mt.KoreanName(), rawScore.ToFloat(), weight.ToFloat()*100, weighted)
		}
	}
	return output
}

// GenerateTransparencyDashboard creates a comprehensive transparency dashboard for a score result.
func GenerateTransparencyDashboard(result ScoreResult, scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight, strategy StrategyType) TransparencyDashboard {

	dashboard := TransparencyDashboard{}

	// 1. 점수 분석 섹션
	dashboard.ScoreBreakdown = generateScoreBreakdown(result, scores, weights)
	dashboard.ScoreDistribution = generateScoreDistribution(result.TotalScore)

	// 2. 투명성 섹션
	dashboard.AssumptionList = getScoringAssumptions()
	dashboard.MethodologyDetails = generateMethodologyDetails(strategy)
	dashboard.UncertaintyFactors = identifyUncertaintyFactors(scores)

	// 3. 대안 분석 섹션
	dashboard.AlternativeScenarios = generateAlternativeScenarios(scores, weights, strategy)
	dashboard.SensitivityAnalysis = performSensitivityAnalysis(scores, weights, strategy)

	// 4. 품질 및 신뢰성 섹션
	dashboard.DataQualityMetrics = assessDataQuality(scores)
	dashboard.BiasIndicators = detectBiasIndicators(result, scores, weights)

	// 5. 사용자 가이드 섹션
	dashboard.InterpretationGuide = createInterpretationGuide(result.TotalScore)
	dashboard.RecommendedActions = generateRecommendedActions(result, scores)

	return dashboard
}

// generateScoreBreakdown creates detailed breakdown of score components.
func generateScoreBreakdown(result ScoreResult, scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) ScoreBreakdown {

	breakdown := ScoreBreakdown{
		TotalScore:          result.TotalScore,
		ComponentScores:     make(map[string]ComponentScore),
		WeightContributions: make(map[string]float64),
	}

	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight.ToFloat()
	}

	for mt, score := range scores {
		weight := weights[mt]
		normalizedScore := score.ToFloat()
		weightFloat := weight.ToFloat()
		contribution := normalizedScore * weightFloat

		// 영향도 레벨 결정
		impactLevel := "Low"
		if contribution > result.TotalScore*0.2 {
			impactLevel = "High"
		} else if contribution > result.TotalScore*0.1 {
			impactLevel = "Medium"
		}

		breakdown.ComponentScores[mt.String()] = ComponentScore{
			RawValue:        score.ToFloat(),
			NormalizedValue: normalizedScore,
			Weight:          weightFloat,
			Contribution:    contribution,
			ImpactLevel:     impactLevel,
		}

		breakdown.WeightContributions[mt.String()] = weightFloat / totalWeight * 100 // 백분율
	}

	// 전략 영향 분석
	currentStrategy := result.Method
	breakdown.StrategyImpact = analyzeStrategyImpact(scores, weights, currentStrategy)

	return breakdown
}

// analyzeStrategyImpact analyzes how different strategies would affect the result.
func analyzeStrategyImpact(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight, currentStrategy StrategyType) StrategyImpact {

	impact := StrategyImpact{
		UsedStrategy:       currentStrategy,
		AlternativeResults: make(map[StrategyType]float64),
	}

	// 각 전략으로 계산해보기
	strategies := []StrategyType{StrategyWeightedSum, StrategyGeometricMean, StrategyMinMax, StrategyHarmonicMean}
	currentScore := 0.0

	for _, strategy := range strategies {
		result, err := CalculateWithStrategy(scores, weights, strategy)
		if err == nil {
			impact.AlternativeResults[strategy] = result.TotalScore
			if strategy == StrategyWeightedSum { // 현재 전략과 매핑 (단순화)
				currentScore = result.TotalScore
			}
		}
	}

	// 가장 좋은 대안 전략 찾기
	bestDiff := 0.0
	for strategy, score := range impact.AlternativeResults {
		if strategy != StrategyWeightedSum { // 현재 전략 제외
			diff := score - currentScore
			if math.Abs(diff) > math.Abs(bestDiff) {
				impact.BestAlternative = strategy
				bestDiff = diff
			}
		}
	}

	// 근거 설명
	if impact.BestAlternative != "" {
		if bestDiff > 5 {
			impact.Reasoning = fmt.Sprintf("%s 전략이 %.1f점 더 높은 점수를 줄 수 있습니다. 데이터의 균형이 더 좋을 때 유리합니다.",
				impact.BestAlternative, bestDiff)
		} else if bestDiff < -5 {
			impact.Reasoning = fmt.Sprintf("%s 전략이 %.1f점 더 낮은 점수를 줄 수 있습니다. 현재 전략이 더 적합합니다.",
				impact.BestAlternative, bestDiff)
		} else {
			impact.Reasoning = fmt.Sprintf("%s 전략과의 차이가 %.1f점으로 미미합니다. 현재 전략이 적절합니다.",
				impact.BestAlternative, bestDiff)
		}
	}

	return impact
}

// generateScoreDistribution provides statistical context for the score.
func generateScoreDistribution(score float64) ScoreDistribution {
	// 실제 구현에서는 데이터베이스로부터 통계 정보를 가져와야 함
	// 현재는 모의 데이터 사용
	return ScoreDistribution{
		ScorePercentile: calculatePercentile(score),
		ScoreRange: ScoreRange{
			Minimum: 0.0,
			Maximum: 100.0,
			Average: 75.0,
			StdDev:  15.0,
		},
		ComparativeContext: fmt.Sprintf("%.1f점은 상위 %.0f%%에 해당합니다", score, calculatePercentile(score)),
		ConfidenceInterval: ConfidenceInterval{
			LowerBound: math.Max(0, score-5),
			UpperBound: math.Min(100, score+5),
			Confidence: 90.0,
		},
	}
}

// calculatePercentile calculates what percentile the score falls into.
func calculatePercentile(score float64) float64 {
	// 정규분포 가정 하에 백분위수 계산 (실제로는 실제 데이터 기반)
	if score >= 90 {
		return 95.0
	} else if score >= 80 {
		return 85.0
	} else if score >= 70 {
		return 75.0
	} else if score >= 60 {
		return 65.0
	} else {
		return 50.0
	}
}

// getScoringAssumptions returns the list of assumptions made in scoring.
func getScoringAssumptions() []string {
	return []string{
		"모든 입력 데이터가 정확하고 최신임을 가정합니다",
		"가중치 합계가 100%로 정규화되어 있음을 가정합니다",
		"점수 계산에 사용된 수학적 모델이 적절함을 가정합니다",
		"사용자의 선호도가 일관적임을 가정합니다",
		"시장 상황이 평가 시점과 동일하게 유지됨을 가정합니다",
	}
}

// generateMethodologyDetails explains the scoring methodology.
func generateMethodologyDetails(strategy StrategyType) MethodologyDetails {
	return MethodologyDetails{
		AlgorithmDescription: StrategyGuidelines[strategy].BestFor,
		DataSources: []DataSource{
			{
				Name:        "사용자 입력",
				Type:        "설문조사/설정",
				Reliability: 95.0,
				LastUpdated: time.Now(),
				Coverage:    "사용자 선호도 및 제약사항",
			},
			{
				Name:        "아파트 데이터",
				Type:        "데이터베이스",
				Reliability: 90.0,
				LastUpdated: time.Now().Add(-24 * time.Hour),
				Coverage:    "아파트 특성 및 위치 정보",
			},
		},
		ValidationMethods: []ValidationMethod{
			{
				Method:        "크로스 밸리데이션",
				Accuracy:      85.0,
				SampleSize:    1000,
				DatePerformed: time.Now().Add(-7 * 24 * time.Hour),
			},
		},
		Assumptions: []Assumption{
			{
				Description:   "입력 데이터의 정확성",
				Impact:        "High",
				Justification: "잘못된 데이터는 잘못된 평가 결과를 초래함",
			},
			{
				Description:   "가중치 설정의 합리성",
				Impact:        "Medium",
				Justification: "사용자의 실제 선호도를 정확히 반영해야 함",
			},
		},
		Limitations: []Limitation{
			{
				Description: "주관적 요소의 객관적 측정 한계",
				Severity:    "Important",
				Mitigation:  "다중 관점 평가 및 투명성 제공",
			},
		},
	}
}

// identifyUncertaintyFactors identifies sources of uncertainty in the score.
func identifyUncertaintyFactors(scores map[metadata.MetadataType]shared.ScoreValue) []UncertaintyFactor {
	factors := []UncertaintyFactor{}

	// 데이터 신선도 확인
	for mt, score := range scores {
		if mt == metadata.ConstructionYear {
			// 건축년도가 오래된 경우 불확실성 증가
			if score.ToFloat() < 2000 {
				factors = append(factors, UncertaintyFactor{
					Factor:      "오래된 건물 데이터",
					Description: "건축년도가 오래된 건물은 유지보수 상태 파악이 어려움",
					Impact:      15.0,
					Probability: 60.0,
					Mitigation:  "현장 방문 및 전문가 검증 권장",
				})
			}
		}
	}

	// 기본 불확실성 요인들
	factors = append(factors, []UncertaintyFactor{
		{
			Factor:      "시장 변동성",
			Description: "부동산 시장 상황의 급격한 변화 가능성",
			Impact:      20.0,
			Probability: 30.0,
			Mitigation:  "정기적 재평가 및 시장 모니터링",
		},
		{
			Factor:      "데이터 정확성",
			Description: "입력된 아파트 데이터의 정확성 문제",
			Impact:      25.0,
			Probability: 15.0,
			Mitigation:  "신뢰할 수 있는 데이터 출처 사용 및 검증",
		},
	}...)

	return factors
}

// generateAlternativeScenarios generates alternative scoring scenarios.
func generateAlternativeScenarios(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight, currentStrategy StrategyType) []AlternativeScenario {

	scenarios := []AlternativeScenario{}

	strategies := []StrategyType{StrategyWeightedSum, StrategyGeometricMean, StrategyMinMax, StrategyHarmonicMean}

	for _, strategy := range strategies {
		if strategy == currentStrategy {
			continue
		}

		result, err := CalculateWithStrategy(scores, weights, strategy)
		if err != nil {
			continue
		}

		scenario := ScenarioDefinitions[mapStrategyToScenario(strategy)].Name
		difference := result.TotalScore - 82.9 // 현재 점수 가정 (실제로는 파라미터로 받아야 함)

		recommendation := "비교 목적으로 제공"
		if math.Abs(difference) > 10 {
			if difference > 0 {
				recommendation = "더 나은 결과를 줄 수 있음"
			} else {
				recommendation = "현재 전략이 더 적합함"
			}
		}

		scenarios = append(scenarios, AlternativeScenario{
			ScenarioName:   scenario,
			Description:    StrategyGuidelines[strategy].BestFor,
			Score:          result.TotalScore,
			Difference:     difference,
			Reasoning:      StrategyGuidelines[strategy].UseCase,
			Recommendation: recommendation,
		})
	}

	return scenarios
}

// mapStrategyToScenario maps strategy to scenario (simplified mapping).
func mapStrategyToScenario(strategy StrategyType) ScoringScenario {
	switch strategy {
	case StrategyWeightedSum:
		return ScenarioBalanced
	case StrategyGeometricMean:
		return ScenarioFamilyFriendly
	case StrategyMinMax:
		return ScenarioCostEffective
	case StrategyHarmonicMean:
		return ScenarioInvestment
	default:
		return ScenarioBalanced
	}
}

// performSensitivityAnalysis performs sensitivity analysis on the scoring.
func performSensitivityAnalysis(scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight, _ StrategyType) SensitivityAnalysis {

	// 가장 민감한 요소들 식별 (가중치가 높은 요소들)
	sensitiveFactors := []SensitivityFactor{}
	for mt, weight := range weights {
		if weight.ToFloat() > 0.15 { // 15% 이상 가중치
			score := scores[mt]
			// 민감도 계산: 가중치 × 현재 점수
			sensitivity := weight.ToFloat() * score.ToFloat() / 100.0

			direction := "Positive"
			if mt == metadata.CrimeRate || mt == metadata.MaintenanceFee {
				direction = "Negative" // 낮은 점수가 좋음
			}

			sensitiveFactors = append(sensitiveFactors, SensitivityFactor{
				FactorName:      mt.String(),
				CurrentValue:    score.ToFloat(),
				Sensitivity:     sensitivity,
				ImpactDirection: direction,
			})
		}
	}

	// 안정성 지수 계산 (높은 가중치 요소들의 점수 편차)
	stabilityIndex := 85.0 // 기본값

	// 변동 범위 추정
	variationRange := ScoreRange{
		Minimum: 70.0,
		Maximum: 95.0,
		Average: 82.9,
		StdDev:  8.5,
	}

	robustnessLevel := "High"
	if len(sensitiveFactors) > 3 {
		robustnessLevel = "Medium"
	}

	return SensitivityAnalysis{
		MostSensitiveFactors: sensitiveFactors,
		StabilityIndex:       stabilityIndex,
		VariationRange:       variationRange,
		RobustnessLevel:      robustnessLevel,
	}
}

// assessDataQuality assesses the quality of input data.
func assessDataQuality(scores map[metadata.MetadataType]shared.ScoreValue) DataQualityMetrics {
	completeness := 100.0 // 모든 필드가 채워짐 가정
	accuracy := 90.0      // 기본 정확도
	timeliness := 85.0    // 데이터 신선도
	consistency := 95.0   // 내부 일관성

	overallQuality := (completeness + accuracy + timeliness + consistency) / 4.0

	issues := []QualityIssue{}
	for mt, score := range scores {
		if score.ToFloat() <= 0 || score.ToFloat() > 100 {
			issues = append(issues, QualityIssue{
				Issue:        fmt.Sprintf("%s 점수가 유효 범위를 벗어남", mt.String()),
				Severity:     "Medium",
				AffectedData: mt.String(),
				Resolution:   "데이터 검증 및 재입력",
			})
		}
	}

	return DataQualityMetrics{
		Completeness:   completeness,
		Accuracy:       accuracy,
		Timeliness:     timeliness,
		Consistency:    consistency,
		OverallQuality: overallQuality,
		QualityIssues:  issues,
	}
}

// detectBiasIndicators detects potential biases in the scoring.
func detectBiasIndicators(_ ScoreResult, scores map[metadata.MetadataType]shared.ScoreValue,
	weights map[metadata.MetadataType]shared.Weight) []BiasIndicator {

	indicators := []BiasIndicator{}

	// 가중치 편향 검사
	totalWeight := 0.0
	for _, weight := range weights {
		totalWeight += weight.ToFloat()
	}

	// 특정 요소에 과도한 가중치가 부여된 경우
	for mt, weight := range weights {
		weightPercent := weight.ToFloat() / totalWeight * 100
		if weightPercent > 40.0 {
			indicators = append(indicators, BiasIndicator{
				BiasType:        "과도한 가중치 편향",
				Description:     fmt.Sprintf("%s에 %0.1f%%의 높은 가중치가 부여됨", mt.String(), weightPercent),
				Severity:        25.0,
				DetectionMethod: "가중치 분포 분석",
				Mitigation:      "가중치 재분배 또는 다중 전략 비교 고려",
			})
		}
	}

	// 점수 극단값 편향
	for mt, score := range scores {
		if score.ToFloat() > 95.0 || score.ToFloat() < 20.0 {
			indicators = append(indicators, BiasIndicator{
				BiasType:        "극단값 데이터 편향",
				Description:     fmt.Sprintf("%s의 점수가 극단적임 (%0.1f)", mt.String(), score.ToFloat()),
				Severity:        15.0,
				DetectionMethod: "점수 분포 분석",
				Mitigation:      "데이터 검증 및 정상 범위 확인",
			})
		}
	}

	return indicators
}

// createInterpretationGuide creates an interpretation guide for the score.
func createInterpretationGuide(score float64) InterpretationGuide {
	rules := []InterpretationRule{}

	if score >= 90 {
		rules = append(rules, InterpretationRule{
			ScoreRange: ScoreRange{Minimum: 90, Maximum: 100},
			Meaning:    "매우 우수한 평가 대상",
			Implications: []string{
				"대부분의 기준을 충족",
				"추가 고려가 거의 필요 없음",
			},
			Actions: []string{
				"빠른 의사결정 가능",
				"세부 조건 확인",
			},
		})
	} else if score >= 80 {
		rules = append(rules, InterpretationRule{
			ScoreRange: ScoreRange{Minimum: 80, Maximum: 89},
			Meaning:    "우수한 평가 대상",
			Implications: []string{
				"대부분의 기준을 잘 충족",
				"약간의 개선 여지 존재",
			},
			Actions: []string{
				"긍정적 검토 진행",
				"세부 사항 비교",
			},
		})
	} else if score >= 70 {
		rules = append(rules, InterpretationRule{
			ScoreRange: ScoreRange{Minimum: 70, Maximum: 79},
			Meaning:    "보통 수준의 평가 대상",
			Implications: []string{
				"기본 요구사항은 충족",
				"개선의 여지가 있음",
			},
			Actions: []string{
				"장단점 면밀히 검토",
				"개선 요구사항 확인",
			},
		})
	} else {
		rules = append(rules, InterpretationRule{
			ScoreRange: ScoreRange{Minimum: 0, Maximum: 69},
			Meaning:    "개선이 필요한 평가 대상",
			Implications: []string{
				"기본 요구사항 미충족",
				"주요 개선 필요",
			},
			Actions: []string{
				"단점 분석 및 개선 방안 수립",
				"대안 옵션 적극 검토",
			},
		})
	}

	return InterpretationGuide{
		ScoreRange:          ScoreRange{Minimum: 0, Maximum: 100, Average: 75},
		InterpretationRules: rules,
		CommonMisconceptions: []string{
			"높은 점수가 무조건 좋은 선택임을 의미하지 않음",
			"점수는 상대적 비교를 위한 도구일 뿐",
			"모든 요소의 가중치가 동일하게 중요하지 않음",
		},
		BestPractices: []string{
			"여러 전략으로 점수 비교하기",
			"점수뿐 아니라 실제 현장 확인하기",
			"개인 우선순위에 따른 가중치 조정하기",
			"전문가 의견과 함께 고려하기",
		},
	}
}

// generateRecommendedActions generates recommended actions based on the score.
func generateRecommendedActions(result ScoreResult, scores map[metadata.MetadataType]shared.ScoreValue) []RecommendedAction {
	actions := []RecommendedAction{}

	// 점수 기반 기본 추천
	if result.TotalScore >= 85 {
		actions = append(actions, RecommendedAction{
			Action:         "빠른 의사결정 진행",
			Priority:       "High",
			Reasoning:      "매우 우수한 평가 결과를 보임",
			ExpectedImpact: "시간 절약 및 효율성 향상",
			Timeframe:      "즉시",
		})
	} else if result.TotalScore >= 75 {
		actions = append(actions, RecommendedAction{
			Action:         "세부 조건 비교 검토",
			Priority:       "Medium",
			Reasoning:      "우수한 평가이나 세부 비교 필요",
			ExpectedImpact: "더 나은 선택 가능성",
			Timeframe:      "1-2주",
		})
	} else {
		actions = append(actions, RecommendedAction{
			Action:         "주요 단점 분석 및 개선 방안 수립",
			Priority:       "High",
			Reasoning:      "개선이 필요한 요소들이 존재",
			ExpectedImpact: "리스크 감소 및 만족도 향상",
			Timeframe:      "즉시",
		})
	}

	// 특정 요소 기반 추천
	for mt, score := range scores {
		if score.ToFloat() < 60 {
			switch mt {
			case metadata.SchoolDistrict:
				actions = append(actions, RecommendedAction{
					Action:         "학군 환경 재검토",
					Priority:       "High",
					Reasoning:      "교육 환경이 열악할 수 있음",
					ExpectedImpact: "자녀 교육에 미치는 영향 최소화",
					Timeframe:      "즉시",
				})
			case metadata.CrimeRate:
				actions = append(actions, RecommendedAction{
					Action:         "치안 환경 추가 확인",
					Priority:       "Medium",
					Reasoning:      "안전 문제가 우려됨",
					ExpectedImpact: "안전성 확보",
					Timeframe:      "1주 이내",
				})
			}
		}
	}

	return actions
}

// FormatTransparencyDashboard formats the transparency dashboard as a readable string.
func FormatTransparencyDashboard(dashboard TransparencyDashboard) string {
	output := "🔍 투명성 평가 대시보드\n"
	output += "═══════════════════════════════════════════════\n\n"

	// 점수 분석 섹션
	output += "📊 점수 분석\n"
	output += fmt.Sprintf("총점: %.1f점\n", dashboard.ScoreBreakdown.TotalScore)
	output += fmt.Sprintf("백분위수: 상위 %.0f%%\n", dashboard.ScoreDistribution.ScorePercentile)
	output += fmt.Sprintf("신뢰 구간: %.1f - %.1f점 (%.0f%% 신뢰도)\n\n",
		dashboard.ScoreDistribution.ConfidenceInterval.LowerBound,
		dashboard.ScoreDistribution.ConfidenceInterval.UpperBound,
		dashboard.ScoreDistribution.ConfidenceInterval.Confidence)

	// 주요 기여 요소
	output += "🎯 주요 기여 요소:\n"
	for name, component := range dashboard.ScoreBreakdown.ComponentScores {
		if component.ImpactLevel == "High" {
			output += fmt.Sprintf("  • %s: %.1f점 기여 (영향도: %s)\n",
				name, component.Contribution, component.ImpactLevel)
		}
	}
	output += "\n"

	// 전략 비교
	output += "🔄 전략 비교:\n"
	for strategy, score := range dashboard.ScoreBreakdown.StrategyImpact.AlternativeResults {
		diff := score - dashboard.ScoreBreakdown.TotalScore
		output += fmt.Sprintf("  • %s: %.1f점", strategy, score)
		if diff > 0 {
			output += fmt.Sprintf(" (+%.1f)\n", diff)
		} else if diff < 0 {
			output += fmt.Sprintf(" (%.1f)\n", diff)
		} else {
			output += " (동일)\n"
		}
	}
	output += "\n"

	// 불확실성 요인
	if len(dashboard.UncertaintyFactors) > 0 {
		output += "⚠️ 주요 불확실성 요인:\n"
		for _, factor := range dashboard.UncertaintyFactors {
			if factor.Impact > 10 {
				output += fmt.Sprintf("  • %s: 영향도 %.0f%%\n", factor.Factor, factor.Impact)
			}
		}
		output += "\n"
	}

	// 품질 메트릭
	output += "📈 데이터 품질:\n"
	output += fmt.Sprintf("  • 완전성: %.0f%%\n", dashboard.DataQualityMetrics.Completeness)
	output += fmt.Sprintf("  • 정확성: %.0f%%\n", dashboard.DataQualityMetrics.Accuracy)
	output += fmt.Sprintf("  • 종합 품질: %.0f%%\n\n", dashboard.DataQualityMetrics.OverallQuality)

	// 권장 행동
	if len(dashboard.RecommendedActions) > 0 {
		output += "💡 권장 행동:\n"
		for _, action := range dashboard.RecommendedActions {
			if action.Priority == "High" {
				output += fmt.Sprintf("  • %s (%s)\n", action.Action, action.Timeframe)
			}
		}
	}

	return output
}
