import { Button } from '@nextui-org/react';
import React from 'react';
import { BsRocketFill } from 'react-icons/bs';
import { FaRegSun } from 'react-icons/fa';
import { MdBarChart } from 'react-icons/md';
import { SlFire } from 'react-icons/sl';

const TrendingBar = () => {
    return (
        <>
            <div className="item-2 items">
                <Button size="md" color="success" variant="flat">
                    <BsRocketFill />
                    Best
                </Button>
                <Button size="md">
                    <SlFire />
                    Hot
                </Button>
                <Button size="md">
                    <FaRegSun />
                    New
                </Button>
                <Button size="md">
                    <MdBarChart />
                    Top
                </Button>
            </div>
        </>
    );
};

export default React.memo(TrendingBar);
