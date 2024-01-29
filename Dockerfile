FROM scratch
COPY go-tables /
ENTRYPOINT ["/go-tables"]
