// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React from 'react';
import {shallow, mount} from 'enzyme';

import AddBookmarkForm from './add_bookmark_form';

describe('components/AddBookmarkForm', () => {
    const baseActions = {
        getBookmark: jest.fn().mockResolvedValue({}),
        getAllLabels: jest.fn(),
        close: jest.fn(),
        save: jest.fn(),
    };
    const baseProps = {
        ...baseActions,
        channelId: 'channelID',
        post: {
            id: 'postID',
            message: 'This is the post message',
        },
        visible: true,
    };

    it('snapshot.. no labels or bookmark on mount', () => {
        const props = {...baseProps};
        const wrapper = mount<AddBookmarkForm>(
            <AddBookmarkForm {...props}/>,
        );
        wrapper.setState({
            submitting: false,
        });
        expect(wrapper).toMatchSnapshot();
    });

    it('snapshot.. Creatable options available', () => {
        const props = {...baseProps};
        const wrapper = mount<AddBookmarkForm>(
            <AddBookmarkForm {...props}/>,
        );
        wrapper.setState({
            submitting: false,
            allLabels: {
                ByID: {
                    lid1: {
                        name: 'lname1',
                    },
                    lid2: {
                        name: 'lname2',
                    },
                },
            },
        });
        expect(wrapper).toMatchSnapshot();
    });

    it('snapshot.. Creatable value preset in bmark', () => {
        const props = {...baseProps};
        const wrapper = mount<AddBookmarkForm>(
            <AddBookmarkForm {...props}/>,
        );
        wrapper.setState({
            submitting: false,
            selectLabelValues: [{value: 'lid1', label: 'lname1'}],
        });
        expect(wrapper).toMatchSnapshot();
    });

    it('snapshot.. title from setState bmark', () => {
        const props = {...baseProps};
        const wrapper = mount<AddBookmarkForm>(
            <AddBookmarkForm {...props}/>,
        );
        wrapper.setState({
            submitting: false,
            title: 'This bookmark has a user title',
        });
        expect(wrapper).toMatchSnapshot();
    });
});
