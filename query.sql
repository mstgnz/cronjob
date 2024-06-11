-- APP_LOG_INSERT
INSERT INTO app_logs (level,message) VALUES ($1,$2);

-- USER_EXISTS_WITH_EMAIL
SELECT count(*) FROM users WHERE email=$1;

-- USER_GET_WITH_ID
SELECT id, fullname, email, is_admin, password FROM users WHERE id=$1 AND deleted_at isnull;

-- USER_GET_WITH_EMAIL
SELECT id, fullname, email, is_admin, password FROM users WHERE email=$1 AND deleted_at isnull;

-- USER_INSERT
INSERT INTO users (fullname,email,password,phone) VALUES ($1,$2,$3,$4) RETURNING id,fullname,email,phone;

-- USER_UPDATE_PASS
UPDATE users SET password=$1, updated_at=$2 WHERE id=$3;

-- USER_LAST_LOGIN
UPDATE users SET last_login=$1 WHERE id=$2;

-- USER_DELETE
UPDATE users SET deleted_at=$1, updated_at=$2 WHERE id=$3;


-- GROUPS
SELECT id,uid,name,active,created_at,updated_at FROM groups WHERE user_id=$1 AND deleted_at isnull;

-- GROUPS_WITH_ID
SELECT id,uid,name,active,created_at,updated_at FROM groups WHERE user_id=$1 AND id=$2 AND deleted_at isnull;

-- GROUP_INSERT
INSERT INTO groups (uid, user_id, name) VALUES (CASE WHEN $1 = 0 THEN NULL ELSE $1 END, $2, $3) RETURNING id,uid,name,active;

-- GROUP_NAME_EXISTS_WITH_USER
SELECT count(*) FROM groups WHERE name=$1 AND user_id=$2 AND deleted_at isnull;

-- GROUP_ID_EXISTS_WITH_USER
SELECT count(*) FROM groups WHERE id=$1 AND user_id=$2 AND deleted_at isnull;

-- GROUP_DELETE
UPDATE groups SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;


-- REQUESTS
SELECT * FROM requests WHERE user_id=$1 AND deleted_at isnull;

-- REQUESTS_WITH_ID
SELECT * FROM requests WHERE user_id=$1 AND id=$2 AND deleted_at isnull;

-- REQUEST_INSERT
INSERT INTO requests (user_id,url,method,content,active) VALUES ($1,$2,$3,$4,$5) RETURNING id,user_id,url,method,content,active;

-- REQUEST_URL_EXISTS_WITH_USER
SELECT count(*) FROM requests WHERE url=$1 AND user_id=$2 AND deleted_at isnull;

-- REQUEST_ID_EXISTS_WITH_USER
SELECT count(*) FROM requests WHERE id=$1 AND user_id=$2 AND deleted_at isnull;

-- REQUEST_DELETE
UPDATE requests SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;


-- REQUEST_HEADERS
SELECT rh.* FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE r.user_id=$1 AND rh.deleted_at isnull;

-- REQUEST_HEADERS_WITH_ID
SELECT rh.* FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE r.user_id=$1 AND rh.id=$2 AND rh.deleted_at isnull;

-- REQUEST_HEADER_INSERT
INSERT INTO request_headers (request_id,key,value,active) VALUES ($1,$2,$3,$4) RETURNING id,request_id,key,value,active;

-- REQUEST_HEADER_EXISTS_WITH_USER
SELECT count(*) FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE rh.key=$1 AND r.user_id=$2 AND rh.deleted_at isnull;

-- REQUEST_HEADER_ID_EXISTS_WITH_USER
SELECT count(*) FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE rh.id=$1 AND r.user_id=$2 AND rh.deleted_at isnull;

-- REQUEST_HEADER_DELETE
UPDATE request_headers SET deleted_at=$1, updated_at=$2 FROM requests 
WHERE requests.id=request_headers.request_id AND request_headers.id=$3 AND requests.user_id=$4;


-- SCHEDULES
SELECT * FROM schedules WHERE user_id=$1 AND deleted_at isnull;

-- SCHEDULES_WITH_ID
SELECT * FROM schedules WHERE user_id=$1 AND id=$2 AND deleted_at isnull;

-- SCHEDULES_INSERT
INSERT INTO schedules (user_id,group_id,request_id,timing,timeout,retries,active) 
VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id,user_id,group_id,request_id,timing,timeout,retries,running,active;

-- SCHEDULES_ID_EXISTS_WITH_USER
SELECT count(*) FROM schedules WHERE id=$1 AND user_id=$2 AND deleted_at isnull;

-- SCHEDULES_DELETE
UPDATE schedules SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;

-- SCHEDULE_LOGS
SELECT * FROM schedule_logs sl JOIN schedules s ON s.id=sl.schedule_id WHERE sl.schedule_id=$1 AND s.user_id=$2;


-- WEBHOOKS
SELECT w.* FROM webhooks w JOIN schedules s ON s.id=w.schedule_id WHERE s.user_id=$1 AND w.deleted_at isnull;

-- WEBHOOKS_WITH_ID
SELECT w.* FROM webhooks w JOIN schedules s ON s.id=w.schedule_id WHERE w.id=$1 AND s.user_id=$2 AND w.deleted_at isnull;

-- WEBHOOK_INSERT
INSERT INTO webhooks (schedule_id,request_id, active) VALUES ($1,$2,$3) RETURNING id,schedule_id,request_id,active;

