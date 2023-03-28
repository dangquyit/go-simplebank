CREATE TABLE "users" (
                         "username" varchar PRIMARY KEY,
                         "hashed_password" varchar NOT NULL,
                         "full_name" varchar NOT NULL,
                         "email" varchar UNIQUE NOT NULL,
                         "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00+00',
                         "created_at" timestamp NOT NULL DEFAULT (now()),
                         "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD CONSTRAINT "accounts_owner_fkey" FOREIGN KEY ("owner") REFERENCES "users" ("username");

CREATE UNIQUE INDEX "accounts_owner_currency_idx" ON "accounts" ("owner", "currency");