package json_test

import (
	. "github.com/andreasf/spotify-weekly-releases/json"

	"encoding/json"
	"github.com/andreasf/spotify-weekly-releases/model"
	"github.com/andreasf/spotify-weekly-releases/test_resources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Schema", func() {
	var fooArtist Artist
	var barArtist Artist

	BeforeEach(func() {
		fooArtist = Artist{
			Id:   "foo-id",
			Name: "foo",
			Uri:  "spotify:artist:foo-id",
			Images: []Image{
				{
					Width:  640,
					Height: 640,
					Url:    "https://image.cdn/foo/1",
				},
				{
					Width:  320,
					Height: 320,
					Url:    "https://image.cdn/foo/2",
				},
				{
					Width:  160,
					Height: 160,
					Url:    "https://image.cdn/foo/3",
				},
			},
		}

		barArtist = Artist{
			Id:   "bar-id",
			Name: "bar",
			Uri:  "spotify:artist:bar-id",
			Images: []Image{
				{
					Width:  235,
					Height: 235,
					Url:    "https://image.cdn/bar/1",
				},
				{
					Width:  4242,
					Height: 4242,
					Url:    "https://image.cdn/bar/2",
				},
			},
		}
	})

	It("Deserializes the followed artists response", func() {
		rawJson := test_resources.LoadResource("../test_resources/me_following_page1.json")
		followed := FollowedArtists{}

		err := json.Unmarshal(rawJson, &followed)

		Expect(err).To(BeNil())
		Expect(followed.Artists.Next).To(Equal("${API_PREFIX}/v1/me/following?type=artist&after=bar-id&limit=20"))
		Expect(followed.Artists.Total).To(Equal(3))
		Expect(followed.Artists.Limit).To(Equal(2))
		Expect(followed.Artists.Href).To(Equal("${API_PREFIX}/v1/me/following?type=artist&limit=2"))

		Expect(followed.Artists.Items).To(HaveLen(2))
		Expect(followed.Artists.Items[0]).To(Equal(fooArtist))
		Expect(followed.Artists.Items[1]).To(Equal(barArtist))
	})

	It("Converts json.Artist to model.Artist", func() {
		Expect(fooArtist.ToModel()).To(Equal(model.Artist{
			Name: "foo",
			Id:   "foo-id",
		}))
	})

	It("Deserializes the artist albums response", func() {
		rawJson := test_resources.LoadResource("../test_resources/artist_albums_page1.json")
		albums := ArtistAlbums{}

		err := json.Unmarshal(rawJson, &albums)

		Expect(err).To(BeNil())
		Expect(albums.Next).To(Equal("${API_PREFIX}/v1/artists/foo-id/albums?offset=2&limit=2&album_type=single,album,compilation,appears_on,ep"))

		Expect(albums.Items).To(HaveLen(2))
		Expect(albums.Items[0]).To(Equal(ArtistAlbum{
			Name:             "Foo, The Album",
			Id:               "foo-album",
			AvailableMarkets: []string{"SG"},
		}))
		Expect(albums.Items[1]).To(Equal(ArtistAlbum{
			Name:             "Bar, The Album",
			Id:               "bar-album",
			AvailableMarkets: []string{"CA", "MX", "US"},
		}))
	})

	It("Deserializes the multiple albums response", func() {
		rawJson := test_resources.LoadResource("../test_resources/multiple_albums.json")
		multipleAlbums := MultipleAlbums{}

		err := json.Unmarshal(rawJson, &multipleAlbums)

		Expect(err).To(BeNil())
		Expect(multipleAlbums.Albums).To(HaveLen(3))
		Expect(multipleAlbums.Albums[0]).To(Equal(blackRadio))
		Expect(multipleAlbums.Albums[1]).To(Equal(dummy))
		Expect(multipleAlbums.Albums[2]).To(Equal(mezzanine))
	})

	It("Converts json.Album to model.Album", func() {
		Expect(mezzanine.ToModel()).To(Equal(model.Album{
			Name:        "Mezzanine",
			Id:          "49MNmJhZQewjt06rpwp6QR",
			ReleaseDate: "1998-04-20",
			Markets: []string{
				"AB", "CD",
			},
			Tracks: []model.Track{
				{
					Name:       "Angel",
					Id:         "7uv632EkfwYhXoqf8rhYrg",
					DurationMs: 379533,
				},
				{
					Name:       "Risingson",
					Id:         "6ggJ6MceyHGWtUg1KLp3M1",
					DurationMs: 298826,
				},
				{
					Name:       "Teardrop",
					Id:         "67Hna13dNDkZvBpTXRIaOJ",
					DurationMs: 330773,
				},
			},
		}))
	})

	It("Deserializes the user profile response", func() {
		rawJson := test_resources.LoadResource("../test_resources/user_profile.json")
		profile := UserProfile{}

		err := json.Unmarshal(rawJson, &profile)

		Expect(err).To(BeNil())
		Expect(profile.Id).To(Equal("user-id"))
		Expect(profile.Country).To(Equal("market-id"))
	})

	It("Converts json.UserProfile to model.UserProfile", func() {
		jsonProfile := &UserProfile{
			Id:      "user-id",
			Country: "market-id",
		}

		Expect(jsonProfile.ToModel()).To(Equal(model.UserProfile{
			Id:      "user-id",
			Country: "market-id",
		}))
	})

	It("Deserializes the create playlist response", func() {
		rawJson := test_resources.LoadResource("../test_resources/create_playlist_response.json")
		response := CreatePlaylistResponse{}

		err := json.Unmarshal(rawJson, &response)

		Expect(err).To(BeNil())
		Expect(response.Id).To(Equal("playlist-id"))
	})

})

