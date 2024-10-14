#!/bin/bash
vendor:
	@echo "Using ART_USERNAME=${ART_USERNAME} ART_TOKEN=${ART_TOKEN} make vendor"
	rm -rvf vendor
	GONOSUMDB=kbtg.tech/* GO111MODULE=on \
	GOPROXY="https://${ART_USERNAME}:${ART_TOKEN}@artifactory.kasikornbank.com:8443/artifactory/api/go/virtual-kpp" \
		go mod vendor -v
