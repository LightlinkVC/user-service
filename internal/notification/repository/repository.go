package repository

import "github.com/lightlink/user-service/internal/notification/domain/dto"

type NotificationRepositoryI interface {
	Send(notification dto.RawNotification) error
}
