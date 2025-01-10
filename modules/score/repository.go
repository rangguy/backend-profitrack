package score

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllScoreRepository() (result []Score, err error)
	CreateScoreRepository(category *Score) (err error)
	DeleteAllScoresRepository() (err error)
}

type scoreRepository struct {
	DB *gorm.DB
}

func NewScoreRepository(db *gorm.DB) Repository {
	return &scoreRepository{
		DB: db,
	}
}

func (r *scoreRepository) GetAllScoreRepository() (result []Score, err error) {
	err = r.DB.Find(&result).Error
	return result, err
}

func (r *scoreRepository) CreateScoreRepository(score *Score) (err error) {
	err = r.DB.Create(&score).Error
	return err
}

func (r *scoreRepository) DeleteAllScoresRepository() (err error) {
	err = r.DB.Where("1 = 1").Delete(&Score{}).Error
	return err
}
