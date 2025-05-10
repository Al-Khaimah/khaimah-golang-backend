package categories

type Category struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	IsNewsIntensive bool   `json:"is_news_intensive"`
}
