// Copyright (c) 2017-present Mattermost, Inc. All Rights Reserved.
// See License for license information.

import React from 'react';

type Props = {
    type?: string;
};

export default class BookmarkIcon extends React.PureComponent<Props> {
    public render() {
        let iconStyle = {};
        if (this.props.type === 'menu') {
            iconStyle = {flex: '0 0 auto', width: '20px', height: '20px', fill: '#0052CC', background: 'white', borderRadius: '50px', padding: '2px'};
        }

        return (
            <span className='MenuItem__icon'>
                <svg
                    aria-hidden='true'
                    focusable='false'
                    role='img'
                    viewBox='0 0 24 24'
                    width='14'
                    height='14'
                    style={iconStyle}
                >
                    <path d='M15 1H5a2 2 0 0 0-2 2v16l7-5 7 5V3a2 2 0 0 0-2-2zm0 14.25l-5-3.5-5 3.5V3h10z'/>
                </svg>
            </span>
        );
    }
}
