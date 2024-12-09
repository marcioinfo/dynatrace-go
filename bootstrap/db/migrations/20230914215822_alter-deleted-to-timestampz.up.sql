-- Alterar coluna `deleted` para `deleted_at` e mudar o tipo na tabela `customers`
ALTER TABLE "customers" ALTER COLUMN "deleted" DROP DEFAULT;
ALTER TABLE "customers"
    RENAME COLUMN "deleted" TO "deleted_at";
ALTER TABLE "customers"
    ALTER COLUMN "deleted_at" TYPE TIMESTAMPTZ USING CASE WHEN "deleted_at" THEN CURRENT_TIMESTAMP ELSE NULL END;
ALTER TABLE "customers" ALTER COLUMN "deleted_at" DROP NOT NULL;

-- Alterar coluna `deleted` para `deleted_at` e mudar o tipo na tabela `cards`
ALTER TABLE "cards" ALTER COLUMN "deleted" DROP DEFAULT;
ALTER TABLE "cards"
    RENAME COLUMN "deleted" TO "deleted_at";
ALTER TABLE "cards"
    ALTER COLUMN "deleted_at" TYPE TIMESTAMPTZ USING CASE WHEN "deleted_at" THEN CURRENT_TIMESTAMP ELSE NULL END;
ALTER TABLE "cards" ALTER COLUMN "deleted_at" DROP NOT NULL;
