package repository 

type MessageNotification struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	EntityId int    `json:"id"`
}


type KafkaMessageEvent string

const (
	KafkaDeleteMessageEvent KafkaMessageEvent = "delete-message"
	KafkaPublishMessageEvent KafkaMessageEvent = "publish-message"
)