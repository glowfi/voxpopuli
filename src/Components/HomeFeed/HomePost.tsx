import { Button, Card, CardBody, CardHeader, Image } from '@nextui-org/react';
import React from 'react';

function App() {
    return (
        <Card className="py-4">
            <CardHeader className="pb-0 pt-2 px-4 flex-col items-start">
                <h3 className="font-bold text-large">Home</h3>
                <p className="text-md uppercase">
                    Your personal VoxPopuli frontpage. Come here to check in
                    with your favorite communities.
                </p>
            </CardHeader>
            <CardBody className="overflow-visible py-2 flex items-center justify-center">
                <Image
                    alt="Card background"
                    className="object-cover rounded-xl"
                    src="https://images.fineartamerica.com/images/artworkimages/mediumlarge/3/1-vox-populi-vox-dei-vidddie-publyshd.jpg"
                    width={270}
                />
            </CardBody>
            <div className="flex flex-col gap-1">
                <Button className="my-3" fullWidth={true}>
                    Create Post
                </Button>
                <Button className="my-3" fullWidth={true}>
                    Create Community
                </Button>
            </div>
        </Card>
    );
}

export default React.memo(App);
