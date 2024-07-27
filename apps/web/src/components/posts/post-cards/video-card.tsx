import VideoPlayer from 'next-video';

interface props {
    link: (string | null)[];
    over18: boolean;
    spoiler: boolean;
}

const Video = ({ link, over18, spoiler }: props) => {
    return (
        <div className="flex justify-center items-center">
            <VideoPlayer
                //@ts-ignore
                src={link[0] === 'NA' ? link[1] : link[1]}
                autoPlay={true}
                controls={true}
            />
        </div>
    );
};

export default Video;
