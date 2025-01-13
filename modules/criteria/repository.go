package criteria

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllCriteriaRepository() (result []Criteria, err error)
	CreateCriteriaRepository(criteria *Criteria) (err error)
	GetCriteriaByIdRepository(criteriaID int) (criteria Criteria, err error)
	UpdateCriteriaRepository(criteria Criteria) (err error)
	DeleteCriteriaRepository(criteria Criteria) (err error)
}

type criteriaRepository struct {
	DB *gorm.DB
}

func NewCriteriaRepository(db *gorm.DB) Repository {
	return &criteriaRepository{
		DB: db,
	}
}

func (r *criteriaRepository) GetAllCriteriaRepository() (result []Criteria, err error) {
	err = r.DB.Order("id ASC").Find(&result).Error
	return result, err
}

func (r *criteriaRepository) CreateCriteriaRepository(criteria *Criteria) (err error) {
	err = r.DB.Create(&criteria).Error
	return err
}

func (r *criteriaRepository) GetCriteriaByIdRepository(criteriaID int) (criteria Criteria, err error) {
	err = r.DB.First(&criteria, criteriaID).Error
	return criteria, err
}

func (r *criteriaRepository) UpdateCriteriaRepository(criteria Criteria) (err error) {
	err = r.DB.Save(&criteria).Error
	return err
}

func (r *criteriaRepository) DeleteCriteriaRepository(criteria Criteria) (err error) {
	err = r.DB.Delete(&criteria).Error
	return err
}
