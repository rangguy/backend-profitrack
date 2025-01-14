package final_score

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllFinalScoreByMethodID(methodID int) (result []FinalScore, err error)
}

type finalScoreRepository struct {
	DB *gorm.DB
}

func NewFinalScoreRepository(db *gorm.DB) Repository {
	return &finalScoreRepository{
		DB: db,
	}
}

func (r *finalScoreRepository) GetAllFinalScoreByMethodID(methodID int) (result []FinalScore, err error) {
	err = r.DB.Where("method_id = ?", methodID).Order("final_score DESC").Find(&result).Error
	return result, err
}
