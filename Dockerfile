FROM ubuntu

ARG TARGETARCH

COPY healthchecker /healthchecker

EXPOSE 8080

ENTRYPOINT [ "/healthchecker" ]
