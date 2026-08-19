package components

import (
	"errors"
	"testing"

	"labor-calculador-4companies/internal/domain/entity"

	"fyne.io/fyne/v2/test"
)

func testCompany() (*entity.Company, error) {
	return entity.LoadCompany(10, "Empresa", "11222333000181")
}

func testEmployee() (*entity.Employee, error) {
	return entity.LoadEmployee(25, 10, "Ana", "Silva", "52998224725")
}

func TestInlineEditableFieldStartsAsLabelAndCanCancel(t *testing.T) {
	app := test.NewApp()
	t.Cleanup(app.Quit)

	field := NewInlineEditableField("Empresa original", 200)
	if field.IsEditing() {
		t.Fatal("field starts in edit mode")
	}

	field.Tapped(nil)
	field.SetEditingText("Empresa alterada")
	if !field.IsEditing() {
		t.Fatal("field did not enter edit mode")
	}
	if got := field.Text(); got != "Empresa alterada" {
		t.Fatalf("editing text = %q, want changed value", got)
	}

	field.CancelEdit()
	if field.IsEditing() {
		t.Fatal("field remains in edit mode after cancel")
	}
	if got := field.Text(); got != "Empresa original" {
		t.Fatalf("cancelled text = %q, want original value", got)
	}
}

func TestInlineEditableFieldCommitUpdatesDisplaySnapshot(t *testing.T) {
	field := NewInlineEditableField("Antes", 200)
	field.BeginEdit()
	field.SetEditingText("Depois")
	field.EndEdit()

	if got := field.Text(); got != "Depois" {
		t.Fatalf("committed text = %q, want Depois", got)
	}

	field.BeginEdit()
	field.SetEditingText("Temporário")
	field.CancelEdit()
	if got := field.Text(); got != "Depois" {
		t.Fatalf("cancelled committed text = %q, want Depois", got)
	}
}

func TestInlineEditableFieldUsesOnlyTheActivePresentationState(t *testing.T) {
	field := NewInlineEditableField("Antes", 200)
	readView := field.content.Objects[0]

	if readView != field.readView {
		t.Fatal("field does not start with its read-only presentation")
	}
	if field.HasChanges() {
		t.Fatal("read-only field reports changes")
	}

	field.BeginEdit()
	if got := field.content.Objects[0]; got != field.editView {
		t.Fatal("field did not replace the read-only presentation")
	}
	if field.HasChanges() {
		t.Fatal("entering edit mode alone reports changes")
	}

	field.SetEditingText("Depois")
	if !field.HasChanges() {
		t.Fatal("changed field does not report changes")
	}
	field.SetEditingText("Antes")
	if field.HasChanges() {
		t.Fatal("field reports changes after returning to its snapshot")
	}
}

func TestInlineEditableFieldShowsValidationErrorNearField(t *testing.T) {
	field := NewInlineEditableField("Empresa", 200)
	err := errors.New("nome obrigatório")
	field.SetValidationError(err)

	if !field.validation.Visible() {
		t.Fatal("validation label is hidden")
	}
	if got := field.validation.Text; got != err.Error() {
		t.Fatalf("validation text = %q, want %q", got, err.Error())
	}

	field.SetValidationError(nil)
	if field.validation.Visible() {
		t.Fatal("validation label remains visible after clearing")
	}
}

func TestCompanyDetailsActivatesAllBusinessFields(t *testing.T) {
	company, err := testCompany()
	if err != nil {
		t.Fatal(err)
	}

	details := NewCompanyDetailsComponent(company)
	details.SetOnActivate(func(target *InlineEditableField) {
		details.BeginEdit(target)
	})
	details.NameField.Tapped(nil)

	if !details.NameField.IsEditing() || !details.CNPJField.IsEditing() {
		t.Fatal("activating one company field did not activate all business fields")
	}
}

func TestEditSessionStartsOnlyAfterAnEditableValueChanges(t *testing.T) {
	company, err := testCompany()
	if err != nil {
		t.Fatal(err)
	}

	details := NewCompanyDetailsComponent(company)
	dirty := false
	details.SetOnActivate(func(target *InlineEditableField) {
		details.BeginEdit(target)
	})
	details.SetOnChanged(func(*InlineEditableField) {
		if details.HasChanges() {
			dirty = true
			return
		}
		dirty = false
	})

	details.NameField.Tapped(nil)
	if dirty {
		t.Fatal("clicking a field activated the dirty session")
	}

	details.NameField.SetEditingText("Empresa alterada")
	if !dirty {
		t.Fatal("changing a field did not activate the dirty session")
	}

	details.NameField.SetEditingText(company.Name())
	if dirty {
		t.Fatal("returning to the original value kept the dirty session active")
	}
}

func TestEmployeeDetailsExposeSeparateNameFields(t *testing.T) {
	employee, err := testEmployee()
	if err != nil {
		t.Fatal(err)
	}

	details := NewEmployeeDetailsComponent(employee, "EMPRESA")
	details.SetOnActivate(func(target *InlineEditableField) {
		details.BeginEdit(target)
	})
	details.FirstNameField.Tapped(nil)

	if !details.FirstNameField.IsEditing() || !details.LastNameField.IsEditing() || !details.CPFField.IsEditing() {
		t.Fatal("activating one employee field did not activate all business fields")
	}
}
