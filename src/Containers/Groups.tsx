import React from 'react';
import { useParams } from 'react-router-dom';
import About from '../Components/Groups/About';
import Banner from '../Components/Groups/Banner';
import Flairs from '../Components/Groups/Flairs';
import GroupPosts from '../Components/Groups/GroupPosts';
import Rules from '../Components/Groups/Rules';
import TrendingBar from '../Components/Groups/TrendingBar';
import '../styles/group.css';

interface Props {
    theme: string;
}

const Groups = ({ theme }: Props) => {
    const params = useParams();
    return (
        <div className="containerg">
            <Banner gname={params.name} />
            <TrendingBar />
            <div className="asideg">
                <About gname={params.name} />
                <Rules gname={params.name} />
                <Flairs gname={params.name} />
            </div>
            <div className="postsg">
                {/* @ts-ignore */}
                <GroupPosts cname="postg" gname={params.name} theme={theme} />
            </div>
        </div>
    );
};

export default React.memo(Groups);
