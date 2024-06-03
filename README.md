# cronjob

Defining cronjob tasks in the crontab on a server and tracking and managing these tasks can be challenging. In this project, you can input cronjob definitions for all your projects and the corresponding URL information. This allows requests to be sent to the specified URLs on your behalf at the designated times.

The time taken by triggered tasks, their responses, are logged, and notifications can be optionally sent to the email addresses you specify.


## atention
This application utilizes a database locking mechanism to prevent duplicate cronjob tasks when running multiple instances in a Kubernetes (k8s) environment.

The instance that first adds a record to the triggered table will be the instance that triggers the cronjob task. Subsequent instances will not be able to create a new record due to an existing record, preventing them from executing the task. The completed task will be unregistered from the triggered table.

Nevertheless, it is strongly advised for users of this project to establish an additional control mechanism on their own systems.