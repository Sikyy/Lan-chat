// 在页面加载时从本地存储获取用户名
var username = localStorage.getItem('username') || 'DefaultUsername';

// 更新页面元素
document.getElementById('username').innerText = username;