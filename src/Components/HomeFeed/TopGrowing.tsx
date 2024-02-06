import {
    Accordion,
    AccordionItem,
    Avatar,
    Button,
    Card,
    CardBody,
    CardFooter,
    CardHeader,
    Chip,
    Divider,
    Link
} from '@nextui-org/react';
import React, { useState } from 'react';

import { FaPeopleGroup } from 'react-icons/fa6';
import { useNavigate } from 'react-router-dom';
import subredditArray from '../../data/subredditArray';

interface Props {
    cname: string;
}

function App({ cname }: Props) {
    const [selectedKeys, setSelectedKeys] = useState<any>(new Set(['-1']));

    //@ts-ignore
    const [top5, setTop5] = useState(subredditArray.slice(0, 5));
    //@ts-ignore
    const [rest5, setRest5] = useState(subredditArray.slice(5));

    const navigate = useNavigate();

    return (
        <Card className={`max-w-[400px] ${cname} w-full`}>
            <CardHeader className="flex gap-3">
                <FaPeopleGroup size={'1.8rem'} />
                <div className="flex flex-col">
                    <p className="text-md">Today's Top Growing Communities</p>
                </div>
            </CardHeader>
            <Divider />
            <CardBody>
                {top5.map((obj, idx) => {
                    return (
                        <div key={idx}>
                            <div className="flex gap-1 asidetopc justify-center items-center">
                                <Avatar isBordered src={obj.logo} />

                                <p className="px-3">
                                    <Link
                                        href={`/g/${obj.title}`}
                                        color="foreground"
                                        underline="hover"
                                    >
                                        g/{obj.title}
                                    </Link>
                                </p>
                                <Button
                                    size="sm"
                                    color="primary"
                                    onClick={() => {
                                        navigate(`/g/${obj.title}`);
                                    }}
                                >
                                    Join
                                </Button>
                            </div>
                            {idx == 4 ? '' : <Divider />}
                        </div>
                    );
                })}
            </CardBody>
            <Divider />
            <CardFooter className={'asidefootermain'}>
                <Accordion
                    selectedKeys={selectedKeys}
                    onSelectionChange={setSelectedKeys}
                >
                    <AccordionItem
                        key="1"
                        aria-label="Accordion 1"
                        title="View All Communities"
                        className={'text-center'}
                    >
                        {rest5.map((obj, idx) => {
                            return (
                                <div key={idx}>
                                    <div className="flex gap-1 asidetopc justify-center">
                                        <Avatar isBordered src={obj.logo} />

                                        <p className="px-3">
                                            <Link
                                                href={`/g/${obj.title}`}
                                                color="foreground"
                                                underline="hover"
                                            >
                                                g/{obj.title}
                                            </Link>
                                        </p>
                                        <Button
                                            size="sm"
                                            color="primary"
                                            onClick={() => {
                                                navigate(`/g/${obj.title}`);
                                            }}
                                        >
                                            Join
                                        </Button>
                                    </div>
                                    <Divider />
                                </div>
                            );
                        })}
                    </AccordionItem>
                </Accordion>
                <div className="asidefooterbadges">
                    <Chip className="asidebadges" size="sm">
                        Crypto
                    </Chip>
                    <Chip className="asidebadges" size="sm">
                        Books
                    </Chip>
                    <Chip className="asidebadges" size="sm">
                        Sports
                    </Chip>
                    <Chip className="asidebadges" size="sm">
                        Gaming
                    </Chip>
                </div>
            </CardFooter>
        </Card>
    );
}

export default React.memo(App);
