package model_test

import (
	. "github.com/andreasf/spotify-release-feed/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model", func() {
	Describe("Album", func() {
		var shortAlbum *Album
		var longAlbum *Album
		var emptyAlbum *Album

		BeforeEach(func() {
			shortAlbum = &Album{
				Tracks: []Track{
					{
						Id: "track-1",
					},
					{
						Id: "track-2",
					},
					{
						Id: "track-3",
					},
				},
			}
			longAlbum = &Album{
				Tracks: []Track{
					{
						Id: "track-1",
					},
					{
						Id: "track-2",
					},
					{
						Id: "track-3",
					},
					{
						Id: "track-4",
					},
				},
			}
			emptyAlbum = &Album{
				Tracks: []Track{},
			}
		})

		Describe("GetSampleTrack", func() {
			It("Returns the 3rd track if the album has more than 3 tracks", func() {
				Expect(longAlbum.GetSampleTrack().Id).To(Equal("track-3"))
			})

			It("Returns the 1st track if the album has up to 3 tracks", func() {
				Expect(shortAlbum.GetSampleTrack().Id).To(Equal("track-1"))
			})

			It("Returns nil if the album has no tracks", func() {
				Expect(emptyAlbum.GetSampleTrack()).To(BeNil())
			})
		})
	})

	Describe("TrackList", func() {
		It("GetUris returns the list of track URIs", func() {
			var tracks TrackList = []Track{
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

			Expect(tracks.GetUris()).To(Equal([]string{
				"spotify:track:track-1",
				"spotify:track:track-2",
				"spotify:track:track-3",
			}))
		})
	})
})
