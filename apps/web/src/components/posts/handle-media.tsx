import Gallery from './post-cards/gallery-card';
import Gif from './post-cards/gif-card';
import ImageCard from './post-cards/image-card';
import Video from './post-cards/video-card';

interface mediadata {
    type: string;
    id: number;
    original_link: string | null;
    videoLink: string | null;
    imageLink: string | null;
    gifLink: string | null;
    gallery: string[][];
    postId: number;
}

interface props {
    mediadata: mediadata;
    over18: boolean;
    spoiler: boolean;
}

const HandleMedia = ({ mediadata, over18, spoiler }: props) => {
    if (mediadata?.type === 'image') {
        let data = [mediadata.imageLink, mediadata.original_link];
        return <ImageCard link={data} over18={over18} spoiler={spoiler} />;
    } else if (mediadata?.type === 'gif') {
        let data = [mediadata.gifLink, mediadata.original_link];
        return <Gif link={data} over18={over18} spoiler={spoiler} />;
    } else if (mediadata?.type === 'gallery') {
        return (
            <Gallery
                links={mediadata.gallery}
                over18={over18}
                spoiler={spoiler}
            />
        );
    } else if (mediadata?.type === 'video') {
        let data = [mediadata.videoLink, mediadata.original_link];

        return <Video link={data} over18={over18} spoiler={spoiler} />;
    }
    return <></>;
};

export default HandleMedia;
