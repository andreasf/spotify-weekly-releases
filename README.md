# Weekly Releases

[![Build Status](https://travis-ci.org/andreasf/spotify-weekly-releases.svg?branch=master)](https://travis-ci.org/andreasf/spotify-weekly-releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://github.com/andreasf/spotify-weekly-releases/blob/master/LICENSE)

Weekly Releases is a weekly playlist that you can subscribe to. It is created from new releases by artists that you follow on Spotify. After subscribing, check "Playlists" in your Spotify account.

# Why?

Spotify has good recommendations, but sometimes I just want to know about new releases from artists I already know. Spotify used to have built-in notifications for this purpose. The feature has mostly been removed from the clients. Currently, what's left is push notifications and emails, which don't seem to work for me.

# Running the web app

TBD

# Command line usage

1. Run `go build` in the cli subfolder
2. Head to the [Spotify Web Console](https://developer.spotify.com/web-api/console/get-users-profile/), click "Get OAuth Token" and copy the token
3. Run `./cli <your access token>`
