DROP INDEX IF EXISTS idx_evo_core_agents_account_id;

ALTER TABLE evo_core_agents DROP COLUMN IF EXISTS account_id;
