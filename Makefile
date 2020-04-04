.PHONY: issue
issue:
	/usr/bin/cert-issuer -c "${CONF_PATH}" > /storage/issuing_service.log 2>&1

.DEFAULT_GOAL := issue