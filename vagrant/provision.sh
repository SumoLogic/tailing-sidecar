#!/bin/bash

set -x

export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get --yes upgrade

apt-get install --yes make

# Install docker
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get install --yes docker-ce docker-ce-cli containerd.io
usermod -aG docker vagrant

# Install k8s
snap install microk8s --classic --channel=1.19/stable
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
