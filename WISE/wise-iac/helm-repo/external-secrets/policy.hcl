path "secret/*" {
  capabilities = ["create", "read", "update", "delete", "list"]
}
path "sys/policies/acl/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
path "auth/*"
{
  capabilities = ["create", "read", "update", "delete", "list", "sudo"]
}
