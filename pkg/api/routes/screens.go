package routes

import (
	"fmt"
	"github.com/bthuilot/pixelate/pkg/api/requests"
	"github.com/bthuilot/pixelate/pkg/api/responses"
	"github.com/bthuilot/pixelate/pkg/display"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

// StopCurrentAgent is the endpoint stop the currently running rendering agent
func ClearScreen(d display.Display) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := d.ClearScreen(); err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, responses.InvalidResponse{
				Success: false,
				Message: fmt.Sprintf("unable to clear screen: %s", err),
			})
			return
		}
		c.JSON(http.StatusOK, responses.ValidResponse[struct{}]{
			Success: true,
		})
	}
}

// UpdateConfig will update the configuration of the currently running agent
func UpdateScreenConfig(d display.Display) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cfg map[string]string
		if c.ShouldBindJSON(&cfg) != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.InvalidResponse{
				Success: false,
				Message: "invalid config, must be  object of string -> string",
			})
			return
		}
		_, name, running := d.CurrentScreen()
		if !running {
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.InvalidResponse{
				Success: false,
				Message: "no service is running",
			})
			return
		}

		if err := d.SetScreenConfig(name, cfg); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.InvalidResponse{
				Success: false,
				Message: fmt.Sprintf("unable to update config: %s", err.Error()),
			})
			return
		}
	}
}

// SetAgent sets the current agent rendering to the display
func SetScreen(d display.Display) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request requests.SetScreenRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			logrus.Warningf("invalid set agent request: %s\n", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.InvalidResponse{
				Success: false,
				Message: "endpoint requires a screen name",
			})
			return
		}
		if err := d.SetScreen(request.Screen); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, responses.InvalidResponse{
				Success: false,
				Message: fmt.Sprintf("unable to start agent %s: %s", request.Screen, err),
			})
			return
		}
		c.JSON(http.StatusAccepted, responses.ValidResponse[struct{}]{
			Success: true,
		})
	}
}
