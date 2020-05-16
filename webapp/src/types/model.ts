export type Bookmark = {
    postID: string;
    title: string;
    create_at: number;
    update_at: number;
    label_ids: string[];
};

export type Label = {
    name: string;
    color: string;
    description: number;
};

export type Labels = {
    [ByID: string]: Label;
};

export type SelectValue = {
    name: string;
    label: string;
};
