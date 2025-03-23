package media

import (
	"slices"
	"testing"

	"github.com/glowfi/voxpopuli/backend/pkg/models"
	"github.com/stretchr/testify/assert"
)

func AssertImageMetadataWithoutTimestamp(t *testing.T, wantImageMetadata, gotImageMetadata models.ImageMetadata) {
	assert.Equal(t, wantImageMetadata.ID, gotImageMetadata.ID, "expected id to match")
	assert.Equal(t, wantImageMetadata.ImageID, gotImageMetadata.ImageID, "expected image id to match")
	assert.Equal(t, wantImageMetadata.Height, gotImageMetadata.Height, "expected height to match")
	assert.Equal(t, wantImageMetadata.Width, gotImageMetadata.Width, "expected width to match")
	assert.Equal(t, wantImageMetadata.Url, gotImageMetadata.Url, "expected url to match")
}

func AssertImageMetadatasWithoutTimestamp(t *testing.T, wantImageMetadatas, gotImageMetadatas []models.ImageMetadata) {
	t.Helper()

	if len(wantImageMetadatas) != len(gotImageMetadatas) {
		t.Fatal("length of wantImageMetadatas and gotImageMetadatas do not match")
	}

	for _, imageMetadata := range wantImageMetadatas {
		idx := slices.IndexFunc(gotImageMetadatas, func(im models.ImageMetadata) bool {
			return im.ID == imageMetadata.ID
		})

		if idx == -1 {
			t.Fatalf("image metadata %v of ID %v is not present in gotImageMetadatas", imageMetadata.Url, imageMetadata.ID)
			return
		}
		AssertImageMetadataWithoutTimestamp(t, imageMetadata, gotImageMetadatas[idx])
	}
}

func AssertImageMetadatasWithTimestamp(t *testing.T, wantImageMetadatas, gotImageMetadatas []models.ImageMetadata) {
	t.Helper()

	for _, imageMetadata := range wantImageMetadatas {
		idx := slices.IndexFunc(gotImageMetadatas, func(im models.ImageMetadata) bool {
			return im.ID == imageMetadata.ID
		})

		if idx == -1 {
			t.Fatalf("image metadata %v of ID %v is not present in gotImageMetadatas", imageMetadata.Url, imageMetadata.ID)
			return
		}
		assert.Equal(t, imageMetadata, gotImageMetadatas[idx], "expected image metadata to match")
	}
}

func AssertImage(t *testing.T, wantImage, gotImage models.Image) {
	assert.Equal(t, wantImage.ID, gotImage.ID, "expected id to match")
	assert.Equal(t, wantImage.MediaID, gotImage.MediaID, "expected media id to match")
	AssertImageMetadatasWithTimestamp(t, wantImage.ImageMetadata, gotImage.ImageMetadata)
}

func AssertImages(t *testing.T, wantImages, gotImages []models.Image) {
	t.Helper()

	if len(wantImages) != len(gotImages) {
		t.Fatal("length of wantImages and gotImages do not match")
	}

	for _, image := range wantImages {
		idx := slices.IndexFunc(gotImages, func(im models.Image) bool {
			return im.ID == image.ID
		})

		if idx == -1 {
			t.Fatalf("image %v of ID %v is not present in gotImages", image.MediaID, image.ID)
			return
		}
		AssertImage(t, image, gotImages[idx])
	}
}

func AssertGifMetadataWithoutTimestamp(t *testing.T, wantGifMetadata, gotGifMetadata models.GifMetadata) {
	assert.Equal(t, wantGifMetadata.ID, gotGifMetadata.ID, "expected id to match")
	assert.Equal(t, wantGifMetadata.GifID, gotGifMetadata.GifID, "expected gif id to match")
	assert.Equal(t, wantGifMetadata.Height, gotGifMetadata.Height, "expected height to match")
	assert.Equal(t, wantGifMetadata.Width, gotGifMetadata.Width, "expected width to match")
	assert.Equal(t, wantGifMetadata.Url, gotGifMetadata.Url, "expected url to match")
}

func AssertGifMetadatasWithoutTimestamp(t *testing.T, wantGifMetadatas, gotGifMetadatas []models.GifMetadata) {
	t.Helper()

	if len(wantGifMetadatas) != len(gotGifMetadatas) {
		t.Fatal("length of wantGifMetadatas and gotGifMetadatas do not match")
	}

	for _, gifMetadata := range wantGifMetadatas {
		idx := slices.IndexFunc(gotGifMetadatas, func(gm models.GifMetadata) bool {
			return gm.ID == gifMetadata.ID
		})

		if idx == -1 {
			t.Fatalf("gif metadata %v of ID %v is not present in gotGifMetadatas", gifMetadata.Url, gifMetadata.ID)
			return
		}
		AssertGifMetadataWithoutTimestamp(t, gifMetadata, gotGifMetadatas[idx])
	}
}

