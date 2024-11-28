// 获取所有具有 'show-datail' 类的按钮元素
import omihttp from "./omihttp.js";

const back = document.querySelector(".back");
const dataContainer = document.querySelector(".data-container");
const detailContainer = document.querySelector(".detail-container");


back.addEventListener('click', function () {
    dataContainer.style.display = "block";
    detailContainer.style.display = "none";

})

// 页面加载时渲染数据
window.addEventListener('DOMContentLoaded', async renderDataContainer);

func renderDataContainer(){
    try {
        // 获取嵌套的 map 数据
        const nodes = await omihttp.get('/GetNodes');
        // 渲染数据
        let lastServerName = "null"; // 用于跟踪上一个行的 ServerName
        let rowClass = "table-row1"
        nodes.forEach((node) => {
            if (node.Type) {
                rowClass = "table-row1"
            } else {
                rowClass = "table-row2"
            }

            // 创建表格行
            const row = document.createElement("div");
            row.className = `table-row ${rowClass}`;
            row.innerHTML = `
        <span class="name">${node.ServerName}</span>
        <span class="addr">${node.Address}</span>
        <span class="weight">${node.Weight}</span>
        <button class="show-datail">显示详情</button>
    `;

            // 添加到容器中
            dataContainer.appendChild(row);

        });
    } catch (error) {
        console.error('Error fetching nodes:', error);
    }
}

const detailServerName = document.querySelector('.detail-serverName');
const detailAddress = document.querySelector('.detail-address');
// 监听点击事件
dataContainer.addEventListener('click', async (event) => {
    // 检查是否点击了 "显示详情" 按钮
    if (event.target.classList.contains('show-datail')) {
        // 获取所在行的元素
        const row = event.target.closest('.table-row');

        // 提取 addr 和 name
        const addr = row.querySelector('.addr').textContent;
        const name = row.querySelector('.name').textContent;
        detailServerName.innerHTML = name
        detailAddress.innerHTML = addr
        dataContainer.style.display = "none";
        detailContainer.style.display = "block";
        try {
            const data = await omihttp.get(`/GetNodeInfo?name=${name}&address=${addr}`);
            renderDetailContainer(data)
        } catch (error) {
            console.error('Error fetching nodes:', error);
        }
    }
});
function renderDetailContainer(data) {
    const commandsContainer = document.querySelector('.commands');
    const detailsContainer = document.querySelector('.details');

    // 清空现有内容
    commandsContainer.innerHTML = '';
    detailsContainer.innerHTML = '';

    // 处理 MessageHandlers，按逗号分割并渲染
    if (data.MessageHandlers) {
        const handlers = data.MessageHandlers.split(',').map(command => command.trim());
        handlers.forEach(command => {
            const commandItem = document.createElement('div');
            commandItem.className = 'command-item';
            commandItem.innerHTML = `
                <div class="command">${command}</div>
                <input type="text" placeholder="Enter message for ${command}" data-command="${command}">
            `;
            commandsContainer.appendChild(commandItem);

            // 获取当前input元素
            const inputElement = commandItem.querySelector('input');

            // 为输入框添加失焦事件 (blur)
            inputElement.addEventListener('blur', function () {
                sendMessage(command, inputElement.value);
            });

            // 为输入框添加回车事件 (keydown)
            inputElement.addEventListener('keydown', function (event) {
                if (event.key === 'Enter') {
                    sendMessage(command, inputElement.value);
                    inputElement.blur();
                }
            });
        });
    }

    // 渲染其他字段到 details
    for (const [key, value] of Object.entries(data)) {
        if (key !== 'MessageHandlers') {
            const detailItem = document.createElement('div');
            detailItem.className = 'detail-item';
            detailItem.innerHTML = `
                <span>${key}</span>
                <span>${value}</span>
            `;
            detailsContainer.appendChild(detailItem);
        }
    }
}

// 发送请求的函数
function sendMessage(command, message) {
    // 这里根据实际情况修改请求的 URL 和数据
    try {
        omihttp.get(`/SendMessage?name=${detailServerName.innerHTML}&address=${detailAddress.innerHTML}&command=${command}&message=${message}`)
    } catch (e) { 
        console.log(e)
    }
}