package bookmarks

import (
	"fmt"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi"
)

// getPermaLink returns a link to a postID
func getPermaLink(siteURL, postID string) string {
	return fmt.Sprintf("%v/_redirect/pl/%v", siteURL, postID)
}

// getSiteURL returns the SiteURL from the config settings
func getSiteURL(api pluginapi.API) string {
	return *api.GetConfig().ServiceSettings.SiteURL
}

// getIconLink returns a markdown link to a postID including a :link: icon
func getIconLink(api pluginapi.API, postID string) string {
	url := getSiteURL(api)
	iconLink := fmt.Sprintf("[:link:](%s)", getPermaLink(url, postID))
	return iconLink
}
