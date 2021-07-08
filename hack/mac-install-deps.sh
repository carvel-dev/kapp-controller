#! /bin/bash
# to recreate this file run :
# cat ../Dockerfile | grep wget -A2 | grep tanzu -A2 | grep -v echo | sed "s/linux/darwin/" | sed 's/RUN //' | grep -v "\-\-" > mac-install-deps.
sh

wget -O- https://github.com/vmware-tanzu/carvel-ytt/releases/download/v0.34.0/ytt-darwin-amd64 > /usr/local/bin/ytt && \
  chmod +x /usr/local/bin/ytt && ytt version
wget -O- https://github.com/vmware-tanzu/carvel-kapp/releases/download/v0.37.0/kapp-darwin-amd64 > /usr/local/bin/kapp && \
  chmod +x /usr/local/bin/kapp && kapp version
wget -O- https://github.com/vmware-tanzu/carvel-kbld/releases/download/v0.30.0/kbld-darwin-amd64 > /usr/local/bin/kbld && \
  chmod +x /usr/local/bin/kbld && kbld version
wget -O- https://github.com/vmware-tanzu/carvel-imgpkg/releases/download/v0.12.0/imgpkg-darwin-amd64 > /usr/local/bin/imgpkg && \
  chmod +x /usr/local/bin/imgpkg && imgpkg version
wget -O- https://github.com/vmware-tanzu/carvel-vendir/releases/download/v0.21.1/vendir-darwin-amd64 > /usr/local/bin/vendir && \
  chmod +x /usr/local/bin/vendir && vendir version
