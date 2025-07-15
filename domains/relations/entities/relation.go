package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/malikhisyam/user-graph-service/domains/users/entities"
	"gorm.io/gorm"
)

type Follows struct {
	ID uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey;column:id"`

	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index:idx_follows_follower_id;column:follower_id"`
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index:idx_follows_following_id;column:following_id"`

	Follower  entities.User `gorm:"foreignKey:FollowerID;references:ID;constraint:OnDelete:CASCADE;"`
	Following entities.User `gorm:"foreignKey:FollowingID;references:ID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time      `gorm:"type:timestamp;column:created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at"`
}


