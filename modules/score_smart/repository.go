package score_smart

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllScoreSmartRepository() (result []ScoreSmart, err error)
	CreateScoreSmartRepository(category *ScoreSmart) (err error)
	DeleteAllScoreSmartsRepository() (err error)
}

type valueRepository struct {
	DB *gorm.DB
}

func NewScoreSmartRepository(db *gorm.DB) Repository {
	return &valueRepository{
		DB: db,
	}
}

func (r *valueRepository) GetAllScoreSmartRepository() (result []ScoreSmart, err error) {
	err = r.DB.Find(&result).Error
	return result, err
}

func (r *valueRepository) CreateScoreSmartRepository(value *ScoreSmart) (err error) {
	err = r.DB.Create(&value).Error
	return err
}

func (r *valueRepository) DeleteAllScoreSmartsRepository() (err error) {
	err = r.DB.Where("1 = 1").Delete(&ScoreSmart{}).Error
	return err
}
