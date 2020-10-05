Create a k8s secret for AWS credentials to fetch SSM params. 
Write the AWS access key and secret key in files "key" and "id" files respectively. 
 
kubectl create secret generic aws-credentials-ssm --from-file=id=./id --from-file=key=./key

update this secret path in values.yaml as below. 
```
# envVarsFromSecret:
#  AWS_ACCESS_KEY_ID:
#    secretKeyRef: aws-credentials
#    key: id
#  AWS_SECRET_ACCESS_KEY:
#    secretKeyRef: aws-credentials
#    key: key

```
Install the helm charts for external secret using 

helm install external-secrets kubernetes-external-secrets/charts/kubernetes-external-secrets/ \
-f kubernetes-external-secrets/charts/kubernetes-external-secrets/values.yaml  \
-f kubernetes-external-secrets/values/sbx.yaml \
--namespace vault \
--create-namespace
