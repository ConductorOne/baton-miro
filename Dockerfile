FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-miro"]
COPY baton-miro /