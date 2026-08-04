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
**Run a single instance.** Overlapping executions of the same schedule are prevented in process: the scheduler wraps every job with `SkipIfStillRunning`, so a run is skipped while the previous one is still going.

There is **no cross-instance locking yet**. The `triggered` table records which schedules are currently executing, but it carries no unique constraint and is not consulted before a job starts, so running several replicas would fire every schedule once per replica. Distributed locking has to be implemented before scaling out.

It is strongly recommended that users implement additional control mechanisms on their own systems.


## Kubernetes Deployment

The application is designed to run in a Kubernetes environment with high availability and scalability in mind. We provide comprehensive Kubernetes configurations and deployment guides in the [k8s](k8s) directory.

Note that the manifests scale horizontally, which the scheduler does not support yet (see Considerations above). Keep the replica count at 1 and the autoscaler disabled until cross-instance locking exists.

Key deployment features include:
- Automatic scaling based on CPU and Memory usage (leave disabled for now)
- Rolling updates for zero-downtime deployments
- Health checks and self-healing capabilities
- SSL/TLS termination with automatic certificate management
- Centralized configuration management

For detailed deployment instructions and configuration options, please refer to the [Kubernetes Deployment Guide](k8s/README.md).


## Contributing
This project is open-source, and contributions are welcome. Feel free to contribute or provide feedback of any kind.


## License
This project is licensed under the Apache License. See the [LICENSE](LICENSE) file for more details.