-- WEBHOOK_ID_EXISTS_WITH_USER
SELECT count(*) FROM webhooks w JOIN schedules s ON s.id=w.schedule_id WHERE w.id=$1 AND s.user_id=$2 AND w.deleted_at isnull;

-- WEBHOOK_UNIQ_EXISTS_WITH_USER
SELECT count(*) FROM webhooks w JOIN schedules s ON s.id=w.schedule_id WHERE w.schedule_id=$1 AND w.request_id=$2 AND s.user_id=$3 AND w.deleted_at isnull;

-- WEBHOOK_DELETE
UPDATE webhooks SET deleted_at=$1, updated_at=$2 FROM schedules 
WHERE schedules.id=webhooks.schedule_id AND webhooks.id=$3 AND schedules.user_id=$4;


-- NOTIFICATIONS
SELECT n.* FROM notifications n JOIN schedules s ON s.id=n.schedule_id WHERE s.user_id=$1 AND n.deleted_at isnull;

-- NOTIFICATIONS_WITH_ID
SELECT n.* FROM notifications n JOIN schedules s ON s.id=n.schedule_id WHERE s.user_id=$1 AND n.id=$2 AND n.deleted_at isnull;

-- NOTIFICATION_INSERT
INSERT INTO notifications 
(schedule_id,is_sms,is_mail,title,content,active) 
VALUES 
($1,$2,$3,$4,$5,$6) RETURNING id,schedule_id,is_sms,is_mail,title,content,active;

-- NOTIFICATION_TITLE_EXISTS_WITH_USER_AND_SCHEDULE
SELECT count(*) FROM notifications n JOIN schedules s ON s.id=n.schedule_id WHERE s.user_id=$1 AND n.title=$2 AND n.schedule_id=$3 AND n.deleted_at isnull;

-- NOTIFICATION_ID_EXISTS_WITH_USER
SELECT count(*) FROM notifications n JOIN schedules s ON s.id=n.schedule_id WHERE n.id=$1 AND s.user_id=$2 AND n.deleted_at isnull;

-- NOTIFICATION_DELETE
UPDATE notifications SET deleted_at=$1, updated_at=$2 FROM schedules 
WHERE schedules.id=notifications.schedule_id AND notifications.id=$3 AND schedules.user_id=$4;


-- NOTIFICATION_EMAILS
SELECT ne.* FROM notify_email ne 
JOIN notifications n ON n.id=ne.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE s.user_id=$1 AND ne.deleted_at isnull;

-- NOTIFICATION_EMAILS_WITH_ID
SELECT ne.* FROM notify_email ne 
JOIN notifications n ON n.id=ne.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE s.user_id=$1 AND ne.id=$2 AND ne.deleted_at isnull;

-- NOTIFICATION_EMAIL_INSERT
INSERT INTO notify_email (notification_id,email,active) VALUES ($1,$2,$3) RETURNING id,notification_id,email,active;

-- NOTIFICATION_EMAIL_EXISTS_WITH_USER
SELECT count(*) FROM notify_email ne 
JOIN notifications n ON n.id=ne.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE s.user_id=$1 AND ne.email=$2 AND n.id=$3 AND ne.deleted_at isnull;

-- NOTIFICATION_EMAIL_ID_EXISTS_WITH_USER
SELECT count(*) FROM notify_email ne 
JOIN notifications n ON n.id=ne.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE ne.id=$1 AND s.user_id=$2 AND ne.deleted_at isnull;

-- NOTIFICATION_EMAIL_DELETE
UPDATE notify_email SET deleted_at=$1, updated_at=$2 FROM notifications 
JOIN schedules ON schedules.id=notifications.schedule_id 
WHERE notifications.id=notify_email.notification_id AND notify_email.id=$3 AND schedules.user_id=$4;


-- NOTIFICATION_SMS
SELECT ns.* FROM notify_sms ns 
JOIN notifications n ON n.id=ns.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE s.user_id=$1 AND ns.deleted_at isnull;

-- NOTIFICATION_SMS_WITH_ID
SELECT ns.* FROM notify_sms ns 
JOIN notifications n ON n.id=ns.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE s.user_id=$1 AND ns.id=$2 AND ns.deleted_at isnull;

-- NOTIFICATION_SMS_INSERT
INSERT INTO notify_sms (notification_id,phone,active) VALUES ($1,$2,$3) RETURNING id,notification_id,phone,active;

-- NOTIFICATION_SMS_PHONE_EXISTS_WITH_USER
SELECT count(*) FROM notify_sms ns 
JOIN notifications n ON n.id=ns.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE s.user_id=$1 AND ns.phone=$2 AND n.id=$3 AND ns.deleted_at isnull;

-- NOTIFICATION_SMS_ID_EXISTS_WITH_USER
SELECT count(*) FROM notify_sms ns 
JOIN notifications n ON n.id=ns.notification_id 
JOIN schedules s ON s.id=n.schedule_id 
WHERE ns.id=$1 AND s.user_id=$2 AND ns.deleted_at isnull;

-- NOTIFICATION_SMS_DELETE
UPDATE notify_sms SET deleted_at=$1, updated_at=$2 FROM notifications 
JOIN schedules ON schedules.id=notifications.schedule_id 
WHERE notifications.id=notify_sms.notification_id AND notify_sms.id=$3 AND schedules.user_id=$4;