package helpers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/keshav743/cinephile/database"
	"github.com/keshav743/cinephile/models"
)

var GenreURL string = "https://api.themoviedb.org/3/genre/movie/list?api_key=c8af2e4fcd4bf5d99fcb9bfa901fc684&language=en-US"

type GenreResponse struct {
	Genres []models.Genre
}

var wg sync.WaitGroup

func PushGenresToDB() {
	genre := new(GenreResponse)

	resp, err := http.Get(GenreURL)
	HandleError(err)

	err = json.NewDecoder(resp.Body).Decode(&genre)
	HandleError(err)

	for i := 0; i < len(genre.Genres); i++ {
		wg.Add(1)
		go func(gen models.Genre) {
			database.Genres.InsertOne(context.TODO(), gen)
			wg.Done()
		}(genre.Genres[i])
	}

	wg.Wait()
}
