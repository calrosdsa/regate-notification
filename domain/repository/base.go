package repository 

type MessageNotification struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	EntityId int    `json:"id"`
}
