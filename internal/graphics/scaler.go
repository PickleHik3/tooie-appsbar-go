package graphics

import (
	"image"

	"golang.org/x/image/draw"
)

// ScaleImage scales an image to the target dimensions using CatmullRom (Lanczos-like).
func ScaleImage(src image.Image, targetW, targetH int) image.Image {
	if targetW <= 0 || targetH <= 0 {
		return src
	}

	dst := image.NewRGBA(image.Rect(0, 0, targetW, targetH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)
	return dst
}

// ScaleImageAspectFit scales an image to fit within target dimensions while preserving aspect ratio.
func ScaleImageAspectFit(src image.Image, maxW, maxH int) image.Image {
	if maxW <= 0 || maxH <= 0 {
		return src
	}

	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	// Calculate scale factor to fit within bounds
	scaleW := float64(maxW) / float64(srcW)
	scaleH := float64(maxH) / float64(srcH)
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	targetW := int(float64(srcW) * scale)
	targetH := int(float64(srcH) * scale)

	if targetW <= 0 {
		targetW = 1
	}
	if targetH <= 0 {
		targetH = 1
	}

	return ScaleImage(src, targetW, targetH)
}
