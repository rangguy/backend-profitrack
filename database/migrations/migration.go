package migrations

import (
	"backend-profitrack/modules/category"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/final_score"
	"backend-profitrack/modules/method"
	"backend-profitrack/modules/product"
	"backend-profitrack/modules/score"
	"backend-profitrack/modules/user"
	"fmt"
	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	var err error
	//err = db.Migrator().DropTable(final_score.FinalScore{}, score.Score{})
	//if err != nil {
	//	panic(err)
	//}
	err = db.AutoMigrate(&user.User{}, &category.Category{}, &product.Product{}, &criteria.Criteria{}, &method.Method{}, &score.Score{}, &final_score.FinalScore{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations Success!")
}
