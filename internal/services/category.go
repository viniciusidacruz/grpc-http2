package services

import (
	"context"
	"io"

	"github.com/viniciusidacruz/grpc-http2/internal/database"
	"github.com/viniciusidacruz/grpc-http2/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)

	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) GetCategories(ctx context.Context, in *pb.Blank) (*pb.Categories, error) {
	categories, err := c.CategoryDB.FindAll()

	if err != nil {
		return nil, err
	}

	categoriesResponse := &pb.Categories{}

	for _, category := range categories {
		categoriesResponse.Categories = append(categoriesResponse.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return categoriesResponse, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.GetCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.FindByCategoryId(in.Id)

	if err != nil {
		return nil, err
	}

	categoryResponse := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryResponse, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.Categories{}

	for {
		request, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(categories)
		}

		if err != nil {
			return err
		}

		category, err := c.CategoryDB.Create(request.Name, request.Description)

		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
}

func (c *CategoryService) CreateCategoryStreamBirectional(stream pb.CategoryService_CreateCategoryStreamBirectionalServer) error {
	for {
		category, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)

		if err != nil {
			return err
		}

		categoryResponse := &pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		}

		err = stream.Send(categoryResponse)

		if err != nil {
			return err
		}
	}
}
