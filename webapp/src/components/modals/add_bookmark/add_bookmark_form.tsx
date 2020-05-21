// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

import React, {PureComponent} from 'react';
import {Modal} from 'react-bootstrap';

import CreatableSelect from 'react-select/creatable';

import {Post} from 'mattermost-redux/types/posts';

import {Bookmark, Labels, SelectValue} from 'types/model';

import FormButton from 'components/form_button';

export type Props = {
    getBookmark: () => void;
    getAllLabels: () => void;
    close: () => void;
    save: () => void;
    channelId: string;
    post: Post;
    visible: boolean;
}

export type State = {
    submitting: boolean;
    bookmark: Bookmark;
    allLabels: Labels;
    title: string;
    bmarkLabelIds: string;
    selectLabelValues: SelectValue[];
};

export default class AddBookmarkForm extends PureComponent<Props, State> {
    state = {
        submitting: false,
        bookmark: null,
        allLabels: null,
        title: '',
        bmarkLabelIds: '',
        selectLabelValues: [],
    };

    componentDidMount() {
        this.initializeLabelSelectValues();
    }

    initializeLabelSelectValues = async () => {
        // get all labels from labels store
        const labelsResult = await this.props.getAllLabels();
        const allLabels = labelsResult.data;
        this.setState({
            allLabels,
        });

        // get bookmark label IDs
        const postId = this.props.post.id;
        const bmarkResult = await this.props.getBookmark(postId);

        // if bmarks result is empty, the post doesn't have a bookmark saved
        if (bmarkResult === {}) {
            return;
        }

        const bmarkLabelIds = bmarkResult.data.label_ids;
        this.setState({
            bookmark: bmarkResult.data,
            title: bmarkResult.data.title,
            bmarkLabelIds,
            submitting: false,
        });

        // intialize select with labels from saved bookmark
        const initialLabelSelectValues = [];
        let value;
        for (value of bmarkLabelIds) {
            const label = allLabels.ByID[value].name;
            initialLabelSelectValues.push({value, label});
        }
        this.setState({selectLabelValues: initialLabelSelectValues});
    }

    handleClose = (e?: Event) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }
        this.props.close();
    };

    handleSubmit = (e?: Event) => {
        if (e && e.preventDefault) {
            e.preventDefault();
        }

        const labelIds = this.state.selectLabelValues.map((selectValue) => {
            return selectValue.value;
        });
        const timestamp = Date.now();
        const bookmark = {
            postid: this.props.post.id,
            title: this.state.title,
            label_ids: labelIds,
            create_at: timestamp,
            update_at: timestamp,
        };

        const currentChannelId = this.props.channelId;
        this.props.save(bookmark, currentChannelId).then((saved) => {
            if (saved.error) {
                this.setState({error: saved.error.message, submitting: false});
            }
        });
        this.props.close();
    };

    handleTitleChange = (e) => {
        this.setState({
            title: e.target.value,
        });
    }

    handleLabelChange = (e) => {
        this.setState({
            selectLabelValues: e,
        });
    }

    getLabelOptions = () => {
        if (this.state.allLabels) {
            const labelIds = Object.keys(this.state.allLabels.ByID);
            const newLabels = labelIds.map((id) => {
                const labelName = this.state.allLabels.ByID[id].name;
                return {value: id, label: labelName};
            });
            return newLabels;
        }
        return {};
    }

    render() {
        const {submitting, title, selectLabelValues} = this.state;
        const {post} = this.props;
        const style = getStyle();

        let postMessageComponent;
        if (post && post.message) {
            const message = post.message;
            postMessageComponent = (
                <div className='form-group'>
                    <label className='control-label'>{'Post Message Preview'}</label>
                    <textarea
                        style={style.textarea}
                        className='form-control'
                        value={message}
                        resize={'none'}
                        disabled={true}
                    />
                </div>
            );
        }

        const titleComponent = (
            <div className='form-group'>
                <label className='control-label'>{'Title'}</label>
                <input
                    type='text'
                    onChange={this.handleTitleChange}
                    className='form-control'
                    value={this.state.title ? this.state.title : ''}
                />
            </div>
        );

        const labelCreateComponent = (
            <div className='form-group'>
                <label className='control-label'>{'Labels'}</label>
                <CreatableSelect
                    isMulti={true}
                    options={this.getLabelOptions()}
                    onChange={this.handleLabelChange}
                    value={selectLabelValues}
                />
            </div>
        );

        return (
            <form
                role='form'
                onSubmit={() => null}
            >
                <Modal.Body ref='modalBody' >
                    {titleComponent}
                    {labelCreateComponent}
                    {postMessageComponent}
                </Modal.Body>
                <Modal.Footer >
                    <FormButton
                        type='button'
                        btnClass='btn-link'
                        defaultMessage='Cancel'
                        onClick={this.handleClose}
                    />
                    <FormButton
                        id='submit-button'
                        type='submit'
                        btnClass='btn btn-primary'
                        saving={submitting}
                        onClick={this.handleSubmit}
                        defaultMessage='Save'
                    />
                </Modal.Footer>
            </form>
        );
    }
}
const getStyle = () => ({
    textarea: {
        resize: 'none',
    },
    modal: {
        height: '100%',
    },
});
