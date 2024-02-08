import React from 'react';
import { useParams } from 'react-router-dom';
import Over from '../Components/Person/Over';
import PersonPost from '../Components/Person/PersonPost';
import PersonComments from '../Components/Person/PersonComments';
import About from '../Components/Person/About';
import Trophies from '../Components/Person/Trophies';
import postsarr from '../data/postsArray';
import '../styles/person.css';
import '../styles/person2.css';

const Person = ({ theme }) => {
    const params = useParams();
    let newArr = postsarr.filter(({ author }) => author == params.name);
    console.log('ASDS', newArr);

    return (
        <>
            {newArr.length == 0 ? (
                <div className="containerpe2">
                    <div className="over2">
                        <Over />
                    </div>
                    <div className="asidepe2">
                        <About pname={params.name} />
                        <Trophies pname={params.name} />
                    </div>
                    <div className="postspe2">
                        <PersonPost
                            cname="postpe"
                            pname={params.name}
                            theme={theme}
                        />
                    </div>
                    <div className="commentspe2">
                        <PersonComments
                            cname="commentpe"
                            pname={params.name}
                            theme={theme}
                        />
                    </div>
                </div>
            ) : (
                <div className="containerpe">
                    <div className="over">
                        <Over />
                    </div>
                    <div className="asidepe">
                        <About pname={params.name} />
                        <Trophies pname={params.name} />
                    </div>
                    <div className="postspe">
                        <PersonPost
                            cname="postpe"
                            pname={params.name}
                            theme={theme}
                        />
                    </div>
                    <div className="commentspe">
                        <PersonComments
                            cname="commentpe"
                            pname={params.name}
                            theme={theme}
                        />
                    </div>
                </div>
            )}
        </>
    );
};

export default React.memo(Person);
