# Build
FROM golang:1.13.6-alpine3.10 as builder
COPY . /app
WORKDIR /app
RUN go build -o /main /app/main.go

# Make blockchain env
FROM alpine:3.10 as cli

COPY ./pkg/cert-issuer /cert-issuer-cli
COPY ./pkg/cert-issuer/conf_regtest.ini /etc/cert-issuer/conf.ini

RUN apk add --update \
        bash \
        ca-certificates \
        curl \
        gcc \
        gmp-dev \
        libffi-dev \
        libressl-dev \
        linux-headers \
        make \
        musl-dev \
        python \
        python3 \
        python3-dev \
        tar \
    && python3 -m ensurepip \
    && pip3 install --upgrade pip setuptools \
    && mkdir -p /etc/cert-issuer/data/unsigned_certificates \
    && mkdir /etc/cert-issuer/data/blockchain_certificates \
    && pip3 install /cert-issuer-cli/. \
    && rm -r /usr/lib/python*/ensurepip \
    && rm -rf /var/cache/apk/* \
    && rm -rf /root/.cache \
    && sed -i.bak s/==1\.0b1/\>=1\.0\.2/g /usr/lib/python3.*/site-packages/merkletools-1.0.2-py3.*.egg-info/requires.txt

# Run
WORKDIR /cert-issuer-cli
RUN python3 setup.py experimental --blockchain=ethereum

WORKDIR /app
COPY ./Makefile .
COPY --from=builder /main .

CMD [ "./main" ]