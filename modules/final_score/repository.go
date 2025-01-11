package final_score

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllFinalScoreByMethodID(methodID int) (result []FinalScore, err error)
	CreateFinalScoreByMethodID(finalScore FinalScore) (err error)
}

type finalScoreRepository struct {
	DB *gorm.DB
}

func NewFinalScoreRepository(db *gorm.DB) Repository {
	return &finalScoreRepository{
		DB: db,
	}
}

func (repo *finalScoreRepository) GetAllFinalScoreByMethodID(methodID int) (result []FinalScore, err error) {
	err = repo.DB.Where("method_id = ?", methodID).Order("product_id ASC").Find(&result).Error
	return result, err
}

func (repo *finalScoreRepository) CreateFinalScoreByMethodID(finalScore FinalScore) (err error) {
	err = repo.DB.Create(&finalScore).Error
	return err
}
