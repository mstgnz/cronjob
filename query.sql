-- USER_EXISTS_WITH_EMAIL
SELECT count(*) FROM users WHERE email=$1;

-- USER_GET_WITH_ID
SELECT id, fullname, email, is_admin, password FROM users WHERE id=$1;

-- USER_GET_WITH_EMAIL
SELECT id, fullname, email, is_admin, password FROM users WHERE email=$1;

-- USER_GET_WITH_SCHEDULE
SELECT u.id, u.fullname, s.* FROM users as u JOIN schedules as s ON s.user_id=u.id WHERE email=$1;

-- USER_UPDATE
UPDATE users SET fullname=$1 updated_at=$2 WHERE id=$3;

-- USER_INSERT
INSERT INTO users (fullname,email,password,is_admin) VALUES ($1,$2,$3,$4);

-- USER_DELETE
UPDATE users SET deleted_at=$1 WHERE id=$2;

-- USER_LAST_LOGIN
UPDATE users SET last_login=$1 WHERE id=$2;

-- SCHEDULES
SELECT * FROM schedules OFFSET $1 LIMIT $2;

-- SCHEDULE_GET_WITH_USER
SELECT * FROM schedules WHERE user_id=$1;

-- SCHEDULE_UPDATE
UPDATE schedules SET timing=$1, active=$2, running=$3, path=$4, updated_at=$5 WHERE id=$6;

-- SCHEDULE_INSERT
INSERT INTO schedules (timing,active,running,path,send_mail,user_id) VALUES ($1,$2,$3,$4,$5,$6);

-- SCHEDULE_DELETE
UPDATE schedules SET deleted_at=$1 WHERE id=$2;

-- SCHEDULE_LOGS
SELECT * FROM schedule_logs OFFSET $1 LIMIT $2;

-- SCHEDULE_LOG_GET_WITH_SCHEDULE
SELECT * FROM schedule_logs WHERE schedule_id=$1;

-- SCHEDULE_LOG_INSERT
INSERT INTO schedule_logs (schedule_id,started_at,finished_at,took,result) VALUES ($1,$2,$3,$4,$5);

-- SCHEDULE_MAIL_GET_WITH_SCHEDULE
SELECT * FROM schedule_mails WHERE schedule_id=$1;

-- SCHEDULE_MAIL_INSERT
INSERT INTO schedule_mails (schedule_id,email) VALUES ($1,$2);

-- SCHEDULE_MAIL_UPDATE
UPDATE schedule_mails SET email=$1 WHERE id=$2;

-- SCHEDULE_MAIL_DELETE
UPDATE schedule_mails SET deleted_at=$1 WHERE id=$1;

-- APP_LOG_INSERT
INSERT INTO app_logs (error,log) VALUES ($1,$2);