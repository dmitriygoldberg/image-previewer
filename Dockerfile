FROM golang:1.17-buster as base

ENV USER=appuser
ENV UID=10001

RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY ./ /app

RUN make build

FROM alpine:3.15
WORKDIR /previewer
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /app/bin/app /previewer/app

RUN chown -R appuser:appuser /previewer
USER appuser:appuser

CMD ["./app"]