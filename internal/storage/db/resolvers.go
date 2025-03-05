package db

//TODO: Check maybe I only need check error in tx.
import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/nabishec/avito_shop_api/internal/model"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (r *Storage) UserNameExist(userName string) (err error) {
	const op = "internal.storage.db.UserNameExist()"

	log.Debug().Msgf("%s started", op)

	queryGetUserID := `SELECT user_id
						FROM Users
						WHERE name = $1`

	var userID uuid.UUID
	err = r.db.QueryRow(queryGetUserID, userName).Scan(&userID)

	if err != nil && err == sql.ErrNoRows {
		return ErrUserNameNotExist
	}

	return err

}

func (r *Storage) UserIDExist(userID uuid.UUID) (err error) {
	const op = "internal.storage.db.UserIDExist()"

	log.Debug().Msgf("%s started", op)

	queryGetUserName := `SELECT name
						FROM Users
						WHERE user_id = $1`

	var userName string
	err = r.db.QueryRow(queryGetUserName, userID).Scan(&userName)

	if err != nil && err == sql.ErrNoRows {
		return ErrUserIDNotExist
	}

	return err

}

func (r *Storage) GetUserID(userAuthData model.AuthRequest) (userID uuid.UUID, err error) {
	const op = "internal.storage.db.GetUserID()"

	log.Debug().Msgf("%s started", op)

	queryGetUserID := `SELECT user_id, password
						FROM Users
						WHERE name = $1`

	var passwordHash string
	err = r.db.QueryRow(queryGetUserID, userAuthData.Name).Scan(&userID, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			userID, err = r.AddUser(userAuthData)
			return
		}

		err = fmt.Errorf("%s:%w", op, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(userAuthData.Password))
	if err != nil {
		err = ErrIncorrectUserPassword
		return
	}

	log.Debug().Msgf("%s's id was found", userAuthData.Name)
	return
}

func (r *Storage) AddUser(userAuthData model.AuthRequest) (userID uuid.UUID, err error) {
	const op = "internal.storage.db.AddUser()"

	log.Debug().Msgf("%s started", op)

	passwordHash, err := createPasswordHash(userAuthData.Password)
	if err != nil {
		log.Error().Err(err).Msg("Failed create password hash")
		err = ErrIncorrectUserPassword
		return
	}

	balanceForNewUser, err := strconv.Atoi(os.Getenv("BALANCE_FOR_NEW_USER"))
	if err != nil {
		balanceForNewUser = 1000
	}
	queryAddUser := `INSERT INTO  Users (user_id, name, password)
						VALUES ($1,$2,$3)`
	queryAddBalanceToUser := `INSERT INTO  Balance (user_id, coins_number)
						VALUES ($1,$2)`

	// create transaction
	tx, err := r.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userID = uuid.New()
	_, err = tx.Exec(queryAddUser, userID, userAuthData.Name, passwordHash)
	if err != nil {
		err = fmt.Errorf("%s:%w", op, err)
		return
	}

	_, err = tx.Exec(queryAddBalanceToUser, userID, balanceForNewUser)
	if err != nil {
		err = fmt.Errorf("%s:%w", op, err)
		return
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("%s:%w", op, err)
		return
	}
	log.Debug().Msgf("user %s add successfully", userAuthData.Name)
	return

}

func createPasswordHash(password string) (string, error) {
	const op = "internal.http_server.hadnlers.auth.createPasswordHash()"

	if len(password) < 1 {
		return "", fmt.Errorf("%s:%s", op, "password is empty")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}
	return string(passwordHash), nil
}

func (r *Storage) GetItemByUser(userID uuid.UUID, item string) error {
	const op = "internal.storage.db.GetItemByUser()"

	log.Debug().Msgf("%s started", op)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryCheckItem := `SELECT item_id
							FROM Items
							WHERE type = $1`

	queryWithdrawingCoinsForBuying := `UPDATE Balance
								SET coins_number = coins_number - (SELECT price FROM Items WHERE type = $1)
								WHERE user_id = $2`

	queryAddInInventory := `INSERT INTO Inventory (user_id, item_id, quantity)  
								VALUES ($1, (SELECT item_id FROM Items WHERE type = $2), $3)
								ON CONFLICT (user_id, item_id)
								DO UPDATE SET quantity = Inventory.quantity + 1;`

	var item_id int
	err = tx.QueryRow(queryCheckItem, item).Scan(&item_id)
	if err != nil && err == sql.ErrNoRows {
		return ErrItemNotExist
	}
	_, err = tx.Exec(queryWithdrawingCoinsForBuying, &item, &userID)
	if err != nil {
		//TODO: Check it working or not
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23514" || pqErr.Constraint == "positive_coins_number" {

				err = ErrNotEnoughCoins
				return err
			}
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	var itemQuantity = 1

	_, err = tx.Exec(queryAddInInventory, &userID, &item, itemQuantity)

	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("User %s bought %s successfully", userID.String(), item)
	return nil
}

