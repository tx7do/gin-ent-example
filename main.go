// @title Gin Ent Example
// @version 1.0
// @description Gin Ent Example

// @contact.name tx7do
// @contact.url https://tx7do.github.io/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// schemes http

package main

import (
	"context"
	"log"
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/lib/pq"
	//_ "github.com/mattn/go-sqlite3"

	_ "gin-ent-example/docs"
	"gin-ent-example/ent"
	"gin-ent-example/ent/user"
)

type Server struct {
	db   *ent.Client
	http *gin.Engine
}

var svr Server

func initDatabase() {
	// init PostgreSQL
	//client, err := ent.Open("postgres", "host=localhost port=5432 user=root password=123456 dbname=test")
	// init MySQL
	client, err := ent.Open("mysql", "root:123456@tcp(localhost:3306)/test?parseTime=True")
	// init SQLite
	//client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	// init Gremlin (AWS Neptune)
	//client, err := ent.Open("gremlin", "http://localhost:8182")
	if err != nil {
		log.Fatal(err)
		return
	}
	svr.db = client

	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func ResponseJSON(c *gin.Context, httpCode, errCode int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code: errCode,
		Msg:  msg,
		Data: data,
	})
	return
}

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, 400
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, 500
	}
	if !check {
		return http.StatusBadRequest, 400
	}

	return http.StatusOK, 200
}

// @Summary create user
// @Produce application/json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Param nickname formData string true "nickname"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /user/create [post]
func handleCreateUser(c *gin.Context) {
	type PostParam struct {
		UserName string `form:"username" json:"username" valid:"Required; MaxSize(50)"`
		Password string `form:"password" json:"password" valid:"Required; MaxSize(50)"`
		Nickname string `form:"nickname" json:"nickname" valid:"Required; MaxSize(50)"`
	}
	var form PostParam

	httpCode, errCode := BindAndValid(c, &form)
	if errCode != 200 {
		ResponseJSON(c, httpCode, errCode, "invalid param", nil)
		return
	}

	usr, err := svr.db.User.
		Create().
		// SetID(0).
		SetUsername(form.UserName).
		SetPassword(form.Password).
		SetNickname(form.Nickname).
		Save(context.Background())
	if err != nil {
		ResponseJSON(c, http.StatusOK, 500, "create user failed: "+err.Error(), nil)
		return
	}

	type ResponseData struct {
		UserID   uint64 `json:"userid"`
		UserName string `json:"username"`
		Nickname string `json:"nickname"`
	}
	var resp ResponseData
	resp.Nickname = form.Nickname
	resp.UserName = form.UserName
	resp.UserID = uint64(usr.ID)

	ResponseJSON(c, http.StatusOK, 200, "", resp)
}

// @Summary get user
// @Produce application/json
// @Param username path string true "username"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /user/{username} [get]
func handleGetUser(c *gin.Context) {
	userName := c.Param("username")
	if len(userName) == 0 {
		ResponseJSON(c, 200, 400, "invalid param", nil)
		return
	}

	usr, _ := svr.db.User.
		Query().
		Where(user.Username(userName)).
		First(context.Background())
	if usr == nil {
		ResponseJSON(c, http.StatusOK, 500, "user doesn't exist", nil)
		return
	}

	type ResponseData struct {
		UserID   uint64 `json:"userid"`
		UserName string `json:"username"`
		Nickname string `json:"nickname"`
	}
	var resp ResponseData
	resp.Nickname = usr.Nickname
	resp.UserName = usr.Username
	resp.UserID = uint64(usr.ID)

	ResponseJSON(c, http.StatusOK, 200, "", resp)
}

// @Summary update user
// @Produce application/json
// @Param username formData string true "username"
// @Param password formData string true "password"
// @Param nickname formData string true "nickname"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /user/update [post]
func handleUpdateUser(c *gin.Context) {
	type PostParam struct {
		UserName string `form:"username" json:"username" valid:"Required; MaxSize(50)"`
		Password string `form:"password" json:"password" valid:"Required; MaxSize(50)"`
		Nickname string `form:"nickname" json:"nickname" valid:"Required; MaxSize(50)"`
	}
	var form PostParam

	httpCode, errCode := BindAndValid(c, &form)
	if errCode != 200 {
		ResponseJSON(c, httpCode, errCode, "invalid param", nil)
		return
	}

	count, _ := svr.db.User.
		Update().
		SetUsername(form.UserName).
		SetPassword(form.Password).
		SetNickname(form.Nickname).
		Where(user.Username(form.UserName)).
		Save(context.Background())

	if count == 0 {
		ResponseJSON(c, http.StatusOK, 500, "update user failed", nil)
		return
	}

	usr, _ := svr.db.User.
		Query().
		Where(user.Username(form.UserName)).
		First(context.Background())
	if usr == nil {
		ResponseJSON(c, http.StatusOK, 500, "user doesn't exist", nil)
		return
	}

	type ResponseData struct {
		UserID   uint64 `json:"userid"`
		UserName string `json:"username"`
		Nickname string `json:"nickname"`
	}
	var resp ResponseData
	resp.Nickname = form.Nickname
	resp.UserName = form.UserName
	resp.UserID = uint64(usr.ID)

	ResponseJSON(c, http.StatusOK, 200, "", resp)
}

// @Summary delete user
// @Produce application/json
// @Param username path string true "username"
// @Success 200 {string} json "{"code":200,"data":{},"msg":"ok"}"
// @Router /user/{username} [delete]
func handleDeleteUser(c *gin.Context) {
	userName := c.Param("username")
	if len(userName) == 0 {
		ResponseJSON(c, 200, 400, "invalid param", nil)
		return
	}

	_, err := svr.db.User.
		Delete().
		Where(user.Username(userName)).
		Exec(context.Background())
	if err != nil {
		ResponseJSON(c, 200, 500, "delete user failed", nil)
		return
	}

	ResponseJSON(c, 200, 200, "delete user ok", nil)
}

func runHttpServer() {
	r := gin.New()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	svr.http = r

	// api doc http://localhost:8080/swagger/index.html
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create
	r.POST("/user/create", handleCreateUser)
	// Read
	r.GET("/user/:username", handleGetUser)
	// Update
	r.POST("/user/update", handleUpdateUser)
	// Delete
	r.DELETE("/user/:username", handleDeleteUser)

	// Listen and serve on 0.0.0.0:8080
	_ = r.Run(":8080")
}

func main() {
	initDatabase()
	runHttpServer()
}
