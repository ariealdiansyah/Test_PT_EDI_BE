package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"test-pelindo/models"

	"log"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func UserServices(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) UserHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.GetUsers(w, r)
	case http.MethodPost:
		s.CreateUser(w, r)
	case http.MethodDelete:
		s.DeleteUser(w, r)
	case http.MethodPut:
		s.UpdateUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Method not allowed")
	}
}

func (s *UserService) GetUsers(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var users []models.User
	result := s.db.Find(&users)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error retrieving users")
		return
	}

	var modifiedUsers []map[string]interface{}
	for _, user := range users {
		modifiedUser := map[string]interface{}{
			"user_id":  user.UserID,
			"username": user.Username,
			"name":     user.Name,
			"password": user.Password,
			"status":   user.Status,
		}
		modifiedUsers = append(modifiedUsers, modifiedUser)
	}

	// Encode data ke dalam format JSON
	jsonDatas, err := json.Marshal(modifiedUsers)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonDatas)
}

func (h *UserService) GetUserDetail(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var user models.User

	result := h.db.First(&user, userID)
	if result.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user ID")
		return
	}

	modifiedUser := map[string]interface{}{
		"user_id":  user.UserID,
		"username": user.Username,
		"name":     user.Name,
		"password": user.Password,
		"status":   user.Status,
	}

	jsonData, err := json.Marshal(modifiedUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error converting users to JSON")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (s *UserService) CreateUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	var users models.User
	err := json.NewDecoder(r.Body).Decode(&users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding request body")
		return
	}
	fmt.Println("r Body", json.NewDecoder(r.Body).Decode(&users))
	fmt.Println("user", &users)
	result := s.db.Create(&users)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error creating users")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "User created successfully")
}

func (s *UserService) UpdateUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Invalid user ID")
		return
	}
	var user models.User
	result := s.db.Where("user_id = ?", userID).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "User not found")
		return
	}

	// Parse the request body
	var updatedUser models.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding request body")
		return
	}

	user.Name = updatedUser.Name
	user.Username = updatedUser.Username
	user.Password = updatedUser.Password
	user.Status = updatedUser.Status

	result = s.db.Save(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error updating user")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User updated successfully")
}

func (s *UserService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	defer r.Body.Close()
	var requestMap map[string]interface{}
	if err := json.Unmarshal(requestBody, &requestMap); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	id, ok := requestMap["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid key format", http.StatusBadRequest)
		return
	}

	userID := uint(id)

	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err := s.db.Delete(&user).Error; err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "User deleted successfully")
}
