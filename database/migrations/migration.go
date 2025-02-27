package migrations

import (
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/criteria_score"
	"backend-profitrack/modules/final_score"
	"backend-profitrack/modules/method"
	"backend-profitrack/modules/product"
	"backend-profitrack/modules/report"
	"backend-profitrack/modules/score"
	"backend-profitrack/modules/user"
	"fmt"
	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	var err error
	//err = db.Migrator().DropTable(&report.Report{}, &report.ReportDetail{})
	//if err != nil {
	//	panic(err)
	//}
	err = db.AutoMigrate(&user.User{}, &product.Product{}, &criteria.Criteria{}, &method.Method{}, &criteria_score.CriteriaScore{}, &score.Score{}, &final_score.FinalScore{}, &report.Report{}, &report.ReportDetail{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations Success!")
}
