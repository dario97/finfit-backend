ALTER TABLE public.expense
    ADD COLUMN currency VARCHAR(3) NOT NULL CHECK ( currency <> '') DEFAULT 'ARS';