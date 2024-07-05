-- APP_LOG_COUNT
SELECT count(*) FROM app_logs;

-- APP_LOG_PAGINATE
SELECT * FROM app_logs ORDER BY id DESC offset $1 LIMIT $2;

-- APP_LOG_INSERT
INSERT INTO app_logs (level,message) VALUES ($1,$2);

-- TRIGGERED_INSERT
INSERT INTO triggered (schedule_id) VALUES ($1);

-- TRIGGERED_DELETE
DELETE FROM triggered WHERE schedule_id=$1;


-- USERS_COUNT
SELECT count(*) FROM users;

-- USERS_PAGINATE
select * from users where fullname ilike $1 or email ilike $1 or phone ilike $1 order by id desc offset $2 limit $3;

-- USER_EXISTS_WITH_ID
SELECT count(*) FROM users WHERE id=$1;

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
UPDATE users SET active=$1, deleted_at=$2, updated_at=$3 WHERE id=$4;


-- GROUPS_COUNT
SELECT count(*) FROM groups WHERE user_id=$1 AND deleted_at isnull;

-- GROUPS_PAGINATE
SELECT g.*, p.name as parent, u.fullname FROM groups g 
JOIN users u ON u.id=g.user_id 
LEFT JOIN groups p ON p.id=g.uid 
WHERE (g.name ilike $1 OR p.name ilike $1 OR u.fullname ilike $1) AND g.user_id=$2 AND g.deleted_at isnull AND p.deleted_at isnull 
ORDER BY g.id DESC offset $3 LIMIT $4;

-- GROUPS
SELECT id,uid,name,active,created_at,updated_at FROM groups WHERE user_id=$1 AND deleted_at isnull;

-- GROUPS_WITH_ID
SELECT id,uid,name,active,created_at,updated_at FROM groups WHERE user_id=$1 AND id=$2 AND deleted_at isnull;

-- GROUP_INSERT
INSERT INTO groups (uid, user_id, name, active) VALUES (CASE WHEN $1 = 0 THEN NULL ELSE $1 END, $2, $3, $4) RETURNING id,uid,name,active;

-- GROUP_NAME_EXISTS_WITH_USER
SELECT count(*) FROM groups WHERE name=$1 AND user_id=$2 AND deleted_at isnull;

-- GROUP_ID_EXISTS_WITH_USER
SELECT count(*) FROM groups WHERE id=$1 AND user_id=$2 AND deleted_at isnull;

-- GROUP_DELETE
UPDATE groups SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;


-- REQUESTS_COUNT
SELECT count(*) FROM requests WHERE user_id=$1 AND deleted_at isnull;

-- REQUESTS_PAGINATE
SELECT r.*, u.fullname FROM requests r 
JOIN users u ON u.id=r.user_id 
WHERE (r.url ilike $1 OR r.method ilike $1 OR r.content::text ilike $1 OR u.fullname ilike $1) AND r.user_id=$2 AND r.deleted_at isnull 
ORDER BY r.id DESC offset $3 LIMIT $4;

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


-- REQUEST_HEADERS_COUNT
SELECT count(rh.*) FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE r.user_id=$1 AND rh.deleted_at isnull;

-- REQUEST_HEADERS_PAGINATE
SELECT rh.*, r.url FROM request_headers rh  
JOIN requests r ON r.id=rh.request_id 
WHERE (rh.key ilike $1 OR rh.value ilike $1 OR r.url ilike $1) AND r.user_id=$2 AND r.deleted_at isnull AND rh.deleted_at isnull 
ORDER BY rh.id DESC offset $3 LIMIT $4;

-- REQUEST_HEADERS
SELECT rh.* FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE r.user_id=$1 AND rh.deleted_at isnull AND r.deleted_at isnull;

-- REQUEST_HEADERS_WITH_ID
SELECT rh.* FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE r.user_id=$1 AND rh.id=$2 AND rh.deleted_at isnull AND r.deleted_at isnull;

-- REQUEST_HEADER_INSERT
INSERT INTO request_headers (request_id,key,value,active) VALUES ($1,$2,$3,$4) RETURNING id,request_id,key,value,active;

