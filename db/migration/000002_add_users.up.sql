
CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "password_hash" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar NOT NULL UNIQUE,
  "password_changed_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamp NOT NULL DEFAULT (now())
)
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");