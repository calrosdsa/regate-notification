package repository

import (
	"context"
	"database/sql"
	"log"
	r "notification/domain/repository"

	"github.com/lib/pq"
)

type salaRepo struct {
	Conn *sql.DB
}

func NewRepository(conn *sql.DB) r.SalaRepository {
	return &salaRepo{
		Conn: conn,
	}
}
func (p salaRepo) GetLastMessagesFromSala(ctx context.Context, id int) (res []r.MessagePayload, err error) {
	query := `select m.id,m.chat_id,m.content,m.created_at,p.nombre,p.apellido,p.profile_photo,m.profile_id,m.reply_to,
	m.type_message,m.sala_id
	from sala_message as m inner join profiles as p on p.profile_id = m.profile_id
	where chat_id = $1
	order by created_at desc limit 3`
	res, err = p.fetchMessagesGrupo(ctx, query, id)
	return
}

func (p *salaRepo) GetFcmTokensUserSalasSala(ctx context.Context,salaId int)(res []r.FcmToken,err error){
	query := `select p.fcm_token,p.profile_id from users_sala as us
	inner join profiles as p on p.profile_id = us.profile_id
	where sala_id = $1`

	// select count(*) from users_sala as us
	// inner join profiles as p on p.profile_id = us.profile_id
	// where sala_id = 101;
	res,err = p.fetchFcmTokens(ctx,query,salaId)
	return
}

func (p *salaRepo) DeleteSala(ctx context.Context,salaId int)(err error){
	query := `delete from salas where estado = $1 and sala_id = $2`
	_,err =  p.Conn.ExecContext(ctx,query,r.SalaUnAvailable,salaId)
	return
}

func (p *salaRepo)GetSalaReservaHora(ctx context.Context,id int)(res r.SalaHora,err error){
	log.Println(id)
	query := `select created_at,horas from salas where sala_id = $1`
	err = p.Conn.QueryRowContext(ctx,query,id).Scan(&res.CreatedAt,pq.Array(&res.Horas))
	if err!= nil {
		log.Println(err,"ERROR-SALA")
	}
	return
}

// func (p *salaRepo) fetchFcmTokens(ctx context.Context, query string, args ...interface{}) (res []r.UserSalaFcmToken, err error) {
// 	rows, err := p.Conn.QueryContext(ctx, query, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		errRow := rows.Close()
// 		if errRow != nil {
// 			log.Println(errRow)
// 		}
// 	}()
// 	res = make([]r.UserSalaFcmToken, 0)
// 	for rows.Next() {
// 		t := r.UserSalaFcmToken{}
// 		err = rows.Scan(
// 			&t.FcmToken,
// 			&t.ProfileId,
// 			&t.Amount,
// 		)
// 		res = append(res, t)
// 	}
// 	return res, nil
// }

func (p *salaRepo) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.MessagePayload, err error) {
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
func (p *salaRepo) fetchFcmTokens(ctx context.Context, query string, args ...interface{}) (res []r.FcmToken, err error) {
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