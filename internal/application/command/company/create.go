package company

import "labor-calculador-4companies/internal/domain/valueobject"

type CreateCompanyCommand struct {
	Name string
	CNPJ valueobject.CNPJ
}
