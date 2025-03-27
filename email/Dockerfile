FROM alpine:latest

RUN mkdir -p /app
RUN mkdir -p /app/templates

COPY bin/emailService /app
COPY email/templates /app/templates

WORKDIR /app

CMD [ "./emailService" ]