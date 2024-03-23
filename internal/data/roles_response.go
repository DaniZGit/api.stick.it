package data

import (
	database "github.com/DaniZGit/api.stick.it/internal/db/generated/models"
	"github.com/gofrs/uuid"
)

type Role struct {
	ID        uuid.NullUUID        `json:"id"`
	Title  string           `json:"title"`
}

type RoleResponse struct {
	Role Role `json:"role"`
}

type RolesResponse struct {
	Metadata Metadata `json:"metadata"`
	Roles []Role `json:"roles"`
}

func BuildRoleResponse(roleRows interface{}, metadata Metadata) any {
	switch value := roleRows.(type) {
		case database.Role:
			return RoleResponse{
				Role: Role{
					ID: uuid.NullUUID{UUID: value.ID, Valid: !value.ID.IsNil()},
					Title: value.Title,
				},
			}
		case []database.GetRolesRow:
			return castToRolesResponse(value, metadata)
	}

	return RoleResponse{}
}

func castToRolesResponse(roleRows []database.GetRolesRow, metadata Metadata) RolesResponse {
	if roleRows == nil || len(roleRows) <= 0 {
		return  RolesResponse{
			Metadata: metadata,
			Roles: []Role{},
		}
	}

	roles := []Role{}
	for _, roleRow := range roleRows {
		role := Role{
			ID: uuid.NullUUID{UUID: roleRow.ID, Valid: !roleRow.ID.IsNil()},
			Title: roleRow.Title,
		}

		roles = append(roles, role)
	}

	return RolesResponse{
		Metadata: metadata,
		Roles: roles,
	}
}