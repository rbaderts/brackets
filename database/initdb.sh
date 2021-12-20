#/bin/sh -x



set POSTGRES_PASSWORD=password
psql --username=keycl0ak --dbname=keycl0ak <<-ENDOFMESSAGE
    create database brackets;
ENDOFMESSAGE
