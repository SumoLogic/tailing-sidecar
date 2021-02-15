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
