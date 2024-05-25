FROM golang:1.22.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO=0 OS=linux ARCH=amd64

RUN CGO_ENABLED=$CGO GOOS=$OS GOARCH=$ARCH go build -o scheduler_app cmd/app/main.go

RUN chmod +x scheduler_app && chmod +x run_sqlite.sh && chmod ugo+rwx -R internal/storage_db

# Run stage
FROM alpine:3.14

RUN apk update && apk upgrade && apk add --no-cache bash=5.1.16-r0 && apk add --no-cache sqlite=3.35.5-r0

ENV USER=docker GROUPNAME=dockergr UID=12345 GID=23456

WORKDIR /home/$USER/app

RUN addgroup --gid "$GID" "$GROUPNAME" \
    && adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$GROUPNAME" \
    --no-create-home \
    --uid "$UID" \
    $USER

USER $USER

COPY --from=builder /app .

ENV TODO_PORT=7540

EXPOSE $TODO_PORT

HEALTHCHECK NONE

CMD [ "./scheduler_app" ]