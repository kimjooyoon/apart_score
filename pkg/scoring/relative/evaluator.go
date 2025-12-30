package relative

import (
	"apart_score/pkg/scoring"
	"math"
	"sort"
)

type DefaultRelativeEvaluator struct{}

func NewDefaultRelativeEvaluator() *DefaultRelativeEvaluator {
	return &DefaultRelativeEvaluator{}
}
func (e *DefaultRelativeEvaluator) EvaluateRelatively(target ApartmentScore, group []ApartmentScore) (*RelativeScore, error) {
	if len(group) == 0 {
		return nil, &EvaluationError{Message: "비교할 아파트가 없습니다"}
	}
	distribution := e.calculateDistribution(group)
	percentileRank := e.CalculatePercentile(target.Score, group)
	regionalRank, groupRank := e.calculateRanks(target, group)
	comparison := e.calculateComparison(target, group)
	return &RelativeScore{
		ApartmentID:       target.ID,
		AbsoluteScore:     target.Score,
		PercentileRank:    percentileRank,
		RegionalRank:      regionalRank,
		GroupRank:         groupRank,
		ScoreDistribution: distribution,
		Comparison:        comparison,
	}, nil
}
func (e *DefaultRelativeEvaluator) CalculatePercentile(score scoring.ScoreValue, group []ApartmentScore) float64 {
	if len(group) == 0 {
		return 0.0
	}
	scores := make([]scoring.ScoreValue, len(group))
	for i, apt := range group {
		scores[i] = apt.Score
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i] < scores[j]
	})
	lower := 0
	for lower < len(scores) && scores[lower] < score {
		lower++
	}
	if lower == 0 {
		return 0.0
	}
	if lower == len(scores) {
		return 100.0
	}
	lowerScore := scores[lower-1]
	upperScore := scores[lower]
	if upperScore == lowerScore {
		return float64(lower-1) / float64(len(scores)-1) * 100.0
	}
	position := float64(lower-1) + (float64(score)-float64(lowerScore))/(float64(upperScore)-float64(lowerScore))
	return position / float64(len(scores)-1) * 100.0
}
func (e *DefaultRelativeEvaluator) FindSimilarApartments(target ApartmentScore, candidates []ApartmentScore, criteria SimilarityCriteria) []ApartmentScore {
	if criteria.MaxResults <= 0 {
		criteria.MaxResults = 10
	}
	if criteria.ScoreRange <= 0 {
		criteria.ScoreRange = 10.0
	}
	var similar []ApartmentScore
	scoreRange := criteria.ScoreRange
	for _, candidate := range candidates {
		if candidate.ID == target.ID {
			continue
		}
		scoreDiff := math.Abs(float64(candidate.Score - target.Score))
		if scoreDiff > scoreRange {
			continue
		}
		locationSimilarity := 1.0
		if criteria.LocationWeight > 0 {
			if target.Location != candidate.Location {
				locationSimilarity = 0.5
			}
		}
		similarityScore := (1.0-scoreDiff/scoreRange)*(1.0-criteria.LocationWeight) +
			locationSimilarity*criteria.LocationWeight
		if similarityScore >= 0.6 {
			similar = append(similar, candidate)
		}
		if len(similar) >= criteria.MaxResults {
			break
		}
	}
	return similar
}
func (e *DefaultRelativeEvaluator) EvaluateApartments(apartments []ApartmentScore, groupCriteria GroupCriteria) ([]RelativeScore, error) {
	if len(apartments) == 0 {
		return nil, &EvaluationError{Message: "평가할 아파트가 없습니다"}
	}
	var results []RelativeScore
	for _, apt := range apartments {
		group := e.filterGroup(apt, apartments, groupCriteria)
		result, err := e.EvaluateRelatively(apt, group)
		if err != nil {
			continue
		}
		results = append(results, *result)
	}
	return results, nil
}
func (e *DefaultRelativeEvaluator) calculateDistribution(apartments []ApartmentScore) ScoreDistribution {
	if len(apartments) == 0 {
		return ScoreDistribution{}
	}
	scores := make([]scoring.ScoreValue, len(apartments))
	sum := 0.0
	min := apartments[0].Score
	max := apartments[0].Score
	for i, apt := range apartments {
		scores[i] = apt.Score
		sum += float64(apt.Score)
		if apt.Score < min {
			min = apt.Score
		}
		if apt.Score > max {
			max = apt.Score
		}
	}
	mean := scoring.ScoreValue(sum / float64(len(apartments)))
	sort.Slice(scores, func(i, j int) bool {
		return scores[i] < scores[j]
	})
	var median scoring.ScoreValue
	if len(scores)%2 == 0 {
		median = (scores[len(scores)/2-1] + scores[len(scores)/2]) / 2
	} else {
		median = scores[len(scores)/2]
	}
	variance := 0.0
	for _, score := range scores {
		variance += math.Pow(float64(score)-float64(mean), 2)
	}
	stdDev := scoring.ScoreValue(math.Sqrt(variance / float64(len(scores))))
	q1 := scores[len(scores)/4]
	q3 := scores[len(scores)*3/4]
	return ScoreDistribution{
		Mean:   mean,
		Median: median,
		StdDev: stdDev,
		Min:    min,
		Max:    max,
		Q1:     q1,
		Q3:     q3,
		Count:  len(apartments),
	}
}
func (e *DefaultRelativeEvaluator) calculateRanks(target ApartmentScore, group []ApartmentScore) (regionalRank, groupRank int) {
	type scorePair struct {
		apartment ApartmentScore
		rank      int
	}
	var pairs []scorePair
	for _, apt := range group {
		pairs = append(pairs, scorePair{apartment: apt, rank: 1})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].apartment.Score > pairs[j].apartment.Score
	})
	for i := range pairs {
		if i > 0 && pairs[i].apartment.Score == pairs[i-1].apartment.Score {
			pairs[i].rank = pairs[i-1].rank
		} else {
			pairs[i].rank = i + 1
		}
		if pairs[i].apartment.ID == target.ID {
			groupRank = pairs[i].rank
			regionalRank = pairs[i].rank
		}
	}
	return regionalRank, groupRank
}
func (e *DefaultRelativeEvaluator) calculateComparison(target ApartmentScore, group []ApartmentScore) ScoreComparison {
	better := 0
	worse := 0
	similar := 0
	for _, apt := range group {
		if apt.ID == target.ID {
			continue
		}
		scoreDiff := math.Abs(float64(apt.Score - target.Score))
		if apt.Score > target.Score {
			better++
		} else if apt.Score < target.Score {
			worse++
		} else if scoreDiff <= 5.0 {
			similar++
		}
	}
	total := len(group) - 1
	rankPercentile := 0.0
	if total > 0 {
		rankPercentile = float64(better) / float64(total) * 100.0
	}
	return ScoreComparison{
		BetterThanCount: better,
		WorseThanCount:  worse,
		SimilarCount:    similar,
		RankPercentile:  rankPercentile,
	}
}
func (e *DefaultRelativeEvaluator) filterGroup(target ApartmentScore, candidates []ApartmentScore, criteria GroupCriteria) []ApartmentScore {
	var group []ApartmentScore
	for _, candidate := range candidates {
		if candidate.Score < criteria.ScoreRange.Min || candidate.Score > criteria.ScoreRange.Max {
			continue
		}
		if criteria.MaxGroupSize > 0 && len(group) >= criteria.MaxGroupSize {
			break
		}
		group = append(group, candidate)
	}
	return group
}

type EvaluationError struct {
	Message string
}

func (e *EvaluationError) Error() string {
	return e.Message
}
