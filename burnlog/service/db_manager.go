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

	orm.RegisterModel(new(model.User))
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
	conn := getRedisConn()
	defer putRedis(conn)
	v, err := redis.Int(conn.Do("sadd", keyUserEmail, user.Email))
	if err != nil {
		return model.ErrSysDb
	}
	if v == 0 {
		return model.ErrEmailExist
	}
	_, err = conn.Do("set", user.Email, user.UserID)

	_, err = conn.Do("hmset", redis.Args{user.UserID}.AddFlat(user)...)
	_, err = conn.Do("hset", user.UserID, keyPassword, password)

	if err != nil {
		fmt.Println(err)
	}
	return model.ErrSuccess
}

func Login(account, password string) (errMsg model.ErrorType, user *model.User) {
	conn := getRedisConn()
	defer putRedis(conn)
	v, err := redis.Int(conn.Do("sismember", keyUserEmail, account))
	if err != nil {
		return model.ErrSysDb, nil
	}
	if v == 0 {
		return model.ErrEmailMiss, nil
	}
	userId, err := redis.String(conn.Do("get", account))
	pasInDb, err := redis.String(conn.Do("hget", userId, keyPassword))

	if err != nil {
		return model.ErrSysDb, nil
	}
	if pasInDb != password {
		return model.ErrSignInPass, nil
	}
	values, err := redis.Values(conn.Do("hgetall", userId))
	if err != nil {
		return model.ErrSysDb, nil
	}
	var theUser model.User

	if err := redis.ScanStruct(values, &theUser); err != nil {
		return model.ErrSysDb, nil
	}
	token := model.GetUserToken(userId)
	theUser.Token = token

	_, err = conn.Do("hmset", userId, "last_signin_time", model.GetNowTime(), "token", token)
	if err != nil {
		fmt.Println(err)
	}
	return model.ErrSuccess, &theUser
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
