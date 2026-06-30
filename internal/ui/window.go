package ui

import (
	"fmt"
	"strings"

	command "labor-calculador-4companies/internal/application/command/company"
	query "labor-calculador-4companies/internal/application/query/company"
	usecase "labor-calculador-4companies/internal/application/usecase/company"
	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/domain/valueobject"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type CompanyWindowUsecases struct {
	Create *usecase.CreateCompanyUsecase
	Update *usecase.UpdateCompanyUsecase
	Delete *usecase.DeleteCompanyUsecase
	Get    *usecase.GetCompanyUsecase
}

type companyWindow struct {
	window fyne.Window

	createCompany *usecase.CreateCompanyUsecase
	updateCompany *usecase.UpdateCompanyUsecase
	deleteCompany *usecase.DeleteCompanyUsecase
	getCompany    *usecase.GetCompanyUsecase

	companies []*entity.Company
	selected  *entity.Company

	searchName *widget.Entry
	searchCNPJ *widget.Entry
	name       *widget.Entry
	cnpj       *widget.Entry
	status     *widget.Label
	selectedID *widget.Label
	list       *widget.List
}

func NewCompanyWindow(app fyne.App, usecases CompanyWindowUsecases) fyne.Window {
	w := app.NewWindow("Empresas")
	screen := &companyWindow{
		window:        w,
		createCompany: usecases.Create,
		updateCompany: usecases.Update,
		deleteCompany: usecases.Delete,
		getCompany:    usecases.Get,
	}

	w.SetContent(screen.buildContent())
	w.Resize(fyne.NewSize(1120, 680))
	screen.refresh()

	return w
}

func (c *companyWindow) buildContent() fyne.CanvasObject {
	c.buildFields()
	c.list = c.buildCompanyList()

	return container.NewBorder(
		c.buildHeader(),
		c.buildFooter(),
		c.buildSidebar(),
		nil,
		container.NewHSplit(c.buildListCard(), c.buildFormCard()),
	)
}

func (c *companyWindow) buildFields() {
	c.searchName = widget.NewEntry()
	c.searchName.SetPlaceHolder("Buscar por nome")
	c.searchCNPJ = widget.NewEntry()
	c.searchCNPJ.SetPlaceHolder("Filtrar por CNPJ")
	c.name = widget.NewEntry()
	c.name.SetPlaceHolder("Razão social")
	c.cnpj = widget.NewEntry()
	c.cnpj.SetPlaceHolder("Somente números ou formatado")
	c.status = widget.NewLabel("Pronto")
	c.selectedID = widget.NewLabel("Nenhuma empresa selecionada")
}

