# cgrep

color enhanced grep

## how to use
```
# regex
git grep --color=always . | cgrep '(".*")' 'yellow' | cgrep '(%[0-9]*[dsf])' 'magenta' | cgrep '([0-9]+)' green

# fixed string
echo 'ヽ(*゜д゜)ノ' | cgrep -F '゜' blue
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
* 一文字ごとに色がついているが，それは不要では?
  * `fzf`は1文字ごとではなくとも該当した文字の色を保持する(`fzy`は保持しない)

## TODO
* ansi文字出力の効率化(1文字ごとに色情報を出力している)
