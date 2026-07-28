package leaderboards

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	store  *store
	logger *slog.Logger
}

func NewService(db *pgxpool.Pool, logger *slog.Logger) *Service {
	return &Service{store: newStore(db), logger: logger}
}

func (s *Service) GetXPPage(ctx context.Context, page int) (Page, error) {
	total, err := s.store.countCharacters(ctx)
	if err != nil {
		return Page{}, err
	}

	page, totalPages := paginate(total, page)

	entries, err := s.store.getXPPage(ctx, PageSize, (page-1)*PageSize)
	if err != nil {
		return Page{}, err
	}

	return Page{Entries: entries, Page: page, TotalPages: totalPages}, nil
}

func paginate(total, page int) (clampedPage, totalPages int) {
	totalPages = (total + PageSize - 1) / PageSize
	if totalPages < 1 {
		totalPages = 1
	}
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}
	return page, totalPages
}
