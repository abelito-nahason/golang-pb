package handler

import (
	"log"
	"net/http"
	"os"
	"pocketbase-seal/helper"
	"pocketbase-seal/model/auth"
	"pocketbase-seal/model/genericresponse"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

type authHandler struct {
	db *pocketbase.PocketBase
}

func NewAuthHandler(db *pocketbase.PocketBase) *authHandler {
	return &authHandler{
		db,
	}
}

func (h *authHandler) RegisterNewUser(e *core.RequestEvent) error {

	user := auth.RegisterInput{}
	getUser := auth.User{}

	if err := e.BindBody(&user); err != nil {
		errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "Failed reading request body", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	empty := helper.CheckEmptyStruct(user)

	if empty {
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Body contains empty field", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	if user.Role != "keeper" && user.Role != "visitor" {
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Invalid role", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	h.db.DB().NewQuery("SELECT email FROM actors WHERE email = {:email}").Bind(dbx.Params{
		"email": user.Email,
	}).One(&getUser)

	if getUser.Email == user.Email {
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Email already exists", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	hash, err := helper.HashPassword(user.Password)
	if err != nil {
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Error hashing password", ResponseCode: 500}
		return e.JSON(http.StatusInternalServerError, errResponse)
	}

	_, err = h.db.DB().NewQuery("INSERT INTO actors(email, password, role) VALUES({:email}, {:password}, {:role})").Bind(dbx.Params{
		"email":    user.Email,
		"password": hash,
		"role":     user.Role,
	}).Execute()

	if err != nil {
		log.Print(err)
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Failed creating new actor", ResponseCode: 500}
		return e.JSON(http.StatusInternalServerError, errResponse)
	}

	response := genericresponse.GenericResponse{ResponseMessage: "Success", ResponseCode: 200}
	return e.JSON(http.StatusOK, response)

}

func (h *authHandler) LoginUser(e *core.RequestEvent) error {
	login := auth.LoginInput{}
	getUser := auth.User{}

	if err := e.BindBody(&login); err != nil {
		errResponse := genericresponse.GenericErrorResponse{Error: err, ResponseMessage: "Failed reading request body", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	empty := helper.CheckEmptyStruct(login)

	if empty {
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Email/password is empty", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	h.db.DB().NewQuery("SELECT email, password, role FROM actors WHERE email = {:email}").Bind(dbx.Params{
		"email": login.Email,
	}).One(&getUser)

	correct := helper.VerifyPassword(login.Password, getUser.Password)

	if !correct {
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Incorrect password", ResponseCode: 400}
		return e.JSON(http.StatusBadRequest, errResponse)
	}

	secretKey := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"email": getUser.Email,
			"role":  getUser.Role,
			"iat":   time.Now().Unix(),
			"exp":   time.Now().Add(time.Hour).Unix(),
		})
	signedToken, err := token.SignedString([]byte(secretKey))

	if err != nil {
		log.Print(err)
		errResponse := genericresponse.GenericErrorResponse{ResponseMessage: "Internal Server Error", ResponseCode: 500}
		return e.JSON(http.StatusInternalServerError, errResponse)
	}

	response := genericresponse.GenericResponse{Data: signedToken, ResponseMessage: "Success", ResponseCode: 200}
	return e.JSON(http.StatusOK, response)
}
