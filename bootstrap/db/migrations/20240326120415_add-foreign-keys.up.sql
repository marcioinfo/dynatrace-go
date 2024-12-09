ALTER TABLE "cards" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id") ON DELETE CASCADE;
ALTER TABLE "card_tokens" ADD FOREIGN KEY ("card_id") REFERENCES "cards" ("id") ON DELETE CASCADE;
