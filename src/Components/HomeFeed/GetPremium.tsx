import { Button, Card, CardBody, CardHeader, Image } from '@nextui-org/react';
import React from 'react';

function App() {
    return (
        <Card className="py-4">
            <CardHeader className="pb-0 pt-2 px-4 flex-col items-start">
                <h4 className="font-bold text-large">VoxPopuli Premium</h4>
                <p className="text-tiny uppercase font-bold">
                    Get the best experience with premium
                </p>
            </CardHeader>
            <CardBody className="overflow-visible py-2">
                <Image
                    alt="Card background"
                    className="object-cover rounded-xl"
                    src="https://nextui.org/images/album-cover.png"
                    width={270}
                />
            </CardBody>
            <Button className="my-3" fullWidth={true}>
                Try Now
            </Button>
        </Card>
    );
}

export default React.memo(App);
