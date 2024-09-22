package repositories

import (
	"context"
	"fmt"

	models "github.com/Orken1119/Ozinshe/internal/models/movie_models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) models.MovieRepository {
	return &MovieRepository{db: db}
}

func (mr *MovieRepository) GetMovieProfile(c context.Context, ID int) (models.MovieInfo, error) {
	movie := models.MovieInfo{}

	query := `
	SELECT  
    mi.id,
    mi.name, 
    mi.year_of_issue, 
    mi.director, 
    mi.producer, 
    mi.minutes, 
    mi.hours, 
    STRING_AGG(DISTINCT g.name, ', ') AS genres,
    mi.views,
    mi.description, 
    mi.image
FROM movie_info mi 
JOIN 
    movie_categories mc ON mi.id = mc.movie_id
JOIN 
    categories c ON mc.category_id = c.id
JOIN 
    movie_genres mg ON mi.id = mg.movie_id
JOIN 
    genres g ON mg.genre_id = g.id
WHERE mi.id = $1 
GROUP BY mi.id, mi.name, mi.year_of_issue, mi.director, mi.producer, mi.minutes, mi.hours, mi.views, mi.description, mi.image;
`

	row := mr.db.QueryRow(c, query, ID)
	err := row.Scan(
		&movie.ID,
		&movie.Name,
		&movie.YearOfIssue,
		&movie.Director,
		&movie.Producer,
		&movie.Minutes,
		&movie.Hours,
		&movie.Genre,
		&movie.Views,
		&movie.Description,
		&movie.Image,
	)
	if err != nil {
		return movie, err
	}
	movie.SimilarMovie, err = mr.GetSimilarMovies(c, ID)
	if err != nil {
		return movie, err
	}
	movie.Seasons, movie.Episodes, err = mr.GetEpisodesCount(c, ID)
	if err != nil {
		return movie, err
	}
	movie.Categories, err = mr.GetCategories(c, ID)
	if err != nil {
		return movie, err
	}
	return movie, nil
}

func (mr *MovieRepository) GetSimilarMovies(c context.Context, ID int) ([]models.MovieProfile, error) {
	var similarMovies []models.MovieProfile
	var arrID []int
	query := `WITH MovieGenres AS (
		SELECT 
			genre_id
		FROM 
			movie_genres
		WHERE 
			movie_id = 2
	),
	SimilarMovies AS (
		SELECT 
			mi.id AS movie_id,
			COUNT(DISTINCT mg.genre_id) AS shared_genre_count
		FROM 
			movie_genres mg
		JOIN 
			MovieGenres mg2 ON mg.genre_id = mg2.genre_id
		JOIN 
			movie_info mi ON mg.movie_id = mi.id
		WHERE 
			mi.id <> $1
		GROUP BY 
			mi.id
		HAVING 
			COUNT(DISTINCT mg.genre_id) >= $1
	)
	SELECT 
		mi.id
	FROM 
		SimilarMovies sm
	JOIN 
		movie_info mi ON sm.movie_id = mi.id
	JOIN 
		movie_genres mg ON mi.id = mg.movie_id
	GROUP BY 
		mi.id
	LIMIT 10;
	`

	rows, err := mr.db.Query(c, query, ID)
	if err != nil {
		return similarMovies, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return similarMovies, err
		}
		arrID = append(arrID, id)
	}
	if err := rows.Err(); err != nil {
		return similarMovies, err
	}

	for _, movieID := range arrID {
		var similarMovie models.MovieProfile
		query2 := `
		SELECT 
		mi.id,
		mi.image,
		mi.name,
		STRING_AGG(c.name, ', ') AS categories
	FROM 
		movie_info mi
	JOIN 
		movie_categories mc ON mi.id = mc.movie_id
	JOIN 
		categories c ON mc.category_id = c.id
	WHERE 
		mi.id = $1
	GROUP BY 
		mi.id, mi.image, mi.name;`

		err := mr.db.QueryRow(c, query2, movieID).Scan(
			&similarMovie.ID,
			&similarMovie.Image,
			&similarMovie.Name,
			&similarMovie.Categories,
		)
		if err != nil {
			return similarMovies, err
		}
		similarMovies = append(similarMovies, similarMovie)
	}

	return similarMovies, nil
}

