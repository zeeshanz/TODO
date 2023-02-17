package tests

import (
	"database/sql"
	"database/sql/driver"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"github.com/zeeshanz/TODO/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
}

type Repo struct {
	db *gorm.DB
}

type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

func (s *RepoTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)
	db, s.mock, err = sqlmock.New()
	s.NoError(err)

	s.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	s.NoError(err)

}

func (s *RepoTestSuite) TestCreateUser() {
	repo := NewRepo(s.db)

	username_test := "harris"
	password_test := "password"
	uuid_test := "a139e8e0-aecc-11ed-afa1-0242ac120002"

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"users\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"uuid\",\"username\",\"password\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\,\\$6\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, uuid_test, username_test, password_test).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	user, err := CreateUser(uuid_test, username_test, password_test, repo.db)

	s.NoError(err)
	s.Equal(uuid_test, user.Uuid)
	s.Equal(username_test, user.Username)
	s.Equal(password_test, user.Password)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RepoTestSuite) TestFindUserEmpty() {
	repo := NewRepo(s.db)
	user_test := "harris"
	user_found := &models.User{}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(user_test).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := FindUser(user_found, user_test, repo.db)

	s.NoError(err)
	s.Equal(user_found.Username, "")
	s.NoError(s.mock.ExpectationsWereMet())

}

func (s *RepoTestSuite) TestFindExistingUser() {
	repo := NewRepo(s.db)

	username_create := "harris"
	password_create := "password"
	uuid_test := "a139e8e0-aecc-11ed-afa1-0242ac120002"

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"users\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"uuid\",\"username\",\"password\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\,\\$6\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, uuid_test, username_create, password_create).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	user_find := "harris"
	user_found := &models.User{}
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT 1`)).
		WithArgs(user_find).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user_created, err_created := CreateUser(uuid_test, username_create, password_create, repo.db)
	err_find := FindUser(user_found, user_find, repo.db)

	s.NoError(err_find)
	s.NoError(err_created)
	s.Equal(username_create, user_created.Username)
	s.Equal(password_create, user_created.Password)
	s.NoError(s.mock.ExpectationsWereMet())

}

func (s *RepoTestSuite) TestCreateTodo() {
	repo := NewRepo(s.db)

	todo_item_test := "Breakfast at 7 am"
	completed_test := false
	uuid_test := "a139e8e0-aecc-11ed-afa1-0242ac120002"
	user_uuid_test := "fcf1050e-e191-45da-894d-4d965fb10c56"

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"todos\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"uuid\",\"todo_item\",\"completed\",\"user_uuid\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\,\\$6\\,\\$7\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, uuid_test, todo_item_test, completed_test, user_uuid_test).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	todo, err := CreateTodo(uuid_test, todo_item_test, completed_test, user_uuid_test, repo.db)

	s.NoError(err)
	s.Equal(uuid_test, todo.Uuid)
	s.Equal(user_uuid_test, todo.UserUuid)
	s.Equal(todo_item_test, todo.TodoItem)
	s.Equal(completed_test, todo.Completed)
	s.NoError(s.mock.ExpectationsWereMet())
}

func (s *RepoTestSuite) TestCreateTodoWrongValue() {
	repo := NewRepo(s.db)

	todo_item_test := "Breakfast at 7 am"
	todo_item_wrong_value_test := "Breakfast at 8 am"
	completed_test := false
	uuid_test := "a139e8e0-aecc-11ed-afa1-0242ac120002"
	user_uuid_test := "fcf1050e-e191-45da-894d-4d965fb10c56"

	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO \"todos\" \\(\"created_at\",\"updated_at\",\"deleted_at\",\"uuid\",\"todo_item\",\"completed\",\"user_uuid\"\\) VALUES \\(\\$1,\\$2,\\$3,\\$4,\\$5\\,\\$6\\,\\$7\\) RETURNING \"id\"").
		WithArgs(anyTime{}, anyTime{}, nil, uuid_test, todo_item_test, completed_test, user_uuid_test).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	todo, err := CreateTodo(uuid_test, todo_item_test, completed_test, user_uuid_test, repo.db)

	s.NoError(err)
	s.Equal(uuid_test, todo.Uuid)
	s.Equal(user_uuid_test, todo.UserUuid)
	s.NotEqual(todo_item_wrong_value_test, todo.TodoItem)
	s.Equal(completed_test, todo.Completed)
	s.NoError(s.mock.ExpectationsWereMet())
}

func TestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTestSuite))
}
