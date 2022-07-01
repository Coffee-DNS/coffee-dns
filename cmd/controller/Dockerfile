FROM gcr.io/distroless/base
WORKDIR /
COPY --chown=nonroot:nonroot controller .
USER nonroot
ENTRYPOINT [ "/controller", "server" ]
