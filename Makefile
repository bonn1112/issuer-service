.PHONY: issue
issue:
	/usr/bin/cert-issuer -c "${CONF_PATH}"

.DEFAULT_GOAL := issue