{% for user, attrs in salt.pillar.get('users', {}).items() %}

{# create user #}
{% if attrs['state'] == "absent" %}

{# if state is absent, can't accept other attrs #}
{{ user }}:
  user.{{ attrs['state'] }}:
    - force: true

{% else %}

{{ user }}:
  user.{{ attrs['state'] }}:
    - fullname: {{ attrs['fullname'] }}
    - shell: {{ attrs['shell'] }}
    - home: /home/{{ user }}
    - createhome: true
    - uid: {{ attrs['uid'] }}

{% endif %}

{% if 'ssh_keys' in attrs %}

{# add ssh keys for user #}
ssh_keys_{{ user }}:
  ssh_auth.present:
    - user: {{ user }}
    - names:
      {% for key in attrs['ssh_keys'] %}
        - {{ key }}
      {% endfor %}

{% endif %}

{% endfor %}
