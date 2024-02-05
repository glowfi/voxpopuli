import { Card, CardBody, CardFooter, Divider, Link } from '@nextui-org/react';

export default function App() {
    return (
        <Card className="max-w-[400px]">
            <CardBody>
                <div className="footercontent px-3 py-3 md:hello">
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
                <div className="footercontent px-3 py-3">
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
                <div className="footercontent px-3 py-3">
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
                No rights reserved. Built for educational purposes.
                <Link isExternal href="https://github.com/nextui-org/nextui">
                    Visit source code on GitHub.
                </Link>
            </CardFooter>
        </Card>
    );
}
