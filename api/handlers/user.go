package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/holys/initials-avatar"
	"github.com/jafarlihi/symposium/api/config"
	"github.com/jafarlihi/symposium/api/logger"
	"github.com/jafarlihi/symposium/api/repositories"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idString := params["id"]
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Provided ID can't be parsed as an integer"}`)
		return
	}
	user, err := repositories.GetUser(uint32(id))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to get the user"}`)
		return
	}
	user.Password = ""
	jsonResult, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to marshal the result to JSON"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(jsonResult))
}

type accountCreationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Add validation
	var acr accountCreationRequest
	err := json.NewDecoder(r.Body).Decode(&acr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Request body couldn't be parsed as JSON"}`)
		return
	}
	if acr.Username == "" || acr.Email == "" || acr.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Username, email, or password field(s) is/are missing"}`)
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(acr.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to hash the password"}`)
		return
	}
	id, err := repositories.CreateUser(acr.Username, acr.Email, string(passwordHash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to create the user"}`)
		return
	}
	avatarConfig := avatar.Config{
		FontFile: "./public/fonts/monofonto.ttf",
		FontSize: 80,
	}
	a := avatar.NewWithConfig(avatarConfig)
	avatarBytes, _ := a.DrawToBytes(acr.Username, 128)
	err = ioutil.WriteFile("./public/avatars/"+strconv.FormatInt(id, 10), avatarBytes, 0644)
	if err != nil {
		logger.Log.Error("Failed to create the user avatar, error: " + err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("avatar")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Failed to open the sent avatar file"}`)
		return
	}
	defer file.Close()
	token := r.FormValue("token")
	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Token field is missing"}`)
		return
	}
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.Jwt.SigningSecret), nil
	})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Failed to parse the token"}`)
		return
	}
	var userID float64
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		userID = claims["userID"].(float64)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Invalid token"}`)
		return
	}
	f, err := os.OpenFile("./public/avatars/"+strconv.Itoa(int(userID)), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to create an avatar file"}`)
		return
	}
	defer f.Close()
	buff := make([]byte, 512)
	_, err = file.Read(buff)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Failed to read the avatar file"}`)
		return
	}
	file.Seek(0, io.SeekStart)
	filetype := http.DetectContentType(buff)
	switch filetype {
	case "image/jpeg", "image/jpg":
		io.Copy(f, file)
	case "image/png":
		io.Copy(f, file)
	default:
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": "Uploaded avatar file is not JPEG or PNG"}`)
		return
	}
	w.WriteHeader(http.StatusOK)
}
