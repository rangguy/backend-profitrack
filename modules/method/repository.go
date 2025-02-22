package method

import "gorm.io/gorm"

type Repository interface {
	GetAllMethodRepository() (result []Method, err error)
	GetMethodByIdRepository(methodID int) (method Method, err error)
	DeleteMethodRepository(method *Method) (err error)
}

type methodRepository struct {
	DB *gorm.DB
}

func NewMethodRepository(db *gorm.DB) Repository {
	repo := &methodRepository{
		DB: db,
	}

	methods, err := repo.GetAllMethodRepository()
	if err != nil {
		return nil
	}

	methodMap := make(map[string]bool)
	for _, method := range methods {
		methodMap[method.Name] = true
	}

	requiredMethods := []string{"SMART", "MOORA"}

	for _, methodName := range requiredMethods {
		if !methodMap[methodName] {
			newMethod := Method{Name: methodName}
			if err = db.Create(&newMethod).Error; err != nil {
				return nil
			}
		}
	}

	return repo
}

func (r *methodRepository) GetAllMethodRepository() (result []Method, err error) {
	err = r.DB.Find(&result).Error
	return
}

func (r *methodRepository) GetMethodByIdRepository(methodID int) (method Method, err error) {
	err = r.DB.First(&method, methodID).Error
	return method, err
}

func (r *methodRepository) DeleteMethodRepository(method *Method) (err error) {
	err = r.DB.Delete(method).Error
	return err
}
