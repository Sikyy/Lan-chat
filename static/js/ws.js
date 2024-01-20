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
        } else if (message.type === "message" && message.channel === "chat") {
            // 处理 chat 频道中的消息
            displayReceivedMessage(message.content);
        } else {
            console.log('未知的消息类型或频道：', message);
        }
    } catch (error) {
        console.error('解析传入消息时出错：', error);
    }
});

// 处理错误
socket.addEventListener('error', (event) => {
    console.error('WebSocket 发生错误：', event);
});

function displayReceivedMessage(message) {
    var chatDiv = document.getElementById('chat');
    var newMessage = document.createElement('div');
    var spanMessage = document.createElement('div');
    var avatarSpan = document.createElement('span');
    var avatarImg = document.createElement('img');

    // 设置头像
    avatarImg.src = 'static/images/avatar.jpg';
    avatarImg.classList.add('avatar');
    avatarSpan.appendChild(avatarImg);
    avatarSpan.style.userSelect = 'none';

    // 创建 span 元素用于包装消息文本
    var messageText = document.createElement('span');
    messageText.innerHTML = message.replace(/\n/g, '<br>');
    messageText.classList.add('message-text');
    messageText.style.whiteSpace = 'pre-line';
    messageText.style.textAlign = 'left';

    // 创建 span 元素用于显示时间戳
    var timestampSpan = document.createElement('span');
    timestampSpan.textContent = getNowFormatDate();
    timestampSpan.classList.add('message-timestamp');
    timestampSpan.style.userSelect = 'none';

    spanMessage.appendChild(messageText);
    spanMessage.appendChild(timestampSpan);
    spanMessage.classList.add('message-span');
    newMessage.appendChild(spanMessage);
    newMessage.appendChild(avatarSpan);
    newMessage.classList.add('message');
    chatDiv.appendChild(newMessage);

    // 滚动聊天框到底部
    chatDiv.scrollTop = chatDiv.scrollHeight;
}