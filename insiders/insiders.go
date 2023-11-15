package insiders

import (
	"context"
	"time"

	"encore.app/hello"
	"encore.dev/pubsub"
	"encore.dev/rlog"
	"encore.dev/storage/sqldb"
)

var peopleDB = sqldb.NewDatabase("people", sqldb.DatabaseConfig{
	Migrations: "./migrations",
})

type RecordResponse struct {
	FirstMet     time.Time
	MeetingCount int
}

//encore:api private path=/record/:name
func Record(ctx context.Context, name string) (*RecordResponse, error) {
	var resp RecordResponse
	err := peopleDB.QueryRow(ctx, `
		INSERT INTO people (name, first_met, meeting_count)
		VALUES ($1, NOW(), 1)
		ON CONFLICT (name) DO UPDATE
		SET meeting_count = people.meeting_count + 1
		RETURNING first_met, meeting_count
	`, name).Scan(&resp.FirstMet, &resp.MeetingCount)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

var _ = pubsub.NewSubscription(hello.Greetings, "listener", pubsub.SubscriptionConfig[*hello.GreetingEvent]{
	Handler: func(ctx context.Context, event *hello.GreetingEvent) error {
		resp, err := Record(ctx, event.Name)
		if err != nil {
			return err
		}
		rlog.Info("successfully recorded greeting", "first_met", resp.FirstMet, "meeting_count", resp.MeetingCount)
		return nil
	},
})
