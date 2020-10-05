route53_record_config:
  file.managed:
    - name: /etc/salt/backup_server.sh
    - source: salt://backup_restore/files/backup_server.sh
    - user: root
    - mode: 700
    - template: jinja
    - context:
      efs_dns_name: {{ salt['environ.get']('EFS_DNS_NAME') }}
      days_to_backup: 15
    - watch_in:
      - cron: backup_server_cron

backup_server_cron:
  cron.present:
    - name: /etc/salt/backup_server.sh
    - minute: 0
    - hour: 05
