package criteria

import (
	"gorm.io/gorm"
	"log"
	"time"
)

type Repository interface {
	CountCriteriaRepository() (total int64, err error)
	GetAllCriteriaRepository() (result []Criteria, err error)
	CreateCriteriaRepository(criteria *Criteria) (err error)
	GetCriteriaByIdRepository(criteriaID int) (criteria Criteria, err error)
	UpdateCriteriaRepository(criteria *Criteria) (err error)
	DeleteCriteriaRepository(criteria *Criteria) (err error)
}

type criteriaRepository struct {
	DB *gorm.DB
}

func NewCriteriaRepository(db *gorm.DB) Repository {
	var count int64
	db.Model(&Criteria{}).Count(&count)

	if count == 0 {
		criteriaList := []Criteria{
			{
				Name:      "Return On Investment",
				Type:      "Benefit",
				Weight:    0.3,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Name:      "Net Profit Margin",
				Type:      "Benefit",
				Weight:    0.4,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Name:      "Rasio Efisiensi",
				Type:      "Cost",
				Weight:    0.3,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		db.Create(&criteriaList)
		log.Println("Criteria created.")
	} else {
		log.Println("Criteria already exists.")
	}

	return &criteriaRepository{
		DB: db,
	}
}

func (r *criteriaRepository) CountCriteriaRepository() (total int64, err error) {
	err = r.DB.Model(&Criteria{}).Count(&total).Error
	return total, err
}

func (r *criteriaRepository) GetAllCriteriaRepository() (result []Criteria, err error) {
	err = r.DB.Order("id ASC").Find(&result).Error
	return result, err
}

func (r *criteriaRepository) CreateCriteriaRepository(criteria *Criteria) (err error) {
	err = r.DB.Create(criteria).Error
	return err
}

func (r *criteriaRepository) GetCriteriaByIdRepository(criteriaID int) (criteria Criteria, err error) {
	err = r.DB.First(&criteria, criteriaID).Error
	return criteria, err
}

func (r *criteriaRepository) UpdateCriteriaRepository(criteria *Criteria) (err error) {
	err = r.DB.Save(criteria).Error
	return err
}

func (r *criteriaRepository) DeleteCriteriaRepository(criteria *Criteria) (err error) {
	err = r.DB.Delete(criteria).Error
	return err
}
