package auth_plugin

import (
	"fmt"
	"os"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/nitric-dev/membrane/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthPlugin - The dev implementation for the Nitric Auth Plugin
type AuthPlugin struct {
	db scribble.Driver
}

// User - The local user entity representation
type User struct {
	id             string `json:id`
	email          string `json:email`
	pwdHashAndSalt string `json:pwdHashAndSalt`
}

// CreateUser - Create a new user using scribble as the DB
func (s *AuthPlugin) CreateUser(tenant string, id string, email string, password string) error {
	collection := fmt.Sprint("auth_%s", tenant)
	var tmpUser User
	err := s.db.Read(collection, id, &tmpUser)
	if os.IsNotExist(err) {
		// We can create the user
		bHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		s.db.Write(collection, id, &User{
			id:             id,
			email:          email,
			pwdHashAndSalt: bHash,
		})
		return nil
	} else if err == nil {
		return fmt.Errorf("User %s already exists")
	}

	return err
}

// New - Instansiate a New concrete dev auth plugin
func New() {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", "/nitric/")
	db, err := scribble.New(dbDir, nil)

	if err != nil {
		return nil, err
	}

	return &AuthPlugin{
		db: db,
	}, nil
}