
FROM golang:1.16-alpine as builder
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
WORKDIR /go/src/app
COPY . .
RUN apk add --no-cache make
RUN CGO_ENABLED=0 make build


FROM redhat/ubi8-minimal
ENV LANG=C.UTF-8 LC_ALL=C.UTF-8
COPY --from=builder /go/src/app/harbor-cli /usr/bin/harbor-cli
RUN \
  chmod +x /usr/bin/harbor-cli

ENTRYPOINT [ "/usr/bin/harbor-cli" ]