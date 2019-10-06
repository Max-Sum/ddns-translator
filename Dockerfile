# build stage
FROM golang:alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ADD . /src
RUN cd /src && go build -o http

# final stage
FROM alpine
LABEL maintainer="Max Sum <max@lolyculture.com>"
WORKDIR /app
COPY --from=build-env /src/http /app/
EXPOSE 80
ENTRYPOINT ./http
