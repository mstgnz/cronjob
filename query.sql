-- APP_LOG_INSERT
INSERT INTO app_logs (level,message) VALUES ($1,$2);

-- USER_EXISTS_WITH_EMAIL
SELECT count(*) FROM users WHERE email=$1;

-- USER_GET_WITH_ID
SELECT id, fullname, email, is_admin, password FROM users WHERE id=$1 and deleted_at isnull;

-- USER_GET_WITH_EMAIL
SELECT id, fullname, email, is_admin, password FROM users WHERE email=$1 and deleted_at isnull;

-- USER_INSERT
INSERT INTO users (fullname,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING id,fullname,email,phone;

-- USER_UPDATE_PASS
UPDATE users SET password=$1, updated_at=$2 WHERE id=$3;

-- USER_LAST_LOGIN
UPDATE users SET last_login=$1 WHERE id=$2;

-- USER_DELETE
UPDATE users SET deleted_at=$1, updated_at=$2 WHERE id=$3;


-- GROUPS
SELECT id,uid,name,active,created_at,updated_at FROM groups where user_id=$1 and deleted_at isnull;

-- GROUPS_WITH_ID
SELECT id,uid,name,active,created_at,updated_at FROM groups where user_id=$1 and id=$2 and deleted_at isnull;

-- GROUP_INSERT
INSERT INTO groups (uid, user_id, name) VALUES (CASE WHEN $1 = 0 THEN NULL ELSE $1 END, $2, $3) RETURNING id,uid,name,active;

-- GROUP_NAME_EXISTS_WITH_USER
SELECT count(*) FROM groups WHERE name=$1 and user_id=$2 and deleted_at isnull;

-- GROUP_ID_EXISTS_WITH_USER
SELECT count(*) FROM groups WHERE id=$1 and user_id=$2 and deleted_at isnull;

-- GROUP_DELETE
UPDATE groups SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;


-- REQUESTS
SELECT * FROM requests where user_id=$1 and deleted_at isnull;

-- REQUESTS_WITH_ID
SELECT * FROM requests where user_id=$1 and id=$2 and deleted_at isnull;

-- REQUEST_INSERT
INSERT INTO requests (user_id,url,method,content,active) VALUES ($1,$2,$3,$4,$5) RETURNING id,user_id,url,method,content,active;

-- REQUEST_URL_EXISTS_WITH_USER
SELECT count(*) FROM requests WHERE url=$1 and user_id=$2 and deleted_at isnull;

-- REQUEST_ID_EXISTS_WITH_USER
SELECT count(*) FROM requests WHERE id=$1 and user_id=$2 and deleted_at isnull;

-- REQUEST_DELETE
UPDATE requests SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;