package transport

import (
	"GoAssignment/internal/authentication"
	"GoAssignment/internal/student"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type StudentService interface {
	CreateStudent(ctx context.Context, std student.Student) (student.Student, error)
	GetStudent(ctx context.Context, ID string) (student.Student, error)
	DeleteStudent(ctx context.Context, ID string) error
	UpdateStudent(ctx context.Context, ID string, newStd student.Student) (student.Student, error)
	GetStudents(ctx context.Context) ([]student.Student, error)
	ReadyCheck(ctx context.Context) error
}

type PostStudentRequest struct {
	ID          string              `json :"id"`
	Fname       string              `json :"fname"`
	Lname       string              `json :"lname"`
	DateOfBirth student.CustomeTime `json :"dateofbirth"`
	Email       string              `json :"email"`
	Address     string              `json :"address"`
	Gender      string              `json :"gender"`
	CreatedBy   string              `json :"createdby"`
	CreatedOn   time.Time           `json :"createdon"`
	UpdatedBy   string              `json :"updatedby"`
	UpdatedOn   time.Time           `json :"updatedon"`
}

func studetFromPostStudentRequest(u PostStudentRequest) student.Student {
	return student.Student{
		Fname:       u.Fname,
		Lname:       u.Lname,
		Email:       u.Email,
		DateOfBirth: u.DateOfBirth,
		Address:     u.Address,
		Gender:      u.Gender,
		CreatedOn:   u.CreatedOn,
		UpdatedBy:   u.UpdatedBy,
		UpdatedOn:   u.UpdatedOn,
	}
}
func (h *Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	var postStdReq PostStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&postStdReq); err != nil {
		return
	}
	validate := validator.New()
	err := validate.Struct(postStdReq)
	if err != nil {
		log.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	std := studetFromPostStudentRequest(postStdReq)
	std, err = h.Service.CreateStudent(r.Context(), std)
	if err != nil {
		log.Error("Failed to create student: ", err)
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(std); err != nil {
		log.Error("Failed to encode response: ", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	std, err := h.Service.GetStudent(r.Context(), id)
	if err != nil {
		if errors.Is(err, student.ErrFetchingStudent) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(std); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]

	if commentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteStudent(r.Context(), commentID)
	if err != nil {
		return
	}

	if err := json.NewEncoder(w).Encode(Response{Message: "Successfully Deleted"}); err != nil {
		panic(err)
	}
}

func (h *Handler) GetStudents(w http.ResponseWriter, r *http.Request) {

	students, err := h.Service.GetStudents(r.Context())
	if err != nil {
		if errors.Is(err, student.ErrFetchingStudent) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(students); err != nil {
		http.Error(w, "Failed to encode students", http.StatusInternalServerError)
	}
}

type UpdateStudentRequest struct {
	Fname       string              `json :"fname"`
	Lname       string              `json :"lname"`
	DateOfBirth student.CustomeTime `json :"dateofbirth"`
	Email       string              `json :"email"`
	Address     string              `json :"address"`
	Gender      string              `json :"gender"`
	CreatedBy   string              `json :"createdby"`
	CreatedOn   time.Time           `json :"createdon"`
	UpdatedBy   string              `json :"updatedby"`
	UpdatedOn   time.Time           `json :"updatedon"`
}

func studentFromUpdateStudentRequest(u UpdateStudentRequest) student.Student {
	return student.Student{
		Fname:       u.Fname,
		Lname:       u.Lname,
		Email:       u.Email,
		DateOfBirth: u.DateOfBirth,
		Address:     u.Address,
		Gender:      u.Gender,
		CreatedOn:   u.CreatedOn,
		UpdatedOn:   u.UpdatedOn,
	}
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	studentID := vars["id"]

	var updateStdRequest UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStdRequest); err != nil {
		return
	}

	validate := validator.New()
	err := validate.Struct(updateStdRequest)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	std := studentFromUpdateStudentRequest(updateStdRequest)

	std, err = h.Service.UpdateStudent(r.Context(), studentID, std)
	if err != nil {
		log.Error(err.Error())
		return
	}
	if err := json.NewEncoder(w).Encode(std); err != nil {
		panic(err)
	}
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID")
	password := r.Header.Get("Password")

	if userID == "" || password == "" {
		http.Error(w, "Missing credentials", http.StatusBadRequest)
		return
	}

	token, err := authentication.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
