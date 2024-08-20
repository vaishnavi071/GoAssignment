package student

import (
	"context"
	"errors"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrFetchingStudent = errors.New("could not fetch student by ID")
	ErrUpdatingStudent = errors.New("could not update student")
	ErrNoStudentFound  = errors.New("no student found")
	ErrDeletingStudent = errors.New("could not delete student")
	ErrNotImplemented  = errors.New("not implemented")
)

type CustomeTime struct {
	time.Time
}

const dateLayout = "02-01-2006"

func (ct *CustomeTime) UnmarshalJSON(b []byte) (err error) {
	dateStr := strings.Trim(string(b), `"`)
	ct.Time, err = time.Parse(dateLayout, dateStr)
	return
}

type Student struct {
	ID          string      `json :"id"`
	Fname       string      `json :"fname"`
	Lname       string      `json :"lname"`
	DateOfBirth CustomeTime `json :"dateofbirth"`
	Email       string      `json :"email"`
	Address     string      `json :"address"`
	Gender      string      `json :"gender"`
	CreatedBy   string      `json :"createdby"`
	CreatedOn   time.Time   `json :"createdon"`
	UpdatedBy   string      `json :"updatedby"`
	UpdatedOn   time.Time   `json :"updatedon"`
}

type StudentStore interface {
	CreateStudent(context.Context, Student) (Student, error)
	GetStudent(context.Context, string) (Student, error)
	DeleteStudent(context.Context, string) error
	UpdateStudent(context.Context, string, Student) (Student, error)
	GetStudents(context.Context) ([]Student, error)
	Ping(context.Context) error
}

type Service struct {
	db StudentStore
}

func NewService(db StudentStore) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateStudent(ctx context.Context, student Student) (Student, error) {
	student, err := s.db.CreateStudent(ctx, student)
	if err != nil {
		log.Errorf("an error occurred adding the Student: %s", err.Error())
	}
	return student, nil

}

func (s *Service) GetStudent(ctx context.Context, ID string) (Student, error) {
	std, err := s.db.GetStudent(ctx, ID)
	if err != nil {
		log.Errorf("an error occured fetching the Student: %s", err.Error())
		return Student{}, ErrFetchingStudent
	}
	return std, nil
}

func (s *Service) DeleteStudent(ctx context.Context, ID string) error {
	return s.db.DeleteStudent(ctx, ID)
}

func (s *Service) UpdateStudent(ctx context.Context, ID string, newStudent Student) (Student, error) {
	std, err := s.db.UpdateStudent(ctx, ID, newStudent)
	if err != nil {
		log.Errorf("an error occurred updating the Student: %s", err.Error())
	}
	return std, nil
}

func (s *Service) GetStudents(ctx context.Context) ([]Student, error) {
	students, err := s.db.GetStudents(ctx)
	if err != nil {
		log.Errorf("an error occurred updating the Student: %s", err.Error())
	}
	return students, nil
}

func (s *Service) ReadyCheck(ctx context.Context) error {
	log.Info("Checking readiness")
	return s.db.Ping(ctx)
}
