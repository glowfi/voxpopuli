import { Image, Link } from '@nextui-org/react';
import { useEffect, useState } from 'react';
import subredditArray from '../../data/subredditArray';
import { Subreddit } from '../../types/types';
const Banner = ({ gname }: any) => {
    const [groupData, setGroupData] = useState<Subreddit | {}>({});

    useEffect(() => {
        const fetchData = () => {
            const data = subredditArray.filter((p) => p.title == gname);
            setGroupData(() => data[0]);
        };
        fetchData();
    }, []);
    return (
        <>
            <div className="banner">
                <Image
                    isZoomed
                    style={{ display: 'block' }}
                    alt="Banner"
                    src={
                        subredditArray.filter((p) => p.title == gname)[0]
                            .bannerUrl
                    }
                />
                <div className="anchors">
                    {/* @ts-ignore */}
                    {groupData?.anchors?.map(({ title }: any, idx) => {
                        return (
                            <Link href="#" color="foreground" key={idx}>
                                {title}
                            </Link>
                        );
                    })}
                </div>
            </div>
        </>
    );
};

export default Banner;
