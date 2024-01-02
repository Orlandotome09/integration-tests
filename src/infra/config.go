package infra

import "github.com/google/uuid"

var (
	//UUID namespace for generate other IDs
	uuidNamespace = uuid.MustParse("55bdda49-b39c-4084-ab69-530266ba7623")
)

func GetUuidNamespace() uuid.UUID {
	return uuidNamespace
}
