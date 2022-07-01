FROM gcr.io/distroless/base
COPY --chown=nonroot:nonroot discovery /discovery
USER nonroot
ENTRYPOINT [ "/discovery" ]
