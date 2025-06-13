FROM gcr.io/distroless/static:nonroot
LABEL org.opencontainers.image.source="https://github.com/gabe565/trmnl-nightscout"
COPY trmnl-nightscout /
ENTRYPOINT ["/trmnl-nightscout"]
