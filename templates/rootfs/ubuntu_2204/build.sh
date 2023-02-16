#!/bin/bash
docker build -t test .
rm rootfs.ext4
truncate -s 1G rootfs.ext4
mkfs.ext4 -E nodiscard rootfs.ext4
mkdir -p rootfs
mount rootfs.ext4 rootfs
docker run --rm -it -v $(pwd)/rootfs:/mnt/rootfs test cp /rootfs /mnt -R

echo "nameserver 8.8.8.8" > ./rootfs/etc/resolv.conf

umount rootfs
rm -d rootfs