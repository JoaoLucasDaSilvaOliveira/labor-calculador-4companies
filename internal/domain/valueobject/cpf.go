package valueobject

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"strings"

	"github.com/paemuri/brdoc"
)

// error declarations
var (
	ErrInvalidCPF = error_factory.NewError("cpf inválido")
)

type CPF string

func NewCPF(rawCPF string) (CPF, error) {
	//remove spaces
	rawCPF = strings.TrimSpace(rawCPF)

	if !brdoc.IsCPF(rawCPF) {
		return "", fmt.Errorf("%w: %s", ErrInvalidCPF, rawCPF)
	}

	// returns a cast of the CPF, using the return value of the function that extracts only the digits from the CPF string.
	return CPF(onlyDigits(rawCPF)), nil
}

func (c CPF) String() string {
	return string(c)
}
