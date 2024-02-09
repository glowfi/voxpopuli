import React from 'react';
import { useParams } from 'react-router-dom';
import Over from '../Components/Person/Over';
import PersonPost from '../Components/Person/PersonPost';
import PersonComments from '../Components/Person/PersonComments';
import About from '../Components/Person/About';
import Trophies from '../Components/Person/Trophies';
import postsarr from '../data/postsArray';
import '../styles/person.css';

const Person = ({ theme }) => {
    const params = useParams();
    let newArr = postsarr.filter(({ author }) => author == params.name);

    let cname = 'containerpe';
    if (newArr.length == 0) {
        cname = 'containerpe2';
    }

    return (
        <div className={cname}>
            <div className="over">
                <Over />
            </div>
            <div className="asidepe">
                <About pname={params.name} />
                <Trophies pname={params.name} />
            </div>
            {newArr.length == 0 ? (
                ''
            ) : (
                <div className="postspe">
                    <PersonPost
                        cname="postpe"
                        pname={params.name}
                        theme={theme}
                    />
                </div>
            )}
            <div className="commentspe">
                <PersonComments
                    cname="commentpe"
                    pname={params.name}
                    theme={theme}
                />
            </div>
        </div>
    );
};

export default React.memo(Person);
