package shinobi

import "context"

type Key struct {
	Key  string `json:"ke"`
	UID  string `json:"uid"`
	IP   string `json:"ip"`
	Code string `json:"code"`
}

type KeysResponse struct {
	OK   bool   `json:"ok"`
	UID  string `json:"uid"`
	Keys []Key  `json:"keys"`
}

func (s *Server) Keys(ctx context.Context, group Group) ([]Key, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultRequestTimeout)
	defer cancel()

	var response KeysResponse

	err := s.sendGetRequest(ctx, "/api/"+string(group)+"/list", &response)

	return response.Keys, err
}
