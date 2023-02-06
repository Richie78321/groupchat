package chatdata

import pb "github.com/Richie78321/groupchat/chatservice"

func MessageToPb(m Message) *pb.Message {
	return &pb.Message{
		Uuid:   m.Id().String(),
		Author: m.Author(),
		Body:   m.Body(),
		// TODO(richie): Update when liking is implemented
		Likers: []*pb.User{},
	}
}
