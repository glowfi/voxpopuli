import { Divider, NextUIProvider } from '@nextui-org/react';
import { Route, Routes, useNavigate } from 'react-router-dom';
import Nav from './Components/Global/Navbar';
import Groups from './Containers/Groups';
import HomeFeed from './Containers/HomeFeed';
import NotFound from './Containers/NotFound';
import Person from './Containers/Person';
import Post from './Containers/Post';
import Submit from './Containers/Submit';
import { useCallback, useState } from 'react';

const App = () => {
    const navigate = useNavigate();
    const [theme, setTheme] = useState('dark');
    const changeTheme = useCallback(() => {
        setTheme(() => (theme == 'light' ? 'dark' : 'light'));
        document.documentElement.setAttribute(
            'data-theme',
            theme == 'light' ? 'dark' : 'light'
        );
    }, [theme]);
    return (
        <NextUIProvider navigate={navigate}>
            <main className={`${theme} text-foreground bg-background`}>
                <Nav changeTheme={changeTheme} theme={theme} />
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
            </main>
        </NextUIProvider>
    );
};

export default App;
