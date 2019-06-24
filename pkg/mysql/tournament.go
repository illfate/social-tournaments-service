package mysql

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/illfate/social-tournaments-service/pkg/sts"
)

// AddTournament adds tournament with passed name and deposit. Return id of this tournament.
func (c *Connector) AddTournament(ctx context.Context, name string, deposit uint64) (int64, error) {
	insert, err := c.db.ExecContext(ctx, `
 INSERT INTO tournaments (name,deposit)
 	  VALUES (?, ?)`,
		name, deposit)
	if err != nil {
		return 0, fmt.Errorf("couldn't add tournament: %s", err)
	}
	id, err := insert.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// GetTournament returns tournament with passed id. If tournament isn't found,
// function returns ErrNotFound.
func (c *Connector) GetTournament(ctx context.Context, id int64) (*sts.Tournament, error) {
	var (
		users    string
		winner   sql.NullInt64
		finished bool
		t        sts.Tournament
	)
	err := c.db.QueryRowContext(ctx, `
    SELECT id, name, deposit, prize, winner, finished, JSON_ARRAYAGG(user_id)
      FROM tournaments
INNER JOIN participants ON id = tournament_id 
     WHERE id = ?
  GROUP BY id`, id).
		Scan(&t.ID, &t.Name, &t.Deposit, &t.Prize, &winner, &finished, &users)
	if err == sql.ErrNoRows {
		return nil, sts.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(users), &t.Users)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal json: %s", err)
	}
	if finished {
		if !winner.Valid {
			return nil, fmt.Errorf("no winner")
		}
		t.Winner = winner.Int64
	}
	return &t, nil
}
