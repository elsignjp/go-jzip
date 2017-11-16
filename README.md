# jzip

Because this is a tool to convert data provided by Japan Post, this section is written only in Japanese.

本ツールは[日本郵便株式会社](http://www.post.japanpost.jp)が公開している[郵便番号データ](http://www.post.japanpost.jp/zipcode/download.html)をJSONやSQL形式に展開するものです。


## Description

公開されているデータのうち、不要な文字列の除去、複数行に分割されているデータの結合等、最低限必要と思われる加工を施し、JSON形式やSQL形式で出力を行います。

## Install

```
go get -u github.com/elsignjp/go-jzip/cmd/jzip
```

## Usage

HELP参照 - `jzip help`

```
NAME:
   jzip - will convert zip codes of Japan Post

USAGE:
   jzip [global options]

VERSION:
   0.0.1

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  Specify the config file to be read
   --json, -j                Output in JSON format
   --sql, -q                 Output in SQL format
   --file value, -f value    Specify the ZIP file to be read
   --out value, -o value     Specify the file path to output
   --table value             Specify the table name to output (format:SQL)
   --fld-local value         Field name of the local public organization code to be output
   --fld-zip value           Field name of the zip code to be output
   --fld-pref value          Field name of prefecture to output
   --fld-city value          Field name of city to output
   --fld-town value          Field name of the town area to output
   --help, -h                show help
   --version, -v             print the version
```

### 設定ファイル
`config`オプションで指定する設定ファイルは下記のフォーマットで作成します。  
設定ファイルの項目はすべてコマンドオプションで指定可能なので、コマンド単体での実行も可能です。

```
ZipFile="ken_all.zip"     #郵便番号データのZIPファイルのパス(ken_all.zipなど)
Format="sql"              #出力形式(json,sql)
Output="jzip.sql"         #出力するファイルパス
TableName="zipcode"       #SQL形式の場合のテーブル名
LocalCodeField="local"    #地方公共団体コードの出力フィールド名
ZipCodeField="zip"        #郵便番号(7桁)の出力フィールド名
PrefecturesField="pref"   #都道府県の出力フィールド名
CityField="city"          #市区町村の出力フィールド名
TownField="town"          #町域の出力フィールド名
```

## Licence

[MIT](https://github.com/elsignjp/go-jzip/blob/master/LICENCE)

## Author

[elsignjp](https://github.com/elsignjp)
