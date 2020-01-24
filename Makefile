CERT_ISSUER_EXECUTABLE?=/usr/bin/cert-issuer

.PHONY: issue
issue:
	${CERT_ISSUER_EXECUTABLE} -c "${CONF_PATH}"

.DEFAULT_GOAL := issue