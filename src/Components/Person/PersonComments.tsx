import {
    Avatar,
    Button,
    Card,
    CardBody,
    CardFooter,
    CardHeader,
    CircularProgress,
    Link
} from '@nextui-org/react';

import React, { useEffect, useRef, useState } from 'react';
import { BiDownvote, BiUpvote } from 'react-icons/bi';
import { FaAward, FaRegBookmark, FaShare } from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';
import { Bounce, ToastContainer, toast } from 'react-toastify';
import postsarr from '../../data/postsArray';

interface Props {
    cname: string;
    pname: string;
    theme: string;
}

const App = ({ cname, pname, theme }: Props) => {
    const [comments, setComments] = useState([]);
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

            let newArr = [];

            for (const items of postsarr) {
                for (const comments of items['comments']) {
                    if (comments.author == pname) {
                        newArr.push({
                            ...comments,
                            subreddit: items['subreddit'],
                            id: items['id']
                        });
                    }
                }
            }
            console.log('hello');
            console.log(newArr);

            // let newArr = postsarr.map(({ comments }) =>
            //     comments.filter(({ author }) => author == pname)
            // );
            // console.log(newArr[0]);
            setHasmore(() => false);

            // Get comments
            // let newArr = [...comments, ...postsarr.slice(idx, endidx)]
            //     .map((value) => ({ value, sort: Math.random() }))
            //     .sort((a, b) => a.sort - b.sort)
            //     .map(({ value }) => value);

            setComments(() => [...newArr]);
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
    }, [comments]);

    return (
        <>
            {comments.map((obj, idx) => {
                return (
                    <Card className={`max-w-[340px] ${cname}`} key={idx}>
                        <CardHeader className="justify-between">
                            <div className="flex gap-5">
                                <div className="flex gap-5">
                                    <Avatar
                                        isBordered
                                        radius="full"
                                        size="md"
                                        src={`/resources/images/avatar${obj.avatar}.png`}
                                    />
                                    <div className="flex flex-col gap-1 items-start justify-center">
                                        <h4 className="text-small font-semibold leading-none text-default-600 px-3">
                                            <Link
                                                href={`/p/${obj.author}`}
                                                color="foreground"
                                                underline="hover"
                                            >
                                                p/{obj.author}
                                            </Link>
                                        </h4>
                                        <h5 className="text-small tracking-tight text-default-400 px-3">
                                            <div className="time-awards flex flex-row gap-1">
                                                <p>{obj.time}</p>
                                            </div>
                                        </h5>
                                    </div>
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
                            <p>{obj.content}</p>
                        </CardBody>
                        <CardFooter className="gap-3 postfooterpe">
                            <div className="flex gap-3 postfootpe items-center justify-center">
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
                                className="flex gap-3 postfootpe items-center justify-center"
                                onClick={(e) => {
                                    e.preventDefault();
                                    console.log('shar');
                                    navigator.clipboard.writeText(
                                        `${window.location.origin}/g/${
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
                                            `/g/${obj.subreddit}/${parseInt(
                                                obj.id
                                            )}`
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
                            <div className="flex gap-3 postfootpe items-center justify-center">
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
};

export default React.memo(App);
