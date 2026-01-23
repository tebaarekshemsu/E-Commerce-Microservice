# Kubernetes Deployment (Dev)

Namespace: `e-commerce`

## Images
These manifests reference local images built by Docker Compose (e.g., `project-broker-service:latest`). On Docker Desktop Kubernetes, the image store is shared and they can be used directly. If using another cluster, push images to a registry and update `image:` fields or use `kubectl set image`.

Build images via compose:

```powershell
cd project
docker compose build
```

## Apply manifests
```powershell
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/postgres-init-configmap.yaml
kubectl apply -f k8s/configmap-env.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/mongo.yaml
kubectl apply -f k8s/rabbitmq.yaml
kubectl apply -f k8s/product.yaml
kubectl apply -f k8s/payment.yaml
kubectl apply -f k8s/user.yaml
kubectl apply -f k8s/order.ya
ml
kubectl apply -f k8s/notification.yaml
kubectl apply -f k8s/broker.yaml
kubectl apply -f k8s/ingress.yaml
```

## Access
- Broker via NodePort: `http://localhost:<nodeport>` or via Ingress `http://localhost` if NGINX Ingress is installed.
- RabbitMQ mgmt: `kubectl port-forward svc/rabbitmq -n e-commerce 15672:15672`
- Postgres: `kubectl port-forward svc/postgres -n e-commerce 5432:5432`
- Mongo: `kubectl port-forward svc/mongo -n e-commerce 27017:27017`

## Notes
- Order-service targets container port 8300; service exposes port 80 mapped to 8300.
- Product/payment services expose HTTP (80) and gRPC ports (9001/9002) separately.
- User-service exposes 8000 HTTP and 9003 gRPC.
- Environment values are in `configmap-env.yaml`; sensitive values should be moved to `Secret`.
- Initialize DB schemas via `postgres-init-configmap.yaml` (adjust to your SQL).
