package lib

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	id, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	InfoLog.Println("Generated UUID:", id)
	return id.String()
}
