/* 创建 WebSocket 连接 */
const socket = new WebSocket("ws://" + window.location.host + "/ws");


/* 监听消息事件 */
socket.onmessage = function (event) {
    const chatDiv = document.getElementById("chat-messages");
    chatDiv.innerHTML += "<p>" + event.data + "</p>";
    chatDiv.scrollTop = chatDiv.scrollHeight;
};

/* 发送按钮点击事件 */
document.getElementById("send-button").addEventListener("click", function () {
    const input = document.getElementById("input");
    const message = input.value;
    socket.send(message);
    input.value = "";
});

/* 输入框按键事件 */
document.getElementById("input").addEventListener("keyup", function (event) {
    if (event.key === "Enter") {
        document.getElementById("send-button").click();
    }
});