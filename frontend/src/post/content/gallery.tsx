'use client';
import { Button } from '@/components/ui/button';
import { cn } from '@/lib/utils';
import {
    ArrowLeft,
    ArrowRight,
    ChevronLeft,
    ChevronRight,
    Maximize2,
    X
} from 'lucide-react';
import Image from 'next/image';
import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { GalleryMetadata } from '../post';

interface Gallery {
    index: number;
    images: GalleryMetadata[];
}

function buildGallery(data: GalleryMetadata[]): Gallery[] {
    const final: Gallery[] = [];
    const obj: { [key: number]: GalleryMetadata[] } = {};

    for (let index = 0; index < data.length; index++) {
        const groupId = data[index].order_index;
        if (groupId in obj) {
            obj[groupId].push(data[index]);
        } else {
            obj[groupId] = [data[index]];
        }
    }

    for (const key in obj) {
        if (obj.hasOwnProperty(key)) {
            final.push({ index: parseInt(key), images: obj[key] });
        }
    }

    return final;
}

function lowerBound(
    arr: GalleryMetadata[],
    y: number,
    n: number
): GalleryMetadata {
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

function getGalleryPics(data: GalleryMetadata[]): GalleryMetadata[] {
    const galleryImages: Gallery[] = buildGallery(data);

    const viewportHeight = window.innerHeight;

    const final: GalleryMetadata[] = [];

    for (let index = 0; index < galleryImages.length; index++) {
        const gallery = galleryImages[index];
        const galleryMetadata = lowerBound(
            gallery.images,
            viewportHeight,
            gallery.images.length - 1
        );
        final.push(galleryMetadata);
    }

    return final;
}

interface GalleryProps {
    data: GalleryMetadata[];
}

const Gallery: React.FC<GalleryProps> = ({ data }) => {
    const pics: GalleryMetadata[] = useMemo(() => {
        return getGalleryPics(data);
    }, [data]);

    const [currentIndex, setCurrentIndex] = useState(0);
    const [isFullscreen, setIsFullscreen] = useState(false);
    const [isLoading, setIsLoading] = useState(true);
    const [touchStart, setTouchStart] = useState(0);
    const [touchEnd, setTouchEnd] = useState(0);

    const totalImages = pics.length;

    const goToNext = useCallback(() => {
        setIsLoading(true);
        setCurrentIndex((prevIndex) =>
            prevIndex === totalImages - 1 ? 0 : prevIndex + 1
        );
    }, [totalImages]);

    const goToPrevious = useCallback(() => {
        setIsLoading(true);
        setCurrentIndex((prevIndex) =>
            prevIndex === 0 ? totalImages - 1 : prevIndex - 1
        );
    }, [totalImages]);

    const goToSlide = (index: number) => {
        setIsLoading(true);
        setCurrentIndex(index);
    };

    const toggleFullscreen = () => {
        setIsFullscreen(!isFullscreen);
    };

    const handleImageLoad = () => {
        setIsLoading(false);
    };

    // Handle keyboard navigation
    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (isFullscreen) {
                if (e.key === 'ArrowRight') goToNext();
                if (e.key === 'ArrowLeft') goToPrevious();
                if (e.key === 'Escape') setIsFullscreen(false);
            }
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [isFullscreen, goToNext, goToPrevious]);

    // Handle touch events for swipe
    const handleTouchStart = (e: React.TouchEvent) => {
        setTouchStart(e.targetTouches[0].clientX);
    };

    const handleTouchMove = (e: React.TouchEvent) => {
        setTouchEnd(e.targetTouches[0].clientX);
    };

    const handleTouchEnd = () => {
        if (touchStart - touchEnd > 50) {
            // Swipe left
            goToNext();
        }

        if (touchStart - touchEnd < -50) {
            // Swipe right
            goToPrevious();
        }
    };

    // Fullscreen gallery modal
    const FullscreenGallery = () => (
        <div className="fixed inset-0 z-50 bg-black flex items-center justify-center">
            <div className="absolute top-4 right-4 z-10 flex gap-2">
                <Button
                    variant="ghost"
                    size="icon"
                    className="rounded-full bg-black/50 text-white hover:bg-black/70"
                    onClick={toggleFullscreen}
                >
                    <X className="h-5 w-5" />
                </Button>
            </div>

            <div
                className="relative w-full h-full flex items-center justify-center"
                onTouchStart={handleTouchStart}
                onTouchMove={handleTouchMove}
                onTouchEnd={handleTouchEnd}
            >
                {isLoading && (
                    <div className="absolute inset-0 flex items-center justify-center">
                        <div className="w-8 h-8 border-4 border-t-primary rounded-full animate-spin"></div>
                    </div>
                )}

                <Image
                    src={pics[currentIndex].url || '/placeholder.svg'}
                    alt={'Not Found'}
                    fill
                    className={cn(
                        'object-contain transition-opacity duration-300',
                        isLoading ? 'opacity-0' : 'opacity-100'
                    )}
                    onLoad={handleImageLoad}
                    priority
                />

                <Button
                    variant="ghost"
                    size="icon"
                    className="absolute left-4 rounded-full bg-black/50 text-white hover:bg-black/70"
                    onClick={goToPrevious}
                >
                    <ArrowLeft className="h-6 w-6" />
                </Button>

                <Button
                    variant="ghost"
                    size="icon"
                    className="absolute right-4 rounded-full bg-black/50 text-white hover:bg-black/70"
                    onClick={goToNext}
                >
                    <ArrowRight className="h-6 w-6" />
                </Button>

                <div className="absolute bottom-4 left-1/2 -translate-x-1/2 bg-black/50 text-white px-3 py-1 rounded-full text-sm">
                    {currentIndex + 1} / {totalImages}
                </div>
            </div>
        </div>
    );

    return (
        <div className="my-2">
            {/* Regular gallery view */}
            <div className="relative rounded-md overflow-hidden bg-slate-100 dark:bg-slate-800">
                <div
                    className="relative aspect-[16/9] w-full"
                    onTouchStart={handleTouchStart}
                    onTouchMove={handleTouchMove}
                    onTouchEnd={handleTouchEnd}
                >
                    {isLoading && (
                        <div className="absolute inset-0 flex items-center justify-center">
                            <div className="w-6 h-6 border-3 border-t-primary rounded-full animate-spin"></div>
                        </div>
                    )}

                    <div className="absolute inset-0 flex items-center justify-center">
                        <Image
                            src={pics[currentIndex].url || '/placeholder.svg'}
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

                    {/* Navigation arrows */}
                    {totalImages > 1 && (
                        <>
                            <Button
                                variant="ghost"
                                size="icon"
                                className="absolute left-2 top-1/2 -translate-y-1/2 rounded-full bg-black/30 hover:bg-black/50 text-white border-none h-8 w-8"
                                onClick={goToPrevious}
                            >
                                <ChevronLeft className="h-5 w-5" />
                            </Button>

                            <Button
                                variant="ghost"
                                size="icon"
                                className="absolute right-2 top-1/2 -translate-y-1/2 rounded-full bg-black/30 hover:bg-black/50 text-white border-none h-8 w-8"
                                onClick={goToNext}
                            >
                                <ChevronRight className="h-5 w-5" />
                            </Button>
                        </>
                    )}

                    {/* Fullscreen button */}
                    <Button
                        variant="ghost"
                        size="icon"
                        className="absolute top-2 right-2 rounded-full bg-black/30 hover:bg-black/50 text-white border-none h-8 w-8"
                        onClick={toggleFullscreen}
                    >
                        <Maximize2 className="h-4 w-4" />
                    </Button>

                    {/* Image counter */}
                    {totalImages > 1 && (
                        <div className="absolute top-2 left-2 bg-black/50 text-white px-2 py-1 rounded-md text-xs">
                            {currentIndex + 1} / {totalImages}
                        </div>
                    )}
                </div>

                {/* Thumbnail indicators */}
                {totalImages > 1 && (
                    <div className="flex justify-center gap-1.5 p-2 bg-slate-200 dark:bg-slate-700">
                        {pics.map((_, index) => (
                            <button
                                key={index}
                                className={cn(
                                    'w-2 h-2 rounded-full transition-colors',
                                    index === currentIndex
                                        ? 'bg-primary'
                                        : 'bg-slate-400 dark:bg-slate-500 hover:bg-slate-500 dark:hover:bg-slate-400'
                                )}
                                onClick={() => goToSlide(index)}
                                aria-label={`Go to slide ${index + 1}`}
                            />
                        ))}
                    </div>
                )}
            </div>

            {/* Fullscreen modal */}
            {isFullscreen && <FullscreenGallery />}
        </div>
    );
};

export default Gallery;
