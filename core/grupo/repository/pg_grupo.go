package repository

import (
	"context"
	"database/sql"
	"log"
	r "notification/domain/repository"
	// "github.com/lib/pq"
)

type grupoRepository struct {
	Conn *sql.DB
}

func NewRepository(sql *sql.DB) r.GrupoRepository {
	return &grupoRepository{
		Conn: sql,
	}
}

func (p grupoRepository) GetLastMessagesFromGroup(ctx context.Context, id int) (res []r.MessagePayload, err error) {
	query := `select m.id,m.chat_id,m.content,m.created_at,p.nombre,p.apellido,p.profile_photo,m.profile_id,m.reply_to,
	m.type_message,m.grupo_id
	 from grupo_message as m inner join profiles as p on p.profile_id = m.profile_id
	where chat_id = $1
	order by created_at desc limit 3`
	res, err = p.fetchMessagesGrupo(ctx, query, id)
	return
}

func (p grupoRepository) GetUsersFromGroup(ctx context.Context, id int) (res []r.FcmToken, err error) {
	query := `select p.fcm_token,p.profile_id from user_grupo as us 
	left join profiles as p on p.profile_id = us.profile_id where grupo_id = $1`
	log.Println("ID", id)
	res, err = p.fetchFcmTokens(ctx, query, id)
	if err != nil {
		log.Println("DEBUG_SQL", err)
	}
	return
}

func (p *grupoRepository) fetchFcmTokens(ctx context.Context, query string, args ...interface{}) (res []r.FcmToken, err error) {
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

func (p *grupoRepository) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.MessagePayload, err error) {
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
	res = make([]r.MessagePayload, 0)
	for rows.Next() {
		t := r.MessagePayload{}
		err = rows.Scan(
			&t.Message.Id,
			&t.Message.ChatId,
			&t.Message.Content,
			&t.Message.CreatedAt,
			&t.Profile.ProfileName,
			&t.Profile.ProfileApellido,
			&t.Profile.ProfilePhoto,
			&t.Profile.ProfileId,
			&t.Message.ReplyTo,
			&t.Message.TypeMessage,
			&t.Message.ParentId,
		)
		res = append(res, t)
	}
	return res, nil
}
