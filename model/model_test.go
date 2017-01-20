package model_test

import (
	. "github.com/andreasf/spotify-weekly-releases/model"

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
		var duplicateTracks TrackList

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
			duplicateTracks = []Track{
				{
					Id:       "track-1",
					ArtistId: "artist-1",
				},
				{
					Id:       "track-1",
					ArtistId: "artist-1",
				},
				{
					Id:       "track-2",
					ArtistId: "artist-2",
				},
				{
					Id:       "track-2",
					ArtistId: "artist-3",
				},
			}

			Expect(tracks.GetUris()).To(Equal([]string{
				"spotify:track:track-1",
				"spotify:track:track-2",
				"spotify:track:track-3",
			}))
		})

		It("RemoveDuplicates removes duplicates based on artist id and track name", func() {
			filteredList := duplicateTracks.RemoveDuplicates()

			Expect(filteredList).To(HaveLen(3))
			Expect(filteredList[0].Id).To(Equal("track-1"))
			Expect(filteredList[1].Id).To(Equal("track-2"))
			Expect(filteredList[2].Id).To(Equal("track-2"))
		})
	})

	Describe("AlbumList", func() {
		var duplicateAlbumList AlbumList

		BeforeEach(func() {
			duplicateAlbumList = []Album{
				{
					Name:        "fooplicate",
					ArtistIds:   []string{"foo-id"},
					Id:          "baz-album-id",
					ReleaseDate: "2017-01-01",
				},
				{
					Name:        "fooplicate",
					Id:          "foo-album-id",
					ArtistIds:   []string{"foo-id"},
					ReleaseDate: "2016-06-02",
				},
				{
					Name:        "barnique",
					Id:          "bar-album-id",
					ArtistIds:   []string{"bar-id", "baz-id"},
					ReleaseDate: "2016-05-23",
				},
			}
		})

		Describe("RemoveDuplicates", func() {
			It("Removes duplicates based on artist id and album name", func() {
				filteredList := duplicateAlbumList.RemoveDuplicates()

				Expect(filteredList).To(HaveLen(2))
				Expect(filteredList[0].Id).To(Equal("baz-album-id"))
				Expect(filteredList[1].Id).To(Equal("bar-album-id"))
			})
		})

		Describe("GetArtistIds", func() {
			It("Returns the list of artists", func() {
				artistList := []string{
					"foo-id",
					"foo-id",
					"bar-id",
					"baz-id",
				}

				Expect(duplicateAlbumList.GetArtistIds()).To(Equal(artistList))
			})
		})

		It("Remove removes the given albums from the list", func() {
			toRemove := Album{
				Id: "baz-album-id",
			}

			newList := duplicateAlbumList.Remove([]Album{toRemove})

			Expect(newList).To(HaveLen(2))
			Expect(newList[0].Id).To(Equal("foo-album-id"))
			Expect(newList[1].Id).To(Equal("bar-album-id"))
		})
	})

	Describe("ArtistList", func() {
		It("GetIds returns list of artist ids", func() {
			var artistList ArtistList
			artistList = []Artist{
				{
					Name: "foo",
					Id:   "foo-id",
				},
				{
					Name: "bar",
					Id:   "bar-id",
				},
				{
					Name: "foo",
					Id:   "foo-id",
				},
			}

			artistIds := []string{
				"foo-id",
				"bar-id",
				"foo-id",
			}

			Expect(artistList.GetIds()).To(Equal(artistIds))
		})
	})
})
