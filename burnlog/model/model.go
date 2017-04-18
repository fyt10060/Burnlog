// model
package model

import (
	"crypto/sha1"
	"fmt"
	"regexp" // 正则表达式包
	"sort"
	"strings"
	"time"
)

const (
	// format secret
	serUser     = "user!@#"
	serToken    = "token@$@"
	serArticle  = ""
	serComment  = ""
	serMaterial = ""
)

type UserAuthority int

const (
	// user authority
	AuSuperAdmin UserAuthority = iota // 0 move up/down admin authority
	AuAdmin                           // 1 delete all articles, move up/down normal authority
	AuAuthor                          // 2 write/modify/delete self articles, delete self article's comment
	AuNormal                          // 3 write/modify/delete self comments
	AuBanned                          // 4 read
)

type User struct {
	Name           string        `json:"name" redis:"name"`
	Authority      UserAuthority `json:"authority redis:"authority"`
	Uid            string        `json:"-" orm:"PK"`
	SignUpTime     int64         `json:"signup_time" redis:"signup_time"`
	LastSignInTime int64         `json:"last_signin_time" redis:"last_signin_time"`
	CommentList    string        `json:"comment_list" redis:"comment_list"`
	ArticleList    string        `json:"article_list" redis:"article_list"`
	Email          string        `json:"email redis:"email"`
}

type UserInfo struct {
	User
	Token string `json:"token" redis:"token"`
}

func CreateNewUser(name, email string) *User {
	now := GetNowTime()
	return &User{
		Name:       name,
		Email:      email,
		SignUpTime: now,
		Uid:        getFormatId(email, now),
		Authority:  AuNormal,
	}
}

type Article struct {
	Title        string `json:"title"`
	Detail       string `json:"detail"`
	CreateTime   int64  `json:"create_time"`
	UpdateTime   int64  `json:"update_time"`
	Author       string `json:"author"`
	CommentList  string `json:"comment_list"`
	MaterialList string `json:"material_list"`
	ArticleID    string `json:"article_id"`
}

type ArticleList struct {
	List []Article `json:"article_list"`
}

type Material struct {
	MediaID string
	URL     string
}

type MaterialList struct {
	List []Material
}

type Comment struct {
	Content        string
	Author         string
	CreatTime      int64
	UpdateTime     int64
	CommentID      string
	SuperCommentID string
}

type CommentList struct {
	List []Comment
}

func GetNewArticle(title, detail, author string, materialList string) *Article {
	now := GetNowTime()
	return &Article{
		Title:        title,
		Detail:       detail,
		Author:       author,
		ArticleID:    getFormatId("article", now),
		CreateTime:   now,
		UpdateTime:   now,
		MaterialList: materialList,
	}
}

func GetNowTime() int64 {
	return time.Now().Unix()
}

func getFormatId(value string, time int64) string {
	theTime := fmt.Sprintf("%d", time)
	array := []string{value, theTime}
	sort.Strings(array)

	idString := strings.Join(array, "")

	shId := sha1.New()
	shId.Write([]byte(idString))

	shByte := shId.Sum(nil)
	id := fmt.Sprintf("%x", shByte)

	return id
}

func GetUserToken(account string) string {
	return getFormatId(account, GetNowTime())
}

func CheckEmailFormat(email string) bool {
	if email != "" {
		if isOk, _ := regexp.MatchString("^[_a-z0-9-]+(\\.[_a-z0-9-]+)*@[a-z0-9-]+(\\.[a-z0-9-]+)*(\\.[a-z]{2,4})$", email); isOk {
			return true
		}
	}
	return false
}
