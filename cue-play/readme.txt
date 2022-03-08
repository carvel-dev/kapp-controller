dk@dks-MacBook-Pro 1-ex % cue export --out yaml ./
no CUE files in .

dk@dks-MacBook-Pro 1-ex % cue export --out yaml ./...
a: true
---
b: true


~~~~~~~~~~~~~~~~~~~~~~~~~~~

dk@dks-MacBook-Pro cue-play % cue export --out yaml ./2-ex
top: true

dk@dks-MacBook-Pro 2-ex % cue export --out yaml .
top: true


~~~~~~~~~~~~~~~~~~~~~~~~~~~

dk@dks-MacBook-Pro cue-play % cue export --out yaml ./3-ex
build constraints exclude all CUE files in ./3-ex:
    cue-play/3-ex/top.cue: no package name

dk@dks-MacBook-Pro cue-play % cue export --out yaml ./3-ex/...
a: true
---
b: true

dk@dks-MacBook-Pro 3-ex % cue export --out yaml
build constraints exclude all CUE files in .:
    cue-play/3-ex/top.cue: no package name


~~~~~~~~~~~~~~~~~~~~~~~~~~

dk@dks-MacBook-Pro cue-play % cue export --out yaml ./4-ex/...
top: true
ex4: true
---
top: true
a: true
ex4: true
---
b: true

dk@dks-MacBook-Pro cue-play % cue export --out yaml ./4-ex/
top: true
ex4: true

~~~~~~~~~~~~~~~~~~~~~~~~~~~


dk@dks-MacBook-Pro cue-play % cue export --out yaml 5-ex
cannot find package "5-ex"

dk@dks-MacBook-Pro cue-play % cue export --out yaml ./5-ex
import failed: cannot find package "top.com/pkga":
    ./5-ex/top.cue:4:2

dk@dks-MacBook-Pro 5-ex % cue export --out yaml ./
top: true
ex5: true
pkga_contents:
  a: true

dk@dks-MacBook-Pro 5-ex % cue export --out yaml
top: true
ex5: true
pkga_contents:
  a: true
