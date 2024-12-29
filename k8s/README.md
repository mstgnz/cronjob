# Kubernetes Deployment Guide

This directory contains Kubernetes configuration files for deploying the Cronjob Manager application. This guide will help you understand and deploy the application in a Kubernetes environment.

## Directory Structure

```
k8s/
├── configmap.yaml    # Application configuration
├── secret.yaml       # Sensitive information
├── deployment.yaml   # Application deployment
├── service.yaml      # Service configuration
├── ingress.yaml      # Ingress configuration
└── hpa.yaml         # Horizontal Pod Autoscaler
```

## Prerequisites

- Kubernetes cluster (v1.20+)
- kubectl CLI tool
- Docker registry access
- Nginx Ingress Controller
- cert-manager (for SSL/TLS)

## Configuration Files

### 1. ConfigMap (configmap.yaml)
Contains non-sensitive configuration data:
- Application settings
- Database connection details
- Redis connection details
- Timezone settings

### 2. Secret (secret.yaml)
Contains sensitive information (base64 encoded):
- Database password
- Redis password
- JWT secret
- SMTP credentials

### 3. Deployment (deployment.yaml)
Defines how the application should be deployed:
- 3 replicas for high availability
- Rolling update strategy
- Resource limits and requests
- Health checks (liveness and readiness probes)
- Environment variables from ConfigMap and Secret

### 4. Service (service.yaml)
Exposes the application within the cluster:
- ClusterIP type service
- Port 80 forwarding to container port 8080

### 5. Ingress (ingress.yaml)
Configures external access:
- SSL/TLS termination
- Domain configuration
- Nginx ingress settings

### 6. HPA (hpa.yaml)
Configures automatic scaling:
- Min 3 replicas
- Max 10 replicas
- CPU and Memory based scaling

## Deployment Steps

1. **Prepare Docker Image**
   ```bash
   # Build the image
   docker build -t your-registry.com/cronjob:latest .
   
   # Push to registry
   docker push your-registry.com/cronjob:latest
   ```

2. **Create Registry Secret**
   ```bash
   kubectl create secret docker-registry regcred \
     --docker-server=your-registry.com \
     --docker-username=your-username \
     --docker-password=your-password
   ```

3. **Update Configuration**
   - Modify `configmap.yaml` with your environment settings
   - Update `secret.yaml` with your encoded secrets:
     ```bash
     # Example: Encoding secrets
     echo -n "your-password" | base64
     ```
   - Update `ingress.yaml` with your domain

4. **Deploy Applications**
   ```bash
   # Apply all configurations
   kubectl apply -f k8s/

   # Or apply individually
   kubectl apply -f k8s/configmap.yaml
   kubectl apply -f k8s/secret.yaml
   kubectl apply -f k8s/deployment.yaml
   kubectl apply -f k8s/service.yaml
   kubectl apply -f k8s/ingress.yaml
   kubectl apply -f k8s/hpa.yaml
   ```

5. **Verify Deployment**
   ```bash
   # Check pods status
   kubectl get pods -l app=cronjob

   # Check service
   kubectl get svc cronjob-service

   # Check ingress
   kubectl get ingress cronjob-ingress

   # Check HPA
   kubectl get hpa cronjob-hpa
   ```

## Monitoring

1. **View Logs**
   ```bash
   # Get pod logs
   kubectl logs -l app=cronjob

   # Follow logs from all pods
   kubectl logs -f -l app=cronjob --all-containers
   ```

2. **Check Resources**
   ```bash
   # Get pod details
   kubectl describe pod -l app=cronjob

   # Check HPA status
   kubectl describe hpa cronjob-hpa
   ```

## Troubleshooting

1. **Pod Issues**
   ```bash
   # Check pod status
   kubectl get pods -l app=cronjob

   # Get pod details
   kubectl describe pod [pod-name]

   # Get pod logs
   kubectl logs [pod-name]
   ```

2. **Service Issues**
   ```bash
   # Check service endpoints
   kubectl get endpoints cronjob-service

   # Test service from another pod
   kubectl run test-pod --rm -it --image=busybox -- wget -qO- http://cronjob-service
   ```

3. **Ingress Issues**
   ```bash
   # Check ingress status
   kubectl describe ingress cronjob-ingress

   # Check ingress controller logs
   kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
   ```

## Scaling

- **Manual Scaling**
  ```bash
  # Scale deployment
  kubectl scale deployment cronjob-deployment --replicas=5
  ```

- **Auto Scaling**
  HPA will automatically scale based on CPU and Memory utilization:
  - Scales up when CPU or Memory > 80%
  - Scales down when CPU and Memory < 80%

## Maintenance

1. **Update Image**
   ```bash
   # Update deployment image
   kubectl set image deployment/cronjob-deployment cronjob=your-registry.com/cronjob:new-tag
   ```

2. **Rolling Restart**
   ```bash
   # Restart all pods
   kubectl rollout restart deployment cronjob-deployment
   ```

3. **Backup Configuration**
   ```bash
   # Export all resources
   kubectl get all -l app=cronjob -o yaml > backup.yaml
   ```

## Security Considerations

1. Always use secrets for sensitive information
2. Regularly rotate credentials
3. Use network policies to restrict traffic
4. Keep the Kubernetes cluster and ingress controller updated
5. Monitor pod security policies
6. Regularly scan container images for vulnerabilities

## Development vs Production

For development environments, consider:
- Reducing replica count
- Lowering resource limits
- Disabling HPA
- Using different ingress settings
- Setting debug level logs 