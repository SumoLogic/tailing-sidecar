{{- define "operator.webhookWithCertManager" }}
apiVersion: admissionregistration.k8s.io/v1beta1
kind: MutatingWebhookConfiguration
metadata:
  annotations:
    cert-manager.io/inject-ca-from: {{ .Release.Namespace }}/tailing-sidecar-serving-cert
  name: tailing-sidecar-mutating-webhook-configuration
  namespace: {{ .Release.Namespace }}
webhooks:
- clientConfig:
    caBundle: Cg==
    service:
      name: {{ include "operator.fullname" . }}
      namespace: {{ .Release.Namespace }}
      path: /add-tailing-sidecars-v1-pod
  failurePolicy:  Ignore
  name: tailing-sidecar.sumologic.com
  rules:
  - apiGroups:
    - ""
    apiVersions:
    - v1
    operations:
    - CREATE
    - UPDATE
    resources:
    - pods
---
apiVersion: cert-manager.io/v1alpha2
kind: Certificate
metadata:
  labels:
    {{- include "operator.labels" . | nindent 4 }}
  name: tailing-sidecar-serving-cert
  namespace: {{ .Release.Namespace }}
spec:
  dnsNames:
  - {{ include "operator.fullname" . }}.{{ .Release.Namespace }}.svc
  - {{ include "operator.fullname" . }}.{{ .Release.Namespace }}.svc.cluster.local
  issuerRef:
    kind: Issuer
    name: tailing-sidecar-selfsigned-issuer
  secretName: webhook-server-cert
---
apiVersion: cert-manager.io/v1alpha2
kind: Issuer
metadata:
  labels:
    {{- include "operator.labels" . | nindent 4 }}
  name: tailing-sidecar-selfsigned-issuer
  namespace: {{ .Release.Namespace }}
spec:
  selfSigned: {}
{{- end }}