-- REQUEST_HEADER_EXISTS_WITH_USER
SELECT count(*) FROM request_headers rh JOIN requests r ON r.id=rh.request_id 
WHERE rh.key=$1 AND r.user_id=$2 AND r.id=$3 AND rh.deleted_at isnull;

-- REQUEST_HEADER_ID_EXISTS_WITH_USER
SELECT count(*) FROM request_headers rh JOIN requests r ON r.id=rh.request_id WHERE rh.id=$1 AND r.user_id=$2 AND rh.deleted_at isnull AND r.deleted_at isnull;

-- REQUEST_HEADER_DELETE
UPDATE request_headers SET deleted_at=$1, updated_at=$2 FROM requests 
WHERE requests.id=request_headers.request_id AND request_headers.id=$3 AND requests.user_id=$4;


-- SCHEDULES_COUNT
SELECT count(*) FROM schedules WHERE user_id=$1 AND deleted_at isnull;

-- SCHEDULES_PAGINATE
SELECT s.*, g.name, r.url, n.title FROM schedules s
JOIN groups g ON g.id=s.group_id
JOIN requests r ON r.id=s.request_id
JOIN notifications n ON n.id=s.notification_id
WHERE s.user_id=$1 AND s.deleted_at isnull AND (s.timing ilike $2 OR g.name ilike $2 OR r.url ilike $2 OR n.title ilike $2)
ORDER BY s.id DESC offset $3 LIMIT $4;

-- SCHEDULES
SELECT s.*, g.name, r.url, n.title FROM schedules s 
JOIN groups g ON g.id=s.group_id 
JOIN requests r ON r.id=s.request_id 
JOIN notifications n ON n.id=s.notification_id 
WHERE s.user_id=$1 AND s.deleted_at isnull;

-- SCHEDULES_WITH_ID
SELECT * FROM schedules WHERE user_id=$1 AND id=$2 AND deleted_at isnull;

-- SCHEDULES_INSERT
INSERT INTO schedules (user_id,group_id,request_id,notification_id,timing,timeout,retries,active) 
VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id,user_id,group_id,request_id,notification_id,timing,timeout,retries,running,active;

-- SCHEDULES_ID_EXISTS_WITH_USER
SELECT count(*) FROM schedules WHERE id=$1 AND user_id=$2 AND deleted_at isnull;

-- SCHEDULES_TIMING_EXISTS_WITH_USER
SELECT count(*) FROM schedules WHERE user_id=$1 AND request_id=$2 AND timing=$3 AND deleted_at isnull;

-- SCHEDULES_DELETE
UPDATE schedules SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;


-- SCHEDULE_LOGS_COUNT
SELECT count(sl.*) FROM schedule_logs sl JOIN schedules s ON s.id=sl.schedule_id WHERE s.user_id=$1 AND s.deleted_at isnull;

-- SCHEDULE_LOGS_PAGINATE
SELECT sl.*, s.timing FROM schedule_logs sl
JOIN schedules s ON s.id=sl.schedule_id
WHERE s.user_id=$1 AND s.deleted_at isnull AND s.timing ilike $2
ORDER BY sl.id DESC offset $3 LIMIT $4;

-- SCHEDULE_LOGS
SELECT * FROM schedule_logs sl JOIN schedules s ON s.id=sl.schedule_id WHERE sl.schedule_id=$1 AND s.user_id=$2;


-- WEBHOOKS_COUNT
SELECT count(w.*) FROM webhooks w JOIN schedules s ON s.id=w.schedule_id WHERE s.user_id=$1 AND w.deleted_at isnull;

-- WEBHOOKS_PAGINATE
SELECT w.*, s.timing, r.url FROM webhooks w
JOIN schedules s ON s.id=w.schedule_id
JOIN requests r ON r.id=w.request_id
WHERE s.user_id=$1 AND w.deleted_at isnull AND (s.timing ilike $2 OR r.url ilike $2)
ORDER BY w.id DESC offset $3 LIMIT $4;

-- WEBHOOKS
SELECT w.*, s.timing, r.url FROM webhooks w
JOIN schedules s ON s.id=w.schedule_id
JOIN requests r ON r.id=w.request_id
WHERE s.user_id=$1 AND w.deleted_at isnull;

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


-- NOTIFICATIONS_COUNT
SELECT count(*) FROM notifications WHERE user_id=$1 AND deleted_at isnull;

