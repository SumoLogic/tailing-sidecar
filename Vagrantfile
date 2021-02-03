# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure('2') do |config|
  config.vm.box = 'ubuntu/bionic64'
  config.disksize.size = '50GB'
  config.vm.box_check_update = false
  config.vm.host_name = 'sumologic-tailing-sidecar'
  config.vm.network :private_network, ip: "192.168.78.67"

  config.vm.provider 'virtualbox' do |vb|
    vb.gui = false
    vb.cpus = 4
    vb.memory = 8192
    vb.name = 'sumologic-tailing-sidecar'
  end

  config.vm.provision 'shell', path: 'vagrant/provision.sh'
  config.vm.synced_folder ".", "/tailing-sidecar"
end
