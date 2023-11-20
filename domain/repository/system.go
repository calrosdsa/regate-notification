package repository

import "context"

type SystemRepository interface {
	// SendNotificationUserBilling(ctx context.Context, d []byte)
	GetUserFcmTokens(ctx context.Context,categories []int)([]FcmToken,error)
}

type SystemUseCase interface {
	SendNotificationDiffusion(ctx context.Context, d []byte)
	// AddConsume(ctx context.Context,d Consumo)
}


type Notification struct {
	Id         int        `json:"id"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	EntityId   int        `json:"entity_id"`
	TypeEntity TypeEntity `json:"type_entity"`
	Read       bool       `json:"read"`
	ProfileId  int        `json:"profile_id"`
	CreatedAt  string     `json:"created_at,omitempty"`
}

type RequestNotificationDiffusion struct {
	Notification Notification `json:"notification"`
	Categories   []int        `json:"categories"`
	// Cities       []string     `json:"cities"`
}