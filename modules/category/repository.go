package category

import (
	"fmt"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllCategoryRepository() (result []Category, err error)
	CreateCategoryRepository(category *Category) (err error)
	GetCategoryByIdRepository(categoryID int) (category Category, err error)
	UpdateCategoryRepository(category Category) (err error)
	DeleteCategoryRepository(category Category) (err error)
}

type categoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) Repository {
	err := db.AutoMigrate(&Category{})
	//err := db.Migrator().DropTable(&Category{})
	if err != nil {
		panic(err)
	}
	fmt.Println("Migrations Category Success!")

	return &categoryRepository{
		DB: db,
	}
}

func (r *categoryRepository) GetAllCategoryRepository() (result []Category, err error) {
	err = r.DB.Find(&result).Error
	return result, err
}

func (r *categoryRepository) CreateCategoryRepository(category *Category) (err error) {
	err = r.DB.Create(&category).Error
	return err
}

func (r *categoryRepository) GetCategoryByIdRepository(categoryID int) (category Category, err error) {
	err = r.DB.First(&category, categoryID).Error
	return category, err
}

func (r *categoryRepository) UpdateCategoryRepository(category Category) (err error) {
	err = r.DB.Save(&category).Error
	return err
}

func (r *categoryRepository) DeleteCategoryRepository(category Category) (err error) {
	err = r.DB.Delete(&category).Error
	return err
}
