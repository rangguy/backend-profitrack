package score

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllScoreByMethodIDRepository(methodID int) (result []Score, err error)
	CreateScoreByMethodIDRepository(methodID int, score *Score) (err error)
	UpdateScoreByMethodIDRepository(methodID int, score *Score) (err error)
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

func (r *scoreRepository) GetAllScoreByMethodIDRepository(methodID int) (result []Score, err error) {
	err = r.DB.Where("method_id = ?", methodID).
		Order("id ASC").
		Find(&result).Error
	return result, err
}

func (r *scoreRepository) CreateScoreByMethodIDRepository(methodID int, score *Score) (err error) {
	err = r.DB.Where("method_id = ?", methodID).Create(&score).Error
	return err
}

func (r *scoreRepository) UpdateScoreByMethodIDRepository(methodID int, score *Score) (err error) {
	err = r.DB.Where("method_id = ?", methodID).Save(&score).Error
	return err
}

func (r *scoreRepository) DeleteAllScoresRepository() (err error) {
	err = r.DB.Where("1 = 1").Delete(&Score{}).Error
	return err
}
