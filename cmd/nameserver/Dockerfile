FROM gcr.io/distroless/base
COPY --chown=nonroot:nonroot nameserver /nameserver
USER nonroot
ENTRYPOINT [ "/nameserver", "server" ]
