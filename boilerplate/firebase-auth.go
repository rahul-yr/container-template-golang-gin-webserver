package boilerplate

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type FirebaseAuth struct {
	ServiceAccount     string
	firebaseApp        *firebase.App
	FirebaseAppOnce    sync.Once
	firebaseClient     *auth.Client
	FirebaseClientOnce sync.Once
}

func (fb *FirebaseAuth) Validate() error {
	if fb.ServiceAccount == "" {
		return errors.New("Missing service account for firebase")
	}
	return nil
}

// FirebaseApp returns a singleton instance of firebase app
func (fb *FirebaseAuth) GetSetupFirebaseApp() *firebase.App {
	fb.FirebaseAppOnce.Do(func() {
		if fb.ServiceAccount == "" {
			log.Fatalf("Missing service account for firebase")
		}
		opt := option.WithCredentialsFile(fb.ServiceAccount)
		app, err := firebase.NewApp(context.Background(), nil, opt)
		if err != nil {
			log.Fatalf("error initializing app: %v", err)
		}
		fb.firebaseApp = app
	})
	return fb.firebaseApp
}

// FirebaseClient returns a singleton instance of firebase client
func (fb *FirebaseAuth) GetSetupFirebaseClient() *auth.Client {
	fb.FirebaseClientOnce.Do(func() {
		app := fb.GetSetupFirebaseApp()
		client, err := app.Auth(context.Background())
		if err != nil {
			log.Fatalf("error getting Auth client: %v\n", err)
		}
		fb.firebaseClient = client
	})
	return fb.firebaseClient
}

func (fb *FirebaseAuth) IsAuthTokenValid(token string) bool {
	// split bearer token
	splitToken := strings.Split(token, "Bearer ")
	token = splitToken[1]
	client := fb.GetSetupFirebaseClient()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := client.VerifyIDToken(ctx, token)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return false
	}
	return true
}

func (fb *FirebaseAuth) FirebaseAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get token from request header
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.AbortWithStatus(401)
			return
		}
		if !fb.IsAuthTokenValid(token) {
			c.AbortWithStatus(401)
			return
		}
		c.Next()
	}
}
