#! /bin/bash
# scrapes Dockerfile for wget commands for fetching current versions of tanzu carvel tools,
#  replaces "linux" with "darwin" so we pull the osx images, cleans up lingering Docker and multiline-grep formatting,
# executes resulting commands (fetch executable, chmod executable, run executable --version to verify)
cat ../Dockerfile | \
grep wget -A2 | \
grep tanzu -A2 | \
grep -v echo | \
sed "s/linux/darwin/" | \
sed 's/RUN //' | \
grep -v "\-\-" | \
bash
