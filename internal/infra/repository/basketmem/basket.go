package basketmem

import (
	"context"
	"sync"

	"github.com/SahandNoey/Cart-Service/internal/domain/model"
	"github.com/SahandNoey/Cart-Service/internal/domain/repository/basketrepo"
)

type Repository struct {
	baskets map[uint64]model.Basket
	lock    sync.RWMutex
}

func New() *Repository {
	return &Repository{
		baskets: make(map[uint64]model.Basket),
		lock:    sync.RWMutex{},
	}
}

func (r *Repository) Create(ctx context.Context, basket model.Basket) error {
	r.lock.RLock()
	if _, ok := r.baskets[basket.Id]; ok {
		return basketrepo.ErrorDuplicateBasketID
	}
	r.lock.RUnlock()

	r.lock.Lock()
	r.baskets[basket.Id] = basket
	r.lock.Unlock()

	return nil
}

func (r *Repository) Get(ctx context.Context, cmd basketrepo.GetCommand) []model.Basket {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var baskets []model.Basket

	if cmd.Id != nil {
		basket, ok := r.baskets[*cmd.Id]
		if !ok {
			return nil
		}

		baskets = []model.Basket{basket}

	} else {
		for _, basket := range r.baskets {
			baskets = append(baskets, basket)
		}
	}

	for i := 0; i < len(baskets); i++ {
		if cmd.CreatedAt != nil {
			if baskets[i].CreatedAt != *cmd.CreatedAt {
				baskets = append(baskets[:i], baskets[i+1:]...)
				i--

				continue
			}
		}

		if cmd.UpdatedAt != nil {
			if baskets[i].UpdatedAt != *cmd.UpdatedAt {
				baskets = append(baskets[:i], baskets[i+1:]...)
				i--

				continue
			}
		}

		if cmd.State != nil {
			if baskets[i].State != *cmd.State {
				baskets = append(baskets[:i], baskets[i+1:]...)
				i--

				continue
			}
		}
	}

	return baskets
}

func (r *Repository) Update(ctx context.Context, basket model.Basket) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.baskets, basket.Id)
	r.baskets[basket.Id] = basket

	return nil
}

func (r *Repository) Delete(ctx context.Context, cmd basketrepo.GetCommand) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.baskets, *cmd.Id)

	return nil
}
