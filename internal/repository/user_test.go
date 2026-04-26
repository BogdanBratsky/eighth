package repository

import (
	"context"
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/BogdanBratsky/eigth/internal/config"
	"github.com/BogdanBratsky/eigth/internal/db"
	"github.com/BogdanBratsky/eigth/internal/model"
	"github.com/joho/godotenv"
)

var (
	testRepo *UserPostgres
	testDB   *sql.DB
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println(err)
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	testDB, err = db.Connect(&cfg.DBConfig)
	if err != nil {
		log.Fatal(err)
	}
	if err := testDB.Ping(); err != nil {
		log.Fatal(err)
	}

	testRepo = NewUserPostgres(testDB)

	code := m.Run()
	testDB.Close()
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	user := &model.User{
		Login:        "test2",
		Email:        "test2",
		PasswordHash: "test",
	}

	// t.Cleanup(func() {
	// 	testDB.ExecContext(ctx, "DELETE FROM users WHERE login = $1", user.Login)
	// })

	err := testRepo.Create(ctx, user)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("user id is %v", user.ID)

	if user.ID == 0 {
		t.Fatal("user did not create")
	}
}

// func TestGetByEmail(t *testing.T) {
// 	ctx := context.Background()

// 	user, err := testRepo.GetByIdentifier(ctx, "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Println(user)
// }

// func TestGetByID(t *testing.T) {
// 	ctx := context.Background()

// 	user, err := testRepo.GetByID(ctx, 1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Println(user)
// }

// func TestUserList(t *testing.T) {
// 	ctx := context.Background()

// 	users, err := testRepo.List(ctx, 10, 0)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	log.Printf("%+v\n", users)
// }
