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
    Link
} from '@nextui-org/react';
import { useEffect, useState } from 'react';
import { FaBirthdayCake, FaCalendarAlt } from 'react-icons/fa';
import { GiCrystalGrowth } from 'react-icons/gi';
import { IoMdAdd } from 'react-icons/io';
import { TbReport } from 'react-icons/tb';
import userarray from '../../data/userArray';

export default function App({ pname }) {
    const [userData, setUserData] = useState({});
    const [isSelected, setIsSelected] = useState<boolean>(true);

    useEffect(() => {
        const fetchData = () => {
            console.log(userarray);
            const data = userarray.filter((p) => p.username == pname);
            console.log(data[0]);
            setUserData(() => data[0]);
        };
        fetchData();
    }, []);

    return (
        <Card className="max-w-[400px] w-full aside-item">
            <div
                className="aboutcl"
                //@ts-ignore
                style={{ backgroundColor: userData?.color }}
            >
                <h1 className="text-center text-white abouthead ">
                    <b>p/{pname}</b>
                </h1>
            </div>
            <Divider />
            <CardHeader className="flex gap-3 flex-col">
                {/* @ts-ignore */}
                <Image
                    className="imgabout"
                    src={`/resources/images/avatar${userData?.avatar}_head.png`}
                />
                <div className="flex flex-col gap-1 justify-center items-center flex-wrap">
                    {/* @ts-ignore */}
                    <p className="text-md">{userData?.username}</p>
                    <div className="created flex flex-row gap-1">
                        <FaCalendarAlt />
                        <p className="text-center text-md">
                            {/* @ts-ignore */}
                            {userData?.age}
                        </p>
                    </div>
                </div>
                <div className="moreinfo">
                    <div className="created flex flex-col gap-1">
                        <div className="flex gap-1 justify-center items-center">
                            <FaBirthdayCake />
                            <p>Cakeday</p>
                        </div>

                        <p className="text-center text-sm">
                            {/* @ts-ignore */}
                            {userData?.cakeday}
                        </p>
                    </div>
                    <div className="created flex flex-col gap-1">
                        <div className="flex gap-1 justify-center items-center">
                            <GiCrystalGrowth />
                            <p>Karma</p>
                        </div>
                        <p className="text-center text-md">
                            {/* @ts-ignore */}
                            {userData?.karma}
                        </p>
                    </div>
                </div>
            </CardHeader>
            <Divider />
            <CardBody>
                <div className="buttons flex flex-wrap gap-3 justify-center items-center">
                    <Button color="default">
                        <TbReport />
                        Report User
                    </Button>
                    <Button color="default">
                        <IoMdAdd />
                        Add as a friend
                    </Button>
                    <Button color="primary" radius="lg" fullWidth>
                        Follow
                    </Button>
                    <Button color="primary" radius="lg" fullWidth>
                        Chat
                    </Button>
                </div>
            </CardBody>
            <Divider />
            <CardFooter>
                <Accordion variant="shadow">
                    <AccordionItem
                        key="1"
                        aria-label="Accordion 1"
                        title="More Options"
                    >
                        <div className="flex flex-wrap flex-col justify-center items-center gap-3">
                            <Link
                                isExternal
                                href="https://github.com/nextui-org/nextui"
                            >
                                Get them help and support
                            </Link>
                            <Link
                                isExternal
                                href="https://github.com/nextui-org/nextui"
                            >
                                Add to custom feed
                            </Link>
                            <Link
                                isExternal
                                href="https://github.com/nextui-org/nextui"
                            >
                                Invite someone to chat
                            </Link>
                        </div>
                    </AccordionItem>
                </Accordion>
            </CardFooter>
        </Card>
    );
}
