package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"simplebank/token"

	"github.com/gin-gonic/gin"
)

const (
	authorisationHeaderKey  = "authorization"
	authorisationTypeBearer = "bearer"
	authorisationPayloadKey = "authorization_payload"
)

// AuthMiddleware creates a gin middleware for authorisation
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorisationHeader := ctx.GetHeader(authorisationHeaderKey)

		if len(authorisationHeader) == 0 {
			err := errors.New("authorisation header is not provided")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		fields := strings.Fields(authorisationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorisation header format")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		authorisationType := strings.ToLower(fields[0])
		if authorisationType != authorisationTypeBearer {
			err := fmt.Errorf("unsupported authorisation type %s", authorisationType)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		ctx.Set(authorisationPayloadKey, payload)
		ctx.Next()
	}
}
