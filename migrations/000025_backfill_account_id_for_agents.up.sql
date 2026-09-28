-- See specs/account-feature-toggles/03-plan.md, Step 6 and Fase 0b
-- (migration 000016). Pre-existing evo_core_agents rows predate account_id
-- and are backfilled onto "Account #1" - the earliest-created row in the
-- accounts table, which this service shares a physical database with (see
-- specs/account-feature-toggles/04-architecture.md and
-- specs/multi-account-tenancy). Mirrors the same "Account #1" backfill
-- convention already used in evo-auth-service-community and
-- evo-ai-crm-community for their own pre-existing rows.
--
-- No-op (0 rows affected) on an installation with no accounts row yet, or
-- with no pre-existing agents - safe to run on a fresh install.
UPDATE evo_core_agents
SET account_id = (SELECT id FROM accounts ORDER BY created_at ASC LIMIT 1)
WHERE account_id IS NULL
  AND EXISTS (SELECT 1 FROM accounts);
