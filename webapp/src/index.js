import React from 'react';

import {FormattedMessage} from 'react-intl';

import manifest from './manifest';
import {id as pluginId} from './manifest';

import {postDropdownSubMenuAction} from './actions';

export default class Plugin {
    // eslint-disable-next-line no-unused-vars
    initialize(registry, store) {
        // @see https://developers.mattermost.com/extend/plugins/webapp/reference/
        // eslint-disable-next-line no-unused-vars
        const {id, rootRegisterMenuItem} = registry.registerPostDropdownSubMenuAction(
            <FormattedMessage
                id='submenu.menu'
                key='submenu.menu'
                defaultMessage='Bookmarks'
            />
        );

        const addItem = (
            <FormattedMessage
                id='submenu.add'
                key='submenu.add'
                defaultMessage='Add'
            />
        );
        rootRegisterMenuItem(
            addItem,
            () => {
                store.dispatch(postDropdownSubMenuAction(addItem));
            }
        );
        const removeItem = (
            <FormattedMessage
                id='submenu.remove'
                key='submenu.remove'
                defaultMessage='Remove'
            />
        );
        rootRegisterMenuItem(
            removeItem,
            () => {
                store.dispatch(postDropdownSubMenuAction(removeItem));
            }
        );

        const quickMarkItem = (
            <FormattedMessage
                id='submenu.quickMark'
                key='submenu.quickMark'
                defaultMessage='QuickMark'
            />
        );
        rootRegisterMenuItem(
            quickMarkItem,
            () => {
                store.dispatch(postDropdownSubMenuAction(quickMarkItem));
            }
        );

    }
}

window.registerPlugin(manifest.id, new Plugin());
