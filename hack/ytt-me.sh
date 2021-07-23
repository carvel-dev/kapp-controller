# ytt-me is intended to be replace ytt -f args... in a chain of piped commands and so does not provide its own shebang
ytt -f config/ -v kapp_controller_version="v`cat cmd/main.go| sed  -ne 's/\tVersion = \"\(.*\)\"/\1/p'`"
