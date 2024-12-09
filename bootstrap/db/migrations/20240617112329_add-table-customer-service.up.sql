CREATE TABLE customer_service (
    id UUID PRIMARY KEY,
    service_id UUID NOT NULL,
    customer_id UUID NOT NULL REFERENCES customers(id),
    name VARCHAR(255) NOT NULL,
    document VARCHAR(100) NOT NULL,
    birth_date DATE,
    email VARCHAR(255),
    phone VARCHAR(50),
    gender VARCHAR(10),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);