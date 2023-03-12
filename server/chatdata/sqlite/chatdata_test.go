package sqlite

import (
	"testing"

	pb "github.com/Richie78321/groupchat/chatservice"
)

func makeChatdata(t *testing.T) *SqliteChatdata {
	chatdata, err := NewSqliteChatdata(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	return chatdata
}

func TestIgnoreDuplicates(t *testing.T) {
	chatdata := makeChatdata(t)
	testEvent1 := &pb.Event{
		Pid:              "server1",
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
		Pid:              "server1",
		SequenceNumber:   0,
		LamportTimestamp: 78321,
		Event:            &pb.Event_MessageAppend{},
	}

	// First event should be successfully consumed.
	if _, err := chatdata.ConsumeEvent(testEvent1); err != nil {
		t.Fatal(err)
	}

	// Second duplicate event should not cause any failures.
	ignored, err := chatdata.ConsumeEvent(testEvent2)
	if err != nil {
		t.Fatal(err)
	}
	if !ignored {
		t.Fatal("Expected the duplicate event to be ignored")
	}
}
