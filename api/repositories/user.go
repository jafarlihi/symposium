package repositories

import (
	"github.com/jafarlihi/symposium/backend/database"
	"github.com/jafarlihi/symposium/backend/logger"
	"github.com/jafarlihi/symposium/backend/models"
)

func GetUser(id uint32) (*models.User, error) {
	sql := "SELECT id, username, email, password, access FROM users where id = $1"
	row := database.Database.QueryRow(sql, id)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Access)
	if err != nil {
		logger.Log.Error("Failed to SELECT a user, error: " + err.Error())
		return nil, err
	}
	return &user, nil
}

func GetUserByUsername(username string) (*models.User, error) {
	sql := "SELECT id, username, email, password, access FROM users where username = $1"
	row := database.Database.QueryRow(sql, username)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Access)
	if err != nil {
		logger.Log.Error("Failed to SELECT a user, error: " + err.Error())
		return nil, err
	}
	return &user, nil
}

func CreateUser(username string, email string, password string) (int64, error) {
	sql := "INSERT INTO users (username, email, password, access) VALUES ($1, $2, $3, 0) RETURNING id"
	var id int64
	err := database.Database.QueryRow(sql, username, email, password).Scan(&id)
	if err != nil {
		logger.Log.Error("Failed to INSERT a new user, error: " + err.Error())
		return 0, err
	}
	return id, nil
}
