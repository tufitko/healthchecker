FROM ubuntu

ARG TARGETARCH

RUN apt-get update && apt-get install ca-certificates -y && update-ca-certificates

COPY healthchecker /healthchecker

EXPOSE 8080

ENTRYPOINT [ "/healthchecker" ]
