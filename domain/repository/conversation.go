package repository

import "context"

type ConversationUseCase interface {
	SendNotificationMessageConversation(ctx context.Context, message []byte) (err error)
}

type ConversationRepository interface {
	GetLastMessagesConversation(ctx context.Context, id int) ([]MessagePayload, error)	
	GetFcmTokenFromEstablecimientoAdmins(ctx context.Context,id int)([]FcmToken,error)
}




type RolUserAdmin int8

const (
	UserRol  = 0
	AdminRol = 1
)