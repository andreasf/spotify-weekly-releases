package services_test

import (
	. "github.com/andreasf/spotify-weekly-releases/services"

	"github.com/andreasf/spotify-weekly-releases/api/apifakes"
	"github.com/andreasf/spotify-weekly-releases/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strconv"
)

var _ = Describe("SpotifyService", func() {
	Describe("GetRecentReleases", func() {
		var expectedAlbums []model.Album
		var followedArtists []model.Artist
		var fooAlbumInfo model.Album
		var client *apifakes.FakeSpotifyConnector
		var service *SpotifyServiceImpl

		BeforeEach(func() {
			expectedAlbums = []model.Album{
				{
					Name:        "foo-album",
					Id:          "foo-album-id",
					ReleaseDate: "2017-01-01",
				},
			}
			followedArtists = []model.Artist{
				{
					Name: "foo",
					Id:   "foo-id",
				},
			}

			fooAlbumInfo = model.Album{
				Name:        "foo-album",
				Id:          "foo-album-id",
				ReleaseDate: "2017-01-01",
			}

			client = &apifakes.FakeSpotifyConnector{}

			client.GetFollowedArtistsReturns(followedArtists, nil)
			client.GetAlbumInfoReturns([]model.Album{fooAlbumInfo}, nil)
			client.GetArtistAlbumsReturns(expectedAlbums, nil)
			client.GetUserProfileReturns(model.UserProfile{
				Id:      "user-id",
				Country: "market-id",
			}, nil)

			service = NewSpotifyService(client)
		})

		It("Gets the user profile in order to filter by country", func() {
			_, err := service.GetRecentReleases("access-token")

			Expect(err).To(BeNil())

			Expect(client.GetUserProfileCallCount()).To(Equal(1))
			Expect(client.GetUserProfileArgsForCall(0)).To(Equal("access-token"))
		})

		It("Returns a list of recent releases for the user's market", func() {
			service := NewSpotifyService(client)

			albums, err := service.GetRecentReleases("access-token")

			Expect(err).To(BeNil())
			Expect(albums).To(Equal(expectedAlbums))

			Expect(client.GetFollowedArtistsCallCount()).To(Equal(1))
			Expect(client.GetArtistAlbumsCallCount()).To(Equal(1))

			token, artistId, market := client.GetArtistAlbumsArgsForCall(0)
			Expect(token).To(Equal("access-token"))
			Expect(artistId).To(Equal("foo-id"))
			Expect(market).To(Equal("market-id"))

			Expect(client.GetAlbumInfoCallCount()).To(Equal(1))

			token1, id1 := client.GetAlbumInfoArgsForCall(0)
			Expect(token1).To(Equal("access-token"))
			Expect(id1).To(Equal([]string{"foo-album-id"}))
		})

		It("Gets album info for 20 albums at a time", func() {
			albums := []model.Album{}
			albumInfos := []model.Album{}
			albumIds := []string{}
			for i := 0; i < 23; i++ {
				albums = append(albums, model.Album{
					Name: "album-" + strconv.Itoa(i),
					Id:   "id-" + strconv.Itoa(i),
				})

				albumInfos = append(albumInfos, model.Album{
					Name:        "album-" + strconv.Itoa(i),
					Id:          "id-" + strconv.Itoa(i),
					ReleaseDate: "2017-01-01",
				})

				albumIds = append(albumIds, "id-"+strconv.Itoa(i))
			}

			client.GetArtistAlbumsReturns(albums, nil)
			client.GetAlbumInfoReturns(albumInfos, nil)

			albums, err := service.GetRecentReleases("access-token")

			Expect(err).To(BeNil())
			Expect(client.GetAlbumInfoCallCount()).To(Equal(2))

			_, ids1 := client.GetAlbumInfoArgsForCall(0)
			_, ids2 := client.GetAlbumInfoArgsForCall(1)
			Expect(ids1).To(Equal(albumIds[0:20]))
			Expect(ids2).To(Equal(albumIds[20:]))
		})

		//It("Does not return releases older than a year", func() {
		//	Fail("not implemented")
		//})
	})

	Describe("CreatePlaylist", func() {
		var tracks []model.Track
		var client *apifakes.FakeSpotifyConnector
		var service *SpotifyServiceImpl

		BeforeEach(func() {
			tracks = []model.Track{
				{
					Id: "track-1",
				},
				{
					Id: "track-2",
				},
				{
					Id: "track-3",
				},
			}

			client = &apifakes.FakeSpotifyConnector{}
			service = NewSpotifyService(client)

			user := model.UserProfile{
				Id: "my-user-id",
			}
			client.GetUserProfileReturns(user, nil)
			client.CreatePlaylistReturns("playlist-id", nil)
		})

		It("Gets the current user's id", func() {
			err := service.CreatePlaylist("access-token", "playlist name", tracks)

			Expect(err).To(BeNil())

			Expect(client.GetUserProfileCallCount()).To(Equal(1))
			Expect(client.GetUserProfileArgsForCall(0)).To(Equal("access-token"))
		})

		It("Creates a new playlist", func() {
			err := service.CreatePlaylist("access-token", "playlist name", tracks)

			Expect(err).To(BeNil())

			Expect(client.CreatePlaylistCallCount()).To(Equal(1))

			token, userId, name := client.CreatePlaylistArgsForCall(0)
			Expect(token).To(Equal("access-token"))
			Expect(userId).To(Equal("my-user-id"))
			Expect(name).To(Equal("playlist name"))
		})

		It("Adds all tracks to the playlist", func() {
			err := service.CreatePlaylist("access-token", "playlist name", tracks)

			Expect(err).To(BeNil())

			Expect(client.AddTracksToPlaylistCallCount()).To(Equal(1))

			token, userId, playlistId, actualTracks := client.AddTracksToPlaylistArgsForCall(0)
			Expect(token).To(Equal("access-token"))
			Expect(userId).To(Equal("my-user-id"))
			Expect(playlistId).To(Equal("playlist-id"))
			Expect(actualTracks).To(Equal(tracks))
		})
	})
})