FROM scratch
COPY kubepf /usr/bin/kubepf
ENTRYPOINT ["/usr/bin/kubepf"]