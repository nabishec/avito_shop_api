package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/rs/zerolog/log"
)

func (r *Database) GetUserID(userAuthData model.AuthRequest) (userID uuid.UUID, err error) {
	const op = "internal.storage.db.GetUserID()"

	queryGetUserID := `SELECT user_id, password
						FROM Users
						WHERE name = $1`

	var password string
	err = r.DB.QueryRow(queryGetUserID, userAuthData.Name).Scan(&userID, &password)
	if err != nil {
		if err == pgx.ErrNoRows {
			userID, err = r.AddUser(userAuthData)
			return
		}

		err = fmt.Errorf("%s:%w", op, err)
		return
	}

	log.Debug().Msgf("%s's id was found", userAuthData.Name)
	return
}

func (r *Database) AddUser(userAuthData model.AuthRequest) (userID uuid.UUID, err error) {
	const op = "internal.storage.db.AddUser()"

	balanceForNewUser, err := strconv.Atoi(os.Getenv("BALANCE_FOR_NEW_USER"))
	if err != nil {
		balanceForNewUser = 1000
	}
	queryAddUser := `INSERT INTO  Users (user_id, name, password)
						VALUES ($1,$2,$3)`
	queryAddBalanceToUser := `INSERT INTO  Users (user_id, coins_number)
						VALUES ($1,$2)`

	//create transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	userID = uuid.New()
	_, err = tx.Exec(queryAddUser, userID, userAuthData.Name, userAuthData.Password)
	if err != nil {
		err = fmt.Errorf("%s:%w", op, err)
		return
	}

	_, err = tx.Exec(queryAddBalanceToUser, userID, balanceForNewUser)
	if err != nil {
		err = fmt.Errorf("%s:%w", op, err)
		return
	}

	log.Debug().Msgf("user %s add successfully", userAuthData.Name)
	return

}

func (r *Database) GetItemByUser(userID uuid.UUID, item string) error {
	const op = "internal.storage.db.GetItemByUser()"

	queryWithdrawingMoney := `UPDATE Balance
								SET coins_number = coins_number - (SELECT price FROM Items WHERE name = $1)
								WHERE user_id = $2`

	_, err := r.DB.Exec(queryWithdrawingMoney, &item, &userID)
	if err != nil {
		//TODO: Check it working or not
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23514" || pqErr.Constraint == "positive_coins_number" {

				err = ErrNotEnoughCoins
				return err
			}
		}

		err = fmt.Errorf("%s:%w", op, err)
		return err
	}

	return nil
}
