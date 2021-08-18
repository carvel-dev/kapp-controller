deploy_secretgen-controller() {
  if [ "$KAPPCTRL_E2E_SECRETGEN_CONTROLLER" == "true" ]; then
    echo "Deploying secretgen-controller..."
    kapp deploy -a sg -f https://raw.githubusercontent.com/vmware-tanzu/carvel-secretgen-controller/develop/alpha-releases/0.4.0-alpha.1.yml -c -y
  else
    echo "Skipping secretgen-controller deployment"
  fi
}
