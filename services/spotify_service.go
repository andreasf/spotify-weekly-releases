package services

import (
	"fmt"
	"github.com/andreasf/spotify-release-feed/api"
	"github.com/andreasf/spotify-release-feed/model"
)

type SpotifyService interface {
	GetRecentReleases(accessToken string) ([]model.Album, error)
	CreatePlaylist(accessToken string, name string, tracks []model.Track) error
}

type SpotifyServiceImpl struct {
	apiClient api.SpotifyConnector
}

const ALBUMS_PER_REQUEST int = 20

func NewSpotifyService(apiClient api.SpotifyConnector) *SpotifyServiceImpl {
	return &SpotifyServiceImpl{
		apiClient: apiClient,
	}
}

func (self *SpotifyServiceImpl) GetRecentReleases(accessToken string) ([]model.Album, error) {
	profile, err := self.apiClient.GetUserProfile(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetRecentReleases: error retrieving user profile: %v", err)
	}

	artists, err := self.apiClient.GetFollowedArtists(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetRecentReleases: error retrieving followed artists: %v", err)
	}

	albums := []model.Album{}
	for _, artist := range artists {
		artistAlbums, err := self.apiClient.GetArtistAlbums(accessToken, artist.Id, profile.Country)
		if err != nil {
			return nil, fmt.Errorf("GetRecentReleases: error retrieving artistAlbums for %s: %v", artist.Id, err)
		}

		albums = append(albums, artistAlbums...)
	}

	albumDetails := make([]model.Album, 0, len(albums))

	numberOfRequests := len(albums) / ALBUMS_PER_REQUEST
	if len(albums)%ALBUMS_PER_REQUEST > 0 {
		numberOfRequests++
	}

	for i := 0; i < numberOfRequests; i++ {
		from := i * ALBUMS_PER_REQUEST
		to := min((i+1)*ALBUMS_PER_REQUEST, len(albums))
		albumSlice := albums[from:to]
		albumIds := getAlbumIds(albumSlice)

		albumInfos, err := self.apiClient.GetAlbumInfo(accessToken, albumIds)
		if err != nil {
			return nil, fmt.Errorf("GetRecentReleases: error retrieving album infos for %s: %v", albumIds, err)
		}

		albumDetails = append(albumDetails, albumInfos...)
	}

	return albumDetails, nil
}

func (self *SpotifyServiceImpl) CreatePlaylist(accessToken string, name string, tracks []model.Track) error {
	userProfile, err := self.apiClient.GetUserProfile(accessToken)
	if err != nil {
		return fmt.Errorf("CreatePlaylist: error retrieving user profile: %v", err)
	}

	playlistId, err := self.apiClient.CreatePlaylist(accessToken, userProfile.Id, name)
	if err != nil {
		return fmt.Errorf("CreatePlaylist: error creating playlist: %v", err)
	}

	err = self.apiClient.AddTracksToPlaylist(accessToken, userProfile.Id, playlistId, tracks)
	if err != nil {
		return fmt.Errorf("CreatePlaylist: error adding tracks: %v", err)
	}

	return nil
}

func getAlbumIds(albums []model.Album) []string {
	ids := make([]string, 0, len(albums))

	for _, album := range albums {
		ids = append(ids, album.Id)
	}

	return ids
}

func min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}