func (r *Storage) GetUserInfo(userID uuid.UUID) (userInfo *model.InfoResponse, err error) {
	const op = "internal.storage.db.GetUserInfo()"
	userInfo = &model.InfoResponse{}
	log.Debug().Msgf("%s started", op)

	tx, err := r.db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// inventory items sents receiveds
	queryGetCoinsNumber := `SELECT coins_number 
								FROM Balance
								WHERE user_id = $1`

	queryGetItemsFromInvetory := `SELECT Items.type, Inventory.quantity
									FROM Inventory 
									LEFT JOIN Items
									ON Inventory.item_id = Items.item_id
									WHERE Inventory.user_id = $1
									ORDER BY create_date DESC`

	queryGetSentByUSer := `SELECT Users.name, Sent.amount
									FROM Sent 
									LEFT JOIN Users
									ON Sent.to_user_id = Users.user_id
									WHERE Sent.user_id = $1
									ORDER BY transaction_date DESC`

	queryGetReceivedByUSer := `SELECT Users.name, Received.amount
									FROM Received 
									LEFT JOIN Users
									ON Received.from_user_id = Users.user_id
									WHERE Received.user_id = $1
									ORDER BY transaction_date DESC`

	err = tx.QueryRow(queryGetCoinsNumber, userID).Scan(&userInfo.Coins)
	if err != nil {
		return nil, fmt.Errorf("%s:%w in coins", op, err)
	}

	err = tx.Select(&userInfo.Inventory, queryGetItemsFromInvetory, userID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w in inventory", op, err)
	}

	err = tx.Select(&userInfo.CoinHistory.Sent, queryGetSentByUSer, userID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w in sent", op, err)
	}

	err = tx.Select(&userInfo.CoinHistory.Received, queryGetReceivedByUSer, userID)
	if err != nil {
		return nil, fmt.Errorf("%s:%w in received", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s:%w in received", op, err)
	}

	log.Debug().Msgf("User %s info getted successfully", userID.String())
	return
}

func (r *Storage) SendCoinsToUser(sendData model.SendCoinRequest, userID uuid.UUID) error {
	const op = "internal.storage.db.TransferCoinsToUser()"

	log.Debug().Msgf("%s started", op)

	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryWithdrawingCoins := `UPDATE Balance
								SET coins_number = coins_number - $1
								WHERE user_id = $2`

	querySendCoins := `INSERT INTO  Sent (user_id, to_user_id, amount)
						VALUES ($1,(SELECT user_id FROM Users WHERE name = $2),$3)
						RETURNING to_user_id`

	queryReceiveCoins := `INSERT INTO  Received (user_id, from_user_id, amount)
						VALUES ($1,$2,$3)`

	queryDepositingCoins := `UPDATE Balance
								SET coins_number = coins_number + $1
								WHERE user_id = $2`

	_, err = tx.Exec(queryWithdrawingCoins, sendData.Amount, userID)
	if err != nil {
		//TODO: Check it working or not
		if pqErr, ok := err.(*pgconn.PgError); ok {
			if pqErr.Code == "23514" || pqErr.ConstraintName == "positive_coins_number" {

				err = ErrNotEnoughCoins
				return err
			}
		}

		return fmt.Errorf("%s:%w", op, err)
	}

	var toUserID uuid.UUID
	err = tx.QueryRow(querySendCoins, userID, sendData.ToUser, sendData.Amount).Scan(&toUserID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = tx.Exec(queryReceiveCoins, toUserID, userID, sendData.Amount)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	_, err = tx.Exec(queryDepositingCoins, sendData.Amount, toUserID)
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	log.Debug().Msgf("Send coins from user %s to user %s has been successfully", userID.String(), toUserID.String())
	return nil
}
