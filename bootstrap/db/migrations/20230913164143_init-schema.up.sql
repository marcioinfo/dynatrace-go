-- Criação das tabelas iniciais

CREATE TABLE IF NOT EXISTS "customers"(
    "id" UUID PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL,
    "document" VARCHAR(255) NOT NULL,
    "birth_date" DATE NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "phone" VARCHAR(255) NOT NULL,
    "gender" VARCHAR(255) NOT NULL,
    "deleted" BOOLEAN NOT NULL DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "customer_tokens"(
    "id" UUID PRIMARY KEY,
    "customer_id" UUID NOT NULL,
    "customer_token" VARCHAR(255) NOT NULL,
    "gateway" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "cards"(
    "id" UUID PRIMARY KEY,
    "customer_id" UUID NOT NULL,
    "fingerprint" VARCHAR(255) NOT NULL,
    "holder" VARCHAR(255) NOT NULL,
    "brand" VARCHAR(255) NOT NULL,
    "first_digits" VARCHAR(255) NOT NULL,
    "last_digits" VARCHAR(255) NOT NULL,
    "expire_year" VARCHAR(255) NOT NULL,
    "expire_month" VARCHAR(255) NOT NULL,
    "deleted" BOOLEAN NOT NULL  DEFAULT false,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "card_tokens"(
    "id" UUID PRIMARY KEY,
    "card_id" UUID NOT NULL,
    "card_token" VARCHAR(255) NOT NULL,
    "gateway" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

ALTER TABLE "customer_tokens" ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id") ON DELETE CASCADE;