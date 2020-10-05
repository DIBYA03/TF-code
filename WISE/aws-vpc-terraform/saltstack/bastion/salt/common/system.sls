system:
  network.system:
    - enabled: True # only neccessary as a bypass for https://github.com/saltstack/salt/issues/6922
    - hostname: '{{ salt['environ.get']('BASTION_HOSTNAME') }}'
    - apply_hostname: True
    - retain_settings: True
