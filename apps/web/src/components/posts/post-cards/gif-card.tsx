import Image from 'next/image';
import { loadingHash } from '../constants';

interface props {
    link: (string | null)[];
    over18: boolean;
    spoiler: boolean;
}
const Gif = ({ link, over18, spoiler }: props) => {
    return (
        <div className="flex justify-center items-center relative w-full mt-3">
            {(over18 || spoiler) && (
                <div className="absolute inset-0 bg-black/70 group-hover:opacity-90 transition-opacity flex items-center justify-center">
                    <h3 className="text-white font-semibold text-lg z-50">
                        {over18 ? (
                            <h2 className="scroll-m-20 pb-2 text-3xl font-bold tracking-tight first:mt-0 text-red-500 hover:cursor-pointer">
                                NSFW
                            </h2>
                        ) : (
                            ''
                        )}
                        {spoiler ? (
                            <h2 className="scroll-m-20 pb-2 text-3xl font-semibold tracking-tight first:mt-0 text-yellow-500 hover:cursor-pointer">
                                Spoiler
                            </h2>
                        ) : (
                            ''
                        )}
                    </h3>
                </div>
            )}
            <Image
                //@ts-ignore
                src={link[0] === 'NA' ? link[1] : link[0]}
                alt="Not Found!"
                style={{ objectFit: 'cover' }}
                width={300}
                height={300}
                className={`object-cover hover:cursor-pointer duration-100 ease-in-out hover:opacity-70 rounded-md transition-all ${over18 ? 'scale-100 blur-xl grayscale' : ''}`}
                placeholder="blur"
                blurDataURL={loadingHash}
            />
        </div>
    );
};

export default Gif;
