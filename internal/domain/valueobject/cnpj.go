package valueobject

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"strings"
	"github.com/paemuri/brdoc"
)

//error declarations
var (
	ErrInvalidCNPJ = error_factory.NewError("cnpj inválido")
)

type CNPJ string

func NewCNPJ(rawCNPJ string) (CNPJ, error) {
	rawCNPJ = strings.TrimSpace(rawCNPJ)

	if !brdoc.IsCNPJ(rawCNPJ) {
		return "", fmt.Errorf("%w: %s", ErrInvalidCNPJ, rawCNPJ)
	}

	return CNPJ(onlyDigits(rawCNPJ)), nil
}

func (c CNPJ) String() string {
	return string(c)
}