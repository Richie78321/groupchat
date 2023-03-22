package sqlite

import (
	"errors"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/Richie78321/groupchat/server/chatdata"
	"gorm.io/gorm"
)

func (c *SqliteChatdata) IncomingEvents() chan<- *pb.Event {
	return c.incomingEvents
}

func (c *SqliteChatdata) OutgoingEvents() <-chan *pb.Event {
	return c.outgoingEvents
}

func (c *SqliteChatdata) SequenceNumberVector() (chatdata.SequenceNumberVector, error) {
	vector := make(chatdata.SequenceNumberVector)
	for _, pid := range c.allPids {
		nextExpected, err := c.nextExpectedSequenceNumber(pid)
		if err != nil {
			return nil, err
		}

		vector[pid] = nextExpected
	}

	return vector, nil
}

func (c *SqliteChatdata) nextExpectedSequenceNumber(pid string) (int64, error) {
	var result struct {
		NextExpected int64 `gorm:"column:next_expected"`
	}

	// Returns the smallest missing sequence number for a specific PID.
	// This query could likely be made more efficient, but it works for now.
	err := c.db.Raw(`
		SELECT MIN(sequence_number + 1) AS next_expected
		FROM events t1
		WHERE pid = ?
		  AND NOT EXISTS (
			SELECT *
			FROM events t2
			WHERE t2.sequence_number = t1.sequence_number + 1 AND t2.pid = ?
		)
	`, pid, pid).Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// There are no events with this PID. Therefore the next-expected
			// sequence number is zero.
			return 0, nil
		}

		return 0, err
	}

	return result.NextExpected, nil
}

func (c *SqliteChatdata) EventDiff(vector chatdata.SequenceNumberVector) ([]*pb.Event, error) {
	if vector == nil {
		vector = make(chatdata.SequenceNumberVector)
	}
	// Fill missing PIDs with a next-expected sequence number of 0
	for _, pid := range c.allPids {
		if _, ok := vector[pid]; ok {
			continue
		}
		vector[pid] = 0
	}

	// Get the missing events for every PID
	eventDiff := make([]*pb.Event, 0)
	for pid, sequence_number := range vector {
		events, err := c.sequenceNumberDiff(pid, sequence_number)
		if err != nil {
			return nil, err
		}

		eventDiff = append(eventDiff, events...)
	}

	return eventDiff, nil
}

// sequenceNumberDiff returns events from the specified PID that occurred on
// or after the specified sequence number.
func (c *SqliteChatdata) sequenceNumberDiff(pid string, sequence_number int64) ([]*pb.Event, error) {
	eventDiff := make([]*pb.Event, 0)

	// Query for MessageAppend events
	messageEvents := make([]MessageEvent, 0)
	err := c.db.Model(&MessageEvent{}).Joins("Event").Where("Event__pid = ? AND Event__sequence_number >= ?", pid, sequence_number).Find(&messageEvents).Error
	if err != nil {
		return nil, err
	}
	for _, messageEvent := range messageEvents {
		eventDiff = append(eventDiff, messageEventToPb(&messageEvent))
	}

	// Query for MessageLike events
	likeEvents := make([]LikeEvent, 0)
	err = c.db.Model(&LikeEvent{}).Joins("Event").Where("Event__pid = ? AND Event__sequence_number >= ?", pid, sequence_number).Find(&likeEvents).Error
	if err != nil {
		return nil, err
	}
	for _, likeEvent := range likeEvents {
		eventDiff = append(eventDiff, likeEventToPb(&likeEvent))
	}

	return eventDiff, nil
}
