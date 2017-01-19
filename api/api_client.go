package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/andreasf/spotify-weekly-releases/cache"
	json2 "github.com/andreasf/spotify-weekly-releases/json"
	"github.com/andreasf/spotify-weekly-releases/model"
	"github.com/andreasf/spotify-weekly-releases/platform"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const TRACKS_PER_REQUEST int = 100

//go:generate counterfeiter . SpotifyConnector
type SpotifyConnector interface {
	GetFollowedArtists(accessToken string) ([]model.Artist, error)
	GetArtistAlbums(accessToken, artistId string, market string) ([]model.Album, error)
	GetAlbumInfo(accessToken string, albumIds []string) ([]model.Album, error)
	GetUserProfile(accessToken string) (model.UserProfile, error)
	CreatePlaylist(accessToken, userId, name string) (string, error)
	AddTracksToPlaylist(accessToken, userId, playlistId string, tracks []model.Track) error
}

type SpotifyApiClient struct {
	urlPrefix   string
	timeWrapper platform.Time
	cache       cache.Cache
}

func NewSpotifyApiClient(apiUrlPrefix string, timeWrapper platform.Time, cache cache.Cache) *SpotifyApiClient {
	return &SpotifyApiClient{
		urlPrefix:   apiUrlPrefix,
		timeWrapper: timeWrapper,
		cache:       cache,
	}
}

func (self *SpotifyApiClient) GetFollowedArtists(accessToken string) ([]model.Artist, error) {
	artists := []model.Artist{}
	nextUrl := self.urlPrefix + "/v1/me/following?type=artist"

	for nextUrl != "" {
		contents, err := self.getWithRateLimiting(accessToken, nextUrl)
		if err != nil {
			return nil, fmt.Errorf("GetFollowedArtists: request error: %v", err)
		}

		followedArtists := json2.FollowedArtists{}
		err = json.Unmarshal(contents, &followedArtists)
		if err != nil {
			return nil, fmt.Errorf("GetFollowedArtists: error deserializing JSON: %v", err)
		}

		nextUrl = followedArtists.Artists.Next

		for _, artist := range followedArtists.Artists.Items {
			artists = append(artists, artist.ToModel())
		}
	}

	return artists, nil
}

func (self *SpotifyApiClient) GetArtistAlbums(accessToken, artistId string, market string) ([]model.Album, error) {
	albums := []model.Album{}
	nextUrl := self.urlPrefix + "/v1/artists/" + artistId + "/albums?album_type=album&market=" + market

	for nextUrl != "" {
		contents, err := self.getWithRateLimitingAndCache(accessToken, nextUrl)
		if err != nil {
			return nil, fmt.Errorf("GetArtistAlbums: request error: %v", err)
		}

		artistAlbums := json2.ArtistAlbums{}
		err = json.Unmarshal(contents, &artistAlbums)
		if err != nil {
			return nil, fmt.Errorf("GetArtistAlbums: error deserializing JSON: %v", err)
		}

		nextUrl = artistAlbums.Next

		for _, album := range artistAlbums.Items {
			albums = append(albums, album.ToModel())
		}
	}

	return albums, nil
}

func (self *SpotifyApiClient) getWithRateLimitingAndCache(accessToken string, url string) ([]byte, error) {
	cached, err := self.cache.Get(url)
	if err == nil {
		return cached, nil
	}

	fromApi, err := self.getWithRateLimiting(accessToken, url)
	if err == nil {
		cacheErr := self.cache.Set(url, fromApi)
		if cacheErr != nil {
			log.Printf("SpotifyApiClient: error caching response: %v", cacheErr)
		}
	}
	return fromApi, err
}

func (self *SpotifyApiClient) getWithRateLimiting(accessToken string, url string) ([]byte, error) {
	return self.requestWithRateLimiting(accessToken, "GET", url, "", nil)
}

func (self *SpotifyApiClient) postWithRateLimiting(accessToken string, url string, contentType string, body []byte) ([]byte, error) {
	return self.requestWithRateLimiting(accessToken, "POST", url, contentType, bytes.NewBuffer(body))
}

func (self *SpotifyApiClient) requestWithRateLimiting(accessToken string, method string, url string, contentType string, body io.Reader) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("requestWithRateLimiting: error creating request: %v", err)
	}

	req.Header.Add("Authorization", "Bearer "+accessToken)

	if contentType != "" {
		req.Header.Add("Content-Type", "application/json")
	}

	for {
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("requestWithRateLimiting: error performing request: %v", err)
		}

		switch resp.StatusCode {
		case 200:
			fallthrough
		case 201:
			log.Printf("requestWithRateLimiting: %s %d %s", method, resp.StatusCode, url)
			contents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("requestWithRateLimiting: error reading response: %v", err)
			}

			return contents, nil

		case 429:
			retryAfter := resp.Header.Get("Retry-After")
			sleepSeconds, err := strconv.Atoi(retryAfter)
			if err != nil {
				sleepSeconds = 1
			}
			log.Printf("requestWithRateLimiting: 429 %s, retrying in %d seconds", url, sleepSeconds)
			self.timeWrapper.Sleep(time.Second * time.Duration(sleepSeconds))

		default:
			log.Printf("requestWithRateLimiting: %d %s", resp.StatusCode, url)
			return nil, fmt.Errorf("requestWithRateLimiting: received %d for GET %s", resp.StatusCode, url)
		}
	}
}

