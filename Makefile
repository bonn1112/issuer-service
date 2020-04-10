.PHONY: issue
issue:
	@if [ "${APP_ENV}" = "test" ]; then\
		echo "/usr/bin/cert-issuer -c ${CONF_PATH}";\
	else\
		/usr/bin/cert-issuer -c "${CONF_PATH}";\
	fi

.PHONY: htmltopdf
htmltopdf:
	node ./pkg/htmltopdf/index.js ${HTML_FILEPATH} ${PDF_FILEPATH}

.DEFAULT_GOAL := issue
