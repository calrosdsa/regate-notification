package repository

import (
	"context"
	"database/sql"
	"log"
	r "notification/domain/repository"

	"github.com/lib/pq"
	// "github.com/lib/pq"
)

type systemRepo struct {
	Conn *sql.DB
}

func NewRepository(sql *sql.DB) r.SystemRepository {
	return &systemRepo{
		Conn: sql,
	}
}


func (p systemRepo) GetUserFcmTokens(ctx context.Context,categories []int) (res []r.FcmToken, err error) {
	query := `select p.fcm_token,p.profile_id from profile_category  as pc
	inner join profiles as p on p.profile_id = pc.profile_id
	where pc.category_id = any($1) 
	group by p.profile_id`
	res, err = p.fetchFcmTokens(ctx, query, pq.Array(categories))
	if err != nil {
		log.Println("DEBUG_SQL", err)
	}
	return
}

func (p *systemRepo) fetchFcmTokens(ctx context.Context, query string, args ...interface{}) (res []r.FcmToken, err error) {
	rows, err := p.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			log.Println(errRow)
		}
	}()
	res = make([]r.FcmToken, 0)
	for rows.Next() {
		t := r.FcmToken{}
		err = rows.Scan(
			&t.FcmToken,
			&t.ProfileId,
		)
		res = append(res, t)
	}
	return res, nil
}

