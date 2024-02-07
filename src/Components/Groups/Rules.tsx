import { Accordion, AccordionItem, Card, Divider } from '@nextui-org/react';
import { useEffect, useState } from 'react';
import subredditArray from '../../data/subredditArray';
import { Subreddit } from '../../types/types';

export default function App({ gname }: any) {
    const [groupData, setGroupData] = useState<Subreddit | {}>({});

    useEffect(() => {
        const fetchData = () => {
            const data = subredditArray.filter((p) => p.title == gname);
            setGroupData(() => data[0]);
        };
        fetchData();
    }, []);
    return (
        <Card className="max-w-[400px] w-full aside-item">
            <div
                className="aboutcl"
                //@ts-ignore
                style={{ backgroundColor: groupData?.headerColor }}
            >
                <h1 className="text-center text-white abouthead">
                    {/* @ts-ignore */}
                    g/{groupData?.title} Rules
                </h1>
            </div>
            <Divider />
            <div className="accord">
                <Accordion variant="shadow">
                    {/* @ts-ignore  */}
                    {groupData?.rules?.map(
                        ({ number, title, desc }: any, idx: any) => {
                            return (
                                <AccordionItem
                                    key={idx}
                                    aria-label={`Accordion ${idx}`}
                                    title={`${number}. ${title}`}
                                >
                                    {desc}
                                </AccordionItem>
                            );
                        }
                    )}
                </Accordion>
            </div>
        </Card>
    );
}
