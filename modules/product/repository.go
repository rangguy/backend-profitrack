package product

import (
	"backend-profitrack/modules/category"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProductRepository() (result []Product, err error)
	CreateProductRepository(product *Product) (err error)
	GetProductByIdRepository(productID int) (product Product, err error)
	UpdateProductRepository(product Product) (err error)
	DeleteProductRepository(product Product) (err error)
	BulkCreateProductRepository(products []Product) error
	GetCategoryByNameRepository(name string) (category.Category, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) Repository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) BulkCreateProductRepository(products []Product) error {
	return r.DB.Create(&products).Error
}

func (r *productRepository) GetCategoryByNameRepository(name string) (category.Category, error) {
	var categoryName category.Category
	err := r.DB.Where("name = ?", name).First(&categoryName).Error
	return categoryName, err
}

func (r *productRepository) GetAllProductRepository() (result []Product, err error) {
	err = r.DB.Preload("Category").Find(&result).Error
	return result, err
}

func (r *productRepository) CreateProductRepository(product *Product) (err error) {
	err = r.DB.Create(&product).Error
	return err
}

func (r *productRepository) GetProductByIdRepository(productID int) (product Product, err error) {
	err = r.DB.Preload("Category").First(&product, productID).Error
	return product, err
}

func (r *productRepository) UpdateProductRepository(product Product) (err error) {
	err = r.DB.Save(&product).Error
	return err
}

func (r *productRepository) DeleteProductRepository(product Product) (err error) {
	err = r.DB.Delete(&product).Error
	return err
}
