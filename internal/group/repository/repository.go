package repository

import "github.com/lightlink/user-service/internal/group/domain/entity"

type GroupRepositoryI interface {
	Create(personalGroupEntity *entity.PersonalGroup) error
}
