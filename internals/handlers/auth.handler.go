// handler --> sama seperti controller, sebagai logika untuk memproses request dan response

package handlers

import (
	"golang_server_bookstore/internals/models"
	"golang_server_bookstore/internals/repositories"
	"log"
	"net/http"

	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	*repositories.AuthRepo
}

func InitAuthHandler(item *repositories.AuthRepo) *AuthHandler {
	return &AuthHandler{item}
}

func (item *AuthHandler) Register(ctx *gin.Context) {

	// ambil body,konversi dari json atau form ke struct
	body := models.AuthModel{}
	if err := ctx.ShouldBind(&body); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		}) // kirim response dalam bentuk json, gin.H untuk membuat map dengan key string & vlaue any
		return
	}

	// check email already used or no
	result, err := item.FindByEmail(body)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if len(result) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Email already used",
		})
		return
	}

	// hashing password
	hash, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// buat data yang baru dengan password yg sudah dihash
	if err := item.SaveUser(models.AuthModel{
		// Id: body.Id,
		Email:    body.Email,
		Password: hash,
	}); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	 // kirim response dalam bentuk json, gin.H untuk membuat map dengan key string & vlaue any
	ctx.JSON(http.StatusCreated, gin.H{
		"messages": "success register",
	})

}

func (item *AuthHandler) Login(ctx *gin.Context) {

}