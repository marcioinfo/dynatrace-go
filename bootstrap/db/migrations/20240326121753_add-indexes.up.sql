CREATE INDEX IF NOT EXISTS "idx_customers_document_deleted_at" ON customers("document", "deleted_at");

CREATE INDEX IF NOT EXISTS "idx_cards_fingerprint" ON cards("fingerprint");
CREATE INDEX IF NOT EXISTS "idx_cards_customer_id" ON cards("customer_id");

CREATE INDEX IF NOT EXISTS "idx_card_tokens_card_id_gateway" ON card_tokens("card_id", "gateway");

CREATE INDEX IF NOT EXISTS "idx_customer_token_customer_id_gateway" ON customer_tokens("customer_id", "gateway");
