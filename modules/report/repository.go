package report

import "gorm.io/gorm"

type Repository interface {
	GetAllReportRepository(methodID int, period string) (result []Report, err error)
	DeleteAllReportRepository(methodID int, period string) (err error)
}

type newReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) Repository {
	return &newReportRepository{db}
}

func (r *newReportRepository) GetAllReportRepository(methodID int, period string) (result []Report, err error) {
	err = r.DB.Preload("Product").Where("method_id = ? AND period = ?", methodID, period).Find(&result).Error
	return result, err
}

func (r *newReportRepository) DeleteAllReportRepository(methodID int, period string) (err error) {
	err = r.DB.Where("method_id = ? AND period = ?", methodID, period).Delete(&Report{}).Error
	return err
}
