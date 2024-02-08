import React from 'react';
import { useParams } from 'react-router-dom';
import About from '../Components/Groups/About';
import Flairs from '../Components/Groups/Flairs';
import Rules from '../Components/Groups/Rules';
import RealPost from '../Components/Post/RealPost';
import '../styles/post.css';

import EditorComponent from '../Components/Post/Editor';
import AllComments from '../Components/Post/AllComments';
import { Button } from '@nextui-org/react';
const Post = ({ theme }: any) => {
    const params = useParams();

    return (
        <div className="containerpo">
            <div className="postspo">
                <RealPost
                    cname="postpo"
                    gname={params.name}
                    pid={params.id}
                    theme={theme}
                />

                <div className="comments">
                    <div className="writecomment">
                        <p className="text-default-900 lg:text-large md:text-medium sm:text-small">
                            <b>Comment as User</b>
                        </p>
                        <EditorComponent />
                        <Button
                            variant="solid"
                            color="primary"
                            style={{ width: '50px', justifyItems: 'end' }}
                            isDisabled
                        >
                            Comment
                        </Button>
                    </div>
                    <div className="allcomments">
                        <p className="text-default-900 lg:text-large md:text-medium sm:text-small px-3 py-3">
                            <b>Comment(s)</b>
                        </p>

                        <AllComments />
                    </div>
                </div>
            </div>
            <div className="asidepo">
                <About gname={params.name} />
                <Rules gname={params.name} />
                <Flairs gname={params.name} />
            </div>
        </div>
    );
};

export default React.memo(Post);
