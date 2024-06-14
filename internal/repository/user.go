package repository
import (
	"context"
	"time"
	
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/Orken1119/Ozinshe/internal/models"
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
	
    err := ur.db.QueryRow(c, userQuery, user.Email, user.Password, user.PhoneNumber, 2, currentTime).Scan(&userID)
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

func (ur *UserRepository) GetProfile(c context.Context, userID int) (models.User, error) {
	user := models.User{}

	query := `SELECT id, email, phone_number, roleid, created_at FROM users WHERE id = $1`
	row := ur.db.QueryRow(c, query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.PhoneNumber, &user.RoleID, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}