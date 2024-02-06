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
                style={{ backgroundColor: groupData.headerColor }}
            >
                <h1 className="text-center text-white abouthead">
                    About Community
                </h1>
            </div>
            <Divider />
            <CardHeader className="flex gap-3">
                <Image className="imgpostg" src={groupData.logo} />
                <div className="flex flex-row gap-3 justify-center items-center flex-wrap">
                    <p className="text-md">{groupData?.title}</p>
                    <Button
                        size="sm"
                        style={{ backgroundColor: groupData.headerColor }}
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
                    Created {groupData.creationDate}
                </p>
            </div>
            <Divider />
            <div className="stats">
                <div className="members">
                    <div className="info ">
                        <IoPeopleSharp />
                        Members
                    </div>
                    {groupData.members}
                </div>
                <div className="online">
                    <div className="info">
                        <GoDotFill color={'#46d160'} />
                        <p>Online</p>
                    </div>
                    {groupData.online}
                </div>
                <div className="rank">
                    <p>Rank</p>
                    <div className="info ">
                        <FaHashtag />2
                    </div>
                </div>
            </div>
            <Divider />
            <CardBody>
                <p>{groupData?.about}</p>
            </CardBody>
            <Divider />
            <div className="flex justify-center py-4">
                <Button
                    variant="flat"
                    size="md"
                    style={{
                        backgroundColor: groupData.headerColor,
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
