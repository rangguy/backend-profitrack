package value

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllValueRepository() (result []Value, err error)
	CreateValueRepository(category *Value) (err error)
	GetValueByIdRepository(categoryID int) (category Value, err error)
	DeleteAllValuesRepository() (err error)
}

type valueRepository struct {
	DB *gorm.DB
}

func NewValueRepository(db *gorm.DB) Repository {
	return &valueRepository{
		DB: db,
	}
}

func (r *valueRepository) GetAllValueRepository() (result []Value, err error) {
	err = r.DB.Find(&result).Error
	return result, err
}

func (r *valueRepository) CreateValueRepository(value *Value) (err error) {
	err = r.DB.Create(&value).Error
	return err
}

func (r *valueRepository) GetValueByIdRepository(valueID int) (value Value, err error) {
	err = r.DB.First(&value, valueID).Error
	return value, err
}

func (r *valueRepository) DeleteAllValuesRepository() (err error) {
	err = r.DB.Where("1 = 1").Delete(&Value{}).Error
	return err
}
