package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/auth"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/models"
	"github.com/vivk-PBL-5-Backend/AuthServer/pkg/parser"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"
)

const (
	STATUS_OK    = "ok"
	STATUS_ERROR = "error"
)

type response struct {
	Status string `json:"status"`
	Msg    string `json:"message,omitempty"`
}

func newResponse(status, msg string) *response {
	return &response{
		Status: status,
		Msg:    msg,
	}
}

type handler struct {
	useCase auth.UseCase
}

func newHandler(useCase auth.UseCase) *handler {
	return &handler{
		useCase: useCase,
	}
}

func (h *handler) signUp(c *gin.Context) {
	inp := new(models.User)

	c.Header("Access-Control-Allow-Origin", "*")

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	if err := h.useCase.SignUp(c.Request.Context(), inp); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, newResponse(STATUS_OK, "confirm user email")) //"user created successfully"))
}

func (h *handler) confirm(c *gin.Context) {
	token := c.Param("token")
	if err := h.useCase.Confirm(c.Request.Context(), token); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, newResponse(STATUS_OK, "user created successfully"))
}

type signInResponse struct {
	*response
	Token string `json:"token,omitempty"`
}

func newSignInResponse(status, msg, token string) *signInResponse {
	return &signInResponse{
		&response{
			Status: status,
			Msg:    msg,
		},
		token,
	}
}

func (h *handler) signIn(c *gin.Context) {
	inp := new(models.User)

	c.Header("Access-Control-Allow-Origin", "*")

	if err := c.BindJSON(inp); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := h.useCase.SignIn(c.Request.Context(), inp)
	if err != nil {
		if err == auth.ErrInvalidAccessToken {
			c.AbortWithStatusJSON(http.StatusBadRequest, newSignInResponse(STATUS_ERROR, err.Error(), ""))
			return
		}

		if err == auth.ErrUserDoesNotExist {
			c.AbortWithStatusJSON(http.StatusBadRequest, newSignInResponse(STATUS_ERROR, err.Error(), ""))
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, newSignInResponse(STATUS_ERROR, err.Error(), ""))
		return
	}

	c.JSON(http.StatusOK, newSignInResponse(STATUS_OK, "", token))
}

func (h *handler) send(c *gin.Context) {
	message := new(models.Message)

	c.Header("Access-Control-Allow-Origin", "*")

	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]

	userID, err := parser.ParseToken(token, []byte(viper.GetString("auth.signing_key")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	if err = c.BindJSON(message); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	message.ID = primitive.NewObjectID()
	message.AuthorID = userID
	message.Date = time.Now()
	if err = h.useCase.Send(c.Request.Context(), message); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, newResponse(STATUS_OK, "user created successfully"))
}

func (h *handler) get(c *gin.Context) {
	message := new(models.Message)

	c.Header("Access-Control-Allow-Origin", "*")

	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]

	userID, err := parser.ParseToken(token, []byte(viper.GetString("auth.signing_key")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	if err = c.BindJSON(message); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	messages, err := h.useCase.Get(c.Request.Context(), userID, message.AuthorID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *handler) addCompanion(c *gin.Context) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]

	userID, err := parser.ParseToken(token, []byte(viper.GetString("auth.signing_key")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	companionID := c.Param("companion")

	err = h.useCase.AddCompanion(c.Request.Context(), userID, companionID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, newResponse(STATUS_OK, "added companion [\""+companionID+"\"]"))
}

func (h *handler) removeCompanion(c *gin.Context) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]

	userID, err := parser.ParseToken(token, []byte(viper.GetString("auth.signing_key")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	companionID := c.Param("companion")

	err = h.useCase.RemoveCompanion(c.Request.Context(), userID, companionID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, newResponse(STATUS_OK, "removed companion [\""+companionID+"\"]"))
}

func (h *handler) getCompanions(c *gin.Context) {
	reqToken := c.Request.Header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	token := splitToken[1]

	userID, err := parser.ParseToken(token, []byte(viper.GetString("auth.signing_key")))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	companions, err := h.useCase.GetCompanions(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, newResponse(STATUS_ERROR, err.Error()))
		return
	}

	c.JSON(http.StatusOK, companions)
}
