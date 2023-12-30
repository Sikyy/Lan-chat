package data

import "encoding/json"

type Data struct {
	Ip       string   `json:"ip"`
	User     string   `json:"user"`
	From     string   `json:"from"`
	Type     string   `json:"type"`
	Content  string   `json:"content"`
	UserList []string `json:"user_list"`
	Username string   `json:"username"`
}

// ToJSON 方法将 Data 结构体转换为 JSON 字节切片
func (d *Data) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}
