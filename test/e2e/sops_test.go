 // Copyright 2020 VMware, Inc.
 // SPDX-License-Identifier: Apache-2.0

package e2e

import (
    "strings"
    "testing"

    "github.com/ghodss/yaml"
    corev1 "k8s.io/api/core/v1"
)

func TestSops(t *testing.T) {
    env := BuildEnv(t)
    logger := Logger{}
    kapp := Kapp{t, env.Namespace, logger}
    sas := ServiceAccounts{env.Namespace}

    yaml1 := `
apiVersion: kappctrl.k14s.io/v1alpha1
kind: App
metadata:
  name: test-sops
  annotations:
    kapp.k14s.io/change-group: kappctrl-e2e.k14s.io/apps
spec:
  serviceAccountName: kappctrl-e2e-ns-sa
  fetch:
  - inline:
      paths:
        config.yml: |
          #@ load("@ytt:data", "data")
          apiVersion: v1
          kind: ConfigMap
          metadata:
            name: configmap
          data:
            key: #@ data.values.value
        values.wsops.yml: |
          data: ENC[AES256_GCM,data:/CF3wqiBqYVhZEWTtiBELgorAZDNJIj/j2Eh4xo/NRInXBv6w6gZ/zXM,iv:roiMIOpvtMZmUJWPrPqMbTQFbUG9tfRDb/YPm6BjDzU=,tag:4nVM1sORfSzk8bXRj5AysQ==,type:str]
          sops:
              kms: []
              gcp_kms: []
              azure_kv: []
              hc_vault: []
              lastmodified: '2020-10-03T01:05:21Z'
              mac: ENC[AES256_GCM,data:fVzy1fhbbRRni4Dtlgh/UchZs+dEMmG9kdI82xcZonTxCBjnXZqSIomSb3rtDxKXMFkgn7UyfrrQzD6duUpF5qiX4o0BTC412mN1LEbsh5LO3Zl4IAgbSxuGdyIljgk9Pzv5dBz7ajoQE8Z5A78e5N3+5DBYf8XxVhgUvBkeRT0=,iv:bZqkgGLI9CLNPTQwU4RBVBCz8IUpzmmYLF1WU32dwEs=,tag:qUpN3Fh5QC5pX3gHULEJng==,type:str]
              pgp:
              -   created_at: '2020-10-03T01:05:21Z'
                  enc: |
                      -----BEGIN PGP MESSAGE-----
                      Version: GnuPG v1

                      hQIMA/7je44gmO38ARAAijjlBfHB/swjSnTJijzddFAGDZKPo7c/SBPYt5ql3eBg
                      dp3SWFihynQtyDbIu/gTHyjkHu/X9Zp1E01UwavgnGrKIl8dYn9OsZXQi1t9zcZF
                      Cm4UQNXSaKV0NATlC6Bz35s2WZDQyRQQOHvItCnPprz68/Om84LanN5+Aj7GHcgg
                      9DyGbt+s0slGfQxi3Sj/3YhcxgFxblJ95yygRkbopfJv9+cVNdTCCAn4gO/TbtMl
                      SDAukhMycbbXf1pxKGnUn07Y3kY39++209B1AQ+V/R+LkjFxP+XOIZSiY0oBxror
                      hXF5lwqMVOFpJeni61iS9qqJ2WzytxAoOaiXz61bHZl0HmWtSZ6/htwDzILPHujG
                      VzO4YpytTZqrhsIVtrVzoQQgcuVmWhHnXFAf3wUx+4AAe5jetfCwgioKJ69KVgzB
                      NwSQIX7OrgsJ+0ypJLMJWBU4J5cfWN5vETIDgKipxN2gOZPdUu3CqBN8qDTB8Vdd
                      oHOAHhEd09JuaJ6CPru3HAe8aUV46QZriK/1C1ogxpxF9OL0WynVUtn7i+bJumtz
                      QfERFMZuaj4aP3N9r49vmvWRIaozEQbatXVKhaibIcJH8pCNbJ8I6ZOWjsOpMomA
                      VfMWJXq3eWwDIjIi7rE7kAzeH+xQp31M3jX53R5huCY1u19MDjjZqn58K8k2Ed7S
                      XAFJPGiVi1g3Pl2H8aNj+/J5HjmDlOCwueCNZTi1rXnywYfrATVoAhOaH49LSdU6
                      IORRfu9qKdVPVHzQj1AEJBEGZXKsx8Xw1q1m0LPQEDuIFtDaVP42jH8HsQK+
                      =do5e
                      -----END PGP MESSAGE-----
                  fp: B464DFD255C6B9F8
              unencrypted_suffix: _unencrypted
              version: 3.6.1
  template:
  - sops:
      pgp:
        privateKeysSecretRef:
          name: pgp-key
  - ytt: {}
  deploy:
  - kapp: {}
---
apiVersion: v1
kind: Secret
metadata:
  name: pgp-key
stringData:
  all.pk: |
    -----BEGIN PGP PRIVATE KEY BLOCK-----
    Version: GnuPG v1

    lQcYBF93yvoBEAC3vDMJY02+1q0liRo8RFDo2T2DAVhm2nHeBTU+CUDpVHN/K0vT
    rrZoy90hLuhI7N6LrCmT9hIrn+bYdaZnakoAIUw+SUfFZsVOfnwTSiya4MoHxb2R
    EPJcR3gASzXNBFsD74w3IpWRrEVqdaK+hk7j11poy9RbOCD5kxId6u2BEDrNQGek
    6Xysp3NUF6+WfpPKoyf5OCddgSOtrwQQr6QW/7FCCcgrijplFx5YBL3OIosWP2mg
    Cx4iijJchYUSmflQNnVa6hhjbecZUs0fMgBb1ceXNn7cjpCQlbyw8OnuCvNvq37I
    PHELO1yn97QbjfQldg4bwZk8f5C5DtaedVT6Vg9PSqFVxuqOka76m1dVZqsqRbvo
    04Uw8elxCfWah16rlbGRMHEBw03ESiLck0Wi6oAAf3a51nuX/ZdNxgCBgfH2QaIb
    CidKbZHjdlDRTx7SRF7VF/w5DGWEGqggQioCt2jUCHoMzO8pbeGxJgV/MhxhX1J3
    Lfc/K32tmWq8W/Z7Tr6qKNC8Tv0aIu11cSnt7f0jyNPCEV5gMrRcAvD4cyeixKSR
    +r3OV0IJbKxtiVifh7pemgrzVrDH5ruPXwdAhJUMut0nRXCVS+WUDuaviESc1yck
    zPVvbslQUk48ESggQWlskGKNJVm62gjLzSmAb9f5QJ4mma+ay6rj5xs12QARAQAB
    AA/+JT8QI4+Pc6fmSs0r1drNgh6D4zpTFuqimaz5mZ1bnNFjZny67uskjEMDjVYK
    fboS9UKN3TJNha1xKSFUffNkk/ksERZe58wJJHvsoCZxu2XlXsT9xFoon39XesvE
    WM7QupAFnymyI2lGWyoEt1XXyUVfTQ5A+sr6mE1xp0H4KqlFGlW7jQlOHlwFu57f
    mAUJ5dLEaDezdzeWKX/otY89lvH3l5kPDJCFfPe+TX6MkaycAIMTYP/P+JWGVw40
    J0yyZ0na6Xa6QfHGHvKTpYbH3tYME4HXHtQBx2Wrbj0wVvK+Xb1owqPKEMpUYrKK
    vuk8fJmdi4/oBgUgK/uvk0ja1FrL6+tzYqhvZA07CsfZG7c36kPMwweM1zNNsa47
    pF0QFualYm9Xt1GtQ9cg2snluha1rgkR2K4XZteSKdsrlp40RgmWDTKTt6+ReQXX
    /3WWkHvRJBX00Sth41z6zalnCAzSLr8TrYWqqvl0RRX0ke+q9rUUbnwgWbKWYM9N
    4kD4Xb5qYMm5oTIhAnp/oe3l2j3Kh0kXTr6qwHgbNXcV7VoKYRcb8trtqACaxBuO
    LjIIfgL66wswjDpTqjeEU8ptQu8ozDEes62ca3AydBxL1spZJVfXAUDB6YYiNmNu
    QLWdEU5YsgZg//VB5ORa2uoVq+qbfj/XeYsVMw1FMlYCXN8IANWufcbWdnjo+Eo/
    UM2BeYpyMzID11mSLGEm2VaFLGHMnpb3VG9qDgU8jxWzgv5j7L0lWIZMzG+Wk/yo
    X8lDVNt4O2JQhzc2LMT1mY/h0kTS9EoJi4RxQRoFOz+HPdDmscEkF8eBNru0oIti
    ++iA9wNlL9JP4tGrAXl4Yb97eG0O/DWrauW2KIiI60gvyA7hdeJ87oHXsLxUPCby
    NlElcSFTpL4vf/VcJBd8q1RIVvRZMvmSdCHZL0Wi0wUDji+bQppEHAEiqKzSBTQV
    V3nVNHbn4/xZoeXVe2l686NtDSFhoNLcQ7acUY40lE3xhn4qMv9YilN2ygNwudxm
    HP4R52cIANwfcZh2XYzDyotRy8IkI3BSkINQzuymuf29HF7LFvMzsMcwOHtFHu+v
    AJyej5aMVJSLdcqrbitIsjTr24zG2gJN/AEzoytcIau1s/ER2mOAmH7kiTX0XkWK
    YgYkwUHMo3wga+CUZNrCw8TTTu0w+61TSRjLmAVDKLPvvunT80QNYB4NzLb9bLFD
    cwdJ/KCrzHDhUDuEiOq45kgOxQGIYvMlEIsHNUkweOfNXqrK0RVzp5FJ26ZD3u4h
    nDiQWjuM/i/EbE9Ew5fHf8oBu605HZvjoZZmej71v5aq+BOFAbvb6Ebd2f2HGcIt
    CdTYq3cL54xuHPVs72DE5JFE4YRZ8L8H/iph3juXZELYoRgvUkAF9gQrU8LkRwfl
    YdNMmZC0YBXcchnOpbZ5FJhq+o4LcPj284h9vIK0RB4Uvj3tFedw5EVbEq59GijW
    /fXQdMAF54nxbeoKPIsWr997jzfZdamVOg5Jq59CSJvmgyDLET3Bz/DGMhwrAW5x
    nNvU2VTlPCRniD1LPKZrhcmhpQjRZCM2C9Mc2PTfsTS3PVKl8geSZEDap2zbS788
    HqKHVSphRqlpjTaEsZGxyDFX6Wa2S5ztQ+1pP6zlVqo0ca/62gejO0Vc2nFOuZnJ
    usSd13ay8+uW92/RuYsj+wTcn8yzbPZ0i3roal984OBWi4e1nOBoB8KIp7QgdGVz
    dCB0ZXN0ICh0ZXN0KSA8dGVzdEB0ZXN0LmNvbT6JAjgEEwECACIFAl93yvoCGwMG
    CwkIBwMCBhUIAgkKCwQWAgMBAh4BAheAAAoJELRk39JVxrn4jRcP/jyWz3wepwuc
    sidx1Vg0nuaMo+rIgVde8uiW9CI1vcpw3YDdRrYdZTi2d1IJGVmSmTkTQdMXjrDL
    hW7GcJWr5fVfEIWuxOxMkCaUU+rVX/7jmN9O6QLEA6OYsYDnlofroHJAyrg8rLOt
    C/L4jKYxCeowCy1Aj/6pbYBBE7R6lnyJFSNo4JBoqrF8an5bFX4/V9Qay2kSbm1q
    G+oYFEVWIrrP3b/Ia1lE/q80iqOEejtRgGhGs2rFfBJvZ/L8fk4/08OXzBmUPmcj
    +JN3+S9RMhLOvEvIkqejEmhdIiyRc1K7rUehkAJyfeInCbyrWdaCzNuTAsJ2Wb/f
    h1jkda+2C8+paJh7QXn59fHLxgqywCaoD/yGXvKDGsQd57XJgw46OJhjklg2PCoS
    8uYvqkAdefp2DLpWLV0U92/Z7Wj35tKqJzVKpN4J9NwR8GEHxNnDPI7Yl+RJElc7
    7b31mkzNyaJ9ngRFgLA427KLcW0R354hSks7A520EoVCS4NjVVH2pyteHN/tQa5J
    NqYSvyypxTzp0iqOMcqz+EiYTHEd4o+T//t556GZ7EWDpPGBIxCi95OqMLp5K/3U
    6IpmV/gPKPRUmZuaqgaVJ9Rr6ao87CpxVrcIcRt+gWK6+1Nw7oSTDpTKL0nuzYnW
    m+tyxcn7Yhdzc9VVTM3ga5LAfccoylhpnQcYBF93yvoBEADlwIMLTpsWftGyc0t0
    c/ikePoxKqaGARMiUCK3rB6bTrT8nRhD1guwevLjux/SurwfR8pMUj4943WY2K9y
    zTuq6XC0LnQQU/oTWm7wpCPtDqGrVwDNDRPQGgD8/vtxRdcurDDuFMaz9S+5paO3
    1KEVV6f7R6JH/Hr2hSMEa9XZ8tc+Hwh8nyJB7tujEZSuGPr6efbg67uTtG4qyvqx
    fbKmLXBRiTDp57d7XU2RfwqHHWqp4dvkwZIVpi1jQluZUHMH9ELChzpFajv3kM2T
    /aSTaV74YGlKnkAq17RqlRfHb07Di/r0EAq8xfj5UIDetTY3KyfX6Th+gb16odZG
    fTBP3vB7hdM/od32D8m4TVg4RdIvY/RanLJfSJMabAMyIjkmOK/UaqE8mNwDKnDH
    1PbbqF/DxG4hflaezXgEbchbgbxd1y5VZvv5gNnWdhdons6kx4kjpm+hxOP2qc3T
    tRlmdMEhfKqSqxXhwXzi6mdmAp555iuHR+TaSpE+9wGGkXpFSo7+M+npsbgE8iWo
    vFZEIEdwVk6YRPniMcDc0lNee1jUAp5e3dD/jrK7xZOjUMH7Lq75afRYOo3etH1E
    HKdTS9E0+Rw5rXYSLDZN8mXJwGzC+jtAy9b0WvQekuvRnvLq4PoFFcvqSZ7A1w/E
    xdNkOkOPF0R/0eylpgxRclo4AwARAQABAA/9HZ7lm5OVfcpOkXSOidkEeXqkvFSv
    W3TQBAso2V1LejfPhbILQCkHIMhOgFB5GIYSsvBigyx8nHtdh0iIdi3s6hVmqRRE
    H8bJOwLcbWdRam1fJ+PzqnwWeIH3FqcZFnrWn0x758tYmDhINZBxZxtMGUr6e/Rk
    VmZAGYBYtL03uP23WmmjLSNRxgZcMs3qyUyshEEtNG+v+Km5JQiTrEZ0aJBQfB1K
    knLZaQJCeeZTHnsLFr6Txw5dyITAMp003132i/6l4hvlG1tIQpwHT8knwAmZwOlN
    Kd2f5ZDMezagS4oXhtvejZDJPEVEhYAnh+RSX70knzrmRF5zEK3EyRm4HJV75Vm+
    leCWk1SbEq5j2P9X5vNRxYWbIueyDYzTzpsWhu7PksTEkF1g/bBVgA0x2eYCMhch
    VJAjpBQhC9XBSEJbZz4weGsRIRKJGAahyVjGjVzdNGkhQ0gOYMXmgFKauBqJND9Q
    HIxB2YeXjI2UXavcRstnnNW8UgZRfHSkJet8ztTNQS6OjKk213saEtPa4ybYAoO9
    mRxLncHZjyjy+gxTS7ZAopC20sjCWyFeWlkKzZIKDurFMgIOF4X5hFNGSZ+uj1zN
    DZN1aPA2IbdXQt78dvsBWWwT8soN6dEayJxfsdhamE9q7L3DKho8vM/AigkziP1I
    h1+ZKxLj5agmDsEIAOvsImx9zgjSb4RY1+0WoYjogFX1A+KZp/1EzvEKPDZdSDFf
    Kcy6zFL1QlSo8fS1TmqecCf4h27IIcg3PjcFja+slT+24BKkzBWU9dmnm9SIv2Hp
    BNVxV1OzfvGRvAMnYVtxiTZErjvRMg2twViM86WovCk6w/N7MNLCXjNY0NL7PdWc
    CzeAXYkYieIpBbaJwxhXfYLHK/bicg9o9I5pVNlm3nuy+6BdJXr5oPAXH9tMML+R
    QXylkmzcduXnrh306xLtU3W9S8DT+LJ5fI+GsdQx2Z2+Jhluo9FCc/D/4H0Y06Kq
    kwIwykHeheLfou86MAFZyBJVCv1GwNiFc1Z/jQsIAPlN8o91Xo4fuhOWyoRV4xLa
    D6tIt0F9VMJefzkp+ux4jIvdXTRcfJNsf253y8aWF1MfaOalpNfC4KXrvP/JSZHM
    mBtK+3a2Y/dNceOvBeAM6TZITDFpiLxcCBetDtCD2q8GgwzxEQXB/HRkehZUbhL5
    Ioo1eOBbC5hDI6OGS1kD8q8eemeYbrdIPbhl5SK2Rmblr9BqOuwQlsprGyF4l0m3
    YZ/qRashalmJPKe2FroFPpebFXZ+qSTT5c4mGXcfvINMDhdVuj1/6TtkSeNjIHWa
    MMKQR7zJ3m9ArZx9e65tMTe2UX7YM5ox/oRSVmF8L0gOG/29LojtJ7MY+Cn2K+kI
    AOO4TlvHs3gLNo4HPWvx5K5++uTSaeSF+jOJJ/ANynLsojzsts4C94h2EoUtJXn8
    /+XzVdu0B0vo39mnis7dfYr3IvO+CyRCg1O94QJB3fRwB6kgIkYoCddAAXt1vTz0
    vTUKfT1QBpetfsySlrayJFBFgCWVKZm0P+ZUcwx6VsLzeNwtAcjARzMjgyoL7csR
    8vgpVdPfbWkVSuwLW1WTTY5OS+U6THgenW8cWK5aIQvuBuvIKkZbKkxNtXtfaay6
    6K5g/EHex4yEdVfx25jBqVi0vz97qVQCq4MN8k27GtiaMO7LFPSnHvaYoId11Bdu
    ZCQ4SAaL/2+ztceGVmFWc/h5jokCHwQYAQIACQUCX3fK+gIbDAAKCRC0ZN/SVca5
    +HtuD/9C9viYi8gnlpXzIUit224tyxlXUDgS5UsiF0K+kRJcuQ3tlaM321IPzcla
    FZeg03VCdlWPkIACf2Mid0g2XAtd/866xlu/T2FJbw1MIae13KuZlpZbh9IZ5Gnf
    aW/AjxNqXxsdfyYFdujLJVwmzx5ePTAiIZvS8869D/F0h2vSc2nHzFT+og6q/1jK
    U/JCryERFKCRqdzH1GKMvDjOgZZanFFixOdvbhn7v3qGolHbTqW9XUbQrgq6SaSa
    IGgGZG0nhWXiKQJSQux+gRycZoeJZrdCN437Y4u/OYZaBhiKcU44YTGPa10pN620
    nY9+8wazPBIifBU7PBJOqcQCM+ldvCKRdCwrFvtKLes+oRXQ58cJKR51VNJjpSp1
    1tx0TcIWH1rYxpekfS0q0OKUKSKg/NUkMHtoecOVZQpKOnUUD06Pu582ooYAA4i0
    2C9uEXdrCOP6ZCNKO26+AD5bukSIN+wKOrsyzmGCxiC2Ae//otvCWlO9+VGLKPxV
    DAyP9yWx7587jfkXSrtGBuToyEJi3ORRzfPBk0gLQZLOcKycniFuLQG0PXPuVDSd
    EhrGdqEDxr+nvKD7POoCSxgm7IORFS+n2lPPlj4xIrY09wfxElR+SBBjBvsQZmXj
    9/zXBe3JOXNsfFYVeNVlkq4BeHWtxz6aOZgcbN2Nd+ykDNCxpQ==
    =2uSU
    -----END PGP PRIVATE KEY BLOCK-----
`+sas.ForNamespaceYAML()

    name := "test-sops"
    cleanUp := func() {
        kapp.Run([]string{"delete", "-a", name})
    }

    cleanUp()
    defer cleanUp()

    logger.Section("deploy", func() {
        kapp.RunWithOpts([]string{"deploy", "-f", "-", "-a", name},
            RunOpts{IntoNs: true, StdinReader: strings.NewReader(yaml1)})
    })

    logger.Section("verify", func() {
        out := kapp.Run([]string{"inspect", "-a", name + "-ctrl", "--raw", "--tty=false", "--filter-kind", "ConfigMap"})

        var cm corev1.ConfigMap

        err := yaml.Unmarshal([]byte(out), &cm)
        if err != nil {
            t.Fatalf("Failed to unmarshal: %s", err)
        }

        if cm.ObjectMeta.Name != "configmap" {
            t.Fatalf(`Expected name to be "configmap" got %#v`, cm.ObjectMeta.Name)
        }

        if cm.Data["key"] != "encrypted-valued" {
            t.Fatalf(`Expected data.key to be "value" got %#v`, cm.Data["key"])
        }
    })
}

/*

$ gpg --gen-key

$ gpg --list-secret-keys --keyid-format LONG
/root/.gnupg/secring.gpg
------------------------
sec   4096R/B464DFD255C6B9F8 2020-10-03
uid                          test test (test) <test@test.com>
ssb   4096R/FEE37B8E2098EDFC 2020-10-03

$ cat foo.yml

data: |
  #@data/values
  ---
  value: encrypted-valued

$ sops --encrypt --pgp B464DFD255C6B9F8 foo.yml

$ gpg --armor --export-secret-keys B464DFD255C6B9F8

*/
