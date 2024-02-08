import { Card, CardBody, Image } from '@nextui-org/react';
import { useEffect, useState } from 'react';
import userarray from '../../data/userArray';

const Trophies = ({ pname }) => {
    const [userData, setUserData] = useState([]);
    useEffect(() => {
        const fetchData = () => {
            const data = userarray.filter((p) => p.username == pname);
            setUserData(() => data[0]);
        };
        fetchData();
    }, []);
    return (
        <Card className="max-w-[400px] w-full">
            <CardBody>
                <h1 className="text-center">
                    <b>Trophy Case ({userData?.trophies?.length})</b>
                </h1>
                <div className="cabinet flex flex-col justify-start text-center">
                    {userData?.trophies?.map((p, idx) => {
                        return (
                            <div
                                className="trophy flex flex-ro gap-3"
                                key={idx}
                                style={{ margin: '1rem' }}
                            >
                                <Image
                                    className="imgtrophies"
                                    src={`/resources/images/${p}.png`}
                                />
                                <p className="text-md">{p}</p>
                            </div>
                        );
                    })}
                </div>
            </CardBody>
        </Card>
    );
};

export default Trophies;
