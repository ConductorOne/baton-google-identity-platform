FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-google-identity-platform"]
COPY baton-google-identity-platform /