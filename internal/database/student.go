package database

import (
	"GoAssignment/internal/contextkey"
	"GoAssignment/internal/student"
	"context"
	"database/sql"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	uuid "github.com/satori/go.uuid"
)

type StudentRow struct {
	ID          string
	Fname       sql.NullString
	Lname       sql.NullString
	Email       sql.NullString
	Gender      sql.NullString
	DateOfBirth sql.NullTime
	Address     sql.NullString
	CreatedBy   sql.NullString
	CreatedOn   sql.NullTime
	UpdatedBy   sql.NullString
	UpdatedOn   sql.NullTime
}

func convertStudentRowToStudent(s StudentRow) student.Student {
	dateOfBirth := student.CustomeTime{Time: s.DateOfBirth.Time}
	return student.Student{
		ID:          s.ID,
		Fname:       s.Fname.String,
		Lname:       s.Lname.String,
		Email:       s.Email.String,
		Gender:      s.Gender.String,
		DateOfBirth: dateOfBirth,
		Address:     s.Address.String,
		CreatedBy:   s.CreatedBy.String,
		CreatedOn:   s.CreatedOn.Time,
		UpdatedBy:   s.UpdatedBy.String,
		UpdatedOn:   s.UpdatedOn.Time,
	}
}

func (d *Database) CreateStudent(ctx context.Context, std student.Student) (student.Student, error) {
	createdby, ok := ctx.Value(contextkey.UserIDKey).(string)
	log.Info(createdby)
	if !ok {
		return student.Student{}, fmt.Errorf("could not retrieve user ID from context")
	}
	std.ID = uuid.NewV4().String()
	std.CreatedBy = createdby
	dateOfBirthNullTime := sql.NullTime{Time: std.DateOfBirth.Time, Valid: !std.DateOfBirth.IsZero()}
	postRow := StudentRow{
		ID:          std.ID,
		Fname:       sql.NullString{String: std.Fname, Valid: true},
		Lname:       sql.NullString{String: std.Lname, Valid: true},
		Email:       sql.NullString{String: std.Email, Valid: true},
		Gender:      sql.NullString{String: std.Gender, Valid: true},
		DateOfBirth: dateOfBirthNullTime,
		Address:     sql.NullString{String: std.Address, Valid: true},
		CreatedBy:   sql.NullString{String: std.CreatedBy, Valid: true},
		CreatedOn:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:   sql.NullString{String: std.UpdatedBy, Valid: true},
		UpdatedOn:   sql.NullTime{Time: time.Now(), Valid: false},
	}

	rows, err := d.Client.NamedExecContext(
		ctx,
		`INSERT INTO student 
        (id, fname, lname, email, gender, dateofbirth, address, createdby, createdon, updatedby, updatedon) 
     VALUES
        (:id, :fname, :lname, :email, :gender, :dateofbirth, :address, :createdby, :createdon, :updatedby, :updatedon)`,
		postRow,
	)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert Student: %w", err)
	}
	insertedID, err := rows.LastInsertId()
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to close rows: %w", err)
	}
	fmt.Printf("New record ID: %d\n", insertedID)
	return std, nil
}

func (d *Database) GetStudent(ctx context.Context, uuid string) (student.Student, error) {
	var stdRow StudentRow
	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, fname, lname, email, gender, dateofbirth, address, createdby, createdon, updatedby, updatedon
		FROM student
		WHERE id = ?`,
		uuid,
	)
	err := row.Scan(&stdRow.ID, &stdRow.Fname, &stdRow.Lname, &stdRow.Email, &stdRow.Gender, &stdRow.DateOfBirth, &stdRow.Address, &stdRow.CreatedBy, &stdRow.CreatedOn, &stdRow.UpdatedBy, &stdRow.UpdatedOn)
	if err != nil {
		return student.Student{}, fmt.Errorf("an error occurred fetching the student by uuid: %w", err)
	}
	return convertStudentRowToStudent(stdRow), nil
}

func (d *Database) GetStudents(ctx context.Context) ([]student.Student, error) {
	var stdRows []StudentRow
	err := d.Client.Select(&stdRows, "SELECT * FROM student LIMIT 10")
	if err != nil {
		return nil, fmt.Errorf("fetchStudents %v", err)
	}
	students := make([]student.Student, len(stdRows))
	for i, stdRow := range stdRows {
		students[i] = student.Student{
			ID:          stdRow.ID,
			Fname:       stdRow.Fname.String,
			Lname:       stdRow.Lname.String,
			Email:       stdRow.Email.String,
			Gender:      stdRow.Gender.String,
			DateOfBirth: student.CustomeTime{Time: stdRow.DateOfBirth.Time},
			Address:     stdRow.Address.String,
			CreatedBy:   stdRow.CreatedBy.String,
			CreatedOn:   stdRow.CreatedOn.Time,
			UpdatedBy:   stdRow.UpdatedBy.String,
			UpdatedOn:   stdRow.UpdatedOn.Time,
		}
	}

	return students, nil
}

func (d *Database) DeleteStudent(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM student WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment from the database: %w", err)
	}
	return nil
}

func (d *Database) UpdateStudent(ctx context.Context, id string, std student.Student) (student.Student, error) {
	updatedby, ok := ctx.Value(contextkey.UserIDKey).(string)
	std.UpdatedBy = updatedby
	if !ok {
		return student.Student{}, fmt.Errorf("could not retrieve user ID from context")
	}
	dateOfBirthNullTime := sql.NullTime{Time: std.DateOfBirth.Time, Valid: !std.DateOfBirth.IsZero()}
	stdRow := StudentRow{
		ID:          id,
		Fname:       sql.NullString{String: std.Fname, Valid: true},
		Lname:       sql.NullString{String: std.Lname, Valid: true},
		Email:       sql.NullString{String: std.Email, Valid: true},
		Gender:      sql.NullString{String: std.Gender, Valid: true},
		DateOfBirth: dateOfBirthNullTime,
		Address:     sql.NullString{String: std.Address, Valid: true},
		CreatedBy:   sql.NullString{String: std.CreatedBy, Valid: true},
		CreatedOn:   sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:   sql.NullString{String: std.UpdatedBy, Valid: true},
		UpdatedOn:   sql.NullTime{Time: time.Now(), Valid: true},
	}

	rows, err := d.Client.NamedExecContext(
		ctx,
		`UPDATE student 
		SET fname = :fname,
			lname = :lname,
			email = :email,
			gender = :gender,
			dateofbirth = :dateofbirth,
			address = :address,
			createdby = :createdby,
			createdon = :createdon,
			updatedby = :updatedby,
			updatedon = :updatedon
		WHERE id = :id`,
		stdRow,
	)
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to insert Student: %w", err)
	}
	insertedID, err := rows.LastInsertId()
	if err != nil {
		return student.Student{}, fmt.Errorf("failed to close rows: %w", err)
	}
	fmt.Printf("New record ID: %d\n", insertedID)

	return convertStudentRowToStudent(stdRow), nil
}
