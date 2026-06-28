package company

import "labor-calculador-4companies/internal/domain/valueobject"

type GetCompanyWithFilter struct {
	IDCompany int
	Name      string
	CNPJ      valueobject.CNPJ
}
