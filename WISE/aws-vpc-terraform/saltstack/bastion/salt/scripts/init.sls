generate_rds_creds_script:
  file.managed:
    - name: /usr/local/bin/generate_rds_creds
    - source: salt://scripts/files/generate_rds_creds
    - user: root
    - mode: 0775
    - template: jinja

generate_pgcli_connect_script:
  file.managed:
    - name: /usr/local/bin/pgcli_connect
    - source: salt://scripts/files/pgcli_connect
    - user: root
    - mode: 0775
    - template: jinja
