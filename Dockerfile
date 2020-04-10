# Build the wkhtmltopdf
FROM surnet/alpine-wkhtmltopdf:3.9-0.12.5-full as wkhtmltopdf-builder

# Build
FROM golang:1.13.6-alpine3.10 as builder
COPY . /app
WORKDIR /app
RUN go build -o /main /app/main.go

COPY docker/fonts/* /usr/share/fonts/truetype/

# Make blockchain env
FROM alpine:3.10 as cli

COPY ./issuing-service/pkg/cert-issuer /cert-issuer-cli
COPY ./issuing-service/pkg/cert-issuer/conf_regtest.ini /etc/cert-issuer/conf.ini

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
# For HTML to PDF converting
        nodejs nodejs-npm \
    && python3 -m ensurepip \
    && pip3 install --upgrade pip setuptools \
    && mkdir -p /etc/cert-issuer/data/unsigned_certificates \
    && mkdir /etc/cert-issuer/data/blockchain_certificates \
    && pip3 install /cert-issuer-cli/. \
    && rm -r /usr/lib/python*/ensurepip \
    && rm -rf /var/cache/apk/* \
    && rm -rf /root/.cache \
    && sed -i.bak s/==1\.0b1/\>=1\.0\.2/g /usr/lib/python3.*/site-packages/merkletools-1.0.2-py3.*.egg-info/requires.txt

# Install dependencies for wkhtmltopdf
RUN apk add --no-cache \
      libstdc++ \
      libx11 \
      libxrender \
      libxext \
      libssl1.1 \
      ca-certificates \
      fontconfig \
      freetype \
      ttf-dejavu \
      ttf-droid \
      ttf-freefont \
      ttf-liberation \
      ttf-ubuntu-font-family \
    && apk add --no-cache --virtual .build-deps \
      msttcorefonts-installer \
# Install microsoft fonts
    && update-ms-fonts \
    && fc-cache -f \
# Clean up when done
    && rm -rf /tmp/* \
    && apk del .build-deps

# Copy wkhtmltopdf files from docker-wkhtmltopdf image
COPY --from=wkhtmltopdf-builder /bin/wkhtmltopdf /bin/wkhtmltopdf
COPY --from=wkhtmltopdf-builder /bin/wkhtmltoimage /bin/wkhtmltoimage
COPY --from=wkhtmltopdf-builder /bin/libwkhtmltox* /bin/

# Run
WORKDIR /cert-issuer-cli
RUN python3 setup.py experimental --blockchain=ethereum

WORKDIR /app
COPY ./Makefile .
COPY --from=builder /main .

CMD [ "./main" ]