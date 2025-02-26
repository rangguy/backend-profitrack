package score

import (
	"backend-profitrack/modules/final_score"
	"backend-profitrack/modules/report"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllScoreByMethodIDRepository(methodID int) (result []Score, err error)
	GetScoreByProductAndCriteriaAndMethodRepository(productID int, criteriaID int, methodID int) (result Score, err error)
	CreateScoreRepository(score *Score) (err error)
	CreateFinalScoreByMethodIDRepository(methodID int, finalScore *final_score.FinalScore) (err error)
	CreateReportFinalScoreByMethodIDRepository(report *report.Report) (err error)
	DeleteAllScoresByMethodIDRepository(methodID int) (err error)
	CreateReportDetailRepository(reportDetail *report.ReportDetail) (err error)
	UpdateScoreByMethodIDRepository(methodID int, score *Score) (err error)
	UpdateFinalScoreByMethodIDRepository(methodID int, finalScore *final_score.FinalScore) (err error)
	DeleteFinalScoreByMethodIDRepository(methodID int) (err error)
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
		Order("product_id ASC, criteria_id ASC").
		Find(&result).Error
	return result, err
}

func (r *scoreRepository) GetScoreByProductAndCriteriaAndMethodRepository(productID int, criteriaID int, methodID int) (result Score, err error) {
	err = r.DB.Where("product_id = ? AND criteria_id = ? AND method_id = ?", productID, criteriaID, methodID).Find(&result).Error
	return result, err
}

func (r *scoreRepository) CreateScoreRepository(score *Score) (err error) {
	err = r.DB.Create(score).Error
	return err
}

func (r *scoreRepository) CreateFinalScoreByMethodIDRepository(methodID int, finalScore *final_score.FinalScore) (err error) {
	err = r.DB.Where("method_id = ?", methodID).Create(finalScore).Error
	return err
}

func (r *scoreRepository) UpdateFinalScoreByMethodIDRepository(methodID int, finalScore *final_score.FinalScore) (err error) {
	err = r.DB.Where("method_id = ?", methodID).Save(finalScore).Error
	return err
}

func (r *scoreRepository) UpdateScoreByMethodIDRepository(methodID int, score *Score) (err error) {
	err = r.DB.Where("method_id = ?", methodID).Save(score).Error
	return err
}

func (r *scoreRepository) DeleteAllScoresByMethodIDRepository(methodID int) (err error) {
	err = r.DB.Where("method_id = ?", methodID).Delete(&Score{}).Error
	return err
}

func (r *scoreRepository) CreateReportFinalScoreByMethodIDRepository(report *report.Report) (err error) {
	err = r.DB.Create(report).Error
	return err
}

func (r *scoreRepository) CreateReportDetailRepository(reportDetail *report.ReportDetail) (err error) {
	err = r.DB.Create(reportDetail).Error
	return err
}

func (r *scoreRepository) DeleteFinalScoreByMethodIDRepository(methodID int) (err error) {
	err = r.DB.Where("method_id = ?", methodID).Delete(&final_score.FinalScore{}).Error
	return err
}
