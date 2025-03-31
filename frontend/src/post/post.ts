export enum MediaType {
    Image = 'image',
    Gif = 'gif',
    Video = 'video',
    Gallery = 'gallery',
    Link = 'link',
    Multi = 'multi',
    Text = 'text'
}

export interface ImageMetadata {
    id: string;
    image_id: string;
    height: number;
    width: number;
    url: string;
    created_at: Date;
    created_at_unix: number;
    updated_at: Date;
}

export interface GifMetadata {
    id: string;
    gif_id: string;
    height: number;
    width: number;
    url: string;
    created_at: Date;
    created_at_unix: number;
    updated_at: Date;
}

export interface GalleryMetadata {
    id: string;
    gallery_id: string;
    order_index: number;
    height: number;
    width: number;
    url: string;
    created_at: Date;
    created_at_unix: number;
    updated_at: Date;
}

export interface Video {
    id: string;
    media_id: string;
    url: string;
    height: number;
    width: number;
    created_at: Date;
    created_at_unix: number;
    updated_at: Date;
}

export interface LinkType {
    id: string;
    media_id: string;
    link: string;
    image: ImageMetadata[];
    created_at: Date;
    created_at_unix: number;
    updated_at: Date;
}

type PostMedia =
    | ImageMetadata
    | GifMetadata
    | GalleryMetadata
    | Video
    | LinkType;

export interface Post {
    id: string;
    author: string;
    author_id: string;
    voxsphere: string;
    voxsphere_id: string;
    title: string;
    text: string;
    text_html: string;
    media_type: MediaType;
    medias: PostMedia[];
    ups: number;
    num_comments: number;
    num_awards: number;
    over18: boolean;
    spoiler: boolean;
    created_at: Date;
    created_at_unix: number;
    updated_at: Date;
}