func (mr *MovieRepository) GetSeason(c context.Context, seasonID int) (models.Season, error) {
	var season models.Season

	query1 := `SELECT season_number FROM seasons WHERE id = $1`
	row := mr.db.QueryRow(c, query1, seasonID)
	err := row.Scan(&season.SeasonNumber)
	if err != nil {
		return season, err
	}

	query := `SELECT id, episode_number, video_path FROM episodes WHERE season_id = $1`
	rows, err := mr.db.Query(c, query, seasonID)
	if err != nil {
		return season, err
	}

	defer rows.Close()
	for rows.Next() {
		var episode models.Episode
		err := rows.Scan(&episode.ID, &episode.EpisodeNumber, &episode.Path)
		if err != nil {
			return season, err
		}
		season.Episodes = append(season.Episodes, episode)
	}
	if err := rows.Err(); err != nil {
		return season, err
	}

	return season, nil
}
func (mr *MovieRepository) GetEpisodes(c context.Context, movieID int) ([]models.Season, error) {
	var seasons []models.Season
	var season models.Season
	query := `select id from seasons where movie_id = $1`
	rows, err := mr.db.Query(c, query, movieID)
	if err != nil {
		return seasons, err
	}

	defer rows.Close()
	for rows.Next() {
		var num int
		err = rows.Scan(&num)
		if err != nil {
			return seasons, err
		}
		season, err = mr.GetSeason(c, num)
		seasons = append(seasons, season)
	}
	return seasons, err
}

