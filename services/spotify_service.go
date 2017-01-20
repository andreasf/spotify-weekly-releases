package services

import (
	"fmt"
	"github.com/andreasf/spotify-weekly-releases/api"
	"github.com/andreasf/spotify-weekly-releases/model"
	"github.com/andreasf/spotify-weekly-releases/platform"
	"time"
)

type SpotifyService interface {
	GetRecentReleases(accessToken string) ([]model.Album, error)
	CreatePlaylist(accessToken string, name string, tracks []model.Track) error
}

type SpotifyServiceImpl struct {
	apiClient   api.SpotifyConnector
	timeWrapper platform.Time
}

const ALBUMS_PER_REQUEST int = 20

func NewSpotifyService(apiClient api.SpotifyConnector, timeWrapper platform.Time) *SpotifyServiceImpl {
	return &SpotifyServiceImpl{
		apiClient:   apiClient,
		timeWrapper: timeWrapper,
	}
}

func (self *SpotifyServiceImpl) GetRecentReleases(accessToken string) ([]model.Album, error) {
	profile, err := self.apiClient.GetUserProfile(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetRecentReleases: error retrieving user profile: %v", err)
	}

	var artists model.ArtistList
	artists, err = self.apiClient.GetFollowedArtists(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetRecentReleases: error retrieving followed artists: %v", err)
	}

	artistIds := artists.GetIds()

	var savedAlbums model.AlbumList
	savedAlbums, err = self.apiClient.GetSavedAlbums(accessToken)
	if err != nil {
		return nil, fmt.Errorf("GetRecentReleases: error retrieving saved albums: %v", err)
	}

	artistIds = append(artistIds, savedAlbums.GetArtistIds()...)

	var albums model.AlbumList
	albums, err = self.getAlbumsForArtists(accessToken, profile.Country, artistIds)
	if err != nil {
		return nil, fmt.Errorf("GetRecentReleases: %v", err)
	}

	albums = albums.Remove(savedAlbums)

	albumDetails, err := self.getAlbumDetails(accessToken, albums)
	if err != nil {
		return nil, fmt.Errorf("GetRecentReleases: %v", err)
	}

	return self.filterByReleaseDate(albumDetails), nil
}

func (self *SpotifyServiceImpl) getAlbumsForArtists(accessToken string, country string, artistIds []string) ([]model.Album, error) {
	visitedArtists := make(map[string]bool)
	visitedAlbums := make(map[string]bool)
	albums := []model.Album{}

	for _, artistId := range artistIds {
		_, visited := visitedArtists[artistId]
		if visited {
			continue
		}

		artistAlbums, err := self.apiClient.GetArtistAlbums(accessToken, artistId, country)
		if err != nil {
			return nil, fmt.Errorf("getAlbumsForArtists: error retrieving artistAlbums for %s: %v", artistId, err)
		}

		for _, album := range artistAlbums {
			_, visited := visitedAlbums[album.Id]
			if visited {
				continue
			}

			albums = append(albums, album)
			visitedAlbums[album.Id] = true
		}

		visitedArtists[artistId] = true
	}

	return albums, nil
}

func (self *SpotifyServiceImpl) getAlbumDetails(accessToken string, albums []model.Album) ([]model.Album, error) {
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
			return nil, fmt.Errorf("getAlbumDetails: error retrieving album infos for %s: %v", albumIds, err)
		}

		albumDetails = append(albumDetails, albumInfos...)
	}

	return albumDetails, nil
}

func (self *SpotifyServiceImpl) filterByReleaseDate(albums []model.Album) []model.Album {
	filteredAlbums := make([]model.Album, 0, len(albums))

	oneYearAgoDate := self.timeWrapper.Now().Add(time.Hour * 24 * (-365))
	oneYearAgo := oneYearAgoDate.Format("2006-01-02")

	for _, album := range albums {
		if album.ReleaseDate >= oneYearAgo {
			filteredAlbums = append(filteredAlbums, album)
		}
	}

	return filteredAlbums
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
