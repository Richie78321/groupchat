package sqlite

import (
	"testing"

	pb "github.com/Richie78321/groupchat/chatservice"
	"github.com/stretchr/testify/assert"
)

func makeChatdata(t *testing.T, pid string) *SqliteChatdata {
	chatdata, err := NewSqliteChatdata(":memory:", pid)
	assert.NoError(t, err)

	return chatdata
}

func TestIgnoreDuplicates(t *testing.T) {
	pid := "server1"
	chatdata := makeChatdata(t, pid)
	testEvent1 := &pb.Event{
		Pid:              pid,
		SequenceNumber:   0,
		LamportTimestamp: 0,
		Event: &pb.Event_MessageAppend{
			MessageAppend: &pb.MessageAppend{
				ChatroomId:  "chatroom",
				MessageUuid: "messageid",
				AuthorId:    "authorid",
				Body:        "message",
			},
		},
	}
	testEvent2 := &pb.Event{
		Pid:              pid,
		SequenceNumber:   0,
		LamportTimestamp: 78321,
		Event:            &pb.Event_MessageAppend{},
	}

	// First event should be successfully consumed.
	_, err := chatdata.ConsumeEvent(testEvent1)
	assert.NoError(t, err)

	// Second duplicate event should not cause any failures.
	ignored, err := chatdata.ConsumeEvent(testEvent2)
	assert.NoError(t, err)
	assert.True(t, ignored, "Expected the duplicate event to be ignored")
}

func TestLTSUpdated(t *testing.T) {
	chatdata := makeChatdata(t, "server1")
	lamportTimestamp := int64(100)

	_, err := chatdata.ConsumeEvent(&pb.Event{
		Pid:              "somepid",
		SequenceNumber:   0,
		LamportTimestamp: lamportTimestamp,
		Event: &pb.Event_MessageAppend{
			MessageAppend: &pb.MessageAppend{
				ChatroomId:  "chatroom",
				MessageUuid: "messageid",
				AuthorId:    "authorid",
				Body:        "message",
			},
		},
	})
	assert.NoError(t, err)

	assert.Equal(t, lamportTimestamp+1, chatdata.nextLamportTimestamp)
}

func TestLoadFromDisk(t *testing.T) {
	pid := "server1"
	chatdata := makeChatdata(t, pid)
	sequenceNumber := int64(100)
	lamportTimestamp := int64(200)
	// Add test events to the database
	err := chatdata.db.Create([]*Event{
		{
			Pid:              pid,
			SequenceNumber:   sequenceNumber,
			LamportTimestamp: lamportTimestamp,
		},
		{
			Pid:              "notserver1",
			SequenceNumber:   3000,
			LamportTimestamp: 4000,
		},
	}).Error
	assert.NoError(t, err)
	// Reset the sequence number and LTS
	chatdata.nextSequenceNumber = 0
	chatdata.nextLamportTimestamp = 0

	chatdata.loadFromDisk()

	assert.Equal(t, sequenceNumber+1, chatdata.nextSequenceNumber)
}
