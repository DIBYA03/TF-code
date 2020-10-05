#! /bin/bash
sudo apt-get update -y
sudo apt-get install docker.io awscli -y
sudo usermod -G docker ubuntu && sudo chmod 6662 /var/run/docker.sock
docker pull rancher/rancher
aws s3 cp s3://terraform.sbx-k8s-rancher-wise-us.states/cert.tar /home/ubuntu/ 
cd /home/ubuntu/
tar -xvf cert.tar
mkdir -p /var/log/racher/ && touch /var/log/rancher/auditlog
docker volume create data
docker volume create log
docker volume inspect data
docker run -d --restart=unless-stopped \
  -p 80:80 -p 443:443 \
  -v /home/ubuntu/fullchain.pem:/etc/rancher/ssl/cert.pem \
  -v /home/ubuntu/privkey.pem:/etc/rancher/ssl/key.pem \
  -v /home/ubuntu/cert.pem:/etc/rancher/ssl/cacerts.pem \
  -v /var/log/rancher/auditlog:/var/log/auditlog \
  -v data:/var/lib/rancher \
  --network=host \
  --name=rancher2 \
  rancher/rancher:latest



