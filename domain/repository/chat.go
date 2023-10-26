package repository

type MessageNotify struct {
	SenderId int `json:"sender_id"`
	Ids  []int  `json:"ids"`
	Data []byte `json:"data"`
}