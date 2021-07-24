# ytt-me is intended to be replace ytt -f args... in a chain of piped commands and so does not provide its own shebang

source $(dirname "$0")/version-me.sh

ytt -f config/ -v kapp_controller_version="$(get_kappctrl_ver)"
