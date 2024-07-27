export interface post {
    author: {
        username: string;
        avatarImg: string;
    };
    awards: {
        title: string;
        _count: {
            postsAwarded: number;
        };
        imageLink: string;
    }[];
    mediaContent: {
        type: string;
        id: number;
        imageLink: string | null;
        original_link: string | null;
        videoLink: string | null;
        gifLink: string | null;
        gallery: string[][];
        postId: number;
    }[];
    postflair: {
        title: string;
        colorHex: string | null;
    }[];
    voxsphere: {
        name: string;
        logoURL: string | null;
    }[];
    id: number;
    authorId: number;
    createdAt: Date;
    updatedAt: Date;
    createdAtUnix: number;
    createdatHuman: string;
    linkType: boolean;
    numComments: number;
    over18: boolean;
    spoiler: boolean;
    title: string;
    text: string;
    textHTML: string | null;
    ups: number;
}

export interface voxsphere {
    name: string;
    logoURL: string | null;
    totalmembers: number;
    totalmembersHuman: string;
}

export interface topic {
    title: string[];
}
