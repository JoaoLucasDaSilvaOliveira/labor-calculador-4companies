package pages

import (
	"errors"
	"fmt"
	"strings"

	command "labor-calculador-4companies/internal/application/command/employee"
	receiptQuery "labor-calculador-4companies/internal/application/query/receipt"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/valueobject"
	"labor-calculador-4companies/internal/ui/components"
	"labor-calculador-4companies/internal/ui/navigation"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type EmployeeFinderByID interface {
	Execute(command.GetEmployeeById) (*entity.Employee, error)
}

type ReceiptFinder interface {
	Execute(receiptQuery.GetReceiptWithFilter) ([]*entity.Receipt, error)
}

type EmployeeUpdater interface {
	Execute(command.UpdateEmployeeCommand) error
}

type EmployeePageDeps struct {
	EmployeeID  int
	Employee    EmployeeFinderByID
	Receipts    ReceiptFinder
	Updater     EmployeeUpdater
	Navigator   navigation.Navigator
	CompanyName string
	Edit        *navigation.EditSession
}

// NewEmployeePage composes screen 1.2 using the prototype's exact nesting,
// dimensions, receipt actions and bottom action area.
func NewEmployeePage(deps EmployeePageDeps) fyne.CanvasObject {
	return newAsyncPrototypeContent("Carregando informações do funcionário...", func() func() fyne.CanvasObject {
		employee, err := deps.Employee.Execute(command.GetEmployeeById{IDEmployee: deps.EmployeeID})
		if err != nil {
			return func() fyne.CanvasObject {
				return components.NewRecoverableErrorState("Não foi possível carregar as informações do funcionário.", nil)
			}
		}

		receipts, err := deps.Receipts.Execute(receiptQuery.GetReceiptWithFilter{IDEmployee: deps.EmployeeID})
		if err != nil {
			return func() fyne.CanvasObject {
				return components.NewRecoverableErrorState("Não foi possível carregar os recibos deste funcionário.", nil)
			}
		}

		return func() fyne.CanvasObject {
			details := components.NewEmployeeDetailsComponent(employee, deps.CompanyName)
			edit := deps.Edit
			if edit == nil {
				edit = &navigation.EditSession{}
			}

			btnGravar := widget.NewButton("GRAVAR", nil)
			btnCancelar := widget.NewButton("CANCELAR", nil)
			btnGravar.Disable()
			btnCancelar.Disable()
			status := widget.NewLabel("")
			status.Wrapping = fyne.TextWrapWord

			editing := false
			saving := false
			saveGeneration := 0

			cancelEdit := func() {
				saveGeneration++
				details.CancelEdit()
				editing = false
				saving = false
				edit.End()
				btnGravar.Disable()
				btnCancelar.Disable()
				status.SetText("")
			}

			syncPendingChanges := func() {
				if details.HasChanges() {
					edit.Begin(cancelEdit)
					btnGravar.Enable()
					btnCancelar.Enable()
					return
				}

				edit.End()
				btnGravar.Disable()
				btnCancelar.Disable()
			}

			beginEdit := func(target *components.InlineEditableField) {
				if saving {
					return
				}

				details.BeginEdit(target)
				editing = true
				status.SetText("")
				syncPendingChanges()
			}
			details.SetOnActivate(beginEdit)
			details.SetOnChanged(func(*components.InlineEditableField) {
				syncPendingChanges()
			})
			btnCancelar.OnTapped = cancelEdit

			btnGravar.OnTapped = func() {
				if !editing || saving || deps.Updater == nil || !details.HasChanges() {
					return
				}

				firstName := strings.TrimSpace(details.FirstNameField.Text())
				lastName := strings.TrimSpace(details.LastNameField.Text())
				cpf, cpfErr := valueobject.NewCPF(details.CPFField.Text())
				valid := true
				if firstName == "" {
					details.FirstNameField.SetValidationError(fmt.Errorf("%w: informe o primeiro nome", entity.ErrInvalidEmployeeFirstName))
					valid = false
				} else {
					details.FirstNameField.SetValidationError(nil)
				}
				if lastName == "" {
					details.LastNameField.SetValidationError(fmt.Errorf("%w: informe o sobrenome", entity.ErrInvalidEmployeeLastName))
					valid = false
				} else {
					details.LastNameField.SetValidationError(nil)
				}
				if cpfErr != nil {
					details.CPFField.SetValidationError(cpfErr)
					valid = false
				} else {
					details.CPFField.SetValidationError(nil)
				}
				if !valid {
					return
				}

				candidate, err := entity.LoadEmployee(
					deps.EmployeeID,
					employee.CompanyID(),
					firstName,
					lastName,
					cpf.String(),
				)
				if err != nil {
					setEmployeeValidationError(details, err)
					return
				}

				commandToSave := command.UpdateEmployeeCommand{
					IDEmployee: deps.EmployeeID,
					CompanyID:  candidate.CompanyID(),
					FirstName:  candidate.FirstName(),
					LastName:   candidate.LastName(),
					CPF:        cpf,
				}
				generation := saveGeneration
				saving = true
				btnGravar.Disable()
				btnCancelar.Disable()
				status.SetText("Salvando alterações...")

				go func() {
					saveErr := deps.Updater.Execute(commandToSave)
					complete := func() {
						if generation != saveGeneration {
							return
						}
						if saveErr != nil {
							saving = false
							btnGravar.Enable()
							btnCancelar.Enable()
							status.SetText("Não foi possível salvar as alterações. Tente novamente.")
							return
						}

						details.FirstNameField.SetText(candidate.FirstName())
						details.LastNameField.SetText(candidate.LastName())
						details.CPFField.SetText(candidate.CPF())
						details.CommitEdit()
						editing = false
						saving = false
						edit.End()
						btnGravar.Disable()
						btnCancelar.Disable()
						status.SetText("Alterações salvas.")
					}
					if fyne.CurrentApp() == nil {
						complete()
						return
					}
					fyne.Do(complete)
				}()
			}

			receiptListComponent := components.NewReceiptsListComponent(receipts, nil, nil, nil)
			actions := container.NewVBox(
				status,
				uiUtils.VPadding(50),
				container.NewHBox(
					layout.NewSpacer(),
					widget.NewButton("NOVO RECIBO", nil),
					uiUtils.HPadding(10),
					btnGravar,
					uiUtils.HPadding(10),
					btnCancelar,
					uiUtils.HPadding(10),
					widget.NewButton("VOLTAR", func() { deps.Navigator.Back() }),
				),
			)
			employeeListComponentBorder := container.NewBorder(
				uiUtils.VPadding(50),
				actions,
				nil,
				nil,
				receiptListComponent,
			)

			insideBorder := container.NewBorder(details.View, nil, nil, nil, employeeListComponentBorder)
			insideBorderWithPadding := container.NewBorder(
				uiUtils.VPadding(10),
				uiUtils.VPadding(10),
				uiUtils.HPadding(10),
				uiUtils.HPadding(10),
				insideBorder,
			)

			return insideBorderWithPadding
		}
	})
}

func setEmployeeValidationError(details *components.EmployeeDetailsComponent, err error) {
	if err == nil {
		return
	}
	switch {
	case errors.Is(err, entity.ErrInvalidEmployeeFirstName):
		details.FirstNameField.SetValidationError(err)
	case errors.Is(err, entity.ErrInvalidEmployeeLastName):
		details.LastNameField.SetValidationError(err)
	default:
		details.CPFField.SetValidationError(err)
	}
}
