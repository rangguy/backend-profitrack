package criteria_score

import "gorm.io/gorm"

type Repository interface {
	GetAllCriteriaScoreRepository() (result []CriteriaScore, err error)
	GetCriteriaScoreByProductIdRepository(productID int) (criteriaScores []CriteriaScore, err error)
	CreateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error)
	UpdateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error)
	DeleteCriteriaScoreRepository(productID int) (err error)
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
	err = r.DB.Order("product_id ASC, criteria_id ASC").Find(&result).Error
	return result, err
}

func (r *criteriaScoreRepository) GetCriteriaScoreByProductIdRepository(productID int) (criteriaScores []CriteriaScore, err error) {
	err = r.DB.Where("product_id = ?", productID).Find(&criteriaScores).Error
	return criteriaScores, err
}

func (r *criteriaScoreRepository) CreateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error) {
	err = r.DB.Create(criteriaScore).Error
	return err
}

func (r *criteriaScoreRepository) UpdateCriteriaScoreRepository(criteriaScore *CriteriaScore) (err error) {
	err = r.DB.Save(criteriaScore).Error
	return err
}

func (r *criteriaScoreRepository) DeleteCriteriaScoreRepository(productID int) (err error) {
	err = r.DB.Where("product_id = ?", productID).Delete(&CriteriaScore{}).Error
	return err
}