-- NOTIFICATIONS_PAGINATE
SELECT n.*, u.fullname FROM notifications n 
JOIN users u ON u.id=n.user_id 
WHERE (n.title ilike $1 OR n.content ilike $1 OR u.fullname ilike $1) AND n.user_id=$2 AND n.deleted_at isnull 
ORDER BY n.id DESC offset $3 LIMIT $4;

-- NOTIFICATIONS
SELECT * FROM notifications WHERE user_id=$1 AND deleted_at isnull;

-- NOTIFICATIONS_WITH_ID
SELECT * FROM notifications WHERE user_id=$1 AND id=$2 AND deleted_at isnull;

-- NOTIFICATION_INSERT
INSERT INTO notifications (user_id,title,content,is_mail,is_message,active) VALUES ($1,$2,$3,$4,$5,$6) RETURNING id,user_id,title,content,is_mail,is_message,active;

-- NOTIFICATION_TITLE_EXISTS_WITH_USER
SELECT count(*) FROM notifications WHERE user_id=$1 AND title=$2 AND deleted_at isnull;

-- NOTIFICATION_ID_EXISTS_WITH_USER
SELECT count(*) FROM notifications WHERE id=$1 AND user_id=$2 AND deleted_at isnull;

-- NOTIFICATION_DELETE
UPDATE notifications SET deleted_at=$1, updated_at=$2 WHERE id=$3 AND user_id=$4;


-- NOTIFICATION_EMAILS_COUNT
SELECT count(ne.*) FROM notify_emails ne JOIN notifications n ON n.id=ne.notification_id WHERE n.user_id=$1 AND ne.deleted_at isnull;

-- NOTIFICATION_EMAILS_PAGINATE
SELECT ne.*, n.title FROM notify_emails ne 
JOIN notifications n ON n.id=ne.notification_id 
WHERE (n.title ilike $1) AND n.user_id=$2 AND ne.deleted_at isnull AND n.deleted_at isnull 
ORDER BY ne.id DESC offset $3 LIMIT $4;

-- NOTIFICATION_EMAILS
SELECT ne.*, n.title FROM notify_emails ne JOIN notifications n ON n.id=ne.notification_id WHERE n.user_id=$1 AND ne.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_EMAILS_WITH_ID
SELECT ne.* FROM notify_emails ne JOIN notifications n ON n.id=ne.notification_id WHERE n.user_id=$1 AND ne.id=$2 AND ne.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_EMAIL_INSERT
INSERT INTO notify_emails (notification_id,email,active) VALUES ($1,$2,$3) RETURNING id,notification_id,email,active;

-- NOTIFICATION_EMAIL_EXISTS_WITH_USER
SELECT count(*) FROM notify_emails ne JOIN notifications n ON n.id=ne.notification_id 
WHERE n.user_id=$1 AND ne.email=$2 AND n.id=$3 AND ne.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_EMAIL_ID_EXISTS_WITH_USER
SELECT count(*) FROM notify_emails ne JOIN notifications n ON n.id=ne.notification_id 
WHERE ne.id=$1 AND n.user_id=$2 AND ne.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_EMAIL_DELETE
UPDATE notify_emails SET deleted_at=$1, updated_at=$2 FROM notifications 
WHERE notifications.id=notify_emails.notification_id AND notify_emails.id=$3 AND notifications.user_id=$4;


-- NOTIFICATION_MESSAGES_COUNT
SELECT count(nm.*) FROM notify_messages nm JOIN notifications n ON n.id=nm.notification_id WHERE n.user_id=$1 AND nm.deleted_at isnull;

-- NOTIFICATION_MESSAGES_PAGINATE
SELECT nm.*, n.title FROM notify_messages nm 
JOIN notifications n ON n.id=nm.notification_id 
WHERE (n.title ilike $1) AND n.user_id=$2 AND nm.deleted_at isnull AND n.deleted_at isnull 
ORDER BY nm.id DESC offset $3 LIMIT $4;

-- NOTIFICATION_MESSAGES
SELECT nm.*, n.title FROM notify_messages nm JOIN notifications n ON n.id=nm.notification_id WHERE n.user_id=$1 AND nm.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_MESSAGE_WITH_ID
SELECT nm.* FROM notify_messages nm JOIN notifications n ON n.id=nm.notification_id WHERE n.user_id=$1 AND nm.id=$2 AND ns.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_MESSAGE_INSERT
INSERT INTO notify_messages (notification_id,phone,active) VALUES ($1,$2,$3) RETURNING id,notification_id,phone,active;

