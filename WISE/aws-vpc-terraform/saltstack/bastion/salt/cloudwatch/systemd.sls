cloudwatch.systemd_unit:
  file.managed:
    - name: /etc/systemd/system/cloudwatch-agent.service
    - source: salt://cloudwatch/files/cloudwatch-agent.service
    - template: jinja
  module.run:
    - name: service.systemctl_reload
    - onchanges:
      - file: cloudwatch.systemd_unit
