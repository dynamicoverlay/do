package resolvers

import (
	"context"
	mocket "github.com/Selvatico/go-mocket"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"imabad.dev/do/api/handlers"
	"testing"
)

func SetupTests() *gorm.DB { // or *gorm.DB
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	// GORM
	db, _ := gorm.Open(mocket.DriverName, "connection_string") // Can be any connection string
	return db
}

var rootResolver RootResolver

func TestUserMutations(t *testing.T) {
	rootResolver = RootResolver{Db: SetupTests()}
	TestUserSignupMutations(t)
	TestRootResolver_Login(t)
	TestRootResolver_UpdateUser(t)
}

func TestUserSignupMutations(t *testing.T) {
	t.Run("Creating a user with the same email address fails", func(t *testing.T) {
		var args = createUserArgs{
			Username: "Rushmead",
			Email:    "stuart@pomeroys.site",
			Password: "test123",
		}
		commonReply := []map[string]interface{}{{"count": 1}}
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT count(*) FROM "users"  WHERE "users"."deleted_at" IS NULL AND`).WithReply(commonReply).OneTime()
		response, err := rootResolver.CreateUser(args)
		if err != nil {
			if err.Error() != "user already exists with email address" {
				t.Errorf("did not error correctly")
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
	t.Run("Creating a user with no password fails", func(t *testing.T) {
		var args = createUserArgs{
			Username: "Rushmead",
			Email:    "stuart@pomeroys.site",
			Password: "",
		}
		response, err := rootResolver.CreateUser(args)
		if err != nil {
			if err.Error() != "missing password" {
				t.Errorf("did not error correctly: %v", err)
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
	t.Run("Creating a user with no email fails", func(t *testing.T) {
		var args = createUserArgs{
			Username: "Rushmead",
			Email:    "",
			Password: "testtest123",
		}
		response, err := rootResolver.CreateUser(args)
		if err != nil {
			if err.Error() != "missing email address" {
				t.Errorf("did not error correctly: %v", err)
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
	t.Run("Creating a user with an invalid email fails", func(t *testing.T) {
		var args = createUserArgs{
			Username: "Rushmead",
			Email:    "stuartpomeroys",
			Password: "testtest123",
		}
		response, err := rootResolver.CreateUser(args)
		if err != nil {
			if err.Error() != "invalid email address" {
				t.Errorf("did not error correctly: %v", err)
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
	t.Run("Creating a user with no username fails", func(t *testing.T) {
		var args = createUserArgs{
			Username: "",
			Email:    "stuart@pomeroys.site",
			Password: "testtest123",
		}
		response, err := rootResolver.CreateUser(args)
		if err != nil {
			if err.Error() != "missing username" {
				t.Errorf("did not error correctly: %v", err)
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
	t.Run("Creating a user", func(t *testing.T) {
		var args = createUserArgs{
			Username: "Rushmead",
			Email:    "stuart@pomeroys.site",
			Password: "test123",
		}
		response, err := rootResolver.CreateUser(args)
		if err != nil {
			t.Error(err)
		} else {
			if len(response.token) <= 0 {
				t.Errorf("did not generate token")
			}
		}
	})
}
func TestRootResolver_Login(t *testing.T) {
	email := "stuart@pomeroys.site"
	password := "test123"
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		t.Error(err)
	}
	hashedPassword := string(bytes)
	commonReply := []map[string]interface{}{{"username": "test-user", "email": email, "password": hashedPassword, "role": "user"}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND`).WithReply(commonReply)
	t.Run("Logging in with the correct details succeeds", func(t *testing.T) {
		response, err := rootResolver.Login(loginUserArgs{
			Email:    email,
			Password: password,
		})
		if err != nil {
			t.Error(err)
		} else {
			if len(response.token) <= 0 {
				t.Errorf("did not generate token")
			}
		}
	})
	t.Run("Logging in with the incorrect password errors", func(t *testing.T) {
		response, err := rootResolver.Login(loginUserArgs{
			Email:    email,
			Password: "12315",
		})
		if err != nil {
			if err.Error() != "incorrect details" {
				t.Error(err)
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
	t.Run("Logging in with no email address fails", func(t *testing.T) {
		response, err := rootResolver.Login(loginUserArgs{
			Email:    "",
			Password: "12315",
		})
		if err != nil {
			if err.Error() != "incorrect details" {
				t.Error(err)
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
	t.Run("Logging in with no password fails", func(t *testing.T) {
		response, err := rootResolver.Login(loginUserArgs{
			Email:    email,
			Password: "",
		})
		if err != nil {
			if err.Error() != "incorrect details" {
				t.Error(err)
			}
		} else {
			if len(response.token) >= 0 {
				t.Errorf("generated token")
			}
		}
	})
}
func TestRootResolver_UpdateUser(t *testing.T) {
	username := "test-user"
	email := "stuart@pomeroys.site"
	password := "test123"
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	userID := 69
	if err != nil {
		t.Error(err)
	}
	c := context.Background()
	hashedPassword := string(bytes)
	commonReply := []map[string]interface{}{{"id": userID, "username": username, "email": email, "password": hashedPassword, "role": "user"}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND`).WithReply(commonReply)
	t.Run("Updating user details without being logged in fails", func(t *testing.T) {
		newUsername := "new-username"
		response, err := rootResolver.UpdateUser(c, updateUserArgs{
			Username: &newUsername,
		})
		if err != nil {
			if err.Error() != "unauthorized" {
				t.Error(err)
			}
		} else {
			if response.Username() == newUsername {
				t.Errorf("updated username")
			}
		}
	})
	t.Run("Updating user details", func(t *testing.T) {
		newUsername := "new-username"
		GlobalMock := mocket.Catcher.Reset()
		mockDb := GlobalMock.NewMock()
		mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 69)) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(commonReply)
		mockDb.WithQuery(`UPDATE "users" SET "updated_at" = ?, "deleted_at" = ?, "username" = ?, "email" = ?, "password" = ?, "role" = ?  WHERE "users"."deleted_at" IS NULL AND "users"."id" = ?`)
		response, err := rootResolver.UpdateUser(context.WithValue(c, handlers.ContextKey("UserID"), &userID), updateUserArgs{
			Username: &newUsername,
		})
		if !mockDb.Triggered {
			t.Errorf("did not execute update query")
		}
		if err != nil {
			t.Error(err)
		} else {
			if response.Username() != newUsername {
				t.Errorf("username did not update")
			}
		}
	})
	t.Run("Updating a users email to already taken email address fails", func(t *testing.T) {
		newEmail := "test@pomeroys.site"
		GlobalMock := mocket.Catcher.Reset()
		mockDb := GlobalMock.NewMock()
		countReply := []map[string]interface{}{{"count": 1}}
		mocket.Catcher.NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 69)) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(commonReply)
		mocket.Catcher.NewMock().WithQuery(`SELECT count(*) FROM "users"  WHERE "users"."deleted_at" IS NULL AND ((email = test@pomeroys.site))`).WithReply(countReply)
		mockDb.WithQuery(`UPDATE "users" SET "updated_at" = ?, "deleted_at" = ?, "username" = ?, "email" = ?, "password" = ?, "role" = ?  WHERE "users"."deleted_at" IS NULL AND "users"."id" = ?`)
		response, err := rootResolver.UpdateUser(context.WithValue(c, handlers.ContextKey("UserID"), &userID), updateUserArgs{
			Email: &newEmail,
		})
		if mockDb.Triggered {
			t.Errorf("executed update query")
		}
		if err == nil {
			if response.Email() == email {
				t.Errorf("email updated")
			}
		}
	})
	t.Run("Updating a users email to an invalid email address fails", func(f *testing.T) {
		newEmail := "testpomeroyssite"
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 69)) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(commonReply)
		response, err := rootResolver.UpdateUser(context.WithValue(c, handlers.ContextKey("UserID"), &userID), updateUserArgs{
			Email: &newEmail,
		})
		if err == nil {
			if response.Email() == email {
				f.Errorf("email updated")
			}
		}
	})
	t.Run("Updating password with no previous password fails", func(t *testing.T) {
		password := "newPassword!"
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 69)) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(commonReply)
		_, err := rootResolver.UpdateUser(context.WithValue(c, handlers.ContextKey("UserID"), &userID), updateUserArgs{
			Password: &password,
		})
		if err != nil {
			if err.Error() != "updating your password requires your current password" {
				t.Errorf("incorrect error")
			}
		} else {
			t.Errorf("did not error")
		}
	})
	t.Run("Updating password where previous password is incorrect should fail", func(t *testing.T) {
		password := "newPassword!"
		previousPassword := "blah"
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 69)) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(commonReply)
		_, err := rootResolver.UpdateUser(context.WithValue(c, handlers.ContextKey("UserID"), &userID), updateUserArgs{
			Password:         &password,
			PreviousPassword: &previousPassword,
		})
		if err != nil {
			if err.Error() != "incorrect password" {
				t.Errorf("incorrect error")
			}
		} else {
			t.Errorf("did not error")
		}
	})
	t.Run("Updating password where previous password is correct shouldn't fail", func(t *testing.T) {
		password := "newPassword!"
		previousPassword := "test123"
		mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "users"  WHERE "users"."deleted_at" IS NULL AND (("users"."id" = 69)) ORDER BY "users"."id" ASC LIMIT 1`).WithReply(commonReply)
		response, err := rootResolver.UpdateUser(context.WithValue(c, handlers.ContextKey("UserID"), &userID), updateUserArgs{
			Password:         &password,
			PreviousPassword: &previousPassword,
		})
		if err != nil {
			t.Error(err)
		} else {
			err := bcrypt.CompareHashAndPassword([]byte(response.U.Password), []byte(password))
			if err != nil {
				t.Errorf("did not update password")
			}
		}
	})
}
