package peer

import (
	"log"
	"time"
)

//go:generate meshRPC expose -P peer -y

type Service interface {
	Hop(hop *HopState) (*HopState, error)
}

func NewService(name string, otherPeers map[string]Service) Service {
	return &peerService{
		name:       name,
		otherPeers: otherPeers,
	}
}

type peerService struct {
	name       string
	otherPeers map[string]Service
}

func (s *peerService) Hop(hop *HopState) (*HopState, error) {
	ts := time.Now()

	hop.Current = ts.UnixNano()
	hop.Last = s.name
	hop.QueueIdx = (hop.QueueIdx + 1) % len(hop.Queue)
	hop.HopCount++
	hop.HopLatency = (hop.Current - hop.Start) / hop.HopCount
	if hop.HopCount == hop.Limit {
		return hop, nil
	} else if hop.HopCount%1000 == 0 {
		log.Println("DBG: %#v", *hop)
	}
	nextSvc := s.otherPeers[hop.Queue[hop.QueueIdx]]
	lastState, err := nextSvc.Hop(hop)
	if err != nil {
		log.Println("[WARN] failed with", err, "at", hop.HopCount, "after", time.Since(ts), "total time",
			time.Since(time.Unix(0, hop.Start)))
		return hop, err
	}

	return lastState, nil
}

type HopState struct {
	Start      int64    // `json:"start"`
	Current    int64    // `json:"current"`
	Last       string   // `json:"last"`
	Queue      []string // `json:"queue"`
	QueueIdx   int      // `json:"pointer"`
	Limit      int64    // `json:"limit"`
	HopCount   int64    // `json:"count"`
	HopLatency int64    // `json:"latency"`
}
