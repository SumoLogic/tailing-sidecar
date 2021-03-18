all: markdownlint yamllint

markdownlint: mdl

mdl:
	mdl --style .markdownlint/style.rb \
		README.md \
		sidecar/README.md \
		operator/README.md \
		operator/docs \
		dev/releasing.md

yamllint:
	yamllint -c .yamllint.yaml \
		operator/examples/

build-push-deploy: build-push-sidecar build-push-deploy-operator

build-push-sidecar:
	$(MAKE) -C sidecar all

build-push-deploy-operator:
	$(MAKE) -C operator all

push-helm-chart:
	./ci/push-helm-chart.sh
