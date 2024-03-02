-- Create enum type "gender"
CREATE TYPE "public"."gender" AS ENUM ('MALE', 'FEMALE');
-- Create enum type "identity_provider"
CREATE TYPE "public"."identity_provider" AS ENUM ('LOCAL', 'GOOGLE', 'APPLE');
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigint NOT NULL,
  "first_name" text NOT NULL,
  "last_name" text NULL,
  "email" text NULL,
  "phone" text NULL,
  "is_email_verified" boolean NULL,
  "is_phone_verified" boolean NULL,
  "is_bot" boolean NULL,
  "gender" "public"."gender" NULL,
  "identity_provider" "public"."identity_provider" NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_email" to table: "users"
CREATE UNIQUE INDEX "idx_email" ON "public"."users" ("email") WHERE (email IS NOT NULL);
-- Create index "idx_phone" to table: "users"
CREATE UNIQUE INDEX "idx_phone" ON "public"."users" ("phone") WHERE (phone IS NOT NULL);
