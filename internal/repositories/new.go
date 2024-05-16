package repositories

import (
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"lab2/internal/entities"
)

type NewRepo struct {
	db *sql.DB
}

func NewNewRepo(db *sql.DB) NewRepo {
	return NewRepo{db: db}
}

func (r NewRepo) GetAll(ctx context.Context, orderBy string, order string) ([]entities.New, error) {
	var news []entities.New

	sqlSelect, _, _ := sq.Select("n.id", "n.title", "n.image", "n.summary", "n.created_at",
		"COUNT(CASE WHEN uhnwr.is_like = TRUE THEN 1 END) AS likes",
		"COUNT(CASE WHEN uhnwr.is_like = FALSE THEN 1 END) AS dislikes",
		"COUNT(DISTINCT uhnwr.user_id) AS views").
		From("news n").
		LeftJoin("users_have_news_with_reactions uhnwr ON n.id = uhnwr.new_id").
		OrderBy(fmt.Sprintf("%s %s", orderBy, order)).
		GroupBy("n.id").
		ToSql()

	rows, err := r.db.QueryContext(ctx, sqlSelect)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var n entities.New

		if err := rows.Scan(&n.Id, &n.Title, &n.Image, &n.Summary, &n.CreatedAt, &n.Likes, &n.Dislikes, &n.Views); err != nil {
			return nil, err
		}

		news = append(news, n)
	}

	return news, nil
}

func (r NewRepo) GetById(ctx context.Context, id uuid.UUID) (entities.New, error) {
	var n entities.New
	row := r.db.QueryRowContext(ctx, `SELECT
    	n.id, n.title, n.image, n.summary, n.created_at,
    	COUNT(CASE WHEN uhnwr.is_like = TRUE THEN 1 END) AS likes,
    	COUNT(CASE WHEN uhnwr.is_like = FALSE THEN 1 END) AS dislikes,
    	COUNT(DISTINCT uhnwr.user_id) AS views
	FROM
    	news n
    LEFT JOIN
    	users_have_news_with_reactions uhnwr
    ON
        n.id = uhnwr.new_id
	WHERE n.id = $1
	GROUP BY n.id;`, id)

	if err := row.Scan(&n.Id, &n.Title, &n.Image, &n.Summary, &n.CreatedAt, &n.Likes, &n.Dislikes, &n.Views); err != nil {
		return entities.New{}, err
	}

	return n, nil
}
