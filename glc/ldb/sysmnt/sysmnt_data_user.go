/**
 * 系统用户
 */
package sysmnt

import (
	"bytes"
	"encoding/gob"
	"encoding/json"

	"github.com/gotoeasy/glang/cmn"
)

type SysUser struct {
	Username    string `json:"username,omitempty"`    // 用户名
	Password    string `json:"password,omitempty"`    // 密码
	Systems     string `json:"systems,omitempty"`     // 能访问的系统（空白或*代表全部，多个用逗号分隔）
	Note        string `json:"note,omitempty"`        // 备注
	CreateDate  string `json:"createdate,omitempty"`  // 创建日期
	UpdateDate  string `json:"updatedate,omitempty"`  // 更新日期
	OldPassword string `json:"oldPassword,omitempty"` // 原密码，修改密码用
	NewPassword string `json:"newPassword,omitempty"` // 新密码，修改密码用
}

func (s *SysUser) ToBytes() []byte {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(s)
	return buffer.Bytes()
}

func (s *SysUser) LoadBytes(data []byte) {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	decoder.Decode(s)
}

func (s *SysUser) ToJson() string {
	bt, _ := json.Marshal(s)
	return cmn.BytesToString(bt)
}
