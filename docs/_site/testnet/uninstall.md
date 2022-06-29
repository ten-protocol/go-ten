# Uninstall
## Uninstall the Docker Engine, CLI, Containerd, and Docker Compose packages:

     sudo apt-get purge docker-ce docker-ce-cli containerd.io docker-compose-plugin

Images, containers, volumes, or customized configuration files on your host are not automatically removed. To delete all images, containers, and volumes:

     sudo rm -rf /var/lib/docker
     sudo rm -rf /var/lib/containerd

You must delete any edited configuration files manually.

## Uninstall SGX 
1. Run the uninstall shell script:

     `sudo /opt/intel/sgxdriver/uninstall.sh`

1. Uninstall the rest of the dependencies:

     `sudo apt purge -y libsgx-enclave-common libsgx-enclave-common-dev libsgx-urts sgx-aesm-service libsgx-uae-service libsgx-launch libsgx-aesm-launch-plugin libsgx-ae-le`