#!/bin/bash
docker build -t test .
dd if=/dev/zero of=rootfs.ext4 bs=1M count=2048
mkfs.ext4 rootfs.ext4
mkdir -p rootfs
mount rootfs.ext4 rootfs
docker run --rm -it -v $(pwd)/rootfs:/mnt/rootfs test cp /rootfs /mnt -R

echo "PermitRootLogin yes" >> ./rootfs/etc/ssh/sshd_config
echo "/dev/vda  /   ext4   defaults    1   1" >> ./rootfs/etc/fstab
echo "nameserver 8.8.8.8" > ./rootfs/etc/resolv.conf
mkdir -p ./rootfs/usr/local/bin
echo """#!/bin/bash

FLAG="/var/log/fc-init.log"
if [[ ! -f \$FLAG ]]; then
    echo "text" > /etc/hostname
    echo "This is the first boot"
    touch "\$FLAG"
    reboot
else
    echo "Do nothing"
fi
""" > ./rootfs/usr/local/bin/fc-init.sh
chmod +x ./rootfs/usr/local/bin/fc-init.sh

umount rootfs
rm -d rootfs