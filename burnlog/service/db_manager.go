// redismanager
package service

import (
	"database/sql"
	"fmt"
	//	"os"
	//	"os/signal"
	//	"syscall"
	//	"time"

	"github.com/astaxie/beego/orm"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"

	"Burnlog/burnlog/model"
)

const (
	// redis keys
	keyUserEmail = "key-user-email"
	keyPassword  = "key-password"
)

var (
	MaxPoolSize = 20
	redisPool   chan redis.Conn
	sqlPool     *sql.DB
)

func init() {
	sqlPool, err := sql.Open("mysql", "burnlog:1234567burnlog@/burnlogDB?charset=utf8")
	if err != nil {
		fmt.Printf("connect to mysql error: %s", err.Error())
	}
	sqlPool.SetMaxOpenConns(1000)
	sqlPool.SetMaxIdleConns(500)
	err = sqlPool.Ping()
	if err != nil {
		panic(err)
	}

	err = orm.RegisterDataBase("default", "mysql", "burnlog:1234567burnlog@/burnlogdb?charset=utf8", 30)
	if err != nil {
		panic(err)
	}

	orm.RegisterModel(new(model.User), new(model.TokenList), new(model.AccountList))
	orm.RunSyncdb("default", false, true)
}

// redis pool related
func putRedis(conn redis.Conn) {
	if redisPool == nil {
		redisPool = make(chan redis.Conn, MaxPoolSize)
	}
	if len(redisPool) >= MaxPoolSize {
		conn.Close()
		return
	}
	fmt.Printf("redis conn pool size: %d\n", len(redisPool))
	redisPool <- conn
}

func InitRedis(address string) redis.Conn {
	if len(redisPool) == 0 {
		redisPool = make(chan redis.Conn, MaxPoolSize)
		go func() {
			for i := 0; i < MaxPoolSize/2; i++ {
				c, err := redis.Dial("tcp", address)
				if err != nil {
					panic(err)
				}
				putRedis(c)
			}
		}()
	}
	return <-redisPool
}

func getRedisConn() redis.Conn {
	return InitRedis("127.0.0.1:6379")
}

// User related
func AddNewUser(user *model.User, password string) (errMsg model.ErrorType) {
	o := orm.NewOrm()
	account := model.AccountList{
		Email:    user.Email,
		Password: password,
		UId:      user.UId,
	}

	if o.Read(&account) == nil { // err = nil means that email exist already
		return model.ErrEmailExist
	} else {
		o.Insert(&account) // insert into account list
	}
	_, err := o.Insert(user) // insert into userinfo list
	if err != nil {
		fmt.Println(err)
		o.Delete(&account)
		return model.ErrSysDb
	}
	return model.ErrSuccess
}

func AddNewToken(accountName, uId string) (token string) {
	tokenList := model.CreateNewToken(accountName, uId)
	o := orm.NewOrm()
	o.Insert(tokenList)
	return tokenList.Token
}

func Login(account, password string) (errMsg model.ErrorType, userInfo *model.UserInfo) {

	o := orm.NewOrm()
	accountCheck := model.AccountList{Email: account}
	err := o.Read(&accountCheck)
	if err == orm.ErrNoRows { // not find this account in db
		return model.ErrEmailMiss, nil
	}

	if accountCheck.Password != password { // password not correct
		return model.ErrSignInPass, nil
	}

	userId := accountCheck.UId

	user := model.User{UId: userId}
	err = o.Read(&user)
	if err == orm.ErrNoRows { // no user info
		return model.ErrSysDb, nil
	}
	defer cacheUser(&user)

	token := AddNewToken(account, userId)

	var userResult = model.UserInfo{
		User:  user,
		Token: token,
	}

	user.LastSignInTime = model.GetNowTime()
	o.Update(&user, "LastSignInTime")

	return model.ErrSuccess, &userResult
}

func cacheUser(user *model.User) {
	conn := getRedisConn()
	defer putRedis(conn)
	//	_, err = conn.Do("hmset", user.UId, "last_signin_time", model.GetNowTime())
	//	if err != nil {
	//		fmt.Println(err)
	//	}
}

func readCacheUser() {
	//	v, err := redis.Int(conn.Do("sismember", keyUserEmail, account))
	//	if err != nil {
	//		return model.ErrSysDb, nil
	//	}
	//	if v == 0 {
	//		return model.ErrEmailMiss, nil
	//	}
	//	userId, err := redis.String(conn.Do("get", account))
	//	pasInDb, err := redis.String(conn.Do("hget", userId, keyPassword))

	//	if err != nil {
	//		return model.ErrSysDb, nil
	//	}
	//	if pasInDb != password {
	//		return model.ErrSignInPass, nil
	//	}
	//	values, err := redis.Values(conn.Do("hgetall", userId))
	//	if err != nil {
	//		return model.ErrSysDb, nil
	//	}
	//	var theUser model.User

	//	if err := redis.ScanStruct(values, &theUser); err != nil {
	//		return model.ErrSysDb, nil
	//	}
}

func cacheAccount(account *model.AccountList) {
	//	conn := getRedisConn()
	//	defer putRedis(conn)
	//	v, err := redis.Int(conn.Do("sadd", keyUserEmail, user.Email))
	//	if err != nil {
	//		return model.ErrSysDb
	//	}
	//	if v == 0 {
	//		return model.ErrEmailExist
	//	}

	//	_, err = conn.Do("set", user.Email, user.UId)

	//	_, err = conn.Do("hmset", redis.Args{user.UId}.AddFlat(user)...)
	//	_, err = conn.Do("hset", user.UId, keyPassword, password)

	//	if err != nil {
	//		fmt.Println(err)
	//	}
}

// Aritcle related
func GetArticleList() {

}

func AddNewArticle() {

}

func CheckAuthority(userId string, wantAuthority model.UserAuthority) (errors model.ErrorType, users *model.User) {
	conn := getRedisConn()
	defer putRedis(conn)

	values, err := redis.Values((conn.Do("hgetall", userId)))
	if err != nil {
		return model.ErrUidMiss, nil
	}
	var user model.User
	if err := redis.ScanStruct(values, &user); err != nil {
		return model.ErrUidMiss, nil
	}
	if user.Authority <= wantAuthority {
		return model.ErrSuccess, &user
	} else {
		return model.ErrAuLimit, nil
	}
}
