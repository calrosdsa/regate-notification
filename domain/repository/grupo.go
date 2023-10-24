package repository

import "context"

type GrupoRepository interface {
	GetLastMessagesFromGroup(ctx context.Context, id int) ([]MessageGroupPayload, error)
	GetUsersFromGroup(ctx context.Context, id int) ([]FcmToken, error)
}

type GrupoUseCase interface {
	// GetLastMessagesFromGroup(ctx context.Context, id int) ([]MessageGrupo, error)
	// GetUsersFromGroup(ctx context.Context, id int) ([]ProfileUser, error)
	SendNotificationToUsersGroup(ctx context.Context, message []byte) (err error)
	SendNotificationSalaCreation(ctx context.Context, payload []byte) (err error)
}
type FcmToken struct {
	FcmToken string
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

type MessageGroupPayload struct {
	Id              int              `json:"id"`
	GrupoId         int              `json:"grupo_id"`
	Content         string           `json:"content"`
	CreatedAt       string           `json:"created_at,omitempty"`
	ProfileName     string           `json:"profile_name"`
	ProfileApellido *string          `json:"profile_apellido"`
	ProfilePhoto    *string          `json:"profile_photo"`
	ProfileId       int              `json:"profile_id"`
	ReplyTo         *int             `json:"reply_to"`
	TypeMessage     GrupoMessageType `json:"type_message"`
}
type GrupoMessageType int8

const (
	TypeMessageCommon      = 0
	TypeMessageInstalacion = 1
	TypeMessageSala        = 2
)

//

