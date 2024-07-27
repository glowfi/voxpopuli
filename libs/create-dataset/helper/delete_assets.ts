import { v2 as cloudinary } from 'cloudinary';

cloudinary.config({
    cloud_name: '',
    api_key: '',
    api_secret: ''
});

cloudinary.api.delete_resources_by_prefix(
    'social-media/media-content/image',
    console.log
);
