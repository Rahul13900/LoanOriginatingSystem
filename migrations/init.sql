-- Loans table
CREATE TABLE IF NOT EXISTS loans (
    loan_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_name TEXT NOT NULL,
    phone TEXT NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    type TEXT NOT NULL,
    status TEXT NOT NULL,
    version BIGINT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Agents table
CREATE TABLE IF NOT EXISTS agents (
    agent_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    manager_id UUID NULL REFERENCES agents(agent_id)
);

-- Agent-Loans assignments
CREATE TABLE IF NOT EXISTS agent_loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agent_id UUID NOT NULL REFERENCES agents(agent_id),
    loan_id UUID NOT NULL REFERENCES loans(loan_id),
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_loans_status ON loans(status);
CREATE INDEX IF NOT EXISTS idx_agent_loans_agent ON agent_loans(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_loans_loan ON agent_loans(loan_id);

-- Extensions (for gen_random_uuid)
CREATE EXTENSION IF NOT EXISTS pgcrypto;