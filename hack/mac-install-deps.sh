#! /bin/bash
# to recreate this file run :
#  cat ../Dockerfile | grep wget | grep tanzu | sed "s/linux/darwin/" | sed 's/RUN //' > mac-install-deps.sh
# and then remove the final `&& \` manually or write one more sed with an end-of-file matcher whatever 

wget -O- https://github.com/vmware-tanzu/carvel-ytt/releases/download/v0.34.0/ytt-darwin-amd64 > /usr/local/bin/ytt && \
wget -O- https://github.com/vmware-tanzu/carvel-kapp/releases/download/v0.37.0/kapp-darwin-amd64 > /usr/local/bin/kapp && \
wget -O- https://github.com/vmware-tanzu/carvel-kbld/releases/download/v0.30.0/kbld-darwin-amd64 > /usr/local/bin/kbld && \
wget -O- https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/v0.12.0/imgpkg-darwin-amd64 > /usr/local/bin/imgpkg && \
wget -O- https://github.com/vmware-tanzu/carvel-vendir/releases/download/v0.21.1/vendir-darwin-amd64 > /usr/local/bin/vendir

