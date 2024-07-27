import { prisma } from '@voxpopuli/db';
import ImageKit from 'imagekit';

// Dummy query
await prisma.$queryRaw`SELECT 1`;

var imagekit = new ImageKit({
    privateKey: 'private_ZEdb1Xkl5F8MTChESue/AMdlC2c=',
    publicKey: 'public_/VwD4BfXGuqXPTQ2HNlq8Si59nE=',
    urlEndpoint: 'https://ik.imagekit.io/gys1ezqjb'
});

// Trophies
const fn1 = async () => {
    let trophies = await prisma.trophies.findMany({});
    for (let index = 0; index < trophies.length; index++) {
        let curr = trophies[index];
        let filename = crypto.randomUUID();

        try {
            let res = await imagekit.upload({
                file: curr.imageLink,
                fileName: filename,
                folder: `social-media/media-content/trophies`
            });

            await prisma.trophies.update({
                where: {
                    id: curr.id
                },
                data: {
                    imageLink: res.url
                }
            });
        } catch (err) {
            let res = await imagekit.upload({
                file: `https://placehold.co/300x300/black/white?text=${curr.title}`,
                fileName: filename,
                folder: `social-media/media-content/trophies`
            });

            await prisma.awards.update({
                where: {
                    id: curr.id
                },
                data: {
                    imageLink: res.url
                }
            });
        }
    }
};

// Awards
const fn2 = async () => {
    let awards = await prisma.awards.findMany({});
    for (let index = 0; index < awards.length; index++) {
        let curr = awards[index];
        let filename = crypto.randomUUID();

        try {
            let res = await imagekit.upload({
                file: curr.imageLink,
                fileName: filename,
                folder: `social-media/media-content/awards`
            });
            await prisma.awards.update({
                where: {
                    id: curr.id
                },
                data: {
                    imageLink: res.url
                }
            });
        } catch (err) {
            let res = await imagekit.upload({
                file: `https://placehold.co/300x300/black/white?text=${curr.title}`,
                fileName: filename,
                folder: `social-media/media-content/awards`
            });

            await prisma.awards.update({
                where: {
                    id: curr.id
                },
                data: {
                    imageLink: res.url
                }
            });
        }
    }
};

// Users-profile-pics
const fn3 = async () => {
    let allpics = [];

    for (let index = 0; index < 300; index++) {
        let filename = crypto.randomUUID();
        let res = await imagekit.upload({
            file: `https://robohash.org/${filename}.png`,
            fileName: filename,
            folder: `social-media/media-content/profilepics`
        });
        allpics.push(res.url);
    }

    let pics = await prisma.user.findMany({});
    for (let index = 0; index < pics.length; index++) {
        let curr = pics[index];
        await prisma.user.update({
            where: {
                id: curr.id
            },
            data: {
                avatarImg: allpics[Math.floor(Math.random() * allpics.length)]
            }
        });
    }
};

// Subreddit Images
const fn4 = async () => {
    let awards = await prisma.voxsphere.findMany({});
    for (let index = 0; index < awards.length; index++) {
        let curr = awards[index];
        let filename = crypto.randomUUID();

        try {
            let res1 = await imagekit.upload({
                file:
                    curr.logoURL ||
                    `https://placehold.co/300x300/black/white?text=${curr.name}`,
                fileName: `logo-${filename}`,
                folder: `social-media/media-content/voxsphere`
            });

            let res2 = await imagekit.upload({
                file:
                    curr.bannerURL ||
                    `https://placehold.co/1000x1000/black/white?text=${curr.name}`,
                fileName: `banner-${filename}`,
                folder: `social-media/media-content/voxsphere`
            });

            await prisma.voxsphere.update({
                where: {
                    id: curr.id
                },
                data: {
                    logoURL: res1.url,
                    bannerURL: res2.url
                }
            });
        } catch (err) {
            let res1 = await imagekit.upload({
                file: `https://placehold.co/300x300/black/white?text=${curr.name}`,
                fileName: `logo-${filename}`,
                folder: `social-media/media-content/voxsphere`
            });
            let res2 = await imagekit.upload({
                file: `https://placehold.co/1000x1000/black/white?text=${curr.name}`,
                fileName: `banner-${filename}`,
                folder: `social-media/media-content/voxsphere`
            });

            await prisma.voxsphere.update({
                where: {
                    id: curr.id
                },
                data: {
                    logoURL: res1.url,
                    bannerURL: res2.url
                }
            });
        }
    }
};

const runall = async () => {
    await fn1();
    await fn2();
    await fn3();
    await fn4();
};

runall().then(() => {
    console.log('Done !');
});
