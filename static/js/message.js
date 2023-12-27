//实现发送消息的功能
var isComposing = false;
var chatBox = document.getElementById('chat-box');

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
  

// 添加事件监听器，当输入框内容变化时调整高度
chatBox.addEventListener('input', function () {
    adjustInputHeight(chatBox);
});

function displayUserMessage(content) {
    // 在界面上显示用户发送的消息的逻辑
}

function adjustInputHeight(inputElement) {
    var computedStyle = window.getComputedStyle(inputElement);
    var lineHeight = parseInt(computedStyle.lineHeight, 10);

    var lines = inputElement.value.split('\n').length;
    var newHeight = lines * lineHeight;

    // 设置最小和最大高度限制
    var minHeight = lineHeight; // 最小高度为一行的高度
    var maxHeight = 200; // 最大高度

    newHeight = Math.min(maxHeight, Math.max(minHeight, newHeight));

    inputElement.style.height = newHeight + 'px';
}

function sendMessage() {
    var messageInput = document.getElementById('chat-box'); // 获取聊天框输入框元素
    var message = messageInput.value.trim(); // 获取输入框的值

    if (message !== '') { // 检查消息是否为空
        var chatDiv = document.getElementById('chat'); // 获取聊天框元素
        var newMessage = document.createElement('div'); // 创建新的消息元素
        var spanMessage = document.createElement('div');//创建新的消息元素用于包装span
        var avatarSpan = document.createElement('span'); // 创建包装头像的 span 元素
        var avatarImg = document.createElement('img'); // 创建头像图片

        // 设置头像
        avatarImg.src = 'static/images/avatar.jpg'; // 设置头像图片的路径
        avatarImg.classList.add('avatar'); // 添加自定义样式类
        avatarSpan.appendChild(avatarImg); // 将头像添加到包装头像的 span 元素中
        avatarSpan.style.userSelect = 'none'; // 不可被选择

        // 创建 span 元素用于包装消息文本
        var messageText = document.createElement('span');
        messageText.innerHTML = message.replace(/\n/g, '<br>'); // 手动替换换行符为 <br>
        messageText.classList.add('message-text'); // 添加自定义样式类
        messageText.style.whiteSpace = 'pre-line'; // 保留换行符和连续空格
        messageText.style.textAlign = 'left';
        

        // 创建 span 元素用于显示时间戳
        var timestampSpan = document.createElement('span');
        timestampSpan.textContent = getNowFormatDate(); // 设置时间戳内容
        timestampSpan.classList.add('message-timestamp'); // 添加自定义样式类
        timestampSpan.style.userSelect = 'none'; // 不可被选择
 
        spanMessage.appendChild(messageText); // 将消息文本添加到新的消息元素中
        spanMessage.appendChild(timestampSpan); // 将时间戳添加到新的消息元素中
        spanMessage.classList.add('message-span'); // 添加自定义样式类
        newMessage.appendChild(spanMessage); // 将 span 元素添加到newMessage中
        // 将头像 span 元素和消息容器添加到新的消息元素中
        newMessage.appendChild(avatarSpan);
        newMessage.classList.add('message'); // 添加新的消息元素样式类
        chatDiv.appendChild(newMessage); // 将消息添加到聊天框中

        // 发送消息到服务器
        socket.send(JSON.stringify({ type: "message", content: message }));

        // 清空输入框
        messageInput.value = '';

        // 将输入框高度设置为初始状态
        adjustInputHeight(messageInput);

        // 添加右键点击事件监听器
        spanMessage.addEventListener('contextmenu', function (event) {
            event.preventDefault(); // 阻止默认的右键菜单

            // 获取鼠标点击位置
            var mouseX = event.clientX;
            var mouseY = event.clientY;

            // 在鼠标位置显示自定义菜单
            showContextMenu(mouseX, mouseY);
        });

        // 滚动聊天框到底部
        chatDiv.scrollTop = chatDiv.scrollHeight;
    }
}

