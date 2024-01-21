package help

import (
	"github.com/google/uuid"
)

// 获取uuid
func CreateUUID() string {
	id := uuid.New()
	return id.String()
}