func AssertGifMetadatasWithTimestamp(t *testing.T, wantGifMetadatas, gotGifMetadatas []models.GifMetadata) {
	t.Helper()

	for _, gifMetadata := range wantGifMetadatas {
		idx := slices.IndexFunc(gotGifMetadatas, func(gm models.GifMetadata) bool {
			return gm.ID == gifMetadata.ID
		})

		if idx == -1 {
			t.Fatalf("gif metadata %v of ID %v is not present in gotGifMetadatas", gifMetadata.Url, gifMetadata.ID)
			return
		}
		assert.Equal(t, gifMetadata, gotGifMetadatas[idx], "expected gif metadata to match")
	}
}

func AssertGif(t *testing.T, wantGif, gotGif models.Gif) {
	assert.Equal(t, wantGif.ID, gotGif.ID, "expected id to match")
	assert.Equal(t, wantGif.MediaID, gotGif.MediaID, "expected media id to match")
	AssertGifMetadatasWithTimestamp(t, wantGif.GifMetadata, gotGif.GifMetadata)
}

func AssertGifs(t *testing.T, wantGifs, gotGifs []models.Gif) {
	t.Helper()

	if len(wantGifs) != len(gotGifs) {
		t.Fatal("length of wantGifs and gotGifs do not match")
	}

	for _, gif := range wantGifs {
		idx := slices.IndexFunc(gotGifs, func(g models.Gif) bool {
			return g.ID == gif.ID
		})

		if idx == -1 {
			t.Fatalf("gif %v of ID %v is not present in gotGifs", gif.MediaID, gif.ID)
			return
		}
		AssertGif(t, gif, gotGifs[idx])
	}
}

func AssertLinkWithoutTimestamp(t *testing.T, wantLink, gotLink models.Link) {
	assert.Equal(t, wantLink.ID, gotLink.ID, "expect link id to match")
	assert.Equal(t, wantLink.MediaID, gotLink.MediaID, "expect link media id to match")
	assert.Equal(t, wantLink.Link, gotLink.Link, "expect link url to match")
}

func AssertLinksWithoutTimestamp(t *testing.T, wantLinks, gotLinks []models.Link) {
	t.Helper()

	if len(wantLinks) != len(gotLinks) {
		t.Fatal("length of wantLinks and gotLinks do not match")
	}

	for _, link := range wantLinks {
		idx := slices.IndexFunc(gotLinks, func(im models.Link) bool {
			return im.ID == link.ID
		})

		if idx == -1 {
			t.Fatalf("link %v of url %v is not present in gotLinks", link.ID, link.Link)
			return
		}
		AssertLinkWithoutTimestamp(t, link, gotLinks[idx])
	}
}

func AssertLinksWitTimestamp(t *testing.T, wantLinks, gotLinks []models.Link) {
	t.Helper()

	if len(wantLinks) != len(gotLinks) {
		t.Fatal("length of wantLinks and gotLinks do not match")
	}

	for _, link := range wantLinks {
		idx := slices.IndexFunc(gotLinks, func(im models.Link) bool {
			return im.ID == link.ID
		})

		if idx == -1 {
			t.Fatalf("link %v of url %v is not present in gotLinks", link.ID, link.Link)
			return
		}
		assert.Equal(t, link, gotLinks[idx], "expect link to match")
	}
}

func AssertVideoWithoutTimestamp(t *testing.T, wantVideo, gotVideo models.Video) {
	assert.Equal(t, wantVideo.ID, gotVideo.ID, "expect video id to match")
	assert.Equal(t, wantVideo.MediaID, gotVideo.MediaID, "expect video media id to match")
	assert.Equal(t, wantVideo.Url, gotVideo.Url, "expect video url to match")
	assert.Equal(t, wantVideo.Height, gotVideo.Height, "expect video height to match")
	assert.Equal(t, wantVideo.Width, gotVideo.Width, "expect video width to match")
}

func AssertVideosWithoutTimestamp(t *testing.T, wantVideos, gotVideos []models.Video) {
	t.Helper()

	if len(wantVideos) != len(gotVideos) {
		t.Fatal("length of wantVideos and gotVideos do not match")
	}

	for _, video := range wantVideos {
		idx := slices.IndexFunc(gotVideos, func(im models.Video) bool {
			return im.ID == video.ID
		})

		if idx == -1 {
			t.Fatalf("video %v of url %v is not present in gotVideos", video.ID, video.Url)
			return
		}
		AssertVideoWithoutTimestamp(t, video, gotVideos[idx])
	}
}

