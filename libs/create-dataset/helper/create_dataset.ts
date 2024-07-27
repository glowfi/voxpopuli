import 'dotenv/config';
import { configDotenv } from 'dotenv';
configDotenv({
    path: '../.env'
});
import fs from 'node:fs';
import path from 'path';
import crypto from 'crypto';
import readable from 'readable-numbers';
import { getLink } from './getlinks';
import { prisma } from '@voxpopuli/db';

// Dummy query
await prisma.$queryRaw`SELECT 1`;

const handleUser = async (data: any, voxsphereData: any) => {
    for (const key in data) {
        let currUser = data[key];

        let userID = currUser['id'];

        let voxspehereIDS_member = [];
        let voxspehereIDS_moderator = [];
        let trophiesIDS = [];

        // Get trophies
        for (let i = 0; i < currUser['trophies'].length; i++) {
            let title = currUser['trophies'][i]['title'].toString();
            let res = await prisma.trophies.findUnique({
                where: {
                    title: title.trim()
                }
            });
            if (res?.id) {
                trophiesIDS.push({ id: res.id });
            }
        }

        // Get member in which voxsphere
        for (let i = 0; i < currUser['subreddits_member'].length; i++) {
            const [id, name] = currUser['subreddits_member'][i];
            let res = await prisma.voxsphere.findUnique({
                where: {
                    name
                }
            });
            if (res?.id) {
                voxspehereIDS_member.push({ id: res.id });
            }
        }

        // Get moderator in which voxsphere
        for (let i = 0; i < voxsphereData.length; i++) {
            let currSpehereModerators = voxsphereData[i]['moderators'];
            for (let j = 0; j < currSpehereModerators.length; j++) {
                if (currSpehereModerators.length > 0) {
                    let [name, id] = currSpehereModerators[j];
                    if (id === userID) {
                        let res = await prisma.voxsphere.findUnique({
                            where: {
                                name: voxsphereData[i]['title'].replace(
                                    'r/',
                                    ''
                                )
                            }
                        });
                        if (res?.id) {
                            voxspehereIDS_moderator.push({ id: res.id });
                        }
                    }
                }
            }
        }

        try {
            await prisma.user.create({
                data: {
                    email: `${crypto.randomUUID()}@xyz.com`,
                    username: currUser['username'],
                    password: 'xyz',
                    cakeDay: currUser['cakeDay'],
                    cakeDayHuman: currUser['cakeDayHuman'],
                    accountAge: currUser['age'],
                    avatarImg: currUser['avatar_img'],
                    bannerImg: currUser['banner_img'],
                    publicDescription: currUser['publicDescription'],
                    over18: currUser['over18'],
                    keycolorHex: currUser['keycolor'],
                    primarycolorHex: currUser['primarycolor'],
                    iconcolorHex: currUser['iconcolor'],
                    suspended: currUser['suspended'],
                    trophies: {
                        connect: trophiesIDS
                    },
                    memberOf: {
                        connect: voxspehereIDS_member
                    },
                    moderatorOf: {
                        connect: voxspehereIDS_moderator
                    }
                }
            });
        } catch (err) {
            throw err;
        }
    }
    console.log('Done inserting users!');
};

// Read the files
function readFile(filePath: string, type: string) {
    try {
        // Read file from file system
        const fileContent = fs.readFileSync(filePath, 'utf8');
        // Parse JSON string to JavaScript object
        if (
            type === 'user' ||
            type === 'post' ||
            type === 'voxspheres' ||
            type === 'trophies' ||
            type === 'awards'
        ) {
            return JSON.parse(fileContent);
        }
        return null;
    } catch (error) {
        console.error('Error reading JSON file:', error);
        throw error; // Rethrow the error to handle it outside if needed
    }
}

const handleTrophies = async (data: any) => {
    let allTrophies = [];
    for (let index = 0; index < data.length; index++) {
        let currTrophy = data[index];
        allTrophies.push({
            title: currTrophy['title'],
            imageLink: currTrophy['image_link'],
            description: currTrophy['description']
        });
    }
    try {
        await prisma.trophies.createMany({
            data: [...allTrophies]
        });
    } catch (err) {
        throw err; // Rethrow the error to handle it outside if needed
    }
};

