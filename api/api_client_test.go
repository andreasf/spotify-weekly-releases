package api_test

import (
	. "github.com/andreasf/release-feed/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/andreasf/release-feed/test_resources"
)

var _ = Describe("SpotifyApiClient", func() {
	Describe("GetFollowedArtists", func() {
		It("Makes a GET request to the endpoint", func() {
			expectedArtists := []Artist{
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
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/me/following"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, test_resources.LoadResource("../test_resources/me_following_page1.json")),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/v1/me/following?type=artist&after=bar-id&limit=20"),
					ghttp.VerifyHeaderKV("Authorization", "Bearer access-token"),
					ghttp.RespondWith(200, test_resources.LoadResource("../test_resources/me_following_page2.json")),
				),
			)

			client := NewSpotifyApiClient(server.URL())

			artists, err := client.GetFollowedArtists("access-token")
			Expect(err).To(BeNil())
			Expect(artists).To(Equal(expectedArtists))
			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})
	})

})
