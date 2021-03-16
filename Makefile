all: markdownlint yamllint

markdownlint: mdl

mdl:
	mdl --style .markdownlint/style.rb \
		README.md \
		sidecar/README.md \
		operator/README.md \
		operator/docs

yamllint:
	yamllint -c .yamllint.yaml \
		operator/examples/

build-push-deploy: build-push-sidecar build-push-deploy-operator

build-push-sidecar:
	$(MAKE) -C sidecar all

build-push-deploy-operator:
	$(MAKE) -C operator all

login:
	echo "${CR_PAT}" | docker login ghcr.io -u "${CR_OWNER}" --password-stdin
