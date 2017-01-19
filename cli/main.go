package main

import (
	"fmt"
	"github.com/andreasf/spotify-weekly-releases/api"
	"github.com/andreasf/spotify-weekly-releases/cache"
	"github.com/andreasf/spotify-weekly-releases/model"
	"github.com/andreasf/spotify-weekly-releases/platform"
	"github.com/andreasf/spotify-weekly-releases/services"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <access token>\n", os.Args[0])
		os.Exit(1)
	}

	accessToken := os.Args[1]

	timeWrapper := &platform.TimeWrapper{}
	cache := cache.NewDiskCache("cache")
	apiClient := api.NewSpotifyApiClient("https://api.spotify.com", timeWrapper, cache)
	service := services.NewSpotifyService(apiClient)

	albums, err := service.GetRecentReleases(accessToken)
	if err != nil {
		fmt.Printf("Error retrieving followed albums: %v", err)
		os.Exit(1)
	}

	tracks := make([]model.Track, 0, len(albums))
	for _, album := range albums {
		tracks = append(tracks, *album.GetSampleTrack())
	}

	fmt.Printf("Creating a playlist from %d releases...\n", len(albums))

	date := time.Now().Format("2006-01-02")
	service.CreatePlaylist(accessToken, "Weekly Releases - "+date, tracks)
}
