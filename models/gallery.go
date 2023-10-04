package models

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gastrader/website/errors"
)

type Image struct {
	Path string
	GalleryID int
	Filename string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	//where to store/locate images
	ImagesDir string

}

func (gs *GalleryService) Create(title string, userID int) (*Gallery, error){
	gallery := Gallery{
		Title: title,
		UserID: userID,
	}
	row := gs.DB.QueryRow(`
		INSERT INTO galleries (title, user_id)
		VALUES ($1, $2) RETURNING id;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err != nil{
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) ByID(id int) (*Gallery, error){
	gallery := Gallery{
		ID: id,
	}
	row := gs.DB.QueryRow(`
		SELECT title, user_id
		FROM galleries
		WHERE id = $1;`, gallery.ID)
	err := row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil	
}

func (gs *GalleryService) ByUserID(userID int) ([]Gallery, error){
	rows, err := gs.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1;`, userID)
	if err != nil{
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next(){
		gallery := Gallery{
			UserID: userID,
		}
		err = rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil{
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil{
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error{
	_, err := gs.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;`, gallery.ID, gallery.Title)
	if err != nil{
		return fmt.Errorf("update gallery: %w", err)
	}
	return nil
}

func (gs *GalleryService) Delete(id int) error{
	_, err := gs.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("delete gallery.: %w", err)
	}
	return nil
}

func (gs *GalleryService) galleryDir(id int) string{
	imagesDir := gs.ImagesDir
	if imagesDir == ""{
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%d", id))
}

func (gs *GalleryService) Images(galleryID int) ([]Image, error){
	globPattern := filepath.Join(gs.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil{
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []Image
	for _, file := range allFiles{
		if hasExtension(file, gs.extensions()){
			images = append(images, Image{
				Path: file,
				Filename: filepath.Base(file),
				GalleryID: galleryID,
			})
		}
	}
	return images, nil
}

func hasExtension(file string, extensions []string) bool{
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

func (gs *GalleryService) Image(galleryID int, filename string) (Image, error){
	imagePath := filepath.Join(gs.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil{
		if errors.Is(err, fs.ErrNotExist){
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("querying for image: %w", err)
	}
	return Image{
		Filename: filename,
		GalleryID: galleryID,
		Path: imagePath,
	}, nil
}

func (gs *GalleryService)extensions() []string{
	return []string{".png", ".jpg", ".jpeg", ".gif", ".webp"}
}