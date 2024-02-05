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
    Image
} from '@nextui-org/react';
import React, { useState } from 'react';

interface Props {
    cname: string;
}

function App({ cname }: Props) {
    const [selectedKeys, setSelectedKeys] = useState<any>(new Set(['-1']));

    return (
        <Card className={`max-w-[400px] ${cname}`}>
            <CardHeader className="flex gap-3">
                <Image
                    alt="nextui logo"
                    height={40}
                    radius="sm"
                    src="data:image/png;base64,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"
                    width={40}
                />
                <div className="flex flex-col">
                    <p className="text-md">Today's Top Growing Communities</p>
                </div>
            </CardHeader>
            <Divider />
            <CardBody>
                <div className="flex gap-1 asidetopc">
                    <Avatar
                        isBordered
                        src="https://i.pravatar.cc/150?u=a042581f4e29026024d"
                    />

                    <p className="px-3">r/community_name</p>
                    <Button size="sm" color="primary">
                        Join
                    </Button>
                </div>
                <Divider />
                <div className="flex gap-1 asidetopc">
                    <Avatar
                        isBordered
                        src="https://i.pravatar.cc/150?u=a04258a2462d826712d"
                    />

                    <p className="px-3">r/community_name</p>
                    <Button size="sm" color="primary">
                        Join
                    </Button>
                </div>
                <Divider />
                <div className="flex gap-1 asidetopc">
                    <Avatar
                        isBordered
                        src="https://i.pravatar.cc/150?u=a04258114e29026302d"
                    />

                    <p className="px-3">r/community_name</p>
                    <Button size="sm" color="primary">
                        Join
                    </Button>
                </div>
                <Divider />
                <div className="flex gap-2 asidetopc">
                    <Avatar
                        isBordered
                        src="https://i.pravatar.cc/150?u=a04258114e29026702d"
                    />

                    <p className="px-3">r/community_name</p>
                    <Button size="sm" color="primary">
                        Join
                    </Button>
                </div>
                <Divider />
                <div className="flex gap-1 asidetopc">
                    <Avatar
                        isBordered
                        src="https://i.pravatar.cc/150?u=a04258114e29026708c"
                    />

                    <p className="px-3">r/community_name</p>
                    <Button size="sm" color="primary">
                        Join
                    </Button>
                </div>
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
                        <div className="flex gap-1 asidetopc">
                            <Avatar
                                isBordered
                                src="https://i.pravatar.cc/150?u=a04258114e29026708c"
                            />

                            <p className="px-3">r/community_name</p>
                            <Button size="sm" color="primary">
                                Join
                            </Button>
                        </div>
                        <div className="flex gap-1 asidetopc">
                            <Avatar
                                isBordered
                                src="https://i.pravatar.cc/150?u=a04258114e29026708c"
                            />

                            <p className="px-3">r/community_name</p>
                            <Button size="sm" color="primary">
                                Join
                            </Button>
                        </div>
                        <div className="flex gap-1 asidetopc">
                            <Avatar
                                isBordered
                                src="https://i.pravatar.cc/150?u=a04258114e29026708c"
                            />

                            <p className="px-3">r/community_name</p>
                            <Button size="sm" color="primary">
                                Join
                            </Button>
                        </div>
                    </AccordionItem>
                </Accordion>
                <div className="asidefooterbadges">
                    <Chip size="sm">Crypto</Chip>
                    <Chip size="sm">Books</Chip>
                    <Chip size="sm">Sports</Chip>
                    <Chip size="sm">Gaming</Chip>
                </div>
            </CardFooter>
        </Card>
    );
}

export default React.memo(App);
