import { Button } from '@nextui-org/react';
import { useState, useCallback } from 'react';

export default function App({}) {
    const [colorstatus, setColorstatus] = useState({
        Overview: 'danger',
        Posts: 'default',
        Comments: 'default'
    });
    console.log(colorstatus);
    const handleClick = useCallback(
        (args: any) => {
            //@ts-ignore
            const type = args;
            let cp = { ...colorstatus };
            for (const key in cp) {
                if (key == args) {
                    //@ts-ignore
                    cp[key] = 'danger';
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
        <div className="flex flex-row justify-center gap-3">
            <Button
                variant="shadow"
                color={colorstatus['Overview']}
                onClick={() => {
                    handleClick('Overview');
                }}
                className="btnanchor"
            >
                Overview
            </Button>
            <Button
                variant="shadow"
                color={colorstatus['Posts']}
                onClick={() => {
                    handleClick('Posts');
                }}
                className="btnanchor"
            >
                Posts
            </Button>
            <Button
                variant="shadow"
                color={colorstatus['Comments']}
                onClick={() => {
                    handleClick('Comments');
                }}
                className="btnanchor"
            >
                Comments
            </Button>
        </div>
    );
}
