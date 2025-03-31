import React, { FunctionComponent } from 'react';

function truncateString(str: string, maxLength: number): string {
    return str.length > maxLength ? str.slice(0, maxLength) + ' ...' : str;
}

interface TextProps {
    data: { text: string };
}

const Text: FunctionComponent<TextProps> = ({ data }) => {
    return (
        <div className="my-2">
            <div className="text-slate-800 dark:text-slate-200 whitespace-pre-line text-xs sm:text-sm md:text-base lg:text-xl line-clamp-2">
                {/* {truncateString(data.text, 150)} */}
                {data.text}
            </div>
        </div>
    );
};

export default Text;
