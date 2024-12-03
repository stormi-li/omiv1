document.addEventListener('DOMContentLoaded', function () {
    const httpInput = document.querySelector('.http-url');
    const httpResponse = document.querySelector('.http-response');
    const websocketInput = document.querySelector('.websocket-url');
    const websocketResponse = document.querySelector('.websocket-response');

    // 隐藏所有返回框
    function hideAllResponses() {
        httpResponse.style.display = 'none';
        websocketResponse.style.display = 'none';
    }

    // 显示 HTTP 响应框并发起 HTTP 请求
    httpInput.addEventListener('focus', function () {
        hideAllResponses(); // 隐藏其他响应框
        httpResponse.style.display = 'block';
    });

    httpInput.addEventListener('keydown', function (event) {
        if (event.key === 'Enter') {
            const url = httpInput.value.trim();
            if (url) {
                fetch(url)
                    .then(response => response.text())
                    .then(data => {
                        appendContent(httpResponse, `请求地址: ${url}<br>返回数据:${data}`);
                    })
                    .catch(err => {
                        appendContent(httpResponse, `请求地址: ${url}<br>请求失败: ${err.message}`);
                    });
            } else {
                appendContent(httpResponse, '请输入有效的 HTTP 地址。');
            }
        }
    });

    httpInput.addEventListener('blur', function () {
        if (!httpInput.value.trim()) {
            hideAllResponses();
        }
    });

    // 显示 WebSocket 响应框并发起 WebSocket 请求
    websocketInput.addEventListener('focus', function () {
        hideAllResponses(); // 隐藏其他响应框
        websocketResponse.style.display = 'block';
    });

    websocketInput.addEventListener('keydown', function (event) {
        if (event.key === 'Enter') {
            const url = websocketInput.value.trim();
            if (url) {
                try {
                    const ws = new WebSocket(url);
                    ws.onopen = () => {
                        appendContent(websocketResponse, `WebSocket 连接已建立: ${url}`);
                    };
                    ws.onmessage = (event) => {
                        appendContent(websocketResponse, `接收到数据:${event.data}`);
                    };
                    ws.onerror = (err) => {
                        appendContent(websocketResponse, `WebSocket 错误: ${err.message}`);
                    };
                    ws.onclose = () => {
                        appendContent(websocketResponse, 'WebSocket 连接已关闭。');
                    };
                } catch (err) {
                    appendContent(websocketResponse, `WebSocket 请求失败: ${err.message}`);
                }
            } else {
                appendContent(websocketResponse, '请输入有效的 WebSocket 地址。');
            }
        }
    });

    websocketInput.addEventListener('blur', function () {
        if (!websocketInput.value.trim()) {
            hideAllResponses();
        }
    });

    // 追加内容到显示框并自动换行，并滚动到最底部
    function appendContent(element, content) {
        element.innerHTML += `<br>${content}`; // 自动拼接换行符
        element.scrollTop = element.scrollHeight; // 滚动到最底部
    }
});