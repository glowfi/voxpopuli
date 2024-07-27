import os
import concurrent.futures
import multiprocessing


###### Image to webp ######
def get_all_files(directory):
    # Create a dictionary to hold the count of each extension
    allImageFiles = []

    # Iterate over all files in the specified directory
    for filename in os.listdir(directory):
        # Get the file extension (without the leading dot)
        _, extension = os.path.splitext(filename)

        # If the extension is already in the dictionary, increment its count
        if extension == ".jpg" or extension == ".jpeg" or extension == ".png":
            allImageFiles.append([os.path.abspath(f"./media/{filename}")])

    return allImageFiles


def convertImagetoWebp(filePath):
    print(f"Processing {filePath} ...")
    os.system(f"mogrify -format webp -quality 90 {filePath}")
    os.system(f"rm -rf {filePath}")
    return f"Done Converting {filePath}! -> webp"


def image_process():
    # Get all Images
    directory_path = "./media/"
    allImageFiles = get_all_files(directory_path)
    args = []
    for image_location in allImageFiles:
        args.append(image_location)

    # Submitting tasks to the executor
    with concurrent.futures.ThreadPoolExecutor(
        max_workers=multiprocessing.cpu_count() * 5
    ) as executor:
        futures = [executor.submit(convertImagetoWebp, *arg) for arg in args]

    # Collecting results from completed tasks
    for future in concurrent.futures.as_completed(futures):
        try:
            result = future.result()
            print(result)
        except Exception as exc:
            print(f"Generated an exception: {exc}")


###### Convert Large gif to small size ######
def optimize_gif(path):
    print(f"Processing gif {path}")
    os.system(
        f'gifsicle --unoptimize {path} | gifsicle --dither --colors 18 --resize-fit-width 512 -O2 `seq -f "#%g" 0 2 213` -o {path}'
    )
    return f"Done {path}"


def find_large_gifs(directory, args):
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(".gif"):
                file_path = os.path.join(root, file)
                file_size = os.path.getsize(file_path)
                if file_size > 10 * 1024 * 1024:  # 10 MB
                    args.append([file_path])


def gif_process():
    with concurrent.futures.ThreadPoolExecutor(
        max_workers=multiprocessing.cpu_count() * 5
    ) as executor:

        args = []
        find_large_gifs("./media/", args)

        futures = [executor.submit(optimize_gif, *arg) for arg in args]

    # Collecting results from completed tasks
    for future in concurrent.futures.as_completed(futures):
        try:
            result = future.result()
            print(result)
        except Exception as exc:
            print(f"Generated an exception: {exc}")


# Process images and gifs
image_process()
gif_process()
