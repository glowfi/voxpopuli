import React from 'react';
import { useParams } from 'react-router-dom';
import '../styles/group.css';
import GroupPosts from '../Components/Groups/GroupPosts';
import TrendingBar from '../Components/Groups/TrendingBar';
import About from '../Components/Groups/About';
import Banner from '../Components/Groups/Banner';
import Rules from '../Components/Groups/Rules';
import Flairs from '../Components/Groups/Flairs';

const Groups = ({ theme }) => {
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
                <GroupPosts cname="postg" gname={params.name} theme={theme} />
            </div>
        </div>
    );
};

export default React.memo(Groups);
