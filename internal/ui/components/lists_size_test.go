package components

import (
	"testing"

	"labor-calculador-4companies/internal/domain/entity"
	"labor-calculador-4companies/internal/ui/assets"

	"fyne.io/fyne/v2/app"
)

func TestEmployeesListComponentKeepsPrototypeHeight(t *testing.T) {
	app.NewWithID("components-lists-size-employees")

	employee, err := entity.LoadEmployee(1, 1, "Joao", "Silva", "529.982.247-25")
	if err != nil {
		t.Fatalf("load employee: %v", err)
	}

	rowHeight := newEmployeeRow().MinSize().Height
	want := rowHeight + 60
	got := NewEmployeesListComponent([]*entity.Employee{employee}, nil).MinSize().Height
	if got != want {
		t.Fatalf("one employee panel height = %v, want %v", got, want)
	}

	tooManyEmployees := make([]*entity.Employee, 30)
	for i := range tooManyEmployees {
		tooManyEmployees[i] = employee
	}
	const maxListHeight float32 = 300
	want = maxListHeight + 60
	got = NewEmployeesListComponent(tooManyEmployees, nil).MinSize().Height
	if got != want {
		t.Fatalf("capped employee panel height = %v, want %v", got, want)
	}
}

func TestReceiptsListComponentKeepsPrototypeHeight(t *testing.T) {
	app.NewWithID("components-lists-size-receipts")

	receipt, err := entity.LoadReceipt(1, 1, "Folha de janeiro", nil)
	if err != nil {
		t.Fatalf("load receipt: %v", err)
	}

	rowHeight := newReceiptRow([]receiptAction{
		{icon: assets.BinocularsIcon},
		{icon: assets.EditIcon},
		{icon: assets.DeleteDocumentIcon},
	}).MinSize().Height
	want := rowHeight + 60
	got := NewReceiptsListComponent([]*entity.Receipt{receipt}, nil, nil, nil).MinSize().Height
	if got != want {
		t.Fatalf("one receipt panel height = %v, want %v", got, want)
	}

	tooManyReceipts := make([]*entity.Receipt, 30)
	for i := range tooManyReceipts {
		tooManyReceipts[i] = receipt
	}
	const maxListHeight float32 = 300
	want = maxListHeight + 60
	got = NewReceiptsListComponent(tooManyReceipts, nil, nil, nil).MinSize().Height
	if got != want {
		t.Fatalf("capped receipt panel height = %v, want %v", got, want)
	}
}
