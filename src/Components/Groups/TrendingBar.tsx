import { Button } from '@nextui-org/react';
import React, { useCallback, useState } from 'react';
import { BsRocketFill } from 'react-icons/bs';
import { FaRegSun } from 'react-icons/fa';
import { MdBarChart } from 'react-icons/md';
import { SlFire } from 'react-icons/sl';

const TrendingBar = () => {
    const [colorstatus, setColorstatus] = useState({
        best: 'success',
        hot: 'default',
        new: 'default',
        top: 'default'
    });
    const handleClick = useCallback(
        (args: any) => {
            const type = args;
            let cp = { ...colorstatus };
            for (const key in cp) {
                if (key == args) {
                    //@ts-ignore
                    cp[key] = 'success';
                } else {
                    //@ts-ignore
                    cp[key] = 'default';
                }
            }
            console.log(cp);
            setColorstatus(() => cp);
        },
        [colorstatus]
    );

    return (
        <>
            <div className="trendingg">
                <Button
                    size="md"
                    //@ts-ignore
                    color={colorstatus['best']}
                    variant="flat"
                    onClick={() => {
                        handleClick('best');
                    }}
                >
                    <BsRocketFill />
                    Best
                </Button>
                <Button
                    size="md"
                    //@ts-ignore
                    color={colorstatus['hot']}
                    variant="flat"
                    onClick={() => {
                        handleClick('hot');
                    }}
                >
                    <SlFire />
                    Hot
                </Button>
                <Button
                    size="md"
                    //@ts-ignore
                    color={colorstatus['new']}
                    variant="flat"
                    onClick={() => {
                        handleClick('new');
                    }}
                >
                    <FaRegSun />
                    New
                </Button>
                <Button
                    size="md"
                    //@ts-ignore
                    color={colorstatus['top']}
                    variant="flat"
                    onClick={() => {
                        handleClick('top');
                    }}
                >
                    <MdBarChart />
                    Top
                </Button>
            </div>
        </>
    );
};

export default React.memo(TrendingBar);