const handleAwards = async (data: any) => {
    let allAwards = [];
    for (let index = 0; index < data.length; index++) {
        let currAward = data[index];
        allAwards.push({
            title: currAward['title'],
            imageLink: currAward['image_link']
        });
    }

    try {
        await prisma.awards.createMany({
            data: [...allAwards]
        });
    } catch (err) {
        throw err; // Rethrow the error to handle it outside if needed
    }
};

const handleVoxSphere = async (data: any) => {
    // Insert topics
    let allTopics = Object.keys(data);
    for (let i = 0; i < allTopics.length; i++) {
        let currVal = structuredClone(allTopics[i]);
        allTopics[i] = { title: currVal };
    }

    try {
        // DISABLE THIS [trophies + awards while appending new data]
        await prisma.topics.createMany({
            data: [...allTopics]
        });

        let allSpeheresdata = Object.values(data).flat();
        for (let index = 0; index < allSpeheresdata.length; index++) {
            let currSphere = allSpeheresdata[index];
            let allRules = [];
            let allFlairs = [];

            for (let i = 0; i < currSphere['rules'].length; i++) {
                let currRule = currSphere['rules'][i];
                allRules.push({
                    title: currRule['rule_title'],
                    description: currRule['rule_desc']
                });
            }
            for (let i = 0; i < currSphere['flairs'].length; i++) {
                let currFlair = currSphere['flairs'][i];
                allFlairs.push({
                    title: currFlair['text'],
                    colorHex: currFlair['color']
                });
            }

            const res = await prisma.topics.findFirst({
                where: {
                    title: currSphere['category']
                }
            });

            try {
                await prisma.voxsphere.create({
                    data: {
                        name: currSphere['title'].replace('r/', ''),
                        topics: {
                            connectOrCreate: {
                                where: {
                                    title: currSphere['category'],
                                    id: res?.id
                                },
                                create: {
                                    title: currSphere['category']
                                }
                            }
                        },
                        about: currSphere['about'],
                        logoURL: currSphere['logoUrl'],
                        bannerURL: currSphere['bannerUrl'],
                        flairs: {
                            createMany: {
                                data: [...allFlairs]
                            }
                        },
                        rules: {
                            createMany: {
                                data: [...allRules]
                            }
                        },
                        anchors: currSphere['anchors'],
                        buttonColorHex: currSphere['buttonColor'],
                        headerColorHex: currSphere['headerColor'],
                        bannerBgColorHex: currSphere['banner_background_color'],
                        createdAtUnix: currSphere['creationDate'],
                        createdatHuman: currSphere['creationDateHuman'],
                        over18: currSphere['over18'],
                        spoilersEnabled: currSphere['spoilers_enabled'],
                        totalmembers: 0,
                        totalmembersHuman: ''
                    }
                });
            } catch (err) {
                throw err; // Rethrow the error to handle it outside if needed
            }
        }
    } catch (err) {
        throw err;
    }
};

