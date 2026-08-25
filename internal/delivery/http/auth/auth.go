package goldgym

import (
	"errors"
	"gold-gym-be/internal/entity/auth/v2"
	"gold-gym-be/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
func (h *Handler) LoginUser(c *gin.Context) {
	// func (h *Handler) DeleteGoldGymGin(c *gin.Context) {
	log.Println("testDelivery")
	resp := response.Response{}
	// defer resp.RenderJSON(w, r)

	// ctx := c.Request.Context()
	ctx := c.Request.Context()

	device := c.GetHeader("User-Agent")

	user, password, ok := c.Request.BasicAuth()
	if !ok {
		log.Printf("[ERROR] %s %s - %s\n", c.Request.Method, c.Request.URL, errors.New("403 Forbidden"))
		return
	}
	log.Println("testDelivery2")

	result, metadata, _, err := h.goldgymSvc.LoginUser(ctx, user, password, c.Request.RemoteAddr, device)
	log.Println("testDelivery3")
	if err != nil {
		// Return error message with HTTP 200 OK
		resp.SetError(err, http.StatusOK)

		log.Printf("[ERROR] %s %s - %s\n", c.Request.Method, c.Request.URL, err.Error())
		return
	}

	resp.Data = result
	resp.Metadata = metadata

	log.Printf("[INFO] %s %s\n", c.Request.Method, c.Request.URL)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var (
		token        auth.Token
		refreshToken string
		err          error
	)
	ctx := c.Request.Context()
	refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}
	token, err = h.goldgymauthSvc.RefreshToken(ctx, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, token)
}

func (h *Handler) Logout(c *gin.Context) {
	var (
		// token        auth.Token
		refreshToken string
		err          error
	)
	ctx := c.Request.Context()
	refreshToken, err = c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}
	err = h.goldgymauthSvc.Logout(ctx, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, err)
}
