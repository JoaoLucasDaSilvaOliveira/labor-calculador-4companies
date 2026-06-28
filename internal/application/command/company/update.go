package company

import (
	"labor-calculador-4companies/internal/domain/valueobject"
)

type UpdateCompanyCommand struct {
	IDCompany int
	Name string
	CNPJ valueobject.CNPJ
}
