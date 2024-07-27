import { prisma } from '@voxpopuli/db';
import fs from 'fs';

// Dummy query
await prisma.$queryRaw`SELECT 1`;

const handleGif = async (mediaLocations: any, data: any, id: number) => {
    if (data in mediaLocations) {
        await prisma.mediaMetadata.update({
            where: {
                id
            },
            data: {
                gifLink: mediaLocations[data]
            }
        });
    }
};
const handleVideo = async (mediaLocations: any, data: any, id: number) => {
    if (data in mediaLocations) {
        await prisma.mediaMetadata.update({
            where: {
                id
            },
            data: {
                videoLink: mediaLocations[data]
            }
        });
    }
};
const handleImage = async (mediaLocations: any, data: any, id: number) => {
    if (data in mediaLocations) {
        await prisma.mediaMetadata.update({
            where: {
                id
            },
            data: {
                imageLink: mediaLocations[data]
            }
        });
    }
};
const handleGallery = async (mediaLocations: any, data: any, id: number) => {
    let newJson = [];
    for (let i = 0; i < data.length; i++) {
        const [local, original] = data[i];
        if (local in mediaLocations) {
            newJson.push([mediaLocations[local], original]);
        } else {
            newJson.push(data[i]);
        }
    }
    await prisma.mediaMetadata.update({
        where: {
            id
        },
        data: {
            gallery: newJson
        }
    });
};

const run1 = async () => {
    let posts = await prisma.posts.findMany({
        include: {
            mediaContent: {
                select: {
                    postId: true,
                    id: true,
                    type: true,
                    gallery: true,
                    gifLink: true,
                    imageLink: true,
                    videoLink: true,
                    original_link: true
                }
            }
        }
    });

    const mediaLocations = JSON.parse(
        fs.readFileSync('./json/assets.json', 'utf8')
    );

    for (let index = 0; index < posts.length; index++) {
        let currPost = posts[index];
        if (
            currPost.mediaContent.length > 0 &&
            currPost.mediaContent[0].type === 'gallery'
        ) {
            let currGallery = currPost.mediaContent[0].gallery;
            await handleGallery(
                mediaLocations,
                currGallery,
                currPost.mediaContent[0].id
            );
        } else if (
            currPost.mediaContent.length > 0 &&
            currPost.mediaContent[0].type === 'gif'
        ) {
            let currGif = currPost.mediaContent[0].gifLink;
            await handleGif(
                mediaLocations,
                currGif,
                currPost.mediaContent[0].id
            );
        } else if (
            currPost.mediaContent.length > 0 &&
            currPost.mediaContent[0].type === 'video'
        ) {
            let currVideo = currPost.mediaContent[0].videoLink;
            await handleVideo(
                mediaLocations,
                currVideo,
                currPost.mediaContent[0].id
            );
        } else if (
            currPost.mediaContent.length > 0 &&
            currPost.mediaContent[0].type === 'image'
        ) {
            let currImage = currPost.mediaContent[0].imageLink;
            await handleImage(
                mediaLocations,
                currImage,
                currPost.mediaContent[0].id
            );
        }
    }
};

(async () => {
    await run1();
    console.log('Done updating media locations!');
})();