func (self *SpotifyApiClient) GetAlbumInfo(accessToken string, albumIds []string) ([]model.Album, error) {
	cachedAlbums, uncachedIds := self.getAlbumsFromCache(albumIds)

	apiAlbums := json2.MultipleAlbums{}
	if len(uncachedIds) > 0 {
		url := self.urlPrefix + "/v1/albums?ids=" + strings.Join(uncachedIds, ",")

		response, err := self.getWithRateLimiting(accessToken, url)
		if err != nil {
			return nil, fmt.Errorf("GetAlbumInfo: request error: %v", err)
		}

		err = json.Unmarshal(response, &apiAlbums)
		if err != nil {
			return nil, fmt.Errorf("GetAlbumInfo: error deserializing JSON: %v", err)
		}

		self.cacheAlbums(apiAlbums)
	}

	allAlbums := make([]model.Album, 0, len(apiAlbums.Albums)+len(cachedAlbums))
	allAlbums = append(allAlbums, apiAlbums.ToModel()...)
	allAlbums = append(allAlbums, cachedAlbums.ToModel()...)
	return allAlbums, nil
}

func (self *SpotifyApiClient) getAlbumsFromCache(albumIds []string) (json2.ArtistAlbumList, []string) {
	var cachedAlbums json2.ArtistAlbumList = []json2.ArtistAlbum{}
	uncachedIds := make([]string, 0, len(albumIds))
	for _, albumId := range albumIds {
		albumBytes, err := self.cache.Get("album:" + albumId)
		if err != nil {
			uncachedIds = append(uncachedIds, albumId)
			continue
		}

		album := json2.ArtistAlbum{}
		err = json.Unmarshal(albumBytes, &album)
		if err != nil {
			log.Printf("GetAlbumInfo: error deserializing cached album: %v", err)
			uncachedIds = append(uncachedIds, albumId)
			continue
		}
		cachedAlbums = append(cachedAlbums, album)
	}
	return cachedAlbums, uncachedIds
}

func (self *SpotifyApiClient) cacheAlbums(apiAlbums json2.MultipleAlbums) {
	for _, album := range apiAlbums.Albums {
		albumJson, err := json.Marshal(album)
		if err != nil {
			log.Printf("GetAlbumInfo: error serializing album for cache: %v", err)
			continue
		}

		self.cache.Set("album:"+album.Id, albumJson)
	}
}

func (self *SpotifyApiClient) GetUserProfile(accessToken string) (model.UserProfile, error) {
	url := self.urlPrefix + "/v1/me"

	response, err := self.getWithRateLimiting(accessToken, url)
	if err != nil {
		return model.UserProfile{}, fmt.Errorf("GetUserProfile: request error: %v", err)
	}

	jsonProfile := &json2.UserProfile{}
	err = json.Unmarshal(response, jsonProfile)
	if err != nil {
		return model.UserProfile{}, fmt.Errorf("GetUserProfile: error deserializing JSON: %v", err)
	}

	return jsonProfile.ToModel(), nil
}

func (self *SpotifyApiClient) CreatePlaylist(accessToken, userId, name string) (string, error) {
	url := fmt.Sprintf("%s/v1/users/%s/playlists", self.urlPrefix, userId)

	request := json2.CreatePlaylistRequest{
		Name:   name,
		Public: false,
	}

	body, err := json.Marshal(&request)
	if err != nil {
		return "", fmt.Errorf("CreatePlaylist: error serializing JSON: %v", err)
	}

	responseBytes, err := self.postWithRateLimiting(accessToken, url, "application/json", body)
	if err != nil {
		return "", fmt.Errorf("CreatePlaylist: request error: %v", err)
	}

	responseJson := json2.CreatePlaylistResponse{}
	err = json.Unmarshal(responseBytes, &responseJson)
	if err != nil {
		return "", fmt.Errorf("CreatePlaylist: error deserializing JSON: %v", err)
	}

	return responseJson.Id, nil
}

func (self *SpotifyApiClient) AddTracksToPlaylist(accessToken, userId, playlistId string, tracks []model.Track) error {
	url := fmt.Sprintf("%s/v1/users/%s/playlists/%s/tracks", self.urlPrefix, userId, playlistId)

	numberOfRequests := len(tracks) / TRACKS_PER_REQUEST
	if len(tracks)%TRACKS_PER_REQUEST > 0 {
		numberOfRequests++
	}

	for i := 0; i < numberOfRequests; i++ {
		from := i * TRACKS_PER_REQUEST
		to := min((i+1)*TRACKS_PER_REQUEST, len(tracks))
		var trackSlice model.TrackList = tracks[from:to]

		request := json2.AddTracksRequest{
			Uris: trackSlice.GetUris(),
		}
		body, err := json.Marshal(&request)
		if err != nil {
			return fmt.Errorf("AddTracksToPlaylist: error serializing JSON: %v", err)
		}

		_, err = self.postWithRateLimiting(accessToken, url, "application/json", body)
		if err != nil {
			return fmt.Errorf("AddTracksToPlaylist: request error: %v", err)
		}
	}

	return nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}

	return b
}
