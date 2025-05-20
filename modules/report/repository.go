package report

import "gorm.io/gorm"

type Repository interface {
	CountReportsRepository() (total int64, err error)
	GetAllReportsRepository() (result []Report, err error)
	GetReportByIDRepository(ID int) (result Report, err error)
	GetAllReportDetailRepository(ID int) (result []ReportDetail, err error)
	DeleteReportRepository(report *Report) (err error)
	DeleteDetailReportRepository(reportID int) (err error)
}

type reportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) Repository {
	return &reportRepository{db}
}

func (r *reportRepository) CountReportsRepository() (total int64, err error) {
	err = r.DB.Model(&Report{}).Count(&total).Error
	return total, err
}

func (r *reportRepository) GetAllReportsRepository() (result []Report, err error) {
	err = r.DB.Order("id DESC").Find(&result).Error
	return
}

func (r *reportRepository) GetReportByIDRepository(ID int) (result Report, err error) {
	err = r.DB.First(&result, "id = ?", ID).Error
	return result, nil
}

func (r *reportRepository) GetAllReportDetailRepository(ID int) (result []ReportDetail, err error) {
	err = r.DB.Preload("Product").Where("report_id = ?", ID).Find(&result).Error
	return result, err
}

func (r *reportRepository) DeleteReportRepository(report *Report) (err error) {
	err = r.DB.Delete(report).Error
	return err
}

func (r *reportRepository) DeleteDetailReportRepository(reportID int) (err error) {
	err = r.DB.Where("report_id = ?", reportID).Delete(&ReportDetail{}).Error
	return err
}
