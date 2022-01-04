package repository

import (
	"fmt"

	"demo/pb"
)

type EventStoreRepository struct{}

func (store EventStoreRepository) CreateEvent(event *pb.Event) error {
	// Insert two rows into the "accounts" table.
	// sql := fmt.Sprintf("INSERT INTO events (id, eventtype, aggregateid, aggregatetype, eventdata, channel)
	// VALUES ('%s','%s','%s','%s','%s','%s')", event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.EventData, event.Channel)
	sql := `
INSERT INTO events (id, eventtype, aggregateid, aggregatetype, eventdata, stream) 
VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(sql, event.EventId, event.EventType, event.AggregateId, event.AggregateType, event.EventData, event.Stream)
	if err != nil {
		return fmt.Errorf("error on insert into events: %w", err)
	}
	return nil
}

func (store EventStoreRepository) GetEvents(filter *pb.EventFilter) []*pb.Event {
	var events []*pb.Event
	return events
}
