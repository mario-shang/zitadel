FROM --platform=$TARGETPLATFORM debian:latest as artifact
ENV ZITADEL_ARGS=
ARG TARGETPLATFORM

RUN apt-get update && apt-get install -y \
    build-essential \
    ca-certificates \
    make \
    wget \
    curl \
    gnupg

ENV GOLANG_VERSION 1.22.3
RUN wget https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GOLANG_VERSION}.linux-amd64.tar.gz && \
    rm go${GOLANG_VERSION}.linux-amd64.tar.gz
ENV GOPATH="/usr/local/go"
ENV PATH="${GOPATH}/bin:${PATH}"

ENV NODE_VERSION 18.x
RUN curl -fsSL https://deb.nodesource.com/setup_${NODE_VERSION} | bash - && \
    apt-get install -y nodejs

RUN npm install -g sass yarn

COPY . /src
WORKDIR /src

RUN make compile

COPY build/entrypoint.sh /app/entrypoint.sh
RUN cp zitadel /app/

RUN useradd -s "" --home / zitadel && \
    chown zitadel /app/zitadel && \
    chmod +x /app/zitadel && \
    chown zitadel /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

WORKDIR /app
ENV PATH="/app:${PATH}"

USER zitadel
ENTRYPOINT ["/app/entrypoint.sh"]

FROM --platform=$TARGETPLATFORM scratch as final
ARG TARGETPLATFORM

COPY --from=artifact /etc/passwd /etc/passwd
COPY --from=artifact /etc/ssl/certs /etc/ssl/certs
COPY --from=artifact /app /app

HEALTHCHECK NONE

USER zitadel
ENTRYPOINT ["/app/zitadel"]