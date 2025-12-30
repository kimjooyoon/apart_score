package apartment

import (
	"apart_score/pkg/metadata"
	"apart_score/pkg/shared"
)

// Apartment은 아파트 정보를 나타내는 엔티티
type Apartment struct {
	ID       string                                `json:"id"`
	Name     string                                `json:"name"`
	Scores   map[metadata.MetadataType]shared.ScoreValue `json:"scores"`
	Location string                                `json:"location"`
}

// NewApartment은 새로운 아파트 인스턴스를 생성합니다.
func NewApartment(id, name, location string) *Apartment {
	return &Apartment{
		ID:       id,
		Name:     name,
		Location: location,
		Scores:   make(map[metadata.MetadataType]shared.ScoreValue),
	}
}

// SetScore은 특정 메타데이터 타입의 점수를 설정합니다.
func (a *Apartment) SetScore(mt metadata.MetadataType, score shared.ScoreValue) {
	if a.Scores == nil {
		a.Scores = make(map[metadata.MetadataType]shared.ScoreValue)
	}
	a.Scores[mt] = score
}

// GetScore은 특정 메타데이터 타입의 점수를 반환합니다.
func (a *Apartment) GetScore(mt metadata.MetadataType) (shared.ScoreValue, bool) {
	score, exists := a.Scores[mt]
	return score, exists
}
