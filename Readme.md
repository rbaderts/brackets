
==== Overview ======


Runs a simple double-elimination tournament of any size.  

==== Stack =====

- Go/Postgres backend
- Vue 3/Quasar UI
- Keycloak for OIDC authentication
- Frontend and backend containers 
- A docker-compose orchestration for running the system

===== Runtime ====

4 discreet components:   Front end, Back end, Postgres and Keycloak Server


===== Docker =====

 - Continas a docker-compose orchestration which deploys 4 containers:

    1.  The brackets_b
    2.  The brackets_fe
    3.  Postgres
    4.  Keycloak Auth Server

 - Preparation:

    1.  Build the backend:  ```./build prod```
    2.  Build the backend docker image:   ```docker build . -t brackets_be
    3.  Build the front end docker image
    4.  Build the custom postgres docker image:
         - cd database
         - docker build . -t brackets_postgres
    5.  docker-compose up
