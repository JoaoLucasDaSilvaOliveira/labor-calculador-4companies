package pages

import (
	"errors"
	"fmt"
	"strings"

	command "labor-calculador-4companies/internal/application/command/company"
	query "labor-calculador-4companies/internal/application/query/employee"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/valueobject"
	"labor-calculador-4companies/internal/ui/components"
	"labor-calculador-4companies/internal/ui/navigation"
	uiUtils "labor-calculador-4companies/internal/ui/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CompanyFinderByID interface {
	Execute(command.GetCompanyById) (*entity.Company, error)
}

type EmployeeFinder interface {
	Execute(query.GetEmployeeWithFilter) ([]*entity.Employee, error)
}

type CompanyUpdater interface {
	Execute(command.UpdateCompanyCommand) error
}

type CompanyPageDeps struct {
	CompanyID     int
	Company       CompanyFinderByID
	Employees     EmployeeFinder
	Updater       CompanyUpdater
	Navigator     navigation.Navigator
	Edit          *navigation.EditSession
	ReloadSidebar func()
}

// NewCompanyPage composes screen 1.1 using the prototype's exact nesting,
// dimensions and bottom action area.
func NewCompanyPage(deps CompanyPageDeps) fyne.CanvasObject {
	return newAsyncPrototypeContent("Carregando informações da empresa...", func() func() fyne.CanvasObject {
		company, err := deps.Company.Execute(command.GetCompanyById{IDCompany: deps.CompanyID})
		if err != nil {
			return func() fyne.CanvasObject {
				return components.NewRecoverableErrorState("Não foi possível carregar as informações da empresa.", nil)
			}
		}

		employees, err := deps.Employees.Execute(query.GetEmployeeWithFilter{CompanyID: deps.CompanyID})
		if err != nil {
			return func() fyne.CanvasObject {
				return components.NewRecoverableErrorState("Não foi possível carregar os funcionários desta empresa.", nil)
			}
		}

		return func() fyne.CanvasObject {
			details := components.NewCompanyDetailsComponent(company)
			currentCompanyName := company.Name()
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

				nameText := strings.TrimSpace(details.NameField.Text())
				cnpj, cnpjErr := valueobject.NewCNPJ(details.CNPJField.Text())
				valid := true
				if nameText == "" {
					details.NameField.SetValidationError(fmt.Errorf("%w: informe o nome da empresa", entity.ErrInvalidCompanyName))
					valid = false
				} else {
					details.NameField.SetValidationError(nil)
				}
				if cnpjErr != nil {
					details.CNPJField.SetValidationError(cnpjErr)
					valid = false
				} else {
					details.CNPJField.SetValidationError(nil)
				}
				if !valid {
					return
				}

				candidate, err := entity.LoadCompany(deps.CompanyID, nameText, cnpj.String())
				if err != nil {
					setCompanyValidationError(details, err)
					return
				}

				commandToSave := command.UpdateCompanyCommand{
					IDCompany: deps.CompanyID,
					Name:      candidate.Name(),
					CNPJ:      cnpj,
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

						details.NameField.SetText(candidate.Name())
						details.CNPJField.SetText(candidate.CNPJ())
						details.CommitEdit()
						currentCompanyName = candidate.Name()
						editing = false
						saving = false
						edit.End()
						btnGravar.Disable()
						btnCancelar.Disable()
						status.SetText("Alterações salvas.")
						if deps.ReloadSidebar != nil {
							deps.ReloadSidebar()
						}
					}
					if fyne.CurrentApp() == nil {
						complete()
						return
					}
					fyne.Do(complete)
				}()
			}

			employeeListComponent := components.NewEmployeesListComponent(employees, func(employeeID int) {
				if err := deps.Navigator.Push(
					navigation.RouteEmployeeDetails,
					navigation.EmployeeDetailsParams{EmployeeID: employeeID, CompanyName: currentCompanyName},
				); err != nil {
					fyne.LogError("Não foi possível abrir o funcionário selecionado.", err)
				}
			})

			actions := container.NewVBox(
				status,
				uiUtils.VPadding(50),
				container.NewHBox(
					layout.NewSpacer(),
					btnGravar,
					uiUtils.HPadding(10),
					btnCancelar,
					uiUtils.HPadding(10),
					widget.NewButton("VOLTAR", func() { resetToQuickAccess(deps.Navigator) }),
				),
			)
			employeeListComponentBorder := container.NewBorder(
				uiUtils.VPadding(50),
				actions,
				nil,
				nil,
				employeeListComponent,
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

func setCompanyValidationError(details *components.CompanyDetailsComponent, err error) {
	if err == nil {
		return
	}
	if errors.Is(err, entity.ErrInvalidCompanyName) {
		details.NameField.SetValidationError(err)
		return
	}
	details.CNPJField.SetValidationError(err)
}

func resetToQuickAccess(navigator navigation.Navigator) {
	if err := navigator.Reset(navigation.RouteQuickAccess, nil); err != nil {
		fyne.LogError("Não foi possível resetar o workspace.", err)
		return
	}
	if err := navigator.Replace(navigation.RouteQuickAccess, nil); err != nil {
		fyne.LogError("Não foi possível abrir o acesso rápido.", err)
	}
}

func NewCompanyRegistrationPage(navigator navigation.Navigator) fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Cadastrar empresa", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	description := widget.NewLabel("Informe os dados da empresa para começar a organizar os cálculos.")
	description.Wrapping = fyne.TextWrapWord
	content := container.NewVBox(
		title,
		description,
		widget.NewSeparator(),
		widget.NewLabel("O formulário de cadastro será implementado nesta página."),
	)
	return newPageWithBackAction(content, navigator)
}

func pageWithBody(header, body fyne.CanvasObject) fyne.CanvasObject {
	return container.NewPadded(container.NewBorder(header, nil, nil, nil, body))
}

func newPageWithBackAction(content fyne.CanvasObject, navigator navigation.Navigator) fyne.CanvasObject {
	return pageWithBody(
		widget.NewButtonWithIcon("Voltar", theme.NavigateBackIcon(), func() { navigator.Back() }),
		container.NewPadded(content),
	)
}
