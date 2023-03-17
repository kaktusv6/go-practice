#!/bin/sh

DIR_MIGRATIONS=migrations

if [ -z "$1" ]
then
  echo "Задайте имя миграции" >> /dev/stderr
  exit 1
fi

goose -dir $DIR_MIGRATIONS create $1 sql
