package personalization

import (
	"apart_score/pkg/scoring"
	"time"
)

type UserProfile struct {
	UserID           string          `json:"user_id"`
	Preferences      UserPreferences `json:"preferences"`
	Behavior         UserBehavior    `json:"behavior"`
	ProfileScore     float64         `json:"profile_score"`
	LastUpdated      time.Time       `json:"last_updated"`
	InteractionCount int             `json:"interaction_count"`
}
type UserPreferences struct {
	LocationPreferences map[string]float64 `json:"location_preferences"`
	FeatureWeights      map[string]float64 `json:"feature_weights"`
	BudgetRange         struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"budget_range"`
	PreferredScoreRange struct {
		Min scoring.ScoreValue `json:"min"`
		Max scoring.ScoreValue `json:"max"`
	} `json:"preferred_score_range"`
	MustHaveFeatures   []string `json:"must_have_features"`
	NiceToHaveFeatures []string `json:"nice_to_have_features"`
}
type UserBehavior struct {
	ViewPatterns        []ViewPattern     `json:"view_patterns"`
	SearchPatterns      []SearchPattern   `json:"search_patterns"`
	InteractionHistory  []UserInteraction `json:"interaction_history"`
	FavoriteLocations   []string          `json:"favorite_locations"`
	PreferredPriceRange struct {
		Min float64 `json:"min"`
		Max float64 `json:"max"`
	} `json:"preferred_price_range"`
	AvgSessionDuration time.Duration `json:"avg_session_duration"`
}
type ViewPattern struct {
	Feature     string        `json:"feature"`
	Frequency   int           `json:"frequency"`
	LastViewed  time.Time     `json:"last_viewed"`
	AvgDuration time.Duration `json:"avg_duration"`
}
type SearchPattern struct {
	Query          string    `json:"query"`
	Frequency      int       `json:"frequency"`
	LastSearched   time.Time `json:"last_searched"`
	ResultsClicked int       `json:"results_clicked"`
}
type UserInteraction struct {
	UserID      string             `json:"user_id"`
	ApartmentID string             `json:"apartment_id"`
	Action      InteractionType    `json:"action"`
	Duration    time.Duration      `json:"duration"`
	ScoreViewed bool               `json:"score_viewed"`
	ContactMade bool               `json:"contact_made"`
	Timestamp   time.Time          `json:"timestamp"`
	Context     InteractionContext `json:"context"`
}
type InteractionType string

const (
	InteractionView    InteractionType = "view"
	InteractionLike    InteractionType = "like"
	InteractionDislike InteractionType = "dislike"
	InteractionSave    InteractionType = "save"
	InteractionShare   InteractionType = "share"
	InteractionContact InteractionType = "contact"
	InteractionCompare InteractionType = "compare"
	InteractionSearch  InteractionType = "search"
)

type InteractionContext struct {
	UserLocation   string `json:"user_location"`
	SessionID      string `json:"session_id"`
	SearchQuery    string `json:"search_query"`
	ReferralSource string `json:"referral_source"`
	DeviceType     string `json:"device_type"`
	TimeOfDay      string `json:"time_of_day"`
}
type Recommendation struct {
	ApartmentID  string   `json:"apartment_id"`
	Score        float64  `json:"score"`
	Confidence   float64  `json:"confidence"`
	Reason       string   `json:"reason"`
	SimilarUsers int      `json:"similar_users"`
	Features     []string `json:"features"`
}
type RecommendationRequest struct {
	UserID       string                `json:"user_id"`
	CandidateIDs []string              `json:"candidate_ids"`
	Limit        int                   `json:"limit"`
	Context      InteractionContext    `json:"context"`
	Filters      RecommendationFilters `json:"filters"`
}
type RecommendationFilters struct {
	MinScore         scoring.ScoreValue `json:"min_score"`
	MaxPrice         float64            `json:"max_price"`
	Locations        []string           `json:"locations"`
	RequiredFeatures []string           `json:"required_features"`
	ExcludeViewed    bool               `json:"exclude_viewed"`
}
type SimilarityScore struct {
	UserID             string    `json:"user_id"`
	Similarity         float64   `json:"similarity"`
	CommonInteractions int       `json:"common_interactions"`
	LastInteraction    time.Time `json:"last_interaction"`
}
type CollaborativeFilter interface {
	FindSimilarUsers(userID string, limit int) ([]SimilarityScore, error)
	PredictRating(userID, itemID string) (float64, error)
	UpdateModel(interaction UserInteraction) error
}
