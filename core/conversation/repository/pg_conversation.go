package repository

import (
	"context"
	"database/sql"
	"log"
	r "notification/domain/repository"
	// "github.com/lib/pq"
)

type conversationRepo struct {
	Conn *sql.DB
}

func NewRepository(sql *sql.DB) r.ConversationRepository {
	return &conversationRepo{
		Conn: sql,
	}
}

func (p *conversationRepo) GetLastMessagesConversation(ctx context.Context, id int) (res []r.MessagePayload, err error) {
	query := `select m.id,m.chat_id,m.content,m.created_at,
	case when m.is_user then ('Tú') else e.name end,
	null,case when m.is_user then p.profile_photo else e.photo end,(0),
	m.reply_to,m.type_message,m.profile_id
	from conversation_message as m 
	left join profiles as p on p.profile_id = m.profile_id
	left join establecimientos as e on e.establecimiento_id = m.establecimiento_id
	where chat_id = $1
	order by created_at desc limit 3`
	res, err = p.fetchMessagesGrupo(ctx, query, id)
	return
}

func (p *conversationRepo) GetFcmTokenFromEstablecimientoAdmins(ctx context.Context,
id int)(res []r.FcmToken,err error){
	query := `select a.fcm_token,a.id from admin_establecimiento as ae
	inner join admins as a on a.admin_id = ae.admin_id
	where ae.establecimiento_id = $1 and a.estado = $3
	union all 
	select a.fcm_token,a.id from admins as a
	where empresa_id = (select empresa_id from establecimientos where establecimiento_id = $1)
	 and rol =$2 and estado = $3;
	`
	res,err = p.fetchFcmTokens(ctx,query,id,r.AdminRol,r.UserAdminEnabled)
	return
}

func (p *conversationRepo) fetchFcmTokens(ctx context.Context, query string, args ...interface{}) (res []r.FcmToken, err error) {
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

func (p *conversationRepo) fetchMessagesGrupo(ctx context.Context, query string, args ...interface{}) (res []r.MessagePayload, err error) {
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
