package models

import (
	"encoding/json"
	"fmt"
	"github.com/yanshicheng/ikubeops-gin-demo/common/model"
	"time"
)

// Hobby 枚举类型，代表不同的运动爱好
type Hobby int

const (
	Basketball Hobby = iota + 1 // 篮球
	Football                    // 足球
	Badminton                   // 羽毛球
)

// String 方法返回 Hobby 枚举值的字符串表示
func (h Hobby) String() string {
	switch h {
	case Basketball:
		return "basketball"
	case Football:
		return "football"
	case Badminton:
		return "badminton"
	default:
		return "unknown"
	}
}

// MarshalJSON 为 Hobby 类型实现自定义的 JSON 序列化方法
func (h Hobby) MarshalJSON() ([]byte, error) {
	return json.Marshal(h.String())
}

// UnmarshalJSON 为 Hobby 类型实现自定义的 JSON 反序列化方法
func (h *Hobby) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("JSON 解析错误: %w", err)
	}
	switch str {
	case "basketball":
		*h = Basketball
	case "football":
		*h = Football
	case "badminton":
		*h = Badminton
	default:
		return fmt.Errorf("无效的 Hobby 值: %s", str)
	}
	return nil
}

// Date 是一个自定义日期类型，用于处理 JSON 中的日期时间格式
type Date time.Time

// MarshalJSON 实现自定义日期时间的 JSON 序列化
func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(d).Format("2006-01-02 15:04:05"))
}

// UnmarshalJSON 实现自定义日期时间的 JSON 反序列化
func (d *Date) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("JSON 解析错误: %w", err)
	}
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return fmt.Errorf("时间解析错误: %w", err)
	}
	*d = Date(t)
	return nil
}

// Book 定义书籍模型
type Book struct {
	model.Model
	BookName string `json:"book_name" binding:"required,alphanumunicode" gorm:"uniqueIndex;size:20;type:varchar(20);not null;column:book_name;comment:书籍名称，全局唯一，只包含中英文字符"`
	BookId   uint   `json:"book_id" binding:"required" gorm:"type:bigint;not null;column:book_id;comment:书籍ID，数字类型"`
	Details  string `json:"details,omitempty" gorm:"type:text;column:details;comment:书籍详情，大文本类型"`
}

// Author 定义作者模型
type Author struct {
	model.Model
	Name     string  `json:"name" binding:"required,alphanumunicode" gorm:"uniqueIndex;size:20;type:varchar(20);not null;column:name;comment:作者名称，全局唯一，只包含中英文字符"`
	Age      int     `json:"age" binding:"required" gorm:"type:int;not null;column:age;comment:作者年龄"`
	Gender   bool    `json:"gender" binding:"required" gorm:"type:boolean;default:true;not null;column:gender;comment:作者性别，默认为 true"`
	Email    string  `json:"email" binding:"required,email" gorm:"uniqueIndex;type:varchar(100);not null;column:email;comment:作者邮箱，需进行邮箱验证"`
	Hobby    Hobby   `json:"hobby" binding:"required" gorm:"type:int;not null;column:hobby;comment:作者的爱好"`
	Birthday Date    `json:"birthday" binding:"required" gorm:"type:date;not null;column:birthday;comment:作者生日"`
	AddAt    Date    `json:"add_at" binding:"required" gorm:"type:datetime;not null;column:add_at;comment:记录添加时间"`
	Address  Address `json:"address" binding:"required" gorm:"foreignKey:AuthorID;column:address_id;comment:作者地址信息，一对一绑定"`
	Books    []Book  `json:"books" gorm:"foreignKey:AuthorID;column:author_id;comment:作者所拥有的书籍"`
}

// Address 定义地址模型
type Address struct {
	model.Model
	AuthorID uint   `json:"author_id" gorm:"primaryKey;type:bigint;not null;column:author_id;comment:关联的作者ID"`
	Country  string `json:"country" binding:"required,alphanumunicode" gorm:"size:6;type:varchar(6);not null;column:country;comment:国家名称，最大长度6"`
	State    string `json:"state" binding:"required,alphanumunicode" gorm:"size:12;type:varchar(12);not null;column:state;comment:省份名称，最大长度12"`
	City     string `json:"city" binding:"required,alphanumunicode" gorm:"size:12;type:varchar(12);not null;column:city;comment:城市名称，最大长度12"`
	Details  string `json:"details" binding:"required,alphanumunicode" gorm:"uniqueIndex:idx_address_details;size:64;type:varchar(64);not null;column:details;comment:详细地址，与市字段联合唯一，最大长度64"`
}

// Tag 定义标签模型
type Tag struct {
	model.Model
	Name  string `json:"name" binding:"required,alphanumunicode" gorm:"uniqueIndex;size:20;type:varchar(20);not null;column:name;comment:标签名称，全局唯一，只包含中英文字符"`
	Books []Book `json:"books" gorm:"many2many:book_tags;column:tag_id;comment:标签与书籍的多对多关系"`
}
