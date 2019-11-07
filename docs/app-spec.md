## App spec

```yaml
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: simple-app
  # This namespace is going to be used as a default namespace during deploy
  namespace: ns
spec:
  # pauses _future_ reconcilation; does _not_ affect
  # currently running reconciliation (optional; default=false)
  paused: true

  # Fetch must have one or more directives
  fetch:
    - inline:
        # specifies content inline within resource;
        # not recommended for sensitive values as CR is not encrypted (optional)
        paths:
          # mapping of paths to their content
          dir/file.ext: file-content
        # specified content via secrets and config maps;
        # data values are recommended to be placed in secrets (optional)
        pathsFrom:
          - secretRef:
              name: secret-name
              # specifies where to place files found in secret (optional)
              directoryPath: dir
          - configMapRef:
              name: cfgmap-name
              # specifies where to place files found in config map (optional)
              directoryPath: dir

    # pulls content from Docker/OCI registry
    - image:
        # Docker image url; unqualified, tagged, or
        # digest references supported (required)
        url: host.com/username/image:v0.1.0
        # secret with auth details (optional)
        secretRef:
          name: secret-name
        # grab only portion of image (optional)
        subPath: inside-dir/dir2

    # uses http library to fetch file
    - http:
        # http and https url are supported;
        # plain file, tgz and tar types are supported (required)
        url: https://host.com/archive.tgz
        # checksum to verify after download (optional)
        sha256: 0a12cdef83...
        # secret to provide auth details (optional)
        secretRef:
          name: secret-name
        # grab only portion of download (optional)
        subPath: inside-dir/dir2

    # uses git to clone repository
    - git:
        # http or ssh urls are supported (required)
        url: https://github.com/k14s/k8s-simple-app-example
        # branch, tag, commit; origin is the name of the remote (required)
        ref: origin/master
        # secret with auth details (optional)
        secretRef:
          name: secret-name
        # grab only portion of repository (optional)
        subPath: config-step-2-template
        # skip lfs download (optional)
        lfsSkipSmudge: true

    # uses helm fetch to fetch specified chart
    - helmChart:
        name: stable/nginx
        # (optional)
        version: "0.1.0"
        # (optional)
        repository:
          url: https://...
          # (optional)
          secretRef:
            name: secret-name

  # Template must have one or more directives
  template:
    - ytt:
        # ignores comments that ytt doesn't recognize
        # (optional; default=false)
        ignoreUnknownComments: true
        # forces strict mode https://github.com/k14s/ytt/blob/master/docs/strict.md
        # (optional; default=false)
        strict: true
        # specify additional files, including data values (optional)
        inline:
          # specifies content inline within resource;
          # not recommended for sensitive values as CR is not encrypted (optional)
          paths:
            # mapping of paths to their content
            dir/file.ext: |
              file-content
              file-content
          # specified content via secrets and config maps;
          # data values are recommended to be placed in secrets (optional)
          pathsFrom:
            - secretRef:
                name: secret-name
                # specifies where to place files found in secret (optional)
                directoryPath: dir
            - configMapRef:
                name: cfgmap-name
                # specifies where to place files found in config map (optional)
                directoryPath: dir

    # use kbld to resolve image references to use digests
    - kbld: {}

    # use helm template command to render helm chart
    - helmTemplate:
        # one or more secrets or config maps that provide values (optional)
        valuesFrom:
          - secretRef:
              name: secret-name
          - configMapRef:
              name: cfgmap-name

  # Deploy must have one directive
  deploy:

    # use kapp to deploy resources
    - kapp:
        # override namespace for all resources (optional)
        intoNs: another-ns1
        # provide custom namespace override mapping (optional)
        mapNs: ["ns1=another-ns1"]
        # pass through options to kapp deploy (optional)
        rawOptions: ["--apply-concurrency=10"]
        # configuration for delete command (optional)
        delete:
          # pass through options to kapp delete (optional)
          rawOptions: ["--apply-ignored=true"]
```
