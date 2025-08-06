package dto

type PaginationRequest struct {
	Page   int `query:"page"`
	Limit  int `query:"limit"`
	Offset int `json:"offset"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

type AnalysisResponse struct {
	DetectedItems       []string `json:"detected_items"`
	UsableIngredients   []string `json:"usable_ingredients"`
	UnusableIngredients []string `json:"unusable_ingredients"`
	Feedback            string   `json:"feedback"`
	CreatedAt           string   `json:"created_at,omitempty"`
}