func AssertVideosWithTimestamp(t *testing.T, wantVideos, gotVideos []models.Video) {
	t.Helper()

	if len(wantVideos) != len(gotVideos) {
		t.Fatal("length of wantVideos and gotVideos do not match")
	}

	for _, video := range wantVideos {
		idx := slices.IndexFunc(gotVideos, func(im models.Video) bool {
			return im.ID == video.ID
		})

		if idx == -1 {
			t.Fatalf("video %v of url %v is not present in gotVideos", video.ID, video.Url)
			return
		}
		assert.Equal(t, video, gotVideos[idx], "expect video to match")
	}
}

func AssertGalleryMetadataWithoutTimestamp(t *testing.T, wantGalleryMetadata, gotGalleryMetadata models.GalleryMetadata) {
	assert.Equal(t, wantGalleryMetadata.ID, gotGalleryMetadata.ID, "expected id to match")
	assert.Equal(t, wantGalleryMetadata.GalleryID, gotGalleryMetadata.GalleryID, "expected gallery id to match")
	assert.Equal(t, wantGalleryMetadata.OrderIndex, gotGalleryMetadata.OrderIndex, "expected order index to match")
	assert.Equal(t, wantGalleryMetadata.Height, gotGalleryMetadata.Height, "expected height to match")
	assert.Equal(t, wantGalleryMetadata.Width, gotGalleryMetadata.Width, "expected width to match")
	assert.Equal(t, wantGalleryMetadata.Url, gotGalleryMetadata.Url, "expected url to match")
}

func AssertGalleryMetadatasWithoutTimestamp(t *testing.T, wantGalleryMetadatas, gotGalleryMetadatas []models.GalleryMetadata) {
	t.Helper()

	if len(wantGalleryMetadatas) != len(gotGalleryMetadatas) {
		t.Fatal("length of wantGalleryMetadatas and gotGalleryMetadatas do not match")
	}

	for _, galleryMetadata := range wantGalleryMetadatas {
		idx := slices.IndexFunc(gotGalleryMetadatas, func(gm models.GalleryMetadata) bool {
			return gm.ID == galleryMetadata.ID
		})

		if idx == -1 {
			t.Fatalf("gallery metadata %v of ID %v is not present in gotGalleryMetadatas", galleryMetadata.Url, galleryMetadata.ID)
			return
		}
		AssertGalleryMetadataWithoutTimestamp(t, galleryMetadata, gotGalleryMetadatas[idx])
	}
}

func AssertGalleryMetadatasWithTimestamp(t *testing.T, wantGalleryMetadatas, gotGalleryMetadatas []models.GalleryMetadata) {
	t.Helper()

	for _, galleryMetadata := range wantGalleryMetadatas {
		idx := slices.IndexFunc(gotGalleryMetadatas, func(gm models.GalleryMetadata) bool {
			return gm.ID == galleryMetadata.ID
		})

		if idx == -1 {
			t.Fatalf("gallery metadata %v of ID %v is not present in gotGalleryMetadatas", galleryMetadata.Url, galleryMetadata.ID)
			return
		}
		assert.Equal(t, galleryMetadata, gotGalleryMetadatas[idx], "expected gallery metadata to match")
	}
}

func AssertGallery(t *testing.T, wantGallery, gotGallery models.Gallery) {
	assert.Equal(t, wantGallery.ID, gotGallery.ID, "expected id to match")
	assert.Equal(t, wantGallery.MediaID, gotGallery.MediaID, "expected media id to match")
	AssertGalleryMetadatasWithTimestamp(t, wantGallery.GalleryMetadata, gotGallery.GalleryMetadata)
}

func AssertGalleries(t *testing.T, wantGalleries, gotGalleries []models.Gallery) {
	t.Helper()

	if len(wantGalleries) != len(gotGalleries) {
		t.Fatal("length of wantGalleries and gotGalleries do not match")
	}

	for _, gallery := range wantGalleries {
		idx := slices.IndexFunc(gotGalleries, func(g models.Gallery) bool {
			return g.ID == gallery.ID
		})

		if idx == -1 {
			t.Fatalf("gallery %v of ID %v is not present in gotGalleries", gallery.MediaID, gallery.ID)
			return
		}
		AssertGallery(t, gallery, gotGalleries[idx])
	}
}
