// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {shallow, mount, render} from 'enzyme';

import AddBookmarkModal from './add_bookmark';

describe('components/AddBookmark', () => {
    const baseActions = {
        getBookmark: jest.fn().mockResolvedValue({
            postID: 'ID1',
            title: 'my title',
            label_ids: ['labelid1'],
        }),
        getAllLabels: jest.fn().mockResolvedValue({
            ByID: {
                labelid1: {
                    name: 'id1name',
                },
                labelid2: {
                    name: 'id1name',
                },
            },

        }),
        close: jest.fn(),
        save: jest.fn(),
    };

    const baseProps = {
        ...baseActions,
        channelId: 'channelID',
        post: {
            id: 'asdfasdf',
            message: 'This is the post message',
        },
        visible: true,
    };

    it('snapshot modal not visible', () => {
        const props = {...baseProps, visible: false};
        const wrapper = shallow(<AddBookmarkModal {...props}/>);
        expect(wrapper).toMatchSnapshot();
    });
    it('snapshot modal visible', () => {
        const props = {...baseProps};
        const wrapper = shallow(<AddBookmarkModal {...props}/>);
        expect(wrapper).toMatchSnapshot();
    });
});
