ALTER TABLE evo_core_agents ADD COLUMN account_id UUID;

CREATE INDEX IF NOT EXISTS idx_evo_core_agents_account_id ON evo_core_agents (account_id);