// 处理键盘按下事件的函数
function handleKeyDown(event) {
    if (event.key === 'Enter') { // 检查按下的键是否为 Enter 键
        // 阻止默认的 Enter 键行为
        event.preventDefault();

        // 获取输入框的值
        var messageInput = document.getElementById('chat-box'); // 获取聊天框输入框元素
        var message = messageInput.value; // 获取输入框的值

        // 检查输入框是否为空
        if (!isComposing && message.trim() !== '') { // 检查是否正在输入中且消息不为空
            // 调用发送消息的函数
            sendMessage();
        }
    }
}

//获取当前时间戳 格式22:01
function getNowFormatDate() {
    var date = new Date();
    var seperator1 = ":";
    var hour = date.getHours();
    var minute = date.getMinutes();
    if (minute < 10) {
        minute = '0' + minute;
    }
    var currentdate = hour + seperator1 + minute;
    return currentdate;
}

// 显示右键菜单的函数
function showContextMenu(x, y) {
    var contextMenu = document.getElementById('context-menu');

    // 设置菜单位置
    contextMenu.style.left = x + 'px';
    contextMenu.style.top = y + 'px';

    // 显示菜单
    contextMenu.style.display = 'block';

    // 添加菜单项的点击事件监听器
    document.getElementById('reply').addEventListener('click', handleReplyClick);
    document.getElementById('translate').addEventListener('click', handleTranslateClick);
    document.getElementById('copy').addEventListener('click', handleCopyClick);
    document.getElementById('forward').addEventListener('click', handleForwardClick);
    document.getElementById('delete').addEventListener('click', handleDeleteClick);

    // 点击其他地方时隐藏菜单
    document.addEventListener('click', hideContextMenu);
}

function hideContextMenu() {
    var contextMenu = document.getElementById('context-menu');
    contextMenu.style.display = 'none';

    // 移除菜单项的点击事件监听器
    document.getElementById('reply').removeEventListener('click', handleReplyClick);
    document.getElementById('translate').removeEventListener('click', handleTranslateClick);
    document.getElementById('copy').removeEventListener('click', handleCopyClick);
    document.getElementById('forward').removeEventListener('click', handleForwardClick);
    document.getElementById('delete').removeEventListener('click', handleDeleteClick);

    // 移除点击其他地方时隐藏菜单的事件监听器
    document.removeEventListener('click', hideContextMenu);
}

// 具名函数处理点击事件
function handleReplyClick() {
    alert('回复被点击');
    // 处理回复的逻辑
}

function handleTranslateClick() {
    alert('翻译被点击');
    // 处理翻译的逻辑
}

function handleCopyClick() {
    // 获取 spanMessage 元素
    var spanMessage = document.querySelector('.message-span');

    // 获取 messageText 元素
    var messageText = spanMessage.querySelector('.message-text');

    // 获取要复制的文本
    var textToCopy = messageText.textContent;

    // 创建一个临时的文本区域元素
    var tempTextArea = document.createElement("textarea");
    tempTextArea.value = textToCopy;

    // 将临时文本区域元素添加到 DOM 中
    document.body.appendChild(tempTextArea);

    // 选中文本
    tempTextArea.select();

    try {
        // 执行复制命令
        document.execCommand('copy');

    } catch (err) {
        // 如果复制命令失败，可以在这里处理异常
        console.error('复制失败:', err);
    } finally {
        // 移除临时文本区域元素
        document.body.removeChild(tempTextArea);
    }
}

function handleForwardClick() {
    alert('转发被点击');
    // 处理转发的逻辑
}

function handleDeleteClick() {
    alert('删除被点击');
    // 处理删除的逻辑
}


// 处理中文输入法完成输入事件的函数
function handleCompositionEnd() {
    isComposing = false;
    
    // 在这里可以执行一些逻辑，比如更新预览等
}

// 处理中文输入法开始输入事件的函数
function handleCompositionStart() {
    isComposing = true;
}

// 在输入框上添加事件监听器
var chatBox = document.getElementById('chat-box'); // 获取聊天框输入框元素
chatBox.addEventListener('compositionend', handleCompositionEnd); // 添加中文输入法完成输入事件监听器
chatBox.addEventListener('compositionstart', handleCompositionStart); // 添加中文输入法开始输入事件监听器
chatBox.addEventListener('keydown', handleKeyDown); // 添加键盘按下事件监听器