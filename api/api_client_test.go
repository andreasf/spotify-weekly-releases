package api_test

import (
	. "github.com/andreasf/spotify-weekly-releases/api"

	json2 "encoding/json"
	"errors"
	"github.com/andreasf/spotify-weekly-releases/cache/cachefakes"
	"github.com/andreasf/spotify-weekly-releases/json"
	"github.com/andreasf/spotify-weekly-releases/model"
	"github.com/andreasf/spotify-weekly-releases/platform/platformfakes"
	"github.com/andreasf/spotify-weekly-releases/test_resources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var _ = Describe("SpotifyApiClient", func() {
	Describe("GetFollowedArtists", func() {
		It("Makes a GET request to the endpoint", func() {
			expectedArtists := []model.Artist{
				{
					Name: "foo",
					Id:   "foo-id",
				},
				{
					Name: "bar",
					Id:   "bar-id",
				},
				{
					Name: "baz",
					Id:   "baz-id",
				},
			}

			server := ghttp.NewServer()

			page1 := test_resources.LoadResource("../test_resources/me_following_page1.json")
			page1 = replaceApiPrefix(page1, server.URL())

			page2 := test_resources.LoadResource("../test_resources/me_following_page2.json")
			page2 = replaceApiPrefix(page2, server.URL())

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/me/following", "type=artist&limit=50"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page1),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/me/following", "type=artist&after=bar-id&limit=20"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page2),
				),
			)

			cache := &cachefakes.FakeCache{}
			timeWrapper := &platformfakes.FakeTime{}
			client := NewSpotifyApiClient(server.URL(), timeWrapper, cache)
			artists, err := client.GetFollowedArtists("access-token")

			Expect(err).To(BeNil())
			Expect(artists).To(Equal(expectedArtists))
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})
	})

	Describe("GetArtistAlbums", func() {
		var expectedAlbums []model.Album
		var server *ghttp.Server
		var page1, page2 []byte
		var timeWrapper *platformfakes.FakeTime
		var cache *cachefakes.FakeCache

		BeforeEach(func() {
			expectedAlbums = []model.Album{
				{
					Name:      "Foo, The Album",
					Id:        "foo-album",
					ArtistIds: []string{"foo-id"},
					Tracks:    []model.Track{},
					Markets:   []string{"SG"},
				},
				{
					Name:      "Bar, The Album",
					Id:        "bar-album",
					ArtistIds: []string{"foo-id"},
					Tracks:    []model.Track{},
					Markets:   []string{"CA", "MX", "US"},
				},
				{
					Name:      "Baz, The Album",
					Id:        "baz-album",
					ArtistIds: []string{"foo-id"},
					Tracks:    []model.Track{},
					Markets:   []string{"SG"},
				},
			}

			server = ghttp.NewServer()

			page1 = test_resources.LoadResource("../test_resources/artist_albums_page1.json")
			page1 = replaceApiPrefix(page1, server.URL())

			page2 = test_resources.LoadResource("../test_resources/artist_albums_page2.json")
			page2 = replaceApiPrefix(page2, server.URL())

			timeWrapper = &platformfakes.FakeTime{}
			cache = &cachefakes.FakeCache{}
			cache.GetReturns(nil, errors.New("not found"))
		})

		It("Makes a GET request to the endpoint", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/artists/foo-id/albums", "album_type=album&limit=50&market=market-id"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page1),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/artists/foo-id/albums", "offset=2&limit=2&album_type=single,album,compilation,appears_on,ep"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page2),
				),
			)

			client := NewSpotifyApiClient(server.URL(), timeWrapper, cache)
			albums, err := client.GetArtistAlbums("access-token", "foo-id", "market-id")

			Expect(err).To(BeNil())
			Expect(albums).To(Equal(expectedAlbums))
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})

		It("Rate-limits requests", func() {
			var retryHeader http.Header
			retryHeader = make(map[string][]string)
			retryHeader.Add("Retry-After", "2")

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/artists/foo-id/albums", "album_type=album&limit=50&market=market-id"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page1),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/artists/foo-id/albums", "offset=2&limit=2&album_type=single,album,compilation,appears_on,ep"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(429, []byte{}, retryHeader),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/artists/foo-id/albums", "offset=2&limit=2&album_type=single,album,compilation,appears_on,ep"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page2),
				),
			)

			client := NewSpotifyApiClient(server.URL(), timeWrapper, cache)
			albums, err := client.GetArtistAlbums("access-token", "foo-id", "market-id")

			Expect(err).To(BeNil())
			Expect(albums).To(Equal(expectedAlbums))
			Expect(server.ReceivedRequests()).Should(HaveLen(3))
			Expect(timeWrapper.SleepCallCount()).To(Equal(1))
			Expect(timeWrapper.SleepArgsForCall(0)).To(Equal(2 * time.Second))
		})

		It("Checks the cache before making HTTP requests", func() {
			cache.GetReturns(page2, nil)
			page2Albums := []model.Album{
				expectedAlbums[2],
			}

			client := NewSpotifyApiClient(server.URL(), timeWrapper, cache)
			albums, err := client.GetArtistAlbums("access-token", "foo-id", "market-id")

			Expect(err).To(BeNil())
			Expect(albums).To(Equal(page2Albums))
			Expect(server.ReceivedRequests()).Should(HaveLen(0))
			Expect(cache.GetCallCount()).To(Equal(1))
			Expect(cache.GetArgsForCall(0)).To(Equal(server.URL() + "/v1/artists/foo-id/albums?album_type=album&limit=50&market=market-id"))
		})

		It("Stores responses in the cache", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/artists/foo-id/albums", "album_type=album&limit=50&market=market-id"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page1),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/artists/foo-id/albums", "offset=2&limit=2&album_type=single,album,compilation,appears_on,ep"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page2),
				),
			)

			client := NewSpotifyApiClient(server.URL(), timeWrapper, cache)
			albums, err := client.GetArtistAlbums("access-token", "foo-id", "market-id")

			Expect(err).To(BeNil())
			Expect(albums).ToNot(BeNil())
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
			Expect(cache.SetCallCount()).To(Equal(2))

			key1, data1 := cache.SetArgsForCall(0)
			key2, data2 := cache.SetArgsForCall(1)
			Expect(key1).To(Equal(server.URL() + "/v1/artists/foo-id/albums?album_type=album&limit=50&market=market-id"))
			Expect(key2).To(Equal(server.URL() + "/v1/artists/foo-id/albums?offset=2&limit=2&album_type=single,album,compilation,appears_on,ep"))
			Expect(data1).To(Equal(page1))
			Expect(data2).To(Equal(page2))
		})
	})

	Describe("GetAlbumInfo", func() {
		var server *ghttp.Server
		var timeWrapper *platformfakes.FakeTime
		var cache *cachefakes.FakeCache
		var response []byte
		var client *SpotifyApiClient
		var albumIds []string

		BeforeEach(func() {
			response = test_resources.LoadResource("../test_resources/multiple_albums.json")

			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/albums", "ids=album-id-1,album-id-2,album-id-3"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, response),
				),
			)

			timeWrapper = &platformfakes.FakeTime{}
			cache = &cachefakes.FakeCache{}
			client = NewSpotifyApiClient(server.URL(), timeWrapper, cache)

			albumIds = []string{
				"album-id-1",
				"album-id-2",
				"album-id-3",
			}
		})

		It("Makes a single GET request", func() {
			_, err := client.GetAlbumInfo("access-token", albumIds)

			Expect(err).To(BeNil())
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("Returns the matching []model.Album", func() {
			albums, err := client.GetAlbumInfo("access-token", albumIds)

			Expect(err).To(BeNil())
			Expect(albums).To(HaveLen(3))
			Expect(albums[0].Name).To(Equal("Black Radio 2 (Deluxe)"))
		})

		Describe("Caching", func() {
			BeforeEach(func() {
				server.Reset()
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/v1/albums", "ids=album-id-1,album-id-3"),
						ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
						ghttp.RespondWith(200, response),
					),
				)

				cache.GetStub = func(key string) ([]byte, error) {
					switch key {
					case "album:album-id-2":
						return test_resources.LoadResource("../test_resources/cached_album.json"), nil
					}

					return nil, errors.New("not found")
				}
			})

			It("Checks if individual albums are cached", func() {
				_, err := client.GetAlbumInfo("access-token", albumIds)
				Expect(err).To(BeNil())

				Expect(cache.GetCallCount()).To(Equal(3))
				Expect(cache.GetArgsForCall(0)).To(Equal("album:album-id-1"))
				Expect(cache.GetArgsForCall(1)).To(Equal("album:album-id-2"))
				Expect(cache.GetArgsForCall(2)).To(Equal("album:album-id-3"))
			})

			It("Does not query the Spotify API for cached albums", func() {
				_, err := client.GetAlbumInfo("access-token", albumIds)

				Expect(err).To(BeNil())
				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			It("Returns album data from both the cache and the HTTP API", func() {
				albums, err := client.GetAlbumInfo("access-token", albumIds)

				Expect(err).To(BeNil())
				Expect(albums).To(HaveLen(4))
			})

			It("Stores individual albums in the cache", func() {
				_, err := client.GetAlbumInfo("access-token", albumIds)

				Expect(err).To(BeNil())
				Expect(cache.SetCallCount()).To(Equal(3))

				assertAlbumCached(cache, 0, "6D6v2UODKxjwVnySdEjpEX")
				assertAlbumCached(cache, 1, "3539EbNgIdEDGBKkUf4wno")
				assertAlbumCached(cache, 2, "49MNmJhZQewjt06rpwp6QR")
			})

			It("Does not call the API if everything is cached", func() {
				server.Reset()

				cache.GetStub = func(key string) ([]byte, error) {
					return test_resources.LoadResource("../test_resources/cached_album.json"), nil
				}

				_, err := client.GetAlbumInfo("access-token", albumIds)

				Expect(err).To(BeNil())
				Expect(cache.SetCallCount()).To(Equal(0))
				Expect(server.ReceivedRequests()).To(HaveLen(0))
			})
		})
	})

	Describe("GetUserProfile", func() {
		var server *ghttp.Server
		var timeWrapper *platformfakes.FakeTime
		var cache *cachefakes.FakeCache
		var client *SpotifyApiClient

		BeforeEach(func() {
			server = ghttp.NewServer()
			timeWrapper = &platformfakes.FakeTime{}
			cache = &cachefakes.FakeCache{}

			response := test_resources.LoadResource("../test_resources/user_profile.json")
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/me", ""),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, response),
				),
			)

			client = NewSpotifyApiClient(server.URL(), timeWrapper, cache)
		})

		It("Calls the HTTP API", func() {
			profile, err := client.GetUserProfile("access-token")

			Expect(err).To(BeNil())
			Expect(profile).To(Equal(model.UserProfile{
				Id:      "user-id",
				Country: "market-id",
			}))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})

	Describe("CreatePlaylist", func() {
		var server *ghttp.Server
		var timeWrapper *platformfakes.FakeTime
		var cache *cachefakes.FakeCache
		var client *SpotifyApiClient

		BeforeEach(func() {
			server = ghttp.NewServer()
			timeWrapper = &platformfakes.FakeTime{}
			cache = &cachefakes.FakeCache{}

			response := test_resources.LoadResource("../test_resources/create_playlist_response.json")
			expectedBody := test_resources.LoadResource("../test_resources/create_playlist_request.json")
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/users/user-id/playlists", ""),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.VerifyJSON(string(expectedBody)),
					ghttp.RespondWith(201, response),
				),
			)

			client = NewSpotifyApiClient(server.URL(), timeWrapper, cache)
		})

		It("POSTs to the HTTP API", func() {
			playlistId, err := client.CreatePlaylist("access-token", "user-id", "playlist name")

			Expect(err).To(BeNil())
			Expect(playlistId).To(Equal("playlist-id"))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})

	Describe("AddTracksToPlaylist", func() {
		var server *ghttp.Server
		var timeWrapper *platformfakes.FakeTime
		var cache *cachefakes.FakeCache
		var client *SpotifyApiClient
		var tracks []model.Track

		BeforeEach(func() {
			server = ghttp.NewServer()
			timeWrapper = &platformfakes.FakeTime{}
			cache = &cachefakes.FakeCache{}

			expectedBody1 := test_resources.LoadResource("../test_resources/add_tracks_request_1.json")
			expectedBody2 := test_resources.LoadResource("../test_resources/add_tracks_request_2.json")
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/users/user-id/playlists/playlist-id/tracks", ""),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.VerifyJSON(string(expectedBody1)),
					ghttp.RespondWith(201, nil),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/v1/users/user-id/playlists/playlist-id/tracks", ""),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.VerifyJSON(string(expectedBody2)),
					ghttp.RespondWith(201, nil),
				),
			)

			client = NewSpotifyApiClient(server.URL(), timeWrapper, cache)

			tracks = make([]model.Track, 0, 123)
			for i := 1; i < 124; i++ {
				tracks = append(tracks, model.Track{
					Id: "track-" + strconv.Itoa(i),
				})
			}
		})

		It("POSTs to the HTTP API, 100 tracks at a time", func() {
			err := client.AddTracksToPlaylist("access-token", "user-id", "playlist-id", tracks)

			Expect(err).To(BeNil())

			Expect(server.ReceivedRequests()).To(HaveLen(2))
		})
	})

	Describe("GetSavedAlbums", func() {
		var cache *cachefakes.FakeCache
		var client *SpotifyApiClient
		var expectedAlbums []model.Album
		var server *ghttp.Server
		var timeWrapper *platformfakes.FakeTime

		BeforeEach(func() {
			expectedAlbums = []model.Album{
				{
					Name: "Groovements",
					Id:   "4xjys0dhhX8AD2Oiz5Y5S6",
					ArtistIds: []string{
						"22KzEvCtrTGf9l6k7zFcdv",
						"2GWMZZQNuU0VZra0suXVph",
						"1eLFONDpKa9ArYaoVjDrKE",
					},
					ReleaseDate: "2016-04-15",

					Tracks: []model.Track{
						{
							Name:       "Winter",
							Id:         "6YeQJTy8BAiTdDUooihG9p",
							ArtistId:   "2GWMZZQNuU0VZra0suXVph",
							DurationMs: 304493,
						},
					},
					Markets: []string{"AD", "AR"},
				},
				{
					Name:      "Invisible Cinema",
					Id:        "3xfueIrMUw57owAiYVKt8S",
					ArtistIds: []string{"22KzEvCtrTGf9l6k7zFcdv"},
					ReleaseDate: "2008-08-19",
					Tracks: []model.Track{
						{
							Name:       "Travelers",
							Id:         "5vRHlju25fl2hY0IPwTLbS",
							ArtistId:   "22KzEvCtrTGf9l6k7zFcdv",
							DurationMs: 334880,
						},
					},
					Markets: []string{"AD", "AR"},
				},
				{
					Name:      "Senzo",
					Id:        "2I3odMRAs5aHC69TMt9qAj",
					ArtistIds: []string{"39mb0I6tdTcCXkeigvzxOJ"},
					ReleaseDate: "2008-09-26",
					Tracks: []model.Track{
						{
							Name:       "Ocean & The River",
							Id:         "4enJS8L7a0w0DVDrQJXbUk",
							ArtistId:   "39mb0I6tdTcCXkeigvzxOJ",
							DurationMs: 165226,
						},
					},
					Markets: []string{"AD", "AT"},
				},
			}

			server = ghttp.NewServer()
			timeWrapper = &platformfakes.FakeTime{}
			cache = &cachefakes.FakeCache{}

			page1 := test_resources.LoadResource("../test_resources/saved_albums_page1.json")
			page1 = replaceApiPrefix(page1, server.URL())

			page2 := test_resources.LoadResource("../test_resources/saved_albums_page2.json")
			page2 = replaceApiPrefix(page2, server.URL())

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/me/albums", "limit=50"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page1),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/me/albums", "offset=2&limit=2"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, page2),
				),
			)

			client = NewSpotifyApiClient(server.URL(), timeWrapper, cache)
		})

		It("GETs from the HTTP API", func() {
			albums, err := client.GetSavedAlbums("access-token")

			Expect(err).To(BeNil())
			Expect(albums).To(Equal(expectedAlbums))

			Expect(server.ReceivedRequests()).To(HaveLen(2))
		})
	})
})

func replaceApiPrefix(jsonBytes []byte, apiPrefix string) []byte {
	return []byte(strings.Replace(string(jsonBytes), "${API_PREFIX}", apiPrefix, 1))
}

func assertAlbumCached(cache *cachefakes.FakeCache, call int, albumId string) {
	key1, data1 := cache.SetArgsForCall(call)
	Expect(key1).To(Equal("album:" + albumId))

	album1 := json.ArtistAlbum{}
	err := json2.Unmarshal(data1, &album1)
	Expect(err).To(BeNil())
	Expect(album1.Id).To(Equal(albumId))
}
