import {
    Avatar,
    Button,
    Card,
    CardBody,
    CardFooter,
    CardHeader,
    Chip,
    CircularProgress,
    Image,
    Link
} from '@nextui-org/react';

import React, { useEffect, useRef, useState } from 'react';
import { BiDownvote, BiUpvote } from 'react-icons/bi';
import { FaAward, FaComment, FaRegBookmark, FaShare } from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';
import { Bounce, ToastContainer, toast } from 'react-toastify';
import postsarr from '../../data/postsArray';
import { Post } from '../../types/types';

interface Props {
    cname: string;
    gname: string;
    theme: string;
}

function App({ cname, gname, theme }: Props) {
    const [posts, setPosts] = useState<[] | Post[]>([]);
    const [hasmore, setHasmore] = useState<Boolean>(true);
    // const [idx, setIdx] = useState(0);
    const ref = useRef<HTMLElement | any>(null);
    const navigate = useNavigate();

    const onIntersection = (entries: any) => {
        const firstEntry = entries[0];
        // const offset = 10;
        if (firstEntry.isIntersecting && hasmore) {
            // let endidx = idx + offset + 1;

            // // Handle limit
            // if (endidx > postsarr.length) {
            //     endidx = postsarr.length;
            //     setHasmore(() => false);
            //     setIdx(postsarr.length);
            // } else {
            //     setHasmore(true);
            //     setIdx(idx + offset + 1);
            // }

            let newArr = postsarr.filter((p) => p.subreddit == gname);
            setHasmore(() => false);

            // Get Posts
            // let newArr = [...posts, ...postsarr.slice(idx, endidx)]
            //     .map((value) => ({ value, sort: Math.random() }))
            //     .sort((a, b) => a.sort - b.sort)
            //     .map(({ value }) => value);

            setPosts(() => [...newArr]);
        }
    };

    useEffect(() => {
        const observer = new IntersectionObserver(onIntersection);
        if (observer && ref.current) {
            observer.observe(ref.current);
        }

        return () => {
            if (observer) {
                observer.disconnect();
            }
        };
    }, [posts]);

    return (
        <>
            {posts.map((obj: Post, idx) => {
                return (
                    <Card className={`max-w-[340px] ${cname}`} key={idx}>
                        <CardHeader
                            className="justify-between"
                            onClick={() => {
                                console.log('hello');
                                navigate(`/g/${obj.subreddit}/${obj.id}`);
                            }}
                        >
                            <div className="flex gap-5">
                                <Avatar
                                    isBordered
                                    radius="full"
                                    size="md"
                                    src={`/resources/images/Communities/${obj.subreddit}/icon.png`}
                                />
                                <div className="flex flex-col gap-1 items-start justify-center">
                                    <h4 className="text-small font-semibold leading-none text-default-600 px-3">
                                        <Link
                                            href={`/g/${obj.subreddit}`}
                                            color="foreground"
                                            underline="hover"
                                        >
                                            g/{obj.subreddit}
                                        </Link>
                                    </h4>
                                    <h5 className="text-small tracking-tight text-default-400 px-3">
                                        <Link
                                            href={`/p/${obj.author}`}
                                            color="foreground"
                                            underline="hover"
                                        >
                                            Posted by p/{obj.author}
                                        </Link>
                                        <div className="time-awards flex flex-row gap-1">
                                            <p>{obj.time}</p>
                                            {obj.awards.map((name, idx) => {
                                                return (
                                                    <Image
                                                        className="imgpost"
                                                        src={`/resources/images/${name}.png`}
                                                        key={idx}
                                                    />
                                                );
                                            })}
                                        </div>
                                    </h5>
                                </div>
                            </div>
                            <div className="flex gap-3 upvotebtng">
                                <Button
                                    size="sm"
                                    radius="sm"
                                    color="default"
                                    variant="flat"
                                >
                                    <BiUpvote size={'1rem'} />
                                </Button>
                                <div>
                                    <p className="text-medium bg-red">
                                        {obj.upvotes}
                                    </p>
                                </div>
                                <Button
                                    size="sm"
                                    radius="sm"
                                    color="default"
                                    variant="flat"
                                >
                                    <BiDownvote size={'1rem'} />
                                </Button>
                            </div>
                        </CardHeader>
                        <CardBody className="px-3 py-0 text-small text-default-400">
                            <div
                                className="postbodyg"
                                onClick={() => {
                                    console.log('hello');
                                    navigate(`/g/${obj.subreddit}/${obj.id}`);
                                }}
                            >
                                <div className="flex flex-row gap-3">
                                    <p className="text-xl postheadingg">
                                        {obj.title}
                                    </p>
                                    <Chip
                                        className="text-default-400 font-semibold flairpostg"
                                        variant="shadow"
                                        style={{
                                            backgroundColor: obj.flair.color,
                                            color: 'white'
                                        }}
                                    >
                                        {obj.flair.title}
                                    </Chip>
                                </div>
                                {obj.src.includes('.png') ||
                                obj.src.includes('.PNG') ||
                                obj.src.includes('.jpg') ||
                                obj.src.includes('.jpeg') ? (
                                    <Image
                                        width={600}
                                        height={350}
                                        alt="NextUI hero Image with delay"
                                        src={obj.src.toString()}
                                    />
                                ) : (
                                    <p className="postfootertxtg">
                                        {obj.src.slice(25)} ...
                                    </p>
                                )}
                            </div>
                        </CardBody>
                        <CardFooter className="gap-3 postfooterg">
                            <div className="flex gap-3 postfootg items-center justify-center mt-3">
                                <p className="font-semibold text-default-400 text-small">
                                    <FaComment size={'1.3rem'} />
                                </p>
                                <Link
                                    href={`/g/${obj.subreddit}/${parseInt(
                                        obj.id
                                    )}`}
                                    color="foreground"
                                    underline="hover"
                                >
                                    {obj.comments.length} comments
                                </Link>
                            </div>
                            <div className="flex gap-3 postfootg items-center justify-center">
                                <p className="font-semibold text-default-400 text-small">
                                    <FaAward size={'1.3rem'} />
                                </p>
                                <Link
                                    href={`/g/${obj.subreddit}/${parseInt(
                                        obj.id
                                    )}`}
                                    color="foreground"
                                    underline="hover"
                                >
                                    Award
                                </Link>
                            </div>

                            <div
                                className="flex gap-3 postfootg items-center justify-center"
                                onClick={(e) => {
                                    e.preventDefault();
                                    navigator.clipboard.writeText(
                                        `${window.location.href}g/${
                                            obj.subreddit
                                        }/${parseInt(obj.id)}`
                                    );
                                    toast.success('Link Copied!');
                                }}
                            >
                                <p className="font-semibold text-default-400 text-small">
                                    <FaShare size={'1.3rem'} />
                                </p>
                                <Link
                                    href={`/g/${obj.subreddit}/${parseInt(
                                        obj.id
                                    )}`}
                                    color="foreground"
                                    underline="hover"
                                    onClick={(e) => {
                                        e.preventDefault();
                                        navigator.clipboard.writeText(
                                            `${window.location.href}g/${
                                                obj.subreddit
                                            }/${parseInt(obj.id)}`
                                        );
                                        toast.success('Link Copied!');
                                    }}
                                >
                                    Share
                                </Link>
                                <ToastContainer
                                    position="bottom-right"
                                    autoClose={5000}
                                    hideProgressBar={false}
                                    newestOnTop={false}
                                    closeOnClick
                                    rtl={false}
                                    pauseOnFocusLoss
                                    draggable
                                    pauseOnHover
                                    theme={theme}
                                    transition={Bounce}
                                />
                            </div>
                            <div className="flex gap-3 postfootg items-center justify-center">
                                <p className="font-semibold text-default-400 text-small">
                                    <FaRegBookmark size={'1.3rem'} />
                                </p>
                                <p className="text-default-400 text-small">
                                    <Link
                                        href={`/g/${obj.subreddit}/${parseInt(
                                            obj.id
                                        )}`}
                                        color="foreground"
                                        underline="hover"
                                    >
                                        Save
                                    </Link>
                                </p>
                            </div>
                        </CardFooter>
                    </Card>
                );
            })}
            {hasmore ? (
                <div className="item-4" ref={ref}>
                    <CircularProgress label="Loading..." />
                </div>
            ) : (
                ''
            )}
        </>
    );
}

export default React.memo(App);
