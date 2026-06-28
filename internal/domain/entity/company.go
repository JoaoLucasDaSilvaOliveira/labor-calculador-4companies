package entity

import (
	"fmt"
	error_factory "labor-calculador-4companies/internal/domain/error"
	"labor-calculador-4companies/internal/domain/valueobject"
	"strings"
)

var (
	ErrInvalidCompanySequencialID = error_factory.NewError("id sequencial da empresa inválido")
	ErrInvalidCompanyName         = error_factory.NewError("nome da empresa inválido")
)

type Company struct {
	sequencialID int
	name         string
	cnpj         string
}

func NewCompany(sequencialID int, name string, cnpj string) (*Company, error) {
	company := &Company{}

	if err := company.SetSequencialID(sequencialID); err != nil {
		return nil, err
	}

	if err := company.SetName(name); err != nil {
		return nil, err
	}

	if err := company.SetCNPJ(cnpj); err != nil {
		return nil, err
	}

	return company, nil
}

func (c *Company) SequencialID() int {
	return c.sequencialID
}

func (c *Company) SetSequencialID(sequencialID int) error {
	if sequencialID <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidCompanySequencialID, sequencialID)
	}

	c.sequencialID = sequencialID
	return nil
}

func (c *Company) Name() string {
	return c.name
}

func (c *Company) SetName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("%w: %s", ErrInvalidCompanyName, name)
	}

	c.name = name
	return nil
}

func (c *Company) CNPJ() string {
	return c.cnpj
}

func (c *Company) SetCNPJ(cnpj string) error {
	validCNPJ, err := valueobject.NewCNPJ(cnpj)
	if err != nil {
		return err
	}

	c.cnpj = validCNPJ.String()
	return nil
}
