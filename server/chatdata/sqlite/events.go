package sqlite

import pb "github.com/Richie78321/groupchat/chatservice"

func messageEventToEventPb(m *MessageEvent) *pb.Event {
	return &pb.Event{
		Pid:              m.Event.Pid,
		SequenceNumber:   m.Event.SequenceNumber,
		LamportTimestamp: m.Event.LamportTimestamp,
		Event: &pb.Event_MessageAppend{
			MessageAppend: &pb.MessageAppend{
				ChatroomId:  m.ChatroomID,
				MessageUuid: m.MessageID,
				AuthorId:    m.AuthorID,
				Body:        m.MessageBody,
			},
		},
	}
}
