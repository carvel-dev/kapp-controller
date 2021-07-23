# ytt-me is intended to be replace ytt -f args... in a chain of piped commands and so does not provide its own shebang

if [ -z ${KAPP_CONTROLLER_VERSION} ]; then source hack/version-me.sh; fi

ytt -f config/ -v kapp_controller_version="${KAPP_CONTROLLER_VERSION}"
