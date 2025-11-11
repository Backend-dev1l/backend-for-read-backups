package interfaces

import (
	"context"
	"test-http/internal/db"
	"test-http/internal/dto"
)

type UserStatisticsService interface {
	Create(ctx context.Context, request dto.CreateStatisticsRequest) (db.UserStatistic, error)
	GetByID(ctx context.Context, request dto.GetStatisticsRequest) (db.UserStatistic, error)
	List(ctx context.Context, request dto.UpdateStatisticsRequest) ([]db.UserStatistic, error)
	Update(ctx context.Context, request dto.UpdateStatisticsRequest) (db.UserStatistic, error)
	Delete(ctx context.Context, request dto.DeleteStatisticsRequest) error
}
