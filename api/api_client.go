package api

import "errors"

//go:generate counterfeiter . SpotifyConnector
type SpotifyConnector interface {
	GetFollowedArtists(accessToken string) ([]Artist, error)
}

type Artist struct {
	Name string
	Id   string
}

type SpotifyApiClient struct {
	urlPrefix string
}

func NewSpotifyApiClient(apiUrlPrefix string) *SpotifyApiClient {
	return &SpotifyApiClient{
		urlPrefix: apiUrlPrefix,
	}
}

func (self *SpotifyApiClient) GetFollowedArtists(accessToken string) ([]Artist, error) {
	return nil, errors.New("not implemented")
}
