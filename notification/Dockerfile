FROM alpine:latest

RUN mkdir /app
RUN mkdir /app/templates

COPY bin/notificationService /app
COPY email/templates /app/templates

WORKDIR /app

CMD [ "./notificationService" ]