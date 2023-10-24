package repository

import "context"

type BillingRepository interface {
	// SendNotificationUserBilling(ctx context.Context, d []byte)
	AddConsume(ctx context.Context,d []Consumo)
}

type BillingUseCase interface {
	SendNotificationUserBilling(ctx context.Context, d []byte)
	// AddConsume(ctx context.Context,d Consumo)
}

type Consumo struct {
	Id         int     `json:"id"`
	ProfileId  int     `json:"profile_id"`
	Amount     float64 `json:"amount"`
	Message    string  `json:"message"`
	TypeEntity int8     `json:"type_entity"`
	IdEnitity  int     `json:"id_entity"`
	CreatedAt  string  `json:"created_at"`
}

// type BillingPayload struct {
// 	ProfileId int `json:"id"`
// 	Message   string `json:"message"`
// }
