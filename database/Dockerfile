FROM alpine:latest

RUN mkdir /app

COPY bin/migrateUp /app
COPY bin/migrateDown /app
COPY bin/databaseService /app

CMD [ "/app/databaseService" ]