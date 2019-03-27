package product

import (
	"shoppingcart/model"
	"shoppingcart/shared"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"sort"
)

const productTable = "Products"

//Service product service
type Service interface {
	//List returns product list
	List(request *ListRequest) *ListResponse
}

type service struct{}

func (s *service) List(request *ListRequest) *ListResponse {
	response := &ListResponse{Response: shared.NewResponse()}
	var err error
	response.Data, err = s.products()
	sort.Sort(shared.NewProductSorter(response.Data, shared.ByProductID))
	response.SetError(err)
	return response
}

func (s *service) products() ([]*model.Product, error) {
	session, err := shared.GetDynamoService()
	if err != nil {
		return nil, err
	}
	input := &dynamodb.ScanInput{
		AttributesToGet: []*string{
			aws.String("ID"),
			aws.String("Name"),
			aws.String("Price"),
			aws.String("Quantity"),
		},
		TableName: aws.String(productTable),
	}
	output, err := session.Scan(input)
	if err != nil {
		return nil, err
	}
	var products = make([]*model.Product, 0)
	for _, item := range output.Items {
		product := &model.Product{}
		if err = dynamodbattribute.UnmarshalMap(item, &product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

//New creates a new list service
func New() Service {
	return &service{}
}
