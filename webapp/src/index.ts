import {Store} from 'redux';
import {PluginRegistry} from 'mattermost-webapp/plugins/registry';

import AddBookmarkModal from 'components/modals/add_bookmark';
import AddBookmarkPostMenuAction from 'components/post_menu_actions/add_bookmark';

import pluginId from 'plugin_id';

import reducer from './reducer';

export default class Plugin {
    initialize(registry: PluginRegistry, store: Store<object>) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        // eslint-disable-next-line no-unused-vars
        registry.registerReducer(reducer);
        registry.registerPostDropdownMenuComponent(AddBookmarkPostMenuAction);
        registry.registerRootComponent(AddBookmarkModal);
    }
}
window.registerPlugin(pluginId, new Plugin());
