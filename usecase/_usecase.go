
package usecase

import (
	"context"

	"oasis/domain"
)

// Usecase represents  usecases
type Usecase struct {
	Repository domain.Repository
}

// NewUsecase return  usecase instance
func NewUsecase(r domain.Repository) *Usecase {
	return &Usecase{
		Repository: r,
	}
}

// Store insert a new  to storage
func (u *Usecase) Store(ctx context.Context,  *domain.) (*domain., error) {
	.PopulateCreateTimeStamp() // generate created_at & updated_at field

	return u.Repository.Store(ctx, )
}

// Fetch list  from storage based on criteria
func (u *Usecase) Fetch(ctx context.Context, ctr *domain.Criteria) ([]*domain., error) {
	return u.Repository.Fetch(ctx, ctr)
}

// Count return count of  from storage based on criteria
func (u *Usecase) Count(ctx context.Context, ctr *domain.Criteria) (int64, error) {
	return u.Repository.Count(ctx, ctr)
}

// FetchOne fetch a  by primary id
func (u *Usecase) FetchOne(ctx context.Context, ctr *domain.Criteria) (*domain., error) {
	return u.Repository.FetchOne(ctx, ctr)
}

// Update update a  record
func (u *Usecase) Update(ctx context.Context,  *domain.) (*domain., error) {
	.PopulateUpdateTimeStamp() // update the updated_at timestamp
	return u.Repository.Update(ctx, )
}

// Delete soft delete a  record
func (u *Usecase) Delete(ctx context.Context, ctr *domain.Criteria) error {
	return u.Repository.Delete(ctx, ctr)
}
