.PHONY: issue
issue:
	@if [ "${APP_ENV}" = "test" ]; then\
		echo "/usr/bin/cert-issuer -c ${CONF_PATH}";\
	else\
		/usr/bin/cert-issuer -c "${CONF_PATH}";\
	fi

.PHONY: htmltopdf
CHROME_BIN?=/usr/bin/chromium-browser
htmltopdf:
	env CHROME_BIN=${CHROME_BIN} node ./pkg/htmltopdf/index.js ${HTML_FILEPATH} ${PDF_FILEPATH}

.DEFAULT_GOAL := issue
