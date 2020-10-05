#!/usr/local/bin/fish

set -gx CONTAINER_LISTEN_PORT 3001
set -gx GRPC_SERVICE_PORT 3001
set -gx HTTPS_LISTEN_PORT 8443
set -gx HTTP_HEALTH_PORT 8445

set -gx ENV_NAME staging
set -gx API_ENV staging
set -gx AWS_S3_BUCKET terraform.wise-us.states
