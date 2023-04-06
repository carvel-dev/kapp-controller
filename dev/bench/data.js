window.BENCHMARK_DATA = {
  "lastUpdate": 1680771121301,
  "repoUrl": "https://github.com/carvel-dev/kapp-controller",
  "entries": {
    "Benchmark": [
      {
        "commit": {
          "author": {
            "name": "vmware-tanzu",
            "username": "vmware-tanzu"
          },
          "committer": {
            "name": "vmware-tanzu",
            "username": "vmware-tanzu"
          },
          "id": "6200ea5083ca1db200b904390e7f2d958c8321ed",
          "message": "pkgr benchmark test",
          "timestamp": "2022-04-06T17:27:49Z",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/pull/629/commits/6200ea5083ca1db200b904390e7f2d958c8321ed"
        },
        "date": 1650324813749,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 93995132319,
            "unit": "ns/op\t        63.23 DeleteSeconds\t        30.72 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 23716216909,
            "unit": "ns/op\t        15.48 DeleteSeconds\t         8.197 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 14601909907,
            "unit": "ns/op\t         9.418 DeleteSeconds\t         5.148 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "name": "vmware-tanzu",
            "username": "vmware-tanzu"
          },
          "committer": {
            "name": "vmware-tanzu",
            "username": "vmware-tanzu"
          },
          "id": "3281853ec58c6c619bdcd0fe5ff2fe5f87db2848",
          "message": "pkgr benchmark test",
          "timestamp": "2022-04-20T16:56:42Z",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/pull/629/commits/3281853ec58c6c619bdcd0fe5ff2fe5f87db2848"
        },
        "date": 1650479550694,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 95229929395,
            "unit": "ns/op\t        64.36 DeleteSeconds\t        30.83 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 23808474566,
            "unit": "ns/op\t        15.56 DeleteSeconds\t         8.203 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 14629206009,
            "unit": "ns/op\t         9.420 DeleteSeconds\t         5.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9df3455df3ef0e2779d47171b1dfeeb53160d8cf",
          "message": "pkgr benchmark test (#629)\n\n* pkgr benchmark test\r\n\r\n* benchmark action: do it on pullrequests even though maybe we wouldn't really want to\r\n\r\n* benchmarks with right argument and version of go\r\n\r\n* benchmark tests get thresholds so on any given test run we can do an absolute time comparison instead of just relative timings between benchmark runs.\r\n\r\n* working on benchmark graphs on github pages\r\n\r\n* write to github pages conditionally on branch name develop only\r\n\r\n* collapse benchmark storage into one block\r\n\r\n* only store benchmark results on develop branch\r\n\r\n* fix the benchmark storage\r\n\r\n* add docs w link to benchmark tests to devmd\r\n\r\n* pin action to sha instead of tag and comment out conditional so we can test\r\n\r\n* restore conditionals",
          "timestamp": "2022-04-20T15:58:40-04:00",
          "tree_id": "3b0a2791411a3236324ef448db255f1a0ae4fd2f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/9df3455df3ef0e2779d47171b1dfeeb53160d8cf"
        },
        "date": 1650485485829,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 94443823631,
            "unit": "ns/op\t        63.52 DeleteSeconds\t        30.88 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 23855612927,
            "unit": "ns/op\t        15.56 DeleteSeconds\t         8.249 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 14642181389,
            "unit": "ns/op\t         9.423 DeleteSeconds\t         5.172 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2935b002334bfcf9515b414641bb0a0dabb088b7",
          "message": "Merge pull request #554 from vmware-tanzu/kctrl-app-commands\n\n`kctrl app` commands",
          "timestamp": "2022-04-21T08:35:22-04:00",
          "tree_id": "2e7834644fe0de9da446a6c764463d5d527b52be",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2935b002334bfcf9515b414641bb0a0dabb088b7"
        },
        "date": 1650545215285,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 94062479801,
            "unit": "ns/op\t        63.23 DeleteSeconds\t        30.80 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 23757471228,
            "unit": "ns/op\t        15.52 DeleteSeconds\t         8.201 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 14610827937,
            "unit": "ns/op\t         9.433 DeleteSeconds\t         5.141 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3f6d2325258ac8a15e3fbfe8ff6959b44b13a0dd",
          "message": "apiserver: custom QPS and Burst to allow high throughput of packages (#635)\n\n* apiserver: custom QPS and Burst to allow high throughput of packages\r\n\r\n* updating benchmark test upper bounds to reflect our speedy new reality\r\n\r\n* smallest numbers that have the same effect",
          "timestamp": "2022-04-21T10:22:06-04:00",
          "tree_id": "1f247d50501b2429ae43ebba8639373800197c51",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3f6d2325258ac8a15e3fbfe8ff6959b44b13a0dd"
        },
        "date": 1650551538244,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36448166003,
            "unit": "ns/op\t        18.93 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9521619928,
            "unit": "ns/op\t         5.280 DeleteSeconds\t         4.207 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6424466071,
            "unit": "ns/op\t         4.247 DeleteSeconds\t         2.142 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "426dc0517b32e623dee854b20eebdb8f43e02c81",
          "message": "remove dead code (#636)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-04-21T13:49:36-04:00",
          "tree_id": "06aeb89d7a851ef802f07348c46cc9adbff03346",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/426dc0517b32e623dee854b20eebdb8f43e02c81"
        },
        "date": 1650563984950,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36454126329,
            "unit": "ns/op\t        18.94 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9432302068,
            "unit": "ns/op\t         5.254 DeleteSeconds\t         4.142 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6423077774,
            "unit": "ns/op\t         4.255 DeleteSeconds\t         2.133 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "600f2a3f155a5afdf93ac9f7daeeb66b62435518",
          "message": "Merge pull request #643 from vmware-tanzu/dependabot/github_actions/reviewdog/action-misspell-1.12",
          "timestamp": "2022-04-26T08:57:25-04:00",
          "tree_id": "3c4e29723a2dfb5f16df7266c74b5c915dc3cc26",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/600f2a3f155a5afdf93ac9f7daeeb66b62435518"
        },
        "date": 1650978486983,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36571937043,
            "unit": "ns/op\t        18.98 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9551408963,
            "unit": "ns/op\t         5.321 DeleteSeconds\t         4.189 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6525658597,
            "unit": "ns/op\t         4.331 DeleteSeconds\t         2.151 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "708b92b0f4ce2673f9d234004d89b5968273cdd7",
          "message": "Merge pull request #641 from vmware-tanzu/dependabot/github_actions/actions/setup-go-3\n\nBump actions/setup-go from 1 to 3",
          "timestamp": "2022-04-26T08:59:31-04:00",
          "tree_id": "8e35206cfed101ab0fdac5c89dc31d1e05230b49",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/708b92b0f4ce2673f9d234004d89b5968273cdd7"
        },
        "date": 1650978707044,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36934301053,
            "unit": "ns/op\t        19.26 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9638468695,
            "unit": "ns/op\t         5.366 DeleteSeconds\t         4.227 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6527526181,
            "unit": "ns/op\t         4.313 DeleteSeconds\t         2.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4c0a29912178cd9337e83c1e8aaba55f20cc3550",
          "message": "Merge pull request #640 from vmware-tanzu/dependabot/github_actions/actions/cache-3\n\nBump actions/cache from 1 to 3",
          "timestamp": "2022-04-26T08:59:59-04:00",
          "tree_id": "eb05d875f77ccc80e0a60672d6958855dbf98716",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4c0a29912178cd9337e83c1e8aaba55f20cc3550"
        },
        "date": 1650978759868,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36836204793,
            "unit": "ns/op\t        19.15 DeleteSeconds\t        17.65 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9621329566,
            "unit": "ns/op\t         5.362 DeleteSeconds\t         4.215 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6583535836,
            "unit": "ns/op\t         4.342 DeleteSeconds\t         2.146 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c5f32fc1202eaa6e2def9191c8399d82a4e4f238",
          "message": "Merge pull request #642 from vmware-tanzu/dependabot/github_actions/github/codeql-action-2",
          "timestamp": "2022-04-27T13:19:30-04:00",
          "tree_id": "0f755f9843056d19cc89e4727acfaa7a5c14d197",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c5f32fc1202eaa6e2def9191c8399d82a4e4f238"
        },
        "date": 1651080576169,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36395157526,
            "unit": "ns/op\t        18.91 DeleteSeconds\t        17.45 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9492121825,
            "unit": "ns/op\t         5.296 DeleteSeconds\t         4.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6413889991,
            "unit": "ns/op\t         4.250 DeleteSeconds\t         2.122 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "42274fbe2fac7c99477516c93044417d464c0840",
          "message": "Merge pull request #645 from vmware-tanzu/dependabot-stops-ignoring-patch-1",
          "timestamp": "2022-04-27T13:20:31-04:00",
          "tree_id": "ac764e65af343900ee067396ce497fd9db69218e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/42274fbe2fac7c99477516c93044417d464c0840"
        },
        "date": 1651080737833,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36899294340,
            "unit": "ns/op\t        19.18 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9643019379,
            "unit": "ns/op\t         5.410 DeleteSeconds\t         4.181 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6524449138,
            "unit": "ns/op\t         4.316 DeleteSeconds\t         2.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "distinct": true,
          "id": "12b66f37ab1d348001acd0253157296c3159c850",
          "message": "Bump k8s.io/apiserver from 0.22.4 to 0.22.9\n\nBumps [k8s.io/apiserver](https://github.com/kubernetes/apiserver) from 0.22.4 to 0.22.9.\n- [Release notes](https://github.com/kubernetes/apiserver/releases)\n- [Commits](https://github.com/kubernetes/apiserver/compare/v0.22.4...v0.22.9)\n\n---\nupdated-dependencies:\n- dependency-name: k8s.io/apiserver\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2022-04-27T17:49:32Z",
          "tree_id": "9ce1d2629f8ae762b34111141524331585fc4fbb",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/12b66f37ab1d348001acd0253157296c3159c850"
        },
        "date": 1651082396719,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36647962144,
            "unit": "ns/op\t        19.05 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9594400778,
            "unit": "ns/op\t         5.374 DeleteSeconds\t         4.177 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6531480074,
            "unit": "ns/op\t         4.351 DeleteSeconds\t         2.135 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "distinct": true,
          "id": "3b02b556f3de6ef9fa794d37768321559f99a918",
          "message": "Bump github.com/stretchr/testify from 1.7.0 to 1.7.1\n\nBumps [github.com/stretchr/testify](https://github.com/stretchr/testify) from 1.7.0 to 1.7.1.\n- [Release notes](https://github.com/stretchr/testify/releases)\n- [Commits](https://github.com/stretchr/testify/compare/v1.7.0...v1.7.1)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/stretchr/testify\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2022-04-27T17:49:43Z",
          "tree_id": "7aa3cf4117765d721d29ece644bfdf9ed7abf61f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3b02b556f3de6ef9fa794d37768321559f99a918"
        },
        "date": 1651082409903,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36606529854,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9577952833,
            "unit": "ns/op\t         5.336 DeleteSeconds\t         4.198 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6475583865,
            "unit": "ns/op\t         4.276 DeleteSeconds\t         2.157 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "distinct": true,
          "id": "bd5628b75aa712d5537ba20705c6e81e67aa3189",
          "message": "Bump golang.org/x/tools from 0.1.5 to 0.1.10\n\nBumps [golang.org/x/tools](https://github.com/golang/tools) from 0.1.5 to 0.1.10.\n- [Release notes](https://github.com/golang/tools/releases)\n- [Commits](https://github.com/golang/tools/compare/v0.1.5...v0.1.10)\n\n---\nupdated-dependencies:\n- dependency-name: golang.org/x/tools\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2022-04-27T18:16:12Z",
          "tree_id": "4d5e0dca311b359551b5f309dfd3ff71fa332ec0",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/bd5628b75aa712d5537ba20705c6e81e67aa3189"
        },
        "date": 1651084110886,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36907305382,
            "unit": "ns/op\t        19.15 DeleteSeconds\t        17.71 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9672478522,
            "unit": "ns/op\t         5.409 DeleteSeconds\t         4.212 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6635186613,
            "unit": "ns/op\t         4.405 DeleteSeconds\t         2.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "43367ce2d476bd35b7c59d72bfc23d03fbcc159c",
          "message": "Merge pull request #637 from vmware-tanzu/kctrl-package\n\n`kctrl`: Adding commands kick, status and pause to `kctrl package installed` command tree",
          "timestamp": "2022-04-28T01:35:57+05:30",
          "tree_id": "8cd69bac23e54d56a58a882e4491dcdf84037578",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/43367ce2d476bd35b7c59d72bfc23d03fbcc159c"
        },
        "date": 1651090654964,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36856125340,
            "unit": "ns/op\t        19.21 DeleteSeconds\t        17.60 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9703534569,
            "unit": "ns/op\t         5.411 DeleteSeconds\t         4.243 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6601508858,
            "unit": "ns/op\t         4.404 DeleteSeconds\t         2.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "distinct": true,
          "id": "a658317283e7a59a0990c1ce0c36729e9d205bf3",
          "message": "Bump github.com/prometheus/client_golang from 1.11.0 to 1.11.1\n\nBumps [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang) from 1.11.0 to 1.11.1.\n- [Release notes](https://github.com/prometheus/client_golang/releases)\n- [Changelog](https://github.com/prometheus/client_golang/blob/main/CHANGELOG.md)\n- [Commits](https://github.com/prometheus/client_golang/compare/v1.11.0...v1.11.1)\n\n---\nupdated-dependencies:\n- dependency-name: github.com/prometheus/client_golang\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2022-04-27T22:53:01Z",
          "tree_id": "0b4391918a81891a78f7707e1ee774f517f777d6",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/a658317283e7a59a0990c1ce0c36729e9d205bf3"
        },
        "date": 1651100675451,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37011372185,
            "unit": "ns/op\t        19.29 DeleteSeconds\t        17.67 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9652940267,
            "unit": "ns/op\t         5.408 DeleteSeconds\t         4.197 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6583752512,
            "unit": "ns/op\t         4.392 DeleteSeconds\t         2.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "distinct": true,
          "id": "4837ff34a7bd2e1def9e423079a6f457d12dd54f",
          "message": "Bump k8s.io/kube-aggregator from 0.22.4 to 0.22.9\n\nBumps [k8s.io/kube-aggregator](https://github.com/kubernetes/kube-aggregator) from 0.22.4 to 0.22.9.\n- [Release notes](https://github.com/kubernetes/kube-aggregator/releases)\n- [Commits](https://github.com/kubernetes/kube-aggregator/compare/v0.22.4...v0.22.9)\n\n---\nupdated-dependencies:\n- dependency-name: k8s.io/kube-aggregator\n  dependency-type: direct:production\n  update-type: version-update:semver-patch\n...\n\nSigned-off-by: dependabot[bot] <support@github.com>",
          "timestamp": "2022-04-28T13:52:14Z",
          "tree_id": "8ca8ab069d594aadcf5f3170d219b452720502ad",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4837ff34a7bd2e1def9e423079a6f457d12dd54f"
        },
        "date": 1651154643167,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36856650638,
            "unit": "ns/op\t        19.15 DeleteSeconds\t        17.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9632350035,
            "unit": "ns/op\t         5.387 DeleteSeconds\t         4.199 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6502862443,
            "unit": "ns/op\t         4.290 DeleteSeconds\t         2.168 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9a85c5cc179ab086fa0f205debb289c64ba5e3a7",
          "message": "Build / download deps in golang image (#651)",
          "timestamp": "2022-04-28T10:53:35-04:00",
          "tree_id": "d3e50d0bf89094f05fcde1d46cde76d262250ab0",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/9a85c5cc179ab086fa0f205debb289c64ba5e3a7"
        },
        "date": 1651158252212,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36638681748,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.60 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9522305217,
            "unit": "ns/op\t         5.317 DeleteSeconds\t         4.167 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6488986384,
            "unit": "ns/op\t         4.298 DeleteSeconds\t         2.139 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cf5d143158dcfccc826aeb57aa4feb5803709427",
          "message": "Merge pull request #656 from vmware-tanzu/dependabot/github_actions/reviewdog/action-misspell-1.12.1\n\nBump reviewdog/action-misspell from 1.12.0 to 1.12.1",
          "timestamp": "2022-04-28T09:57:03-06:00",
          "tree_id": "c31cd9cd2546f07ecb14a94e120e230dd875a719",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cf5d143158dcfccc826aeb57aa4feb5803709427"
        },
        "date": 1651162238185,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37338682300,
            "unit": "ns/op\t        19.43 DeleteSeconds\t        17.85 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9793923623,
            "unit": "ns/op\t         5.493 DeleteSeconds\t         4.229 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6635586537,
            "unit": "ns/op\t         4.356 DeleteSeconds\t         2.223 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "95c7c3817ac71b8a0a809e956318321ea394ed10",
          "message": "Bump vendir to v0.27.0 (#658)",
          "timestamp": "2022-04-29T11:35:09-04:00",
          "tree_id": "806e8f390cb9c71d1790f3ea461a9083e3fc70d1",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/95c7c3817ac71b8a0a809e956318321ea394ed10"
        },
        "date": 1651247283664,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36942382399,
            "unit": "ns/op\t        19.20 DeleteSeconds\t        17.70 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9634018426,
            "unit": "ns/op\t         5.393 DeleteSeconds\t         4.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6608723450,
            "unit": "ns/op\t         4.374 DeleteSeconds\t         2.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d465359a34bac4fc1abb533002d959ec178df9d7",
          "message": "Merge pull request #663 from vmware-tanzu/bump-kapp-v0.47.0",
          "timestamp": "2022-05-05T10:44:52-04:00",
          "tree_id": "bc2eb87c895f1b36e1cdb9750f960a2dcebbd85d",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d465359a34bac4fc1abb533002d959ec178df9d7"
        },
        "date": 1651762640166,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36954605367,
            "unit": "ns/op\t        19.29 DeleteSeconds\t        17.62 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9626167830,
            "unit": "ns/op\t         5.370 DeleteSeconds\t         4.208 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6646270895,
            "unit": "ns/op\t         4.344 DeleteSeconds\t         2.250 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "da3a696293a0862d84308e4283c397538a011bbe",
          "message": "Merge pull request #668 from vmware-tanzu/dependabot/github_actions/slackapi/slack-github-action-1.19.0\n\nBump slackapi/slack-github-action from 1.18.0 to 1.19.0",
          "timestamp": "2022-05-09T15:56:40-06:00",
          "tree_id": "1fe5aac6610abd3bf0246e26a4bd8fb97d072692",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/da3a696293a0862d84308e4283c397538a011bbe"
        },
        "date": 1652134020704,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36520290882,
            "unit": "ns/op\t        18.99 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9509231077,
            "unit": "ns/op\t         5.290 DeleteSeconds\t         4.178 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6511665838,
            "unit": "ns/op\t         4.337 DeleteSeconds\t         2.133 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cec79c2a5d0f6843e17d8b6c642580adf4b8e3b4",
          "message": "Merge pull request #667 from vmware-tanzu/dependabot/github_actions/peter-evans/create-pull-request-4.0.3\n\nBump peter-evans/create-pull-request from 4.0.2 to 4.0.3",
          "timestamp": "2022-05-09T16:00:14-06:00",
          "tree_id": "2bbd2ed2c33f08c646985a5d088e9cd75db7101d",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cec79c2a5d0f6843e17d8b6c642580adf4b8e3b4"
        },
        "date": 1652134213238,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36490042991,
            "unit": "ns/op\t        18.93 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9518559417,
            "unit": "ns/op\t         5.318 DeleteSeconds\t         4.160 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6520150834,
            "unit": "ns/op\t         4.335 DeleteSeconds\t         2.141 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ff2868a7f4490e2d82ad449b29ac85978e72a885",
          "message": "Merge pull request #669 from vmware-tanzu/dependabot/github_actions/docker/login-action-2\n\nBump docker/login-action from 1 to 2",
          "timestamp": "2022-05-09T16:27:28-06:00",
          "tree_id": "35c1fbe426579c67b1f8cac392c74fe1385a02bc",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ff2868a7f4490e2d82ad449b29ac85978e72a885"
        },
        "date": 1652135949284,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36816820899,
            "unit": "ns/op\t        19.12 DeleteSeconds\t        17.65 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9661418378,
            "unit": "ns/op\t         5.388 DeleteSeconds\t         4.227 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6687331107,
            "unit": "ns/op\t         4.437 DeleteSeconds\t         2.177 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "89922ea418c31e052e06cecddf42102b8e41e82c",
          "message": "Remove buildid= workaround (#681)\n\nfixed in go 1.14 https://github.com/golang/go/commit/aa680c0c49b55722a72ad3772e590cd2f9af541d",
          "timestamp": "2022-05-11T08:16:45-04:00",
          "tree_id": "637acb3794b9c066367397686397cf95afe3ad4c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/89922ea418c31e052e06cecddf42102b8e41e82c"
        },
        "date": 1652272215803,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37184159402,
            "unit": "ns/op\t        19.41 DeleteSeconds\t        17.70 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9678717440,
            "unit": "ns/op\t         5.400 DeleteSeconds\t         4.227 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6752736702,
            "unit": "ns/op\t         4.485 DeleteSeconds\t         2.212 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "866424d03cea80e5cb35598f2ab449defe007d87",
          "message": "Merge pull request #685 from vmware-tanzu/bump-imgpkg-v0.29.0\n\nBump imgpkg to v0.29.0",
          "timestamp": "2022-05-12T10:45:43-06:00",
          "tree_id": "1786a548bf2f4988140985fc5cf9dbbf8750c3a8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/866424d03cea80e5cb35598f2ab449defe007d87"
        },
        "date": 1652374569399,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36535148095,
            "unit": "ns/op\t        18.99 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9504140350,
            "unit": "ns/op\t         5.302 DeleteSeconds\t         4.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6433826683,
            "unit": "ns/op\t         4.255 DeleteSeconds\t         2.139 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3cec05cac8a49e4183ee0ff7ee94ec55a277133d",
          "message": "Merge pull request #693 from vmware-tanzu/inject-cmd-runner\n\nInject cmd runner",
          "timestamp": "2022-05-20T16:02:44-07:00",
          "tree_id": "e686029f3c9708a00b48bb2fcf50f4cdc14ce0f7",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3cec05cac8a49e4183ee0ff7ee94ec55a277133d"
        },
        "date": 1653088426673,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36773934345,
            "unit": "ns/op\t        19.09 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9608171319,
            "unit": "ns/op\t         5.353 DeleteSeconds\t         4.204 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6593205607,
            "unit": "ns/op\t         4.353 DeleteSeconds\t         2.192 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f9beed4c0a46a3af3618d8b6fe89a76518fcac13",
          "message": "sops: bump to 3.7.3 (#696)",
          "timestamp": "2022-05-23T15:13:18-04:00",
          "tree_id": "99ea261ee77eded38c93bf37b0ddd99a8664379f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f9beed4c0a46a3af3618d8b6fe89a76518fcac13"
        },
        "date": 1653333795736,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36435613031,
            "unit": "ns/op\t        18.91 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9501534427,
            "unit": "ns/op\t         5.297 DeleteSeconds\t         4.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6515398956,
            "unit": "ns/op\t         4.344 DeleteSeconds\t         2.134 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "bbdb44b376ee6e198aba8c2163524f309c0ddc9a",
          "message": "Merge pull request #698 from vmware-tanzu/dependabot/github_actions/reviewdog/action-misspell-1.12.2\n\nBump reviewdog/action-misspell from 1.12.1 to 1.12.2",
          "timestamp": "2022-05-23T16:28:59-06:00",
          "tree_id": "c9cc8b5b55cf68aec99adf85b5b34abaaac6e1a5",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/bbdb44b376ee6e198aba8c2163524f309c0ddc9a"
        },
        "date": 1653345564429,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36619588204,
            "unit": "ns/op\t        19.06 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9607161798,
            "unit": "ns/op\t         5.379 DeleteSeconds\t         4.190 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6445591458,
            "unit": "ns/op\t         4.272 DeleteSeconds\t         2.135 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7080dc36c0bdcb62f78053024156bb8c37bc5063",
          "message": "[k8 1.24] Use TokenRequest API to get SA token (#695)\n\nCo-authored-by: Neil Hickey <nhickey@vmware.com>",
          "timestamp": "2022-05-23T19:24:07-04:00",
          "tree_id": "02ae0be7486938548e6928bb395c2c00901be1de",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7080dc36c0bdcb62f78053024156bb8c37bc5063"
        },
        "date": 1653349071877,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37430056409,
            "unit": "ns/op\t        19.54 DeleteSeconds\t        17.83 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9812601204,
            "unit": "ns/op\t         5.485 DeleteSeconds\t         4.254 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6779125244,
            "unit": "ns/op\t         4.507 DeleteSeconds\t         2.210 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "59fc2c1397ff448409ba9b5f36137247a4ad2492",
          "message": "Merge pull request #692 from vmware-tanzu/pkg-update-install\n\nkctrl: Remove `install` option in package installed update command",
          "timestamp": "2022-05-24T12:15:19+02:00",
          "tree_id": "e6dae66f8e8839cbc04ac69dc3882be35fe95fe3",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/59fc2c1397ff448409ba9b5f36137247a4ad2492"
        },
        "date": 1653388094221,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37144113037,
            "unit": "ns/op\t        19.39 DeleteSeconds\t        17.70 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9762874613,
            "unit": "ns/op\t         5.412 DeleteSeconds\t         4.296 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6671985463,
            "unit": "ns/op\t         4.397 DeleteSeconds\t         2.221 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fa18786913f6dff2af3cd9e97db7dde6fbbe45c9",
          "message": "support global kapp rawOptions and set default --app-changes-max-to-keep to 5 (#694)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-05-24T12:26:25-04:00",
          "tree_id": "240e8fe82b6bd4f3de57cc0df03d99040850df39",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/fa18786913f6dff2af3cd9e97db7dde6fbbe45c9"
        },
        "date": 1653410197288,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36584205540,
            "unit": "ns/op\t        19.02 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9507992873,
            "unit": "ns/op\t         5.314 DeleteSeconds\t         4.156 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6459052120,
            "unit": "ns/op\t         4.284 DeleteSeconds\t         2.136 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3cf77bdc169a9da947794cb0ca12323b66d7ef87",
          "message": "Merge pull request #703 from vmware-tanzu/bump-kapp-v0.48.0\n\nBump kapp to v0.48.0",
          "timestamp": "2022-05-26T11:23:35-06:00",
          "tree_id": "4a047fb144bf50c8fa20db1c3d96ae72fc8d1170",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3cf77bdc169a9da947794cb0ca12323b66d7ef87"
        },
        "date": 1653586426640,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36550437595,
            "unit": "ns/op\t        18.98 DeleteSeconds\t        17.53 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9506486106,
            "unit": "ns/op\t         5.304 DeleteSeconds\t         4.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6455390925,
            "unit": "ns/op\t         4.263 DeleteSeconds\t         2.148 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3beaddc443c6ec2ab2b1a0677684fb1c2d6d7f35",
          "message": "Merge pull request #711 from vmware-tanzu/kctrl-examples\n\nAdd support for positional arguments in package installed status command",
          "timestamp": "2022-05-31T10:31:56+05:30",
          "tree_id": "1aac1ab5401f225dc33828b58af41c018f5f624a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3beaddc443c6ec2ab2b1a0677684fb1c2d6d7f35"
        },
        "date": 1653974068288,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36944340327,
            "unit": "ns/op\t        19.21 DeleteSeconds\t        17.68 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9758356739,
            "unit": "ns/op\t         5.512 DeleteSeconds\t         4.195 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6600464357,
            "unit": "ns/op\t         4.330 DeleteSeconds\t         2.200 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d912f0484a9abeeddaf9e70e3617087c7bea691b",
          "message": "Merge pull request #683 from vmware-tanzu/kctrl-drop-values-flag\n\nAdd `--values` to `kctrl package installed update`",
          "timestamp": "2022-05-31T10:31:43+05:30",
          "tree_id": "03a0f108cd3e8031c5cbb4606f7ef56becee7cc2",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d912f0484a9abeeddaf9e70e3617087c7bea691b"
        },
        "date": 1653974086535,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37093245506,
            "unit": "ns/op\t        19.32 DeleteSeconds\t        17.72 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9720343596,
            "unit": "ns/op\t         5.406 DeleteSeconds\t         4.253 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6626451628,
            "unit": "ns/op\t         4.349 DeleteSeconds\t         2.224 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c2660e2fee78272342f63394431faeb4caba145b",
          "message": "Ensure that pkgi status tail is picked up after kicking it (#702)\n\n* Ensure that stale conditions in PackageInstalls copied from underlying AppCR are not picked up by kctrl\r\n* Poll for underlying App CR and ensure it matches the latest generation\r\n* Double poll interval so that we are not heavier on the api-server\r\n* These changes are to be reverted on resolution of https://github.com/vmware-tanzu/carvel-kapp-controller/issues/639\r\n\r\n* Ensure that underlying App CR is paused before triggering reconciliation\r\n\r\n* Make package installed kick tests stricter",
          "timestamp": "2022-05-31T12:24:01-04:00",
          "tree_id": "0db62e310bc59908c7d900772617ad023ca346db",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c2660e2fee78272342f63394431faeb4caba145b"
        },
        "date": 1654014857759,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36628078980,
            "unit": "ns/op\t        19.06 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9562094551,
            "unit": "ns/op\t         5.326 DeleteSeconds\t         4.194 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6453482320,
            "unit": "ns/op\t         4.269 DeleteSeconds\t         2.143 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b7f2002d5f48228aa89b979fd731744757c85846",
          "message": "Merge pull request #709 from vmware-tanzu/bump-ytt-v0.41.1\n\nBump ytt to v0.41.1",
          "timestamp": "2022-05-31T14:32:18-06:00",
          "tree_id": "3d8528ff2a14046ff1d9bf330da3604fb33b4bac",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/b7f2002d5f48228aa89b979fd731744757c85846"
        },
        "date": 1654029837250,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36889137552,
            "unit": "ns/op\t        19.18 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9653905684,
            "unit": "ns/op\t         5.403 DeleteSeconds\t         4.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6526080642,
            "unit": "ns/op\t         4.316 DeleteSeconds\t         2.161 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "72ecfa083349bc0a0a16dae8463df496baed6108",
          "message": "Merge pull request #708 from vmware-tanzu/nh-sa-token-review\n\nUse UID of ServiceAccount for token cache",
          "timestamp": "2022-05-31T14:31:41-06:00",
          "tree_id": "34b2903d6cb3614b4cf69d83b829dcb0781f0f2c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/72ecfa083349bc0a0a16dae8463df496baed6108"
        },
        "date": 1654029838946,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36821899081,
            "unit": "ns/op\t        19.19 DeleteSeconds\t        17.59 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9601859973,
            "unit": "ns/op\t         5.360 DeleteSeconds\t         4.199 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6518151470,
            "unit": "ns/op\t         4.299 DeleteSeconds\t         2.161 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3b4e5b00755bd098944b509df7d5a1b0500d3df5",
          "message": "Merge pull request #712 from vmware-tanzu/dependabot/github_actions/benchmark-action/github-action-benchmark-1.14.0\n\nBump benchmark-action/github-action-benchmark from 1.13.0 to 1.14.0",
          "timestamp": "2022-05-31T17:04:42-06:00",
          "tree_id": "a4a9cea33a18b87381f7694f98e2066b4e328c9b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3b4e5b00755bd098944b509df7d5a1b0500d3df5"
        },
        "date": 1654039124564,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37234552907,
            "unit": "ns/op\t        19.42 DeleteSeconds\t        17.76 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9775640087,
            "unit": "ns/op\t         5.516 DeleteSeconds\t         4.201 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6760887039,
            "unit": "ns/op\t         4.509 DeleteSeconds\t         2.194 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "772c1c29a1af0edad6a21066f570f82fb28a6c4e",
          "message": "update example/test for helm fetching (#716)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-06-02T15:23:03-04:00",
          "tree_id": "ad246315e734c8bf2acc7cd0a271c516ae4ec6e8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/772c1c29a1af0edad6a21066f570f82fb28a6c4e"
        },
        "date": 1654198656599,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37565505538,
            "unit": "ns/op\t        19.65 DeleteSeconds\t        17.84 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9832086572,
            "unit": "ns/op\t         5.495 DeleteSeconds\t         4.273 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6826489722,
            "unit": "ns/op\t         4.517 DeleteSeconds\t         2.247 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "526b10b6decd603de87605e16634df8c8f7d13b1",
          "message": "Ensure that kctrl picks up status tail after secrets are updated (#713)\n\n* Ensure that kctrl picks up status tail after secrets are updated. Stricter tests.\r\n\r\n* Add logging messages while pausing and resuming reconciliation",
          "timestamp": "2022-06-06T16:16:31+05:30",
          "tree_id": "99e0e5296e991b9b327ade98daf21a30912f62ab",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/526b10b6decd603de87605e16634df8c8f7d13b1"
        },
        "date": 1654513006819,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36604218523,
            "unit": "ns/op\t        19.05 DeleteSeconds\t        17.51 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9542257119,
            "unit": "ns/op\t         5.305 DeleteSeconds\t         4.196 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6508661667,
            "unit": "ns/op\t         4.294 DeleteSeconds\t         2.143 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "86e9eb13c1d178afe0cfa5b40450daa62983b698",
          "message": "Merge pull request #697 from vmware-tanzu/sidecarexec\n\nmove some App CR reconciliation parts into sidecar",
          "timestamp": "2022-06-06T05:58:40-07:00",
          "tree_id": "7a880513d798276bf8b07070773aa655783f28c3",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/86e9eb13c1d178afe0cfa5b40450daa62983b698"
        },
        "date": 1654521040893,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36474905557,
            "unit": "ns/op\t        18.97 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9439112708,
            "unit": "ns/op\t         5.278 DeleteSeconds\t         4.126 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6421867080,
            "unit": "ns/op\t         4.257 DeleteSeconds\t         2.121 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "811315f16b12f24487b30e2a593e27e051821158",
          "message": "Merge pull request #699 from vmware-tanzu/nh-upgrade-go-1.18\n\nUpgrade GoLang version to 1.18.x",
          "timestamp": "2022-06-06T14:04:12-04:00",
          "tree_id": "fdbc21e17d0e68b3c52f26cbfa425ef7a6b4cd4a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/811315f16b12f24487b30e2a593e27e051821158"
        },
        "date": 1654539408224,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37002223733,
            "unit": "ns/op\t        19.32 DeleteSeconds\t        17.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9635303948,
            "unit": "ns/op\t         5.413 DeleteSeconds\t         4.176 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6645973872,
            "unit": "ns/op\t         4.389 DeleteSeconds\t         2.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "bc807c4df615d953a39b911292ddbb3bd347f625",
          "message": "Merge pull request #718 from vmware-tanzu/dependabot/docker/golang-1.18.3\n\nBump golang from 1.17.9 to 1.18.3",
          "timestamp": "2022-06-06T16:35:57-04:00",
          "tree_id": "619bbd2a8062daba41a1effeb6733b793ac3287b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/bc807c4df615d953a39b911292ddbb3bd347f625"
        },
        "date": 1654548522361,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37196353381,
            "unit": "ns/op\t        19.37 DeleteSeconds\t        17.77 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9706622473,
            "unit": "ns/op\t         5.422 DeleteSeconds\t         4.220 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6687499450,
            "unit": "ns/op\t         4.403 DeleteSeconds\t         2.231 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e680c2e725c380b2584559051c72e087c20b0ff8",
          "message": "Cleanup before running package available test (#723)",
          "timestamp": "2022-06-07T15:03:52+05:30",
          "tree_id": "18246080fd14c0663c9794de64bf9cf82e5ab9b4",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/e680c2e725c380b2584559051c72e087c20b0ff8"
        },
        "date": 1654595087865,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36487112504,
            "unit": "ns/op\t        18.91 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9514791166,
            "unit": "ns/op\t         5.284 DeleteSeconds\t         4.193 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6402914909,
            "unit": "ns/op\t         4.239 DeleteSeconds\t         2.128 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "799672bc102662495e5ec43aac0c8aa2b880dbb8",
          "message": "Install/lock dependencies via config file (#721)\n\n* Install/lock dependencies via config file\r\n* fix kc binary names",
          "timestamp": "2022-06-07T10:22:21-04:00",
          "tree_id": "2ba7c458bdc1434b925d9ab2d6db905c115ba236",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/799672bc102662495e5ec43aac0c8aa2b880dbb8"
        },
        "date": 1654612457309,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36811947317,
            "unit": "ns/op\t        19.14 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9571864011,
            "unit": "ns/op\t         5.331 DeleteSeconds\t         4.192 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6531431323,
            "unit": "ns/op\t         4.311 DeleteSeconds\t         2.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fcd7c69ace7f08b5ab4bc9eb9a2f863d7f7c8dee",
          "message": "Merge pull request #719 from vmware-tanzu/dependabot/github_actions/peter-evans/create-pull-request-4.0.4\n\nBump peter-evans/create-pull-request from 4.0.3 to 4.0.4",
          "timestamp": "2022-06-07T10:47:48-06:00",
          "tree_id": "9d6d6da4d70a9784a9856a57354f3532c079dc30",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/fcd7c69ace7f08b5ab4bc9eb9a2f863d7f7c8dee"
        },
        "date": 1654621227864,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37038882140,
            "unit": "ns/op\t        19.24 DeleteSeconds\t        17.75 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9697316068,
            "unit": "ns/op\t         5.428 DeleteSeconds\t         4.219 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6586128344,
            "unit": "ns/op\t         4.366 DeleteSeconds\t         2.174 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "bd66239a7343d1dedb4d1fbe4dcf5e7968f40476",
          "message": "Merge pull request #727 from vmware-tanzu/bump-dependencies\n\nBump dependencies",
          "timestamp": "2022-06-07T18:02:23-04:00",
          "tree_id": "cf4d34afe94fdb1e71b526fe8ce8e270e1dfc935",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/bd66239a7343d1dedb4d1fbe4dcf5e7968f40476"
        },
        "date": 1654640144900,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37185516620,
            "unit": "ns/op\t        19.38 DeleteSeconds\t        17.75 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9762978045,
            "unit": "ns/op\t         5.446 DeleteSeconds\t         4.267 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6674678499,
            "unit": "ns/op\t         4.431 DeleteSeconds\t         2.188 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2fab826ade7077952eca4ebb6ec6b7520cd698dc",
          "message": "Merge pull request #704 from vmware-tanzu/dependabot/go_modules/k8s.io/kube-aggregator-0.22.10",
          "timestamp": "2022-06-08T09:07:12-04:00",
          "tree_id": "179bf9a566c27793adfef6f4afc297cc09b195ac",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2fab826ade7077952eca4ebb6ec6b7520cd698dc"
        },
        "date": 1654694412160,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37106661347,
            "unit": "ns/op\t        19.30 DeleteSeconds\t        17.73 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9712917532,
            "unit": "ns/op\t         5.408 DeleteSeconds\t         4.242 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6619403322,
            "unit": "ns/op\t         4.352 DeleteSeconds\t         2.199 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "34875791846ca5093ac5facd770d3ec3486ba947",
          "message": "pkgi waits for app to reconcile latest generation (#726)\n\nRight now we end up observing an old generation of the app's status.\r\nThis update means we will go back to the \"Reconciling\" state while\r\nwaiting for the app reconciler to sync to the latest desired state.\r\n\r\nObservedGeneration has this semi-helpful comment on it:\r\n> Populated based on metadata.generation when controller observes a change\r\n> to the resource; if this value is out of data, other status fields do\r\n> not reflect latest state",
          "timestamp": "2022-06-08T11:10:28-04:00",
          "tree_id": "af7b05a937ef7baa9113120e94a1822872a9a4ed",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/34875791846ca5093ac5facd770d3ec3486ba947"
        },
        "date": 1654701659577,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36423730968,
            "unit": "ns/op\t        18.91 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9548355774,
            "unit": "ns/op\t         5.332 DeleteSeconds\t         4.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6406396388,
            "unit": "ns/op\t         4.236 DeleteSeconds\t         2.133 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "davidblum@users.noreply.github.com",
            "name": "David Blum",
            "username": "davidblum"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a3a04afac3433abd121db2d6143743c8008da919",
          "message": "Add documentation for running kapp-controller in KIND (#724)\n\n* Add documentation for running kapp-controller in KIND\r\n\r\n* typo, update readme",
          "timestamp": "2022-06-08T12:05:38-04:00",
          "tree_id": "5a27326418447473abb5fac273104c50a649d301",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/a3a04afac3433abd121db2d6143743c8008da919"
        },
        "date": 1654705068454,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36741137656,
            "unit": "ns/op\t        19.08 DeleteSeconds\t        17.61 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9566259395,
            "unit": "ns/op\t         5.335 DeleteSeconds\t         4.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6524210189,
            "unit": "ns/op\t         4.330 DeleteSeconds\t         2.145 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d6c7d28734f061acdb9be2a1b88af7fa5601e6f9",
          "message": "rebase rule and e2e tests for pkgrs with identical pkgs (#657)\n\n* Allow identical Packages from different Repos\r\n\r\npkgr templating applies a rebase rule that inserts a noop\r\nannotation on a package coming from a PKGR in cases where\r\na package with identical name and contents is already provided\r\nby a different repo.\r\n\r\n- revision annotation allows changes to the package yaml without\r\n  changing the version\r\n- packages that are not identical will still fail to reconcile\r\n\r\n* allow for upgrades from old versions of kc at the expense of taking ownership of standalone packages",
          "timestamp": "2022-06-08T19:08:02-04:00",
          "tree_id": "09f1e40734c6d571f55b2f2d0f259e84e8eb70bc",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d6c7d28734f061acdb9be2a1b88af7fa5601e6f9"
        },
        "date": 1654730370136,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36838873935,
            "unit": "ns/op\t        19.19 DeleteSeconds\t        17.59 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9571353324,
            "unit": "ns/op\t         5.313 DeleteSeconds\t         4.207 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6548914547,
            "unit": "ns/op\t         4.313 DeleteSeconds\t         2.192 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7d8793060f6a9a1c913c847c770e21506ca78ddd",
          "message": "Merge pull request #730 from vmware-tanzu/bump-dependencies\n\nBump dependencies",
          "timestamp": "2022-06-09T10:48:35-04:00",
          "tree_id": "62161d650f778f90999b06a995a4347c3690b170",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7d8793060f6a9a1c913c847c770e21506ca78ddd"
        },
        "date": 1654786892819,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37057197148,
            "unit": "ns/op\t        19.31 DeleteSeconds\t        17.70 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9739129908,
            "unit": "ns/op\t         5.427 DeleteSeconds\t         4.260 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6660095347,
            "unit": "ns/op\t         4.406 DeleteSeconds\t         2.205 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d49800d7d5024058911e096cfd7487f49ccf0c5f",
          "message": "Bump dependencies (#731)",
          "timestamp": "2022-06-09T10:04:55-07:00",
          "tree_id": "f40d8d1d2bebf2effbfd4c147ec003ef084d3539",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d49800d7d5024058911e096cfd7487f49ccf0c5f"
        },
        "date": 1654794953195,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36626500337,
            "unit": "ns/op\t        19.03 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9519553424,
            "unit": "ns/op\t         5.310 DeleteSeconds\t         4.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6507459578,
            "unit": "ns/op\t         4.321 DeleteSeconds\t         2.145 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7023c865b0dafd52ca9485804d8720b9f74d39a1",
          "message": "Merge pull request #680 from benmoss/multi-arch\n\nAdd arm64 builds",
          "timestamp": "2022-06-09T13:28:32-07:00",
          "tree_id": "f1e603cce3f7e7368ccdac29b0c90f1dc9c4889a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7023c865b0dafd52ca9485804d8720b9f74d39a1"
        },
        "date": 1654807197730,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37011275053,
            "unit": "ns/op\t        19.27 DeleteSeconds\t        17.69 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9629351077,
            "unit": "ns/op\t         5.364 DeleteSeconds\t         4.217 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6559534100,
            "unit": "ns/op\t         4.302 DeleteSeconds\t         2.204 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0f13453db4957439e29519d61b2ffa10bd33d06b",
          "message": "Use the built-in minikube (#734)",
          "timestamp": "2022-06-09T16:32:00-04:00",
          "tree_id": "a932f543d033354aef0eb091851bdb447274859b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/0f13453db4957439e29519d61b2ffa10bd33d06b"
        },
        "date": 1654807397600,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37093229810,
            "unit": "ns/op\t        19.34 DeleteSeconds\t        17.69 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9742420332,
            "unit": "ns/op\t         5.389 DeleteSeconds\t         4.296 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6595219474,
            "unit": "ns/op\t         4.321 DeleteSeconds\t         2.187 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9c27bad3769bd86fa3409ee7c8ca0baa867d9226",
          "message": "correct automated usage of kbld in PackageRepository (when imgpkgBundle is specified) (#737)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-06-10T14:42:55-04:00",
          "tree_id": "eb40697eaafe504c2236b83288bc9c0d89a35cae",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/9c27bad3769bd86fa3409ee7c8ca0baa867d9226"
        },
        "date": 1654887206874,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36643157783,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.57 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9538181732,
            "unit": "ns/op\t         5.339 DeleteSeconds\t         4.158 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6565298856,
            "unit": "ns/op\t         4.371 DeleteSeconds\t         2.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f43b217d44eb9a5468d4502ae5c5e27338ac812a",
          "message": "Fix trivy scan (#741)",
          "timestamp": "2022-06-13T10:29:29-04:00",
          "tree_id": "118a30e0676589c453e5ac533ddf58f64093b6cc",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f43b217d44eb9a5468d4502ae5c5e27338ac812a"
        },
        "date": 1655131136722,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36497949203,
            "unit": "ns/op\t        18.96 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9505756162,
            "unit": "ns/op\t         5.310 DeleteSeconds\t         4.158 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6473463731,
            "unit": "ns/op\t         4.269 DeleteSeconds\t         2.165 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1ece87d94d99fd0327acc70abb3e71fd2682fb5a",
          "message": "Merge pull request #715 from vmware-tanzu/kctrl-ux-enhancements-1\n\nBunch of UX fixes tweaking output",
          "timestamp": "2022-06-13T07:28:58-07:00",
          "tree_id": "3317d3ee4228cddee5b3b76b287142640166237e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/1ece87d94d99fd0327acc70abb3e71fd2682fb5a"
        },
        "date": 1655131141739,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36600010638,
            "unit": "ns/op\t        18.99 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9590178046,
            "unit": "ns/op\t         5.382 DeleteSeconds\t         4.167 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6507111710,
            "unit": "ns/op\t         4.291 DeleteSeconds\t         2.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d10acb8168331c6a724cf6c6060c26bb3a3cebe1",
          "message": "kind action: bump to 1.3 (#742)",
          "timestamp": "2022-06-14T14:15:50-04:00",
          "tree_id": "49bc77a73ce3a09004c9650206dc2979125039f2",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d10acb8168331c6a724cf6c6060c26bb3a3cebe1"
        },
        "date": 1655231114055,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36438688912,
            "unit": "ns/op\t        18.92 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9512965251,
            "unit": "ns/op\t         5.308 DeleteSeconds\t         4.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6517477110,
            "unit": "ns/op\t         4.328 DeleteSeconds\t         2.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4dd0ab3e6d10a7c6a967437e0a4ee691d06abea4",
          "message": "Merge pull request #738 from vmware-tanzu/e2e-controller-config-scoped",
          "timestamp": "2022-06-15T09:08:32-04:00",
          "tree_id": "86f6223d51fa8769c0fad217d7fc2cbe1242a6e0",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4dd0ab3e6d10a7c6a967437e0a4ee691d06abea4"
        },
        "date": 1655299095136,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36660584865,
            "unit": "ns/op\t        19.01 DeleteSeconds\t        17.61 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9536669213,
            "unit": "ns/op\t         5.329 DeleteSeconds\t         4.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6521999267,
            "unit": "ns/op\t         4.351 DeleteSeconds\t         2.134 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4532073f483c0cfb5a7d331a9eae2e87bf65add6",
          "message": "Enhance tty experience (#743)\n\n* Set tty to be an alternate flag instead of global\r\n\r\nSet default value of tty to true for add/update, delete, pause, kick and status commands. For rest of the commands, default value should be false.\r\n\r\n* Make cmd configuration functions private",
          "timestamp": "2022-06-15T10:20:32-04:00",
          "tree_id": "70bbc9a72e25b15864c0df6e2e3e7268681375bc",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4532073f483c0cfb5a7d331a9eae2e87bf65add6"
        },
        "date": 1655303423608,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36459723420,
            "unit": "ns/op\t        18.94 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9520327856,
            "unit": "ns/op\t         5.281 DeleteSeconds\t         4.198 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6478436406,
            "unit": "ns/op\t         4.308 DeleteSeconds\t         2.127 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "08910dda8875ec83fe81b44706b21f50cb44a955",
          "message": "avoid having separate binary just for sidecarexec (#747)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-06-16T12:51:32-04:00",
          "tree_id": "06d18a2a903b354c4a1bfad2c9ec58b6e6276ebe",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/08910dda8875ec83fe81b44706b21f50cb44a955"
        },
        "date": 1655398951379,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36766415402,
            "unit": "ns/op\t        19.14 DeleteSeconds\t        17.57 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9626711041,
            "unit": "ns/op\t         5.359 DeleteSeconds\t         4.223 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6578137485,
            "unit": "ns/op\t         4.369 DeleteSeconds\t         2.165 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c414ff7f75791e69845688d227dfb844ddfb703a",
          "message": "PackageInstall: check if semver constraints are nil so we error instead of panic (#745)",
          "timestamp": "2022-06-16T13:35:40-04:00",
          "tree_id": "b2867a53d0d428bf231f20c52436ca33f7823da7",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c414ff7f75791e69845688d227dfb844ddfb703a"
        },
        "date": 1655401525836,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36617763172,
            "unit": "ns/op\t        19.04 DeleteSeconds\t        17.53 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9559536053,
            "unit": "ns/op\t         5.317 DeleteSeconds\t         4.200 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6462957516,
            "unit": "ns/op\t         4.277 DeleteSeconds\t         2.137 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "eea919aaaf205916548039168bd284e402e82087",
          "message": "support certificate reloading even when /etc is on a different mount from /tmp (#751)\n\n* improve TestConfig_TrustCACerts as it may flake when service->deployment is racing to completion\r\n\r\n* create tmp certs bundle in certs directory\r\n\r\nif created in /tmp, rename call may fail since /tmp and /etc are not guaranteed to be from the same mount\r\n\r\n* Test_PackageInstalled_FromPackageInstall_DeletionFailureBlocks: some prints and logic to hopefully help see why it's so flakey\r\n\r\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>\r\nCo-authored-by: Joe Kimmel <jkimmel@vmware.com>",
          "timestamp": "2022-06-17T11:40:37-07:00",
          "tree_id": "77f6f561231227f908d7a449b13b119e684135a5",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/eea919aaaf205916548039168bd284e402e82087"
        },
        "date": 1655491808550,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36515466638,
            "unit": "ns/op\t        18.98 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9516154693,
            "unit": "ns/op\t         5.301 DeleteSeconds\t         4.177 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6433148742,
            "unit": "ns/op\t         4.251 DeleteSeconds\t         2.139 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e5f0471e720e30dd1ce8464a5667040085a0d5e1",
          "message": "Merge pull request #746 from benmoss/fix-dev-deploy",
          "timestamp": "2022-06-17T15:00:26-04:00",
          "tree_id": "b5862695d45c0dcde21b90968128fa80102a1927",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/e5f0471e720e30dd1ce8464a5667040085a0d5e1"
        },
        "date": 1655493051719,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36628229030,
            "unit": "ns/op\t        18.99 DeleteSeconds\t        17.59 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9565566561,
            "unit": "ns/op\t         5.333 DeleteSeconds\t         4.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6533653390,
            "unit": "ns/op\t         4.331 DeleteSeconds\t         2.159 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c2fe3ced170ad3070a07469875a2efab7daa41b8",
          "message": "Merge pull request #754 from benmoss/fix-dev-deploy",
          "timestamp": "2022-06-17T15:50:20-04:00",
          "tree_id": "8194f52d8539d20d11a9c23c0d3f63add2ceac02",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c2fe3ced170ad3070a07469875a2efab7daa41b8"
        },
        "date": 1655495989974,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36471168153,
            "unit": "ns/op\t        18.96 DeleteSeconds\t        17.46 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9533343063,
            "unit": "ns/op\t         5.329 DeleteSeconds\t         4.165 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6417027465,
            "unit": "ns/op\t         4.250 DeleteSeconds\t         2.129 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1d229748b11f6fa20c4e7ae177bc6c69c910d06f",
          "message": "hack/test-e2e: no more errors on unset variables (#752)",
          "timestamp": "2022-06-21T10:58:01-04:00",
          "tree_id": "d8499b2c211a4031e8da2f9fb7951017a0d02fd5",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/1d229748b11f6fa20c4e7ae177bc6c69c910d06f"
        },
        "date": 1655824037739,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36445692818,
            "unit": "ns/op\t        18.87 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9495432644,
            "unit": "ns/op\t         5.304 DeleteSeconds\t         4.151 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6489881281,
            "unit": "ns/op\t         4.309 DeleteSeconds\t         2.143 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ddfb811901428f0bc459ecf45432dbaf2e88067c",
          "message": "correct flaky Test_PackageInstallAndRepo_CanAuthenticateToPrivateRepository_UsingPlaceholderSecret test (#758)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-06-22T07:06:12-04:00",
          "tree_id": "6c48041b74b6904100ac7ff25cc578a28223148b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ddfb811901428f0bc459ecf45432dbaf2e88067c"
        },
        "date": 1655896532111,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36407228046,
            "unit": "ns/op\t        18.90 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9545516398,
            "unit": "ns/op\t         5.308 DeleteSeconds\t         4.199 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6425505803,
            "unit": "ns/op\t         4.245 DeleteSeconds\t         2.138 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9bd165d7b68c48be3ca1ed0b1b137fc0ed5d41ac",
          "message": "Fix flakiness of private registry auth test (#759)\n\nvendir doesn't retry, so we need to make sure that the registry service\r\nis up and responsive before we deploy anything to kapp-controller",
          "timestamp": "2022-06-23T12:36:42-04:00",
          "tree_id": "b1dfd7182ff2f6c3953422af8a0441e18f69ff4e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/9bd165d7b68c48be3ca1ed0b1b137fc0ed5d41ac"
        },
        "date": 1656002781905,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36537466427,
            "unit": "ns/op\t        18.97 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9553719992,
            "unit": "ns/op\t         5.322 DeleteSeconds\t         4.192 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6535293482,
            "unit": "ns/op\t         4.346 DeleteSeconds\t         2.140 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7dd8177ca515604464541b4c3d8b9fd89d33aacf",
          "message": "Merge pull request #744 from vmware-tanzu/659-add-downward-api\n\nadd downward api",
          "timestamp": "2022-06-23T12:39:01-04:00",
          "tree_id": "a9c927d9ff1042a404a3b45a03aef1f5a81f23e4",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7dd8177ca515604464541b4c3d8b9fd89d33aacf"
        },
        "date": 1656002972844,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36755703529,
            "unit": "ns/op\t        19.11 DeleteSeconds\t        17.60 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9654480530,
            "unit": "ns/op\t         5.429 DeleteSeconds\t         4.181 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6587056084,
            "unit": "ns/op\t         4.376 DeleteSeconds\t         2.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8e67cd7f576179dbe7681dabc92f624c5a243ce0",
          "message": "Update release workflow to have ./ before artefacts (#761)\n\nThis format is required by the carvel-release-scripts and carvel-setup-action",
          "timestamp": "2022-06-27T07:52:21-04:00",
          "tree_id": "eb04e4c34f2179e78331d8d57326375d20782d0f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8e67cd7f576179dbe7681dabc92f624c5a243ce0"
        },
        "date": 1656331435513,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37009780783,
            "unit": "ns/op\t        19.26 DeleteSeconds\t        17.70 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9775177904,
            "unit": "ns/op\t         5.509 DeleteSeconds\t         4.217 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6726261420,
            "unit": "ns/op\t         4.455 DeleteSeconds\t         2.219 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "5e13daa01f0352c3a2e63719c93f047ce4f9ec2c",
          "message": "Add formatting for checksums in draft release body (#767)",
          "timestamp": "2022-06-28T14:48:17-04:00",
          "tree_id": "3c4039778d17f78ac5b9a1bf8f2b277077ff6e8e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/5e13daa01f0352c3a2e63719c93f047ce4f9ec2c"
        },
        "date": 1656442679149,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36339504691,
            "unit": "ns/op\t        18.83 DeleteSeconds\t        17.46 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9589238809,
            "unit": "ns/op\t         5.345 DeleteSeconds\t         4.211 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6391121525,
            "unit": "ns/op\t         4.237 DeleteSeconds\t         2.120 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "bae14f114fa22e6412d56e4dc9d833c261ad3f12",
          "message": "Merge pull request #764 from vmware-tanzu/kctrl-repo-tail\n\nAdd tailing behaviour to package repo and add a package repo kick command",
          "timestamp": "2022-06-28T15:47:28-04:00",
          "tree_id": "eead2b5808e6bd9ca3edfa0758ae53e6aa43e864",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/bae14f114fa22e6412d56e4dc9d833c261ad3f12"
        },
        "date": 1656446349923,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37553097759,
            "unit": "ns/op\t        19.43 DeleteSeconds\t        18.09 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9687594231,
            "unit": "ns/op\t         5.391 DeleteSeconds\t         4.265 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7762077293,
            "unit": "ns/op\t         4.298 DeleteSeconds\t         3.428 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "97f07296e191857493a9f5b6e0b61bc55f36f5c6",
          "message": "Merge pull request #762 from vmware-tanzu/dependabot/go_modules/github.com/stretchr/testify-1.7.5\n\nBump github.com/stretchr/testify from 1.7.1 to 1.7.5",
          "timestamp": "2022-06-28T16:22:15-06:00",
          "tree_id": "14a0bcac7af449780ad151997c8d3cc4ec965184",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/97f07296e191857493a9f5b6e0b61bc55f36f5c6"
        },
        "date": 1656455677124,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37292059308,
            "unit": "ns/op\t        19.55 DeleteSeconds\t        17.68 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9660506541,
            "unit": "ns/op\t         5.393 DeleteSeconds\t         4.200 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6621589413,
            "unit": "ns/op\t         4.377 DeleteSeconds\t         2.196 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c5de7b74cc34eded1d1a0784ed2f8767079dc77a",
          "message": "Merge pull request #750 from vmware-tanzu/dependabot/go_modules/k8s.io/kube-aggregator-0.22.11\n\nBump k8s.io/kube-aggregator from 0.22.10 to 0.22.11",
          "timestamp": "2022-06-28T16:54:04-06:00",
          "tree_id": "a4820b7901064d66d06bc42bc15df46d6fbb6fd7",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c5de7b74cc34eded1d1a0784ed2f8767079dc77a"
        },
        "date": 1656457566012,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37032977901,
            "unit": "ns/op\t        19.24 DeleteSeconds\t        17.75 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9658547344,
            "unit": "ns/op\t         5.408 DeleteSeconds\t         4.203 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6579643413,
            "unit": "ns/op\t         4.369 DeleteSeconds\t         2.161 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8a2c65438b5b02bcf66af789492d3efb026524f3",
          "message": "Merge pull request #735 from vmware-tanzu/dependabot/go_modules/golang.org/x/tools-0.1.11\n\nBump golang.org/x/tools from 0.1.10 to 0.1.11",
          "timestamp": "2022-06-29T09:25:46-06:00",
          "tree_id": "fc4657a292418e506d1f420e8d71ab43ae72e000",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8a2c65438b5b02bcf66af789492d3efb026524f3"
        },
        "date": 1656517438334,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36686896875,
            "unit": "ns/op\t        19.02 DeleteSeconds\t        17.60 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9502522484,
            "unit": "ns/op\t         5.287 DeleteSeconds\t         4.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6449011607,
            "unit": "ns/op\t         4.264 DeleteSeconds\t         2.146 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "58ce5c8bd33e9b930e0d58a5299b2d563b24c0dc",
          "message": "Bump dependencies (#769)",
          "timestamp": "2022-06-30T15:16:02-04:00",
          "tree_id": "6a1ac657fbb96cd07eefd0ffe25861c37f613db8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/58ce5c8bd33e9b930e0d58a5299b2d563b24c0dc"
        },
        "date": 1656617163026,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36643163977,
            "unit": "ns/op\t        19.09 DeleteSeconds\t        17.51 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9494006870,
            "unit": "ns/op\t         5.296 DeleteSeconds\t         4.161 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6408793696,
            "unit": "ns/op\t         4.245 DeleteSeconds\t         2.125 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7cd7d43cb545028ef91f8da390711dc270ea6b0a",
          "message": "rename kctrl github test action to specify that it's for kctrl (#773)",
          "timestamp": "2022-07-05T14:03:48-04:00",
          "tree_id": "1723c72bf46f88780ec7b42141f6c19448b9413b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7cd7d43cb545028ef91f8da390711dc270ea6b0a"
        },
        "date": 1657044797964,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36458194774,
            "unit": "ns/op\t        18.92 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9460359733,
            "unit": "ns/op\t         5.264 DeleteSeconds\t         4.159 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6414619650,
            "unit": "ns/op\t         4.247 DeleteSeconds\t         2.125 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f206bc98b01b78bcb135e1c4b811c3e6e79bd0c9",
          "message": "use cache mount in Dockerfile (#748)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-07-05T14:06:20-04:00",
          "tree_id": "7bd2bb67a80ded8f2e369b4568026b30736ec0a2",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f206bc98b01b78bcb135e1c4b811c3e6e79bd0c9"
        },
        "date": 1657045080470,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36960265042,
            "unit": "ns/op\t        19.26 DeleteSeconds\t        17.65 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9598462168,
            "unit": "ns/op\t         5.346 DeleteSeconds\t         4.208 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6537059271,
            "unit": "ns/op\t         4.315 DeleteSeconds\t         2.172 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "5f167fc1de1861cb546109dd343b67e360279eb2",
          "message": "Rename KC owned apps from `x-ctrl` to `x.app` or `x.pkgr` (#665)\n\n* Add support for `--prev-app` on kapp deploy/delete\r\n\r\n- PackageRepo will be suffixed with .pkgr\r\n- Apps will be suffixed with .app (pkgi, appcr)\r\n- Replace hardcoded -ctrl in tests with .app\r\n- Add e2e test for migration\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\n* Move comments to logger.Section\r\n\r\n- Make config variable longer to adhere to GoLang practices\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>",
          "timestamp": "2022-07-06T10:53:08-04:00",
          "tree_id": "85bee40aa5dc0feb050ddf639c62162c2f6bd5f5",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/5f167fc1de1861cb546109dd343b67e360279eb2"
        },
        "date": 1657119775029,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36699595669,
            "unit": "ns/op\t        18.99 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9472875487,
            "unit": "ns/op\t         5.263 DeleteSeconds\t         4.171 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6454494121,
            "unit": "ns/op\t         4.254 DeleteSeconds\t         2.147 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a228a48374c812588cb5854cf17c5d994bdb373d",
          "message": "Merge pull request #775 from vmware-tanzu/pkg-repo-sidecar\n\nexecute pkg repo fetching in the sidecar",
          "timestamp": "2022-07-07T09:43:16-04:00",
          "tree_id": "3794712ca0db86a787eb2b1b9a26038809a3069a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/a228a48374c812588cb5854cf17c5d994bdb373d"
        },
        "date": 1657202002978,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36613660147,
            "unit": "ns/op\t        19.06 DeleteSeconds\t        17.51 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9699032290,
            "unit": "ns/op\t         5.429 DeleteSeconds\t         4.176 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6451443488,
            "unit": "ns/op\t         4.262 DeleteSeconds\t         2.148 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ae8e99da4ac8c8c168921f6dc0c190bc98cf4f70",
          "message": "Disallow use of shared namespaces for package installs (#757)",
          "timestamp": "2022-07-07T11:39:20-04:00",
          "tree_id": "45d00b51b512ea3343300ff1f64ada80811da32c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ae8e99da4ac8c8c168921f6dc0c190bc98cf4f70"
        },
        "date": 1657208916251,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36355101082,
            "unit": "ns/op\t        18.85 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9528264892,
            "unit": "ns/op\t         5.343 DeleteSeconds\t         4.149 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6510815687,
            "unit": "ns/op\t         4.323 DeleteSeconds\t         2.126 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d17c2bdebdd7852a0173f7bc70c2e58253e4999a",
          "message": "Bump dependencies (#778)",
          "timestamp": "2022-07-08T11:41:47-04:00",
          "tree_id": "f4b66292da00104af6bad43e9a76a6d8d0dc27a6",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d17c2bdebdd7852a0173f7bc70c2e58253e4999a"
        },
        "date": 1657295504297,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36562995942,
            "unit": "ns/op\t        19.06 DeleteSeconds\t        17.46 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9463187363,
            "unit": "ns/op\t         5.271 DeleteSeconds\t         4.154 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6431695470,
            "unit": "ns/op\t         4.254 DeleteSeconds\t         2.142 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8766a1d4915982a3ba5e152faf97983b8f9da28e",
          "message": "Remove dep on go-containerregistry (#779)\n\nThis library brings in a ton of transitive deps unfortunately, which is\r\nannoying for spurious CVE reports. The logic we were relying on from the\r\nlibrary is incredibly tiny.",
          "timestamp": "2022-07-11T11:25:55-04:00",
          "tree_id": "65de180b32165c326dd2178c5b137aaa721d1f81",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8766a1d4915982a3ba5e152faf97983b8f9da28e"
        },
        "date": 1657554078819,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 39047504052,
            "unit": "ns/op\t        19.80 DeleteSeconds\t        19.13 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9798034219,
            "unit": "ns/op\t         5.440 DeleteSeconds\t         4.304 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6679215314,
            "unit": "ns/op\t         4.421 DeleteSeconds\t         2.205 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8409bab24e7337ac3fd7527edd8fdb03f58b0c59",
          "message": "Remove sed (#780)",
          "timestamp": "2022-07-11T19:07:29-04:00",
          "tree_id": "53315d7c1e9b8c17bf1970b77196571d597f1dd0",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8409bab24e7337ac3fd7527edd8fdb03f58b0c59"
        },
        "date": 1657581421204,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36810889751,
            "unit": "ns/op\t        19.22 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9524039481,
            "unit": "ns/op\t         5.326 DeleteSeconds\t         4.158 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6430161606,
            "unit": "ns/op\t         4.252 DeleteSeconds\t         2.135 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ce55e80d878a2b7ac252059f22c4514d950411c2",
          "message": "Fix flaky test (#783)\n\nTest_PackageInstall_UsesExistingAppWithSameName would flake because the app\r\nreconciler is adding finalizers to the app at the same time the PKGI reconciler\r\nwants to update the app.\r\n\r\nIf the PKGI reconciler sees the app before the app reconciler finishes, there's\r\na race for who updates it first. If the PKGI reconciler loses that race, it\r\ngets angry and puts a failure status message on the PKGI, and that fails the\r\n`kapp deploy`, which fails the test.\r\n\r\nBy just waiting for the app to get ReconcileSucceeded, we can guarantee that\r\nthe app reconciler has finished and avoid the race. Ran it 100 times and it\r\ndidn't flake.",
          "timestamp": "2022-07-13T13:11:58-04:00",
          "tree_id": "096c9e1bf15d35f271afc98aa4b87a7c44d4dca2",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ce55e80d878a2b7ac252059f22c4514d950411c2"
        },
        "date": 1657733041582,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38032391381,
            "unit": "ns/op\t        19.20 DeleteSeconds\t        18.78 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9831702376,
            "unit": "ns/op\t         5.542 DeleteSeconds\t         4.230 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6591159874,
            "unit": "ns/op\t         4.340 DeleteSeconds\t         2.204 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d286a0f02a9942e19f3311ea37bcb757eb6b74ca",
          "message": "Add release-published workflow (#782)",
          "timestamp": "2022-07-13T13:13:25-04:00",
          "tree_id": "541987ea52d3701bbed052413b27e867f32b4ab2",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d286a0f02a9942e19f3311ea37bcb757eb6b74ca"
        },
        "date": 1657733086606,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36984619285,
            "unit": "ns/op\t        19.35 DeleteSeconds\t        17.58 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9643714192,
            "unit": "ns/op\t         5.371 DeleteSeconds\t         4.205 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6741640793,
            "unit": "ns/op\t         4.493 DeleteSeconds\t         2.199 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7877611779a1c62ba9773173962fc8c0038974d9",
          "message": "Bump dependencies (#787)",
          "timestamp": "2022-07-15T11:57:08-04:00",
          "tree_id": "f72b100daf4bb26af370f3ed8df7c2a9e8338021",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7877611779a1c62ba9773173962fc8c0038974d9"
        },
        "date": 1657901190693,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36340535142,
            "unit": "ns/op\t        18.84 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9542809179,
            "unit": "ns/op\t         5.353 DeleteSeconds\t         4.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6489411680,
            "unit": "ns/op\t         4.292 DeleteSeconds\t         2.139 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cc0bf278d76d9da33f8df96afb2bb9aa5059c3d0",
          "message": "Merge pull request #786 from vmware-tanzu/dependabot/docker/golang-1.18.4\n\nBump golang from 1.18.3 to 1.18.4",
          "timestamp": "2022-07-15T14:53:48-06:00",
          "tree_id": "d51f8edb1e6a1c4186a744cfec80ede0ca258d46",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cc0bf278d76d9da33f8df96afb2bb9aa5059c3d0"
        },
        "date": 1657919006711,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36500256839,
            "unit": "ns/op\t        18.94 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9471095640,
            "unit": "ns/op\t         5.288 DeleteSeconds\t         4.146 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6452427533,
            "unit": "ns/op\t         4.239 DeleteSeconds\t         2.176 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3fdf3b78d1d9989824b353c4eacc5551bf96e64b",
          "message": "Merge pull request #794 from vmware-tanzu/dependabot/github_actions/actions/stale-5.1.0\n\nBump actions/stale from 5.0.0 to 5.1.0",
          "timestamp": "2022-07-19T15:41:28-06:00",
          "tree_id": "332cfff5452c3e47e3d5ed2b993b99e850f7cc6a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3fdf3b78d1d9989824b353c4eacc5551bf96e64b"
        },
        "date": 1658267581813,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37249717067,
            "unit": "ns/op\t        19.47 DeleteSeconds\t        17.73 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9595482019,
            "unit": "ns/op\t         5.341 DeleteSeconds\t         4.200 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6646514712,
            "unit": "ns/op\t         4.399 DeleteSeconds\t         2.197 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "ahurley@vmware.com",
            "name": "Aaron Hurley",
            "username": "aaronshurley"
          },
          "committer": {
            "email": "ahurley@vmware.com",
            "name": "Aaron Hurley",
            "username": "aaronshurley"
          },
          "distinct": true,
          "id": "7821b4fbe84d48a87ae85877a477d2b0915e1390",
          "message": "Add workflow to add new issues and prs to project",
          "timestamp": "2022-07-26T16:43:14-07:00",
          "tree_id": "8989930ac9fd78ed1a7ffe686f37d16db522f6e6",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7821b4fbe84d48a87ae85877a477d2b0915e1390"
        },
        "date": 1658879625165,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37733003233,
            "unit": "ns/op\t        19.09 DeleteSeconds\t        18.59 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9541254002,
            "unit": "ns/op\t         5.326 DeleteSeconds\t         4.175 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6471745258,
            "unit": "ns/op\t         4.278 DeleteSeconds\t         2.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1c921ec806409c0e158ccd683e9408061e2a8b9c",
          "message": "Merge pull request #804 from slapula/multiarch-build-fix",
          "timestamp": "2022-08-01T10:27:58-04:00",
          "tree_id": "0527bb7a47ded0512e2b5f305cd8b49158506d8c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/1c921ec806409c0e158ccd683e9408061e2a8b9c"
        },
        "date": 1659364640149,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36516961899,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9514580782,
            "unit": "ns/op\t         5.332 DeleteSeconds\t         4.145 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6410404359,
            "unit": "ns/op\t         4.240 DeleteSeconds\t         2.131 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "aaronshurley@users.noreply.github.com",
            "name": "Aaron Hurley",
            "username": "aaronshurley"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "359cb36e89494f0d9a5d85ab4102f676e678da15",
          "message": "Merge pull request #808 from benmoss/fix-add-to-project\n\nuse pull_request_target to allow access to secrets",
          "timestamp": "2022-08-01T11:24:48-07:00",
          "tree_id": "b911af4836f8a439f558ad0ea86f4687fad4b2fb",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/359cb36e89494f0d9a5d85ab4102f676e678da15"
        },
        "date": 1659378969554,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37317424403,
            "unit": "ns/op\t        19.47 DeleteSeconds\t        17.79 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9590749691,
            "unit": "ns/op\t         5.352 DeleteSeconds\t         4.186 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6524362118,
            "unit": "ns/op\t         4.300 DeleteSeconds\t         2.174 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "068c8fd2b0715514d0b24032a1927a6c1bfad7e6",
          "message": "Merge pull request #793 from vmware-tanzu/bump-dependencies",
          "timestamp": "2022-08-01T16:27:51-04:00",
          "tree_id": "4511a1f1af44f30a0afeed3c4f736bd9b47a9874",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/068c8fd2b0715514d0b24032a1927a6c1bfad7e6"
        },
        "date": 1659386251215,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36557403189,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.51 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9598199061,
            "unit": "ns/op\t         5.382 DeleteSeconds\t         4.181 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6420960448,
            "unit": "ns/op\t         4.251 DeleteSeconds\t         2.124 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "5a1e9bb7b7dbd787395a574133ce697197ab1dcf",
          "message": "Bump dependencies (#817)",
          "timestamp": "2022-08-03T10:02:26-04:00",
          "tree_id": "bbf873a91c7a642c25f2f23be563a7aca4a311fe",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/5a1e9bb7b7dbd787395a574133ce697197ab1dcf"
        },
        "date": 1659536049028,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36960064957,
            "unit": "ns/op\t        19.25 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9584309252,
            "unit": "ns/op\t         5.354 DeleteSeconds\t         4.182 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6514667564,
            "unit": "ns/op\t         4.323 DeleteSeconds\t         2.140 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "benm@vmware.com",
            "name": "Ben Moss",
            "username": "benmoss"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "699e59824c2eff853557939a21c091423820477e",
          "message": "Check if pkgrs contain a packages directory (#818)\n\n* Check if pkgrs contain a packages directory\r\n\r\nFail more gracefully than ytt does if called on a non-existent directory",
          "timestamp": "2022-08-04T11:04:54-04:00",
          "tree_id": "9f4466a7e073fdea41ab97c18db5d1e7e047c0bc",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/699e59824c2eff853557939a21c091423820477e"
        },
        "date": 1659626168622,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36888881717,
            "unit": "ns/op\t        19.21 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9596203614,
            "unit": "ns/op\t         5.360 DeleteSeconds\t         4.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6640011760,
            "unit": "ns/op\t         4.327 DeleteSeconds\t         2.262 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4b468248fb8a1da6a7d7417228e0cb7ce41fd924",
          "message": "Bump actions/stale from 5.1.0 to 5.1.1 (#814)\n\nBumps [actions/stale](https://github.com/actions/stale) from 5.1.0 to 5.1.1.\r\n- [Release notes](https://github.com/actions/stale/releases)\r\n- [Changelog](https://github.com/actions/stale/blob/main/CHANGELOG.md)\r\n- [Commits](https://github.com/actions/stale/compare/532554b8a8498a0e006fbcde824b048728c4178f...9c1b1c6e115ca2af09755448e0dbba24e5061cc8)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: actions/stale\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-08-05T09:33:56-07:00",
          "tree_id": "1f180fdfe1273fb23c29e6620df05f4ad5bbf1be",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4b468248fb8a1da6a7d7417228e0cb7ce41fd924"
        },
        "date": 1659717800844,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36556349503,
            "unit": "ns/op\t        18.97 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9562907970,
            "unit": "ns/op\t         5.352 DeleteSeconds\t         4.156 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6460268431,
            "unit": "ns/op\t         4.244 DeleteSeconds\t         2.136 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "36f7a2a2cdbfbc0bf37ba7d027853afb5fbbbf10",
          "message": "Read default ca cert data from os env KAPPCTRL_KUBERNETES_CA_DATA (#819)\n\nThis allows kctrl to inject CA data into the reconciler when dev deploy\r\nruns it locally to mimic the controller.",
          "timestamp": "2022-08-05T12:39:53-04:00",
          "tree_id": "2e09585dcc84de78cb5cb1ee34fbd8b38f71c130",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/36f7a2a2cdbfbc0bf37ba7d027853afb5fbbbf10"
        },
        "date": 1659718423365,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36904449942,
            "unit": "ns/op\t        19.24 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9664373773,
            "unit": "ns/op\t         5.378 DeleteSeconds\t         4.247 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6548997710,
            "unit": "ns/op\t         4.326 DeleteSeconds\t         2.182 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ea5d78ef2d666a6d45fd2a58387fdc119c558be1",
          "message": "Bump dependencies (#820)",
          "timestamp": "2022-08-08T11:56:25-04:00",
          "tree_id": "7961fc13a60f79daa18898bc5fcdc46c0711c4d0",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ea5d78ef2d666a6d45fd2a58387fdc119c558be1"
        },
        "date": 1659974776529,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36622878198,
            "unit": "ns/op\t        19.08 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9472934481,
            "unit": "ns/op\t         5.265 DeleteSeconds\t         4.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6430921572,
            "unit": "ns/op\t         4.249 DeleteSeconds\t         2.144 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e58623822d99e95a00e2dea7b507b05d7b954f8f",
          "message": "Bump slackapi/slack-github-action from 1.19.0 to 1.21.0 (#813)\n\nBumps [slackapi/slack-github-action](https://github.com/slackapi/slack-github-action) from 1.19.0 to 1.21.0.\r\n- [Release notes](https://github.com/slackapi/slack-github-action/releases)\r\n- [Commits](https://github.com/slackapi/slack-github-action/compare/v1.19.0...v1.21.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: slackapi/slack-github-action\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-08-08T11:43:47-07:00",
          "tree_id": "9f3710b9690e670b11600a2e01f5fd350772131b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/e58623822d99e95a00e2dea7b507b05d7b954f8f"
        },
        "date": 1659984830784,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36758553354,
            "unit": "ns/op\t        19.03 DeleteSeconds\t        17.68 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9508029475,
            "unit": "ns/op\t         5.304 DeleteSeconds\t         4.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6564043885,
            "unit": "ns/op\t         4.307 DeleteSeconds\t         2.152 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "13bcb05a1953b1afbb7a696163c27876c18071e5",
          "message": "Bump golang.org/x/tools from 0.1.11 to 0.1.12 (#800)\n\nBumps [golang.org/x/tools](https://github.com/golang/tools) from 0.1.11 to 0.1.12.\r\n- [Release notes](https://github.com/golang/tools/releases)\r\n- [Commits](https://github.com/golang/tools/compare/v0.1.11...v0.1.12)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: golang.org/x/tools\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-08-08T13:21:23-07:00",
          "tree_id": "3959946fa07de7e1591a01dc229c9d6904841be4",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/13bcb05a1953b1afbb7a696163c27876c18071e5"
        },
        "date": 1659990661388,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36617395246,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.58 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9583025003,
            "unit": "ns/op\t         5.338 DeleteSeconds\t         4.187 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6408404563,
            "unit": "ns/op\t         4.243 DeleteSeconds\t         2.126 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b53fadb426f7623dab6a66b049fcf1a50945aced",
          "message": "Bump k8s.io/kube-aggregator from 0.22.11 to 0.22.12 (#791)\n\nBumps [k8s.io/kube-aggregator](https://github.com/kubernetes/kube-aggregator) from 0.22.11 to 0.22.12.\r\n- [Release notes](https://github.com/kubernetes/kube-aggregator/releases)\r\n- [Commits](https://github.com/kubernetes/kube-aggregator/compare/v0.22.11...v0.22.12)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: k8s.io/kube-aggregator\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-08-08T13:22:22-07:00",
          "tree_id": "0d031023927a3a3d80e284c56ad40f43c1d49f9b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/b53fadb426f7623dab6a66b049fcf1a50945aced"
        },
        "date": 1659990720057,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36460712422,
            "unit": "ns/op\t        18.89 DeleteSeconds\t        17.53 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9594255934,
            "unit": "ns/op\t         5.368 DeleteSeconds\t         4.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6413502424,
            "unit": "ns/op\t         4.241 DeleteSeconds\t         2.132 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "aaronshurley@users.noreply.github.com",
            "name": "Aaron Hurley",
            "username": "aaronshurley"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cd46fa2fea53c48621373d4ab496e920ab45de70",
          "message": "Update trivy-scan alert slack channel (#821)",
          "timestamp": "2022-08-09T15:02:58-04:00",
          "tree_id": "b2dca111e4362600f6e447071209997919ab9b80",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cd46fa2fea53c48621373d4ab496e920ab45de70"
        },
        "date": 1660072373946,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36450636541,
            "unit": "ns/op\t        18.82 DeleteSeconds\t        17.59 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9452639047,
            "unit": "ns/op\t         5.268 DeleteSeconds\t         4.144 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6408984984,
            "unit": "ns/op\t         4.251 DeleteSeconds\t         2.120 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "61ce44784f54b3b401416889284f28f5ce530678",
          "message": "add kctrl dev command (#638)\n\n* prep cli/vendor for dev deploy\r\n\r\n* introduce app dev deploy\r\n\r\n* add packageinstall support for dev\r\n\r\n* support in mem secret creation\r\n\r\n* move app dev deploy to dev deploy\r\n\r\n* support local fetch and kbld build\r\n\r\n* configure kubernetes dst\r\n\r\n* add debug logs to show what commands are running\r\n\r\n* add examples/cert-manager-tce-pkg\r\n\r\n* expose CreateToken in minimal dev deploy core client\r\n\r\n* WIP print errors in AppTailer\r\n\r\n* WIP provide k8s ca cert over env var\r\n\r\n* extract local deploy operations into a package\r\n\r\n* inject cmd runner\r\n\r\n* dev deploy => dev. Added a test for dev command.\r\n\r\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>\r\nCo-authored-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>\r\nCo-authored-by: Soumik Majumder <soumikm@vmware.com>",
          "timestamp": "2022-08-11T17:56:15+05:30",
          "tree_id": "370180ebdcd0abc6ca498ae59eabc05e36b93f91",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/61ce44784f54b3b401416889284f28f5ce530678"
        },
        "date": 1660221371112,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36770698757,
            "unit": "ns/op\t        19.07 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9500636254,
            "unit": "ns/op\t         5.291 DeleteSeconds\t         4.171 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6448394808,
            "unit": "ns/op\t         4.266 DeleteSeconds\t         2.139 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7f275ed6b635be93a6d2c22f022cea8c3d291bce",
          "message": "Merge pull request #824 from vmware-tanzu/bump-dependencies\n\nBump dependencies",
          "timestamp": "2022-08-11T10:37:47-06:00",
          "tree_id": "0d0d24127f6628c287862e25fe2e16b52b8594c8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7f275ed6b635be93a6d2c22f022cea8c3d291bce"
        },
        "date": 1660236569611,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36946418082,
            "unit": "ns/op\t        19.25 DeleteSeconds\t        17.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9587065694,
            "unit": "ns/op\t         5.340 DeleteSeconds\t         4.200 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6665276229,
            "unit": "ns/op\t         4.423 DeleteSeconds\t         2.191 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4d3e251b256db552c8457330de3e6e5ae4cbe882",
          "message": "Packages can constrain k8s and kc versions (#798)\n\n- kubernetes and kapp-controller version constraints can be overridden by annotations\r\n- range version selection of a Package chooses the highest version of\r\n  that package which also satisfies kc and k8s version constraints\r\n- Error messages if no package satisfies constraints provides detail of\r\n  which constraints failed.\r\n- simplified error message for the case where you just have zero\r\n  packages.\r\n- kapp-controller version is threaded through into the PKGI reconciler\r\n- refactor factories out to main thread, pass them down via dep injection (fixes a bug in the app factory where service account\r\n  token cache was reinitialized each reconcile\r\n- if PKGI specifies a different cluster, that cluster's version of k8s\r\n  is checked for the constraints (manually verified - no automated test\r\n  of this)\r\n- Error on no packages found\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\nCo-authored-by: Neil Hickey <nhickey@vmware.com>",
          "timestamp": "2022-08-11T15:02:17-04:00",
          "tree_id": "b463600ab52fcb75e39b5319e914532b1c0c3786",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4d3e251b256db552c8457330de3e6e5ae4cbe882"
        },
        "date": 1660245114592,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36404179668,
            "unit": "ns/op\t        18.89 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9490609527,
            "unit": "ns/op\t         5.277 DeleteSeconds\t         4.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6424828917,
            "unit": "ns/op\t         4.245 DeleteSeconds\t         2.141 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cd5e07973c67cd31b08dc73c12763eaeeef652cd",
          "message": "Upgrade GoLang from 1.18 to 1.19 (#822)\n\n* Upgrade GoLang from 1.18 to 1.19\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\n* Re-run generators\r\n\r\n- Upgrade go-lint for go 1.19 support\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>",
          "timestamp": "2022-08-11T15:03:01-04:00",
          "tree_id": "e15e95a3ff9be3bf93892254d77809212664de41",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cd5e07973c67cd31b08dc73c12763eaeeef652cd"
        },
        "date": 1660245242371,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37202962177,
            "unit": "ns/op\t        19.35 DeleteSeconds\t        17.78 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9863945011,
            "unit": "ns/op\t         5.506 DeleteSeconds\t         4.296 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7592908348,
            "unit": "ns/op\t         4.366 DeleteSeconds\t         3.180 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "leonde@vmware.com",
            "name": "Dennis Leon",
            "username": "DennisDenuto"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "04c7e27c3f1a1252fbc7edd45d3de7a7f8e9f709",
          "message": "feat: Surface App resources associated to a deploy (#799)\n\n- Uses kapp metadata file to list on the AppCR status the app label, namespaces and GK's\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>",
          "timestamp": "2022-08-11T15:05:07-04:00",
          "tree_id": "0c6b5fc94d10f6a897ff57d20e6a3d20f84c1130",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/04c7e27c3f1a1252fbc7edd45d3de7a7f8e9f709"
        },
        "date": 1660245421611,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37327775144,
            "unit": "ns/op\t        19.53 DeleteSeconds\t        17.75 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9608945510,
            "unit": "ns/op\t         5.360 DeleteSeconds\t         4.201 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6564816676,
            "unit": "ns/op\t         4.321 DeleteSeconds\t         2.183 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7c77dfac8d113d2d6502d8c9d974e219d4c1ceca",
          "message": "Add package authoring commands(#766)\n\n* introduce app dev deploy\r\n\r\n* add packageinstall support for dev\r\n\r\n* support in mem secret creation\r\n\r\n* move app dev deploy to dev deploy\r\n\r\n* add debug logs to show what commands are running\r\n\r\n* expose CreateToken in minimal dev deploy core client\r\n\r\n* extract local deploy operations into a package\r\n\r\n* introduce app dev deploy\r\n\r\n* add packageinstall support for dev\r\n\r\n* support in mem secret creation\r\n\r\n* move app dev deploy to dev deploy\r\n\r\n* configure kubernetes dst\r\n\r\n* add debug logs to show what commands are running\r\n\r\n* expose CreateToken in minimal dev deploy core client\r\n\r\n* WIP provide k8s ca cert over env var\r\n\r\n* extract local deploy operations into a package\r\n\r\n* Bump kapp controller. tidy and vendor dependencies\r\n\r\n* Adding app init and pkg init\r\n\r\n* Fixing misspelling github action\r\n\r\nFixing misspelling github action\r\n\r\n* Fixing the case where vendir.yml doesnt exist\r\n\r\n* Adopting review comments\r\n\r\n* refactored informational text for pkg and app init\r\n\r\n* Add kctrl package release command\r\n\r\n* Add kctrl package release command\r\n\r\n* Adopting review comments.\r\n\r\nAdopting review comments.\r\n\r\n* Adding App template transform\r\n\r\nAdding App template transform\r\n\r\n* Changes to release command to ensure it produces usable bundles and package resources\r\n\r\n* Added test case\r\n\r\nAdded test case\r\n\r\n* Making template section simpler\r\n\r\nMaking template section simpler\r\n\r\n* Adopted Text comments.\r\n\r\nAdopted Text comments.\r\n\r\n* Adopted review comment.\r\n\r\nAdopted review comment.\r\n\r\n* pkg repo release\r\n\r\n* pkg repo release\r\n\r\n* Formatting release command output. copy-to => repo-output\r\n\r\n* Adding the question for includePaths in case of Local directory.\r\n\r\nAdding the question for includePaths in case of Local directory.\r\n\r\n* Fixed bug.\r\n\r\nFixed bug.\r\n\r\n* Fixed path issue.\r\n\r\nFixed path issue.\r\n\r\n* Move release logic to app release. Cleanup artifact generation. Add repo output flag.\r\n\r\n* Add release section to package build. Spell fixes.\r\n\r\n* Add release section to package build. Spell fixes.\r\n\r\n* Fixed test case\r\n\r\nFixed test case\r\n\r\n* Add test for package release command. Set up local registry for tests\r\n\r\n* added e2e test for pkg repo release\r\n\r\n* Refactoring\r\n\r\nRefactoring\r\n\r\n* Add ValuesSchemaGen to generate calues schema for packages\r\n\r\n* Stricter checks on binary names. Refactoring and fixing typos. Kbld paths are clobbered if user wants to use lockfile.\r\n\r\n* Use otiai10/copy for copying directories. Tidy vendored files. Cleanup imgpkg runner.\r\n\r\n* Tighten package init test. Remove unused constant.\r\n\r\n* Adopting review comments.\r\n\r\nAdopting review comments.\r\n\r\n* dir copy using otiai10/copy in pkg repo release\r\n\r\n* Fix bad filepath base checks. Fix typos.\r\n\r\n* Refactor package init\r\n\r\n* Remove helmVersion from vendirExpectedOutput in e2e test\r\n\r\n* Refactor package init and release e2e test\r\n\r\n* Ensure that folder structure is retained while copying over files\r\n\r\n* Refactoring and clean up\r\n\r\n* Ensure that package-resources generated enables dev deploy\r\n\r\n* Update workflow to use hack script\r\n\r\n* Update workflow to use hack script\r\n\r\n* Fixing rerun case of Local Directory\r\n\r\nFixing rerun case of Local Directory\r\n\r\n* Ensure repo release creates valid repo bundle\r\n\r\n* Copy over pkg and pkg metadata from pkg-resources. Mark commands as experimental.\r\n\r\n* push the correct tag in repo release\r\n\r\n* Rebase on develop+kctrl-dev-deploy. Vendor and tidy.\r\n\r\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>\r\nCo-authored-by: Soumik Majumder <soumikm@vmware.com>\r\nCo-authored-by: Yash Sethiya <ysethiya@Yashs-MacBook-Pro.local>\r\nCo-authored-by: sethiyash <yashsethiya97@gmail.com>\r\nCo-authored-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>",
          "timestamp": "2022-08-16T13:28:24+05:30",
          "tree_id": "6bb7010cd425a09a6e04760447938239d6e942d6",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7c77dfac8d113d2d6502d8c9d974e219d4c1ceca"
        },
        "date": 1660637415463,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37178055613,
            "unit": "ns/op\t        19.41 DeleteSeconds\t        17.71 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9667199436,
            "unit": "ns/op\t         5.402 DeleteSeconds\t         4.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6643977419,
            "unit": "ns/op\t         4.367 DeleteSeconds\t         2.190 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "941bfdbc976b4d0f689ff2a56e221d3f045703c4",
          "message": "Update build script to vendor and tidy before building (#829)",
          "timestamp": "2022-08-16T14:10:13+05:30",
          "tree_id": "cd0043d4ae7776fda5fc0e2683597c224b821d89",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/941bfdbc976b4d0f689ff2a56e221d3f045703c4"
        },
        "date": 1660639842267,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36661954883,
            "unit": "ns/op\t        19.08 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9589329495,
            "unit": "ns/op\t         5.371 DeleteSeconds\t         4.180 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6451523299,
            "unit": "ns/op\t         4.268 DeleteSeconds\t         2.137 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "aaronshurley@users.noreply.github.com",
            "name": "Aaron Hurley",
            "username": "aaronshurley"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f96ecd0102a90fb32c2c9b4c3890392589c1a7fb",
          "message": "Update backlog link",
          "timestamp": "2022-08-16T11:48:33-07:00",
          "tree_id": "c277b805a2b103abc9a55c4d9d1f8025bed3f7f8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f96ecd0102a90fb32c2c9b4c3890392589c1a7fb"
        },
        "date": 1660676345612,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36611529604,
            "unit": "ns/op\t        19.01 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9676008483,
            "unit": "ns/op\t         5.424 DeleteSeconds\t         4.197 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6492832543,
            "unit": "ns/op\t         4.281 DeleteSeconds\t         2.158 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "73543413caf5f532fb3dd4d13b6d0658df1984c4",
          "message": "Kctrl pkg init release text rephrase (#830)\n\n* Rephrasing the text.\r\n\r\nRephrasing the text.\r\n\r\n* Fixing the test case\r\n\r\nFixing the test case\r\n\r\n* Fixing the test case\r\n\r\n* Fixing package Repo Test case\r\n\r\nFixing package Repo Test case\r\n\r\n* Removed Step\r\n\r\nRemoved Step\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Increasing the sleep\r\n\r\n* Removing extra line from text\r\n\r\nRemoving extra line from text\r\n\r\n* Removing fmt\r\n\r\n* Printing the interactive text in test\r\n\r\nPrinting the interactive text in test\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Update package_authoring_e2e_test.go",
          "timestamp": "2022-08-18T02:44:41+05:30",
          "tree_id": "d2a1eb4db930e02ece13400fea52ae3d8a3cb78c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/73543413caf5f532fb3dd4d13b6d0658df1984c4"
        },
        "date": 1660771467394,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36509278170,
            "unit": "ns/op\t        18.97 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9490914197,
            "unit": "ns/op\t         5.289 DeleteSeconds\t         4.159 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6446979442,
            "unit": "ns/op\t         4.253 DeleteSeconds\t         2.138 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ae4bff9c3a5f589040eab71fd12edfd35b846133",
          "message": "Ensure that default version is valid semver. Add error check in case of malformed package-build (#834)",
          "timestamp": "2022-08-18T17:24:55+05:30",
          "tree_id": "f240a0d2b04881f144f541a6acf6394106dfcc97",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ae4bff9c3a5f589040eab71fd12edfd35b846133"
        },
        "date": 1660824438792,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37013891562,
            "unit": "ns/op\t        19.29 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9755719509,
            "unit": "ns/op\t         5.498 DeleteSeconds\t         4.198 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6621318464,
            "unit": "ns/op\t         4.361 DeleteSeconds\t         2.180 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "soumikm@vmware.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "distinct": true,
          "id": "59e69562f0e253439781a6196387fe77ce3db4ae",
          "message": "Add timeout to prompt output tests. Remove now unused code.",
          "timestamp": "2022-08-22T16:06:19+05:30",
          "tree_id": "4a4b33ee3bed480cc3c045c2739481be61f8b08f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/59e69562f0e253439781a6196387fe77ce3db4ae"
        },
        "date": 1661165264517,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36798425296,
            "unit": "ns/op\t        19.20 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9546099253,
            "unit": "ns/op\t         5.337 DeleteSeconds\t         4.169 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6621500283,
            "unit": "ns/op\t         4.430 DeleteSeconds\t         2.144 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f114256a67649bda8d9ccdfb0865216d280cdb93",
          "message": "Merge pull request #835 from vmware-tanzu/dependabot/go_modules/k8s.io/kube-aggregator-0.22.13\n\nBump k8s.io/kube-aggregator from 0.22.12 to 0.22.13",
          "timestamp": "2022-08-22T11:39:07-06:00",
          "tree_id": "5261423c1db414e20b3f55563e2dd2a41968c704",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f114256a67649bda8d9ccdfb0865216d280cdb93"
        },
        "date": 1661190560110,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36467853578,
            "unit": "ns/op\t        18.94 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9543277277,
            "unit": "ns/op\t         5.336 DeleteSeconds\t         4.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6444829300,
            "unit": "ns/op\t         4.272 DeleteSeconds\t         2.131 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "iamlizhiyong@outlook.com",
            "name": "Zhiyong Li",
            "username": "showpune"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e997c38972eb8c3b24368c29e814c59d1ff8a81f",
          "message": "Namespace isolation (#826)\n\n* Add Namespace isolation\r\n\r\n* Fix the issue from merge comment\r\n\r\n* Remove the unncessary change\r\n\r\n* remove teh namespace deploy script\r\n\r\n* Change to use KAPPCTRL_START_API_SERVER  env\r\n\r\n* Use flag instead of env\r\n\r\n* Use flag instead of env\r\n\r\n* rename feature flag to start-api-server",
          "timestamp": "2022-08-22T10:40:52-07:00",
          "tree_id": "9584142407d17d8ccd64d0864180a716e173bc45",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/e997c38972eb8c3b24368c29e814c59d1ff8a81f"
        },
        "date": 1661190680714,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36529225752,
            "unit": "ns/op\t        18.96 DeleteSeconds\t        17.53 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9676484943,
            "unit": "ns/op\t         5.456 DeleteSeconds\t         4.175 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6499051006,
            "unit": "ns/op\t         4.289 DeleteSeconds\t         2.137 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "63df2dbb56af2c3f948a54044cfe1de74e0b04ae",
          "message": "Merge pull request #839 from vmware-tanzu/nh-fix-case-start-api-server\n\nfix case in `start-api-server` flag",
          "timestamp": "2022-08-22T12:13:11-06:00",
          "tree_id": "b89981ea72d14eefafb8704c2d6560042983c473",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/63df2dbb56af2c3f948a54044cfe1de74e0b04ae"
        },
        "date": 1661192601263,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36742883328,
            "unit": "ns/op\t        19.09 DeleteSeconds\t        17.57 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9704037126,
            "unit": "ns/op\t         5.433 DeleteSeconds\t         4.219 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6498540097,
            "unit": "ns/op\t         4.305 DeleteSeconds\t         2.149 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fe0dd27e0854be97e45b6dc08f29b8799b832436",
          "message": "Merge pull request #841 from vmware-tanzu/dependabot/github_actions/peter-evans/create-pull-request-4.1.1\n\nBump peter-evans/create-pull-request from 4.0.4 to 4.1.1",
          "timestamp": "2022-08-22T18:39:14-06:00",
          "tree_id": "d298adcbea491130950cab356921e12203f4b974",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/fe0dd27e0854be97e45b6dc08f29b8799b832436"
        },
        "date": 1661215758064,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36795805375,
            "unit": "ns/op\t        19.16 DeleteSeconds\t        17.60 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9576354909,
            "unit": "ns/op\t         5.350 DeleteSeconds\t         4.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6600505108,
            "unit": "ns/op\t         4.369 DeleteSeconds\t         2.155 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "soumikm@vmware.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "distinct": true,
          "id": "688596406c6d13596ae5bf072ccde99677b14d62",
          "message": "Handle errors after successful reconciliation. Format zero timestamps better.",
          "timestamp": "2022-08-29T11:05:25+05:30",
          "tree_id": "6599367931f1daa5def4df32aec0777bafcf657a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/688596406c6d13596ae5bf072ccde99677b14d62"
        },
        "date": 1661752063959,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37028527972,
            "unit": "ns/op\t        19.25 DeleteSeconds\t        17.73 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9679885676,
            "unit": "ns/op\t         5.409 DeleteSeconds\t         4.211 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6614795488,
            "unit": "ns/op\t         4.382 DeleteSeconds\t         2.173 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "24ff00581ec57dea9430e66aab4582711914fbeb",
          "message": "Merge pull request #840 from vmware-tanzu/dependabot/github_actions/softprops/action-gh-release-1e07f4398721186383de40550babbdf2b84acfc5\n\nBump softprops/action-gh-release from 17cd0d34deddf848fc0e7d9be5202c148c270a0a to 1",
          "timestamp": "2022-08-29T13:14:37-06:00",
          "tree_id": "714071160832162c4d6ae48ac0036de722a5c4f5",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/24ff00581ec57dea9430e66aab4582711914fbeb"
        },
        "date": 1661801189054,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37281763590,
            "unit": "ns/op\t        19.43 DeleteSeconds\t        17.79 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9613331803,
            "unit": "ns/op\t         5.371 DeleteSeconds\t         4.192 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6552334462,
            "unit": "ns/op\t         4.322 DeleteSeconds\t         2.178 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "joaod@vmware.com",
            "name": "João Pereira",
            "username": "joaopapereira"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ee8e7af525c1eb6a2b866115e5e1f4e792720f17",
          "message": "Fixes dev-deploy script (#843)\n\nSigned-off-by: João Pereira <joaod@vmware.com>\r\n\r\nSigned-off-by: João Pereira <joaod@vmware.com>",
          "timestamp": "2022-08-29T13:40:15-07:00",
          "tree_id": "16eb1b57bc6889a1da01a26073d9935977d2c203",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ee8e7af525c1eb6a2b866115e5e1f4e792720f17"
        },
        "date": 1661806373891,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38256523417,
            "unit": "ns/op\t        19.39 DeleteSeconds\t        18.77 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9791111186,
            "unit": "ns/op\t         5.406 DeleteSeconds\t         4.333 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7781573372,
            "unit": "ns/op\t         4.528 DeleteSeconds\t         3.203 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a2953f466e05227184e48fa88f25fc7bbb49750e",
          "message": "Merge pull request #836 from vmware-tanzu/bump-dependencies\n\nBump dependencies",
          "timestamp": "2022-08-29T14:43:29-06:00",
          "tree_id": "43f76b79870fa9ed8e896d55bc7fe1d435e757b9",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/a2953f466e05227184e48fa88f25fc7bbb49750e"
        },
        "date": 1661806571111,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37981451531,
            "unit": "ns/op\t        19.14 DeleteSeconds\t        18.79 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9620261611,
            "unit": "ns/op\t         5.363 DeleteSeconds\t         4.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6684033589,
            "unit": "ns/op\t         4.406 DeleteSeconds\t         2.191 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "soumikm@vmware.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "distinct": true,
          "id": "3fa5e21db1fcc2edd46f5effaeff057b7ac1de15",
          "message": "Ensure that updated package specs are copied over from package spec",
          "timestamp": "2022-09-01T10:17:10+05:30",
          "tree_id": "637c26a14b7fba1690059c958ac9d57fdf44abaf",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3fa5e21db1fcc2edd46f5effaeff057b7ac1de15"
        },
        "date": 1662008219565,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36439393058,
            "unit": "ns/op\t        18.85 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9503238804,
            "unit": "ns/op\t         5.324 DeleteSeconds\t         4.141 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6438168572,
            "unit": "ns/op\t         4.226 DeleteSeconds\t         2.171 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "soumikm@vmware.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "distinct": true,
          "id": "2a832b079e04eae5ef40d38ac3752da0ed77178b",
          "message": "Fix un-fmt'd files",
          "timestamp": "2022-09-06T17:25:23+05:30",
          "tree_id": "bf651d9ccd39454827c6efe4c325539819ef61aa",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2a832b079e04eae5ef40d38ac3752da0ed77178b"
        },
        "date": 1662466126972,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38369044549,
            "unit": "ns/op\t        19.51 DeleteSeconds\t        18.79 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9730528354,
            "unit": "ns/op\t         5.423 DeleteSeconds\t         4.254 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6696436890,
            "unit": "ns/op\t         4.402 DeleteSeconds\t         2.219 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "870f5a1049190121451eb79f8a1835b9e960649f",
          "message": "Pass additional information to downward API (k8s version, kc version, k8s g/v) (#846)\n\n* WIP: Retrieve version information prior to fet,temp,deploy steps\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\n* Fix linter / import statements\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\n* Address review comments\r\n\r\n- renaming some things\r\n- moved away from a values factory back to values struct\r\n- minor fixups\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\n* Fixup errors in app_template\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\n* Add memoized fetching of versions\r\n\r\n- Add template() test to validate memoizing works\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\n* remove unnecessary memoization\r\n\r\n- memoizing within packageinstall didnt actually memoize\r\n- memoizing within componentinfo was too aggressive so controller would not receive updated version after cluster is updated\r\n\r\n* rename kubernetesGroupVersions to kubernetesAPIs\r\n\r\n* use array as type of values for kubernetesAPIs\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-09-07T14:24:30-04:00",
          "tree_id": "44f498040c6217a0f7fc965bf4f4bbcd730f0463",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/870f5a1049190121451eb79f8a1835b9e960649f"
        },
        "date": 1662575676374,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36465320966,
            "unit": "ns/op\t        18.91 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9462344267,
            "unit": "ns/op\t         5.273 DeleteSeconds\t         4.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6416754079,
            "unit": "ns/op\t         4.242 DeleteSeconds\t         2.134 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "aaronshurley@users.noreply.github.com",
            "name": "Aaron Hurley",
            "username": "aaronshurley"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c9b4dbaea3bb898e1e2e675cc9395299a0faec3d",
          "message": "update project docs (#858)",
          "timestamp": "2022-09-07T12:02:15-07:00",
          "tree_id": "19921dce26d3a962f17e1d0f9f528dc853758693",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c9b4dbaea3bb898e1e2e675cc9395299a0faec3d"
        },
        "date": 1662577983089,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37727375670,
            "unit": "ns/op\t        19.11 DeleteSeconds\t        18.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9605606013,
            "unit": "ns/op\t         5.381 DeleteSeconds\t         4.169 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6575538648,
            "unit": "ns/op\t         4.304 DeleteSeconds\t         2.231 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3f1c263369c0a543834cb5b4f1796c3c354458f0",
          "message": "Merge pull request #859 from vmware-tanzu/dk-min-app-sync-period\n\nintroduce min app sync period",
          "timestamp": "2022-09-09T12:09:22-04:00",
          "tree_id": "644464892095bd5e1bbcac33b4993ac1a2898baf",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3f1c263369c0a543834cb5b4f1796c3c354458f0"
        },
        "date": 1662740561257,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37323172960,
            "unit": "ns/op\t        19.40 DeleteSeconds\t        17.85 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9733531760,
            "unit": "ns/op\t         5.429 DeleteSeconds\t         4.252 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6675231351,
            "unit": "ns/op\t         4.373 DeleteSeconds\t         2.249 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "80be13ad451b6d06e9fd8e2846c240f376a5e6ab",
          "message": "never report the kubernetes version with the pre or buildmeta (#862)\n\n* never report the kubernetes version with the pre or buildmeta\r\n\r\n* test should pass but its a wip bc weve reduced test coverage\r\n\r\n* adds new test for component_info\r\n\r\nasserts prerelease gets scrubbed at that level",
          "timestamp": "2022-09-13T11:04:41-04:00",
          "tree_id": "a224ea954dad370be8b72ac42a4612cbbf563b57",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/80be13ad451b6d06e9fd8e2846c240f376a5e6ab"
        },
        "date": 1663082071654,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36489057263,
            "unit": "ns/op\t        18.94 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9508824500,
            "unit": "ns/op\t         5.304 DeleteSeconds\t         4.160 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6629166632,
            "unit": "ns/op\t         4.346 DeleteSeconds\t         2.240 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c0eaf6da0d2c3f4a12b06a720b8b2e6a256ee752",
          "message": "add package details to child app cr annotations so that they can be used in app cr downward api (#864)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-09-13T11:05:53-04:00",
          "tree_id": "859b781b011bf201278452a4347fa51e34ef4162",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c0eaf6da0d2c3f4a12b06a720b8b2e6a256ee752"
        },
        "date": 1663082163577,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36466094757,
            "unit": "ns/op\t        18.92 DeleteSeconds\t        17.51 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9534263451,
            "unit": "ns/op\t         5.343 DeleteSeconds\t         4.152 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6466250278,
            "unit": "ns/op\t         4.275 DeleteSeconds\t         2.152 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "68ecccad508e4498188631d31089747ce5743fd8",
          "message": "changes to make it easier to see where tests fail (#865)",
          "timestamp": "2022-09-13T14:39:44-07:00",
          "tree_id": "0a877a0ede16bd9970cda966275222d75ee92806",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/68ecccad508e4498188631d31089747ce5743fd8"
        },
        "date": 1663105950638,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37167839938,
            "unit": "ns/op\t        19.47 DeleteSeconds\t        17.65 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9698607891,
            "unit": "ns/op\t         5.455 DeleteSeconds\t         4.192 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6542682964,
            "unit": "ns/op\t         4.345 DeleteSeconds\t         2.155 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7d217d8238f4fa8a53d2052bb05246998b915644",
          "message": "template values AsPaths: clearer error msg (#871)",
          "timestamp": "2022-09-14T13:28:41-04:00",
          "tree_id": "d7be04ab6641dd9b8747566b7eb5187fea13faf9",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7d217d8238f4fa8a53d2052bb05246998b915644"
        },
        "date": 1663177215788,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37121786757,
            "unit": "ns/op\t        19.38 DeleteSeconds\t        17.67 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9692134154,
            "unit": "ns/op\t         5.435 DeleteSeconds\t         4.209 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6540709763,
            "unit": "ns/op\t         4.320 DeleteSeconds\t         2.168 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "723ecded6b59d2bc4f2c61bf85e0dbc155400750",
          "message": "adds example using simple-app and downwardAPI (#870)",
          "timestamp": "2022-09-14T13:29:23-04:00",
          "tree_id": "c671219b5d5f2a71dc4a28ffae1f874386a0a28f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/723ecded6b59d2bc4f2c61bf85e0dbc155400750"
        },
        "date": 1663177279958,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37163307935,
            "unit": "ns/op\t        19.45 DeleteSeconds\t        17.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9697082326,
            "unit": "ns/op\t         5.410 DeleteSeconds\t         4.228 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6605900613,
            "unit": "ns/op\t         4.393 DeleteSeconds\t         2.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "soumikm@vmware.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "distinct": true,
          "id": "ccdb1be9bdbff5ed826ebf0b8a5bf629c16e3a41",
          "message": "Ensure that ytt overlays secrets are garbage collected. Refactor lengthy conditionals",
          "timestamp": "2022-09-14T23:21:21+05:30",
          "tree_id": "ce94970b768e62166ce4566fb28a26f651b0ae3b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ccdb1be9bdbff5ed826ebf0b8a5bf629c16e3a41"
        },
        "date": 1663178608865,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37387033452,
            "unit": "ns/op\t        19.71 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9647542903,
            "unit": "ns/op\t         5.390 DeleteSeconds\t         4.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6602486961,
            "unit": "ns/op\t         4.357 DeleteSeconds\t         2.190 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "560aaea259bc7deef959fca65f0b70506712330f",
          "message": "Merge pull request #876 from vmware-tanzu/depsup\n\nbump k8s libraries to 1.25",
          "timestamp": "2022-09-15T18:21:12-04:00",
          "tree_id": "d06f525953d5ab72862bf21f5a45eb2c65055bf1",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/560aaea259bc7deef959fca65f0b70506712330f"
        },
        "date": 1663281149561,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36700000125,
            "unit": "ns/op\t        19.11 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9620625546,
            "unit": "ns/op\t         5.382 DeleteSeconds\t         4.161 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6499953564,
            "unit": "ns/op\t         4.269 DeleteSeconds\t         2.187 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "5a87cb30016d322d4276bf46df61eaca20088703",
          "message": "test k8s 1.20 and 1.25 (#880)",
          "timestamp": "2022-09-15T19:30:54-04:00",
          "tree_id": "d3c2e885d8b21659600fac59ed95575d19f712a1",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/5a87cb30016d322d4276bf46df61eaca20088703"
        },
        "date": 1663285345458,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36742342972,
            "unit": "ns/op\t        19.07 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9637171559,
            "unit": "ns/op\t         5.403 DeleteSeconds\t         4.188 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6461986422,
            "unit": "ns/op\t         4.274 DeleteSeconds\t         2.141 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2d3ed642504c47881bca29305f04b45b5342b9cc",
          "message": "Deflake package repo tests. Using only necessary checks (#866)",
          "timestamp": "2022-09-16T08:43:49-04:00",
          "tree_id": "478bfcc1a021da480d288b03dee41286046d854e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2d3ed642504c47881bca29305f04b45b5342b9cc"
        },
        "date": 1663332982975,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37072876510,
            "unit": "ns/op\t        19.26 DeleteSeconds\t        17.74 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9707658904,
            "unit": "ns/op\t         5.417 DeleteSeconds\t         4.214 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6722036090,
            "unit": "ns/op\t         4.367 DeleteSeconds\t         2.299 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2f3288952790273095601900f953ae08e52eab8c",
          "message": "Adding examples of package authoring (#861)\n\n* Adding examples of package authoring\r\n\r\nAdding examples of package authoring\r\n\r\n* Moved to another folder\r\n\r\nMoved to another folder\r\n\r\n* Adopting comments.\r\n\r\nAdopting comments.",
          "timestamp": "2022-09-19T13:04:42+05:30",
          "tree_id": "ae2a4a511c69c3a4ae2a450a1a2554b88265d51b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2f3288952790273095601900f953ae08e52eab8c"
        },
        "date": 1663573507801,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36497691147,
            "unit": "ns/op\t        18.90 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9535286992,
            "unit": "ns/op\t         5.346 DeleteSeconds\t         4.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6423418592,
            "unit": "ns/op\t         4.249 DeleteSeconds\t         2.132 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "yashsethiya97@gmail.com",
            "name": "Yash Sethiya",
            "username": "sethiyash"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b533f41b0d99d50313b5aff7ab284dcad2f7ab1e",
          "message": "kctrl: Tightning  up the pkg authoring e2e testcases (#850)\n\n* added pkg authoring e2e testcases flow\r\n\r\n* commented installing pkg\r\n\r\n* added git repo flow\r\n\r\n* cleaning up installed pkg with defer\r\n\r\n* using simple-app for git repo flow\r\n\r\n* adding expected outputs in testcases struct\r\n\r\n* adopted nits",
          "timestamp": "2022-09-19T13:18:55+05:30",
          "tree_id": "d040a5759af4767f1701136e2b7bba7fddfc46fd",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/b533f41b0d99d50313b5aff7ab284dcad2f7ab1e"
        },
        "date": 1663574454937,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37172053228,
            "unit": "ns/op\t        19.30 DeleteSeconds\t        17.81 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9619585418,
            "unit": "ns/op\t         5.356 DeleteSeconds\t         4.207 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6548834420,
            "unit": "ns/op\t         4.317 DeleteSeconds\t         2.183 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "53a8b441c112c25e1b15d260a17602d135f86fa6",
          "message": "Add new kapp flags to allowed change opts (#887)\n\n- --exit-early-on-apply-error\r\n- --exit-early-on-wait-error",
          "timestamp": "2022-09-19T12:04:10-04:00",
          "tree_id": "83911c64cb28019ddcfa4e7ce320ce3adf8cce39",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/53a8b441c112c25e1b15d260a17602d135f86fa6"
        },
        "date": 1663604084968,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36528870222,
            "unit": "ns/op\t        18.93 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9543386954,
            "unit": "ns/op\t         5.350 DeleteSeconds\t         4.155 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6468148161,
            "unit": "ns/op\t         4.278 DeleteSeconds\t         2.138 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b6153d09c5f94e932bbdd03bd2c65270387f15c4",
          "message": "clean up sidecarexec socket file in case of previous unclean process termination (#881)\n\nCo-authored-by: Dmitriy Kalinin <dkalinin@vmware.com>",
          "timestamp": "2022-09-19T12:11:09-04:00",
          "tree_id": "429c50ffd5c00752c8aff16eb43eb0bbffaf434a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/b6153d09c5f94e932bbdd03bd2c65270387f15c4"
        },
        "date": 1663604500343,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36485285001,
            "unit": "ns/op\t        18.91 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9540810629,
            "unit": "ns/op\t         5.347 DeleteSeconds\t         4.155 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6465152859,
            "unit": "ns/op\t         4.259 DeleteSeconds\t         2.136 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "882eeb50389ea1d18ee4e1278d27719adf785523",
          "message": "configurable tls cipher suites (#882)",
          "timestamp": "2022-09-19T12:07:20-04:00",
          "tree_id": "8688064c4cca447b3ef97d7b8c32a01bcfdfe5fc",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/882eeb50389ea1d18ee4e1278d27719adf785523"
        },
        "date": 1663604550665,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37715593497,
            "unit": "ns/op\t        19.58 DeleteSeconds\t        18.07 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9789044440,
            "unit": "ns/op\t         5.423 DeleteSeconds\t         4.301 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6751839123,
            "unit": "ns/op\t         4.463 DeleteSeconds\t         2.239 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2299295239a426ba6db05d2874944b4e67aa80df",
          "message": "Fix bash completion for kctrl (#889)\n\n- Do not print Succeeded for the help command\r\n- Use SetOut from cobra to set output for cmd.Help()",
          "timestamp": "2022-09-19T13:35:43-04:00",
          "tree_id": "66887d9801eeb6b9cb2840e255fc5788fb131a7e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2299295239a426ba6db05d2874944b4e67aa80df"
        },
        "date": 1663609580027,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36569708861,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.53 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9578577917,
            "unit": "ns/op\t         5.370 DeleteSeconds\t         4.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6471259746,
            "unit": "ns/op\t         4.288 DeleteSeconds\t         2.144 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ca9fa4e3e4641dcf022ae64ad9dd0fc02bf99f12",
          "message": "Bump slackapi/slack-github-action from 1.21.0 to 1.22.0 (#879)\n\nBumps [slackapi/slack-github-action](https://github.com/slackapi/slack-github-action) from 1.21.0 to 1.22.0.\r\n- [Release notes](https://github.com/slackapi/slack-github-action/releases)\r\n- [Commits](https://github.com/slackapi/slack-github-action/compare/v1.21.0...v1.22.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: slackapi/slack-github-action\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-19T11:08:49-07:00",
          "tree_id": "fd73802efd84d3c2506d2baa478196580201ba97",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ca9fa4e3e4641dcf022ae64ad9dd0fc02bf99f12"
        },
        "date": 1663611554175,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36621163471,
            "unit": "ns/op\t        19.04 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9550863843,
            "unit": "ns/op\t         5.357 DeleteSeconds\t         4.154 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6493662290,
            "unit": "ns/op\t         4.297 DeleteSeconds\t         2.132 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d5af2b75f38c483929a0d9d5b8079558325b0532",
          "message": "Bump k8s.io/kube-aggregator from 0.22.13 to 0.22.14 (#883)\n\nBumps [k8s.io/kube-aggregator](https://github.com/kubernetes/kube-aggregator) from 0.22.13 to 0.22.14.\r\n- [Release notes](https://github.com/kubernetes/kube-aggregator/releases)\r\n- [Commits](https://github.com/kubernetes/kube-aggregator/compare/v0.22.13...v0.22.14)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: k8s.io/kube-aggregator\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-19T11:42:44-07:00",
          "tree_id": "3f81059d24559d3e617b162f836b07dabd903475",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d5af2b75f38c483929a0d9d5b8079558325b0532"
        },
        "date": 1663613582147,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36440714135,
            "unit": "ns/op\t        18.84 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9446603995,
            "unit": "ns/op\t         5.260 DeleteSeconds\t         4.148 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6380474913,
            "unit": "ns/op\t         4.224 DeleteSeconds\t         2.120 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6d0bf536e2d95a6eb1e0b7ed6ba2c230f05035ec",
          "message": "Bump golang from 1.19.0 to 1.19.1 (#857)\n\nBumps golang from 1.19.0 to 1.19.1.\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: golang\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-19T11:57:17-07:00",
          "tree_id": "51d6a238e86d7746a76b9e6ecd408db29d1b2ac5",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/6d0bf536e2d95a6eb1e0b7ed6ba2c230f05035ec"
        },
        "date": 1663614694452,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37149776792,
            "unit": "ns/op\t        19.31 DeleteSeconds\t        17.79 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9725758448,
            "unit": "ns/op\t         5.422 DeleteSeconds\t         4.254 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6587401011,
            "unit": "ns/op\t         4.352 DeleteSeconds\t         2.185 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2e523c0c504dc133ec5f324307c96fd00297eb8b",
          "message": "Bump actions/checkout from 3.0.1 to 3.0.2 (#869)\n\nBumps [actions/checkout](https://github.com/actions/checkout) from 3.0.1 to 3.0.2.\r\n- [Release notes](https://github.com/actions/checkout/releases)\r\n- [Changelog](https://github.com/actions/checkout/blob/main/CHANGELOG.md)\r\n- [Commits](https://github.com/actions/checkout/compare/v3.0.1...v3.0.2)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: actions/checkout\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-19T14:59:10-07:00",
          "tree_id": "dfa8a58700cddb0755cf0a0856187226240c78f1",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2e523c0c504dc133ec5f324307c96fd00297eb8b"
        },
        "date": 1663625567669,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38106696921,
            "unit": "ns/op\t        19.31 DeleteSeconds\t        18.74 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9744962531,
            "unit": "ns/op\t         5.459 DeleteSeconds\t         4.235 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6721326452,
            "unit": "ns/op\t         4.428 DeleteSeconds\t         2.208 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "010256002aa5f7c52af8f2feafe528ea9bf3cb07",
          "message": "Bump k8s.io/component-base from 0.25.0 to 0.25.1 (#890)\n\nBumps [k8s.io/component-base](https://github.com/kubernetes/component-base) from 0.25.0 to 0.25.1.\r\n- [Release notes](https://github.com/kubernetes/component-base/releases)\r\n- [Commits](https://github.com/kubernetes/component-base/compare/v0.25.0...v0.25.1)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: k8s.io/component-base\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-19T15:43:36-07:00",
          "tree_id": "eedb3276cc688a08141168ee87806d3bf8bfacf3",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/010256002aa5f7c52af8f2feafe528ea9bf3cb07"
        },
        "date": 1663628235895,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37228895120,
            "unit": "ns/op\t        19.41 DeleteSeconds\t        17.76 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9916312551,
            "unit": "ns/op\t         5.641 DeleteSeconds\t         4.228 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6613354530,
            "unit": "ns/op\t         4.388 DeleteSeconds\t         2.177 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jryan@pivotal.io",
            "name": "John S. Ryan",
            "username": "pivotaljohn"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1f9c74299ca0c49ba386178767dacfaddf07d5da",
          "message": "Bump ytt, kbld, kapp (not imgpkg,vendir) (#891)\n\nimgpkg and vendir are not ready to be released.\r\n\r\nCo-authored-by: John Ryan <jtigger@infosysengr.com>",
          "timestamp": "2022-09-19T16:42:59-07:00",
          "tree_id": "058a7cc818889152a46bc67f62a0977d1fefa580",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/1f9c74299ca0c49ba386178767dacfaddf07d5da"
        },
        "date": 1663631737171,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37226761338,
            "unit": "ns/op\t        19.39 DeleteSeconds\t        17.76 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9777850576,
            "unit": "ns/op\t         5.433 DeleteSeconds\t         4.288 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6609487697,
            "unit": "ns/op\t         4.343 DeleteSeconds\t         2.211 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c3eb24d9ee59ab3bc48cdb18c87186676ecadba4",
          "message": "Ensure that  For option helm Chart from Git, pkg init sync's properly (#852)\n\n* Fixing the bug - For option helm Chart from Git, pkg init is not syncing the helm charts.\r\n\r\n* Update init.go\r\n\r\n* Adopting review comments.\r\n\r\nAdopting review comments.\r\n\r\n* Adding test case\r\n\r\nAdding test case\r\n\r\n* Update package_authoring_e2e_test.go",
          "timestamp": "2022-09-20T08:32:34+05:30",
          "tree_id": "4ab72603101f41625c0b4d987e6b50cbb14e5b1a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c3eb24d9ee59ab3bc48cdb18c87186676ecadba4"
        },
        "date": 1663643580439,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36579934114,
            "unit": "ns/op\t        19.04 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9655249220,
            "unit": "ns/op\t         5.381 DeleteSeconds\t         4.231 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6428575124,
            "unit": "ns/op\t         4.252 DeleteSeconds\t         2.137 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "feaf16d331c32bbaaa877205eb1634f1d1c19b1d",
          "message": "Fixing kctrl dev failing on GKE (#885)\n\n* Fixing kctrl dev failing on GKE\r\n\r\nFixing kctrl dev failing on GKE\r\n\r\n* Update detailed_cmd_runner.go\r\n\r\n* Update detailed_cmd_runner.go\r\n\r\n* Adopting review comments\r\n\r\nAdopting review comments",
          "timestamp": "2022-09-20T09:35:30+05:30",
          "tree_id": "beeb640462fafc0ead06a5f94268991e78ee63bf",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/feaf16d331c32bbaaa877205eb1634f1d1c19b1d"
        },
        "date": 1663647555038,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38091205856,
            "unit": "ns/op\t        19.26 DeleteSeconds\t        18.76 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9895118302,
            "unit": "ns/op\t         5.620 DeleteSeconds\t         4.212 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6645763690,
            "unit": "ns/op\t         4.380 DeleteSeconds\t         2.211 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e0954d8081f80ce0c7cee85f18f8568531c4749c",
          "message": "Bump dependencies for cli (#888)\n\n* Bump k8s libraries\r\n\r\n* Bump go version",
          "timestamp": "2022-09-20T12:12:59+05:30",
          "tree_id": "aaaa46e86e580138036715887cbba0c24ea12575",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/e0954d8081f80ce0c7cee85f18f8568531c4749c"
        },
        "date": 1663656937326,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37179267522,
            "unit": "ns/op\t        19.29 DeleteSeconds\t        17.82 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9710710397,
            "unit": "ns/op\t         5.424 DeleteSeconds\t         4.231 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6698107411,
            "unit": "ns/op\t         4.426 DeleteSeconds\t         2.203 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8b093174e9139a277101a8104091a6070e6eb77b",
          "message": "Restructure help sections for all commands (#860)\n\n* Restructure help sections for all commands\r\n\r\n* Fix help test\r\n\r\n* Renaming annotation key. Deferring removal of app init",
          "timestamp": "2022-09-20T14:10:21+05:30",
          "tree_id": "c9d451cc1ea72a4fe776123e8b77886a788b09c7",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8b093174e9139a277101a8104091a6070e6eb77b"
        },
        "date": 1663663892446,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36881701204,
            "unit": "ns/op\t        19.28 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9546099524,
            "unit": "ns/op\t         5.331 DeleteSeconds\t         4.169 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6603231043,
            "unit": "ns/op\t         4.376 DeleteSeconds\t         2.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7841dd203a35aee3d7fdaabcdd3fa3914feafdef",
          "message": "Bump k8s.io/kube-aggregator from 0.22.14 to 0.22.15 (#895)\n\nBumps [k8s.io/kube-aggregator](https://github.com/kubernetes/kube-aggregator) from 0.22.14 to 0.22.15.\r\n- [Release notes](https://github.com/kubernetes/kube-aggregator/releases)\r\n- [Commits](https://github.com/kubernetes/kube-aggregator/compare/v0.22.14...v0.22.15)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: k8s.io/kube-aggregator\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-23T12:07:54-07:00",
          "tree_id": "87c07beec0b0a1e0868c57b8dfbcee7b8b543451",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7841dd203a35aee3d7fdaabcdd3fa3914feafdef"
        },
        "date": 1663960697056,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35635078214,
            "unit": "ns/op\t        18.00 DeleteSeconds\t        17.59 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8454617679,
            "unit": "ns/op\t         4.262 DeleteSeconds\t         4.149 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6401437856,
            "unit": "ns/op\t         4.232 DeleteSeconds\t         2.132 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e0b396525002694623c11594a7ae5b9aa7558f74",
          "message": "Bump helm/kind-action from 1.3.0 to 1.4.0 (#900)\n\nBumps [helm/kind-action](https://github.com/helm/kind-action) from 1.3.0 to 1.4.0.\r\n- [Release notes](https://github.com/helm/kind-action/releases)\r\n- [Commits](https://github.com/helm/kind-action/compare/v1.3.0...v1.4.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: helm/kind-action\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-minor\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-27T10:05:16-07:00",
          "tree_id": "ca1e3379ca914b619fbc8a3cf5c2e44c363dd441",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/e0b396525002694623c11594a7ae5b9aa7558f74"
        },
        "date": 1664298978984,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37887300769,
            "unit": "ns/op\t        19.19 DeleteSeconds\t        18.65 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8517028224,
            "unit": "ns/op\t         4.312 DeleteSeconds\t         4.164 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6509694925,
            "unit": "ns/op\t         4.304 DeleteSeconds\t         2.145 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "daefad6a35fa41cc7b07304202e84ecebef06e5e",
          "message": "Fixes trivy installation (#905)",
          "timestamp": "2022-09-27T11:08:07-07:00",
          "tree_id": "ac989128f380273c8a56a1ce373aecd9bf94fa6c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/daefad6a35fa41cc7b07304202e84ecebef06e5e"
        },
        "date": 1664302858234,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38145611084,
            "unit": "ns/op\t        19.38 DeleteSeconds\t        18.71 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8720357835,
            "unit": "ns/op\t         4.429 DeleteSeconds\t         4.224 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6559801207,
            "unit": "ns/op\t         4.332 DeleteSeconds\t         2.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a9148583b9e7f9cf02da5f2d7c94e33edc2daca4",
          "message": "Bump peter-evans/create-pull-request from 4.1.1 to 4.1.2 (#901)\n\nBumps [peter-evans/create-pull-request](https://github.com/peter-evans/create-pull-request) from 4.1.1 to 4.1.2.\r\n- [Release notes](https://github.com/peter-evans/create-pull-request/releases)\r\n- [Commits](https://github.com/peter-evans/create-pull-request/compare/18f90432bedd2afd6a825469ffd38aa24712a91d...171dd555b9ab6b18fa02519fdfacbb8bf671e1b4)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: peter-evans/create-pull-request\r\n  dependency-type: direct:production\r\n  update-type: version-update:semver-patch\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2022-09-27T17:39:22-07:00",
          "tree_id": "d1f9c03127f96d78ba76631f6092e5cf0dcfc2c9",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/a9148583b9e7f9cf02da5f2d7c94e33edc2daca4"
        },
        "date": 1664326180942,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36540524667,
            "unit": "ns/op\t        19.02 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8495956005,
            "unit": "ns/op\t         4.287 DeleteSeconds\t         4.168 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6419913048,
            "unit": "ns/op\t         4.251 DeleteSeconds\t         2.129 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fa11cbc367023cd544021efcbc930c378408bad3",
          "message": "fix spelling of global in controller (#911)",
          "timestamp": "2022-09-28T19:54:46-04:00",
          "tree_id": "ab0b037b2e4298e52ff0dc5f1b4f46e32d5cca17",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/fa11cbc367023cd544021efcbc930c378408bad3"
        },
        "date": 1664409883973,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36385375880,
            "unit": "ns/op\t        18.87 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8441456347,
            "unit": "ns/op\t         4.258 DeleteSeconds\t         4.142 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6451974307,
            "unit": "ns/op\t         4.287 DeleteSeconds\t         2.127 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "11ae520e7971acbd6329c39cbb38cc1d98079c39",
          "message": "Merge pull request #896 from vmware-tanzu/dependabot/go_modules/k8s.io/component-base-0.25.2\n\nBump k8s.io/component-base from 0.25.1 to 0.25.2",
          "timestamp": "2022-09-29T11:53:05-06:00",
          "tree_id": "d43add01acdf4f6bdc48d191fee7e5f645b6e951",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/11ae520e7971acbd6329c39cbb38cc1d98079c39"
        },
        "date": 1664474591700,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36508946459,
            "unit": "ns/op\t        18.98 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8471288156,
            "unit": "ns/op\t         4.275 DeleteSeconds\t         4.159 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6480372885,
            "unit": "ns/op\t         4.243 DeleteSeconds\t         2.194 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fab549af1df974c2bf6da85f78cdf65f5ea972a7",
          "message": "Merge pull request #922 from vmware-tanzu/dependabot/docker/golang-1.19.2\n\nBump golang from 1.19.1 to 1.19.2",
          "timestamp": "2022-10-05T10:44:10-06:00",
          "tree_id": "f03b43e19b52018ad1d85c271b538ea91c0d7fb8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/fab549af1df974c2bf6da85f78cdf65f5ea972a7"
        },
        "date": 1664989974242,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35566358307,
            "unit": "ns/op\t        17.97 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8757853592,
            "unit": "ns/op\t         4.542 DeleteSeconds\t         4.173 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6411282210,
            "unit": "ns/op\t         4.235 DeleteSeconds\t         2.133 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "727113ead6dfac282115376102cb8a27d6c4a496",
          "message": "Merge pull request #916 from vmware-tanzu/dependabot/github_actions/peter-evans/create-pull-request-4.1.3\n\nBump peter-evans/create-pull-request from 4.1.2 to 4.1.3",
          "timestamp": "2022-10-05T10:44:48-06:00",
          "tree_id": "7323da0870af9655e4d8aa5ef6a6300d4df87cc1",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/727113ead6dfac282115376102cb8a27d6c4a496"
        },
        "date": 1664989981787,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36562018687,
            "unit": "ns/op\t        18.98 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8491759353,
            "unit": "ns/op\t         4.286 DeleteSeconds\t         4.168 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6457589339,
            "unit": "ns/op\t         4.269 DeleteSeconds\t         2.131 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "110162484be8dc0327d0825b3e47a77543c2a0b5",
          "message": "Merge pull request #915 from vmware-tanzu/nh-add-release-checklist\n\nAdd release checklist Issue Template",
          "timestamp": "2022-10-05T10:37:37-06:00",
          "tree_id": "18213a7a60f173d20dd8f858576f0d744e4e14ac",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/110162484be8dc0327d0825b3e47a77543c2a0b5"
        },
        "date": 1664990012223,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38717397096,
            "unit": "ns/op\t        19.74 DeleteSeconds\t        18.87 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8852140157,
            "unit": "ns/op\t         4.485 DeleteSeconds\t         4.303 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7747934661,
            "unit": "ns/op\t         4.442 DeleteSeconds\t         3.244 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ee66652ee34d4a82304032fd2401ba0f72043ddb",
          "message": "Merge pull request #923 from vmware-tanzu/dependabot/github_actions/actions/checkout-3.1.0\n\nBump actions/checkout from 3.0.2 to 3.1.0",
          "timestamp": "2022-10-05T10:57:54-06:00",
          "tree_id": "399fc5720321a547cb1d5267b066957e9c1583c6",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ee66652ee34d4a82304032fd2401ba0f72043ddb"
        },
        "date": 1664990456668,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36596412099,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8545140212,
            "unit": "ns/op\t         4.330 DeleteSeconds\t         4.168 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6479107852,
            "unit": "ns/op\t         4.277 DeleteSeconds\t         2.159 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "jryan@pivotal.io",
            "name": "John S. Ryan",
            "username": "pivotaljohn"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cd69a5f9e5660fc2679e8f8a39392b702b2bdac7",
          "message": "Use all component info fields in Downward API example (#892)\n\n* Use all component info fields in Downward API example\r\n\r\n* Add example output for api-versions\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\n\r\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\r\nCo-authored-by: John Ryan <jtigger@infosysengr.com>\r\nCo-authored-by: Neil Hickey <nhickey@vmware.com>",
          "timestamp": "2022-10-10T11:13:57-04:00",
          "tree_id": "158d5ea6db42263a7fcc4f3d6d89b7d67790b1ce",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cd69a5f9e5660fc2679e8f8a39392b702b2bdac7"
        },
        "date": 1665415474633,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35472509265,
            "unit": "ns/op\t        17.96 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8485920481,
            "unit": "ns/op\t         4.281 DeleteSeconds\t         4.164 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6464953680,
            "unit": "ns/op\t         4.293 DeleteSeconds\t         2.130 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f17a931d4f01920d4c8aedebe9e138816b78cb42",
          "message": "kctrl: Generating openapi schema from helm values file. (#904)\n\n* Initial Implementation for generating openapi schema from helm values.yaml file\r\n\r\nInitial Implementation for generating openapi schema from helm values.yaml file\r\n\r\n* Fixing test case\r\n\r\nFixing test case\r\n\r\n* Fixing E2e Test case\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Adding unit test cases.\r\n\r\nAdding unit test cases.\r\n\r\n* Update helm_openapi_schema_gen.go\r\n\r\n* Update helm_openapi_schema_gen.go\r\n\r\n* Adopted review comments.\r\n\r\nAdopted review comments.\r\n\r\n* Adopting review comments\r\n\r\nAdopting review comments\r\n\r\n* Update helm_openapi_schema_gen.go\r\n\r\n* Adopting review comments\r\n\r\n* Update release.go\r\n\r\n* Modifying behavior in case > 1 items present in YAML array\r\n\r\nModifying behavior in case > 1 items present in YAML array\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Update helm_openapi_schema_gen.go\r\n\r\n* Update helm_openapi_schema_gen.go\r\n\r\n* Add openapi schema flag (#914)\r\n\r\n* Adding openapi-schema flag to pkg release command\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Adopting review comments\r\n\r\n* Override openapi only when openapi-schema flagg is set to true.\r\n\r\nOverride openapi only when openapi-schema flagg is set to true.\r\n\r\n* Update release.go",
          "timestamp": "2022-10-13T14:24:46+05:30",
          "tree_id": "802abe5a829a46c9b99c6168d41c7d245f59505e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f17a931d4f01920d4c8aedebe9e138816b78cb42"
        },
        "date": 1665652032011,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37027199344,
            "unit": "ns/op\t        19.25 DeleteSeconds\t        17.73 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8650370239,
            "unit": "ns/op\t         4.348 DeleteSeconds\t         4.242 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6533674256,
            "unit": "ns/op\t         4.309 DeleteSeconds\t         2.167 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "dbde3d601f4b543854b8b6eb54681ea4a28c3331",
          "message": "Instantiate coreClient while using reconciler only for the dev command (#939)",
          "timestamp": "2022-10-18T11:06:14+05:30",
          "tree_id": "27ffe6be07f467e5ffdc4556cfa6d77c28e9537b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/dbde3d601f4b543854b8b6eb54681ea4a28c3331"
        },
        "date": 1666072112995,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36822366615,
            "unit": "ns/op\t        18.12 DeleteSeconds\t        18.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8693877425,
            "unit": "ns/op\t         4.444 DeleteSeconds\t         4.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6529529375,
            "unit": "ns/op\t         4.323 DeleteSeconds\t         2.154 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d6b68623f3336af0947fced9a6ae797f275c64bd",
          "message": "Adding a tag flag to override default imgpkg bundle tag (#930)\n\n* Adding a tag flag to override default imgpkg bundle tag\r\n\r\nAdding a tag flag to override default imgpkg bundle tag\r\n\r\n* Adding test case\r\n\r\nAdding test case\r\n\r\n* Adopting comments\r\n\r\nAdopting comments\r\n\r\n* Adopting comments\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Update package_authoring_e2e_test.go\r\n\r\n* Update package_authoring_e2e_test.go",
          "timestamp": "2022-10-18T13:16:33+05:30",
          "tree_id": "26f5f35fb5ad721807bafab4209646a2e77fed66",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d6b68623f3336af0947fced9a6ae797f275c64bd"
        },
        "date": 1666079916973,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37110487963,
            "unit": "ns/op\t        19.42 DeleteSeconds\t        17.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8561202025,
            "unit": "ns/op\t         4.327 DeleteSeconds\t         4.190 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6487215415,
            "unit": "ns/op\t         4.291 DeleteSeconds\t         2.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "995597990c2a8b3d181d4ddce3d33eb704e5fa6d",
          "message": "Adding tag as annotation for pkg repo release command. (#938)\n\n* Adding tag as annotation for pkg repo release command.\r\n\r\nAdding tag as annotation for pkg repo release command.\r\n\r\n* Update package_repo_release_test.go\r\n\r\n* Adopting review comments\r\n\r\nAdopting review comments\r\n\r\n* Adopting reviews\r\n\r\n* Update package_repo_release_test.go",
          "timestamp": "2022-10-18T13:29:49+05:30",
          "tree_id": "8277e097467fcad017d5d0c92ecad2588f4126f7",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/995597990c2a8b3d181d4ddce3d33eb704e5fa6d"
        },
        "date": 1666080610723,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36507449910,
            "unit": "ns/op\t        18.92 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8501330660,
            "unit": "ns/op\t         4.259 DeleteSeconds\t         4.204 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6569984611,
            "unit": "ns/op\t         4.381 DeleteSeconds\t         2.145 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f54408f0bf378c18df9556c81b24a9b3959bcfb2",
          "message": "Refresh package install after pausing it (#929)\n\nAfter a package install is paused successfully, the observedGeneration is updated resulting in conflicts while updating the package install",
          "timestamp": "2022-10-19T09:59:29+05:30",
          "tree_id": "31b768171a72e627aa5ab793a1873a0ef9403dbc",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f54408f0bf378c18df9556c81b24a9b3959bcfb2"
        },
        "date": 1666154391938,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35688087668,
            "unit": "ns/op\t        18.16 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8460735008,
            "unit": "ns/op\t         4.270 DeleteSeconds\t         4.151 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6456875580,
            "unit": "ns/op\t         4.293 DeleteSeconds\t         2.122 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "050ba8ab3e592c814b53947f44ce418fe9ce1a8b",
          "message": "Add missing package condition while waiting for app pause (#944)\n\nWhen a package install for which the installed package version is now removed from the cluster, the app cr is never paused and so we need to check for the failing condition in package install",
          "timestamp": "2022-10-20T11:20:11+05:30",
          "tree_id": "79c72d962a29e83e06ce40e547fa34bf702d6e7d",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/050ba8ab3e592c814b53947f44ce418fe9ce1a8b"
        },
        "date": 1666245624975,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36478840671,
            "unit": "ns/op\t        18.94 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8469204544,
            "unit": "ns/op\t         4.268 DeleteSeconds\t         4.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6470042474,
            "unit": "ns/op\t         4.306 DeleteSeconds\t         2.125 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "498483538975cb854bd1164ce45b5ea477edadd1",
          "message": "Refactor app/pkg init logic (#913)\n\n* Scope logic dealing with vendir.yml to VendirConfig\r\n\r\n* Refactor how Fetch is configured. Move build scoped fucntions to build.\r\n\r\n* Split VendirStep into VendirRunner and VendirConfigBuilder. Remove deadcode\r\n\r\n* Refactor git step. Remove deadcode.\r\n\r\n* Refactor GithubStep\r\n\r\n* Remove unused create step and refactor necessarily in app init\r\n\r\n* Refactor GitStep\r\n\r\n* Dedup app init logic. configureAppBuild => getAppBuildName\r\n\r\n* Refactor TemplateStep\r\n\r\n* Remove step interface and move Build interface\r\n\r\n* Refactoring package init\r\n\r\n* Add missing deffered cleanup\r\n\r\n* Remove file_utils. Use builins instead\r\n\r\n* Add check while running vendir sync to handle local directory case\r\n\r\n* Merge VendirCOnfiguration into FetchConfiguation\r\n\r\n* Move PackageBuild and AppBuild to buildconfigs\r\n\r\n* Move vendir config and related objects to sources package\r\n\r\n* Move init command files out of a separate package\r\n\r\n* Remove dependency of annotation for storing fetch mode\r\n\r\n* Moving constants to appropriate locations\r\n\r\n* Returning a non-pointer value from GetExport as it is always dereferenced. Removing stale comments\r\n\r\n* Stricter checks before running vendir sync. Making not exists check on files cleaner\r\n\r\n* Move source specific configuration to source.go. fetch.go => source.go. Refactoring init files\r\n\r\n* Remove unnecessary dependency on carvel-kapp-controller/.../exec package\r\n\r\n* Use vendirConfig.Contents while configuring package build instead of passing contents down the function tree\r\n\r\n* Move logic for initialising deploy section to build interface. Remove duplicate dependencies\r\n\r\n* TemplateConfiguration => Template. Othere refactoring.",
          "timestamp": "2022-10-20T16:55:29+05:30",
          "tree_id": "be91724ed25ef9a772beb32ae6469a6f5748610f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/498483538975cb854bd1164ce45b5ea477edadd1"
        },
        "date": 1666265892111,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36992952971,
            "unit": "ns/op\t        19.23 DeleteSeconds\t        17.71 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8593540756,
            "unit": "ns/op\t         4.352 DeleteSeconds\t         4.196 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6523328603,
            "unit": "ns/op\t         4.321 DeleteSeconds\t         2.154 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "519ecab2f1f8b2a5b669cbbe92e1d69ec9567a86",
          "message": "Use carvel setup action in test-kctrl-gh (#945)",
          "timestamp": "2022-10-20T21:36:27+05:30",
          "tree_id": "b361b81c01390e68db50d5b379db2b7e7e52df3c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/519ecab2f1f8b2a5b669cbbe92e1d69ec9567a86"
        },
        "date": 1666282624692,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35407020195,
            "unit": "ns/op\t        17.89 DeleteSeconds\t        17.45 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8451640995,
            "unit": "ns/op\t         4.258 DeleteSeconds\t         4.154 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6385936936,
            "unit": "ns/op\t         4.222 DeleteSeconds\t         2.124 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "544229852ae693682a456bbd5ae793444b1bc9a5",
          "message": "Merge pull request #948 from vmware-tanzu/dependabot/github_actions/peter-evans/create-pull-request-4.2.0\n\nBump peter-evans/create-pull-request from 4.1.3 to 4.2.0",
          "timestamp": "2022-10-20T12:26:41-06:00",
          "tree_id": "2ce085cd5ca39791f85a99166f932fbbf402c77c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/544229852ae693682a456bbd5ae793444b1bc9a5"
        },
        "date": 1666291020295,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35348109177,
            "unit": "ns/op\t        17.84 DeleteSeconds\t        17.45 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8519258415,
            "unit": "ns/op\t         4.243 DeleteSeconds\t         4.234 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6393130921,
            "unit": "ns/op\t         4.225 DeleteSeconds\t         2.129 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9fc702b4f48eba7ee61003636d1f176f5881b1b5",
          "message": "Merge pull request #940 from vmware-tanzu/dependabot/github_actions/slackapi/slack-github-action-1.23.0\n\nBump slackapi/slack-github-action from 1.22.0 to 1.23.0",
          "timestamp": "2022-10-20T12:27:12-06:00",
          "tree_id": "0d00e057fa95d28e6dfd68a3ae51e7c0f2d0aac8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/9fc702b4f48eba7ee61003636d1f176f5881b1b5"
        },
        "date": 1666291231086,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38261518589,
            "unit": "ns/op\t        19.45 DeleteSeconds\t        18.76 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8655025404,
            "unit": "ns/op\t         4.398 DeleteSeconds\t         4.202 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7620285596,
            "unit": "ns/op\t         4.381 DeleteSeconds\t         3.179 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ea3dd3fd2fa92267d14193da7c29a3212fd55ea6",
          "message": "Merge pull request #937 from vmware-tanzu/dependabot/go_modules/k8s.io/component-base-0.25.3\n\nBump k8s.io/component-base from 0.25.2 to 0.25.3",
          "timestamp": "2022-10-20T13:32:03-06:00",
          "tree_id": "55190c0042bdaae6a5e0782077b6127e2c1c2199",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ea3dd3fd2fa92267d14193da7c29a3212fd55ea6"
        },
        "date": 1666294965198,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35643345777,
            "unit": "ns/op\t        18.14 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8514473522,
            "unit": "ns/op\t         4.298 DeleteSeconds\t         4.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6415424554,
            "unit": "ns/op\t         4.249 DeleteSeconds\t         2.126 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "joaod@vmware.com",
            "name": "João Pereira",
            "username": "joaopapereira"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "054ad9af5337ccaced3998a0b69c0a6d0e993dc9",
          "message": "Some tweaks on the release process (#950)\n\nSigned-off-by: João Pereira <joaod@vmware.com>\r\n\r\nSigned-off-by: João Pereira <joaod@vmware.com>",
          "timestamp": "2022-10-21T14:14:59-04:00",
          "tree_id": "4f298dc01b293b653f24e18235604200f9ea99b3",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/054ad9af5337ccaced3998a0b69c0a6d0e993dc9"
        },
        "date": 1666376714887,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35452291265,
            "unit": "ns/op\t        17.92 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8483988222,
            "unit": "ns/op\t         4.283 DeleteSeconds\t         4.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6435505014,
            "unit": "ns/op\t         4.259 DeleteSeconds\t         2.138 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "joaod@vmware.com",
            "name": "João Pereira",
            "username": "joaopapereira"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c6bd2ecae7374c1b07102445b0e3d327994e1888",
          "message": "Merge pull request #942 from vmware-tanzu/bump-deps\n\nBump dependencies",
          "timestamp": "2022-10-25T13:18:37-05:00",
          "tree_id": "0d991acee53657cf9659bae2847883572b7dbaa0",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c6bd2ecae7374c1b07102445b0e3d327994e1888"
        },
        "date": 1666722653417,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38009009769,
            "unit": "ns/op\t        19.18 DeleteSeconds\t        18.78 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8597387378,
            "unit": "ns/op\t         4.337 DeleteSeconds\t         4.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6550127007,
            "unit": "ns/op\t         4.328 DeleteSeconds\t         2.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "cppforlife@gmail.com",
            "name": "Dmitriy Kalinin",
            "username": "cppforlife"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "eb31a8f5338c82f068627b72579b7ff3d8bc33e4",
          "message": "Merge pull request #909 from vmware-tanzu/vendir-caching\n\nActivate caching of images and bundles",
          "timestamp": "2022-10-26T15:08:21-04:00",
          "tree_id": "acd93c223997ac0af23d3b3f7816bd65eaa5139e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/eb31a8f5338c82f068627b72579b7ff3d8bc33e4"
        },
        "date": 1666811933116,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36602453223,
            "unit": "ns/op\t        18.98 DeleteSeconds\t        17.58 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8489019909,
            "unit": "ns/op\t         4.278 DeleteSeconds\t         4.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6464989333,
            "unit": "ns/op\t         4.269 DeleteSeconds\t         2.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c834810997d4b593a0f8b4874e9eb78a0261d86f",
          "message": "Reword informational text in authoring commands (#956)",
          "timestamp": "2022-10-27T03:07:45+05:30",
          "tree_id": "1f63f7c27ffd318d8a266bb252221ddc56409091",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c834810997d4b593a0f8b4874e9eb78a0261d86f"
        },
        "date": 1666821042887,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36012471238,
            "unit": "ns/op\t        18.29 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9738439101,
            "unit": "ns/op\t         5.445 DeleteSeconds\t         4.240 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6580555074,
            "unit": "ns/op\t         4.338 DeleteSeconds\t         2.179 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3e5633a5b5352fa45d613f23b7d85df60e980ee7",
          "message": "Fix nil check for GithubRelease fetch mode (#972)",
          "timestamp": "2022-11-15T14:52:16+05:30",
          "tree_id": "9658ad65a9631e72732cfb5e2efb0e830c215c69",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/3e5633a5b5352fa45d613f23b7d85df60e980ee7"
        },
        "date": 1668504796193,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36675198703,
            "unit": "ns/op\t        19.06 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8506583603,
            "unit": "ns/op\t         4.280 DeleteSeconds\t         4.176 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6440231966,
            "unit": "ns/op\t         4.257 DeleteSeconds\t         2.144 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "joaod@vmware.com",
            "name": "João Pereira",
            "username": "joaopapereira"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8c154f09d42f05ad07cef610586db474404c58d6",
          "message": "Merge pull request #955 from vmware-tanzu/bump-x-text\n\nBump x/text to version 0.3.8",
          "timestamp": "2022-11-15T16:15:38-06:00",
          "tree_id": "332895010002acc587c65be9337931612ff550fd",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8c154f09d42f05ad07cef610586db474404c58d6"
        },
        "date": 1668551173260,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35486439516,
            "unit": "ns/op\t        17.94 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8636993019,
            "unit": "ns/op\t         4.359 DeleteSeconds\t         4.238 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6437324410,
            "unit": "ns/op\t         4.258 DeleteSeconds\t         2.138 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7e318f4b5567e2f398fe4b06e8d0d54cef81db10",
          "message": "Set wait-check-interval to 1s for benchmark tests (#979)\n\nWith kapp v0.54.0 wait-check-interval is set to 3s which increases the time to wait for package repositories, hence we need to explicitly set it to 1s",
          "timestamp": "2022-11-24T15:54:38+05:30",
          "tree_id": "7a756cff29234e5722c486a090cba4bc99e4e491",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7e318f4b5567e2f398fe4b06e8d0d54cef81db10"
        },
        "date": 1669286203373,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36244518875,
            "unit": "ns/op\t        18.50 DeleteSeconds\t        17.67 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8608082542,
            "unit": "ns/op\t         4.351 DeleteSeconds\t         4.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6539948394,
            "unit": "ns/op\t         4.301 DeleteSeconds\t         2.189 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4c7230e75c03f29e5ade77f0fc0d9600f7801daf",
          "message": "Set kapp wait-check-interval to 1s for e2e tests (#985)",
          "timestamp": "2022-11-24T19:26:15+05:30",
          "tree_id": "1f2a101463a2d543461dcad06aafbd27def89dc7",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4c7230e75c03f29e5ade77f0fc0d9600f7801daf"
        },
        "date": 1669298881986,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37804620321,
            "unit": "ns/op\t        19.11 DeleteSeconds\t        18.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8579286689,
            "unit": "ns/op\t         4.328 DeleteSeconds\t         4.200 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6586729283,
            "unit": "ns/op\t         4.340 DeleteSeconds\t         2.197 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8456447f4ca4005b4064aee06fd3e12c35c2d71f",
          "message": "Install imgpkg before running release workflow (#986)",
          "timestamp": "2022-11-24T22:36:44+05:30",
          "tree_id": "8b02eb4fde11a939aa18efd956aea8b21ac4b51f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8456447f4ca4005b4064aee06fd3e12c35c2d71f"
        },
        "date": 1669310229308,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37190168135,
            "unit": "ns/op\t        19.64 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8533807384,
            "unit": "ns/op\t         4.314 DeleteSeconds\t         4.175 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6440344626,
            "unit": "ns/op\t         4.262 DeleteSeconds\t         2.132 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "73c7c61aa13e796f681f314078dd2b8fcb026b65",
          "message": "Bump golang.org/x/net (#987)",
          "timestamp": "2022-11-25T08:51:52-05:00",
          "tree_id": "54cfd65620c303e67ee41f69eb1b074195afa220",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/73c7c61aa13e796f681f314078dd2b8fcb026b65"
        },
        "date": 1669385013753,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36777247629,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.74 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8652964711,
            "unit": "ns/op\t         4.379 DeleteSeconds\t         4.215 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6552901502,
            "unit": "ns/op\t         4.342 DeleteSeconds\t         2.164 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "98549719+carvel-bot@users.noreply.github.com",
            "name": "Carvel Bot",
            "username": "carvel-bot"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "08265e3a1b9952c7ec55c49a8bfea76e4af50487",
          "message": "Bump dependencies (#988)",
          "timestamp": "2022-11-26T15:05:04+05:30",
          "tree_id": "341b1ac54c601ec1f8e104056cb5b3e12beeaf9f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/08265e3a1b9952c7ec55c49a8bfea76e4af50487"
        },
        "date": 1669455943838,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36859495158,
            "unit": "ns/op\t        19.29 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8569163033,
            "unit": "ns/op\t         4.333 DeleteSeconds\t         4.190 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6468223406,
            "unit": "ns/op\t         4.281 DeleteSeconds\t         2.142 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "aaronshurley@users.noreply.github.com",
            "name": "Aaron Hurley",
            "username": "aaronshurley"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cc210091e62457fe479e0edad4e060d523ab5397",
          "message": "Introduce inclusive naming check and update existing language (#977)\n\nWe are following the Inclusive Naming Initiative's guidance as that is what the\r\nCNCF supports.",
          "timestamp": "2022-11-28T09:27:14-06:00",
          "tree_id": "bf5578d2ff7133d57dcf3110b8411b021eb49871",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cc210091e62457fe479e0edad4e060d523ab5397"
        },
        "date": 1669649906367,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35530732861,
            "unit": "ns/op\t        17.98 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8537122568,
            "unit": "ns/op\t         4.329 DeleteSeconds\t         4.169 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6464156443,
            "unit": "ns/op\t         4.253 DeleteSeconds\t         2.167 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1b8b19aa339d8885ccc2821668ce15facba34371",
          "message": "Merge pull request #951 from vmware-tanzu/dependabot/go_modules/github.com/stretchr/testify-1.8.1\n\nBump github.com/stretchr/testify from 1.8.0 to 1.8.1",
          "timestamp": "2022-11-28T16:41:53-07:00",
          "tree_id": "52dd8d48edfa4431b55d281603ba1540107f2212",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/1b8b19aa339d8885ccc2821668ce15facba34371"
        },
        "date": 1669679516820,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36432001207,
            "unit": "ns/op\t        18.82 DeleteSeconds\t        17.57 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8478863384,
            "unit": "ns/op\t         4.265 DeleteSeconds\t         4.172 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6555828568,
            "unit": "ns/op\t         4.321 DeleteSeconds\t         2.160 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "1843f99c53df901c1287612d15c991a49e28ff39",
          "message": "Add trivy scan for kctrl (#969)\n\n* Fix trivy scan: support multiline strings with set-output\r\n\r\n* Add trivy scan for kctrl\r\n\r\n* Use -o=json instead of --to-json\r\n\r\nDownload only release.yml instead of all release artefacts",
          "timestamp": "2022-11-30T01:52:57+05:30",
          "tree_id": "d4ca2d5fcb0cd71465c4e63e79a656b72b0a15a0",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/1843f99c53df901c1287612d15c991a49e28ff39"
        },
        "date": 1669754189093,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38799798702,
            "unit": "ns/op\t        19.89 DeleteSeconds\t        18.85 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 9765535497,
            "unit": "ns/op\t         5.442 DeleteSeconds\t         4.254 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6638290975,
            "unit": "ns/op\t         4.365 DeleteSeconds\t         2.212 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "43b643866b75f8c19e1144e330069550b47d89f0",
          "message": "Remove version constraint on kind tests (#992)",
          "timestamp": "2022-11-29T16:04:15-05:00",
          "tree_id": "46d7ca1615ae09973d416f94c90d97d4f5711f77",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/43b643866b75f8c19e1144e330069550b47d89f0"
        },
        "date": 1669756484735,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36791136865,
            "unit": "ns/op\t        18.97 DeleteSeconds\t        17.78 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8504645702,
            "unit": "ns/op\t         4.283 DeleteSeconds\t         4.175 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6446492949,
            "unit": "ns/op\t         4.263 DeleteSeconds\t         2.146 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7b606d03033d4d70168d8be45b62e8cebd9f8658",
          "message": "Bump cobrautil in cli (#994)",
          "timestamp": "2022-12-01T19:04:14+05:30",
          "tree_id": "3a4acb2644ee52125990e1926da977940043535a",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7b606d03033d4d70168d8be45b62e8cebd9f8658"
        },
        "date": 1669902292520,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36585078137,
            "unit": "ns/op\t        18.99 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8554517846,
            "unit": "ns/op\t         4.331 DeleteSeconds\t         4.178 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6597616553,
            "unit": "ns/op\t         4.385 DeleteSeconds\t         2.161 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "ThomasVitale@users.noreply.github.com",
            "name": "Thomas Vitale",
            "username": "ThomasVitale"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9693c1b246b6b94f8c638b68d05f829274df5bce",
          "message": "Improve kctrl dev command description (#1000)\n\n* Update the description for the kctrl dev command.\r\n* Update the description for the -f flag.\r\n\r\nFixes gh-964",
          "timestamp": "2022-12-02T13:37:17+05:30",
          "tree_id": "8bd01cc64947c4cd1b5e6ddc6b0669720d0bc37e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/9693c1b246b6b94f8c638b68d05f829274df5bce"
        },
        "date": 1669969038658,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36902734404,
            "unit": "ns/op\t        19.32 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8471431941,
            "unit": "ns/op\t         4.272 DeleteSeconds\t         4.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6454251932,
            "unit": "ns/op\t         4.270 DeleteSeconds\t         2.140 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2c51b3ec201869f9a43e740b79320f709ea4aa56",
          "message": "Merge pull request #993 from vmware-tanzu/trivy-scan\n\nAdd space for if condition in trivy-scan",
          "timestamp": "2022-12-05T11:20:30-07:00",
          "tree_id": "3273c8cc4bb05d4e747d44cbc437407f6b3f6470",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2c51b3ec201869f9a43e740b79320f709ea4aa56"
        },
        "date": 1670265091092,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36413487021,
            "unit": "ns/op\t        18.53 DeleteSeconds\t        17.83 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8562732806,
            "unit": "ns/op\t         4.337 DeleteSeconds\t         4.176 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6556021562,
            "unit": "ns/op\t         4.325 DeleteSeconds\t         2.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7b9994cfd0bd1d9962f2989e429a22d6d0b8ba80",
          "message": "Merge pull request #973 from vmware-tanzu/dependabot/go_modules/k8s.io/component-base-0.25.4\n\nBump k8s.io/component-base from 0.25.3 to 0.25.4",
          "timestamp": "2022-12-05T11:37:36-07:00",
          "tree_id": "7dda6bda9705ae45f7bc4ad51106dabc8f512330",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/7b9994cfd0bd1d9962f2989e429a22d6d0b8ba80"
        },
        "date": 1670266085555,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35573437167,
            "unit": "ns/op\t        18.04 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8584019388,
            "unit": "ns/op\t         4.341 DeleteSeconds\t         4.172 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6499393554,
            "unit": "ns/op\t         4.314 DeleteSeconds\t         2.143 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "753ec43be7f18ee99332d59d18aee1aa3279223a",
          "message": "Merge pull request #968 from vmware-tanzu/dependabot/go_modules/sigs.k8s.io/controller-runtime-0.13.1\n\nBump sigs.k8s.io/controller-runtime from 0.13.0 to 0.13.1",
          "timestamp": "2022-12-05T11:36:47-07:00",
          "tree_id": "2ee7fd15fe1778eaa0ee01584ad0604bb8dbc6a8",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/753ec43be7f18ee99332d59d18aee1aa3279223a"
        },
        "date": 1670266181839,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37898853087,
            "unit": "ns/op\t        20.07 DeleteSeconds\t        17.78 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8679795060,
            "unit": "ns/op\t         4.393 DeleteSeconds\t         4.234 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6636009776,
            "unit": "ns/op\t         4.402 DeleteSeconds\t         2.179 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "94a49386256191f702e03063f10c3e19583c4d29",
          "message": "Merge pull request #975 from vmware-tanzu/dependabot/go_modules/k8s.io/kube-aggregator-0.22.16\n\nBump k8s.io/kube-aggregator from 0.22.15 to 0.22.16",
          "timestamp": "2022-12-05T14:26:16-07:00",
          "tree_id": "30750434d8873c71a1d3234b470c1d5b19189a96",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/94a49386256191f702e03063f10c3e19583c4d29"
        },
        "date": 1670276220860,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36718406683,
            "unit": "ns/op\t        19.01 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8741934646,
            "unit": "ns/op\t         4.504 DeleteSeconds\t         4.195 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6514998013,
            "unit": "ns/op\t         4.302 DeleteSeconds\t         2.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cd5cbe7b4495c999ff7845fe9239873820e1d287",
          "message": "Merge pull request #990 from vmware-tanzu/dependabot/github_actions/peter-evans/create-pull-request-4.2.3\n\nBump peter-evans/create-pull-request from 4.2.0 to 4.2.3",
          "timestamp": "2022-12-05T14:26:48-07:00",
          "tree_id": "482938f78961904e55fba394358a40982a6e1cbe",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/cd5cbe7b4495c999ff7845fe9239873820e1d287"
        },
        "date": 1670276307926,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35850500150,
            "unit": "ns/op\t        18.14 DeleteSeconds\t        17.65 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8760887697,
            "unit": "ns/op\t         4.466 DeleteSeconds\t         4.225 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6553020932,
            "unit": "ns/op\t         4.314 DeleteSeconds\t         2.191 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f0264678917dded01dd1936839176f3ee6d3e34d",
          "message": "Merge pull request #1003 from vmware-tanzu/dependabot/github_actions/softprops/action-gh-release-0.1.15\n\nBump softprops/action-gh-release from 0.1.14 to 0.1.15",
          "timestamp": "2022-12-05T14:29:47-07:00",
          "tree_id": "45a528f4873db135363cb82a687c237cf7f9cc90",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f0264678917dded01dd1936839176f3ee6d3e34d"
        },
        "date": 1670276391210,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36597987247,
            "unit": "ns/op\t        18.96 DeleteSeconds\t        17.58 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8535039727,
            "unit": "ns/op\t         4.309 DeleteSeconds\t         4.185 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6470175329,
            "unit": "ns/op\t         4.267 DeleteSeconds\t         2.148 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4f0cda217970340a440dcba5989b996e72e67bfa",
          "message": "Merge pull request #967 from vmware-tanzu/dependabot/github_actions/benchmark-action/github-action-benchmark-1.15.0\n\nBump benchmark-action/github-action-benchmark from 1.14.0 to 1.15.0",
          "timestamp": "2022-12-05T14:27:16-07:00",
          "tree_id": "ed6fa708fa128a20ed9a0f8fcdfcc7d583efc87f",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4f0cda217970340a440dcba5989b996e72e67bfa"
        },
        "date": 1670276418605,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37379469275,
            "unit": "ns/op\t        19.58 DeleteSeconds\t        17.75 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8817282335,
            "unit": "ns/op\t         4.504 DeleteSeconds\t         4.218 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6643089335,
            "unit": "ns/op\t         4.376 DeleteSeconds\t         2.203 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6706205ca85e3f4e9dc305fe47fc581f07337ca3",
          "message": "Set seccompProfile type for both containers (#999)\n\nboth = kapp-controller and kapp-controller-sidecarexec",
          "timestamp": "2022-12-07T11:19:55-05:00",
          "tree_id": "7c909eea6d98ba305e69514cff15cbb16a224362",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/6706205ca85e3f4e9dc305fe47fc581f07337ca3"
        },
        "date": 1670430625745,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37045289274,
            "unit": "ns/op\t        19.48 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8537708318,
            "unit": "ns/op\t         4.311 DeleteSeconds\t         4.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6477540185,
            "unit": "ns/op\t         4.285 DeleteSeconds\t         2.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "e44812bd0fbc8f67c9a692163c9d08fafe4c6287",
          "message": "Update secure namespace flag name in a hint (#1025)",
          "timestamp": "2022-12-13T20:30:20+05:30",
          "tree_id": "25442dc53db5512f86bdbf208f9671e9b43462e2",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/e44812bd0fbc8f67c9a692163c9d08fafe4c6287"
        },
        "date": 1670944239323,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36597911010,
            "unit": "ns/op\t        19.03 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8487722588,
            "unit": "ns/op\t         4.290 DeleteSeconds\t         4.154 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6426696745,
            "unit": "ns/op\t         4.248 DeleteSeconds\t         2.136 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0e47da970e9fd015ea7022eb2f29ea39198557f1",
          "message": "Bump golang.org/x/net in cli (#1026)",
          "timestamp": "2022-12-13T20:31:16+05:30",
          "tree_id": "15f3468106647258e531eefc69e6924f446a9763",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/0e47da970e9fd015ea7022eb2f29ea39198557f1"
        },
        "date": 1670944298581,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37042031129,
            "unit": "ns/op\t        19.02 DeleteSeconds\t        17.98 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8476003872,
            "unit": "ns/op\t         4.284 DeleteSeconds\t         4.151 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6473996595,
            "unit": "ns/op\t         4.284 DeleteSeconds\t         2.148 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "0956ea38597decadc7f98f85d15ed9eb1cfad6e3",
          "message": "Merge pull request #1024 from vmware-tanzu/bump-go-1.19.3\n\nBump go 1.19.3",
          "timestamp": "2022-12-13T12:07:07-07:00",
          "tree_id": "7448fcc49d62f5bb7af6ced6ba2363a010efdde9",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/0956ea38597decadc7f98f85d15ed9eb1cfad6e3"
        },
        "date": 1670959119070,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37105832440,
            "unit": "ns/op\t        18.40 DeleteSeconds\t        18.65 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8585458630,
            "unit": "ns/op\t         4.328 DeleteSeconds\t         4.207 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6543104580,
            "unit": "ns/op\t         4.300 DeleteSeconds\t         2.192 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "843e3cd4fd9d47f8a0fcc32438e2f5340020ac0b",
          "message": "Merge pull request #1007 from vmware-tanzu/dependabot/go_modules/github.com/vmware-tanzu/carvel-vendir-0.30.1\n\nBump github.com/vmware-tanzu/carvel-vendir from 0.30.0 to 0.30.1",
          "timestamp": "2022-12-13T12:17:39-07:00",
          "tree_id": "3ac9e1bc3f5b48b3a2dc3736ab847660fd08ac32",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/843e3cd4fd9d47f8a0fcc32438e2f5340020ac0b"
        },
        "date": 1670959650956,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36733316865,
            "unit": "ns/op\t        18.97 DeleteSeconds\t        17.73 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8450875053,
            "unit": "ns/op\t         4.255 DeleteSeconds\t         4.158 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6425604774,
            "unit": "ns/op\t         4.234 DeleteSeconds\t         2.141 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "d54145f59b5de8e3abc7d9d1f8bba56ff1cf4f1e",
          "message": "Merge pull request #991 from vmware-tanzu/fix-release-note-paths\n\nRemove file path on automated release note generated",
          "timestamp": "2022-12-13T12:25:01-07:00",
          "tree_id": "df7954b39594c9dcc3cde305fb30fb537e5b5b36",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/d54145f59b5de8e3abc7d9d1f8bba56ff1cf4f1e"
        },
        "date": 1670960089339,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35424262384,
            "unit": "ns/op\t        17.88 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8535283813,
            "unit": "ns/op\t         4.337 DeleteSeconds\t         4.157 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6386190305,
            "unit": "ns/op\t         4.228 DeleteSeconds\t         2.119 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "49b9ce1ef81388ac1b21fe9a1f8830b166295008",
          "message": "Merge pull request #981 from vmware-tanzu/dependabot/github_actions/reviewdog/action-misspell-1.12.3\n\nBump reviewdog/action-misspell from 1.12.2 to 1.12.3",
          "timestamp": "2022-12-13T12:22:05-07:00",
          "tree_id": "fc90ce78e1a6c4c809f4253eca37ee9746db75bd",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/49b9ce1ef81388ac1b21fe9a1f8830b166295008"
        },
        "date": 1670960136335,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38806580582,
            "unit": "ns/op\t        19.96 DeleteSeconds\t        18.77 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8806817597,
            "unit": "ns/op\t         4.487 DeleteSeconds\t         4.259 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6642678445,
            "unit": "ns/op\t         4.379 DeleteSeconds\t         2.195 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4fd75f560f198cf43532fa396c43f53f3210a922",
          "message": "Merge pull request #1022 from vmware-tanzu/dependabot/github_actions/actions/stale-6.0.1\n\nBump actions/stale from 5.1.1 to 6.0.1",
          "timestamp": "2022-12-13T12:26:41-07:00",
          "tree_id": "c6f9ecb14acbc2490461cf98c139b8ecb46d2664",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4fd75f560f198cf43532fa396c43f53f3210a922"
        },
        "date": 1670960234807,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36663957021,
            "unit": "ns/op\t        19.07 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8734708033,
            "unit": "ns/op\t         4.495 DeleteSeconds\t         4.184 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6508145049,
            "unit": "ns/op\t         4.276 DeleteSeconds\t         2.183 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c9dffa605d4e2995d68ccb9d27a8ffa97b83a16e",
          "message": "Merge pull request #1020 from vmware-tanzu/dependabot/go_modules/k8s.io/kube-aggregator-0.22.17\n\nBump k8s.io/kube-aggregator from 0.22.16 to 0.22.17",
          "timestamp": "2022-12-13T12:25:38-07:00",
          "tree_id": "245bd1a97bdaaae6f7b98ab9a6027f9f9428ed88",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c9dffa605d4e2995d68ccb9d27a8ffa97b83a16e"
        },
        "date": 1670960262855,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37195277789,
            "unit": "ns/op\t        19.48 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8651068273,
            "unit": "ns/op\t         4.391 DeleteSeconds\t         4.203 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6531666791,
            "unit": "ns/op\t         4.318 DeleteSeconds\t         2.160 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "4b2f93d8eec3dd19aa10626b354f091b0f6e704e",
          "message": "Merge pull request #1009 from vmware-tanzu/bump-dependencies\n\nBump dependencies",
          "timestamp": "2022-12-13T12:29:24-07:00",
          "tree_id": "150f932ffed1f3dec0861815a488bf0f208cec1b",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/4b2f93d8eec3dd19aa10626b354f091b0f6e704e"
        },
        "date": 1670960373798,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36724448368,
            "unit": "ns/op\t        19.21 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8538129075,
            "unit": "ns/op\t         4.315 DeleteSeconds\t         4.181 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6432623677,
            "unit": "ns/op\t         4.252 DeleteSeconds\t         2.138 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "422f186917e3d9527fd7a253794dd87819a819e5",
          "message": "Merge pull request #982 from vmware-tanzu/dependabot/github_actions/actions/add-to-project-0.4.0\n\nBump actions/add-to-project from 0.3.0 to 0.4.0",
          "timestamp": "2022-12-13T14:00:37-07:00",
          "tree_id": "26f217d524fca4b64d9e81a95b26b60b1ed0dd3e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/422f186917e3d9527fd7a253794dd87819a819e5"
        },
        "date": 1670965952190,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37040713076,
            "unit": "ns/op\t        18.10 DeleteSeconds\t        18.88 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8564667602,
            "unit": "ns/op\t         4.321 DeleteSeconds\t         4.197 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6518576686,
            "unit": "ns/op\t         4.305 DeleteSeconds\t         2.164 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "30009b4f1c890528166e1c66ad0a264adb336720",
          "message": "Merge pull request #1018 from vmware-tanzu/dependabot/go_modules/k8s.io/component-base-0.25.5\n\nBump k8s.io/component-base from 0.25.4 to 0.25.5",
          "timestamp": "2022-12-13T14:01:03-07:00",
          "tree_id": "89f245af49a5013ac98013cb5febc1888b6b462d",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/30009b4f1c890528166e1c66ad0a264adb336720"
        },
        "date": 1670966072764,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37727193699,
            "unit": "ns/op\t        18.64 DeleteSeconds\t        19.03 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8727134120,
            "unit": "ns/op\t         4.446 DeleteSeconds\t         4.225 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6648066372,
            "unit": "ns/op\t         4.385 DeleteSeconds\t         2.204 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "ryanjo@vmware.com",
            "name": "John Ryan",
            "username": "pivotaljohn"
          },
          "committer": {
            "email": "ryanjo@vmware.com",
            "name": "John Ryan",
            "username": "pivotaljohn"
          },
          "distinct": true,
          "id": "1afc52e00b1cda77886a22bccad051b351cd337a",
          "message": "Fix typo in \"release process\" workflow\n\nSigned-off-by: Neil Hickey <nhickey@vmware.com>\nSigned-off-by: Varsha Munishwar <vmunishwar@vmware.com>",
          "timestamp": "2022-12-13T14:26:10-08:00",
          "tree_id": "0b618c6020c03f79fe2fc2b228f6e1110f109e46",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/1afc52e00b1cda77886a22bccad051b351cd337a"
        },
        "date": 1670971119236,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35831233790,
            "unit": "ns/op\t        18.31 DeleteSeconds\t        17.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8548395724,
            "unit": "ns/op\t         4.321 DeleteSeconds\t         4.172 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6446645043,
            "unit": "ns/op\t         4.261 DeleteSeconds\t         2.144 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "6b0843a3db36784fea47d34307deb994de107cb2",
          "message": "Fix app status diff, to ensure that deploy output is deduped (#1013)\n\n* Ensure that deploy output is cached to fix deploy output diffing\r\n\r\n* Harden package install tests to ensure that app statuses are diffed adequately",
          "timestamp": "2022-12-14T11:57:43+05:30",
          "tree_id": "c6ef03c105efb1f9b7a237df184d227ba05dba33",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/6b0843a3db36784fea47d34307deb994de107cb2"
        },
        "date": 1670999868051,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37025596166,
            "unit": "ns/op\t        19.43 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8477347028,
            "unit": "ns/op\t         4.284 DeleteSeconds\t         4.151 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6430730626,
            "unit": "ns/op\t         4.260 DeleteSeconds\t         2.127 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "c5f626e3f6a85d689621ae25f4edb0e57d8c3e55",
          "message": "Merge pull request #989 from vmware-tanzu/update-test-workflows\n\nDo not run/stop workflows that are not required",
          "timestamp": "2022-12-14T10:14:00-07:00",
          "tree_id": "36434af2d232045af999f46177502dfddd0550f5",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/c5f626e3f6a85d689621ae25f4edb0e57d8c3e55"
        },
        "date": 1671038807250,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36468908044,
            "unit": "ns/op\t        18.64 DeleteSeconds\t        17.76 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8784005190,
            "unit": "ns/op\t         4.383 DeleteSeconds\t         4.340 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7732190397,
            "unit": "ns/op\t         4.480 DeleteSeconds\t         3.201 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fea2c7cae46f3809ec9ab8e9d27161b713db10f5",
          "message": "Merge pull request #1033 from vmware-tanzu/nh-bump-redis-chart-latest\n\nBump testing redis chart to latest [17.3.17]",
          "timestamp": "2022-12-22T10:21:13-07:00",
          "tree_id": "6927a1a1ccb7c3a0509eed865e725242fde11119",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/fea2c7cae46f3809ec9ab8e9d27161b713db10f5"
        },
        "date": 1671730307074,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36843259569,
            "unit": "ns/op\t        19.30 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8586132039,
            "unit": "ns/op\t         4.363 DeleteSeconds\t         4.174 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6814628307,
            "unit": "ns/op\t         4.560 DeleteSeconds\t         2.206 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f88670e53d679a1ae16c2eccb9abb4397af2d0bd",
          "message": "Merge pull request #1030 from vmware-tanzu/fix-trivy-scan-1\n\nFix warning in trivy scan and add error check to the command",
          "timestamp": "2022-12-22T11:11:22-07:00",
          "tree_id": "5a68349715da9e91bae2922aa1791eeaedc8d66d",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f88670e53d679a1ae16c2eccb9abb4397af2d0bd"
        },
        "date": 1671733319147,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36672472054,
            "unit": "ns/op\t        18.89 DeleteSeconds\t        17.74 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8458544284,
            "unit": "ns/op\t         4.268 DeleteSeconds\t         4.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6737991429,
            "unit": "ns/op\t         4.544 DeleteSeconds\t         2.149 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8dde54286976fb485c35e0906a3b75115a0128cb",
          "message": "Merge pull request #1027 from vmware-tanzu/bump-dependencies\n\nBump dependencies",
          "timestamp": "2022-12-22T11:11:59-07:00",
          "tree_id": "608df58498f2ef96a16ef94c83eccd51f4c46fac",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/8dde54286976fb485c35e0906a3b75115a0128cb"
        },
        "date": 1671733328252,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36816167944,
            "unit": "ns/op\t        19.28 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8489462167,
            "unit": "ns/op\t         4.281 DeleteSeconds\t         4.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6423293925,
            "unit": "ns/op\t         4.243 DeleteSeconds\t         2.142 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f7e02c8595a5248252b9d36de8856fa489e01886",
          "message": "Merge pull request #1031 from vmware-tanzu/dependabot/github_actions/helm/kind-action-1.5.0\n\nBump helm/kind-action from 1.4.0 to 1.5.0",
          "timestamp": "2022-12-22T11:12:54-07:00",
          "tree_id": "24c523a2d9eabb86f3828ca93892f4134b96a145",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/f7e02c8595a5248252b9d36de8856fa489e01886"
        },
        "date": 1671733380041,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36443786127,
            "unit": "ns/op\t        18.91 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8484638080,
            "unit": "ns/op\t         4.271 DeleteSeconds\t         4.171 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6422362187,
            "unit": "ns/op\t         4.233 DeleteSeconds\t         2.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "782b007253ffdad06013c952aaea8c7598130deb",
          "message": "Merge pull request #1023 from vmware-tanzu/dependabot/github_actions/actions/checkout-3.2.0\n\nBump actions/checkout from 3.1.0 to 3.2.0",
          "timestamp": "2022-12-22T15:19:15-07:00",
          "tree_id": "e84a69a8a21bff6779eecab53d0ab5f74cf2a047",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/782b007253ffdad06013c952aaea8c7598130deb"
        },
        "date": 1671748248900,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35581308638,
            "unit": "ns/op\t        18.02 DeleteSeconds\t        17.51 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8511484395,
            "unit": "ns/op\t         4.292 DeleteSeconds\t         4.177 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6481370003,
            "unit": "ns/op\t         4.281 DeleteSeconds\t         2.155 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "joaod@vmware.com",
            "name": "João Pereira",
            "username": "joaopapereira"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "dcaf6e006a29472183258b9290cd665d2697cfc3",
          "message": "Merge pull request #1034 from vmware-tanzu/nh-fix-trivy-output\n\nFix trivy output formatting",
          "timestamp": "2023-01-02T10:06:05-06:00",
          "tree_id": "95d1c7a0eff9e1879353dd5e9b420223fb53b42e",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/dcaf6e006a29472183258b9290cd665d2697cfc3"
        },
        "date": 1672676390560,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37343293056,
            "unit": "ns/op\t        19.46 DeleteSeconds\t        17.81 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8739560669,
            "unit": "ns/op\t         4.450 DeleteSeconds\t         4.239 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6716709625,
            "unit": "ns/op\t         4.477 DeleteSeconds\t         2.185 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2dd4e3d6baf3f3469e6f85253865bced42252a3f",
          "message": "Merge pull request #1036 from vmware-tanzu/dependabot/github_actions/actions/stale-7.0.0\n\nBump actions/stale from 6.0.1 to 7.0.0",
          "timestamp": "2023-01-03T12:54:43-07:00",
          "tree_id": "b8ded8a2742b02e8e34d36f9db84d44c5561f6a3",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/2dd4e3d6baf3f3469e6f85253865bced42252a3f"
        },
        "date": 1672776420680,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36841447196,
            "unit": "ns/op\t        19.13 DeleteSeconds\t        17.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8714117869,
            "unit": "ns/op\t         4.387 DeleteSeconds\t         4.275 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6524168911,
            "unit": "ns/op\t         4.304 DeleteSeconds\t         2.173 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "60dc4c13796d92e81281503b6aa58e350b468449",
          "message": "Update hint in error messages (#1040)",
          "timestamp": "2023-01-05T23:36:00+05:30",
          "tree_id": "f3eabc909b12d5d8562e4e4fc33ee30df492ba9c",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/60dc4c13796d92e81281503b6aa58e350b468449"
        },
        "date": 1672942738522,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37091995231,
            "unit": "ns/op\t        19.29 DeleteSeconds\t        17.75 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8663337581,
            "unit": "ns/op\t         4.373 DeleteSeconds\t         4.236 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6586119265,
            "unit": "ns/op\t         4.341 DeleteSeconds\t         2.183 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "38600853+kumaritanushree@users.noreply.github.com",
            "name": "kumari tanushree",
            "username": "kumaritanushree"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "ad24bdc3e41219b176f88280a469e9fd21979339",
          "message": "Updated repo url and package name to generic name in examples (#1039)\n\n* updating repo url to generic name as TCE is getting deprecated and updated package name to a generic name as well\r\n\r\n* updated namespace flag to package available cmd\r\n\r\n* removed -n flag from cmd added in example\r\n\r\n* reverted unwanted change\r\n\r\nCo-authored-by: kumari tanushree <ktanushree@vmware.com>",
          "timestamp": "2023-01-06T10:06:44+05:30",
          "tree_id": "84906f760e6714b2705bb63e5dcc01e450b43ef7",
          "url": "https://github.com/vmware-tanzu/carvel-kapp-controller/commit/ad24bdc3e41219b176f88280a469e9fd21979339"
        },
        "date": 1672980419075,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36723494856,
            "unit": "ns/op\t        19.23 DeleteSeconds\t        17.44 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8522641312,
            "unit": "ns/op\t         4.326 DeleteSeconds\t         4.158 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6418738587,
            "unit": "ns/op\t         4.250 DeleteSeconds\t         2.132 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "63347c2522ec44373742a71e0553043c6a0a82a0",
          "message": "Merge pull request #1053 from carvel-dev/seccomp-profile\n\nDo not set seccompProfile",
          "timestamp": "2023-01-23T11:54:31-07:00",
          "tree_id": "411cd3b38323016a6f6da43d3212ed60ad7707df",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/63347c2522ec44373742a71e0553043c6a0a82a0"
        },
        "date": 1674500804997,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36667058934,
            "unit": "ns/op\t        18.94 DeleteSeconds\t        17.67 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8619664217,
            "unit": "ns/op\t         4.341 DeleteSeconds\t         4.224 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6584309329,
            "unit": "ns/op\t         4.336 DeleteSeconds\t         2.198 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "543679435db906c15af95de9c0bf3b2f1b7dccfa",
          "message": "Merge pull request #1062 from carvel-dev/nh-fix-add-to-issues\n\nFix add-to-issues to point to carvel-dev",
          "timestamp": "2023-01-23T13:52:51-07:00",
          "tree_id": "efa7d515e68d9baecc50885c2654b59ed43f5bd3",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/543679435db906c15af95de9c0bf3b2f1b7dccfa"
        },
        "date": 1674507826682,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36784162889,
            "unit": "ns/op\t        19.22 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8544781266,
            "unit": "ns/op\t         4.318 DeleteSeconds\t         4.181 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6495382518,
            "unit": "ns/op\t         4.297 DeleteSeconds\t         2.157 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9004b3990c0ce1342f0a783590b2b84ac0f64af5",
          "message": "Merge pull request #1045 from carvel-dev/dependabot/github_actions/actions/checkout-3.3.0\n\nBump actions/checkout from 3.2.0 to 3.3.0",
          "timestamp": "2023-01-23T13:53:25-07:00",
          "tree_id": "463b1fbbc518a3b065972a1a16595bf0575ebcd0",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/9004b3990c0ce1342f0a783590b2b84ac0f64af5"
        },
        "date": 1674508061037,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37433409815,
            "unit": "ns/op\t        19.67 DeleteSeconds\t        17.70 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8714071383,
            "unit": "ns/op\t         4.420 DeleteSeconds\t         4.244 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6639771771,
            "unit": "ns/op\t         4.404 DeleteSeconds\t         2.182 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "b68563c6749fd357968034f0af655a0066971ca2",
          "message": "Merge pull request #1037 from carvel-dev/bump-dependencies\n\nBump dependencies",
          "timestamp": "2023-01-23T16:04:15-07:00",
          "tree_id": "74a6113674d3a95cf31fbf07020121765c988340",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/b68563c6749fd357968034f0af655a0066971ca2"
        },
        "date": 1674515642985,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35590516942,
            "unit": "ns/op\t        18.08 DeleteSeconds\t        17.47 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8465912557,
            "unit": "ns/op\t         4.259 DeleteSeconds\t         4.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6410658674,
            "unit": "ns/op\t         4.238 DeleteSeconds\t         2.136 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "cecb740c5c999ae4b52fb54f318b00d303d40cf2",
          "message": "Merge pull request #1060 from carvel-dev/dependabot/go_modules/k8s.io/apiserver-0.25.6\n\nBump k8s.io/apiserver from 0.25.0 to 0.25.6",
          "timestamp": "2023-01-23T16:01:10-07:00",
          "tree_id": "3781867ce52a7445bc824c4587dc701cec207573",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/cecb740c5c999ae4b52fb54f318b00d303d40cf2"
        },
        "date": 1674515665611,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38584536134,
            "unit": "ns/op\t        19.74 DeleteSeconds\t        18.78 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8754334385,
            "unit": "ns/op\t         4.431 DeleteSeconds\t         4.257 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6635553935,
            "unit": "ns/op\t         4.390 DeleteSeconds\t         2.186 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "13ff0cbefa297138b3a15d46146adb3d7fd4003b",
          "message": "Add signed-off to each commit from the bot",
          "timestamp": "2023-01-23T16:10:25-07:00",
          "tree_id": "db1948e5b501baf33840a56e6cd012141f49bcea",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/13ff0cbefa297138b3a15d46146adb3d7fd4003b"
        },
        "date": 1674516025067,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36720065880,
            "unit": "ns/op\t        18.93 DeleteSeconds\t        17.74 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8553384743,
            "unit": "ns/op\t         4.352 DeleteSeconds\t         4.158 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6449558356,
            "unit": "ns/op\t         4.282 DeleteSeconds\t         2.130 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8dadce32112113631c769347b716f42b2499c4de",
          "message": "Add signed-off-by for carvel-bot",
          "timestamp": "2023-01-23T16:12:29-07:00",
          "tree_id": "e133a720e78b6b0ba1fdad8f756db7a80aa05c58",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/8dadce32112113631c769347b716f42b2499c4de"
        },
        "date": 1674516169661,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36562261391,
            "unit": "ns/op\t        18.98 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8633606419,
            "unit": "ns/op\t         4.301 DeleteSeconds\t         4.288 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6473900267,
            "unit": "ns/op\t         4.271 DeleteSeconds\t         2.155 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a69036a6467ae923813a7b8482eb51985c0be197",
          "message": "Fix typo in \"dependency updater\" workflow",
          "timestamp": "2023-01-23T16:14:50-07:00",
          "tree_id": "57919989a5791fc94634fa761c00389704f0f6fc",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/a69036a6467ae923813a7b8482eb51985c0be197"
        },
        "date": 1674516301390,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35433257790,
            "unit": "ns/op\t        17.85 DeleteSeconds\t        17.54 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8485949897,
            "unit": "ns/op\t         4.272 DeleteSeconds\t         4.174 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6525833109,
            "unit": "ns/op\t         4.309 DeleteSeconds\t         2.156 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "86852107+joe-kimmel-vmw@users.noreply.github.com",
            "name": "Joe Kimmel",
            "username": "joe-kimmel-vmw"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9afffabb6ea453ae050c95a6dc25533c22a746e4",
          "message": "dev.md - fix link to benchmark after repo donation (#1070)\n\nSigned-off-by: Joe Kimmel <jkimmel@vmware.com>\r\n\r\nSigned-off-by: Joe Kimmel <jkimmel@vmware.com>",
          "timestamp": "2023-01-26T11:48:42-05:00",
          "tree_id": "4fe17e4532c35405bb6acffe0f48c8c423af2a0f",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/9afffabb6ea453ae050c95a6dc25533c22a746e4"
        },
        "date": 1674752536985,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38301289991,
            "unit": "ns/op\t        19.42 DeleteSeconds\t        18.83 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8966633407,
            "unit": "ns/op\t         4.624 DeleteSeconds\t         4.269 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6663823121,
            "unit": "ns/op\t         4.395 DeleteSeconds\t         2.204 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "joaod@vmware.com",
            "name": "João Pereira",
            "username": "joaopapereira"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "50caee3ca668c2d9cb84739174b5c8110cfc4b59",
          "message": "Sidecar execution honor sidecar environment variables instead of main pod (#1068)\n\nSigned-off-by: João Pereira <joaod@vmware.com>\r\n\r\nSigned-off-by: João Pereira <joaod@vmware.com>",
          "timestamp": "2023-01-26T20:02:27-05:00",
          "tree_id": "013d03dcccd84d21f0ae84ddea2e86de78fef7bd",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/50caee3ca668c2d9cb84739174b5c8110cfc4b59"
        },
        "date": 1674781960028,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36636445078,
            "unit": "ns/op\t        18.96 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8518164789,
            "unit": "ns/op\t         4.303 DeleteSeconds\t         4.175 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6477002429,
            "unit": "ns/op\t         4.278 DeleteSeconds\t         2.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "26af4b0e4e0dd1b9cde987704233436737615a4c",
          "message": "Merge pull request #1073 from carvel-dev/go-bump\n\nBump go 1.19.5 to develop",
          "timestamp": "2023-02-03T15:07:07-07:00",
          "tree_id": "77056e448f1e2d6590bfb5c8183c1bb3520bfa48",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/26af4b0e4e0dd1b9cde987704233436737615a4c"
        },
        "date": 1675462643001,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36745749112,
            "unit": "ns/op\t        19.00 DeleteSeconds\t        17.70 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8525140038,
            "unit": "ns/op\t         4.300 DeleteSeconds\t         4.173 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6539578728,
            "unit": "ns/op\t         4.291 DeleteSeconds\t         2.197 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "fd680ce61b8af3781b236f69a4ac0f1a08a453e6",
          "message": "Merge pull request #1065 from carvel-dev/bump-dependencies\n\nBump dependencies",
          "timestamp": "2023-02-03T15:07:40-07:00",
          "tree_id": "0fcfde004a0a3b62d6adf23c93072b04eec90c03",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/fd680ce61b8af3781b236f69a4ac0f1a08a453e6"
        },
        "date": 1675462674476,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37541349621,
            "unit": "ns/op\t        19.02 DeleteSeconds\t        18.48 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8536635847,
            "unit": "ns/op\t         4.300 DeleteSeconds\t         4.191 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6538332406,
            "unit": "ns/op\t         4.272 DeleteSeconds\t         2.218 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "joaod@vmware.com",
            "name": "João Pereira",
            "username": "joaopapereira"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "2a6d8da190be2ac38820706ea889413ebb7beaf0",
          "message": "Change org and repository name in develop (#1101)\n\nSigned-off-by: João Pereira <joaod@vmware.com>\r\nCo-authored-by: Varsha Munishwar <vmunishwar@vmware.com>",
          "timestamp": "2023-02-21T17:28:30-05:00",
          "tree_id": "fa9066b2787823159379d5cb146b2a82d4b82edc",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/2a6d8da190be2ac38820706ea889413ebb7beaf0"
        },
        "date": 1677019108666,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35382216039,
            "unit": "ns/op\t        17.88 DeleteSeconds\t        17.46 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8482547683,
            "unit": "ns/op\t         4.290 DeleteSeconds\t         4.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6458949533,
            "unit": "ns/op\t         4.258 DeleteSeconds\t         2.159 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "49699333+dependabot[bot]@users.noreply.github.com",
            "name": "dependabot[bot]",
            "username": "dependabot[bot]"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a26d873f5c53ec36cafac48dd198a6c6ef00c52e",
          "message": "Bump golang.org/x/net from 0.4.0 to 0.7.0 in /cli (#1102)\n\nBumps [golang.org/x/net](https://github.com/golang/net) from 0.4.0 to 0.7.0.\r\n- [Release notes](https://github.com/golang/net/releases)\r\n- [Commits](https://github.com/golang/net/compare/v0.4.0...v0.7.0)\r\n\r\n---\r\nupdated-dependencies:\r\n- dependency-name: golang.org/x/net\r\n  dependency-type: indirect\r\n...\r\n\r\nSigned-off-by: dependabot[bot] <support@github.com>\r\nCo-authored-by: dependabot[bot] <49699333+dependabot[bot]@users.noreply.github.com>",
          "timestamp": "2023-02-22T10:50:58+05:30",
          "tree_id": "93d9dbc7d3a8696468ee586e479350f2cd8132c6",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/a26d873f5c53ec36cafac48dd198a6c6ef00c52e"
        },
        "date": 1677043866889,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36571814408,
            "unit": "ns/op\t        19.04 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8574677794,
            "unit": "ns/op\t         4.363 DeleteSeconds\t         4.167 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 5466984566,
            "unit": "ns/op\t         3.264 DeleteSeconds\t         2.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "f1c91e05e12510bec00f2ff2459439326d67592d",
          "message": "Merge pull request #1066 from carvel-dev/dependabot/go_modules/k8s.io/code-generator-0.25.6\n\nBump k8s.io/code-generator from 0.25.0 to 0.25.6",
          "timestamp": "2023-02-22T11:08:35-07:00",
          "tree_id": "2f1effc435db9e4841bb2703795be1927194d295",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/f1c91e05e12510bec00f2ff2459439326d67592d"
        },
        "date": 1677090038370,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36705939434,
            "unit": "ns/op\t        19.07 DeleteSeconds\t        17.58 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8578270853,
            "unit": "ns/op\t         4.335 DeleteSeconds\t         4.191 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6484060691,
            "unit": "ns/op\t         4.285 DeleteSeconds\t         2.147 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "33b7fabd319231b3dd85554b6e31e70c0cfc994c",
          "message": "allow minor bumps to depedencies via dependabot",
          "timestamp": "2023-02-22T11:13:06-07:00",
          "tree_id": "07d5a1992bfdc8c42f85eaa9fcb5eb63d62dffd4",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/33b7fabd319231b3dd85554b6e31e70c0cfc994c"
        },
        "date": 1677090188666,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35415438679,
            "unit": "ns/op\t        17.87 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8495718573,
            "unit": "ns/op\t         4.284 DeleteSeconds\t         4.168 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6432314697,
            "unit": "ns/op\t         4.257 DeleteSeconds\t         2.134 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "8d8de5c2fb080a23db604e1c4f4c0f29a4b13ca0",
          "message": "Merge pull request #1098 from carvel-dev/bump-dependencies\n\nBump dependencies",
          "timestamp": "2023-02-22T13:15:50-07:00",
          "tree_id": "7a7e47b48be32fd3256f22e99d4aa57befaa5162",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/8d8de5c2fb080a23db604e1c4f4c0f29a4b13ca0"
        },
        "date": 1677097577265,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35509973377,
            "unit": "ns/op\t        17.91 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8637626725,
            "unit": "ns/op\t         4.418 DeleteSeconds\t         4.174 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6472593650,
            "unit": "ns/op\t         4.275 DeleteSeconds\t         2.150 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "dcf92debf7567107b5af598cb930c949bd72ed72",
          "message": "Merge pull request #1107 from carvel-dev/dependabot/go_modules/golang.org/x/tools-0.6.0\n\nBump golang.org/x/tools from 0.1.12 to 0.6.0",
          "timestamp": "2023-02-22T13:15:24-07:00",
          "tree_id": "b3147a652e88b2847d373d53a631f675b051f757",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/dcf92debf7567107b5af598cb930c949bd72ed72"
        },
        "date": 1677097661789,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35755379542,
            "unit": "ns/op\t        18.07 DeleteSeconds\t        17.64 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8591711544,
            "unit": "ns/op\t         4.353 DeleteSeconds\t         4.193 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7556562155,
            "unit": "ns/op\t         4.331 DeleteSeconds\t         3.169 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "distinct": true,
          "id": "dc9017a0886b809b21fc6dccf70d8ff04f63c29b",
          "message": "Deflake TestConfig_TrustCACerts\n\nSigned-off-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>",
          "timestamp": "2023-03-01T10:57:33+05:30",
          "tree_id": "60b0540fb463ac7bc76c0e787bd9001b45e17188",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/dc9017a0886b809b21fc6dccf70d8ff04f63c29b"
        },
        "date": 1677649102291,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36565988952,
            "unit": "ns/op\t        18.99 DeleteSeconds\t        17.53 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8626260033,
            "unit": "ns/op\t         4.405 DeleteSeconds\t         4.166 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6595135179,
            "unit": "ns/op\t         4.385 DeleteSeconds\t         2.162 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9682ece7c8379252065a8a37afce7af57d5de27d",
          "message": "Do not expose development values to kc package bundle (#1111)\n\nStructure the config values to have a clear separation of values\r\nUpdate dev-deploy.sh to use these config values\r\n\r\nSigned-off-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>",
          "timestamp": "2023-03-02T19:48:02+05:30",
          "tree_id": "41806ff8634a1bbbb146258f7810b63bd8f11d83",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/9682ece7c8379252065a8a37afce7af57d5de27d"
        },
        "date": 1677767454256,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 38130987265,
            "unit": "ns/op\t        19.32 DeleteSeconds\t        18.74 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8734704011,
            "unit": "ns/op\t         4.409 DeleteSeconds\t         4.225 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 7690574322,
            "unit": "ns/op\t         4.371 DeleteSeconds\t         3.266 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "18197babab94ab91012b6ffe8507ea71aadeaff3",
          "message": "Merge pull request #1116 from carvel-dev/dependabot/github_actions/actions/add-to-project-0.4.1\n\nBump actions/add-to-project from 0.4.0 to 0.4.1",
          "timestamp": "2023-03-02T13:23:55-07:00",
          "tree_id": "ff3a04946288766382aecc0035613503dec7ce13",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/18197babab94ab91012b6ffe8507ea71aadeaff3"
        },
        "date": 1677789279067,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36573840764,
            "unit": "ns/op\t        19.01 DeleteSeconds\t        17.52 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8537694307,
            "unit": "ns/op\t         4.329 DeleteSeconds\t         4.167 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 5478190182,
            "unit": "ns/op\t         3.278 DeleteSeconds\t         2.159 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9870aee8b5b80cbd0d6f218461171aa423c6adb0",
          "message": "Merge pull request #1106 from carvel-dev/dependabot/go_modules/github.com/prometheus/client_golang-1.14.0\n\nBump github.com/prometheus/client_golang from 1.12.2 to 1.14.0",
          "timestamp": "2023-03-02T13:24:29-07:00",
          "tree_id": "54175d212d6a5c64f99fd7180eddca0d805dae0d",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/9870aee8b5b80cbd0d6f218461171aa423c6adb0"
        },
        "date": 1677789359936,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36989562009,
            "unit": "ns/op\t        19.31 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8583307794,
            "unit": "ns/op\t         4.334 DeleteSeconds\t         4.201 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6521104595,
            "unit": "ns/op\t         4.302 DeleteSeconds\t         2.170 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "distinct": true,
          "id": "4485e29a2b218b333d9a99e0547606351356f7d5",
          "message": "Add example for pkg repo kick\n\nSigned-off-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>",
          "timestamp": "2023-03-06T11:18:16+05:30",
          "tree_id": "e358684a2ee2748a08318552b5b7ff8bbb46a393",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/4485e29a2b218b333d9a99e0547606351356f7d5"
        },
        "date": 1678082322149,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37722488906,
            "unit": "ns/op\t        19.11 DeleteSeconds\t        18.57 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8527624946,
            "unit": "ns/op\t         4.306 DeleteSeconds\t         4.179 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6470364605,
            "unit": "ns/op\t         4.276 DeleteSeconds\t         2.141 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "soumikm@vmware.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "distinct": true,
          "id": "ab03cc28b25b8d5e50d1196250f431f7690e5318",
          "message": "Add tests for package repo dry-run. Fix logic for creating RBAC resources.\n\nSigned-off-by: Soumik Majumder <soumikm@vmware.com>",
          "timestamp": "2023-03-07T17:34:35+05:30",
          "tree_id": "d33e78505866f3041ec67143b12f79b196856a82",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/ab03cc28b25b8d5e50d1196250f431f7690e5318"
        },
        "date": 1678191372248,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37134457519,
            "unit": "ns/op\t        18.21 DeleteSeconds\t        18.88 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8625474699,
            "unit": "ns/op\t         4.370 DeleteSeconds\t         4.202 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6751103722,
            "unit": "ns/op\t         4.511 DeleteSeconds\t         2.177 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "33070011+100mik@users.noreply.github.com",
            "name": "Soumik Majumder",
            "username": "100mik"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "3af991eafb55d94ab590cc4190b3343526205fa5",
          "message": "Allow disabling ytt validations while building packages (#1077)\n\n* Disable ytt validations while building packages\r\n\r\nSigned-off-by: Soumik Majumder <soumikm@vmware.com>\r\n\r\n* Add test to ensure that kctrl disables ytt validations while releasing packages\r\n\r\nSigned-off-by: Soumik Majumder <soumikm@vmware.com>\r\n\r\n* Add flag to disable ytt validations while releasing package\r\n\r\nSigned-off-by: Soumik Majumder <soumikm@vmware.com>\r\n\r\n* Add test case for using ytt validations while building packages\r\n\r\nSigned-off-by: Soumik Majumder <soumikm@vmware.com>\r\n\r\n---------\r\n\r\nSigned-off-by: Soumik Majumder <soumikm@vmware.com>",
          "timestamp": "2023-03-07T18:10:47+05:30",
          "tree_id": "296fbfa6ff76cdc855b10afeda1aa98e1bbe1650",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/3af991eafb55d94ab590cc4190b3343526205fa5"
        },
        "date": 1678193685025,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37492231456,
            "unit": "ns/op\t        18.56 DeleteSeconds\t        18.88 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8847715912,
            "unit": "ns/op\t         4.510 DeleteSeconds\t         4.273 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6870659040,
            "unit": "ns/op\t         4.528 DeleteSeconds\t         2.228 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "943bfa9cad4da152f68883db49669bb933b4702d",
          "message": "Merge pull request #1123 from carvel-dev/dependabot/go_modules/github.com/vmware-tanzu/carvel-vendir-0.33.1\n\nBump github.com/vmware-tanzu/carvel-vendir from 0.30.1 to 0.33.1",
          "timestamp": "2023-03-07T09:58:16-07:00",
          "tree_id": "c5c4b4be9df1e396524f014d93402344e26e6408",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/943bfa9cad4da152f68883db49669bb933b4702d"
        },
        "date": 1678208927321,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35521433054,
            "unit": "ns/op\t        17.97 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8508095997,
            "unit": "ns/op\t         4.291 DeleteSeconds\t         4.171 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6475200319,
            "unit": "ns/op\t         4.277 DeleteSeconds\t         2.149 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "039ce70667c6879669ad1aa09d0f706be33a0994",
          "message": "Print errors while parsing default values for a pkg (#1041)\n\nDo not ignore errors are from `saveDefaultValuesFileOutput`.",
          "timestamp": "2023-03-07T23:49:34+05:30",
          "tree_id": "8d304c7773ea18453b96b3a936bb433194709f16",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/039ce70667c6879669ad1aa09d0f706be33a0994"
        },
        "date": 1678213791885,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36571336051,
            "unit": "ns/op\t        19.04 DeleteSeconds\t        17.49 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8512787910,
            "unit": "ns/op\t         4.295 DeleteSeconds\t         4.169 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6610432251,
            "unit": "ns/op\t         4.378 DeleteSeconds\t         2.161 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "7980c08c7ef990d7d73be357925ea162f1463f5d",
          "message": "Bump kapp to v0.55.0 (#1125)\n\n* Bump kapp to v0.55.0\r\n\r\nSigned-off-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>\r\n\r\n* Update kapp error message in tests\r\n\r\nWith kapp v0.55.0 we display the usefulErrorMessage as part of the error.\r\n\r\nSigned-off-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>\r\n\r\n---------\r\n\r\nSigned-off-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>",
          "timestamp": "2023-03-09T01:20:14+05:30",
          "tree_id": "d4909d6090e4ad56b6822fbcfa5494e9517a06ba",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/7980c08c7ef990d7d73be357925ea162f1463f5d"
        },
        "date": 1678305616193,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36422961800,
            "unit": "ns/op\t        18.93 DeleteSeconds\t        17.45 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8513360789,
            "unit": "ns/op\t         4.314 DeleteSeconds\t         4.153 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6446699872,
            "unit": "ns/op\t         4.259 DeleteSeconds\t         2.142 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "8457124+praveenrewar@users.noreply.github.com",
            "name": "Praveen Rewar",
            "username": "praveenrewar"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "a22e85996b0cd11dd176c536a366f23e91f95bdc",
          "message": "Add check for annotations field while looking for pkg_repo_ann (#1127)\n\nSigned-off-by: Praveen Rewar <8457124+praveenrewar@users.noreply.github.com>",
          "timestamp": "2023-03-09T10:44:40+05:30",
          "tree_id": "9ac0e08b381fdaf36abf79c08b81705ffc759a4e",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/a22e85996b0cd11dd176c536a366f23e91f95bdc"
        },
        "date": 1678339539501,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36914069894,
            "unit": "ns/op\t        19.30 DeleteSeconds\t        17.56 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8721269207,
            "unit": "ns/op\t         4.450 DeleteSeconds\t         4.217 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6522481855,
            "unit": "ns/op\t         4.298 DeleteSeconds\t         2.171 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "922e32fbd868552f0b5a3a3dae0b7c5c6edbfb2a",
          "message": "Merge pull request #1133 from carvel-dev/upgrade-go-and-dependencies\n\nUpdating go version and dependencies",
          "timestamp": "2023-03-09T11:23:57-07:00",
          "tree_id": "e93ae8ad112f40d0bde6c810892855c67e1966ed",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/922e32fbd868552f0b5a3a3dae0b7c5c6edbfb2a"
        },
        "date": 1678386997897,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36037653826,
            "unit": "ns/op\t        18.31 DeleteSeconds\t        17.66 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8693157992,
            "unit": "ns/op\t         4.411 DeleteSeconds\t         4.221 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6653617799,
            "unit": "ns/op\t         4.405 DeleteSeconds\t         2.199 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "75b5f47f302eb95e96966ed2daccc65c46747126",
          "message": "Merge pull request #1122 from carvel-dev/dependabot/github_actions/benchmark-action/github-action-benchmark-1.16.1\n\nBump benchmark-action/github-action-benchmark from 1.15.0 to 1.16.1",
          "timestamp": "2023-03-09T14:12:53-07:00",
          "tree_id": "c6094904f676263d9364aff9b03084e6bf7af806",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/75b5f47f302eb95e96966ed2daccc65c46747126"
        },
        "date": 1678397081418,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35894887975,
            "unit": "ns/op\t        18.29 DeleteSeconds\t        17.55 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8594912032,
            "unit": "ns/op\t         4.333 DeleteSeconds\t         4.205 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6549056752,
            "unit": "ns/op\t         4.308 DeleteSeconds\t         2.188 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "55523204+rohitagg2020@users.noreply.github.com",
            "name": "rohitagg2020",
            "username": "rohitagg2020"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "53f77542f1afbb113c9ad7458d40b1e416ee9b5b",
          "message": "Merge pull request #1120 from praveenrewar/update-package-values\n\nExpose values in kapp-controller package",
          "timestamp": "2023-03-10T15:43:32+05:30",
          "tree_id": "bb264b5b871d2c38340e370a008adac95a4d751c",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/53f77542f1afbb113c9ad7458d40b1e416ee9b5b"
        },
        "date": 1678443865427,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36874510976,
            "unit": "ns/op\t        19.19 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8606800607,
            "unit": "ns/op\t         4.346 DeleteSeconds\t         4.213 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 5578976736,
            "unit": "ns/op\t         3.359 DeleteSeconds\t         2.172 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "nhickey@vmware.com",
            "name": "Neil Hickey",
            "username": "neil-hickey"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "edac0f1ec76d53b09368f7b900aa841a907b8d1d",
          "message": "Merge pull request #1130 from carvel-dev/dependabot/go_modules/k8s.io/klog/v2-2.90.1\n\nBump k8s.io/klog/v2 from 2.70.1 to 2.90.1",
          "timestamp": "2023-03-10T12:07:39-07:00",
          "tree_id": "f09be6255c31210e23a458e839c970e215dce67f",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/edac0f1ec76d53b09368f7b900aa841a907b8d1d"
        },
        "date": 1678476070379,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 37398325920,
            "unit": "ns/op\t        19.59 DeleteSeconds\t        17.75 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8770318305,
            "unit": "ns/op\t         4.494 DeleteSeconds\t         4.214 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6759388647,
            "unit": "ns/op\t         4.441 DeleteSeconds\t         2.215 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "ThomasVitale@users.noreply.github.com",
            "name": "Thomas Vitale",
            "username": "ThomasVitale"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "15c92dedd8e726d8f8281e50be42f497e4339584",
          "message": "`kctrl`: Flag to create namespace when adding new repo (#1113)\n\n* kctrl: Flag to create namespace when adding repo\r\n\r\nWhen adding a new package repository to a cluster, it's now possible\r\nto create the installation namespace automatically\r\nby specifying the \"--create-namespace\" flag.\r\n\r\nFixes gh-1001\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>\r\n\r\n* Improve error handling\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>\r\n\r\n* Optimize status messages for namespace creation\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>\r\n\r\n* Add cleanup after tests\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>\r\n\r\n* Update test cleanup\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>\r\n\r\n* Use cleanup function for new tests\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>\r\n\r\n* Re-use existing namespace in test\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>\r\n\r\n---------\r\n\r\nSigned-off-by: Thomas Vitale <ThomasVitale@users.noreply.github.com>",
          "timestamp": "2023-04-05T21:16:48+05:30",
          "tree_id": "48b4ae83f8069b1678cf15c125b1e645037ad738",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/15c92dedd8e726d8f8281e50be42f497e4339584"
        },
        "date": 1680710434686,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36340507489,
            "unit": "ns/op\t        18.59 DeleteSeconds\t        17.69 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8796131110,
            "unit": "ns/op\t         4.479 DeleteSeconds\t         4.242 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6675035241,
            "unit": "ns/op\t         4.392 DeleteSeconds\t         2.213 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "clebs@users.noreply.github.com",
            "name": "Borja Clemente",
            "username": "clebs"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "9a0282f460933a0240faf4f8fd1f69711861ccbf",
          "message": "Fix panic calling `tanzu package installed status` (#1161)\n\n* Fix panic calling tanzu package installed status\r\n\r\nCalling tanzu package installed status without any arguments causes the\r\nprogram to panic.\r\n\r\nSigned-off-by: Borja Clemente <cborja@vmware.com>\r\n\r\n* Apply review feedback\r\n\r\nSigned-off-by: Borja Clemente <cborja@vmware.com>\r\n\r\n---------\r\n\r\nSigned-off-by: Borja Clemente <cborja@vmware.com>",
          "timestamp": "2023-04-05T23:21:40+05:30",
          "tree_id": "6c6cff30a8279f6a8ca45100d47a7f871345185e",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/9a0282f460933a0240faf4f8fd1f69711861ccbf"
        },
        "date": 1680717788721,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 36837939494,
            "unit": "ns/op\t        19.16 DeleteSeconds\t        17.63 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8644547376,
            "unit": "ns/op\t         4.343 DeleteSeconds\t         4.248 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6522204197,
            "unit": "ns/op\t         4.306 DeleteSeconds\t         2.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      },
      {
        "commit": {
          "author": {
            "email": "clebs@users.noreply.github.com",
            "name": "Borja Clemente",
            "username": "clebs"
          },
          "committer": {
            "email": "noreply@github.com",
            "name": "GitHub",
            "username": "web-flow"
          },
          "distinct": true,
          "id": "5b29ad419e759d362a3bb35ef6a89858a953b960",
          "message": "Fix panic reading empty args slice (#1163)\n\nSigned-off-by: Borja Clemente <cborja@vmware.com>",
          "timestamp": "2023-04-06T14:11:49+05:30",
          "tree_id": "aedab09f86fb40131c016a11beed3ddb3b0e05cd",
          "url": "https://github.com/carvel-dev/kapp-controller/commit/5b29ad419e759d362a3bb35ef6a89858a953b960"
        },
        "date": 1680771120247,
        "tool": "go",
        "benches": [
          {
            "name": "Benchmark_pkgr_with_500_packages",
            "value": 35553712855,
            "unit": "ns/op\t        17.97 DeleteSeconds\t        17.50 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_100_packages",
            "value": 8649724662,
            "unit": "ns/op\t         4.304 DeleteSeconds\t         4.297 DeploySeconds",
            "extra": "1 times\n2 procs"
          },
          {
            "name": "Benchmark_pkgr_with_50_packages",
            "value": 6482522999,
            "unit": "ns/op\t         4.269 DeleteSeconds\t         2.163 DeploySeconds",
            "extra": "1 times\n2 procs"
          }
        ]
      }
    ]
  }
}