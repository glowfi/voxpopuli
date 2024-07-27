import { execSync } from 'child_process';
import fs from 'fs';
import path from 'path';
require('dotenv').config();

// utility functions
function getFileSizeInMB(filePath: string) {
    try {
        const stats = fs.statSync(filePath);
        const fileSizeInBytes = stats.size;
        return fileSizeInBytes / (1024 * 1024);
    } catch (err) {
        console.log('Seeing the file first time ...');
        // console.error(`Error getting file size: ${err}`);
        return 10000;
    }
}

const handleImageResizing_convert_to_webp = (
    filename: string,
    fileExtension: string,
    filePath: string
) => {
    let idx = 0;
    while (true) {
        if (parseFloat(getFileSizeInMB(`media/${filename}.webp`)) < 10) {
            execSync(
                `mogrify -format webp -quality ${100 - idx} "${filePath}"`
            );
            execSync(`rm -rf 'media/${filename}.${fileExtension}'`);
            break;
        }
        execSync(`mogrify -format webp -quality ${100 - idx} "${filePath}"`);
        idx = idx + 5;
    }
};

const handleGifResizing = (
    filename: string,
    fileExtension: string,
    filePath: string
) => {
    let idx = 0;
    while (true) {
        if (parseFloat(getFileSizeInMB(`media/${filename}-lossy.gif`)) < 10) {
            execSync(
                `gifsicle -O3 --lossy=${25 + idx} -o '${filePath}' '${filePath}'`
            );
            execSync(`rm -rf 'media/${filename}-lossy.gif'`);

            break;
        }
        execSync(
            `gifsicle -O3 --lossy=${25 + idx} -o 'media/${filename}-lossy.gif' '${filePath}'`
        );
        idx = idx + 5;
    }
};

// verify if video size under 100 mb and gif/image under 10 mb limit
const checkFileSize = async (
    filePath: string,
    filename: string,
    file_location: string
) => {
    const fileExtension = filePath?.split('.').pop()?.toLowerCase();

    try {
        // Calculate the file size in MB
        console.log(filePath, filename, file_location);
        const fileSizeInMB = getFileSizeInMB(filePath);
        console.log(`File Size: ${fileSizeInMB.toString(6)} MB`);

        // Check if the file size is greater than 100 MB
        if (fileExtension === 'mp4') {
            if (fileSizeInMB > 100) {
                // Log the video resized
                execSync(`touch vid-resize.txt`);
                execSync(
                    `echo '${filePath} Location: ${file_location}' >> vid-resize.txt`
                );
                return ['mp4', false];
            }
            return ['mp4', false];
        } else {
            // Size greater than 10 mb(images+gif)
            if (fileSizeInMB > 10) {
                // Log the gif/image resized
                let cmd = `File Size: ${fileSizeInMB.toFixed(4)} MB Location: ${file_location}`;
                execSync(`touch img-resize.txt`);
                execSync(`echo '${cmd}' >> img-resize.txt`);

                if (fileExtension === 'gif') {
                    //resize gif
                    if (fileSizeInMB < 40) {
                        handleGifResizing(filename, fileExtension, filePath);
                    }
                    return ['gif', true];
                } else {
                    //resize image and convert them to webp
                    handleImageResizing_convert_to_webp(
                        filename,
                        fileExtension,
                        filePath
                    );
                    return ['notwebp', true];
                }
            } else {
                // Size lesser than 10 mb(images+gif+video)

                // Convert images to webp If size lesser than 10 mb
                if (
                    fileExtension === 'png' ||
                    fileExtension === 'jpeg' ||
                    fileExtension === 'jpg'
                ) {
                    // just convert images to webp
                    handleImageResizing_convert_to_webp(
                        filename,
                        fileExtension,
                        filePath
                    );
                    return ['notwebp', false];
                }
                // If size lesser than 10 mb just send the raw gif or video
                console.log('The file size is not greater than 10 MB.');
                return [fileExtension, false];
            }
        }
    } catch (error) {
        console.error('Error accessing the file:', error);
        throw error;
    }
};

// download video
const downloadVideo = async (file_location: string, filename: string) => {
    try {
        execSync(`yt-dlp -o '${filename}' '${file_location}'`);
        console.log('Donwloaded!');
        return false;
    } catch (err) {
        // handle 403 forbidden video
        execSync(`touch err.txt`);
        execSync(`echo '${file_location} ${err}' >> err.txt`);
        return true;
    }
};

// download image
const download_gif_image_video = async (
    file_location: string,
    type: string
) => {
    let ext = file_location.split('/')[3].split('?')[0].split('.')[1];
    let filename = crypto.randomUUID();
    let filename_with_ext = 'media/' + filename + '.' + ext;
    let iserr = false;

    if (
        type === 'video' &&
        JSON.stringify(file_location.split('/')).includes('m3u8')
    ) {
        // handle video
        console.log(file_location);
        iserr = await downloadVideo(
            file_location,
            'media/' + filename + '.mp4'
        );
        filename_with_ext = 'media/' + filename + '.mp4';
    } else {
        // handle image/gif
        try {
            execSync(
                `aria2c -j 16 -x 16 -s 16 -k 1M '${file_location}' -o '${filename_with_ext}'`
            );
        } catch (err) {
            // handle 403 forbidden image/gif
            execSync(`touch err.txt`);
            execSync(`echo '${file_location} ${err}' >> err.txt`);
            iserr = true;
        }
    }

    // if current media is not forbidden 403
    if (!iserr) {
        if (
            ext === 'jpg' ||
            ext === 'jpeg' ||
            ext === 'webp' ||
            ext === 'png'
        ) {
            return ['media/' + filename + '.webp'];
        }
        return [filename_with_ext];
    }
    // if current media is forbidden 403 return NA
    return ['NA'];
};

// Get media links
export const getLink = async (type: string, link: string) => {
    let val = await download_gif_image_video(link, type);

    if (val[0] === 'NA') {
        return ['NA', val[1]];
    } else {
        let absPath = path.resolve(`./${val[0]}`);
        return [absPath, val[1]];
    }
};
