# cronjob

Defining cronjob tasks in the crontab on a server and tracking and managing these tasks can be challenging. In this project, you can input cronjob definitions for all your projects and the corresponding URL information. This allows requests to be sent to the specified URLs on your behalf at the designated times.

The time taken by triggered tasks, their responses, are logged, and notifications can be optionally sent to the email addresses you specify.

This project is user-centric. Each user has their own projects, and therefore, the base URL definition is determined based on user information. When a user wants to add a scheduled task, they only need to enter the path after the domain part of the URL. The system combines this with the base URL to send requests to the resulting URL.