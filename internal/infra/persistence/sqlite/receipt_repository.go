package sqlite

import (
	"fmt"
	query "labor-calculador-4companies/internal/application/query/receipt"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/repository"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type ReceiptRepository struct {
	db *gorm.DB
}

type receiptModel struct {
	ID         int `gorm:"column:id;primaryKey;autoIncrement"`
	IDEmployee int `gorm:"column:employee_id"`
}

func (receiptModel) TableName() string {
	return "receipt"
}

type receiptItemModel struct {
	ReceiptID    int    `gorm:"column:receipt_id"`
	Position     int    `gorm:"column:position"`
	Description  string `gorm:"column:description"`
	DecimalValue string `gorm:"column:decimal_value"`
	ValueType    int    `gorm:"column:value_type"`
}

func (receiptItemModel) TableName() string {
	return "receipt_items"
}

type persistedReceiptItem struct {
	value entity.ReceiptValueAndDescription
}

func (i persistedReceiptItem) Value() entity.ReceiptValueAndDescription {
	return i.value
}

func NewReceiptRepository(db *gorm.DB) repository.ReceiptRepository {
	return &ReceiptRepository{db: db}
}

func (r *ReceiptRepository) Create(receipt *entity.Receipt) error {
	model := receiptModel{IDEmployee: receipt.GetEmployeeId()}

	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&model).Error; err != nil {
			return err
		}

		return createReceiptItems(tx, model.ID, receipt.GetInformationsFromItems())
	}); err != nil {
		return err
	}

	return receipt.SetId(model.ID)
}

func (r *ReceiptRepository) Update(receipt *entity.Receipt) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		model := receiptModel{
			ID:         receipt.GetId(),
			IDEmployee: receipt.GetEmployeeId(),
		}
		if err := tx.Save(&model).Error; err != nil {
			return err
		}

		if err := tx.Where("receipt_id = ?", receipt.GetId()).Delete(&receiptItemModel{}).Error; err != nil {
			return err
		}

		return createReceiptItems(tx, receipt.GetId(), receipt.GetInformationsFromItems())
	})
}

func (r *ReceiptRepository) Delete(sequencialID int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("receipt_id = ?", sequencialID).Delete(&receiptItemModel{}).Error; err != nil {
			return err
		}

		return tx.Delete(&receiptModel{}, sequencialID).Error
	})
}

func (r *ReceiptRepository) GetByEmployeeIdWithFilter(filter query.GetReceiptWithFilter) ([]*entity.Receipt, error) {
	if filter.IDEmployee <= 0 {
		return nil, fmt.Errorf("é preciso informar o id do empregado")
	}

	dbQuery := r.db.Model(&receiptModel{}).Where("receipt.employee_id = ?", filter.IDEmployee)
	if filter.Description != "" {
		dbQuery = dbQuery.
			Joins("JOIN receipt_items ON receipt_items.receipt_id = receipt.id").
			Where("receipt_items.description LIKE ? COLLATE NOCASE", "%"+filter.Description+"%")
	}

	var models []receiptModel
	if err := dbQuery.Distinct().Find(&models).Error; err != nil {
		return nil, err
	}

	receipts := make([]*entity.Receipt, 0, len(models))
	for _, model := range models {
		receipt, err := r.GetByID(model.ID)
		if err != nil {
			return nil, err
		}

		receipts = append(receipts, receipt)
	}

	return receipts, nil
}

func (r *ReceiptRepository) GetByID(sequencialID int) (*entity.Receipt, error) {
	var model receiptModel
	if err := r.db.First(&model, sequencialID).Error; err != nil {
		return nil, err
	}

	var itemModels []receiptItemModel
	if err := r.db.Where("receipt_id = ?", model.ID).Order("position ASC").Find(&itemModels).Error; err != nil {
		return nil, err
	}

	items, err := toReceiptItems(itemModels)
	if err != nil {
		return nil, err
	}

	return entity.LoadReceipt(model.ID, model.IDEmployee, items)
}

func createReceiptItems(db *gorm.DB, receiptID int, values []entity.ReceiptValueAndDescription) error {
	if len(values) == 0 {
		return nil
	}

	models := make([]receiptItemModel, 0, len(values))
	for position, value := range values {
		models = append(models, receiptItemModel{
			ReceiptID:    receiptID,
			Position:     position,
			Description:  value.Description,
			DecimalValue: value.DecimalValue.String(),
			ValueType:    int(value.ValueType),
		})
	}

	return db.Create(&models).Error
}

func toReceiptItems(models []receiptItemModel) ([]entity.ReceiptItem, error) {
	items := make([]entity.ReceiptItem, 0, len(models))

	for _, model := range models {
		amount, err := decimal.NewFromString(model.DecimalValue)
		if err != nil {
			return nil, fmt.Errorf("valor decimal inválido no item do recibo: %w", err)
		}

		items = append(items, persistedReceiptItem{
			value: entity.ReceiptValueAndDescription{
				Description:  model.Description,
				DecimalValue: amount,
				ValueType:    entity.ReceiptValueType(model.ValueType),
			},
		})
	}

	return items, nil
}
