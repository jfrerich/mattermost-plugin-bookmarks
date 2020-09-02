package bookmarks

import (
	"fmt"

	"github.com/jfrerich/mattermost-plugin-bookmarks/server/pluginapi"
	"github.com/jfrerich/mattermost-plugin-bookmarks/server/utils"
)

// getPermaLink returns a link to a postID
func getPermaLink(siteURL, postID string) string {
	return fmt.Sprintf("%v/_redirect/pl/%v", siteURL, postID)
}

// getIconLink returns a markdown link to a postID including a :link: icon
func getIconLink(api pluginapi.API, postID string) string {
	url := utils.GetSiteURL(api)
	iconLink := fmt.Sprintf("[:link:](%s)", getPermaLink(url, postID))
	return iconLink
}
