from imagekitio import ImageKit
from imagekitio.client import UploadFileRequestOptions
import os
import concurrent.futures
import multiprocessing
import json

# import cloudinary.uploader


imagekit = ImageKit(
    # private_key="",
    # public_key="",
    # url_endpoint="",
    private_key="private_ZEdb1Xkl5F8MTChESue/AMdlC2c=",
    public_key="public_/VwD4BfXGuqXPTQ2HNlq8Si59nE=",
    url_endpoint="https://ik.imagekit.io/gys1ezqjb",
)

# cloudinary.config(
#     cloud_name="",
#     api_key="",
#     api_secret="",
#     secure=True,
# )

master_data = {}


def get_all_files(directory):
    # Create a dictionary to hold the count of each extension
    allImageFiles = []

    for filename in os.listdir(directory):
        _, extension = os.path.splitext(filename)
        if (
            extension == ".jpg"
            or extension == ".png"
            or extension == ".jpeg"
            or extension == ".webp"
        ):
            allImageFiles.append([os.path.abspath(f"./media/{filename}"), "image"])
        elif extension == ".gif":
            allImageFiles.append([os.path.abspath(f"./media/{filename}"), "gif"])
        elif extension == ".mp4":
            allImageFiles.append([os.path.abspath(f"./media/{filename}"), "video"])

    return allImageFiles


def upload_file(path: str, folderName: str):
    filename = os.path.abspath(path).split("/")[-1].split(".")[0]
    ext = os.path.abspath(path).split("/")[-1].split(".")[-1]
    # result = {"secure_url": ""}

    print(f"Uploading {path} ...")

    upload = imagekit.upload(
        file=open(path, "rb"),
        file_name=f"{filename}.{ext}",
        options=UploadFileRequestOptions(
            folder=f"social-media/media-content/{folderName}"
        ),
    )
    master_data[os.path.abspath(path)] = upload.response_metadata.raw["url"]
    return f'Done {upload.response_metadata.raw["url"]} !'

    # if ext == "mp4":
    #     result = cloudinary.uploader.upload_large(
    #         path,
    #         folder=f"social-media/media-content/{folderName}",
    #         resource_type="video",
    #     )
    #     master_data[os.path.abspath(path)] = result["secure_url"]
    #     return f'Done {result["secure_url"]} !'
    # else:
    #     result = cloudinary.uploader.upload(
    #         path,
    #         folder=f"social-media/media-content/{folderName}",
    #     )
    #     master_data[os.path.abspath(path)] = result["secure_url"]
    #     return f'Done {result["secure_url"]} !'


def upload_media():
    with concurrent.futures.ThreadPoolExecutor(
        max_workers=multiprocessing.cpu_count() + 5
    ) as executor:

        args = get_all_files("./media/")
        # print(args)

        futures = [executor.submit(upload_file, *arg) for arg in args]

    # Collecting results from completed tasks
    for future in concurrent.futures.as_completed(futures):
        try:
            result = future.result()
            print(result)
        except Exception as exc:
            print(f"Generated an exception: {exc}")


upload_media()
with open("assets.json", "w") as fp:
    json.dump(master_data, fp)
