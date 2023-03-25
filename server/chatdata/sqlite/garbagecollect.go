package sqlite

import (
	"encoding/json"
	"os"
	"time"
)

func (c *SqliteChatdata) loadEventMetadata() error {
	metadataFile, err := os.ReadFile(c.metadataPath)
	if os.IsNotExist(err) {
		// Metadata file does not exist. No loading required.
		return nil
	}

	err = json.Unmarshal(metadataFile, &c.eventMetadata)
	if err != nil {
		return err
	}

	return nil
}

func (c *SqliteChatdata) saveEventMetadata() error {
	metadataBytes, err := json.Marshal(c.eventMetadata)
	if err != nil {
		return err
	}

	metadataFile, err := os.Create(c.metadataPath)
	if err != nil {
		return err
	}

	metadataFile.Write(metadataBytes)
	if err != nil {
		return err
	}

	return nil
}

func (c *SqliteChatdata) garbageCollectRoutine() {
	for {
		time.Sleep(garbageCollectSleep)
		if err := c.garbageCollect(); err != nil {
			// Failures during garbage collection are fatal
			c.log.Fatalf("%v", err)
		}

		c.log.Printf("Garbage collected up to: %v", c.eventMetadata.GarbageCollectedTo)
	}
}

func (c *SqliteChatdata) garbageCollect() error {
	// Acquire global lock, as garbage collection will mutate the event database
	// and event metadata.
	c.globalLock.Lock()
	defer c.globalLock.Unlock()

	newGarbageCollectedTo := make(map[string]int64)
	for _, pid := range c.allPids {
		garbageCollectedTo, err := c.garbageCollectPid(pid)
		if err != nil {
			return err
		}

		newGarbageCollectedTo[pid] = garbageCollectedTo
	}

	// Update the metadata and store the changes.
	c.eventMetadata.GarbageCollectedTo = newGarbageCollectedTo

	if err := c.saveEventMetadata(); err != nil {
		return err
	}

	return nil
}

func (c *SqliteChatdata) garbageCollectPid(pid string) (int64, error) {
	garbageCollectUpTo, err := c.nextExpectedSequenceNumber(pid)
	if err != nil {
		return 0, err
	}
	// Do not garbage collect the last contiguous message.
	garbageCollectUpTo--

	if err := c.garbageCollectLikeEvents(pid, garbageCollectUpTo); err != nil {
		return 0, err
	}

	return garbageCollectUpTo, nil
}

func (c *SqliteChatdata) garbageCollectLikeEvents(pid string, garbageCollectUpTo int64) error {
	// Get all of the events to check for garbage collection.
	eventsToCheck := make([]*LikeEvent, 0)
	err := c.db.Model(&LikeEvent{}).
		Joins("Event").
		Where("Event__pid = ? AND Event__sequence_number <= ?", pid, garbageCollectUpTo).
		Find(&eventsToCheck).Error
	if err != nil {
		return err
	}

	likeEventsPruned := 0
	for _, likeEvent := range eventsToCheck {
		// Check if a like event with a higher precedence exists in the database.
		highestPrecendence := LikeEvent{}
		result := causalOrdering(
			c.db.Model(&LikeEvent{}).
				Joins("Event").
				Where("Event__chatroom_id = ? AND message_id = ? AND liker_id = ?", likeEvent.Event.ChatroomID, likeEvent.MessageID, likeEvent.LikerID).
				Limit(1),
		).Find(&highestPrecendence)
		if result.Error != nil {
			return err
		}
		if result.RowsAffected <= 0 || highestPrecendence.ID == likeEvent.ID {
			continue
		}

		// There exists a like event with a higher precedence. Delete this
		// LikeEvent.
		tx := c.db.Begin()
		if err := c.db.Model(&LikeEvent{}).Delete(likeEvent).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := c.db.Model(&Event{}).Delete(likeEvent.Event).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit().Error; err != nil {
			return err
		}

		likeEventsPruned++
	}

	c.log.Printf("Garbage collection pruned %d LikeEvents", likeEventsPruned)

	return nil
}
