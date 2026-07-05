package model

// OrderPaidEvent представляет событие об успешной оплате заказа.
// Генерируется на уровне сервиса и передаётся в Kafka-продюсер.
type OrderPaidEvent struct {
	EventUUID       string
	OrderUUID       string
	UserUUID        string
	PaymentMethod   string
	TransactionUUID string
}

// ShipAssembledEvent представляет событие о завершении сборки корабля.
// Потребляется из Kafka-топика и используется для обновления статуса заказа.
type ShipAssembledEvent struct {
	EventUUID    string
	OrderUUID    string
	UserUUID     string
	BuildTimeSec int64
}
