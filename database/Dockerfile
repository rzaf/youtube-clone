FROM alpine:latest

RUN mkdir -p /app

COPY bin/migrateUp /app
COPY bin/migrateDown /app
COPY bin/databaseService /app

CMD [ "/app/databaseService" ]