package movie_models

import "context"

type MovieInfo struct {
	ID           uint           `json:"id"`
	Name         string         `json:"name"`
	Categories   []Category     `json:"categories"`
	YearOfIssue  int            `json:"year_of_issue"`
	Director     string         `json:"director"`
	Producer     string         `json:"producer"`
	Minutes      int            `json:"minutes"`
	Hours        int            `json:"hours"`
	Genre        string         `json:"genre"`
	Views        int            `json:"views"`
	Description  string         `json:"description"`
	Image        string         `json:"image"`
	SimilarMovie []MovieProfile `json:"similar_movie"`
	Seasons      int            `json:"seasons"`
	Episodes     int            `json:"episodes"`
}

type Category struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
}

type Movies struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Season    int    `json:"season"`
	Episode   int    `json:"episode"`
	VideoPath string `json:"video_path"`
}

type MovieProfile struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Categories string `json:"categories"`
	Image      string `json:"image"`
}

type Season struct {
	SeasonNumber int       `json:"season_number"`
	Episodes     []Episode `json:"episodes"`
}

type Episode struct {
	ID            int    `json:"id"`
	Path          string `json:"path"`
	EpisodeNumber int    `json:"episode_number"`
}

type WatchedMovie struct {
	ID      uint    `json:"id"`
	UserID  uint    `json:"user_id"`
	MovieID uint    `json:"movie_id"`
	Episode Episode `json:"episodes"`
}

type FoundMovies struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Categories string `json:"categories"`
	Image      string `json:"image"`
}
type SearchRequest struct {
	Name string `form:"name"`
}

type MovieRepository interface {
	GetMovieProfile(c context.Context, ID int) (MovieInfo, error)
	GetSimilarMovies(c context.Context, ID int) ([]MovieProfile, error)
	GetSeason(c context.Context, seasonID int) (Season, error)
	GetEpisodes(c context.Context, movieID int) ([]Season, error)
	GetWatchedMoviesByUser(ctx context.Context, userID uint) ([]Episode, error)
	AddWatchedMovie(c context.Context, userID int, episodeID int) error
	GetByCategory(c context.Context, category string) ([]MovieProfile, error)
	FindMovie(c context.Context, name string) ([]MovieProfile, error)
	GetTopMovies(c context.Context) ([]MovieProfile, error)
	AddInFavorites(c context.Context, movieId int, userId int) error
	GetFavoriteMoviesByUser(c context.Context, userID uint) ([]MovieProfile, error)
	GetEpisodesCount(c context.Context, seasonID int) (int, int, error)
	GetCategories(c context.Context, movieID int) ([]Category, error)
}
