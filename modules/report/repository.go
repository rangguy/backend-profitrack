package report

import "gorm.io/gorm"

type Repository interface {
	GetAllReportsRepository() (result []Report, err error)
	GetReportByIDRepository(ID int) (result Report, err error)
	GetAllReportDetailRepository(ID int) (result []ReportDetail, err error)
	DeleteReportRepository(report *Report) (err error)
	DeleteDetailReportRepository(reportID int) (err error)
}

type newReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) Repository {
	return &newReportRepository{db}
}

func (r *newReportRepository) GetAllReportsRepository() (result []Report, err error) {
	err = r.DB.Find(&result).Error
	return
}

func (r *newReportRepository) GetReportByIDRepository(ID int) (result Report, err error) {
	err = r.DB.First(&result, "id = ?", ID).Error
	return result, nil
}

func (r *newReportRepository) GetAllReportDetailRepository(ID int) (result []ReportDetail, err error) {
	err = r.DB.Preload("Product").Where("report_id = ?", ID).Find(&result).Error
	return result, err
}

func (r *newReportRepository) DeleteReportRepository(report *Report) (err error) {
	err = r.DB.Delete(report).Error
	return err
}

func (r *newReportRepository) DeleteDetailReportRepository(reportID int) (err error) {
	err = r.DB.Where("report_id = ?", reportID).Delete(&ReportDetail{}).Error
	return err
}
