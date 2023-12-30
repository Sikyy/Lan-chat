// login.js

function handleLogin(event) {
    event.preventDefault(); // 阻止表单默认提交行为
  
    var username = document.getElementById('username').value;
    console.log('Username:', username); // 添加调试输出
  
    // 处理登录逻辑，例如将用户名保存到本地存储
    localStorage.setItem('username', username);
  
    // 使用前端路由进行跳转
    // 假设使用的是 "/chat" 作为聊天界面的路由
    window.location.href = "/chat";
  }
  