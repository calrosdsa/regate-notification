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

type TypeEntity int8

const (
	ENTITY_NONE    = 0
	ENTITY_SALA    = 1
	ENTITY_GRUPO   = 2
	ENTITY_ACCOUNT = 3
	ENTITY_BILLING = 4
	ENTITY_RESERVA = 5
	ENTITY_ESTABLECIMIENTO = 6
	ENTITY_URI = 7
	ENTITY_SALA_COMPLETE = 8

)