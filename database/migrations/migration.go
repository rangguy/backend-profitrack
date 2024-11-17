package migrations

import (
	"fmt"
	"gorm.io/gorm"
	"profitrack/modules/category"
	"profitrack/modules/criteria"
	"profitrack/modules/product"
	"profitrack/modules/user"
	"profitrack/modules/value"
)

func Migrations(db *gorm.DB) {
	var err error
	//err = db.Migrator().DropTable(product.Product{})
	//if err != nil {
	//	panic(err)
	//}
	err = db.AutoMigrate(&user.User{}, &category.Category{}, &product.Product{}, &criteria.Criteria{}, &value.Value{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations Success!")
}
