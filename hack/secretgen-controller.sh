deploy_secretgen-controller() {
  if [ "$KAPPCTRL_E2E_SECRETGEN_CONTROLLER" == "true" ]; then
    echo "Deploying secretgen-controller..."
    kapp deploy -a sg -f https://github.com/vmware-tanzu/carvel-secretgen-controller/releases/download/v0.7.1/release.yml -c -y
  else
    echo "Skipping secretgen-controller deployment"
  fi
}
