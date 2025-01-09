package migrations

import (
	"backend-profitrack/modules/category"
	"backend-profitrack/modules/criteria"
	"backend-profitrack/modules/product"
	"backend-profitrack/modules/score_smart"
	"backend-profitrack/modules/user"
	"fmt"
	"gorm.io/gorm"
)

func Migrations(db *gorm.DB) {
	var err error
	//err = db.Migrator().DropTable(score_smart.ScoreSmart{})
	//if err != nil {
	//	panic(err)
	//}
	err = db.AutoMigrate(&user.User{}, &category.Category{}, &product.Product{}, &criteria.Criteria{}, &score_smart.ScoreSmart{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations Success!")
}
