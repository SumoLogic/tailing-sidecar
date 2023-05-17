NAMESPACE ?= tailing-sidecar-system
RELEASE ?= tailing-sidecar
HELM_CHART ?= helm/tailing-sidecar-operator

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
