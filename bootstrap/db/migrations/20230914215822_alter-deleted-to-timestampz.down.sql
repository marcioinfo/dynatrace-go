ALTER TABLE "customers"
    ALTER COLUMN "deleted_at" TYPE BOOLEAN USING CASE WHEN "deleted_at" IS NOT NULL THEN TRUE ELSE FALSE END;
ALTER TABLE "customers"
    RENAME COLUMN "deleted_at" TO "deleted";

ALTER TABLE "cards"
    ALTER COLUMN "deleted_at" TYPE BOOLEAN USING CASE WHEN "deleted_at" IS NOT NULL THEN TRUE ELSE FALSE END;
ALTER TABLE "cards"
    RENAME COLUMN "deleted_at" TO "deleted";
