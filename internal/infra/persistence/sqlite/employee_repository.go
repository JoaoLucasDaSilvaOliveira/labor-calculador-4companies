package sqlite

import (
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/repository"
	"labor-calculador-4companies/internal/domain/valueobject"
	"time"

	"gorm.io/gorm"
)

type EmployeeRepository struct {
	db *gorm.DB
}

type employeeModel struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	CPF       string    `gorm:"column:cpf"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (employeeModel) TableName() string {
	return "employee"
}

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (r *EmployeeRepository) Create(firstName string, lastName string, cpf valueobject.CPF) error {
	model := employeeModel{
		FirstName: firstName,
		LastName:  lastName,
		CPF:       cpf.String(),
	}

	return r.db.Create(&model).Error
}

func (r *EmployeeRepository) Update(employee *entity.Employee) error {
	model := employeeModel{
		ID:        employee.GetId(),
		FirstName: employee.FirstName(),
		LastName:  employee.LastName(),
		CPF:       employee.CPF(),
	}

	return r.db.Save(&model).Error
}

func (r *EmployeeRepository) Delete(sequencialID int) error {
	return r.db.Delete(&employeeModel{}, sequencialID).Error
}

func (r *EmployeeRepository) Get(filter query.GetEmployeeWithFilter) ([]*entity.Employee, error) {
	dbQuery := r.db.Model(&employeeModel{})

	if filter.IDEmployee > 0 {
		dbQuery = dbQuery.Where("id = ?", filter.IDEmployee)
	}

	if filter.FirstName != "" {
		dbQuery = dbQuery.Where("first_name LIKE ? COLLATE NOCASE", "%"+filter.FirstName+"%")
	}

	if filter.LastName != "" {
		dbQuery = dbQuery.Where("last_name LIKE ? COLLATE NOCASE", "%"+filter.LastName+"%")
	}

	if filter.CPF != "" {
		dbQuery = dbQuery.Where("cpf = ?", filter.CPF.String())
	}

	var models []employeeModel
	if err := dbQuery.Find(&models).Error; err != nil {
		return nil, err
	}

	return toEmployeeEntities(models)
}

func (r *EmployeeRepository) GetByID(sequencialID int) (*entity.Employee, error) {
	var model employeeModel
	if err := r.db.First(&model, sequencialID).Error; err != nil {
		return nil, err
	}

	return toEmployeeEntity(model)
}

func toEmployeeEntities(models []employeeModel) ([]*entity.Employee, error) {
	employees := make([]*entity.Employee, 0, len(models))

	for _, model := range models {
		employee, err := toEmployeeEntity(model)
		if err != nil {
			return nil, err
		}

		employees = append(employees, employee)
	}

	return employees, nil
}

func toEmployeeEntity(model employeeModel) (*entity.Employee, error) {
	return entity.NewEmployee(model.ID, model.FirstName, model.LastName, model.CPF)
}
