// 在页面加载时从本地存储获取用户名
var username = localStorage.getItem('username') || 'DefaultUsername';

// 更新页面元素
document.getElementById('username').innerText = username;


const connectionCountElement = document.getElementById("connectionCount");

var requestBody = {
    setname: "people",
  };

// 定期获取连接数
function fetchConnectionCount() {
    fetch(`http://localhost:8000/getSetPeopleNum`, {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestBody),
  })
    .then(response => response.json())
    .then(data => {
        // 更新连接数显示
        connectionCountElement.innerText = `${data.num}`;
    })
    .catch(error => {
        console.error('Error fetching connection count:', error);
    });
}

// 初始获取一次连接数
fetchConnectionCount();

// 每隔一段时间获取一次连接数（可以根据实际需求调整时间间隔）
setInterval(fetchConnectionCount, 1000); // 每1秒更新一次连接数
