package chatdata

import pb "github.com/Richie78321/groupchat/chatservice"

func MessageListToPb(messages []Message) []*pb.Message {
	pbMessages := make([]*pb.Message, len(messages))
	for i, message := range messages {
		pbMessages[i] = MessageToPb(message)
	}

	return pbMessages
}

func MessageToPb(m Message) *pb.Message {
	return &pb.Message{
		Uuid:   m.Id().String(),
		Author: m.Author(),
		Body:   m.Body(),
		Likers: m.Likers(),
	}
}
