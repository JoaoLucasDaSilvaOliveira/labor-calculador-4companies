CREATE TABLE IF NOT EXISTS company (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    cnpj TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_company_updated_at
AFTER UPDATE ON company
FOR EACH ROW
BEGIN
    UPDATE company
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;

CREATE TABLE IF NOT EXISTS employee (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL DEFAULT 0,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    cpf TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER IF NOT EXISTS update_employee_updated_at
AFTER UPDATE ON employee
FOR EACH ROW
BEGIN
    UPDATE employee
    SET updated_at = CURRENT_TIMESTAMP
    WHERE id = OLD.id;
END;

CREATE INDEX IF NOT EXISTS idx_employee_company_id ON employee(company_id);

CREATE TABLE IF NOT EXISTS receipt (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    employee_id INTEGER NOT NULL,
    sumary_description TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (employee_id) REFERENCES employee(id)
);

CREATE TABLE IF NOT EXISTS receipt_items (
    receipt_id INTEGER NOT NULL,
    position INTEGER NOT NULL,
    description TEXT NOT NULL,
    decimal_value TEXT NOT NULL,
    value_type INTEGER NOT NULL,
    FOREIGN KEY (receipt_id) REFERENCES receipt(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_receipt_items_receipt_id ON receipt_items(receipt_id);
CREATE INDEX IF NOT EXISTS idx_receipt_employee_id ON receipt(employee_id);
