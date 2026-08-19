package services

// Универсальный Payload для пакетного сохранения любой сущности
type GenericBulkSavePayload[T any] struct {
	Create []T   `json:"create"`
	Update []T   `json:"update"`
	Delete []int `json:"delete"`
}

// Универсальный ответ
type GenericBulkSaveResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    []T    `json:"data"`
}
