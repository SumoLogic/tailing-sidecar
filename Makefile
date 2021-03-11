all: markdownlint yamllint

markdownlint: mdl

mdl:
	mdl --style .markdownlint/style.rb \
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
