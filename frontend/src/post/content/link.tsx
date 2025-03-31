import React from 'react';
import { ImageMetadata, LinkType } from '../post';
import { ExternalLink } from 'lucide-react';

interface LinkVizProps {
    data: LinkType;
    isFullPage: boolean;
}

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

function getLinkThumbnailPic(data: ImageMetadata[]): ImageMetadata | null {
    if (!data) {
        return null;
    }

    const viewportHeight = window.innerHeight;
    return lowerBound(data, viewportHeight, data.length - 1);
}

const LinkViz: React.FunctionComponent<LinkVizProps> = ({
    data,
    isFullPage = false
}) => {
    const thumbnail = getLinkThumbnailPic(data.image);

    const content = (
        <>
            {thumbnail && (
                <div className="sm:w-1/3 bg-slate-100 dark:bg-slate-700">
                    <div
                        className="w-full h-full min-h-[120px] bg-cover bg-center"
                        style={{ backgroundImage: `url(${thumbnail.url})` }}
                    />
                </div>
            )}

            <div
                className={`p-4 flex flex-col justify-between ${thumbnail ? 'sm:w-2/3' : 'w-full'}`}
            >
                <div>
                    <h3 className="font-medium text-slate-900 dark:text-slate-100 flex items-center gap-2">
                        {/* {data.title} */}
                        <ExternalLink className="h-4 w-4 text-slate-400" />
                    </h3>
                    {/* <p className="text-sm text-slate-500 dark:text-slate-400 mt-1 line-clamp-2">{data.description}</p> */}
                    <div className="mt-2 text-xs text-slate-400 truncate">
                        {data.link}
                    </div>
                </div>

                {/* <div className="mt-2 text-xs text-slate-400 truncate"> */}
                {/*     {data.link} */}
                {/* </div> */}
            </div>
        </>
    );

    // If we're on the full post page or inside a card that's already a link,
    // render the content without wrapping it in another link
    // if (!isFullPage) {
    //     return (
    //         <div className="border rounded-md my-2 overflow-hidden hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors">
    //             <div className="flex flex-col sm:flex-row items-stretch">
    //                 {content}
    //             </div>
    //         </div>
    //     );
    // }

    // On the full post page, make the content clickable
    return (
        <div className="border rounded-md my-2 overflow-hidden hover:bg-slate-50 dark:hover:bg-slate-700/50 transition-colors">
            <a
                href={data.link}
                target="_blank"
                rel="noopener noreferrer"
                className="flex flex-col sm:flex-row items-stretch"
            >
                {content}
            </a>
        </div>
    );
};

export default LinkViz;
