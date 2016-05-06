FROM alpine:latest

RUN mkdir -p /opt/k8sbroker /opt/k8sbroker/secrets-templates/ /opt/k8sbroker/service-templates/
ADD broker.elf /opt/k8sbroker/broker.elf
ADD service-templates/ /opt/k8sbroker/service-templates/
ADD secrets-templates/ /opt/k8sbroker/secrets-templates/

EXPOSE 80

WORKDIR /opt/k8sbroker/

ENTRYPOINT [ "/opt/k8sbroker/broker.elf" ]
