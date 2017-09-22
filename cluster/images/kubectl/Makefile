TEMP_DIR:=$(shell mktemp -d -t kubectlXXXXXX)

ifndef VERSION
    $(error VERSION is undefined)
endif

build:
	cp -r ./* ${TEMP_DIR}
	cp ../../../_output/dockerized/bin/linux/amd64/kubectl ${TEMP_DIR}
	docker build --pull -t kubectl-amd64:${VERSION} ${TEMP_DIR}
	rm -rf "${TEMP_DIR}"
