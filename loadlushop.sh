#!/bin/bash

sudo ctr -n k8s.io images import /home/zzx/lucien/images/mysql-8.tar
sudo ctr -n k8s.io images import /home/zzx/lucien/images/redis-7-alpine.tar
sudo ctr -n k8s.io images import /home/zzx/lucien/images/apache-rocketmq-5.3.3.tar
sudo ctr -n k8s.io images import /home/zzx/lucien/images/rocketmq-dashboard-2.0.1.tar
sudo ctr -n k8s.io images import /home/zzx/lucien/images/nacos-server-2.3.2.tar
sudo ctr -n k8s.io images import /home/zzx/lucien/images/lushop-images.tar
sudo ctr -n k8s.io images import /home/zzx/lucien/images/busybox.tar