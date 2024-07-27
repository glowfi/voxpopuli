'use client';
import {
    Credenza,
    CredenzaTrigger,
    CredenzaContent,
    CredenzaHeader,
    CredenzaTitle,
    CredenzaDescription,
    CredenzaBody,
    CredenzaFooter,
    CredenzaClose
} from '@/components/ui/credenza';
import React, { Dispatch, SetStateAction } from 'react';

const MediaViewer = ({
    open,
    setIsopen
}: {
    open: boolean;
    setIsopen: Dispatch<SetStateAction<boolean>>;
}) => {
    return (
        <Credenza open={open}>
            <CredenzaTrigger asChild>
                <button className="hidden">Open modal</button>
            </CredenzaTrigger>
            <CredenzaContent>
                <CredenzaHeader>
                    <CredenzaTitle>Credenza</CredenzaTitle>
                    <CredenzaDescription>
                        A responsive modal component for shadcn/ui.
                    </CredenzaDescription>
                </CredenzaHeader>
                <CredenzaBody>
                    This component is built using shadcn/ui&apos;s dialog and
                    drawer component, which is built on top of Vaul.
                </CredenzaBody>
                <CredenzaFooter>
                    <CredenzaClose asChild>
                        <button
                            onClick={() => {
                                setIsopen(false);
                            }}
                        >
                            Close
                        </button>
                    </CredenzaClose>
                </CredenzaFooter>
            </CredenzaContent>
        </Credenza>
    );
};

export default React.memo(MediaViewer);
