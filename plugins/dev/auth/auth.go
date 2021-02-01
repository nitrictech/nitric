package auth_plugin

import (
	"encoding/json"
	"fmt"
	"os"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/nitric-dev/membrane/plugins/dev/ifaces"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
	"golang.org/x/crypto/bcrypt"
)

// AuthPlugin - The dev implementation for the Nitric Auth Plugin
type AuthPlugin struct {
	db ifaces.ScribbleIface
}

// User - The local user entity representation
type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	PwdHashAndSalt string `json:"pwdHashAndSalt"`
}

// CreateUser - Create a new user using scribble as the DB
func (s *AuthPlugin) CreateUser(tenant string, id string, email string, password string) error {
	collection := fmt.Sprintf("auth_%s", tenant)
	// tmpUsers := make([]User, 0)
	if usersStrs, err := s.db.ReadAll(collection); err == nil {
		var tmpUsr User
		for _, usrStr := range usersStrs {
			if err := json.Unmarshal([]byte(usrStr), &tmpUsr); err == nil {
				if tmpUsr.ID == id {
					return fmt.Errorf("User with id %s, already exists", email)
				}

				if tmpUsr.Email == email {
					return fmt.Errorf("User with email %s, already exists", email)
				}
			} else {
				return err
			}
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// We can create the user
	bHash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	s.db.Write(collection, id, &User{
		ID:             id,
		Email:          email,
		PwdHashAndSalt: string(bHash),
	})

	return nil
}

// New - Instansiate a New concrete dev auth plugin
func New() (sdk.AuthPlugin, error) {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", "/nitric/")
	db, err := scribble.New(dbDir, nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &AuthPlugin{
		db: db,
	}, nil
}

func NewWithDriver(driver ifaces.ScribbleIface) sdk.AuthPlugin {
	return &AuthPlugin{
		db: driver,
	}
}
