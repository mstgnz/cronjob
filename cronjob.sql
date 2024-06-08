-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS app_logs_id_seq;

-- Table Definition
CREATE TABLE "public"."app_logs" (
    "id" int8 NOT NULL DEFAULT nextval('app_logs_id_seq'::regclass),
    "level" varchar NOT NULL,
    "message" text NOT NULL,
    "created_at" timestamp NOT NULL DEFAULT now(),
    PRIMARY KEY ("id")
);

-- Column Comment
COMMENT ON COLUMN "public"."app_logs"."level" IS 'info, error, warning, debug';

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS groups_id_seq;

-- Table Definition
CREATE TABLE "public"."groups" (
    "id" int4 NOT NULL DEFAULT nextval('groups_id_seq'::regclass),
    "uid" int4,
    "user_id" int4 NOT NULL,
    "name" varchar NOT NULL,
    "active" bool NOT NULL DEFAULT true,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- Column Comment
COMMENT ON COLUMN "public"."groups"."uid" IS 'parent id';

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS notifications_id_seq;

-- Table Definition
CREATE TABLE "public"."notifications" (
    "id" int4 NOT NULL DEFAULT nextval('notifications_id_seq'::regclass),
    "schedule_id" int4 NOT NULL,
    "is_sms" bool NOT NULL DEFAULT false,
    "is_mail" bool NOT NULL DEFAULT false,
    "title" varchar NOT NULL,
    "content" varchar NOT NULL,
    "active" bool NOT NULL DEFAULT true,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS notify_email_id_seq;

-- Table Definition
CREATE TABLE "public"."notify_email" (
    "id" int4 NOT NULL DEFAULT nextval('notify_email_id_seq'::regclass),
    "notification_id" int4 NOT NULL,
    "email" varchar NOT NULL,
    "active" bool DEFAULT true,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS notify_sms_id_seq;

-- Table Definition
CREATE TABLE "public"."notify_sms" (
    "id" int4 NOT NULL DEFAULT nextval('notify_sms_id_seq'::regclass),
    "notification_id" int4 NOT NULL,
    "phone" varchar NOT NULL,
    "active" bool DEFAULT true,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS request_headers_id_seq;

-- Table Definition
CREATE TABLE "public"."request_headers" (
    "id" int4 NOT NULL DEFAULT nextval('request_headers_id_seq'::regclass),
    "request_id" int4 NOT NULL,
    "header" varchar NOT NULL,
    "active" bool NOT NULL DEFAULT true,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS requests_id_seq;

-- Table Definition
CREATE TABLE "public"."requests" (
    "id" int4 NOT NULL DEFAULT nextval('requests_id_seq'::regclass),
    "user_id" int4 NOT NULL,
    "url" varchar NOT NULL,
    "method" varchar NOT NULL CHECK ((method)::text = ANY (ARRAY['GET'::text, 'POST'::text, 'PUT'::text, 'DELETE'::text, 'PATCH'::text])),
    "content" jsonb,
    "active" bool NOT NULL DEFAULT true,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- Column Comment
COMMENT ON COLUMN "public"."requests"."method" IS 'GET-POST-PUT-DELETE-PATCH';

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS schedule_logs_id_seq;

-- Table Definition
CREATE TABLE "public"."schedule_logs" (
    "id" int4 NOT NULL DEFAULT nextval('schedule_logs_id_seq'::regclass),
    "schedule_id" int4 NOT NULL,
    "started_at" timestamp NOT NULL,
    "finished_at" timestamp NOT NULL,
    "took" float4 NOT NULL,
    "result" text NOT NULL,
    PRIMARY KEY ("id")
);

-- Column Comment
COMMENT ON COLUMN "public"."schedule_logs"."took" IS 'processing time';
COMMENT ON COLUMN "public"."schedule_logs"."result" IS 'endpoint response';

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS schedules_id_seq;

-- Table Definition
CREATE TABLE "public"."schedules" (
    "id" int4 NOT NULL DEFAULT nextval('schedules_id_seq'::regclass),
    "user_id" int4 NOT NULL,
    "group_id" int2 NOT NULL,
    "request_id" int4 NOT NULL,
    "timing" varchar NOT NULL,
    "timeout" int2 DEFAULT 0,
    "retries" int2 DEFAULT 0,
    "running" bool NOT NULL DEFAULT false,
    "active" bool NOT NULL DEFAULT true,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- Column Comment
COMMENT ON COLUMN "public"."schedules"."timing" IS '* * * * *';
COMMENT ON COLUMN "public"."schedules"."timeout" IS 'Enter the maximum time allowed for jobs to complete, 0 to disable';
COMMENT ON COLUMN "public"."schedules"."retries" IS 'Select the number of retries to be attempted before an error is reported';
COMMENT ON COLUMN "public"."schedules"."running" IS 'will be true when triggered and false again when it is done.';
COMMENT ON COLUMN "public"."schedules"."active" IS 'if active, it will be triggered in due time';

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS triggered_id_seq;

-- Table Definition
CREATE TABLE "public"."triggered" (
    "id" int4 NOT NULL DEFAULT nextval('triggered_id_seq'::regclass),
    "schedule_id" int4 NOT NULL,
    PRIMARY KEY ("id")
);

-- Column Comment
COMMENT ON COLUMN "public"."triggered"."schedule_id" IS 'database lock will be used, there can be only one schedule_id record. it will be deleted when the process is complete.';

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS users_id_seq;

-- Table Definition
CREATE TABLE "public"."users" (
    "id" int4 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "fullname" varchar NOT NULL,
    "email" varchar NOT NULL,
    "password" varchar NOT NULL,
    "phone" varchar NOT NULL,
    "is_admin" bool DEFAULT false,
    "active" bool DEFAULT true,
    "last_login" timestamp,
    "created_at" timestamp DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- This script only contains the table creation statements and does not fully represent the table in the database. It's still missing: indices, triggers. Do not use it as a backup.

-- Sequence and defined type
CREATE SEQUENCE IF NOT EXISTS webhooks_id_seq;

-- Table Definition
CREATE TABLE "public"."webhooks" (
    "id" int4 NOT NULL DEFAULT nextval('webhooks_id_seq'::regclass),
    "schedule_id" int4 NOT NULL,
    "request_id" int4 NOT NULL,
    "active" bool,
    "created_at" timestamp NOT NULL DEFAULT now(),
    "updated_at" timestamp,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

ALTER TABLE "public"."groups" ADD FOREIGN KEY ("uid") REFERENCES "public"."groups"("id");
ALTER TABLE "public"."groups" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."notifications" ADD FOREIGN KEY ("schedule_id") REFERENCES "public"."schedules"("id");
ALTER TABLE "public"."notify_email" ADD FOREIGN KEY ("notification_id") REFERENCES "public"."notifications"("id");
ALTER TABLE "public"."notify_sms" ADD FOREIGN KEY ("notification_id") REFERENCES "public"."notifications"("id");
ALTER TABLE "public"."request_headers" ADD FOREIGN KEY ("request_id") REFERENCES "public"."requests"("id");
ALTER TABLE "public"."schedule_logs" ADD FOREIGN KEY ("schedule_id") REFERENCES "public"."schedules"("id") ON DELETE CASCADE;
ALTER TABLE "public"."schedules" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");
ALTER TABLE "public"."schedules" ADD FOREIGN KEY ("group_id") REFERENCES "public"."groups"("id");
ALTER TABLE "public"."schedules" ADD FOREIGN KEY ("request_id") REFERENCES "public"."requests"("id");
ALTER TABLE "public"."triggered" ADD FOREIGN KEY ("schedule_id") REFERENCES "public"."schedules"("id") ON DELETE CASCADE;
ALTER TABLE "public"."webhooks" ADD FOREIGN KEY ("schedule_id") REFERENCES "public"."schedules"("id");
ALTER TABLE "public"."webhooks" ADD FOREIGN KEY ("request_id") REFERENCES "public"."requests"("id");
ALTER TABLE "public"."requests" ADD FOREIGN KEY ("user_id") REFERENCES "public"."users"("id");
