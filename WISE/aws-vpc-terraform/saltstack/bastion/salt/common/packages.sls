common.packages:
  pkg.installed:
    - pkgs:
      - awscli
      - curl
      - make
      - nfs-common
      - pgcli
      - postgresql-client
      - unzip
      - vim
      - wget
      - zip
      - zsh

install_golang:
  pkgrepo.managed:
    - ppa: longsleep/golang-backports
  pkg.latest:
    - name: golang

install_terraform:
  cmd.run:
    - names:
      - wget https://releases.hashicorp.com/terraform/0.11.14/terraform_0.11.14_linux_amd64.zip
      - unzip -o ./terraform_0.11.14_linux_amd64.zip -d /usr/local/bin/
