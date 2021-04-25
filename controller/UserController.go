package controller

import (
	"Demo/common"
	"Demo/dto"
	"Demo/model"
	"Demo/response"
	"Demo/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

//注册
func Register(c *gin.Context) {
	DB := common.GetDb()
	//1、获取参数
	json := make(map[string]string)
	c.BindJSON(&json)
	name := json["name"]
	telephone := json["telephone"]
	password := json["password"]
	//2、数据验证
	if len(telephone) != 11 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "密码少于6位")
		return
	}
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)
	//3、判断手机号是否存在
	if isTelephoneExit(DB, telephone) {
		response.Response(c, http.StatusUnprocessableEntity, 422, nil, "用户已注册")
		return
	}


	//4、创建用户
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Response(c, http.StatusInternalServerError, 500, nil, "手机号必须为11位")
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)
	//5、返回结果
	response.Success(c, nil, "注册成功")
}

//登录
func Login(ctx *gin.Context)  {
	DB := common.GetDb()
	//获取参数
	json := make(map[string]string)
	_ = ctx.BindJSON(&json)
	telephone := json["telephone"]
	password := json["password"]
	log.Println(telephone, password)
	//数据验证
	if len(telephone) != 11 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "手机号必须为11位")
		return
	}
	if len(password) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码少于6位")
		return
	}

	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在")
		return
	}

	//判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil{
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		return
	}

	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Fail(ctx, nil, "系统异常")
		log.Printf("token generate error : %v", err)
		return
	}


	//返回结果
	response.Success(ctx, gin.H{"token":token}, "登陆成功")

}

//获取用户信息
func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user":dto.ToUserDto(user.(model.User))}, "success")
}

//判断电话是否存在
func isTelephoneExit(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}

	return false
}

