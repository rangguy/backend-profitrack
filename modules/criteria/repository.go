package criteria

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllCriteriaRepository() (result []Criteria, err error)
	CreateCriteriaRepository(category *Criteria) (err error)
	GetCriteriaByIdRepository(categoryID int) (category Criteria, err error)
	UpdateCriteriaRepository(category Criteria) (err error)
	DeleteCriteriaRepository(category Criteria) (err error)
}

type criteriaRepository struct {
	DB *gorm.DB
}

func NewCriteriaRepository(db *gorm.DB) Repository {
	err := db.Migrator().DropTable(&Criteria{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations Drop Criteria Success!")

	err = db.AutoMigrate(&Criteria{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations Up Criteria Success!")

	return &criteriaRepository{
		DB: db,
	}
}

func (r *criteriaRepository) GetAllCriteriaRepository() (result []Criteria, err error) {
	err = r.DB.Find(&result).Error
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
