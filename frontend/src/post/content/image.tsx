import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import { X, Maximize2 } from 'lucide-react';
import React, { useState } from 'react';
import { ImageMetadata } from '../post';
import Image from 'next/image';

function lowerBound(arr: ImageMetadata[], y: number, n: number): ImageMetadata {
    let st = 0;
    let en = n - 1;
    let ans = n;

    while (st <= en) {
        const mid = Math.floor(st + (en - st) / 2);

        if (arr[mid].height >= y) {
            ans = mid;
            en = mid - 1;
        } else {
            st = mid + 1;
        }
    }

    return arr[ans];
}

function getImagePic(data: ImageMetadata[]): ImageMetadata {
    const viewportHeight = window.innerHeight;
    return lowerBound(data, viewportHeight, data.length - 1);
}

interface ImageVizProps {
    data: ImageMetadata[];
}

const ImageViz: React.FunctionComponent<ImageProps> = ({ data }) => {
    const bestImage: ImageMetadata = getImagePic(data);

    const [isFullscreen, setIsFullscreen] = useState(false);
    const [isLoading, setIsLoading] = useState(true);

    const toggleFullscreen = () => {
        setIsFullscreen(!isFullscreen);
    };

    const handleImageLoad = () => {
        setIsLoading(false);
    };

    // Fullscreen image modal
    const FullscreenImage = () => (
        <div className="fixed inset-0 z-50 bg-black flex items-center justify-center">
            <div className="absolute top-4 right-4 z-10">
                <Button
                    variant="ghost"
                    size="icon"
                    className="rounded-full bg-black/50 text-white hover:bg-black/70"
                    onClick={toggleFullscreen}
                >
                    <X className="h-5 w-5" />
                </Button>
            </div>

            <div className="relative w-full h-full flex items-center justify-center">
                {isLoading && (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-8 h-8 border-4 border-t-primary rounded-full animate-spin"></div>
                    </div>
                )}

                <Image
                    src={bestImage.url || '/placeholder.svg'}
                    alt={'Not Found'}
                    fill
                    className={cn(
                        'object-contain transition-opacity duration-300',
                        isLoading ? 'opacity-0' : 'opacity-100'
                    )}
                    onLoad={handleImageLoad}
                    priority
                />
            </div>
        </div>
    );

    return (
        <div className="relative overflow-hidden rounded-md my-2">
            <div className="relative aspect-[16/9] w-full">
                {isLoading && (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-6 h-6 border-3 border-t-primary rounded-full animate-spin"></div>
                    </div>
                )}

                <div className="absolute inset-0 flex items-center justify-center">
                    <Image
                        src={bestImage.url || '/placeholder.svg'}
                        alt={'Not Found'}
                        fill
                        className={cn(
                            'object-contain transition-opacity duration-300',
                            isLoading ? 'opacity-0' : 'opacity-100'
                        )}
                        onLoad={handleImageLoad}
                    />
                </div>

                {/* Fullscreen button */}
                <Button
                    variant="ghost"
                    size="icon"
                    className="absolute top-2 right-2 rounded-full bg-black/30 hover:bg-black/50 text-white border-none h-8 w-8"
                    onClick={toggleFullscreen}
                >
                    <Maximize2 className="h-4 w-4" />
                </Button>
            </div>

            {/* Fullscreen modal */}
            {isFullscreen && <FullscreenImage />}
        </div>
    );
};

export default ImageViz;
