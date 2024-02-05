import { CircularProgress, Divider } from '@nextui-org/react';
import Footer from '../Components/HomeFeed/Footer';
import GetPremium from '../Components/HomeFeed/GetPremium';
import HomePost from '../Components/HomeFeed/HomePost';
import Nav from '../Components/HomeFeed/Navbar';
import Post from '../Components/HomeFeed/Post';
import TopGrowing from '../Components/HomeFeed/TopGrowing';
import TrendingBar from '../Components/HomeFeed/TrendingBar';
import '../styles/homefeed.css';

const App = () => {
    return (
        <>
            <Nav />
            <Divider className="my-5" />
            <div className="container">
                <TrendingBar />
                <div className="item-3 items">
                    <TopGrowing cname="aside-item" />
                    <GetPremium />
                    <HomePost />
                    <Footer />
                </div>
                <div className="items posts">
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                    <Post cname="post" />
                </div>
                <div className="progress item-4 flex justify-center">
                    <CircularProgress label="Loading..." />
                </div>
            </div>
        </>
    );
};

export default App;
