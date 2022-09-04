package shinobi

import (
	"context"
)

type MonitorStatus string

const (
	Watching = MonitorStatus("Watching")
)

type MonitorMode string

const (
	Start = MonitorMode("start")
	Stop  = MonitorMode("stop")
)

type MID string

type Monitor struct {
	MID    MID           `json:"mid"`
	Name   string        `json:"name"`
	Status MonitorStatus `json:"status"`
	Mode   MonitorMode   `json:"mode"`
	Code   int           `json:"code,string"`
}

func (s *Server) Monitors(ctx context.Context, group Group) ([]Monitor, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultRequestTimeout)
	defer cancel()

	var monitors []Monitor

	err := s.sendGetRequest(ctx, "/monitor/"+string(group), &monitors)

	return monitors, err
}

func (s *Server) MonitorByID(ctx context.Context, group Group, mid MID) (Monitor, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultRequestTimeout)
	defer cancel()

	var monitors []Monitor

	err := s.sendGetRequest(ctx, "/monitor/"+string(group)+"/"+string(mid), &monitors)
	if err != nil {
		return Monitor{}, nil
	}

	if len(monitors) != 0 {
		return Monitor{}, ErrInconsistentResponseFormat
	}

	return monitors[0], nil
}
