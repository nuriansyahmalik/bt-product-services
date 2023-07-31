package users

import (
	"database/sql"
	"fmt"

	"github.com/evermos/boilerplate-go/infras"
	"github.com/evermos/boilerplate-go/shared/failure"
	"github.com/evermos/boilerplate-go/shared/logger"
	"github.com/gofrs/uuid"
)

var (
	userQueries = struct {
		selectUser string
		insertUser string
		updateUser string
		deleteUser string
	}{
		selectUser: `
			SELECT
			    u.userId,
			    u.username,
			    u.email,
			    u.userType,
			    u.created,
				u.created_by,
				u.updated,
				u.updated_by,
				u.deleted,
				u.deleted_by
			FROM user u`,
		insertUser: `INSERT INTO user (userId,username, email, userType, createdAt, createdBy)
		VALUE (:userId,:username, :email, :userType, :createdAt, :createdBy)`,

		updateUser: `
			UPDATE user
            SET 
                username = :username,
                userType = :userType, 
                updatedAt = :updatedAt, 
                updatedBy = :updatedBy
            WHERE userId = :userId`,
	}
)

type UserRepository interface {
	Create(user User) (err error)
	ExistsByID(id uuid.UUID) (exists bool, err error)
	ResolveByID(id uuid.UUID) (user User, err error)
	Update(user User) (err error)
}

type UserRepositoryMysql struct {
	DB *infras.MySQLConn
}

func ProvideUserRepositoryMySQL(db *infras.MySQLConn) *UserRepositoryMysql {
	return &UserRepositoryMysql{DB: db}
}

func (u *UserRepositoryMysql) Create(user User) (err error) {
	exists, err := u.ExistsByID(user.ID)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}

	if exists {
		err = failure.Conflict("create", "user", "already exists")
		logger.ErrorWithStack(err)
		return
	}
	stmt, err := u.DB.Write.PrepareNamed(userQueries.insertUser)
	if err != nil {
		logger.ErrorWithStack(err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(user)
	if err != nil {
		logger.ErrorWithStack(err)
	}
	return
}

func (u *UserRepositoryMysql) ResolveByID(id uuid.UUID) (user User, err error) {
	row := u.DB.Read.QueryRow(userQueries.selectUser+" WHERE userId = ?", id.String())
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.UserType,
		&user.CreatedAt,
		&user.CreatedBy,
		&user.UpdatedAt,
		&user.UpdatedBy,
		&user.Deleted,
		&user.DeletedBy,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user not found: %w", err)
		}
		return User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return user, nil

}
func (r *UserRepositoryMysql) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Read.Get(
		&exists,
		"SELECT COUNT(userId) FROM user WHERE user.userId = ?",
		id.String())
	if err != nil {
		logger.ErrorWithStack(err)
	}

	return
}
func (u *UserRepositoryMysql) Update(user User) (err error) {
	return err
}
