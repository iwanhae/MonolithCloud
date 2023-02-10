FROM ubuntu:22.04 AS guest
RUN apt-get update 
RUN apt install udev kmod -y
RUN apt install openssh-server sudo cron -y
RUN apt install -y vim htop

RUN crontab -l | { cat; echo "@reboot bash /usr/local/bin/fc-init.sh"; } | crontab -
RUN echo "ubuntu\nubuntu" | passwd

FROM ubuntu:22.04
COPY --from=guest / /rootfs