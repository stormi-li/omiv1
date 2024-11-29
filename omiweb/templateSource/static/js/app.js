// 跳转到请求页面
function redirectToRequest() {
    window.location.href = "request_test.html"; // 替换为实际请求页面的 URL
}

// HTTP 请求处理
async function handleSearch(event) {
    if (event.key === "Enter") {
        const query = document.querySelector(".search-box").value;
        try {
            const response = await fetch(query);
            if (response.ok) {
                const data = await response.text();
                alert(`HTTP 响应: ${data}`);
            } else {
                alert("HTTP 请求失败，请稍后重试！");
            }
        } catch (error) {
            alert("HTTP 请求出错！", error);
        }
    }
}

// WebSocket 请求处理
let websocket;

function handleWebSocket(event) {
    if (event.key === "Enter") {
        const query = document.querySelector(".ws-search-box").value;
        if (!websocket || websocket.readyState !== WebSocket.OPEN) {
            websocket = new WebSocket(query);

            websocket.onopen = () => {
                websocket.send("Hello WebSocket!");
            };

            websocket.onmessage = (event) => {
                alert(`WebSocket 响应: ${event.data}`);
            };

            websocket.onerror = (error) => {
                alert(`WebSocket 错误: ${error.message}`);
            };

            websocket.onclose = () => {
                console.log("WebSocket 连接已关闭");
            };
        } else {
            websocket.send("Hello WebSocket!");
        }
    }
}

// 跳转功能处理
function handleNavigate(event) {
    if (event.key === "Enter") {
        const query = document.querySelector(".navigate-box").value;
        if (query) {
            window.location.href = query; // 跳转到输入的 URL
        } else {
            alert("请输入有效的 URL！");
        }
    }
}
