// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {PureComponent} from 'react';

import {Post} from 'mattermost-redux/types/posts';

import BookmarkIcon from 'components/icon';

export type Props = {
    open: () => void;
    post: Post;
}

export default class AddBookmarkPostMenuAction extends PureComponent<Props, null> {
    handleClick = (e: React.MouseEvent) => {
        const {open, post} = this.props;
        e.preventDefault();
        open(post.id);
    }
    render() {
        const content = (
            <button
                className='style--none'
                role='presentation'
                onClick={this.handleClick}
            >
                <BookmarkIcon type='menu'/>
                {'Add Bookmark'}
            </button>
        );
        return (
            <React.Fragment>
                <li
                    className='MenuItem'
                    role='menuitem'
                >
                    {content}
                </li>
            </React.Fragment>
        );
    }
}
