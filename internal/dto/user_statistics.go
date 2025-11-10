package dto

type CreateStatisticsRequest struct {
	UserID            string  `json:"user_id" validate:"required,uuid"`
	TotalWordsLearned int32   `json:"total_words_learned" validate:"gte=0"`
	Accuracy          float64 `json:"accuracy" validate:"gte=0,lte=100"`
	TotalTime         int32   `json:"total_time" validate:"gte=0"`
}

type UpdateStatisticsRequest struct {
	UserID            string  `json:"user_id" validate:"required,uuid"`
	TotalWordsLearned int32   `json:"total_words_learned" validate:"gte=0"`
	Accuracy          float64 `json:"accuracy" validate:"gte=0,lte=100"`
	TotalTime         int32   `json:"total_time" validate:"gte=0"`
}

type ListStatisticsRequest struct {
	Limit  int32 `json:"limit" validate:"gte=0,lte=100"`
	Offset int32 `json:"offset" validate:"gte=0"`
}

type GetStatisticsRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

type DeleteStatisticsRequest struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}
