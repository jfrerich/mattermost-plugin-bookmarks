<img src="https://github.com/jfrerich/mattermost-plugin-bookmarks/blob/master/assets/profile.png?raw=true" width="75" height="75" alt="bookmarks">

# Mattermost Bookmarks Plugin

[![CircleCI](https://circleci.com/gh/jfrerich/mattermost-plugin-bookmarks.svg?style=shield)](https://circleci.com/gh/jfrerich/mattermost-plugin-bookmarks)
[![Go Report Card](https://goreportcard.com/badge/github.com/jfrerich/mattermost-plugin-bookmarks)](https://goreportcard.com/report/github.com/jfrerich/mattermost-plugin-bookmarks)
<!-- [![codecov](https://codecov.io/gh/jfrerich/mattermost-plugin-bookmarks/branch/master/graph/badge.svg)](https://codecov.io/gh/jfrerich/mattermost-plugin-bookmarks) -->

The bookmarks plugin provides advanced options for users to bookmark posts in [Mattermost](https://mattermost.com).

Mattermost allows users to flag a post (similar to bookmarking), but you cannot arrange, group, sort, or view a condensed list of the flags. The bookmarks plugin allows for bookmarking posts and adding personalized titles which allows the user to add context to a post message.

Addiitionally, the plugin adds slash commands which provide methods to add, view, and remove bookmarks. The `bookmarks view` command prints a condensed view of the bookmarks allowing a user to easily scan bookmark titles.


## Slash Commands

### Currently Implemented

##### Add a bookmark

```
/bookmarks add <permalink> <bookmark_title>
/bookmarks add <post_id> <bookmark_title>
    - bookmark a post by providing a post_id or the post permalink
    - optionally, provide a bookmark_title
        - if user no title is provided, the title will be the first 30 characters
          of the post message
```

##### View a bookmark

```
/bookmarks view
    - view all saved bookmark titles

/bookmarks view <permalink>
/bookmarks view <post_id>
    - Bookmarks Bot will post an ephemeral message of the post message
```

##### Remove a bookmark

```
/bookmarks remove <permalink>
/bookmarks remove <post_id>
    - remove a bookmark from your saved bookmarks
```

### ScreenShots

##### Add a bookmark

`/bookmarks add http://localhost:8065/demoteam/pl/5p4xi5hqmjddzfgggtqafk4iga ThisPostHasEmojisAndCodeBlock`
![bookmarks add post](./assets/commandAddPost.png)

##### View a bookmark

`/bookmarks view`

![bookmarks view](./assets/commandView.png)

`/bookmarks view http://localhost:8065/demoteam/pl/75ga1c6pm7n48en8sshn9bgjhy`

![bookmarks view post](./assets/commandViewWithPostID.png)

##### Remove a bookmark

`/bookmarks remove http://localhost:8065/demoteam/pl/75ga1c6pm7n48en8sshn9bgjhy`

![bookmarks remove post](./assets/commandRemovePost.png)

### Future Implementations

* `/bookmarks add <permalink> <title> <labels>` - bookmark a post with optional labels.
  * if labels omitted, `unlabeled` autoadded
* `/bookmarks label <post_id> <labels>` - add labels to a bookmark
  * if labels omitted, unlabeled autoadded
* `/bookmarks label add <labels>` - create a new label
* `/bookmarks label list` - list all labels (include number of bookmarks per label)
* `/bookmarks view <label>` - view bookmarks with optional labels for filtering
* `/bookmarks remove label <label>` - remove label from all bookmarks
* `/bookmarks rename <label-old> <label-new>`- rename a label

## UI Enhancements

The following UI Enhancements are planned for future release.

* post action menu
*   * `bookmark/add` (submenu) - same action as /edit but when post_id has not not been bookmarked
    * `bookmark/labels` (submenu) - shows submenus to quickly add / remove labels from current post
    * `bookmark/edit` (submenu) - open modal showing previously saved bookmark
*   * `quickmark` - quickly bookmark the current post without labels (similar to Mattermost flag option)

### Future Implementations

To learn more about plugins, see [Mattermost plugin documentation](https://developers.mattermost.com/extend/plugins/).
