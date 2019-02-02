# cgrep

color enhanced grep

## how to use
```
git grep --color=always . | cgrep '(".*")' 'yellow' | cgrep '(%[0-9]*[dsf])' 'magenta' | cgrep '([0-9]+)' green
```

## how to install
```
go get -u "github.com/umaumax/cgrep"
```

## color format info
* [mgutz/ansi: Small, fast library to create ANSI colored strings and codes\. \[go, golang\]]( https://github.com/mgutz/ansi )

## NOTE
* 色が重なった場合には後勝ち
* `()`に対応して，色がつく，`,`区切りで色を指定
  * 色指定を空文字にすると色設定をskip

## TODO
* ansi文字出力の効率化
* helpの記述
* exampleの記述
