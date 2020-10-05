Configure the kube config for the k8s cluster in your machine with sufficient privileges. 

switch to the chart directory of respective tool, and run `helm template <toolname> ./ ` to render the template on local.

Confirm the changes, and run `helm install <toolname> <--set to override any values> ./ ` to apply the changes, (use helm upgrade command to make changes in already deployed chart)

Note: Helm version: Version:"v3.2.1" +
