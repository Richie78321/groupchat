package chatdata

import pb "github.com/Richie78321/groupchat/chatservice"

type SequenceNumberVector map[string]int64

type EventSynchronizer interface {
	// IncomingEvents that should be consumed and persisted.
	IncomingEvents() chan<- *pb.Event
	// OutgoingEvents that should be broadcasted.
	OutgoingEvents() <-chan *pb.Event

	// SequenceNumberVector is a map of process PIDs to the next-expected
	// event sequence number from that process.
	SequenceNumberVector() SequenceNumberVector
	// EventDiff returns the events that have not yet been received by the given
	// sequence number vector.
	EventDiff(SequenceNumberVector) []*pb.Event
}
