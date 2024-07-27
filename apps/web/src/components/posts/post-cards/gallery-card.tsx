import { Card, CardContent } from '@/components/ui/card';
import {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious
} from '@/components/ui/carousel';
import Image from 'next/image';
import ImageCard from './image-card';
import { loadingHash } from '../constants';
interface props {
    links: string[][];
    over18: boolean;
    spoiler: boolean;
}

const Gallery = ({ links, over18, spoiler }: props) => {
    if (links.length === 1) {
        return <ImageCard {...{ link: links[0], over18, spoiler }} />;
    } else {
        return (
            <div className="flex justify-center items-center">
                <Carousel className="w-full max-w-md">
                    <CarouselContent>
                        {links.map((p, index) => (
                            <CarouselItem
                                key={index}
                                className="pl-1 md:basis-1/2 lg:basis-1/3"
                            >
                                <div className="p-1">
                                    <Card>
                                        <CardContent className="flex aspect-square items-center justify-center p-6">
                                            <Image
                                                src={
                                                    p[0] === 'NA' ? p[1] : p[0]
                                                }
                                                alt={p[1]}
                                                width={'300'}
                                                height={'300'}
                                                style={{ objectFit: 'cover' }}
                                                className="h-full w-full duration-3 ease-in-out hover:opacity-70 rounded-md transition-all"
                                                placeholder="blur"
                                                blurDataURL={loadingHash}
                                            />
                                        </CardContent>
                                    </Card>
                                </div>
                            </CarouselItem>
                        ))}
                    </CarouselContent>
                    <CarouselPrevious />
                    <CarouselNext />
                </Carousel>
            </div>
        );
    }
};

export default Gallery;
