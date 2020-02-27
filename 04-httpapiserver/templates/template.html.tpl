<!doctype html>
<html lang="ja">

<head>
    <meta charset="UTF-8">
    <meta name="description" content="celan">
    <meta name="viewport" content="width=device-width">
    <title>{{.HostName}}</title>
    <script>
        window.onload = function () {
            setInterval("update()", 1000);
            update();
        }

        function update() {
            // https://ja.javascript.info/xmlhttprequest

            // 1. new XMLHttpRequest オブジェクト作成
            let xhr = new XMLHttpRequest();

            // 2. 設定: URL .../api に対する GET-リクエスト
            xhr.open('GET', '/api');

            // 3. ネットワーク経由でリクエスト送信
            xhr.send();

            // 4. レスポンスを受け取った後のコールバック
            xhr.onload = function () {
                if (xhr.status != 200) {

                    // エラーの場合
                    alert(`Error ${xhr.status}: ${xhr.statusText}`); // e.g. 404: Not Found

                } else {
                    // 成功
                    // alert(`Done, got ${xhr.response.length} bytes`); // responseText is the server

                    var data = xhr.response;

                    var dom = document.getElementById('message');
                    dom.innerHTML = data;

                    document.getElementById('time').innerHTML = getDateString(new Date());
                }
            };
        }

        // https://tagamidaiki.com/javascript-0-chink/
        function getDateString(date) {

            var yyyy = date.getFullYear();
            var mm = toDoubleDigits(date.getMonth() + 1);
            var dd = toDoubleDigits(date.getDate());
            var hh = toDoubleDigits(date.getHours());
            var mi = toDoubleDigits(date.getMinutes());
            var sq = toDoubleDigits(date.getSeconds());
            return hh + ':' + mi + ':' + sq;
        };

        function toDoubleDigits(num) {
            num += "";
            if (num.length === 1) {
                num = "0" + num;
            }
            return num;
        }
    </script>
</head>

<body>
    <span id="time">xxx:xxx:xxx</span> 時点での {{.HostName}} のプロセス<br>
    <br>
    <font color="red">
        <div id="message"></div>
    </font>
    <br>
    <br>
    <br>
    <br>
</body>

</html>