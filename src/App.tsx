import { Divider, NextUIProvider } from '@nextui-org/react';
import { Route, Routes, useNavigate } from 'react-router-dom';
import Nav from './Components/Global/Navbar';
import Groups from './Containers/Groups';
import HomeFeed from './Containers/HomeFeed';
import NotFound from './Containers/NotFound';
import Person from './Containers/Person';
import Post from './Containers/Post';
import Submit from './Containers/Submit';

const App = () => {
    const navigate = useNavigate();
    return (
        <NextUIProvider navigate={navigate}>
            <Nav />
            <Divider className="my-5" />
            <Routes>
                <Route path="/" element={<HomeFeed />} />
                <Route path="/g">
                    <Route path=":name" element={<Groups />} />
                    <Route path=":name/:id" element={<Post />} />
                </Route>
                <Route path="/p">
                    <Route path=":name" element={<Person />} />
                </Route>

                <Route path="/submit" element={<Submit />} />
                <Route path="*" element={<NotFound />} />
            </Routes>
        </NextUIProvider>
    );
};

export default App;
