const socket = new WebSocket("ws://" + window.location.host + "/ws");

// 处理连接打开事件
socket.addEventListener('open', (event) => {
    console.log('WebSocket连接已打开');
});

// 处理连接关闭事件
socket.addEventListener('close', (event) => {
    console.log('WebSocket连接已关闭');
});

// 处理接收到的消息
socket.addEventListener('message', (event) => {
    try {
        const message = JSON.parse(event.data);

        if (message.type === "connectionCount") {
            const connectionCountElement = document.getElementById("connectionCount");
            connectionCountElement.innerText = `当前在线连接数：${message.content}`;
        } else if (message.type === "userMessage") {
            // 处理用户发送的消息
            displayUserMessage(message.content);
        }
    } catch (error) {
        console.error('解析传入消息时出错：', error);
    }
});

// 处理错误
socket.addEventListener('error', (event) => {
    console.error('WebSocket 发生错误：', event);
});