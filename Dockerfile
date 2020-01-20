FROM seegno/bitcoind:0.13-alpine as cli

ENV GOLANG_VERSION 1.13.6

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
        git \
        go \
    && python3 -m ensurepip \
    && pip3 install --upgrade pip setuptools \
    && mkdir -p /etc/cert-issuer/data/unsigned_certificates \
    && mkdir /etc/cert-issuer/data/blockchain_certificates \
    && mkdir ~/.bitcoin \
    && echo $'rpcuser=foo\nrpcpassword=bar\nrpcport=8332\nregtest=1\nrelaypriority=0\nrpcallowip=127.0.0.1\nrpcconnect=127.0.0.1\n' > /root/.bitcoin/bitcoin.conf \
    && pip3 install /cert-issuer-cli/. \
    && rm -r /usr/lib/python*/ensurepip \
    && rm -rf /var/cache/apk/* \
    && rm -rf /root/.cache \
    && sed -i.bak s/==1\.0b1/\>=1\.0\.2/g /usr/lib/python3.*/site-packages/merkletools-1.0.2-py3.*.egg-info/requires.txt

CMD bitcoind -daemon && bash



# Golang Dockerfile

RUN export \
		GOOS="$(go env GOOS)" \
		GOARCH="$(go env GOARCH)" \
		GOHOSTOS="$(go env GOHOSTOS)" \
		GOHOSTARCH="$(go env GOHOSTARCH)" \
	; \
	apkArch="$(apk --print-arch)"; \
		case "$apkArch" in \
			armhf) export GOARM='6' ;; \
			armv7) export GOARM='7' ;; \
			x86) export GO386='387' ;; \
		esac; \
		\
		wget -O go.tgz "https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz"; \
		echo 'aae5be954bdc40bcf8006eb77e8d8a5dde412722bc8effcdaf9772620d06420c *go.tgz' | sha256sum -c -; \
		tar -C /usr/local -xzf go.tgz; \
		rm go.tgz; \
		\
		cd /usr/local/go/src; \
		./make.bash; \
		\
		rm -rf \
			/usr/local/go/pkg/bootstrap \
			/usr/local/go/pkg/obj \
		; \
		apk del .build-deps; \
		\
		export PATH="/usr/local/go/bin:$PATH"; \
		go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"



# Base logic

ADD ./storage /storage

WORKDIR /cert-issuer-cli
RUN python3 setup.py experimental --blockchain=ethereum

ADD ./issuing-service /cert-issuer
WORKDIR /cert-issuer

RUN go get github.com/oxequa/realize

CMD [ "realize", "start" ]