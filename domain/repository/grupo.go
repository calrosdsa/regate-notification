package repository

import "context"

type GrupoRepository interface {
	GetLastMessagesFromGroup(ctx context.Context, id int) ([]MessagePayload, error)
	GetUsersFromGroup(ctx context.Context, id int) ([]FcmToken, error)
}

type GrupoUseCase interface {
	// GetLastMessagesFromGroup(ctx context.Context, id int) ([]MessageGrupo, error)
	// GetUsersFromGroup(ctx context.Context, id int) ([]ProfileUser, error)
	SendNotificationToUsersGroup(ctx context.Context, message []byte) (err error)	
	SendNotificationSalaCreation(ctx context.Context, payload []byte) (err error)
}
type FcmToken struct {
	FcmToken  *string
	ProfileId int
}
type ProfileUser struct {
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	// Genero   string     `json:"genero"`
	// BirthDay *time.Time `json:"birthDay"`
	ProfilePhoto *string `json:"profile_photo"`
	//only for user grupo table
	UserGrupoId int     `json:"user_grupo_id,omitempty"`
	FcmToken    *string `json:"fcm_token"`
}
type ProfileBase struct {
	ProfileName     string  `json:"name"`
	ProfileApellido *string `json:"apellido"`
	ProfilePhoto    *string `json:"profile_photo"`
	ProfileId       int     `json:"id"`
}

type Message struct {
	Id          int              `json:"id"`
	LocalId     int64            `json:"local_id"`
	ChatId      int              `json:"chat_id"`
	ProfileId   int              `json:"profile_id"`
	TypeMessage GrupoMessageType `json:"type_message"`
	Content     string           `json:"content"`
	Data        *string          `json:"data"`
	CreatedAt   string           `json:"created_at,omitempty"`
	ParentId    int              `json:"parent_id"`
	IsUser      bool             `json:"is_user"`
	ReplyTo     *int             `json:"reply_to"`
	// ReplyMessage ReplyMessage     `json:"reply_message"`
}

type MessagePayload struct {
	Message Message     `json:"message"`
	Profile ProfileBase `json:"profile"`
}
type GrupoMessageType int8

const (
	TypeMessageCommon      = 0
	TypeMessageInstalacion = 1
	TypeMessageSala        = 2
)

//
