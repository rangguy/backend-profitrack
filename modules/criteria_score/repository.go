package criteria_score

import "gorm.io/gorm"

type Repository interface {
	GetAllCriteriaScoreRepository() (result []CriteriaScore, err error)
	GetCriteriaScoreByIdRepository(criteriaScoreID int) (criteriaScore CriteriaScore, err error)
	CreateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error)
	UpdateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error)
	DeleteCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error)
}

type criteriaScoreRepository struct {
	DB *gorm.DB
}

func NewCriteriaScoreRepository(db *gorm.DB) Repository {
	return &criteriaScoreRepository{
		DB: db,
	}
}

func (r *criteriaScoreRepository) GetAllCriteriaScoreRepository() (result []CriteriaScore, err error) {
	err = r.DB.Order("id ASC").Find(&result).Error
	return result, err
}

func (r *criteriaScoreRepository) GetCriteriaScoreByIdRepository(criteriaScoreID int) (criteriaScore CriteriaScore, err error) {
	err = r.DB.First(&criteriaScore, criteriaScoreID).Error
	return criteriaScore, err
}

func (r *criteriaScoreRepository) CreateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error) {
	err = r.DB.Create(&criteriaScore).Error
	return err
}

func (r *criteriaScoreRepository) UpdateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error) {
	err = r.DB.Save(&criteriaScore).Error
	return err
}

func (r *criteriaScoreRepository) DeleteCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error) {
	err = r.DB.Delete(criteriaScore).Error
	return err
}
