CREATE TABLE IF NOT EXISTS expense(
    id uuid PRIMARY KEY,
    expense_type_id uuid NOT NULL,
    amount decimal NOT NULL,
    description VARCHAR(40),
    expense_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT fk_expense_type
        FOREIGN KEY(expense_type_id)
            REFERENCES expense_type(id)
);