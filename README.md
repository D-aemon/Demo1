# Demo1
一个基于vue+go-gin的jwt登录验证小Demo。

## describe(描述)
demo1关注的是基于jwt的后端实现，将不同模块进行了分离与重构。通过注册端口，用户信息和加密后的密码将被存入数据库。通过登录端口的验证后，将派发token。而获取用户信息端口在验证token后将会返回该用户的信息。
demo1对数据库、路由、中间件、返回请求、工具类、数据模型和配置文件均进行了不同程度的封装。简单修改即可开始搭建其他业务逻辑。

## USE（如何使用）

#### step1 
在拉取demo1后，需要修改位于==config/application.yml==的配置信息
```
server:
  port: 8080
datasource:
  driverName: mysql
  host: 127.0.0.1
  port: 3306
  database: your database
  username: your username
  password: your password
  charset: utf8
```
#### step 2
修改位于==common/jwt.go==下的jwtkey和claims相关信息
```
var jwtKey = []byte("a_key_by_D_aemon")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	//设置过期时间
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	//创建认证
	claims := &Claims{
		UserId:         user.ID,
		StandardClaims: jwt.StandardClaims{
			//过期时间
			ExpiresAt: expirationTime.Unix(),
			//发放时间
			IssuedAt: time.Now().Unix(),
			//发放者
			Issuer: "D_aemon",
			//主题
			Subject: "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
```

#### step3
启动！
```
go run main.go
```

## API 文档

| 方法   | url | 请求参数 | 返回参数 | 
| -------- | :----------: | :---------:|---------------|
|POST     |  /api/auth/register | name, telephone, password | code, data, msg |
|POST     |  /api/auth/login | telephone, password | code, data, msg |
|POST     |  /api/auth/info | token | code, data, msg |

### api详解
----
> **/api/auth/register** --用户注册

请求示例
```
{
  "name": "zhangsan", //name可有可无，若为空则随机生成
  "telephone": 12345678911, //telephone的长度需为11位
  "password": "123456" //password的长度需大于6位
}
```

成功示例
```
{
  "code": 200,
  "data": "",
  "msg": "注册成功"
}
```
-----
> **/api/auth/login** --> 用户登录

请求示例
```
{
  "telephone": 12345678911, //11位电话
  "password": "123456" //大于6位
}
```
成功示例
```
{
  "code": 200,
  "data": {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQsImV4cCI6MTYxODk5MzU3MywiaWF0IjoxNjE4Mzg4NzczLCJpc3MiOiJ6cHBqIiwic3ViIjoidXNlciB0b2tlbiJ9.SYCsoxkhdq1Ru-ZVb6AFic_0QQems4JuKLXxzs1Q6wM"},
  "msg": "登陆成功"
}
```

-----
> **/api/auth/info** --> 获取用户信息

请求实例 (json方式)
```
{
  "token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOjQsImV4cCI6MTYxODk5MzU3MywiaWF0IjoxNjE4Mzg4NzczLCJpc3MiOiJ6cHBqIiwic3ViIjoidXNlciB0b2tlbiJ9.SYCsoxkhdq1Ru-ZVb6AFic_0QQems4JuKLXxzs1Q6wM"
}
```
请求示例 ---- 设置header:{'Authorization': 'Bearer ' + token})
```
//以axios为例  
axios({
    method: 'post',
    url: url,
    headers: { 'Authorization': 'Bearer ' + token }
  }).then(res => {
   console.log(res.data.data.user)
  })
```
返回示例
```
{
  "code": 200,
  "data":{
	"user":{
		"name":"张三", 
		"telephone":"12345678911"
		}
	 },
  "msg": "successs"
}
```
