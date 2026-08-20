package main

type ResolutionTier struct {
	Height         int
	MinBitrateKbps int
}

var resolutionLadder = []ResolutionTier{
	{2160, 35000}, // 4K
	{1440, 16000}, // 2K
	{1080, 8000},  // 1080p
	{720, 5000},   // 720p
	{480, 2500},   // 480p
	{360, 1000},   // 360p
}

func pickResolution(originalHeight int, targetBitrateKbps int) (height int, ok bool) {
	var last ResolutionTier
	for _, tier := range resolutionLadder {
		if tier.Height > originalHeight {
			continue // never upscale
		}
		last = tier
		if targetBitrateKbps >= tier.MinBitrateKbps {
			return tier.Height, true
		}
	}
	if last.Height == 0 {
		// original is smaller than every tier in the ladder
		return originalHeight, true
	}
	return last.Height, false
}
