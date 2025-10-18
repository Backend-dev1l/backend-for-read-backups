-- Create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "username" text NOT NULL,
  "email" citext NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
-- Create index "idx_user_email" to table: "users"
CREATE INDEX "idx_user_email" ON "users" ("email");
-- Create "user_progress" table
CREATE TABLE "user_progress" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "word_id" uuid NOT NULL,
  "correct_count" integer NOT NULL DEFAULT 0,
  "incorrect_count" integer NOT NULL DEFAULT 0,
  "last_attempt" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "user_progress_user_id_word_id_key" UNIQUE ("user_id", "word_id"),
  CONSTRAINT "user_progress_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_progress_correct_count_check" CHECK (correct_count >= 0),
  CONSTRAINT "user_progress_incorrect_count_check" CHECK (incorrect_count >= 0)
);
-- Create "user_sessions" table
CREATE TABLE "user_sessions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "started_at" timestamptz NOT NULL DEFAULT now(),
  "ended_at" timestamptz NULL,
  "status" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "user_sessions_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_sessions_status_check" CHECK (status = ANY (ARRAY['active'::text, 'completed'::text, 'abandoned'::text]))
);
-- Create index "idx_user_sessions_active" to table: "user_sessions"
CREATE INDEX "idx_user_sessions_active" ON "user_sessions" ("status") WHERE (status = 'active'::text);
-- Create index "idx_user_sessions_user_id" to table: "user_sessions"
CREATE INDEX "idx_user_sessions_user_id" ON "user_sessions" ("user_id");
-- Create "user_statistics" table
CREATE TABLE "user_statistics" (
  "user_id" uuid NOT NULL,
  "total_words_learned" integer NOT NULL DEFAULT 0,
  "accuracy" numeric(5,2) NOT NULL DEFAULT 0,
  "total_time" integer NOT NULL DEFAULT 0,
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("user_id"),
  CONSTRAINT "user_statistics_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "user_statistics_accuracy_check" CHECK ((accuracy >= (0)::numeric) AND (accuracy <= (100)::numeric)),
  CONSTRAINT "user_statistics_total_time_check" CHECK (total_time >= 0),
  CONSTRAINT "user_statistics_total_words_learned_check" CHECK (total_words_learned >= 0)
);
-- Create "user_word_sets" table
CREATE TABLE "user_word_sets" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "word_set_id" uuid NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "user_word_sets_user_id_word_set_id_key" UNIQUE ("user_id", "word_set_id"),
  CONSTRAINT "user_word_sets_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
