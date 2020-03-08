<img src="https://github.com/jfrerich/mattermost-plugin-bookmarks/blob/master/assets/profile.png?raw=true" width="75" height="75" alt="bookmarks">

# Mattermost Bookmarks Plugin

The bookmarks plugin provides advanced options for users to bookmark posts in [Mattermost](https://mattermost.com).

Mattermost allows to users flag a post (similar to bookmarking), but you cannot arrange, group, sort, or view a condensed list of the flags. The bookmarks plugin utilizes a labeling method for bookmarking posts.  A single post can have multiple labels attached to it.

## Slash Commands

### Currently Implemented

### Future Implementations

**`/bookmark add`**
* `/bookmark add <post_id> <labels>` - bookmark a post_id with optional labels. if labels omitted, `unlabeled` autoadded
* `/bookmark add label <label>` - create a new label

**`/bookmark list`**

* `/bookmark list bookmarks <label>` - list bookmarks with optional labels for filtering
* `/bookmark list labels` - list all labels

**`/bookmark remove`**

* `/bookmark remove <post_id> <labels>` - remove labels from bookmarked post_id. if labels omitted remove post_id from bookmarks
* `/bookmark remove label <label>` - remove label from all bookmarks

**`/bookmark rename`**

* `/bookmark rename <label-old> <label-new>`- rename a label

* `/bookmark list groups`
    * list all groups


## UI Enhancements

The following UI Enhancements are planned for future release.

* post action menu
*   * `bookmark/add` (submenu) - same action as /edit but when post_id has not not been bookmarked
    * `bookmark/labels` (submenu) - shows submenus to quickly add / remove labels from current post
    * `bookmark/edit` (submenu) - open modal showing previously saved bookmark
*   * `quickmark` - quickly bookmark the current post without labels (similar to Mattermost flag option)

### Future Implementations

To learn more about plugins, see [Mattermost plugin documentation](https://developers.mattermost.com/extend/plugins/).
