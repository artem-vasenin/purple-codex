package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Gorm *gorm.DB
	SQL  *sql.DB
}
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email        string    `gorm:"uniqueIndex;not null"`
	FirstName    string    `gorm:"not null;default:''"`
	LastName     string    `gorm:"not null;default:''"`
	Phone        string    `gorm:"not null;default:''"`
	AvatarData   []byte    `gorm:"type:bytea"`
	AvatarType   string    `gorm:"not null;default:''"`
	PasswordHash string    `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
type RefreshSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;index;not null"`
	TokenHash string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	RevokedAt *time.Time
	CreatedAt time.Time
}

func Open(dsn string) (*DB, error) {
	g, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	sqlDB, err := g.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	return &DB{Gorm: g, SQL: sqlDB}, nil
}

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) CreateUser(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}
func (r *Repository) FindUser(ctx context.Context, email string) (User, error) {
	var user User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}
func (r *Repository) FindUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	var user User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	return user, err
}
func (r *Repository) UpdateUserProfile(ctx context.Context, id uuid.UUID, firstName, lastName, phone string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(map[string]any{"first_name": firstName, "last_name": lastName, "phone": phone}).Error
}
func (r *Repository) UpdateUserAvatar(ctx context.Context, id uuid.UUID, data []byte, contentType string) error {
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Updates(map[string]any{"avatar_data": data, "avatar_type": contentType}).Error
}
func (r *Repository) FindUserAvatar(ctx context.Context, id uuid.UUID) ([]byte, string, error) {
	var user User
	if err := r.db.WithContext(ctx).Select("avatar_data", "avatar_type").First(&user, "id = ?", id).Error; err != nil {
		return nil, "", err
	}
	return user.AvatarData, user.AvatarType, nil
}
func (r *Repository) CreateSession(ctx context.Context, session *RefreshSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}
func (r *Repository) FindSession(ctx context.Context, hash string) (RefreshSession, error) {
	var s RefreshSession
	err := r.db.WithContext(ctx).Where("token_hash = ?", hash).First(&s).Error
	return s, err
}
func (r *Repository) RotateSession(ctx context.Context, oldID uuid.UUID, next *RefreshSession) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UTC()
		result := tx.Model(&RefreshSession{}).Where("id = ? AND revoked_at IS NULL AND expires_at > ?", oldID, now).Update("revoked_at", now)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return gorm.ErrRecordNotFound
		}
		return tx.Create(next).Error
	})
}
func (r *Repository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	now := time.Now().UTC()
	return r.db.WithContext(ctx).Model(&RefreshSession{}).Where("id = ? AND revoked_at IS NULL", id).Update("revoked_at", now).Error
}
