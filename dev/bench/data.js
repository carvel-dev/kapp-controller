window.BENCHMARK_DATA = {
  "lastUpdate": 1655493052647,
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
      }
    ]
  }
}