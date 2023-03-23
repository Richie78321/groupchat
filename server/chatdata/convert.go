package chatdata

import pb "github.com/Richie78321/groupchat/chatservice"

func MessageListToPb(messages []Message) ([]*pb.Message, error) {
	pbMessages := make([]*pb.Message, len(messages))
	for i, message := range messages {
		m, err := MessageToPb(message)
		if err != nil {
			return nil, err
		}

		pbMessages[i] = m
	}

	return pbMessages, nil
}

func MessageToPb(m Message) (*pb.Message, error) {
	likers, err := m.Likers()
	if err != nil {
		return nil, err
	}

	return &pb.Message{
		Uuid:   m.Id().String(),
		Author: m.Author(),
		Body:   m.Body(),
		Likers: likers,
	}, nil
}
