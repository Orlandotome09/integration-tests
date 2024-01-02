FROM postgres:12
RUN  rm -rf /tempDB
RUN  mkdir /tempDB
COPY 01-init.sql /docker-entrypoint-initdb.d/