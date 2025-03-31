import React, { FunctionComponent, useRef, useState } from 'react';
import { Video } from '../post';
import { Button } from '@/components/ui/button';
import { Pause, Play, Volume2, VolumeX } from 'lucide-react';
import { Slider } from '@/components/ui/slider';

interface VideoVizProps {
    data: Video;
}

const VideoViz: FunctionComponent<VideoVizProps> = ({ data }) => {
    const [isPlaying, setIsPlaying] = useState(false);
    const [isMuted, setIsMuted] = useState(false);
    const [progress, setProgress] = useState(0);
    const videoRef = useRef<HTMLVideoElement>(null);

    const togglePlay = () => {
        if (videoRef.current) {
            if (isPlaying) {
                videoRef.current.pause();
            } else {
                videoRef.current.play();
            }
            setIsPlaying(!isPlaying);
        }
    };

    const toggleMute = () => {
        if (videoRef.current) {
            videoRef.current.muted = !isMuted;
            setIsMuted(!isMuted);
        }
    };

    const handleTimeUpdate = () => {
        if (videoRef.current) {
            const progress =
                (videoRef.current.currentTime / videoRef.current.duration) *
                100;
            setProgress(progress);
        }
    };

    const handleProgressChange = (value: number[]) => {
        if (videoRef.current) {
            const newTime = (value[0] / 100) * videoRef.current.duration;
            videoRef.current.currentTime = newTime;
            setProgress(value[0]);
        }
    };

    return (
        <div className="relative overflow-hidden rounded-md my-2">
            <div className="relative aspect-video w-full overflow-hidden bg-black">
                <video
                    ref={videoRef}
                    src={data.url}
                    autoPlay={true}
                    className="w-full h-full object-contain"
                    onTimeUpdate={handleTimeUpdate}
                    onEnded={() => setIsPlaying(false)}
                />

                <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 to-transparent p-3">
                    <div className="flex items-center gap-2">
                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-white rounded-full"
                            onClick={togglePlay}
                        >
                            {isPlaying ? (
                                <Pause className="h-4 w-4" />
                            ) : (
                                <Play className="h-4 w-4" />
                            )}
                        </Button>

                        <div className="flex-1 mx-2">
                            <Slider
                                value={[progress]}
                                min={0}
                                max={100}
                                step={0.1}
                                onValueChange={handleProgressChange}
                                className="cursor-pointer"
                            />
                        </div>

                        <Button
                            variant="ghost"
                            size="icon"
                            className="h-8 w-8 text-white rounded-full"
                            onClick={toggleMute}
                        >
                            {isMuted ? (
                                <VolumeX className="h-4 w-4" />
                            ) : (
                                <Volume2 className="h-4 w-4" />
                            )}
                        </Button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default VideoViz;
