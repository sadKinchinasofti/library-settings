package database

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type LibRepo interface {
	InsertBook(l *Library) error
	GetBook() []GetLibrary
	GetBookByID(id string) []GetLibrary
	DeleteBook(l *GetLibrary) error
	UpdateBook(l *GetLibrary) error
}

type repolibrary struct{}

func getSession() *session.Session {
	config := &aws.Config{
		Region:   aws.String("ap-south-1"),
		Endpoint: aws.String("http://localhost:8000"),
	}

	sess := session.Must(session.NewSession(config))
	return sess
}

func createTable() {
	session := getSession()
	svc := dynamodb.New(session)
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("BookID"),
				AttributeType: aws.String("N"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("BookID"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
		TableName: aws.String("Library"),
	}

	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fmt.Println("Table created successfully")
}

func (re *repolibrary) InsertBook(l *Library) error {
	//createTable()
	session := getSession()
	svc := dynamodb.New(session)
	av, err := dynamodbattribute.MarshalMap(l)
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("Library"),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (re *repolibrary) GetBook() []GetLibrary {
	session := getSession()
	svc := dynamodb.New(session)
	var l []GetLibrary
	// proj := expression.NamesList(expression.Name("BookID"), expression.Name("BookName"))
	// expr, err := expression.NewBuilder().WithProjection(proj).Build()
	// if err != nil {
	// 	fmt.Println("Got error building expression:")
	// 	fmt.Println(err.Error())
	// 	os.Exit(1)
	// }
	// params := &dynamodb.ScanInput{
	// 	ExpressionAttributeNames:  expr.Names(),
	// 	ExpressionAttributeValues: expr.Values(),
	// 	ProjectionExpression:      expr.Projection(),
	// 	TableName:                 aws.String("Library"),
	// }
	params := &dynamodb.ScanInput{
		TableName: aws.String("Library"),
	}
	result, errresult := svc.Scan(params)
	if errresult != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((errresult.Error()))
		os.Exit(1)
	}
	for _, i := range result.Items {
		lib := GetLibrary{}

		err := dynamodbattribute.UnmarshalMap(i, &lib)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		l = append(l, lib)
	}
	return l

}

func (re *repolibrary) GetBookByID(id string) []GetLibrary {
	session := getSession()
	svc := dynamodb.New(session)
	var l1 []GetLibrary
	//bookid := string(l.BookID)
	//bookname := l.BookName
	result, errresult := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("Library"),
		Key: map[string]*dynamodb.AttributeValue{
			"BookID": {
				N: aws.String(id),
			},
		},
	})
	if errresult != nil {
		panic(errresult)
	}
	lib := GetLibrary{}
	err := dynamodbattribute.UnmarshalMap(result.Item, &lib)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	l1 = append(l1, lib)
	return l1
}

func (re *repolibrary) DeleteBook(l *GetLibrary) error {
	session := getSession()
	svc := dynamodb.New(session)
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String("Library"),
		Key: map[string]*dynamodb.AttributeValue{
			"BookID": {
				N: aws.String(l.BookID),
			},
		},
	}

	_, errd := svc.DeleteItem(input)
	if errd != nil {
		fmt.Println("Got error calling DeleteItem")
		fmt.Println(errd.Error())
		return errd
	}
	return nil

}

func (re *repolibrary) UpdateBook(l *GetLibrary) error {
	session := getSession()
	svc := dynamodb.New(session)

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":BookAuthor": {
				S: aws.String(l.BookAuthor),
			},
		},
		TableName: aws.String("Library"),
		Key: map[string]*dynamodb.AttributeValue{
			"BookID": {
				N: aws.String(l.BookID),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set BookAuthor = :BookAuthor"),
	}

	_, erru := svc.UpdateItem(input)
	if erru != nil {
		fmt.Println(erru.Error())
		return erru
	}

	return nil

}
func NewLibRepo() LibRepo {
	return &repolibrary{}
}
