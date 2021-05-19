FROM scratch
COPY doxctl /usr/bin/doxctl
ENTRYPOINT ["/usr/bin/doxctl"]
