FROM            ubuntu:14.04
MAINTAINER      Stephen <steve@steve.com>

ENV DATABASE_HOST="localhost" \
    DATABASE_NAME="slack" \
    DATABASE_USER="sack" \
    DATABASE_PASSWORD=""

RUN apt-get update && apt-get install -y ca-certificates curl
COPY backend-docker backend

EXPOSE 8080
HEALTHCHECK CMD curl -f -s http://localhost:8080/ping || exit 1
ENTRYPOINT ["/backend"]
