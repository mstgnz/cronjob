-- APP_LOG_INSERT
INSERT INTO app_logs (level,message) VALUES ($1,$2);

-- USER_EXISTS_WITH_EMAIL
SELECT count(*) FROM users WHERE email=$1;

-- USER_GET_WITH_ID
SELECT id, fullname, email, is_admin, password FROM users WHERE id=$1;

-- USER_GET_WITH_EMAIL
SELECT id, fullname, email, is_admin, password FROM users WHERE email=$1;

-- USER_INSERT
INSERT INTO users (fullname,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING id,fullname,email,phone;

-- USER_UPDATE_PASS
UPDATE users SET password=$1 WHERE id=$2;

-- USER_LAST_LOGIN
UPDATE users SET last_login=$1 WHERE id=$2;

-- USER_DELETE
UPDATE users SET deleted_at=$1 WHERE id=$2;



-- GROUP_INSERT
INSERT INTO groups (uid, user_id, name) VALUES (CASE WHEN $1 = 0 THEN NULL ELSE $1 END, $2, $3) RETURNING id;

-- GROUP_EXISTS_WITH_USER
SELECT count(*) FROM groups WHERE name=$1 and user_id=$2;