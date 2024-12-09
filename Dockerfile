FROM alpine:latest AS release-stage
RUN apk add --no-cache tzdata

RUN adduser --disabled-password app

RUN mkdir /db
RUN chown -R app:app /db

RUN addgroup app dialout

WORKDIR /home/app

USER app

COPY --chown=app:app --chmod=0755 /app app

ENTRYPOINT app/wetterserver
