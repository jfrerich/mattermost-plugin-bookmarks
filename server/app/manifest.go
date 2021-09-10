package app

import "github.com/mattermost/mattermost-plugin-apps/apps"

func GetManifest() apps.Manifest {
	return apps.Manifest{
		AppID:       "com.mattermost.bookmarks",
		PluginID:    "com.mattermost.bookmarks",
		DisplayName: "Bookmarks App",
		HomepageURL: "https://github.com/jfrerich/mattermost-plugin-bookmarks",
		AppType:     apps.AppTypePlugin,
		// Icon:        "bookmarks.jpg",
		RequestedPermissions: apps.Permissions{
			apps.PermissionActAsBot,
		},
		RequestedLocations: apps.Locations{
			apps.LocationCommand,
			apps.LocationPostMenu,
		},
	}
}
