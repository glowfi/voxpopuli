import {
    Accordion,
    AccordionItem,
    Button,
    Card,
    CardBody,
    CardFooter,
    CardHeader,
    Divider,
    Image,
    Switch
} from '@nextui-org/react';
import { useEffect, useState } from 'react';
import { FaCalendarAlt, FaHashtag } from 'react-icons/fa';
import { GoDotFill } from 'react-icons/go';
import { IoPeopleSharp } from 'react-icons/io5';
import subredditArray from '../../data/subredditArray';
import { Subreddit } from '../../types/types';

export default function App({ gname }: any) {
    const [groupData, setGroupData] = useState<Subreddit | {}>({});
    const [isSelected, setIsSelected] = useState<boolean>(true);

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
                //@ts-ignore
                style={{ backgroundColor: groupData?.headerColor }}
            >
                <h1 className="text-center text-white abouthead">
                    About Community
                </h1>
            </div>
            <Divider />
            <CardHeader className="flex gap-3">
                {/* @ts-ignore */}
                <Image className="imgpostg" src={groupData?.logo} />
                <div className="flex flex-row gap-3 justify-center items-center flex-wrap">
                    {/* @ts-ignore */}
                    <p className="text-md">{groupData?.title}</p>
                    <Button
                        size="sm"
                        //@ts-ignore
                        style={{ backgroundColor: groupData?.headerColor }}
                        variant="flat"
                        onClick={() => {}}
                    >
                        Join
                    </Button>
                </div>
            </CardHeader>
            <div className="created">
                <FaCalendarAlt />
                <p className="text-center text-md">
                    {/* @ts-ignore */}
                    Created {groupData?.creationDate}
                </p>
            </div>
            <Divider />
            <div className="stats">
                <div className="members flex flex-col justify-center items-center gap-1">
                    <div className="info ">
                        <IoPeopleSharp />
                        Members
                    </div>
                    {/* @ts-ignore */}
                    {groupData?.members}
                </div>
                <div className="online flex flex-col justify-center items-center gap-1">
                    <div className="info">
                        <GoDotFill color={'#46d160'} />
                        <p>Online</p>
                    </div>
                    {/* @ts-ignore */}
                    {groupData?.online}
                </div>
                <div className="rank flex flex-col justify-center items-center gap-1">
                    <p>Rank</p>
                    <div className="info ">
                        <FaHashtag />2
                    </div>
                </div>
            </div>
            <Divider />
            <CardBody>
                {/* @ts-ignore */}
                <p>{groupData?.about}</p>
            </CardBody>
            <Divider />
            <div className="flex justify-center py-4">
                <Button
                    variant="flat"
                    size="md"
                    style={{
                        // @ts-ignore *
                        backgroundColor: groupData?.headerColor,
                        width: '90%'
                    }}
                >
                    Create Post
                </Button>
            </div>
            <Divider />
            <CardFooter>
                <Accordion variant="shadow">
                    <AccordionItem
                        key="1"
                        aria-label="Accordion 1"
                        title="Community Theme"
                    >
                        <div className="flex flex-col gap-2">
                            <Switch
                                isSelected={isSelected}
                                onValueChange={setIsSelected}
                            >
                                Community Theme
                            </Switch>
                            <p className="text-small text-default-500">
                                Selected: {isSelected ? 'true' : 'false'}
                            </p>
                        </div>
                    </AccordionItem>
                </Accordion>
            </CardFooter>
        </Card>
    );
}
