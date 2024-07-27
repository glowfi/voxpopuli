import PopularCommunities from './popular-communities';
import Posts from './posts';
import Sortposts from './sortposts';
import Trendingtopics from './trending-topics';

const Allposts = () => {
    return (
        <div className="flex flex-col min-h-screen justify-center items-center">
            <div className="flex flex-1 gap-4 p-3 sm:p-6 justify-between">
                <div className="hidden w-80 space-y-4 xl:block">
                    <Trendingtopics />
                </div>
                <div className="flex flex-col flex-wrap">
                    <Sortposts />
                    <Posts />
                </div>
                <div className="hidden w-80 space-y-4 xl:block">
                    <PopularCommunities />
                </div>
            </div>
        </div>
    );
};

export default Allposts;
