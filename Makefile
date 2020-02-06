.PHONY: issue
issue:
	@if [ "${APP_ENV}" = "test" ]; then\
		echo "/usr/bin/cert-issuer -c ${CONF_PATH}";\
	else\
		/usr/bin/cert-issuer -c "${CONF_PATH}";\
	fi

.DEFAULT_GOAL := issue

qwe:
	@echo "qwe"