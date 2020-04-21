import {GlobalState} from 'mattermost-redux/types/store';

import pluginId from './plugin_id';

const getPluginState = (state: GlobalState) => state['plugins-' + pluginId] || {};

export const addBookmarksModalState = (state: GlobalState) => getPluginState(state).addBookmarksModalVisible;

export const getAddBookmarksModalPostId = (state: GlobalState) => getPluginState(state).addBookmarkModalForPostId;
