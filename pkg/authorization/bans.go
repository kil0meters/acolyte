package authorization

import (
	"log"
	"time"

	"github.com/kil0meters/acolyte/pkg/database"
)

// Ban a struct representing a ban
type Ban struct {
	AccountID string    `db:"account_id"`
	BanTime   time.Time `db:"ban_time"`
	UnbanTime time.Time `db:"unban_time"`
	BanReason string    `db:"ban_reason"`
	BannedBy  string    `db:"banned_by"` // account id
}

// GetBans returns an array of all bans
func GetBans() []*Ban {
	rows, err := database.DB.Queryx("SELECT account.account_id, (SELECT TOP 1 ")
	if err != nil {
		log.Println(err)
		return nil
	}

	bans := make([]*Ban, 0)

	for rows.Next() {
		ban := Ban{}
		rows.StructScan(&ban)

		bans = append(bans, &ban)
	}

	return bans
}

// CheckBansJob starts a job to check if bans run out
func CheckBansJob() {
	checkBans()
	ticker := time.NewTicker(5 * time.Minute)
	go func(ticker *time.Ticker) {
		for {
			select {
			case <-ticker.C:
				checkBans()
			}
		}
	}(ticker)
}

func checkBans() { // TODO:  this approach scales hoorribly with a large databse
	accounts := GetAccounts()

	log.Println("Checking for unabns...")

	for _, account := range accounts {
		ban, err := account.GetBanInfo()
		if err != nil {
			continue
		}

		if account.ID == ban.AccountID {
			if account.Permissions == Banned && time.Now().After(ban.UnbanTime) {
				log.Printf("User %s is no longer banned", account.Username)
				account.Unban()
			}
		}
	}
}

// Ban bans a user for time.Duration
func (account *Account) Ban(bannedBy *Account, reason string, duration time.Duration) error {
	account.Permissions = Banned

	_, err := database.DB.Exec("UPDATE accounts SET permissions = $1 WHERE account_id = $2", account.Permissions, account.ID)
	if err != nil {
		return err
	}

	_, err = database.DB.Exec("INSERT INTO bans (account_id, unban_time, ban_reason, banned_by) VALUES ($1, $2, $3, $4)", account.ID, time.Now().Add(duration), reason, bannedBy.ID)
	if err != nil {
		return err
	}

	return nil
}

// Unban bans a user for time.Duration
func (account *Account) Unban() error {
	account.Permissions = Banned

	_, err := database.DB.Exec("UPDATE accounts SET permissions = $1 WHERE account_id = $2", Standard, account.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetBanInfo returns information about the most recent ban of a user
func (account *Account) GetBanInfo() (*Ban, error) {
	rows, err := database.DB.Queryx("SELECT * FROM bans WHERE account_id = $1 ORDER BY ban_time DESC", account.ID)
	if err != nil {
		return nil, err
	}

	ban := Ban{}

	rows.Next()
	err = rows.StructScan(&ban)
	if err != nil {
		return nil, err
	}

	return &ban, nil
}
