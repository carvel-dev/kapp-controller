window.BENCHMARK_DATA = {
  "lastUpdate": 1662575677578,
  "repoUrl": "https://github.com/vmware-tanzu/carvel-kapp-controller",
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
      }
    ]
  }
}