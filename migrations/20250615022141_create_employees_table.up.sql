CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE employee_role_enum AS ENUM (
    'FIELD_VALIDATOR',
    'FIELD_OFFICER',
    'ADMIN'
    );

CREATE TABLE employees (
                           id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                           username VARCHAR(50) NOT NULL UNIQUE,
                           email VARCHAR(100) NOT NULL UNIQUE,
                           password_hash VARCHAR(255) NOT NULL,
                           full_name VARCHAR(100) NOT NULL,
                           phone_number VARCHAR(15),
                           employee_role employee_role_enum NOT NULL,
                           department VARCHAR(50),
                           is_active BOOLEAN DEFAULT true,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_employees_username ON employees(username);
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_role ON employees(employee_role);
CREATE INDEX idx_employees_is_active ON employees(is_active);