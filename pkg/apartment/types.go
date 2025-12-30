package apartment

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
)

type Apartment struct {
	ID       string                                      `json:"id"`
	Name     string                                      `json:"name"`
	Scores   map[metadata.MetadataType]shared.ScoreValue `json:"scores"`
	Location string                                      `json:"location"`
}

func NewApartment(id, name, location string) *Apartment {
	return &Apartment{
		ID:       id,
		Name:     name,
		Location: location,
		Scores:   make(map[metadata.MetadataType]shared.ScoreValue),
	}
}
func (a *Apartment) SetScore(mt metadata.MetadataType, score shared.ScoreValue) {
	if a.Scores == nil {
		a.Scores = make(map[metadata.MetadataType]shared.ScoreValue)
	}
	a.Scores[mt] = score
}
func (a *Apartment) GetScore(mt metadata.MetadataType) (shared.ScoreValue, bool) {
	score, exists := a.Scores[mt]
	return score, exists
}
