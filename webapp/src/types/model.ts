
export type Bookmark = {
    postID: string;
    title: string;
    createAt: number;
    modifiedAt: number;
    labelIds: string[];
};

export type Label = {
    name: string;
    color: string;
    description: number;
};
