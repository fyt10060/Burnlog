# Burnlog
## Introduction
A beginner's blog server trial by golang and docker.
Learn server programming from zero

## 0.0.3

### 版本更新

通过beego orm 实现了MySQL的链接，将原本的读取和存储redis的登录和注册的API改为了MySQL的实现，还没有实现缓存

### 关于token
首先用于区分用户的肯定是UID，UID也同时是user表的主键，那么如果按照简单的方法，肯定是直接就用UID来区分用户，用UID来获取这个用户的info、article、comment、material等等，这个也会是比较快的，可是，真的应该这么做么？如果我们返回给前端UID，那就相当于把整个用户最关键的key交了出去，所以这里我打算把用户数据的架构改成这样

* user
	* UID 'pk'
	* other info
* email_list
	* email 'pk'
	* password
	* UID
* token_list
	* token 'pk'
	* UID 
	* create_time

即登录的时候去查询email_list,核对密码和获取UID，得到UID在底层生成一个新的token，存入token_list, 当然每次调用signin的API我们肯定是会新生成token的，所以可能会有很多token对应同一个UID，我们只要同时跑一个定时任务，每一段时间清理一下token_list清除我们认为过期的token就好了。当然这里会有这么一个问题，虽然我们可以每一个UID只对应一个token，这样看起来也更美好，因为一个账户只会有一个token，token_list这张表会比多个token对应同个UID小好几倍，查询的速度也会更快，可是这种时候真的好么？这里其实可以稍微聊一下后端的设计思想，在我们公司里带我的后端师傅也说过类似的，后端与前端最大的差别就是前端是每一份程序就对应一个用户，不论是网页还是APP，而后端的每一份程序对应着成千上万的用户，这些用户里可能会有各种客户端，那比如我们认为token7天就过期了，用户在手机APP端，网页端分别登录了，网页端为先，2天之后手机端登录了，说好了7天过期，结果手机客户端5天就过期无法登录了，那又该如何呢？

### 小知识点
* 给一个struct中的匿名字段赋值，赋值的是另一个struct的copy，即更改A的匿名字段的值，并不会影响赋值的B的值，所以在登录操作的时候缓存我们用defer语句来做，更改了lastSignInTime之后再缓存。
* 可以去了解下golang的struct的tag语法，可以实现一些比较棒的语法糖。