package repository

import (
	"context"
	"database/sql"
	r "notification/domain/repository"
)

type utilRepo struct {
	Conn *sql.DB
}

func NewRepo(conn *sql.DB) r.UtilRepository {
	return &utilRepo{
		Conn: conn,
	}
}

func (p *utilRepo)GetProfileFcmToken(ctx context.Context,id int)(res string,err error){
	query := `select fcm_token from profiles where profile_id = $1`
	err = p.Conn.QueryRowContext(ctx,query,id).Scan(&res)
	return
}