func (mr *MovieRepository) AddWatchedMovie(c context.Context, userID int, episodeID int) error {

	query := `INSERT INTO watched_movies (user_id, episode_id) VALUES ($1, $2)
	ON CONFLICT (user_id, episode_id) DO NOTHING;`
	_, err := mr.db.Exec(c, query, userID, episodeID)
	if err != nil {
		return err
	}

	var seasonID int
	query3 := `select season_id from episodes where id = $1`
	row1 := mr.db.QueryRow(c, query3, episodeID)
	err = row1.Scan(&seasonID)
	if err != nil {
		return err
	}

	var movieID int
	query4 := `select movie_id from seasons where id = $1`
	row2 := mr.db.QueryRow(c, query4, seasonID)
	err = row2.Scan(&movieID)
	if err != nil {
		return err
	}

	var view int
	query1 := `SELECT views FROM movie_info WHERE id = $1`
	row := mr.db.QueryRow(c, query1, movieID)

	err = row.Scan(&view)
	if err != nil {
		return err
	}

	view++

	query2 := fmt.Sprintf("UPDATE movie_info SET views = %d WHERE id = %d", view, movieID)
	_, err = mr.db.Exec(c, query2)
	if err != nil {
		return err
	}

	return nil
}
func (mr *MovieRepository) GetWatchedMoviesByUser(c context.Context, userID uint) ([]models.Episode, error) {

	var episodes []models.Episode

	query := `SELECT episode_id, e.episode_number, e.video_path 
	FROM watched_movies w 
	join episodes e on w.episode_id = e.id
	where user_id = $1`
	rows, err := mr.db.Query(c, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var episode models.Episode
		err := rows.Scan(&episode.ID, &episode.EpisodeNumber, &episode.Path)
		if err != nil {
			return nil, err
		}
		episodes = append(episodes, episode)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return episodes, nil

}

func (mr *MovieRepository) GetByCategory(c context.Context, category string) ([]models.MovieProfile, error) {
	var movies []models.MovieProfile

	query := `SELECT 
    m.id, 
    m.name, 
    STRING_AGG(DISTINCT c.name, ', ') AS categories, 
    m.image
FROM 
    movie_info m
JOIN 
    movie_categories mc ON m.id = mc.movie_id
JOIN 
    categories c ON mc.category_id = c.id
WHERE 
    m.id IN (
        SELECT 
            m.id
        FROM 
            movie_info m
        JOIN 
            movie_categories mc ON m.id = mc.movie_id
        JOIN 
            categories c ON mc.category_id = c.id
        WHERE 
            c.name = $1
    )
GROUP BY 
    m.id, m.name, m.image;

	`
	rows, err := mr.db.Query(c, query, category)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var movie models.MovieProfile
		err := rows.Scan(&movie.ID, &movie.Name, &movie.Categories, &movie.Image)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return movies, nil
}
func (mr *MovieRepository) FindMovie(c context.Context, name string) ([]models.MovieProfile, error) {
	var movies []models.MovieProfile
	query := `select mi.id, 
	mi.name,
	 STRING_AGG(DISTINCT c.name, ', ') AS categories,
	  mi.image
	   FROM
	    movie_info mi 
		JOIN 
            movie_categories mc ON mi.id = mc.movie_id
        JOIN 
            categories c ON mc.category_id = c.id
		where
		 mi.name like $1 
		 GROUP BY mi.id, mi.name, mi.image;`
	rows, err := mr.db.Query(c, query, name+"%")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var movie models.MovieProfile
		err := rows.Scan(&movie.ID, &movie.Name, &movie.Categories, &movie.Image)
		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
func (mr *MovieRepository) GetTopMovies(c context.Context) ([]models.MovieProfile, error) {
	var topMovies []models.MovieProfile
	query := `
	SELECT mi.id, 
	       mi.name,
		   mi.image,
		   STRING_AGG(DISTINCT c.name, ', ') AS categories
    FROM movie_info mi
	JOIN 
            movie_categories mc ON mi.id = mc.movie_id
    JOIN 
            categories c ON mc.category_id = c.id
	GROUP BY 
	mi.id, mi.name, mi.image
	ORDER BY 
	mi.views
	LIMIT 50 `
	rows, err := mr.db.Query(c, query)
	if err != nil {
		return topMovies, err
	}

	defer rows.Close()
	for rows.Next() {
		var topMovie models.MovieProfile
		err := rows.Scan(&topMovie.ID, &topMovie.Name, &topMovie.Image, &topMovie.Categories)
		if err != nil {
			return topMovies, err
		}
		topMovies = append(topMovies, topMovie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return topMovies, nil
}
func (mr *MovieRepository) AddInFavorites(c context.Context, movieId int, userId int) error {
	query := `insert into favorite_movies (movie_id, user_id)
	values($1,$2)`
	_, err := mr.db.Exec(c, query, movieId, userId)
	if err != nil {
		return err
	}

	return nil
}
func (mr *MovieRepository) GetFavoriteMoviesByUser(c context.Context, userID uint) ([]models.MovieProfile, error) {
	var favoriteMovies []models.MovieProfile
	var movieID []int
	query := `select movie_id
	from
	favorite_movies 
where user_id = $1
		`

	rows, err := mr.db.Query(c, query, userID)
	if err != nil {
		return favoriteMovies, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return favoriteMovies, err
		}
		movieID = append(movieID, id)
	}
	if err := rows.Err(); err != nil {
		return favoriteMovies, err
	}

	for _, each := range movieID {
		var favoriteMovie models.MovieProfile
		query2 := `
		SELECT 
		mi.id,
		mi.image,
		mi.name,
		STRING_AGG(c.name, ', ') AS categories
	FROM 
		movie_info mi
	JOIN 
		movie_categories mc ON mi.id = mc.movie_id
	JOIN 
		categories c ON mc.category_id = c.id
	WHERE 
		mi.id = $1
	GROUP BY 
		mi.id, mi.image, mi.name;`

		err := mr.db.QueryRow(c, query2, each).Scan(
			&favoriteMovie.ID,
			&favoriteMovie.Image,
			&favoriteMovie.Name,
			&favoriteMovie.Categories,
		)
		if err != nil {
			return favoriteMovies, err
		}
		favoriteMovies = append(favoriteMovies, favoriteMovie)
	}

	return favoriteMovies, nil
}
func (mr *MovieRepository) GetEpisodesCount(c context.Context, movieID int) (int, int, error) {
	var seasonsCount int
	var episodesCount int
	query := `SELECT COUNT(*) 
	FROM seasons
	WHERE movie_id = $1;
	`
	row := mr.db.QueryRow(c, query, movieID)

	err := row.Scan(&seasonsCount)
	if err != nil {
		return seasonsCount, episodesCount, err
	}

	query1 := `select count(*)
	from episodes e
	join seasons s on e.season_id = s.id
	where s.movie_id = $1`
	row1 := mr.db.QueryRow(c, query1, movieID)

	err = row1.Scan(&episodesCount)
	if err != nil {
		return seasonsCount, episodesCount, err
	}
	return seasonsCount, episodesCount, nil
}
func (mr *MovieRepository) GetCategories(c context.Context, movieID int) ([]models.Category, error) {
	var category models.Category
	var categories []models.Category
	query := `select c.id,c.name from categories c 
	join movie_categories m on c.id = m.category_id
	where m.movie_id = $1`

	rows, err := mr.db.Query(c, query, movieID)
	if err != nil {
		return categories, err
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&category.ID, &category.Category)
		if err != nil {
			return categories, nil
		}
		categories = append(categories, category)
	}

	return categories, nil

}
