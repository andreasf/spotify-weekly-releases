// This file was generated by counterfeiter
package apifakes

import (
	"sync"

	"github.com/andreasf/spotify-weekly-releases/api"
	"github.com/andreasf/spotify-weekly-releases/model"
)

type FakeSpotifyConnector struct {
	AddTracksToPlaylistStub        func(accessToken, userId, playlistId string, tracks []model.Track) error
	addTracksToPlaylistMutex       sync.RWMutex
	addTracksToPlaylistArgsForCall []struct {
		accessToken string
		userId      string
		playlistId  string
		tracks      []model.Track
	}
	addTracksToPlaylistReturns struct {
		result1 error
	}
	CreatePlaylistStub        func(accessToken, userId, name string) (string, error)
	createPlaylistMutex       sync.RWMutex
	createPlaylistArgsForCall []struct {
		accessToken string
		userId      string
		name        string
	}
	createPlaylistReturns struct {
		result1 string
		result2 error
	}
	GetAlbumInfoStub        func(accessToken string, albumIds []string) ([]model.Album, error)
	getAlbumInfoMutex       sync.RWMutex
	getAlbumInfoArgsForCall []struct {
		accessToken string
		albumIds    []string
	}
	getAlbumInfoReturns struct {
		result1 []model.Album
		result2 error
	}
	GetArtistAlbumsStub        func(accessToken, artistId string, market string) ([]model.Album, error)
	getArtistAlbumsMutex       sync.RWMutex
	getArtistAlbumsArgsForCall []struct {
		accessToken string
		artistId    string
		market      string
	}
	getArtistAlbumsReturns struct {
		result1 []model.Album
		result2 error
	}
	GetFollowedArtistsStub        func(accessToken string) ([]model.Artist, error)
	getFollowedArtistsMutex       sync.RWMutex
	getFollowedArtistsArgsForCall []struct {
		accessToken string
	}
	getFollowedArtistsReturns struct {
		result1 []model.Artist
		result2 error
	}
	GetSavedAlbumsStub        func(accessToken string) ([]model.Album, error)
	getSavedAlbumsMutex       sync.RWMutex
	getSavedAlbumsArgsForCall []struct {
		accessToken string
	}
	getSavedAlbumsReturns struct {
		result1 []model.Album
		result2 error
	}
	GetUserProfileStub        func(accessToken string) (model.UserProfile, error)
	getUserProfileMutex       sync.RWMutex
	getUserProfileArgsForCall []struct {
		accessToken string
	}
	getUserProfileReturns struct {
		result1 model.UserProfile
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSpotifyConnector) AddTracksToPlaylist(accessToken string, userId string, playlistId string, tracks []model.Track) error {
	var tracksCopy []model.Track
	if tracks != nil {
		tracksCopy = make([]model.Track, len(tracks))
		copy(tracksCopy, tracks)
	}
	fake.addTracksToPlaylistMutex.Lock()
	fake.addTracksToPlaylistArgsForCall = append(fake.addTracksToPlaylistArgsForCall, struct {
		accessToken string
		userId      string
		playlistId  string
		tracks      []model.Track
	}{accessToken, userId, playlistId, tracksCopy})
	fake.recordInvocation("AddTracksToPlaylist", []interface{}{accessToken, userId, playlistId, tracksCopy})
	fake.addTracksToPlaylistMutex.Unlock()
	if fake.AddTracksToPlaylistStub != nil {
		return fake.AddTracksToPlaylistStub(accessToken, userId, playlistId, tracks)
	}
	return fake.addTracksToPlaylistReturns.result1
}

func (fake *FakeSpotifyConnector) AddTracksToPlaylistCallCount() int {
	fake.addTracksToPlaylistMutex.RLock()
	defer fake.addTracksToPlaylistMutex.RUnlock()
	return len(fake.addTracksToPlaylistArgsForCall)
}

func (fake *FakeSpotifyConnector) AddTracksToPlaylistArgsForCall(i int) (string, string, string, []model.Track) {
	fake.addTracksToPlaylistMutex.RLock()
	defer fake.addTracksToPlaylistMutex.RUnlock()
	return fake.addTracksToPlaylistArgsForCall[i].accessToken, fake.addTracksToPlaylistArgsForCall[i].userId, fake.addTracksToPlaylistArgsForCall[i].playlistId, fake.addTracksToPlaylistArgsForCall[i].tracks
}

func (fake *FakeSpotifyConnector) AddTracksToPlaylistReturns(result1 error) {
	fake.AddTracksToPlaylistStub = nil
	fake.addTracksToPlaylistReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSpotifyConnector) CreatePlaylist(accessToken string, userId string, name string) (string, error) {
	fake.createPlaylistMutex.Lock()
	fake.createPlaylistArgsForCall = append(fake.createPlaylistArgsForCall, struct {
		accessToken string
		userId      string
		name        string
	}{accessToken, userId, name})
	fake.recordInvocation("CreatePlaylist", []interface{}{accessToken, userId, name})
	fake.createPlaylistMutex.Unlock()
	if fake.CreatePlaylistStub != nil {
		return fake.CreatePlaylistStub(accessToken, userId, name)
	}
	return fake.createPlaylistReturns.result1, fake.createPlaylistReturns.result2
}

func (fake *FakeSpotifyConnector) CreatePlaylistCallCount() int {
	fake.createPlaylistMutex.RLock()
	defer fake.createPlaylistMutex.RUnlock()
	return len(fake.createPlaylistArgsForCall)
}

func (fake *FakeSpotifyConnector) CreatePlaylistArgsForCall(i int) (string, string, string) {
	fake.createPlaylistMutex.RLock()
	defer fake.createPlaylistMutex.RUnlock()
	return fake.createPlaylistArgsForCall[i].accessToken, fake.createPlaylistArgsForCall[i].userId, fake.createPlaylistArgsForCall[i].name
}

func (fake *FakeSpotifyConnector) CreatePlaylistReturns(result1 string, result2 error) {
	fake.CreatePlaylistStub = nil
	fake.createPlaylistReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeSpotifyConnector) GetAlbumInfo(accessToken string, albumIds []string) ([]model.Album, error) {
	var albumIdsCopy []string
	if albumIds != nil {
		albumIdsCopy = make([]string, len(albumIds))
		copy(albumIdsCopy, albumIds)
	}
	fake.getAlbumInfoMutex.Lock()
	fake.getAlbumInfoArgsForCall = append(fake.getAlbumInfoArgsForCall, struct {
		accessToken string
		albumIds    []string
	}{accessToken, albumIdsCopy})
	fake.recordInvocation("GetAlbumInfo", []interface{}{accessToken, albumIdsCopy})
	fake.getAlbumInfoMutex.Unlock()
	if fake.GetAlbumInfoStub != nil {
		return fake.GetAlbumInfoStub(accessToken, albumIds)
	}
	return fake.getAlbumInfoReturns.result1, fake.getAlbumInfoReturns.result2
}

func (fake *FakeSpotifyConnector) GetAlbumInfoCallCount() int {
	fake.getAlbumInfoMutex.RLock()
	defer fake.getAlbumInfoMutex.RUnlock()
	return len(fake.getAlbumInfoArgsForCall)
}

func (fake *FakeSpotifyConnector) GetAlbumInfoArgsForCall(i int) (string, []string) {
	fake.getAlbumInfoMutex.RLock()
	defer fake.getAlbumInfoMutex.RUnlock()
	return fake.getAlbumInfoArgsForCall[i].accessToken, fake.getAlbumInfoArgsForCall[i].albumIds
}

func (fake *FakeSpotifyConnector) GetAlbumInfoReturns(result1 []model.Album, result2 error) {
	fake.GetAlbumInfoStub = nil
	fake.getAlbumInfoReturns = struct {
		result1 []model.Album
		result2 error
	}{result1, result2}
}

func (fake *FakeSpotifyConnector) GetArtistAlbums(accessToken string, artistId string, market string) ([]model.Album, error) {
	fake.getArtistAlbumsMutex.Lock()
	fake.getArtistAlbumsArgsForCall = append(fake.getArtistAlbumsArgsForCall, struct {
		accessToken string
		artistId    string
		market      string
	}{accessToken, artistId, market})
	fake.recordInvocation("GetArtistAlbums", []interface{}{accessToken, artistId, market})
	fake.getArtistAlbumsMutex.Unlock()
	if fake.GetArtistAlbumsStub != nil {
		return fake.GetArtistAlbumsStub(accessToken, artistId, market)
	}
	return fake.getArtistAlbumsReturns.result1, fake.getArtistAlbumsReturns.result2
}

func (fake *FakeSpotifyConnector) GetArtistAlbumsCallCount() int {
	fake.getArtistAlbumsMutex.RLock()
	defer fake.getArtistAlbumsMutex.RUnlock()
	return len(fake.getArtistAlbumsArgsForCall)
}

func (fake *FakeSpotifyConnector) GetArtistAlbumsArgsForCall(i int) (string, string, string) {
	fake.getArtistAlbumsMutex.RLock()
	defer fake.getArtistAlbumsMutex.RUnlock()
	return fake.getArtistAlbumsArgsForCall[i].accessToken, fake.getArtistAlbumsArgsForCall[i].artistId, fake.getArtistAlbumsArgsForCall[i].market
}

func (fake *FakeSpotifyConnector) GetArtistAlbumsReturns(result1 []model.Album, result2 error) {
	fake.GetArtistAlbumsStub = nil
	fake.getArtistAlbumsReturns = struct {
		result1 []model.Album
		result2 error
	}{result1, result2}
}

func (fake *FakeSpotifyConnector) GetFollowedArtists(accessToken string) ([]model.Artist, error) {
	fake.getFollowedArtistsMutex.Lock()
	fake.getFollowedArtistsArgsForCall = append(fake.getFollowedArtistsArgsForCall, struct {
		accessToken string
	}{accessToken})
	fake.recordInvocation("GetFollowedArtists", []interface{}{accessToken})
	fake.getFollowedArtistsMutex.Unlock()
	if fake.GetFollowedArtistsStub != nil {
		return fake.GetFollowedArtistsStub(accessToken)
	}
	return fake.getFollowedArtistsReturns.result1, fake.getFollowedArtistsReturns.result2
}

func (fake *FakeSpotifyConnector) GetFollowedArtistsCallCount() int {
	fake.getFollowedArtistsMutex.RLock()
	defer fake.getFollowedArtistsMutex.RUnlock()
	return len(fake.getFollowedArtistsArgsForCall)
}

func (fake *FakeSpotifyConnector) GetFollowedArtistsArgsForCall(i int) string {
	fake.getFollowedArtistsMutex.RLock()
	defer fake.getFollowedArtistsMutex.RUnlock()
	return fake.getFollowedArtistsArgsForCall[i].accessToken
}

func (fake *FakeSpotifyConnector) GetFollowedArtistsReturns(result1 []model.Artist, result2 error) {
	fake.GetFollowedArtistsStub = nil
	fake.getFollowedArtistsReturns = struct {
		result1 []model.Artist
		result2 error
	}{result1, result2}
}

func (fake *FakeSpotifyConnector) GetSavedAlbums(accessToken string) ([]model.Album, error) {
	fake.getSavedAlbumsMutex.Lock()
	fake.getSavedAlbumsArgsForCall = append(fake.getSavedAlbumsArgsForCall, struct {
		accessToken string
	}{accessToken})
	fake.recordInvocation("GetSavedAlbums", []interface{}{accessToken})
	fake.getSavedAlbumsMutex.Unlock()
	if fake.GetSavedAlbumsStub != nil {
		return fake.GetSavedAlbumsStub(accessToken)
	}
	return fake.getSavedAlbumsReturns.result1, fake.getSavedAlbumsReturns.result2
}

func (fake *FakeSpotifyConnector) GetSavedAlbumsCallCount() int {
	fake.getSavedAlbumsMutex.RLock()
	defer fake.getSavedAlbumsMutex.RUnlock()
	return len(fake.getSavedAlbumsArgsForCall)
}

func (fake *FakeSpotifyConnector) GetSavedAlbumsArgsForCall(i int) string {
	fake.getSavedAlbumsMutex.RLock()
	defer fake.getSavedAlbumsMutex.RUnlock()
	return fake.getSavedAlbumsArgsForCall[i].accessToken
}

func (fake *FakeSpotifyConnector) GetSavedAlbumsReturns(result1 []model.Album, result2 error) {
	fake.GetSavedAlbumsStub = nil
	fake.getSavedAlbumsReturns = struct {
		result1 []model.Album
		result2 error
	}{result1, result2}
}

func (fake *FakeSpotifyConnector) GetUserProfile(accessToken string) (model.UserProfile, error) {
	fake.getUserProfileMutex.Lock()
	fake.getUserProfileArgsForCall = append(fake.getUserProfileArgsForCall, struct {
		accessToken string
	}{accessToken})
	fake.recordInvocation("GetUserProfile", []interface{}{accessToken})
	fake.getUserProfileMutex.Unlock()
	if fake.GetUserProfileStub != nil {
		return fake.GetUserProfileStub(accessToken)
	}
	return fake.getUserProfileReturns.result1, fake.getUserProfileReturns.result2
}

func (fake *FakeSpotifyConnector) GetUserProfileCallCount() int {
	fake.getUserProfileMutex.RLock()
	defer fake.getUserProfileMutex.RUnlock()
	return len(fake.getUserProfileArgsForCall)
}

func (fake *FakeSpotifyConnector) GetUserProfileArgsForCall(i int) string {
	fake.getUserProfileMutex.RLock()
	defer fake.getUserProfileMutex.RUnlock()
	return fake.getUserProfileArgsForCall[i].accessToken
}

func (fake *FakeSpotifyConnector) GetUserProfileReturns(result1 model.UserProfile, result2 error) {
	fake.GetUserProfileStub = nil
	fake.getUserProfileReturns = struct {
		result1 model.UserProfile
		result2 error
	}{result1, result2}
}

func (fake *FakeSpotifyConnector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addTracksToPlaylistMutex.RLock()
	defer fake.addTracksToPlaylistMutex.RUnlock()
	fake.createPlaylistMutex.RLock()
	defer fake.createPlaylistMutex.RUnlock()
	fake.getAlbumInfoMutex.RLock()
	defer fake.getAlbumInfoMutex.RUnlock()
	fake.getArtistAlbumsMutex.RLock()
	defer fake.getArtistAlbumsMutex.RUnlock()
	fake.getFollowedArtistsMutex.RLock()
	defer fake.getFollowedArtistsMutex.RUnlock()
	fake.getSavedAlbumsMutex.RLock()
	defer fake.getSavedAlbumsMutex.RUnlock()
	fake.getUserProfileMutex.RLock()
	defer fake.getUserProfileMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeSpotifyConnector) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ api.SpotifyConnector = new(FakeSpotifyConnector)
