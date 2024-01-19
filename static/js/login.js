// login.js

function handleLogin(event) {
  event.preventDefault();

  var username = document.getElementById('username').value;

  var requestBody = {
    username: username,
    setname: "people",
  };

  fetch(`http://localhost:8000/login/${username}`, {
      method: 'POST',
      headers: {
          'Content-Type': 'application/json',
      },
      body: JSON.stringify(requestBody),
  })
  .then(response => response.json())
  .then(data => {
      console.log('Response from server:', data);

      if (data.success) {
          console.log(`用户 ${username} 登录成功，集合 ${requestBody.setname} 内已录入`);
          
          // 更新本地存储中的用户名
          localStorage.setItem('username', username);

          // 页面跳转逻辑放在这里
          window.location.href = "/chat";
      } else {
          console.error('用户登录失败:', data.message);
      }
  })
  .catch(error => {
      console.error('Error during login:', error);
  });
}
