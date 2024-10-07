package repositories

import (
	"context"
	"fmt"
	"time"

	models "github.com/Orken1119/Ozinshe/internal/models/auth_models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) models.UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) CreateUser(c context.Context, user models.UserRequest) (int, error) {
	var userID int
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	userQuery := `INSERT INTO users(
		email, password, phone_number, roleid, created_at)
		VALUES ($1, $2, $3, $4, $5) returning id;`

	err := ur.db.QueryRow(c, userQuery, user.Email, user.Password.Password, user.PhoneNumber, 2, currentTime).Scan(&userID)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (ur *UserRepository) GetUserByEmail(c context.Context, email string) (models.User, error) {
	user := models.User{}

	query := `SELECT id, email, password, phone_number, roleid, created_at FROM users where email = $1`
	row := ur.db.QueryRow(c, query, email)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.PhoneNumber, &user.RoleID, &user.CreatedAt)

	if err != nil {
		return user, err
	}
	return user, nil

}

func (ur *UserRepository) GetUserByID(c context.Context, userID int) (models.User, error) {
	user := models.User{}

	query := `SELECT id, email, password, phone_number, roleid, created_at FROM users where id = $1`
	row := ur.db.QueryRow(c, query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.Password, &user.PhoneNumber, &user.RoleID, &user.CreatedAt)

	if err != nil {
		return user, err
	}
	return user, nil

}

func (ur *UserRepository) GetProfile(c context.Context, userID int) (models.UserProfile, error) {
	user := models.UserProfile{}

	query := `SELECT id,user_name, email, phone_number, birthday FROM users WHERE id = $1`
	row := ur.db.QueryRow(c, query, userID)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PhoneNumber, &user.Birthday)
	if err != nil {
		return user, err
	}

	return user, nil
}
func (ur *UserRepository) CreatePasswordResetCode(c context.Context, email string, code string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	query := `INSERT INTO passwordresetcode(email, code, createdAt)
	VALUES ($1, $2, $3);`

	_, err := ur.db.Exec(c, query, email, code, currentTime)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) GetCodeByEmail(c context.Context, email string) (string, error) {
	var code string

	query := `SELECT code FROM PasswordResetCode where email = $1`
	row := ur.db.QueryRow(c, query, email)
	err := row.Scan(&code)

	if err != nil {
		return "", err
	}
	return code, nil

}

func (ur *UserRepository) ChangeForgottenPassword(c context.Context, code string, email string, newPassword string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	var storedCode string

	codeQuery := `SELECT code FROM PasswordResetCode WHERE email = $1`

	err := ur.db.QueryRow(c, codeQuery, email).Scan(&storedCode)

	if err != nil {
		return fmt.Errorf("failed to retrieve reset code: %v", err)
	}

	if storedCode != code {
		return fmt.Errorf("invalid reset code")
	}

	updateQuery := `UPDATE users SET password = $1, created_at = $2 WHERE email = $3;`
	_, err = ur.db.Exec(c, updateQuery, newPassword, currentTime, email)
	if err != nil {
		return err
	}

	return nil
}

func (ur *UserRepository) EditPersonalData(c context.Context, userID int, name string, email string, phoneNumber string, birthday string) error {
	query := `UPDATE users
	SET 
		user_name = $1,
		email = $2,
		phone_number = $3,
		birthday = $4
	WHERE 
	    id = $5;
	`
	_, err := ur.db.Exec(c, query, name, email, phoneNumber, birthday, userID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) ChangePassword(c context.Context, userID int, password string) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	query := `UPDATE users
	SET
	password = $1,
	created_at = $2
	where
	id = $3`
	_, err := ur.db.Exec(c, query, password, currentTime, userID)
	if err != nil {
		return err
	}

	return nil
}
