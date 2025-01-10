package method

import "gorm.io/gorm"

type Repository interface {
	GetAllMethodRepository() (result []Method, err error)
	CreateMethodRepository(method *Method) (err error)
	GetMethodByIdRepository(methodID int) (method Method, err error)
	UpdateMethodRepository(method Method) (err error)
	DeleteMethodRepository(method Method) (err error)
}

type methodRepository struct {
	DB *gorm.DB
}

func NewMethodRepository(db *gorm.DB) Repository {
	return &methodRepository{
		DB: db,
	}
}

func (r *methodRepository) GetAllMethodRepository() (result []Method, err error) {
	err = r.DB.Find(&result).Error
	return
}

func (r *methodRepository) CreateMethodRepository(method *Method) (err error) {
	err = r.DB.Create(method).Error
	return err
}

func (r *methodRepository) GetMethodByIdRepository(methodID int) (method Method, err error) {
	err = r.DB.First(&method, methodID).Error
	return method, err
}

func (r *methodRepository) UpdateMethodRepository(method Method) (err error) {
	err = r.DB.Save(&method).Error
	return err
}

func (r *methodRepository) DeleteMethodRepository(method Method) (err error) {
	err = r.DB.Delete(&method).Error
	return err
}
