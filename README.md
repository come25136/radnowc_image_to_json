# radnowc image to json
[気象庁｜高解像度降水ナウキャスト](https://www.jma.go.jp/jp/realtimerad/index.html)タイル(ラスターデータ)を数値データに変換するサーバーです

# 使い方
http://localhost:8000/?date=202011200500&x=34&y=31&z=6 的な感じで使います

## 起動引数
address = :3000

## パラメーター
date = YYYYMMDDHHmm (mmは5分単位)
x = 19-44 (ナウキャストマップの西から東にかけて数値が大きくなります)
y = 18-46 (ナウキャストマップの北から南にかけて数値が大きくなります)
z = 6 (ナウキャストの仕様的には1-6ですが、XYのバリデートの都合上6固定です)

### 備考
x,yはchromeなどの開発者ツールを使用して、ナウキャストのタイル番号を調べてください (こういうやつ↓)
https://www.jma.go.jp/jp/realtimerad/highresorad_tile/HRKSNC/{date}/{date}/zoom{z}/{x}_{y}.png

1タイル = 256px x 256px

## レスポンス
入る値は
```
0
1
5
10
20
30
50
80
100
```
の9種類で、単位はmm/hです(ナウキャストマップ右下に表示されているサンプルと同じ)

### サンプル
JS的に書くと`data[yPixel][xPixel]`な感じにピクセルデータが入っています

<details>
<summary>めっちゃ長いレスポンスデータ</summary>

```
Chrome のレンダリングが落ちるので省略
```
</details>

# 座標変換について
## ピクセル単位で地図座標に変換する式
レスポンスをmapboxなどにプロットする場合は、この式を使用しないと位置がズレます

おまじない的な数式を[豐多摩研究所](https://twitter.com/intent/user?user_id=1529995508)さんに教えていただきました
> 左上基準なので、タイルzoom{z}/{x}_{y}.pngの画像の左からu, 上からv番目（いずれも0から数えて）のピクセルについて、
> 中心の経度 = (x + u / 256 + 1 / 512) / (2 ^ z) × 70度 + 100度
> 中心の緯度 = 61度 - 54度 × (y + v / 256 + 1 / 512) / (2 ^ z)

u,vが分かりにくいですが、要はタイル内のピクセル位置(0-255)です

### プロット例
見栄えの都合上、市松模様にしています
[imvue.csb.app](https://imvue.csb.app/)
![image](https://user-images.githubusercontent.com/12409412/99955940-800d5a80-2dc8-11eb-94c7-1213504e4d67.png)
![image](https://user-images.githubusercontent.com/12409412/99956025-a8955480-2dc8-11eb-9ee5-770235648399.png)
