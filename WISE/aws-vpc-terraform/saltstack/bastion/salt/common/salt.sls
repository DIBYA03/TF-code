# Required for salt stack masterless
salt_daemon.disable:
  service.dead:
    - name: salt-minion
    - enable: false