func (c *companyWindow) buildHeader() fyne.CanvasObject {
	title := widget.NewLabelWithStyle("Empresas", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabel("Cadastro, consulta e manutenção de empresas")

	newButton := widget.NewButtonWithIcon("Nova empresa", theme.ContentAddIcon(), c.clearForm)
	newButton.Importance = widget.HighImportance
	refreshButton := widget.NewButtonWithIcon("Atualizar", theme.ViewRefreshIcon(), c.refresh)

	return container.NewBorder(
		nil,
		nil,
		container.NewVBox(title, subtitle),
		container.NewHBox(refreshButton, newButton),
		nil,
	)
}

func (c *companyWindow) buildSidebar() fyne.CanvasObject {
	companies := widget.NewButtonWithIcon("Empresas", theme.HomeIcon(), func() {})
	companies.Importance = widget.HighImportance
	reports := widget.NewButtonWithIcon("Relatórios", theme.DocumentIcon(), func() {
		c.setStatus("Relatórios ainda não implementado")
	})
	settings := widget.NewButtonWithIcon("Configurações", theme.SettingsIcon(), func() {
		c.setStatus("Configurações ainda não implementado")
	})

	menu := container.NewVBox(
		widget.NewLabelWithStyle("Menu", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		companies,
		reports,
		settings,
		widget.NewSeparator(),
		widget.NewLabel("Opções"),
		widget.NewButtonWithIcon("Limpar filtros", theme.ContentClearIcon(), c.clearFilters),
		widget.NewButtonWithIcon("Limpar formulário", theme.DeleteIcon(), c.clearForm),
	)

	return c.buildCard("", "", menu)
}

func (c *companyWindow) buildListCard() fyne.CanvasObject {
	searchButton := widget.NewButtonWithIcon("Buscar", theme.SearchIcon(), c.refresh)
	searchButton.Importance = widget.HighImportance
	clearButton := widget.NewButtonWithIcon("Limpar", theme.ContentClearIcon(), c.clearFilters)

	filters := container.NewBorder(
		nil,
		nil,
		nil,
		container.NewHBox(searchButton, clearButton),
		container.NewGridWithColumns(2, c.searchName, c.searchCNPJ),
	)

	content := container.NewBorder(
		container.NewVBox(filters, widget.NewSeparator()),
		nil,
		nil,
		nil,
		c.list,
	)

	return c.buildCard("Lista de empresas", "Selecione uma empresa para editar", content)
}

func (c *companyWindow) buildFormCard() fyne.CanvasObject {
	return c.buildCard("Dados da empresa", "Preencha os campos e escolha uma ação", c.buildForm())
}

func (c *companyWindow) buildForm() fyne.CanvasObject {
	form := widget.NewForm(
		widget.NewFormItem("Nome", c.name),
		widget.NewFormItem("CNPJ", c.cnpj),
	)

	createButton := widget.NewButtonWithIcon("Criar", theme.ContentAddIcon(), c.create)
	createButton.Importance = widget.HighImportance
	updateButton := widget.NewButtonWithIcon("Salvar", theme.DocumentSaveIcon(), c.update)
	deleteButton := widget.NewButtonWithIcon("Excluir", theme.DeleteIcon(), c.delete)
	deleteButton.Importance = widget.DangerImportance
	newButton := widget.NewButtonWithIcon("Novo", theme.FileIcon(), c.clearForm)

	actions := container.NewGridWithColumns(2, createButton, updateButton)
	secondary := container.NewGridWithColumns(2, newButton, deleteButton)

	return container.NewBorder(
		container.NewVBox(c.selectedID, widget.NewSeparator()),
		container.NewVBox(actions, secondary),
		nil,
		nil,
		form,
	)
}

func (c *companyWindow) buildFooter() fyne.CanvasObject {
	return container.NewBorder(
		widget.NewSeparator(),
		nil,
		widget.NewIcon(theme.InfoIcon()),
		nil,
		c.status,
	)
}

func (c *companyWindow) buildCard(title, subtitle string, content fyne.CanvasObject) fyne.CanvasObject {
	return widget.NewCard(title, subtitle, container.NewPadded(content))
}

func (c *companyWindow) buildCompanyList() *widget.List {
	list := widget.NewList(
		func() int { return len(c.companies) },
		func() fyne.CanvasObject {
			name := widget.NewLabelWithStyle("Nome", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			doc := widget.NewLabel("CNPJ")
			id := widget.NewLabel("#")
			return container.NewBorder(nil, nil, id, nil, container.NewVBox(name, doc))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			row := item.(*fyne.Container)
			idLabel := row.Objects[1].(*widget.Label)
			labels := row.Objects[0].(*fyne.Container).Objects
			company := c.companies[id]
			idLabel.SetText(fmt.Sprintf("#%d", company.SequencialID()))
			labels[0].(*widget.Label).SetText(company.Name())
			labels[1].(*widget.Label).SetText(formatCNPJ(company.CNPJ()))
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		c.selected = c.companies[id]
		c.name.SetText(c.selected.Name())
		c.cnpj.SetText(c.selected.CNPJ())
		c.selectedID.SetText(fmt.Sprintf("Editando empresa #%d", c.selected.SequencialID()))
		c.setStatus(fmt.Sprintf("Empresa #%d selecionada", c.selected.SequencialID()))
	}
	return list
}

func (c *companyWindow) refresh() {
	filter, err := c.filter()
	if err != nil {
		c.showError(err)
		return
	}

	companies, err := c.getCompany.Execute(filter)
	if err != nil {
		c.showError(err)
		return
	}

	c.companies = companies
	c.selected = nil
	if c.list != nil {
		c.list.UnselectAll()
		c.list.Refresh()
	}
	c.setStatus(fmt.Sprintf("%d empresa(s) encontrada(s)", len(companies)))
}

func (c *companyWindow) filter() (query.GetCompanyWithFilter, error) {
	filter := query.GetCompanyWithFilter{Name: strings.TrimSpace(c.searchName.Text)}
	if strings.TrimSpace(c.searchCNPJ.Text) != "" {
		cnpj, err := valueobject.NewCNPJ(c.searchCNPJ.Text)
		if err != nil {
			return filter, err
		}
		filter.CNPJ = cnpj
	}
	return filter, nil
}

func (c *companyWindow) create() {
	name, cnpj, err := c.formValues()
	if err != nil {
		c.showError(err)
		return
	}

	if err := c.createCompany.Execute(command.CreateCompanyCommand{Name: name, CNPJ: cnpj}); err != nil {
		c.showError(err)
		return
	}

	c.clearForm()
	c.refresh()
	c.setStatus("Empresa criada")
}

func (c *companyWindow) update() {
	if c.selected == nil {
		c.showError(fmt.Errorf("selecione uma empresa para atualizar"))
		return
	}

	name, cnpj, err := c.formValues()
	if err != nil {
		c.showError(err)
		return
	}

	err = c.updateCompany.Execute(command.UpdateCompanyCommand{
		IDCompany: c.selected.SequencialID(),
		Name:      name,
		CNPJ:      cnpj,
	})
	if err != nil {
		c.showError(err)
		return
	}

	c.refresh()
	c.setStatus("Empresa atualizada")
}

func (c *companyWindow) delete() {
	if c.selected == nil {
		c.showError(fmt.Errorf("selecione uma empresa para excluir"))
		return
	}

	id := c.selected.SequencialID()
	dialog.ShowConfirm("Excluir empresa", "Excluir "+c.selected.Name()+"?", func(confirmed bool) {
		if !confirmed {
			return
		}
		if err := c.deleteCompany.Execute(command.DeleteCompanyCommand{IDCompany: id}); err != nil {
			c.showError(err)
			return
		}
		c.clearForm()
		c.refresh()
		c.setStatus("Empresa excluída")
	}, c.window)
}

func (c *companyWindow) formValues() (string, valueobject.CNPJ, error) {
	name := strings.TrimSpace(c.name.Text)
	if name == "" {
		return "", "", fmt.Errorf("nome da empresa é obrigatório")
	}

	cnpj, err := valueobject.NewCNPJ(c.cnpj.Text)
	return name, cnpj, err
}

func (c *companyWindow) clearFilters() {
	c.searchName.SetText("")
	c.searchCNPJ.SetText("")
	c.refresh()
}

func (c *companyWindow) clearForm() {
	c.selected = nil
	c.name.SetText("")
	c.cnpj.SetText("")
	c.selectedID.SetText("Nenhuma empresa selecionada")
	if c.list != nil {
		c.list.UnselectAll()
	}
	c.setStatus("Pronto")
}

func (c *companyWindow) showError(err error) {
	c.setStatus(err.Error())
	dialog.ShowError(err, c.window)
}

func (c *companyWindow) setStatus(message string) {
	if c.status != nil {
		c.status.SetText(message)
	}
}

func formatCNPJ(raw string) string {
	digits := strings.TrimSpace(raw)
	if len(digits) != 14 {
		return raw
	}

	parts := []string{
		digits[0:2],
		digits[2:5],
		digits[5:8],
		digits[8:12],
		digits[12:14],
	}

	return parts[0] + "." + parts[1] + "." + parts[2] + "/" + parts[3] + "-" + parts[4]
}
