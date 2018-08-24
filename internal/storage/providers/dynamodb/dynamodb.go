package dynamodb

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
	"github.com/crit/inventory/internal/errors"
	"github.com/crit/inventory/internal/storage"
)

var _ storage.Storage = &db{} // compile check

// todo: replace with session injection
func New(region, table string) (storage.Storage, error) {
	if table == "" {
		return nil, errors.String(500, "DynamoDB table name empty")
	}

	local, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		return nil, toError(err)
	}

	return &db{
		svc:   dynamodb.New(local),
		table: aws.String(table),
	}, nil
}

type db struct {
	table *string
	svc   *dynamodb.DynamoDB
}

type item struct {
	ID   string
	Type string
}

type itemData struct {
	ID   string
	Type string
	Data []byte
}

func (db *db) Write(model storage.Writer) error {
	value, err := dynamodbattribute.MarshalMap(itemData{
		ID:   model.Key(),
		Type: model.Type(),
		Data: model.Data(),
	})

	if err != nil {
		return toError(err)
	}

	_, err = db.svc.PutItem(&dynamodb.PutItemInput{
		TableName: db.table,
		Item:      value,
	})

	if err != nil {
		return toError(err)
	}

	return nil
}

func (db *db) Read(model storage.Reader) error {
	value, err := dynamodbattribute.MarshalMap(item{
		ID:   model.Key(),
		Type: model.Type(),
	})

	if err != nil {
		return toError(err)
	}

	out, err := db.svc.GetItem(&dynamodb.GetItemInput{
		TableName: db.table,
		Key:       value,
	})

	if err != nil {
		return toError(err)
	}

	if out.Item == nil {
		return errors.String(http.StatusNotFound, "not found")
	}

	model.SetData(out.Item["Data"].B)

	return nil
}

func (db *db) Delete(model storage.Meta) error {
	value, err := dynamodbattribute.MarshalMap(item{
		ID:   model.Key(),
		Type: model.Type(),
	})

	if err != nil {
		return toError(err)
	}

	_, err = db.svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: db.table,
		Key:       value,
	})

	if err != nil {
		return toError(err)
	}

	return nil
}

func (db *db) List(model storage.Mapper) error {
	filter := expression.Name("Type").Equal(expression.Value(model.Type()))
	projection := expression.NamesList(expression.Name("Key"), expression.Name("Data"))
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()

	if err != nil {
		return toError(err)
	}

	input := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 db.table,
	}

	var pagedErr error

	err = db.svc.ScanPages(input, func(out *dynamodb.ScanOutput, lastPage bool) bool {
		if lastPage {
			model.SetCount(*out.Count)
		}

		for _, item := range out.Items {
			var local itemData
			err := dynamodbattribute.UnmarshalMap(item, &local)

			if err != nil {
				pagedErr = err
				return false // do not continue
			}

			model.Append(local.ID, model.Type(), local.Data)
		}

		return true // continue
	})

	if err != nil {
		return toError(err)
	}

	if pagedErr != nil {
		return toError(err)
	}

	return nil
}
