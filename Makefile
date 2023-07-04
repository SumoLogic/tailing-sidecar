NAMESPACE ?= tailing-sidecar-system
RELEASE ?= tailing-sidecar
HELM_CHART ?= helm/tailing-sidecar-operator
KUTTL_CONFIG ?= kuttl-test.yaml

all: markdownlint yamllint

markdownlint: mdl

mdl:
	mdl --style .markdownlint/style.rb \
		README.md \
		sidecar/README.md \
		operator/README.md \
		docs/*.md

yamllint:
	yamllint -c .yamllint.yaml \
		operator/examples/

login-ecr:
	aws ecr-public get-login-password --region us-east-1 \
	| docker login --username AWS --password-stdin $(ECR_URL)

.PHONY: e2e
e2e: IMG="registry.localhost:5000/sumologic/tailing-sidecar-operator:test"
e2e: TAILING_SIDECAR_IMG = "registry.localhost:5000/sumologic/tailing-sidecar:test"
e2e:
	$(MAKE) -C ./sidecar build TAG=$(TAILING_SIDECAR_IMG)
	$(MAKE) -C ./operator docker-build IMG=$(IMG) TAILING_SIDECAR_IMG=$(TAILING_SIDECAR_IMG)
	kubectl-kuttl test --config $(KUTTL_CONFIG)

.PHONY: e2e-helm
e2e-helm: KUTTL_CONFIG = kuttl-test-helm.yaml
e2e-helm: e2e

.PHONY: e2e-helm-certmanager
e2e-helm-certmanager: KUTTL_CONFIG = kuttl-test-helm-certmanager.yaml
e2e-helm-certmanager: e2e

.PHONY: e2e-helm-custom-configuration
e2e-helm-custom-configuration: KUTTL_CONFIG = kuttl-test-helm-custom-configuration.yaml
e2e-helm-custom-configuration: e2e

build-push-deploy: build-push-sidecar build-push-deploy-operator

build-push-sidecar:
	$(MAKE) -C sidecar all

build-push-deploy-operator:
	$(MAKE) -C operator all

push-helm-chart:
	./ci/push-helm-chart.sh

helm-upgrade:
	helm upgrade --install $(RELEASE) \
		--namespace $(NAMESPACE) \
		--create-namespace \
		$(HELM_CHART)

helm-dry-run:
	helm install --dry-run $(RELEASE) \
		--namespace $(NAMESPACE) \
		$(HELM_CHART)

helm-delete:
	helm delete $(RELEASE) --namespace $(NAMESPACE)

deploy-examples:
	$(MAKE) -C operator deploy-examples

check-examples:
	$(MAKE) -C operator check-examples

teardown-examples:
	$(MAKE) -C operator teardown-examples
