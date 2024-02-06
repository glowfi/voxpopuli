import { Card, CardFooter, Chip } from '@nextui-org/react';
import { useEffect, useState } from 'react';
import subredditArray from '../../data/subredditArray';
import { Subreddit } from '../../types/types';

const Flairs = ({ gname }: any) => {
    const [groupData, setGroupData] = useState<Subreddit | {}>({});

    useEffect(() => {
        const fetchData = () => {
            const data = subredditArray.filter((p) => p.title == gname);
            setGroupData(() => data[0]);
        };
        fetchData();
    }, []);

    return (
        <Card className="max-w-[400px] w-full aside-item">
            <div
                className="aboutcl"
                style={{ backgroundColor: groupData?.headerColor }}
            >
                <h1 className="text-center text-white abouthead">
                    Search By falirs
                </h1>
            </div>
            <CardFooter>
                <div className="flairs flex flex-row gap-2 flex-wrap">
                    {groupData?.flairs?.map(({ title, color }, idx) => {
                        return (
                            <Chip
                                key={idx}
                                variant="shadow"
                                classNames={{
                                    base: 'bg-gradient-to-br from-indigo-500 to-pink-500 border-small border-white/50 shadow-pink-500/30',
                                    content:
                                        'drop-shadow shadow-black text-white'
                                }}
                                style={{ backgroundColor: color }}
                            >
                                {title}
                            </Chip>
                        );
                    })}
                </div>
            </CardFooter>
        </Card>
    );
};

export default Flairs;
