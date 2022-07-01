FROM gcr.io/distroless/base
COPY --chown=nonroot:nonroot app /app
USER nonroot
ENTRYPOINT [ "/app", "server" ]
