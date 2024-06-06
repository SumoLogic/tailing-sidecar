#!/bin/bash

set -x

export DEBIAN_FRONTEND=noninteractive
GO_VERSION="1.21.5"
HELM_VERSION=v3.5.2
KUTTL_VERSION=0.15.0
MICROK8S_VERSION=1.27
ARCH="$(dpkg --print-architecture)"
KUTTL_ARCH="${ARCH}"
if [[ "${KUTTL_ARCH}" == "amd64" ]]; then
  KUTTL_ARCH="x86_64";
fi

apt-get update
apt-get --yes upgrade
apt-get install --yes make gcc

# Install docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
   "deb [arch=${ARCH}] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get install --yes docker-ce docker-ce-cli containerd.io
usermod -aG docker vagrant

# Install k8s
snap install microk8s --classic --channel=${MICROK8S_VERSION}/stable
microk8s.status --wait-ready
ufw allow in on cbr0
ufw allow out on cbr0
ufw default allow routed

microk8s enable registry
microk8s enable storage
microk8s enable dns

microk8s.kubectl config view --raw > /tailing-sidecar/.kube-config

snap alias microk8s.kubectl kubectl

usermod -a -G microk8s vagrant

echo "export KUBECONFIG=/var/snap/microk8s/current/credentials/kubelet.config" >> /home/vagrant/.bashrc

# Install go
wget "https://golang.org/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz"
tar -C /usr/local -xzf "go${GO_VERSION}.linux-${ARCH}.tar.gz"
rm "go${GO_VERSION}.linux-${ARCH}.tar.gz"
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.bashrc

# Install operator SDK
curl -LO "https://github.com/operator-framework/operator-sdk/releases/latest/download/operator-sdk_linux_${ARCH}"
chmod +x "operator-sdk_linux_${ARCH}"
mv "operator-sdk_linux_${ARCH}" /usr/local/bin/operator-sdk

# Install kustomize
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
mv kustomize /usr/local/bin/

# Install Helm
mkdir /opt/helm3
curl "https://get.helm.sh/helm-${HELM_VERSION}-linux-${ARCH}.tar.gz" | tar -xz -C /opt/helm3
ln -s "/opt/helm3/linux-${ARCH}/helm" /usr/bin/helm3
ln -s /usr/bin/helm3 /usr/bin/helm

# Check if k8s is ready
while true; do
  kubectl -n kube-system get services 1>/dev/null 2>&1 && break
  echo 'Waiting for k8s server'
  sleep 1
done

# Deploy cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.11.0/cert-manager.yaml

# Check if cert-manager is ready
# NOTICE: kubectl wait is not used due to unexpected errors
# https://github.com/ubuntu/microk8s/issues/1710
while true; do
  if [ $( kubectl get pods -n cert-manager  | grep Running | wc -l ) -eq 3 ]; then
    echo 'cert-manager is ready'
    break
  else
    echo 'Waiting for cert-manager'
  fi
  sleep 5
done

# Install kuttl
curl -L "https://github.com/kudobuilder/kuttl/releases/download/v${KUTTL_VERSION}/kubectl-kuttl_${KUTTL_VERSION}_linux_${KUTTL_ARCH}" --output kubectl-kuttl
chmod +x kubectl-kuttl
mv kubectl-kuttl /usr/local/bin/kubectl-kuttl

# For AMD64 / x86_64
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-${ARCH}
chmod +x kind
mv kind /usr/local/bin/kind
