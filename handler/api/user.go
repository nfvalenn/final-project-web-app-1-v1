package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UserAPI interface {
	Login(w http.ResponseWriter, r *http.Request)
	Register(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)
}

type userAPI struct {
	userService service.UserService
}

func NewUserAPI(userService service.UserService) *userAPI {
	return &userAPI{userService}
}

func (u *userAPI) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.UserLogin

	err := json.NewDecoder(r.Body).Decode(&user) 
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	if user.Email == "" || user.Password == ""{
		w.WriteHeader(400)
		p := entity.ErrorResponse{Error: "email or password is empty"}
		json, _ := json.Marshal(p)
		w.Write(json)
		return

	}

	id, err := u.userService.Login(r.Context(), &entity.User{
		Email: user.Email,
		Password: user.Password,
	})

	if err != nil {
		w.WriteHeader(500)
		r := entity.ErrorResponse{Error: "error internal server"}
		json, _ := json.Marshal(r)
		w.Write(json)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "user_id", 
		Value: strconv.Itoa(id),
		Path: "/",
		Expires: time.Now().Add(5 * time.Minute),
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": id,
		"message": "login success",
	})
}

func (u *userAPI) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.UserRegister

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}


	if user.Fullname == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(400)
		x := entity.ErrorResponse{Error: "register data is empty"}
		json, _ := json.Marshal(x)
		w.Write(json)
		return
	}
	y := entity.User{
		Fullname: user.Fullname,
		Email: user.Email,
		Password: user.Password,
	}
	i, err := u.userService.Register(r.Context(), &y)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.ErrorResponse{Error: "error internal server"})
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id" : i.ID,
		"message": "register success",
	})
	}


func (u *userAPI) Logout(w http.ResponseWriter, r *http.Request) {
	//userId := r.Context().Value("id")

	//api.
	http.SetCookie(w, &http.Cookie{
		Name: "user_id",
		Value: "",
		Path: "/",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "logout success",
	})


}
func (u *userAPI) Delete(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("user_id")

	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("user_id is empty"))
		return
	}

	deleteUserId, _ := strconv.Atoi(userId)

	err := u.userService.Delete(r.Context(), int(deleteUserId))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "delete success"})
}
