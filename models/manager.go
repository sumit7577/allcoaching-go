package models

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"golang.org/x/crypto/bcrypt"
)

type UserManager struct {
	User *User
}

// sanitizeName sanitizes the name by converting it to lowercase and removing extra spaces.
func sanitizeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// generateRandomUsername generates a random username by appending random numbers to the sanitized name.
func generateRandomUsername(name string) string {
	sanitizedName := sanitizeName(name)
	rand.Seed(time.Now().UnixNano())
	randomNumber := rand.Intn(10000)
	username := fmt.Sprintf("%s_%d", sanitizedName, randomNumber)

	return username
}

// CreateUser creates and saves a new user with hashed password.
func (um *UserManager) CreateUser() (*User, error) {
	o := orm.NewOrm()
	err := o.DoTx(func(ctx context.Context, txOrm orm.TxOrmer) error {
		hashedPassword, err := hashPassword(um.User.Password)

		if err != nil {
			return err
		}
		username := generateRandomUsername(um.User.Name)
		um.User.Password = hashedPassword
		um.User.Username = username
		id, err := txOrm.Insert(um.User)
		um.User.Id = uint64(id)

		if err != nil {
			return err // This will trigger a rollback
		}
		return nil
		/*authToken, nil := generateAuthToken()
		authEntry := &AuthToken{
			Key:     authToken,
			User:    um.User,
			Created: time.Now(),
		}
		_, err = txOrm.Insert(authEntry)
		fmt.Printf("authEntry: %v\n", authEntry)
		fmt.Printf("authEntry: %v\n", um.User)

		if err != nil {
			return err
		}*/

	})

	if err != nil {
		return nil, err
	}

	return um.User, nil
}

// hashPassword hashes a plain-text password using bcrypt.
func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func generateAuthToken() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	return token, nil
}

// AuthenticateUser verifies the username and password for a user.
func (um *UserManager) AuthenticateUser(username, password string) (*User, error) {
	o := orm.NewOrm()
	user := &User{}
	err := o.QueryTable(user).Filter("Username", username).One(user)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}
