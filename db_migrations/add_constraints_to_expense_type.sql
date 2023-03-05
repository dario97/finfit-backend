ALTER TABLE public.expense_type
    ALTER COLUMN name TYPE VARCHAR(32),
    ALTER COLUMN name SET NOT NULL,
    add constraint expense_type_name_unique_constraint unique (name);