async function handlePosts(data: any) {
    for (let index = 0; index < data.length; index++) {
        let currPost = data[index];

        //Find author
        let getAuthor = await prisma.user.findUnique({
            where: {
                username: currPost['author']
            }
        });

        //Find awards
        let allAwardsID = [];
        for (let i = 0; i < currPost['awards'].length; i++) {
            const { title } = currPost['awards'][i];
            let res = await prisma.awards.findUnique({
                where: {
                    title
                }
            });
            if (res?.id) {
                allAwardsID.push({ id: res.id });
            }
        }

        // Find topic
        let getTopic = await prisma.topics.findUnique({
            where: {
                title: currPost['category']
            }
        });

        // mediacontent [Done in the body] gif video image multi
        let gallery = [];
        let imageLink = '';
        let videoLink = '';
        let gifLink = '';
        let original_link = '';

        const gif_gallery_image_video_in_media_content_type = (
            type: string
        ) => {
            let ls = ['gif', 'gallery', 'image', 'video'];
            for (let index = 0; index < ls.length; index++) {
                if (ls[index] === type) {
                    return true;
                }
            }
            return false;
        };

        let type = null;
        if ('media_content' in currPost) {
            type = currPost['media_content']['type'];
        }

        if (type && gif_gallery_image_video_in_media_content_type(type)) {
            let postMediaData = currPost['media_content']['data'];

            const getVideoLink = async (data: any) => {
                const { dash_url, hls_url } = data;
                let res = await getLink('video', hls_url ? hls_url : dash_url);
                original_link = hls_url;
                return res[0];
            };
            const getGifLink = async (data: any) => {
                let { url } = data[data.length - 1];
                let res = await getLink('gif', url);
                original_link = url;
                return res[0];
            };
            const getImageLink = async (data: any) => {
                console.log('CURR DATA', data);
                let url = '';
                let res = null;
                if ('url' in data[data.length - 1]) {
                    url = data[data.length - 1]['url'];
                    res = await getLink('image', url);
                    original_link = url;
                } else if ('u' in data[data.length - 1]) {
                    url = data[data.length - 1]['u'];
                    res = await getLink('image', url);
                    original_link = url;
                } else if ('gif' in data[data.length - 1]) {
                    url = data[data.length - 1]['gif'];
                    res = await getLink('gif', url);
                    original_link = url;
                } else if ('mp4' in data[data.length - 1]) {
                    url = data[data.length - 1]['mp4'];
                    res = await getLink('gif', url);
                    original_link = url;
                }
                console.log('IF  ERROR:', url, data);
                return res[0];
            };
            const getGalleryLink = async (data: any) => {
                console.log(data);
                let allLinks = [];
                for (let i = 0; i < data.length; i++) {
                    let { id, pics } = data[i];
                    allLinks.push([await getImageLink(pics), original_link]);
                }
                console.log(allLinks);
                return allLinks;
            };

            // Handle Multi
            if (type === 'multi') {
                for (let i = 0; i < postMediaData.length; i++) {
                    let currData = postMediaData[i];
                    if (currData['type'] === 'video') {
                        videoLink = await getVideoLink(currData['data']);
                    } else if (currData['type'] === 'image') {
                        imageLink = await getImageLink(currData['data']);
                    } else if (currData['type'] === 'gif') {
                        gifLink = await getGifLink(currData['data']);
                    }
                }
            } else if (type === 'video') {
                videoLink = await getVideoLink(postMediaData);
            } else if (type === 'image') {
                imageLink = await getImageLink(postMediaData);
            } else if (type === 'gif') {
                gifLink = await getGifLink(postMediaData);
            } else if (type === 'gallery') {
                gallery = await getGalleryLink(postMediaData);
            }
        }

        // comments
        const parseComments = async (data: any, parentID: any, postID: any) => {
            let commentIDS = [];
            for (let i = 0; i < data.length; i++) {
                const {
                    author,
                    comment,
                    comment_html,
                    comment_ups,
                    parent_comment_id,
                    replies,
                    created_utc
                } = data[i];

                let authorID = await prisma.user.findUnique({
                    where: {
                        username: author
                    }
                });

                let res = await prisma.comments.create({
                    data: {
                        comment,
                        commentHTML: comment_html,
                        createdAtUnix: created_utc,
                        ups: comment_ups,
                        post: {
                            connect: {
                                id: postID
                            }
                        },
                        author: {
                            connect: {
                                id: authorID?.id
                            }
                        },
                        ...(parent_comment_id === 'isParent'
                            ? {}
                            : {
                                  parentComment: {
                                      connect: {
                                          id: parentID
                                      }
                                  }
                              })
                    }
                });

                data = parseComments(replies, res.id, postID);

                if (data.length > 0) {
                    // Connect replies
                    await prisma.comments.update({
                        where: {
                            id: res.id,
                            postId: postID
                        },
                        data: {
                            replies: {
                                connect: data
                            }
                        }
                    });
                }
                commentIDS.push({ id: res.id });
            }
            return commentIDS;
        };

        // voxsphere
        let getVoxSphere = await prisma.voxsphere.findUnique({
            where: {
                name: currPost['subreddit']
            }
        });

        // flair
        let resFlair = null;
        if (currPost['postflair']) {
            resFlair = await prisma.flairs.findFirst({
                where: {
                    voxsphereId: getVoxSphere?.id,
                    title: currPost['postflair'],
                    colorHex: currPost['postflaircolor']
                }
            });
        }

        // media_content

        // If its a link type and has some metadata too
        let media_content = null;
        if (
            currPost['link_type'] === true &&
            type &&
            gif_gallery_image_video_in_media_content_type(type)
        ) {
            media_content = {
                type,
                original_link: original_link,
                gallery: gallery,
                imageLink: imageLink,
                videoLink: videoLink,
                gifLink: gifLink
            };
        }
        // If its a link type but does not have metadata
        else if (
            currPost['link_type'] === true &&
            type &&
            !gif_gallery_image_video_in_media_content_type(type)
        ) {
            media_content = null;
        }
        // If its not a link type and have metadata
        else if (
            currPost['link_type'] == false &&
            type &&
            gif_gallery_image_video_in_media_content_type(type)
        ) {
            media_content = {
                type,
                original_link: original_link,
                gallery: gallery,
                imageLink: imageLink,
                videoLink: videoLink,
                gifLink: gifLink
            };
        }
        // Else
        else {
            media_content = null;
        }

        try {
            let res = await prisma.posts.create({
                data: {
                    author: {
                        connect: {
                            id: getAuthor?.id,
                            username: getAuthor?.username,
                            email: getAuthor?.email
                        }
                    },
                    voxsphere: {
                        connect: [{ id: getVoxSphere?.id }]
                    },
                    topic: {
                        connect: [{ id: getTopic?.id }]
                    },
                    awards: { connect: [...allAwardsID] },
                    ups: currPost['ups'],
                    text: currPost['text'],
                    textHTML: currPost['text_html'],
                    over18: currPost['over_18'],
                    title: currPost['title'],
                    spoiler: currPost['spoiler'],
                    linkType: currPost['link_type'],
                    createdatHuman: currPost['createdatHuman'],
                    createdAtUnix: currPost['createdat'],
                    numComments: currPost['num_comments'],
                    ...(media_content
                        ? {
                              mediaContent: {
                                  create: media_content
                              }
                          }
                        : {})
                }
            });

            if (resFlair) {
                await prisma.posts.update({
                    where: {
                        id: res.id
                    },
                    data: {
                        postflair: {
                            connect: {
                                id: resFlair.id
                            }
                        }
                    }
                });
            }

            await prisma.posts.update({
                where: {
                    id: res.id
                },
                data: {
                    comments: {
                        connect: await parseComments(
                            currPost['comments'],
                            'isParent',
                            res.id
                        )
                    }
                }
            });
        } catch (err) {
            console.log(err);
            throw err;
        }
    }
    console.log('Done inserting posts');
}

