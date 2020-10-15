package services

import (
	"github.com/XMatrixStudio/Coffee/App/models"
	"github.com/kataras/iris/v12"
)

// 添加电影
func (s *contentService) AddMovie(ctx iris.Context, id string, data ContentData) error {

	return nil
}

func (s *contentService) GetMovieByUser(ownID string, public bool) []models.Content {
	return s.Model.GetContentByOwnAndType(ownID, models.TypeFilm, public)
}
