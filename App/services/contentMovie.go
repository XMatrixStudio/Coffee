package services

import (
	"github.com/kataras/iris"
	"github.com/XMatrixStudio/Coffee/App/models"
)

// 添加电影
func (s *contentService)AddMovie(ctx iris.Context, id string, data ContentData) error {

	return nil
}

func (s *contentService)GetMovieByUser(ownID string, public bool) []models.Content {
	return s.Model.GetContentByOwnAndType(ownID, models.TypeFilm, public)
}