install_cloudwatch_agent:
  pkg.installed:
    - sources:
      - cloudwatchagent: https://s3.amazonaws.com/amazoncloudwatch-agent/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb
    - check_cmd:
      - "dpkg -l | grep cloudwatch"

cloudwatch_agent_conf:
  file.managed:
    - name: /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json
    - source: salt://cloudwatch/files/amazon-cloudwatch-agent.json
    - user: root

cloudwatch_agent_service:
  file.managed:
    - name: /lib/systemd/system/cloudwatch-agent.service
    - source: salt://cloudwatch/files/cloudwatch-agent.service
    - user: root

cloudwatch-agent:
  service.running:
    - enable: True
    - reload: True
