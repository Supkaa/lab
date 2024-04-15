package repositories

import (
	"cmp"
	"context"
	"lab2/internal/entities"
	"slices"
	"sync"
	"time"
)

type NewCache struct {
	nr   NewRepo
	news map[int]entities.New
	mu   sync.RWMutex
}

func NewNewCache(newRepo NewRepo) *NewCache {
	c := &NewCache{nr: newRepo, news: make(map[int]entities.New)}

	go c.cronUpdateCache()

	return c
}

func (c *NewCache) GetAll(ctx context.Context, orderBy string, order string) ([]entities.New, error) {
	if len(c.news) == 0 {
		dbNews, err := c.nr.GetAll(ctx, orderBy, order)

		if err != nil {
			return nil, err
		}

		c.rewrite(dbNews)

		return dbNews, nil
	}

	var news []entities.New

	for _, n := range c.news {
		news = append(news, n)
	}

	switch orderBy {
	case "created_at":
		if order == "ASC" {
			slices.SortFunc(news, func(a, b entities.New) int {
				if a.CreatedAt.Before(b.CreatedAt) {
					return 1
				}

				return 0
			})
		} else {
			slices.SortFunc(news, func(a, b entities.New) int {
				if a.CreatedAt.After(b.CreatedAt) {
					return 1
				}

				return 0
			})
		}
	case "likes":
		if order == "ASC" {
			slices.SortFunc(news, func(a, b entities.New) int {
				return cmp.Compare(a.Likes, b.Likes)
			})
		} else {
			slices.SortFunc(news, func(a, b entities.New) int {
				return cmp.Compare(b.Likes, a.Likes)
			})
		}
	case "dislikes":
		if order == "ASC" {
			slices.SortFunc(news, func(a, b entities.New) int {
				return cmp.Compare(a.Dislikes, b.Dislikes)
			})
		} else {
			slices.SortFunc(news, func(a, b entities.New) int {
				return cmp.Compare(b.Dislikes, a.Dislikes)
			})
		}
	case "views":
		if order == "ASC" {
			slices.SortFunc(news, func(a, b entities.New) int {
				return cmp.Compare(a.Views, b.Views)
			})
		} else {
			slices.SortFunc(news, func(a, b entities.New) int {
				return cmp.Compare(b.Views, a.Views)
			})
		}
	}

	return news, nil
}

func (c *NewCache) GetById(ctx context.Context, id int) (entities.New, error) {
	if n, found := c.news[id]; found {
		return n, nil
	}

	dbNew, err := c.nr.GetById(ctx, id)

	if err != nil {
		return entities.New{}, err
	}

	c.add(dbNew)

	return dbNew, nil
}

func (c *NewCache) rewrite(news []entities.New) {
	newNews := make(map[int]entities.New)

	for _, n := range news {
		newNews[n.Id] = n
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.news = newNews
}

func (c *NewCache) add(n entities.New) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.news[n.Id] = n
}

func (c *NewCache) cronUpdateCache() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dbNews, err := c.nr.GetAll(context.Background(), "created_at", "DESC")

			if err != nil {
				continue
			}

			c.rewrite(dbNews)
		}
	}
}
