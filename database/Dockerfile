FROM postgres

COPY initdb.sh /docker-entrypoint-initdb.d/
COPY postgresql.conf /etc/postgresql/
USER root
RUN  chmod 666 /etc/postgresql/postgresql.conf



CMD ["postgres", "-c", "config_file=/etc/postgresql/postgresql.conf"]

