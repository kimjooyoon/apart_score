package main

import (
	"fmt"
	"apart_score/pkg/scoring"
)

func main() {
	weights := scoring.GetScenarioWeights(scoring.ScenarioBalanced)
	var total float64
	for mt, w := range weights {
		fmt.Printf("%s: %.3f\n", mt.KoreanName(), w)
		total += float64(w)
	}
	fmt.Printf("Total weight: %.3f\n", total)
}
