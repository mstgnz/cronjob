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
The application can run as several replicas. Duplicate execution is prevented on two levels:

- **Across instances**, by the `triggered` table. Before a schedule runs, the instance claims a row keyed by `schedule_id`; the primary key makes that claim exclusive, so exactly one replica proceeds and the rest skip the tick. The row is deleted when the run finishes.
- **Within an instance**, by wrapping every job with `SkipIfStillRunning`, so a slow run is never overlapped by the next tick.

Each lock carries a lease (`expires_at`) sized from the schedule's own timeout and retry count. A lock whose lease has passed can be taken over, so an instance that crashes mid-run cannot block its schedules forever. The lease is a backstop, not a timeout: a healthy run releases its lock as soon as it completes. Time comes from the database (`now()`), so clock drift between instances does not affect the lock.

Two consequences worth knowing:

- If a run genuinely outlives its lease, another instance may take the lock over and run the same job again. Executions are therefore at-least-once, not exactly-once. Give slow jobs a realistic `timeout` so the lease is sized correctly.
- The nightly log cleanup is not locked. It is a plain `DELETE` by age, so replicas running it at the same time is harmless.

It is strongly recommended that users implement additional control mechanisms on their own systems.


## Kubernetes Deployment

The application is designed to run in a Kubernetes environment with high availability and scalability in mind. We provide comprehensive Kubernetes configurations and deployment guides in the [k8s](k8s) directory.

Key deployment features include:
- High availability with multiple replicas
- Automatic scaling based on CPU and Memory usage
- Rolling updates for zero-downtime deployments
- Health checks and self-healing capabilities
- SSL/TLS termination with automatic certificate management
- Centralized configuration management

For detailed deployment instructions and configuration options, please refer to the [Kubernetes Deployment Guide](k8s/README.md).


## Contributing
This project is open-source, and contributions are welcome. Feel free to contribute or provide feedback of any kind.


## License
This project is licensed under the Apache License. See the [LICENSE](LICENSE) file for more details.