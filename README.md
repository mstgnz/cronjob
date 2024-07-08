# Cronjob
### Timed Trigger

It is difficult to manage cronjob tasks on a server. This project provides an environment where you can define your cronjob tasks. It allows you to manage and track your scheduled tasks in a more flexible way.


## Features
- Cronjob Definitions: Define cronjob tasks for each of your projects.
- URL Information: Specify which URL to send requests to for each cronjob.
- Task Scheduling: Determine when tasks will execute.
- Logging: Detailed logs include task execution times and responses.
- Email Notifications: Optionally send notifications to specified email addresses.
- Message Notifications: Optionally send notifications to specified phone number.


## Considerations
This application utilizes a database locking mechanism to prevent duplicated cronjob tasks when running multiple instances in a Kubernetes (k8s) environment. 

The instance that first adds a record to the "triggered" table will be the one to trigger the cronjob task. Subsequent instances will be prevented from creating a new record due to an existing one, thereby preventing duplicate task executions. Completed tasks are then removed from the "triggered" table.

However, it is strongly recommended that users implement additional control mechanisms on their own systems.


## Contributing
This project is open-source, and contributions are welcome. Feel free to contribute or provide feedback of any kind.


## License
This project is licensed under the Apache License. See the [LICENSE](LICENSE) file for more details.