import { Card, CardBody, CardFooter, Divider, Link } from '@nextui-org/react';
import React from 'react';

function App() {
    return (
        <Card className="max-w-[400px] w-full">
            <CardBody>
                <div className="footercontent px-1 py-2">
                    <div>
                        <Link isExternal href="">
                            Help
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Imprint
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Transparency Report
                        </Link>
                    </div>
                </div>
                <Divider />
                <div className="footercontent px-3 py-4">
                    <div>
                        <Link isExternal href="">
                            User Agreement
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Privacy Policy
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Content Policy
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Code of Conduct
                        </Link>
                    </div>
                </div>
                <Divider />
                <div className="footercontent px-2 py-3">
                    <div>
                        <Link isExternal href="">
                            English
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Francais
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Deutsch
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Italiano
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Espanol
                        </Link>
                    </div>
                    <div>
                        <Link isExternal href="">
                            Portugues
                        </Link>
                    </div>
                </div>
            </CardBody>
            <Divider />
            <CardFooter
                className={'flex flex-col gap-3 items-center justify-center'}
            >
                <div>
                    <p className="text-center">
                        No rights reserved. Built for educational purposes.
                    </p>
                </div>
                <Link
                    isExternal
                    href="https://github.com/glowfi/voxpopuli"
                    showAnchorIcon
                >
                    Visit source code on GitHub.
                </Link>
            </CardFooter>
        </Card>
    );
}

export default React.memo(App);
