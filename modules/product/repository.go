package product

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetAllProductRepository() (result []Product, err error)
	CreateProductRepository(product *Product) (err error)
	GetProductByIdRepository(productID int) (product Product, err error)
	UpdateProductRepository(product Product) (err error)
	DeleteProductRepository(product Product) (err error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) Repository {
	return &productRepository{
		DB: db,
	}
}

func (r *productRepository) GetAllProductRepository() (result []Product, err error) {
	err = r.DB.Find(&result).Error
	return result, err
}

func (r *productRepository) CreateProductRepository(product *Product) (err error) {
	err = r.DB.Create(&product).Error
	return err
}

func (r *productRepository) GetProductByIdRepository(productID int) (product Product, err error) {
	err = r.DB.First(&product, productID).Error
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