var blackRadio ArtistAlbum = ArtistAlbum{
	Name:                 "Black Radio 2 (Deluxe)",
	Id:                   "6D6v2UODKxjwVnySdEjpEX",
	ReleaseDate:          "2013-01-01",
	ReleaseDatePrecision: "day",
	AvailableMarkets:     []string{"AD", "AR"},
	Tracks: Tracks{
		Items: []Track{
			{
				Id:          "2fus3smj2B653XXRGoLiYd",
				Name:        "Baby Tonight - Black Radio 2 Theme/Mic Check 2",
				TrackNumber: 1,
				DurationMs:  263653,
			},
			{
				Id:          "1NlTw8XgQrrWJEnXDlk7iq",
				Name:        "I Stand Alone",
				TrackNumber: 2,
				DurationMs:  293960,
			},
			{
				Id:          "5VgshfaTxvVsaJUanqkz8u",
				Name:        "What Are We Doing",
				TrackNumber: 3,
				DurationMs:  214920,
			},
		},
	},
}

var mezzanine ArtistAlbum = ArtistAlbum{
	Name:                 "Mezzanine",
	Id:                   "49MNmJhZQewjt06rpwp6QR",
	ReleaseDate:          "1998-04-20",
	ReleaseDatePrecision: "day",
	AvailableMarkets: []string{
		"AB", "CD",
	},
	Tracks: Tracks{
		Items: []Track{
			{
				Id:          "7uv632EkfwYhXoqf8rhYrg",
				Name:        "Angel",
				TrackNumber: 1,
				DurationMs:  379533,
			},
			{
				Id:          "6ggJ6MceyHGWtUg1KLp3M1",
				Name:        "Risingson",
				TrackNumber: 2,
				DurationMs:  298826,
			},
			{
				Id:          "67Hna13dNDkZvBpTXRIaOJ",
				Name:        "Teardrop",
				TrackNumber: 3,
				DurationMs:  330773,
			},
		},
	},
}

var dummy ArtistAlbum = ArtistAlbum{
	Name:                 "Dummy (Non UK Version)",
	Id:                   "3539EbNgIdEDGBKkUf4wno",
	ReleaseDate:          "1994-01-01",
	ReleaseDatePrecision: "day",
	AvailableMarkets:     []string{"CA", "MX", "US"},
	Tracks: Tracks{
		Items: []Track{
			{
				Id:          "2O6X9nPVVQSefg3xOQAo5u",
				Name:        "Mysterons",
				TrackNumber: 1,
				DurationMs:  306200,
			},
			{
				Id:          "6vTtCOimcPs5H1Jr9d0Aep",
				Name:        "Sour Times",
				TrackNumber: 2,
				DurationMs:  254000,
			},
			{
				Id:          "6pW8YspamPCxUwgvYttTSc",
				Name:        "Strangers",
				TrackNumber: 3,
				DurationMs:  238000,
			},
		},
	},
}
