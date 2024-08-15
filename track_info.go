package main

type TrackInfo struct {
	Result struct {
		Track struct {
			Title      string `json:"title"`
			DurationMs int64  `json:"durationMs"`
			Artists    []struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`
			Albums []struct {
				Id    int    `json:"id"`
				Title string `json:"title"`
				Year  int    `json:"year"`
				Genre string `json:"genre"`
			} `json:"albums"`
		} `json:"track"`
	} `json:"result"`
}
