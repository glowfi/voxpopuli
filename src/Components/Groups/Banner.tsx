import { Button, Image, Link } from '@nextui-org/react';
import { useCallback, useEffect, useState } from 'react';
import subredditArray from '../../data/subredditArray';
import { Subreddit } from '../../types/types';
const Banner = ({ gname }: any) => {
    const [groupData, setGroupData] = useState<Subreddit | {}>({});
    const [colorstatus, setColorstatus] = useState({});
    const handleClick = useCallback(
        (args: any) => {
            //@ts-ignore
            const type = args;
            let cp = { ...colorstatus };
            for (const key in cp) {
                if (key == args) {
                    //@ts-ignore
                    cp[key] = 'danger';
                } else {
                    //@ts-ignore
                    cp[key] = 'default';
                }
            }
            console.log(cp);
            setColorstatus(() => cp);
        },
        [colorstatus]
    );

    useEffect(() => {
        const fetchData = () => {
            const data = subredditArray.filter((p) => p.title == gname);
            setGroupData(() => data[0]);
            setColorstatus(() => {
                const obj = {};
                data[0].anchors.map(({ title }, idx: any) => {
                    if (idx == 0) {
                        obj[title] = 'danger';
                    } else {
                        obj[title] = 'default';
                    }
                    return obj;
                });

                console.log(obj);

                return { ...obj };
            });
            // console.log(colorstatus);
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
                            <Button
                                variant="shadow"
                                color={colorstatus[title]}
                                onClick={() => {
                                    handleClick(title);
                                }}
                                key={idx}
                                className="btnanchor"
                            >
                                {title}
                            </Button>
                        );
                    })}
                </div>
            </div>
        </>
    );
};

export default Banner;