-- NOTIFICATION_MESSAGE_PHONE_EXISTS_WITH_USER
SELECT count(*) FROM notify_messages nm JOIN notifications n ON n.id=nm.notification_id 
WHERE n.user_id=$1 AND nm.phone=$2 AND n.id=$3 AND nm.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_MESSAGE_ID_EXISTS_WITH_USER
SELECT count(*) FROM notify_messages nm JOIN notifications n ON n.id=nm.notification_id 
WHERE nm.id=$1 AND n.user_id=$2 AND nm.deleted_at isnull AND n.deleted_at isnull;

-- NOTIFICATION_MESSAGE_DELETE
UPDATE notify_messages SET deleted_at=$1, updated_at=$2 FROM notifications 
WHERE notifications.id=notify_messages.notification_id AND notify_messages.id=$3 AND notifications.user_id=$4;


-- SCHEDULE_MAPS
WITH schedule_lists AS (
    SELECT
        s.*,
        json_build_object(
            'id', u.id,
            'fullname', u.fullname,
            'email', u.email,
            'phone', u.phone
        ) as user,
        json_build_object(
            'id', g.id,
            'uid', g.uid,
            'name', g.name,
            'active', g.active,
            'parent', (
                SELECT json_build_object(
                    'id', p.id,
                    'name', p.name 
                )
                FROM groups p 
                WHERE p.id = g.uid
            )
        ) as group,
        json_build_object(
            'id', r.id,
            'user_id', r.user_id,
            'url', r.url,
            'method', r.method,
            'content', r.content,
            'active', r.active,
            'headers', (
                SELECT json_agg(
                    json_build_object(
                        'id', rh.id,
                        'key', rh.key,
                        'value', rh.value,
                        'active', rh.active
                    )
                )
                FROM request_headers rh 
                WHERE rh.request_id = r.id
            )
        ) as request,
        json_build_object(
            'id', n.id,
            'user_id', n.user_id,
            'title', n.title,
            'content', n.content,
            'is_mail', n.is_mail,
            'is_message', n.is_mail,
            'active', n.active,
            'emails', (
                SELECT json_agg(
                    json_build_object(
                        'id', ne.id,
                        'email', ne.email,
                        'active', ne.active
                    )
                )
                FROM notify_emails ne 
                WHERE ne.notification_id = n.id
            ),
            'messages', (
                SELECT json_agg(
                    json_build_object(
                        'id', nm.id,
                        'phone', nm.phone,
                        'active', nm.active
                    )
                )
                FROM notify_messages nm 
                WHERE nm.notification_id = n.id
            )
        ) as notification,
        json_agg(
            json_build_object(
                'id', w.id,
                'schedule_id', w.schedule_id,
                'request_id', w.request_id,
                'active', w.active,
                'requests', (
                    SELECT json_agg(
                        json_build_object(
                            'id', r.id,
                            'url', r.url,
                            'content', r.content,
                            'active', r.active,
                            'headers', (
                                SELECT json_agg(
                                    json_build_object(
                                        'id', rh.id,
                                        'key', rh.key,
                                        'value', rh.value,
                                        'active', rh.active
                                    )
                                )
                                FROM request_headers rh 
                                WHERE rh.request_id = r.id
                            )
                        )
                    )
                    FROM requests r 
                    WHERE r.id = w.request_id
                )
            )
        ) as webhooks
    FROM schedules as s
    JOIN users AS u ON u.id = s.user_id
    JOIN groups AS g ON g.id = s.group_id
    JOIN requests AS r ON r.id = s.request_id
    JOIN notifications n on n.id=s.notification_id
    LEFT JOIN webhooks w on w.schedule_id=s.id
    GROUP BY s.id, u.id, g.id, r.id, n.id
)
SELECT * FROM schedule_lists
WHERE user_id=$1 AND deleted_at isnull 
AND (timing ilike $2 OR "group"->>'name' ilike $2 OR request->>'url' ilike $2 OR notification->>'title' ilike $2)
ORDER BY id DESC offset $3 LIMIT $4;