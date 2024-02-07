import React from 'react';
import Footer from '../Components/HomeFeed/Footer';
import GetPremium from '../Components/HomeFeed/GetPremium';
import HomePost from '../Components/HomeFeed/HomePost';
import Post from '../Components/HomeFeed/Post';
import TopGrowing from '../Components/HomeFeed/TopGrowing';
import TrendingBar from '../Components/HomeFeed/TrendingBar';
import '../styles/homefeed.css';

interface Props {
    theme: string;
}
const App = ({ theme }: Props) => {
    return (
        <>
            <div className="container">
                <TrendingBar />
                <div className="item-3 items">
                    <TopGrowing cname="aside-item" />
                    <HomePost />
                    <GetPremium />
                    <Footer />
                </div>
                <div className="items posts">
                    <Post cname="post" theme={theme} />
                </div>
            </div>
        </>
    );
};

export default React.memo(App);
