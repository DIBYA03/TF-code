Execute 
```
sh deploy.yaml

```

to deploy the certmanager. 

Notes: update the mail id in cluster-issuer-prod.yaml and cluster-issuer-staging.yaml to a valid mail id. 
Expected helm version v3+ , and kubectl installed with kubeconfig configured to your expected cluster context

To avoid helm install the remote chart, replace the command with  

```
 kubectl apply -f helm-manifest.yaml
```


To create an ingress, add below annotations to make the cert-manager automatically create certificates for you and the tls secrets as well

```
metadata:
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-staging ## or prod based on your env
    kubernetes.io/ingress.class: nginx-external  ## use public ingress class
    kubernetes.io/tls-acme: "true"
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/whitelist-source-range: 0.0.0.0/0
  name: ingress-name 
  namespace: your-namespace

```


sample ingress file: 
Here, the tls secret event.sbx.wise.us will be automatically created. 


```
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: service-event
  labels:
    helm.sh/chart: service-event-0.1.0
    app.kubernetes.io/name: service-event
    app.kubernetes.io/instance: service-event
    app.kubernetes.io/version: "1.16.0"
    app.kubernetes.io/managed-by: Helm
  annotations:
    kubernetes.io/ingress.class: nginx-external
    kubernetes.io/tls-acme: "true"
    cert-manager.io/cluster-issuer: letsencrypt-staging
    nginx.ingress.kubernetes.io/rewrite-target: /
    nginx.ingress.kubernetes.io/whitelist-source-range: 0.0.0.0/0
spec:
  tls:
    - hosts:
        - "event.sbx.wise.us"
      secretName: event.sbx.wise.us
  rules:
    - host: "event.sbx.wise.us"
      http:
        paths:
          - path: /
            backend:
              serviceName: service-event
              servicePort: 80

```