async function update_member_count(voxspheresdata: any, usersdata: any) {
    let subReddits = Object.values(voxspheresdata).flat();
    let allUsers = Object.values(usersdata).flat();

    for (let i = 0; i < subReddits.length; i++) {
        let { id, title } = subReddits[i];
        let mem_count = 0;
        // Go through all the users and check subreddit id present in subreddit meembers
        for (let index = 0; index < allUsers.length; index++) {
            let currUserSUbredditMember = allUsers[index]['subreddits_member'];
            for (let k = 0; k < currUserSUbredditMember.length; k++) {
                const [sub_id, name] = currUserSUbredditMember[k];
                if (id === sub_id) {
                    mem_count++;
                }
            }
        }
        let new_title = title.replace('r/', '');
        try {
            console.log('Updating ...');
            await prisma.voxsphere.update({
                where: {
                    name: new_title
                },
                data: {
                    totalmembers: mem_count,
                    totalmembersHuman: readable(mem_count)
                }
            });
        } catch (err) {
            throw err;
        }
    }
}

async function run_all() {
    // DISABLE THIS [trophies + awards while appending new data]

    // Handle trophies
    let trophiesdata = readFile(
        path.resolve('./json/trophies.json'),
        'trophies'
    );
    if (trophiesdata) {
        await handleTrophies(trophiesdata);
    }

    // Handle awards
    let awardsdata = readFile(path.resolve('./json/awards.json'), 'awards');
    if (awardsdata) {
        await handleAwards(awardsdata);
    }

    // Handle voxsphere
    let voxspheresdata = readFile(
        path.resolve('./json/subreddits.json'),
        'voxspheres'
    );
    if (voxspheresdata) {
        await handleVoxSphere(voxspheresdata);
    }

    // Handle Users
    let usersdata = readFile(path.resolve('./json/users.json'), 'user');
    if (usersdata) {
        await handleUser(usersdata, Object.values(voxspheresdata).flat());
    }

    // Handle Posts
    let postsdata = readFile(path.resolve('./json/posts.json'), 'post');
    if (postsdata) {
        handlePosts(postsdata);
    }

    // Update member count
    await update_member_count(voxspheresdata, usersdata);
    console.log('Done Updating member count');
}

run_all();
