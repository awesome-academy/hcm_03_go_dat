package repository

import (
	"context"
	"hotel-management/internal/models"

	"gorm.io/gorm"
)

type RoomRepository interface {
	CreateRoom(ctx context.Context, room *models.Room) error
	CreateRoomImage(ctx context.Context, roomImage *models.RoomImage) error
	GetAllRooms(ctx context.Context) ([]models.Room, error)
	FindRoomByID(ctx context.Context, id int) (*models.Room, error)
	FindRoomImageByID(ctx context.Context, id int) (*models.RoomImage, error)
	GetDB() *gorm.DB
	CreateRoomTx(ctx context.Context, tx *gorm.DB, room *models.Room) error
	CreateRoomImageTx(ctx context.Context, tx *gorm.DB, roomImage *models.RoomImage) error
	UpdateRoomTx(ctx context.Context, tx *gorm.DB, room *models.Room) error
	DeleteRoomImageTx(ctx context.Context, tx *gorm.DB, id int) error
	DeleteRoomTx(ctx context.Context, tx *gorm.DB, id int) error
	FindRoomImageByRoomID(ctx context.Context, id int) ([]models.RoomImage, error)
	DeleteRoomImagesByRoomIDTx(ctx context.Context, tx *gorm.DB, roomID int) error
}

type roomRepository struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepository{db: db}
}

func (r *roomRepository) CreateRoom(ctx context.Context, room *models.Room) error {
	err := r.db.WithContext(ctx).Create(room).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) CreateRoomImage(ctx context.Context, roomImage *models.RoomImage) error {
	err := r.db.WithContext(ctx).Create(roomImage).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) CreateRoomTx(ctx context.Context, tx *gorm.DB, room *models.Room) error {
	err := tx.WithContext(ctx).Create(room).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) CreateRoomImageTx(ctx context.Context, tx *gorm.DB, roomImage *models.RoomImage) error {
	err := tx.WithContext(ctx).Create(roomImage).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) DeleteRoomImageTx(ctx context.Context, tx *gorm.DB, id int) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&models.RoomImage{}).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *roomRepository) DeleteRoomImagesByRoomIDTx(ctx context.Context, tx *gorm.DB, roomID int) error {
	return tx.WithContext(ctx).Where("room_id = ?", roomID).Delete(&models.RoomImage{}).Error
}
func (r *roomRepository) DeleteRoomTx(ctx context.Context, tx *gorm.DB, id int) error {
	err := tx.WithContext(ctx).Where("id = ?", id).Delete(&models.Room{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) UpdateRoomTx(ctx context.Context, tx *gorm.DB, room *models.Room) error {
	err := tx.Model(&room).Select(
		"Name", "Type", "PricePerNight", "BedNum", "HasAircon",
		"ViewType", "Description", "IsAvailable",
	).Updates(&room).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *roomRepository) GetAllRooms(ctx context.Context) ([]models.Room, error) {
	var rooms []models.Room
	err := r.db.WithContext(ctx).Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepository) FindRoomByID(ctx context.Context, id int) (*models.Room, error) {
	var room models.Room
	err := r.db.WithContext(ctx).Preload("Images").Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *roomRepository) FindRoomImageByID(ctx context.Context, id int) (*models.RoomImage, error) {
	var roomImage models.RoomImage
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&roomImage).Error
	if err != nil {
		return nil, err
	}
	return &roomImage, nil
}

func (r *roomRepository) FindRoomImageByRoomID(ctx context.Context, id int) ([]models.RoomImage, error) {
	var roomImage []models.RoomImage
	err := r.db.WithContext(ctx).Where("room_id = ?", id).First(&roomImage).Error
	if err != nil {
		return nil, err
	}
	return roomImage, nil
}
