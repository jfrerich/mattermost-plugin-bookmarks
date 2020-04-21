
export type Bookmark = {
    postID: string;
    title: string;
    createAt: number;
    modifiedAt: number;
    labelIDs: string[];
};

export type Label = {
    name: string;
    color: string;
    description: number;
};
