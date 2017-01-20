package json

import (
	"github.com/andreasf/spotify-weekly-releases/model"
)

type FollowedArtists struct {
	Artists PaginatedArtists `json:"artists"`
}

type PaginatedArtists struct {
	Items Artists `json:"items"`
	Next  string  `json:"next"`
	Total int     `json:"total"`
	Limit int     `json:"limit"`
	Href  string  `json:"href"`
}

type Artists []Artist

type Artist struct {
	Id     string  `json:"id"`
	Images []Image `json:"images"`
	Name   string  `json:"name"`
	Uri    string  `json:"uri"`
}

type Image struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Url    string `json:"url"`
}

func (self Artist) ToModel() model.Artist {
	return model.Artist{
		Id:   self.Id,
		Name: self.Name,
	}
}

type ArtistAlbums struct {
	Items []ArtistAlbum `json:"items"`
	Next  string        `json:"next"`
}

type ArtistAlbum struct {
	Id                   string   `json:"id"`
	Artists              []Artist `json:"artists"`
	Name                 string   `json:"name"`
	ReleaseDate          string   `json:"release_date"`
	ReleaseDatePrecision string   `json:"release_date_precision"`
	AvailableMarkets     []string `json:"available_markets"`
	Tracks               Tracks   `json:"tracks"`
}

func (self ArtistAlbum) ToModel() model.Album {
	return model.Album{
		Id:          self.Id,
		ArtistId:    self.Artists[0].Id,
		Name:        self.Name,
		ReleaseDate: self.ReleaseDate,
		Markets:     self.AvailableMarkets,
		Tracks:      self.Tracks.ToModel(),
	}
}

type Tracks struct {
	Items []Track `json:"items"`
}

func (self Tracks) ToModel() []model.Track {
	tracks := make([]model.Track, 0, len(self.Items))

	for _, track := range self.Items {
		tracks = append(tracks, track.ToModel())
	}

	return tracks
}

type Track struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	TrackNumber int    `json:"track_number"`
	DurationMs  int    `json:"duration_ms"`
}

func (self Track) ToModel() model.Track {
	return model.Track{
		Name:       self.Name,
		Id:         self.Id,
		DurationMs: self.DurationMs,
	}
}

type MultipleAlbums struct {
	Albums []ArtistAlbum `json:"albums"`
}

func (self MultipleAlbums) ToModel() []model.Album {
	albums := make([]model.Album, 0, len(self.Albums))

	for _, album := range self.Albums {
		albums = append(albums, album.ToModel())
	}

	return albums
}

type ArtistAlbumList []ArtistAlbum

func (self ArtistAlbumList) ToModel() []model.Album {
	albums := make([]model.Album, 0, len(self))

	for _, album := range self {
		albums = append(albums, album.ToModel())
	}

	return albums
}

type UserProfile struct {
	Id      string `json:"id"`
	Country string `json:"country"`
}

func (self *UserProfile) ToModel() model.UserProfile {
	return model.UserProfile{
		Id:      self.Id,
		Country: self.Country,
	}
}

type CreatePlaylistRequest struct {
	Name   string `json:"name"`
	Public bool   `json:"public"`
}

type CreatePlaylistResponse struct {
	Id string `json:"id"`
}

type AddTracksRequest struct {
	Uris []string `json:"uris"`
}
