base:
  '*':
    - common
    - lvim

  'roles:webserver':
    - match: compound
    - webserver.nginx

  'roles:database':
    - match: compound
    - database.postgresql
