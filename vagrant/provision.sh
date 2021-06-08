#!/bin/bash

set -x

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get --yes upgrade

apt-get install --yes make gcc

# Install docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get install --yes docker-ce docker-ce-cli containerd.io
usermod -aG docker vagrant

# Install k8s
snap install microk8s --classic --channel=1.20/stable
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
GO_VERSION="1.15.7"
wget "https://golang.org/dl/go${GO_VERSION}.linux-amd64.tar.gz"
tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> /home/vagrant/.bashrc

# Install operator SDK
curl -LO "https://github.com/operator-framework/operator-sdk/releases/latest/download/operator-sdk_linux_amd64"
chmod +x operator-sdk_linux_amd64
mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk

# Install kustomize
curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"  | bash
mv kustomize /usr/local/bin/

# Install Helm
HELM_VERSION=v3.5.2
mkdir /opt/helm3
curl "https://get.helm.sh/helm-${HELM_VERSION}-linux-amd64.tar.gz" | tar -xz -C /opt/helm3
ln -s /opt/helm3/linux-amd64/helm /usr/bin/helm3
ln -s /usr/bin/helm3 /usr/bin/helm

# Check if k8s is ready
while true; do
  kubectl -n kube-system get services 1>/dev/null 2>&1 && break
  echo 'Waiting for k8s server'
  sleep 1
done

# Deploy cert-manager
kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.1.0/cert-manager.yaml

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
