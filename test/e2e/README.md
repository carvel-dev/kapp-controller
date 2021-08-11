## End-to-end test

To run

```bash
$ ./hack/test-e2e.sh
$ ./hack/test-e2e.sh -run TestVersion
```

See `./test/e2e/env.go` for required environment variables for some tests.

If running tests for kapp-controller's integration with secretgen-controller, 
you will need to set the environment variable `KAPPCTRL_E2E_SECRETGEN_CONTROLLER=true`. 
